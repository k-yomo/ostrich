package index

import (
	"encoding/json"
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

func NewIndexFromMetas(directory directory.Directory, metas *IndexMeta, inventory *SegmentMetaInventory) *Index {
	return &Index{
		directory: directory,
		schema:    metas.Schema,
		inventory: inventory,
		mu:        &sync.Mutex{},
	}
}

func (i *Index) LoadMetas() (*IndexMeta, error) {
	metaData, err := i.directory.AtomicRead(metaFileName)
	if err != nil {
		return nil, err
	}
	var indexMeta IndexMeta
	if err := json.Unmarshal(metaData, &indexMeta); err != nil {
		return nil, err
	}

	return &indexMeta, nil
}

func (i *Index) SearchableSegments() ([]*Segment, error) {
	meta, err := i.LoadMetas()
	if err != nil {
		return nil, err
	}
	segments := make([]*Segment, 0, len(meta.Segments))
	for _, segmentMeta := range meta.Segments {
		segments = append(segments, newSegment(i, segmentMeta))
	}
	return segments, nil
}

func (i *Index) NewSegment() *Segment {
	segmentMeta := i.inventory.NewSegmentMeta(NewSegmentID(), 0)
	return newSegment(i, segmentMeta)
}

func (i *Index) Directory() directory.Directory {
	return i.directory
}

func (i *Index) Schema() *schema.Schema {
	return i.schema
}
