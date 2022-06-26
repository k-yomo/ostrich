package indexer

import (
	"github.com/k-yomo/ostrich/index"
	"github.com/k-yomo/ostrich/internal/opstamp"
	"sort"
	"sync"
)

type SegmentUpdater struct {
	indexMeta      *index.IndexMeta
	index          *index.Index
	segmentManager *SegmentManager

	activeMeta *index.IndexMeta

	sync.RWMutex
}

func NewSegmentUpdater(idx *index.Index, indexMeta *index.IndexMeta, stamper *opstamp.Stamper) *SegmentUpdater {
	segmentManager := NewSegmentManager(indexMeta.Segments)

	return &SegmentUpdater{
		indexMeta:      indexMeta,
		index:          idx,
		segmentManager: segmentManager,
	}
}

func (s *SegmentUpdater) AddSegment(segmentEntry *SegmentEntry) {
	s.segmentManager.addSegment(segmentEntry)
}

func (s *SegmentUpdater) Commit(opStamp opstamp.OpStamp) error {
	segmentEntries := s.segmentManager.segmentEntries()
	s.segmentManager.commit(segmentEntries)
	return s.saveMetas(opStamp)
}

func (s *SegmentUpdater) saveMetas(opStamp opstamp.OpStamp) error {
	directory := s.index.Directory()
	committedSegmentMetas := s.segmentManager.committedSegmentMetas()
	// We sort segment_readers by number of documents.
	sort.Slice(committedSegmentMetas, func(i, j int) bool {
		return committedSegmentMetas[i].MaxDoc > committedSegmentMetas[j].MaxDoc
	})

	indexMeta := &index.IndexMeta{
		Segments: committedSegmentMetas,
		Schema:   s.index.Schema(),
		Opstamp:  opStamp,
	}
	if err := index.SaveMetas(indexMeta, directory); err != nil {
		return err
	}
	s.storeMeta(indexMeta)

	return nil
}

func (s *SegmentUpdater) storeMeta(indexMeta *index.IndexMeta) {
	s.Lock()
	defer s.Unlock()
	s.activeMeta = indexMeta
}
