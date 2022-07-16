package collector

import "github.com/k-yomo/ostrich/index"

type Collector[T any] interface {
	CollectSegment(segmentReader *index.SegmentReader) []T
	MergeResults(results [][]T) []T
}
