package collector

import (
	"math"

	"github.com/k-yomo/ostrich/index"
	"github.com/k-yomo/ostrich/pkg/heap"
	"github.com/k-yomo/ostrich/reader"
	"github.com/k-yomo/ostrich/schema"
)

type TopDocsCollector struct {
	limit  int
	offset int
}

type TopDocsResult struct {
	DocAddress index.DocAddress
	Score      float64
}

func NewTopDocsCollector(limit int, offset int) reader.Collector[[]*TopDocsResult] {
	return &TopDocsCollector{
		limit:  limit,
		offset: offset,
	}
}

func (t *TopDocsCollector) CollectSegment(w reader.Weight, segmentOrd int, segmentReader *reader.SegmentReader) ([]*TopDocsResult, error) {
	heapLimit := t.limit + t.offset
	topCollector := heap.NewLimitHeap[*TopDocsResult](heapLimit, func(a, b *TopDocsResult) bool {
		return a.Score < b.Score
	})

	threshold := math.SmallestNonzeroFloat64
	err := w.ForEachPruning(threshold, segmentReader, func(docID schema.DocID, score float64) float64 {
		topCollector.Push(&TopDocsResult{
			DocAddress: index.DocAddress{
				SegmentOrd: segmentOrd,
				DocID:      docID,
			},
			Score: score,
		})
		if topCollector.Len() == heapLimit {
			threshold = (*topCollector.Peek()).Score
		}
		return threshold
	})

	if err != nil {
		return nil, err
	}

	return topCollector.TopN(t.limit, t.offset), nil
}

func (t *TopDocsCollector) MergeResults(results [][]*TopDocsResult) []*TopDocsResult {
	topCollector := heap.NewLimitHeap[*TopDocsResult](t.limit+t.offset, func(a, b *TopDocsResult) bool {
		return a.Score < b.Score
	})
	for _, result := range results {
		for _, hit := range result {
			topCollector.Push(hit)
		}
	}

	return topCollector.TopN(t.limit, t.offset)
}

func (t *TopDocsCollector) SegmentCollector(segmentOrd int) reader.SegmentCollector[[]*TopDocsResult] {
	return newTopDocsSegmentCollector(segmentOrd, t.limit, t.offset)
}

type topDocsSegmentCollector struct {
	topCollector *heap.LimitHeap[*TopDocsResult]
	limit        int
	offset       int
	segmentOrd   int
}

func newTopDocsSegmentCollector(segmentOrd int, limit int, offset int) *topDocsSegmentCollector {
	return &topDocsSegmentCollector{
		topCollector: heap.NewLimitHeap[*TopDocsResult](limit+offset, func(a, b *TopDocsResult) bool {
			return a.Score < b.Score
		}),
		limit:      limit,
		offset:     offset,
		segmentOrd: segmentOrd,
	}
}

func (t topDocsSegmentCollector) Collect(docID schema.DocID, score float64) {
	t.topCollector.Push(&TopDocsResult{
		DocAddress: index.DocAddress{
			DocID:      docID,
			SegmentOrd: t.segmentOrd,
		},
		Score: score,
	})
}

func (t topDocsSegmentCollector) Result() []*TopDocsResult {
	return t.topCollector.TopN(t.limit, t.offset)
}
