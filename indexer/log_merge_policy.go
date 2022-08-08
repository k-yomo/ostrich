package indexer

import "github.com/k-yomo/ostrich/index"

type LogMergePolicy struct {
}

func NewLogMergePolicy() *LogMergePolicy {
	return &LogMergePolicy{}
}

func (n *LogMergePolicy) ComputeMergeCandidates(_ []*index.SegmentMeta) [][]index.SegmentID {
	return nil
}
