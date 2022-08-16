package directory

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"strings"
	"sync"

	"github.com/k-yomo/ostrich/internal/logging"
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
func (m *ManagedDirectory) Delete(path string) error {
	return m.directory.Delete(path)
}

func (m *ManagedDirectory) GarbageCollect(livingFilePaths []string) error {
	m.MetaInformation.Lock()
	defer m.MetaInformation.Unlock()

	livingFileMap := make(map[string]struct{}, len(livingFilePaths))
	for _, livingFile := range livingFilePaths {
		livingFileMap[livingFile] = struct{}{}
	}

	var filesToDelete []string
	for _, managedPath := range m.MetaInformation.ManagedPaths {
		if _, ok := livingFileMap[managedPath]; !ok {
			filesToDelete = append(filesToDelete, managedPath)
		}
	}

	var deletedFiles, failedFiles []string
	for _, fileToDelete := range filesToDelete {
		if err := m.Delete(fileToDelete); err != nil {
			failedFiles = append(failedFiles, fileToDelete)
			continue
		}
		logging.Logger().Debug("deleted", "path", fileToDelete)
		deletedFiles = append(deletedFiles, fileToDelete)
	}

	if len(deletedFiles) == 0 {
		return nil
	}

	for _, deletedFile := range deletedFiles {
		m.MetaInformation.RemovePath(deletedFile)
	}
	if err := m.saveManagedPaths(); err != nil {
		return err
	}

	return nil
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

	if err := m.saveManagedPaths(); err != nil {
		m.MetaInformation.RemovePath(path)
		return err
	}

	return nil
}

func (m *ManagedDirectory) saveManagedPaths() error {
	managedPathsJSON, err := json.Marshal(m.MetaInformation.ManagedPaths)
	if err != nil {
		return fmt.Errorf("marshal managed paths: %w", err)
	}
	if err := m.directory.AtomicWrite(managedFilePath, managedPathsJSON); err != nil {
		return fmt.Errorf("write managed paths: %w", err)
	}
	return nil
}

// Filenames that starts by a "." are not managed.
func isManageableFile(path string) bool {
	return !strings.HasPrefix(path, ".")
}
