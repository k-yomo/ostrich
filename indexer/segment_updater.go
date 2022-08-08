package indexer

import (
	"fmt"
	"github.com/k-yomo/ostrich/index"
	"github.com/k-yomo/ostrich/internal/opstamp"
	"sort"
	"sync"
)

type SegmentUpdater struct {
	indexMeta      *index.IndexMeta
	index          *index.Index
	segmentManager *SegmentManager
	mergePolicy    MergePolicy

	activeMeta *index.IndexMeta

	stamper *opstamp.Stamper

	sync.RWMutex
}

func NewSegmentUpdater(idx *index.Index, indexMeta *index.IndexMeta, stamper *opstamp.Stamper) *SegmentUpdater {
	segmentManager := NewSegmentManager(indexMeta.Segments)

	return &SegmentUpdater{
		indexMeta:      indexMeta,
		index:          idx,
		segmentManager: segmentManager,
		mergePolicy:    NewLogMergePolicy(),
		stamper:        stamper,
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

func (s *SegmentUpdater) ConsiderMergeOptions() {
	committedSegments, uncommittedSegments := s.mergeableSegments()

	var mergeOperations []*MergeOperation

	curOpstamp := s.stamper.Stamp()
	committedMergeCandidates := s.mergePolicy.ComputeMergeCandidates(committedSegments)
	for _, segmentIDs := range committedMergeCandidates {
		mergeOperations = append(mergeOperations, NewMergeOperation(curOpstamp, segmentIDs))
	}

	uncommittedMergeCandidates := s.mergePolicy.ComputeMergeCandidates(uncommittedSegments)
	for _, segmentIDs := range uncommittedMergeCandidates {
		mergeOperations = append(mergeOperations, NewMergeOperation(s.activeMeta.Opstamp, segmentIDs))
	}

	for _, mergeOperation := range mergeOperations {
		s.startMerge(mergeOperation)
	}
}

func (s *SegmentUpdater) mergeableSegments() ([]*index.SegmentMeta, []*index.SegmentMeta) {
	// TODO: pass segment ids in merge
	return s.segmentManager.mergeableSegments([]index.SegmentID{})
}

func (s *SegmentUpdater) startMerge(operation *MergeOperation) ([]*SegmentEntry, error) {
	segmentEntries := s.segmentManager.segmentEntriesForMerge(operation.segmentIDs)
}

func merge(idx *index.Index, segmentEntries []*SegmentEntry, targetOpStamp opstamp.OpStamp) (*SegmentEntry, error) {
	mergedSegment := idx.NewSegment()

	segments := make([]*index.Segment, 0, len(segmentEntries))
	for _, segmentEntry := range segmentEntries {
		segments = append(segments, idx.Segment(segmentEntry.meta))
	}

	indexMerger, err := NewIndexMerger(idx.Schema(), segments)
	if err != nil {
		return nil, fmt.Errorf("initialize index merger: %w", err)
	}

	segmentSerializer, err := NewSegmentSerializer(mergedSegment)
	if err != nil {
		return nil, fmt.Errorf("initialize segment serializer: %w", err)
	}
	numDocs, err := indexMerger.Write(segmentSerializer)
	if err != nil {
		return nil, fmt.Errorf("merger write: %w", err)
	}

	segmentMeta := idx.NewSegmentMeta(mergedSegment.Meta().SegmentID, numDocs)
	return NewSegmentEntry(segmentMeta), nil
}
