package indexer

import (
	"fmt"
	"log"
	"sort"
	"sync"

	"github.com/k-yomo/ostrich/index"
	"github.com/k-yomo/ostrich/internal/opstamp"
)

type SegmentUpdater struct {
	index          *index.Index
	activeMeta     *index.IndexMeta
	segmentManager *SegmentManager
	mergePolicy    MergePolicy

	stamper *opstamp.Stamper

	segmentIDsInMerge []index.SegmentID

	sync.RWMutex
}

func NewSegmentUpdater(idx *index.Index, indexMeta *index.IndexMeta, stamper *opstamp.Stamper) *SegmentUpdater {
	segmentManager := NewSegmentManager(indexMeta.Segments)

	return &SegmentUpdater{
		activeMeta:     indexMeta,
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
	s.Lock()
	defer s.Unlock()

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
	if err := index.SaveMeta(indexMeta, directory); err != nil {
		return fmt.Errorf("save meta: %w", err)
	}
	s.storeMeta(indexMeta)

	return nil
}

func (s *SegmentUpdater) storeMeta(indexMeta *index.IndexMeta) {
	s.activeMeta = indexMeta
}

func (s *SegmentUpdater) considerMergeOptions() {
	s.Lock()
	curOpstamp := s.stamper.Stamp()
	committedSegments, uncommittedSegments := s.mergeableSegments()

	var mergeOperations []*MergeOperation
	committedMergeCandidates := s.mergePolicy.ComputeMergeCandidates(committedSegments)
	for _, segmentIDs := range committedMergeCandidates {
		mergeOperations = append(mergeOperations, NewMergeOperation(curOpstamp, segmentIDs))
		s.segmentIDsInMerge = append(s.segmentIDsInMerge, segmentIDs...)
	}

	uncommittedMergeCandidates := s.mergePolicy.ComputeMergeCandidates(uncommittedSegments)
	for _, segmentIDs := range uncommittedMergeCandidates {
		mergeOperations = append(mergeOperations, NewMergeOperation(s.activeMeta.Opstamp, segmentIDs))
		s.segmentIDsInMerge = append(s.segmentIDsInMerge, segmentIDs...)
	}
	s.Unlock()

	for _, mergeOperation := range mergeOperations {
		if _, err := s.startMerge(mergeOperation); err != nil {
			log.Printf("failed to merge: %v\n", err)
		}
	}
}

func (s *SegmentUpdater) mergeableSegments() ([]*index.SegmentMeta, []*index.SegmentMeta) {
	return s.segmentManager.mergeableSegments(s.segmentIDsInMerge)
}

func (s *SegmentUpdater) startMerge(operation *MergeOperation) (*index.SegmentMeta, error) {
	segmentEntries := s.segmentManager.segmentEntriesForMerge(operation.segmentIDs)
	mergedSegmentEntry, err := merge(s.index, segmentEntries, operation.targetOpStamp)
	if err != nil {
		return nil, err
	}

	s.Lock()
	segmentStatus := s.segmentManager.endMerge(operation.segmentIDs, mergedSegmentEntry)
	if segmentStatus == SegmentStatusCommitted {
		if err := s.saveMetas(s.activeMeta.Opstamp); err != nil {
			return nil, err
		}
	}
	s.index.RemoveSegments(operation.segmentIDs)

	newSegmentIDsInMerge := make([]index.SegmentID, 0, len(s.segmentIDsInMerge))
	processedSegmentIDMap := make(map[index.SegmentID]struct{})
	for _, segmentID := range operation.segmentIDs {
		processedSegmentIDMap[segmentID] = struct{}{}
	}
	for _, segmentID := range s.segmentIDsInMerge {
		if _, ok := processedSegmentIDMap[segmentID]; !ok {
			newSegmentIDsInMerge = append(newSegmentIDsInMerge, segmentID)
		}
	}
	s.segmentIDsInMerge = newSegmentIDsInMerge
	s.Unlock()

	s.considerMergeOptions()

	if err := s.garbageCollectFiles(); err != nil {
		// logging?
	}
	return mergedSegmentEntry.meta, nil
}

func (s *SegmentUpdater) garbageCollectFiles() error {
	return s.index.Directory().GarbageCollect(append(s.index.ListSegmentFilePaths(), index.MetaFileName))
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
	docNum, err := indexMerger.Write(segmentSerializer)
	if err != nil {
		return nil, fmt.Errorf("merger write: %w", err)
	}

	segmentMeta := idx.NewSegmentMeta(mergedSegment.Meta().SegmentID, docNum)
	return NewSegmentEntry(segmentMeta), nil
}
