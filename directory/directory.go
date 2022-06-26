package directory

import (
	"io"
)

type WriteCloseFlasher interface {
	io.WriteCloser
	Flush() error
}

type Directory interface {
	// TODO: Should be able to seek
	Read(path string) (io.ReadCloser, error)
	AtomicRead(path string) ([]byte, error)
	OpenWrite(path string) (WriteCloseFlasher, error)
	AtomicWrite(path string, data []byte) error
	Exists(path string) (bool, error)
}
