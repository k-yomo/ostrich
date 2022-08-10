package indexer

import (
	"github.com/k-yomo/ostrich/index"
	"github.com/k-yomo/ostrich/internal/opstamp"
)

type MergeOperation struct {
	targetOpStamp opstamp.OpStamp
	segmentIDs    []index.SegmentID
}

func NewMergeOperation(targetOpStamp opstamp.OpStamp, segmentIDs []index.SegmentID) *MergeOperation {
	return &MergeOperation{
		targetOpStamp: targetOpStamp,
		segmentIDs:    segmentIDs,
	}
}
