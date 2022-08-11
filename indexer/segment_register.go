package indexer

import (
	"sort"
	"sync"

	"github.com/k-yomo/ostrich/index"
	"github.com/k-yomo/ostrich/pkg/list"
)

type SegmentRegister struct {
	segmentStatus map[index.SegmentID]*SegmentEntry
	mu            *sync.Mutex
}

func newSegmentRegister() *SegmentRegister {
	return &SegmentRegister{
		segmentStatus: make(map[index.SegmentID]*SegmentEntry),
	}
}

func newSegmentRegisterFromSegmentMetas(segmentMetas []*index.SegmentMeta) *SegmentRegister {
	segmentStatus := make(map[index.SegmentID]*SegmentEntry)
	for _, segmentMeta := range segmentMetas {
		segmentStatus[segmentMeta.SegmentID] = NewSegmentEntry(segmentMeta)
	}

	return &SegmentRegister{
		segmentStatus: segmentStatus,
	}
}

func (s *SegmentRegister) segmentMetas() []*index.SegmentMeta {
	segmentMetas := make([]*index.SegmentMeta, 0, len(s.segmentStatus))
	for _, entry := range s.segmentStatus {
		segmentMetas = append(segmentMetas, entry.meta)
	}
	sort.Slice(segmentMetas, func(i, j int) bool {
		return segmentMetas[i].SegmentID < segmentMetas[j].SegmentID
	})
	return segmentMetas
}

func (s *SegmentRegister) segmentEntries() []*SegmentEntry {
	entries := make([]*SegmentEntry, 0, len(s.segmentStatus))
	for _, entry := range s.segmentStatus {
		entries = append(entries, entry)
	}
	return entries
}

func (s *SegmentRegister) mergeableSegments(inMergeSegmentIDs []index.SegmentID) []*index.SegmentMeta {
	mergeableSegments := make([]*index.SegmentMeta, 0, len(s.segmentStatus))
	for _, entry := range s.segmentStatus {
		if !list.Contains(inMergeSegmentIDs, entry.SegmentID()) {
			mergeableSegments = append(mergeableSegments, entry.meta)
		}
	}
	return mergeableSegments
}

func (s *SegmentRegister) containsAll(segmentIDs []index.SegmentID) bool {
	for _, segmentID := range segmentIDs {
		if _, ok := s.segmentStatus[segmentID]; !ok {
			return false
		}
	}
	return true
}

func (s *SegmentRegister) addSegmentEntry(segmentEntry *SegmentEntry) {
	s.mu.Lock()
	defer s.mu.Unlock()

	segmentID := segmentEntry.SegmentID()
	s.segmentStatus[segmentID] = segmentEntry
}

func (s *SegmentRegister) removeSegmentEntry(segmentID index.SegmentID) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.segmentStatus, segmentID)
}

func (s *SegmentRegister) clear() {
	s.segmentStatus = map[index.SegmentID]*SegmentEntry{}
}
