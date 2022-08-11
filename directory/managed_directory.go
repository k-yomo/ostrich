package directory

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"strings"
	"sync"
)

const managedFilePath = ".managed.json"

type ManagedDirectory struct {
	directory       Directory
	MetaInformation *MetaInformation
}

type MetaInformation struct {
	sync.RWMutex
	ManagedPaths []string
}

func (m *MetaInformation) AddPath(path string) (added bool) {
	for _, managedPath := range m.ManagedPaths {
		if path == managedPath {
			return false
		}
	}

	m.ManagedPaths = append(m.ManagedPaths, path)
	return true
}

func (m *MetaInformation) RemovePath(path string) {
	newPaths := make([]string, 0, len(m.ManagedPaths))
	for _, managedPath := range m.ManagedPaths {
		if path == managedPath {
			continue
		}
		newPaths = append(newPaths, managedPath)
	}
	m.ManagedPaths = newPaths
}

func NewManagedDirectory(dir Directory) (*ManagedDirectory, error) {
	bytes, err := dir.AtomicRead(managedFilePath)
	if errors.Is(err, fs.ErrNotExist) {
		return &ManagedDirectory{
			directory:       dir,
			MetaInformation: &MetaInformation{},
		}, nil
	}
	if err != nil {
		return nil, err
	}
	var managedFiles []string
	if err := json.Unmarshal(bytes, &managedFiles); err != nil {
		return nil, err
	}

	return &ManagedDirectory{
		directory: dir,
		MetaInformation: &MetaInformation{
			ManagedPaths: managedFiles,
		},
	}, nil
}

func (m *ManagedDirectory) OpenRead(path string) (*FileSlice, error) {
	return m.directory.OpenRead(path)
}
func (m *ManagedDirectory) AtomicRead(path string) ([]byte, error) {
	return m.directory.AtomicRead(path)
}
func (m *ManagedDirectory) OpenWrite(path string) (WriteCloseSyncer, error) {
	if err := m.registerFileAsManaged(path); err != nil {
		return nil, err
	}
	return m.directory.OpenWrite(path)
}
func (m *ManagedDirectory) AtomicWrite(path string, data []byte) error {
	if err := m.registerFileAsManaged(path); err != nil {
		return fmt.Errorf("register file: %w", err)
	}
	return m.directory.AtomicWrite(path, data)
}
func (m *ManagedDirectory) Exists(path string) (bool, error) {
	return m.directory.Exists(path)
}

func (m *ManagedDirectory) registerFileAsManaged(path string) error {
	if !isManageableFile(path) {
		return nil
	}

	m.MetaInformation.Lock()
	defer m.MetaInformation.Unlock()

	if !m.MetaInformation.AddPath(path) {
		return nil
	}

	managedPathsJSON, err := json.Marshal(m.MetaInformation.ManagedPaths)
	if err != nil {
		m.MetaInformation.RemovePath(path)
		return fmt.Errorf("marshal managed paths: %w", err)
	}
	if err := m.directory.AtomicWrite(managedFilePath, managedPathsJSON); err != nil {
		m.MetaInformation.RemovePath(path)
		return fmt.Errorf("write managed paths: %w", err)
	}

	return nil
}

// Filenames that starts by a "." are not managed.
func isManageableFile(path string) bool {
	return !strings.HasPrefix(path, ".")
}
