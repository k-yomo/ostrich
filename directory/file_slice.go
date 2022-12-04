package directory

import (
	"fmt"
	"io"

	"github.com/k-yomo/ostrich/pkg/xrange"
)

type Reader interface {
	io.Reader
	io.ReaderAt
}

type FileHandle interface {
	Reader
	Len() int
}

type FileSlice struct {
	data   Reader
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

func (f *FileSlice) Slice(r xrange.Range) *FileSlice {
	return &FileSlice{
		data: io.NewSectionReader(f.data, int64(r.From), int64(r.Len())),
		siz:  r.Len(),
		closer: func() error {
			panic("Close is called for child file slice, only root file slice can be closed")
		},
	}
}

func (f *FileSlice) Read(from, size int) ([]byte, error) {
	buf := make([]byte, size)
	if _, err := f.data.ReadAt(buf, int64(from)); err != nil {
		return nil, fmt.Errorf("read data: %w", err)
	}
	return buf, nil
}

func (f *FileSlice) Len() int {
	return f.siz
}

func (f *FileSlice) Reader() io.Reader {
	return f.data
}

// Close closes data file
// Close must be called only for the top level file slice, and must not for the child sliced one.
func (f *FileSlice) Close() error {
	return f.closer()
}
