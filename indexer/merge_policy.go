package indexer

import "github.com/k-yomo/ostrich/index"

type MergePolicy interface {
	ComputeMergeCandidates(segmentMetas []*index.SegmentMeta) [][]index.SegmentID
}

type NoMergePolicy struct {
}

func (n *NoMergePolicy) ComputeMergeCandidates(_ []*index.SegmentMeta) [][]index.SegmentID {
	return nil
}
