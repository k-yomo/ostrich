package index

import "github.com/k-yomo/ostrich/pkg/uuid"

type SegmentID string

func NewSegmentID() SegmentID {
	return SegmentID(uuid.Generate())
}
