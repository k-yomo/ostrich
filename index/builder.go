package index

import (
	"errors"
	"fmt"

	"github.com/k-yomo/ostrich/directory"
	"github.com/k-yomo/ostrich/schema"
)

type Builder struct {
	schema *schema.Schema
}

func NewBuilder(indexSchema *schema.Schema) *Builder {
	return &Builder{schema: indexSchema}
}

func (b *Builder) CreateInDir(path string) (*Index, error) {
	mmapDirectory, err := directory.NewMMapDirectory(path)
	if err != nil {
		return nil, fmt.Errorf("open mmap directory: %v", err)
	}
	exists, err := mmapDirectory.Exists(metaFileName)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("previous index is in the directory")
	} else {
		return b.create(mmapDirectory)
	}
}

func (b *Builder) OpenOrCreate(path string) (*Index, error) {
	mmapDirectory, err := directory.NewMMapDirectory(path)
	if err != nil {
		return nil, fmt.Errorf("open mmap directory: %v", err)
	}
	exists, err := mmapDirectory.Exists(metaFileName)
	if err != nil {
		return nil, err
	}
	if exists {
		return OpenIndex(mmapDirectory)
	} else {
		return b.create(mmapDirectory)
	}
}

func (b *Builder) create(dir directory.Directory) (*Index, error) {
	managedDirectory, err := directory.NewManagedDirectory(dir)
	if err != nil {
		return nil, err
	}
	if err := b.saveNewMeta(managedDirectory); err != nil {
		return nil, err
	}
	indexMeta := NewIndexMeta(b.schema)
	return NewIndexFromMeta(managedDirectory, indexMeta, &SegmentMetaInventory{}), nil
}

func (b *Builder) saveNewMeta(dir directory.Directory) error {
	return SaveMeta(NewIndexMeta(b.schema), dir)
}
