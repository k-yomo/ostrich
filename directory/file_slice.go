package directory

import (
	"fmt"
	"io"
)

type FileHandle interface {
	io.ReaderAt
	Len() int
}

type FileSlice struct {
	data   io.ReaderAt
	siz    int
	closer func() error
}

func NewFileSlice(data FileHandle, closer func() error) *FileSlice {
	return &FileSlice{
		data:   data,
		siz:    data.Len(),
		closer: closer,
	}
}

func (f *FileSlice) Slice(from, to int) *FileSlice {
	return &FileSlice{
		data: io.NewSectionReader(f.data, int64(from), int64(to-from)),
		siz:  to - from,
		closer: func() error {
			panic("Close is called for child file slice, only root file slice can be closed")
		},
	}
}

func (f *FileSlice) Read(from, to int) ([]byte, error) {
	buf := make([]byte, to-from)
	if _, err := f.data.ReadAt(buf, int64(from)); err != nil {
		return nil, fmt.Errorf("read data: %w", err)
	}
	return buf, nil
}

func (f *FileSlice) Len() int {
	return f.siz
}

func (f *FileSlice) Reader() *DataReader {
	return &DataReader{
		data: f,
	}
}

// Close closes data file
// Close must be called only for the top level file slice, and must not for the child sliced one.
func (f *FileSlice) Close() error {
	return f.closer()
}

type DataReader struct {
	data *FileSlice
	n    int
}

func (d *DataReader) Read(p []byte) (n int, err error) {
	if d.n >= d.data.Len() {
		return 0, io.EOF
	}
	from := d.n
	to := d.n + len(p)
	if to > d.data.Len() {
		to = d.data.Len()
	}
	data, err := d.data.Read(from, to)
	if err != nil {
		return 0, err
	}
	copy(p, data)
	d.n = to
	return to - from, nil
}
