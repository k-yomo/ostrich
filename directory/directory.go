package directory

import (
	"io"
)

type WriteCloseSyncer interface {
	io.WriteCloser
	Sync() error
}

type Directory interface {
	OpenRead(path string) (*FileSlice, error)
	AtomicRead(path string) ([]byte, error)
	OpenWrite(path string) (WriteCloseSyncer, error)
	AtomicWrite(path string, data []byte) error
	Exists(path string) (bool, error)
	Delete(path string) error
}
