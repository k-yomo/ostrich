package indexer

import (
	"sync"

	"github.com/k-yomo/ostrich/index"
)

type SegmentStatus int

const (
	SegmentStatusUnknown = iota
	SegmentStatusUncommitted
	SegmentStatusCommitted
)

type SegmentManager struct {
	mu        sync.RWMutex
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
	s.mu.Lock()
	defer s.mu.Unlock()

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

func (s *SegmentManager) segmentStatus(segmentIDs []index.SegmentID) SegmentStatus {
	if s.registers.committed.containsAll(segmentIDs) {
		return SegmentStatusCommitted
	} else if s.registers.uncommitted.containsAll(segmentIDs) {
		return SegmentStatusUncommitted
	} else {
		return SegmentStatusUnknown
	}
}

func (s *SegmentManager) targetRegister(status SegmentStatus) *SegmentRegister {
	switch status {
	case SegmentStatusCommitted:
		return s.registers.committed
	case SegmentStatusUncommitted:
		return s.registers.uncommitted
	default:
		panic("no target register for unknown segment status")
	}
}

func (s *SegmentManager) mergeableSegments(inMergeSegmentIDs []index.SegmentID) (commited []*index.SegmentMeta, uncommited []*index.SegmentMeta) {
	return s.registers.committed.mergeableSegments(inMergeSegmentIDs),
		s.registers.uncommitted.mergeableSegments(inMergeSegmentIDs)
}

func (s *SegmentManager) segmentEntriesForMerge(segmentIDs []index.SegmentID) []*SegmentEntry {
	segmentEntries := make([]*SegmentEntry, 0, len(segmentIDs))
	targetRegister := s.targetRegister(s.segmentStatus(segmentIDs))
	for _, segmentID := range segmentIDs {
		segmentEntries = append(segmentEntries, targetRegister.segmentStatus[segmentID])
	}
	return segmentEntries
}

func (s *SegmentManager) endMerge(beforeMergeSegmentIDs []index.SegmentID, mergedSegmentEntry *SegmentEntry) SegmentStatus {
	s.mu.Lock()
	defer s.mu.Unlock()

	segmentStatus := s.segmentStatus(beforeMergeSegmentIDs)
	targetRegister := s.targetRegister(segmentStatus)
	for _, segmentID := range beforeMergeSegmentIDs {
		targetRegister.removeSegmentEntry(segmentID)
	}
	targetRegister.addSegmentEntry(mergedSegmentEntry)

	return segmentStatus
}
