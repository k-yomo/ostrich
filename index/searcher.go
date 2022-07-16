package index

import (
	"github.com/k-yomo/ostrich/schema"
)

type Collector[T any] interface {
	CollectSegment(w Weight, segmentOrd int, segmentReader *SegmentReader) T
	MergeResults(results []T) T
}

type Query interface {
	Weight(searcher *Searcher, scoringEnabled bool) Weight
}

type Weight interface {
	Scorer(segmentReader *SegmentReader) Scorer
	ForEachPruning(threshold float64, reader *SegmentReader, callback func(docID DocID, score float64) float64)
}

type Scorer interface {
	DocSet
	Score() float64
}

type Searcher struct {
	schema         *schema.Schema
	index          *Index
	segmentReaders []*SegmentReader
}

func NewSearcher(schema *schema.Schema, idx *Index, segmentReaders []*SegmentReader) *Searcher {
	return &Searcher{
		schema:         schema,
		index:          idx,
		segmentReaders: segmentReaders,
	}
}

func Search[T any](searcher *Searcher, q Query, c Collector[T]) T {
	results := make([]T, 0, len(searcher.segmentReaders))
	weight := q.Weight(searcher, false)
	for i, segmentReader := range searcher.segmentReaders {
		results = append(results, c.CollectSegment(weight, i, segmentReader))
	}
	return c.MergeResults(results)
}
