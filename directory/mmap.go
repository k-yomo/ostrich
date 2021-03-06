package directory

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/edsrzf/mmap-go"
)

var _ Directory = &mmapDirectory{}

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

func (m *mmapDirectory) OpenRead(path string) (ReaderCloser, error) {
	f, err := os.Open(m.buildPath(path))
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	mem, err := mmap.Map(f, mmap.RDONLY, 0)
	if err != nil {
		return nil, fmt.Errorf("mmap file: %w", err)
	}
	return newMmapIO(mem), nil
}

func (m *mmapDirectory) Read(path string) (io.ReadCloser, error) {
	f, err := os.Open(m.buildPath(path))
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	mem, err := mmap.Map(f, mmap.RDWR, 0)
	if err != nil {
		return nil, fmt.Errorf("mmap file: %w", err)
	}
	return newMmapIO(mem), nil
}
func (m *mmapDirectory) AtomicRead(path string) ([]byte, error) {
	f, err := os.Open(m.buildPath(path))
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	bytes, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (m *mmapDirectory) OpenWrite(path string) (WriteCloseSyncer, error) {
	f, err := os.OpenFile(m.buildPath(path), os.O_RDWR|os.O_CREATE, 0766)
	if err != nil {
		return nil, err
	}
	return f, err
}

func (m *mmapDirectory) AtomicWrite(path string, data []byte) error {
	f, err := os.Create(m.buildPath(path))
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
	*bytes.Reader
	mmap mmap.MMap
}

func newMmapIO(m mmap.MMap) *mmapIO {
	return &mmapIO{
		Reader: bytes.NewReader(m),
		mmap:   m,
	}
}

func (m *mmapIO) Close() error {
	return m.mmap.Unmap()
}
