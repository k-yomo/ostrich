package directory

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"sync"
)

var _ Directory = &memoryDirectory{}

type memoryDirectory struct {
	// TODO: replace with sync.Map
	pathBytesMap map[string][]byte
	mu           *sync.RWMutex
}

func NewMemoryDirectory() *memoryDirectory {
	return &memoryDirectory{
		pathBytesMap: make(map[string][]byte),
		mu:           &sync.RWMutex{},
	}
}

func (m *memoryDirectory) OpenRead(path string) (ReaderCloser, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	b, ok := m.pathBytesMap[path]
	if !ok {
		return nil, fmt.Errorf("path '%s' does not exist", path)
	}
	return newMemoryBytes(b, m.mu, nil), nil
}

func (m *memoryDirectory) Read(path string) (io.ReadCloser, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	b, ok := m.pathBytesMap[path]
	if !ok {
		return nil, fmt.Errorf("path '%s' does not exist", path)
	}
	return newMemoryBytes(b, m.mu, nil), nil
}

func (m *memoryDirectory) AtomicRead(path string) ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	b, ok := m.pathBytesMap[path]
	if !ok {
		return nil, fs.ErrNotExist
	}
	return b, nil
}

func (m *memoryDirectory) OpenWrite(path string) (WriteCloseSyncer, error) {
	b, ok := m.pathBytesMap[path]
	if !ok {
		b = []byte{}
		m.pathBytesMap[path] = b
	}
	set := func(b []byte) {
		m.mu.Lock()
		defer m.mu.Unlock()
		m.pathBytesMap[path] = b
	}
	return newMemoryBytes(b, m.mu, set), nil
}

func (m *memoryDirectory) AtomicWrite(path string, data []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.pathBytesMap[path] = data
	return nil
}

func (m *memoryDirectory) Exists(path string) (bool, error) {
	_, ok := m.pathBytesMap[path]
	return ok, nil
}

type memoryIO struct {
	*bytes.Reader
	mu  *sync.RWMutex
	set func(b []byte)
}

func newMemoryBytes(b []byte, lock *sync.RWMutex, set func(b []byte)) *memoryIO {
	return &memoryIO{
		Reader: bytes.NewReader(b),
		mu:     lock,
		set:    set,
	}
}

func (m *memoryIO) Write(p []byte) (n int, err error) {
	m.set(p)
	return len(p), nil
}

func (m *memoryIO) Sync() error {
	return nil
}

func (m *memoryIO) Close() error {
	return nil
}
