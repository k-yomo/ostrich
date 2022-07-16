package index

import (
	"github.com/k-yomo/ostrich/collector"
	"github.com/k-yomo/ostrich/query"
	"github.com/k-yomo/ostrich/schema"
)

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

func Search[T any](searcher *Searcher, q query.Query, c collector.Collector[T]) []T {
	results := make([][]T, 0, len(searcher.segmentReaders))
	for _, segmentReader := range searcher.segmentReaders {
		results = append(results, c.CollectSegment(segmentReader))
	}
	return c.MergeResults(results)
}
