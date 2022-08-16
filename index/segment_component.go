package index

type SegmentComponent int

const (
	SegmentComponentPostings SegmentComponent = iota + 1
	SegmentComponentTerms
	SegmentComponentStore
	SegmentComponentDelete
)

var segmentComponents = []SegmentComponent{
	SegmentComponentPostings,
	SegmentComponentTerms,
	SegmentComponentStore,
	// SegmentComponentDelete,
}
