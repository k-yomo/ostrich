package index

type SegmentComponent int

const (
	SegmentComponentPostings SegmentComponent = iota + 1
	SegmentComponentTerms
	SegmentComponentStore
	SegmentComponentDelete
)
