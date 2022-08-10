package index

import (
	"sync"

	"github.com/k-yomo/ostrich/directory"
	"github.com/k-yomo/ostrich/schema"
)

type Index struct {
	directory directory.Directory
	schema    *schema.Schema
	inventory *SegmentMetaInventory

	mu *sync.Mutex
}

func NewIndexFromMeta(directory directory.Directory, meta *IndexMeta, inventory *SegmentMetaInventory) *Index {
	return &Index{
		directory: directory,
		schema:    meta.Schema,
		inventory: inventory,
		mu:        &sync.Mutex{},
	}
}

func OpenIndex(dir directory.Directory) (*Index, error) {
	managedDirectory, err := directory.NewManagedDirectory(dir)
	if err != nil {
		return nil, err
	}
	inventory := &SegmentMetaInventory{}
	meta, err := LoadMeta(dir, inventory)
	if err != nil {
		return nil, err
	}
	return NewIndexFromMeta(managedDirectory, meta, inventory), nil
}

func (i *Index) LoadMeta() (*IndexMeta, error) {
	return LoadMeta(i.directory, i.inventory)
}

func (i *Index) SearchableSegments() ([]*Segment, error) {
	meta, err := LoadMeta(i.directory, i.inventory)
	if err != nil {
		return nil, err
	}
	segments := make([]*Segment, 0, len(meta.Segments))
	for _, segmentMeta := range meta.Segments {
		segments = append(segments, newSegment(i, segmentMeta))
	}
	return segments, nil
}

func (i *Index) NewSegmentMeta(segmentID SegmentID, maxDoc schema.DocID) *SegmentMeta {
	return i.inventory.NewSegmentMeta(segmentID, maxDoc)
}

func (i *Index) NewSegment() *Segment {
	segmentMeta := i.inventory.NewSegmentMeta(NewSegmentID(), 0)
	return newSegment(i, segmentMeta)
}

func (i *Index) Segment(segmentMeta *SegmentMeta) *Segment {
	return newSegment(i, segmentMeta)
}

func (i *Index) Directory() directory.Directory {
	return i.directory
}

func (i *Index) Schema() *schema.Schema {
	return i.schema
}
