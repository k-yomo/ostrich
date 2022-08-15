package collector

import (
	"math"

	"github.com/k-yomo/ostrich/index"
	"github.com/k-yomo/ostrich/pkg/heap"
	"github.com/k-yomo/ostrich/pkg/list"
	"github.com/k-yomo/ostrich/reader"
	"github.com/k-yomo/ostrich/schema"
)

type TopScoreCollector struct {
	limit  int
	offset int
}

type TopScoreResult struct {
	DocAddress index.DocAddress
	Score      float64
}

func (t *TopScoreResult) Less(x *TopScoreResult) bool {
	if t.Score == x.Score {
		return t.DocAddress.DocID < t.DocAddress.DocID
	}
	return t.Score < x.Score
}

func NewTopScoreCollector(limit int, offset int) reader.Collector[[]*TopScoreResult] {
	return &TopScoreCollector{
		limit:  limit,
		offset: offset,
	}
}

func (t *TopScoreCollector) CollectSegment(w reader.Weight, segmentOrd int, segmentReader *reader.SegmentReader) ([]*TopScoreResult, error) {
	heapLimit := t.limit + t.offset
	topCollector := heap.NewLimitHeap[*TopScoreResult](heapLimit, func(a, b *TopScoreResult) bool {
		if a.Score == b.Score {
		}
		return a.Score < b.Score
	})

	threshold := math.SmallestNonzeroFloat64
	err := w.ForEachPruning(threshold, segmentReader, func(docID schema.DocID, score float64) float64 {
		topCollector.Push(&TopScoreResult{
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

	results := topCollector.ToArray()
	return list.TakeN(results, t.limit, t.offset), nil
}

func (t *TopScoreCollector) MergeResults(results [][]*TopScoreResult) []*TopScoreResult {
	topCollector := heap.NewLimitHeap[*TopScoreResult](t.limit+t.offset, func(a, b *TopScoreResult) bool {
		return a.Less(b)
	})
	for _, result := range results {
		for _, hit := range result {
			topCollector.Push(hit)
		}
	}

	return list.TakeN(topCollector.ToArray(), t.limit, t.offset)
}

func (t *TopScoreCollector) SegmentCollector(segmentOrd int) reader.SegmentCollector[[]*TopScoreResult] {
	return newTopScoreSegmentCollector(segmentOrd, t.limit, t.offset)
}

type topDocsSegmentCollector struct {
	topCollector *heap.LimitHeap[*TopScoreResult]
	limit        int
	offset       int
	segmentOrd   int
}

func newTopScoreSegmentCollector(segmentOrd int, limit int, offset int) *topDocsSegmentCollector {
	return &topDocsSegmentCollector{
		topCollector: heap.NewLimitHeap[*TopScoreResult](limit+offset, func(a, b *TopScoreResult) bool {
			return a.Score > b.Score
		}),
		limit:      limit,
		offset:     offset,
		segmentOrd: segmentOrd,
	}
}

func (t topDocsSegmentCollector) Collect(docID schema.DocID, score float64) {
	t.topCollector.Push(&TopScoreResult{
		DocAddress: index.DocAddress{
			DocID:      docID,
			SegmentOrd: t.segmentOrd,
		},
		Score: score,
	})
}

func (t topDocsSegmentCollector) Result() []*TopScoreResult {
	return list.TakeN(t.topCollector.ToArray(), t.limit, t.offset)
}
