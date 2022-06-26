package indexer

import (
	"github.com/k-yomo/ostrich/index"
	"sort"
)

type SegmentRegister struct {
	segmentStatus map[index.SegmentID]*SegmentEntry
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

func (s *SegmentRegister) addSegmentEntry(segmentEntry *SegmentEntry) {
	segmentID := segmentEntry.SegmentID()
	s.segmentStatus[segmentID] = segmentEntry
}

func (s *SegmentRegister) removeSegmentEntry(segmentID index.SegmentID) {
	delete(s.segmentStatus, segmentID)
}

func (s *SegmentRegister) clear() {
	s.segmentStatus = map[index.SegmentID]*SegmentEntry{}
}
