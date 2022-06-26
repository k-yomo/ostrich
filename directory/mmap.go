package directory

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/edsrzf/mmap-go"
)

type mmapDirectory struct {
	rootPath string
}

func NewMMapDirectory(rootPath string) (*mmapDirectory, error) {
	fi, err := os.Stat(rootPath)
	if err != nil {
		return nil, err
	}
	if !fi.IsDir() {
		return nil, fmt.Errorf("'%s' is not directory", rootPath)
	}
	return &mmapDirectory{
		rootPath: rootPath,
	}, nil
}

func (m *mmapDirectory) Read(path string) (io.ReadCloser, error) {
	f, err := os.Open(m.buildPath(path))
	if err != nil {
		return nil, fmt.Errorf("open file: %v", err)
	}
	mem, err := mmap.Map(f, mmap.RDWR, 0)
	if err != nil {
		return nil, fmt.Errorf("mmap file: %v", err)
	}
	return newMmapIO(mem, m.buildPath(path)), nil
}
func (m *mmapDirectory) AtomicRead(path string) ([]byte, error) {
	f, err := os.Open(m.buildPath(path))
	if err != nil {
		return nil, fmt.Errorf("open file: %v", err)
	}
	defer f.Close()

	bytes, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (m *mmapDirectory) OpenWrite(path string) (WriteCloseFlasher, error) {
	f, err := os.OpenFile(m.buildPath(path), os.O_RDWR|os.O_CREATE, 0766)
	if err != nil {
		return nil, err
	}
	// Empty file causes mmap error
	if stat, err := f.Stat(); err != nil {
		return nil, err
	} else {
		if stat.Size() == 0 {
			if _, err := f.WriteString("\n"); err != nil {
				return nil, err
			}
		}
	}
	mem, err := mmap.Map(f, mmap.RDWR, 0)
	if err != nil {
		return nil, err
	}

	return newMmapIO(mem, m.buildPath(path)), nil
}

func (m *mmapDirectory) AtomicWrite(path string, data []byte) error {
	f, err := os.OpenFile(m.buildPath(path), os.O_RDWR|os.O_CREATE, 0766)
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	return err
}

func (m *mmapDirectory) Exists(path string) (bool, error) {
	_, err := os.Stat(m.buildPath(path))
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (m *mmapDirectory) buildPath(path string) string {
	return fmt.Sprintf("%s/%s", m.rootPath, path)
}

type mmapIO struct {
	mmap mmap.MMap
	path string
}

func newMmapIO(m mmap.MMap, path string) *mmapIO {
	return &mmapIO{
		mmap: m,
		path: path,
	}
}

func (m *mmapIO) Read(p []byte) (n int, err error) {
	return bytes.NewReader(m.mmap).Read(p)
}

func (m *mmapIO) Write(p []byte) (n int, err error) {
	f, err := os.OpenFile(m.path, os.O_RDWR|os.O_CREATE, 0660)
	if err != nil {
		return 0, err
	}
	return f.Write(p)
}

func (m *mmapIO) Flush() error {
	return m.mmap.Flush()
}

func (m *mmapIO) Close() error {
	return m.mmap.Unmap()
}
