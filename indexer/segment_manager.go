package indexer

import (
	"sync"

	"github.com/k-yomo/ostrich/index"
)

type SegmentManager struct {
	mu        sync.Mutex
	registers *SegmentRegisters
}

type SegmentRegisters struct {
	uncommitted *SegmentRegister
	committed   *SegmentRegister
}

func NewSegmentManager(segmentMetas []*index.SegmentMeta) *SegmentManager {
	return &SegmentManager{
		registers: &SegmentRegisters{
			uncommitted: newSegmentRegister(),
			committed:   newSegmentRegisterFromSegmentMetas(segmentMetas),
		},
	}
}

func (s *SegmentManager) segmentEntries() []*SegmentEntry {
	return append(s.registers.uncommitted.segmentEntries(), s.registers.committed.segmentEntries()...)
}

func (s *SegmentManager) addSegment(segmentEntry *SegmentEntry) {
	s.registers.uncommitted.addSegmentEntry(segmentEntry)
}

func (s *SegmentManager) commit(segmentEntries []*SegmentEntry) {
	s.registers.committed.clear()
	s.registers.uncommitted.clear()
	for _, segmentEntry := range segmentEntries {
		s.registers.committed.addSegmentEntry(segmentEntry)
	}
}

func (s *SegmentManager) committedSegmentMetas() []*index.SegmentMeta {
	s.removeEmptySegments()
	return s.registers.committed.segmentMetas()
}

func (s *SegmentManager) removeEmptySegments() {
	for _, segmentEntry := range s.registers.committed.segmentEntries() {
		if segmentEntry.meta.DocNum() == 0 {
			s.registers.committed.removeSegmentEntry(segmentEntry.SegmentID())
		}
	}
}
