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
	topCollector := heap.NewHeap[*TopDocsResult](func(a, b *TopDocsResult) bool {
		return a.Score < b.Score
	})

	heapLen := t.limit + t.offset
	threshold := math.SmallestNonzeroFloat64
	err := w.ForEachPruning(threshold, segmentReader, func(docID schema.DocID, score float64) float64 {
		if topCollector.Len() < heapLen {
			topCollector.Push(&TopDocsResult{
				DocAddress: index.DocAddress{
					SegmentOrd: segmentOrd,
					DocID:      docID,
				},
				Score: score,
			})
		} else {
			if head := topCollector.Peek(); score > (*head).Score {
				*head = &TopDocsResult{
					DocAddress: index.DocAddress{
						SegmentOrd: segmentOrd,
						DocID:      docID,
					},
					Score: score,
				}
			}
		}
		if topCollector.Len() == heapLen {
			threshold = (*topCollector.Peek()).Score
		}
		return threshold
	})

	if err != nil {
		return nil, err
	}

	return t.getTopDocsResultsFromHeap(topCollector), nil
}

func (t *TopDocsCollector) MergeResults(results [][]*TopDocsResult) []*TopDocsResult {
	topCollector := heap.NewHeap[*TopDocsResult](func(a, b *TopDocsResult) bool {
		return a.Score < b.Score
	})
	for _, result := range results {
		for _, hit := range result {
			if topCollector.Len() < t.limit+t.limit {
				topCollector.Push(hit)
			} else {
				if head := topCollector.Peek(); hit.Score > (*head).Score {
					*head = hit
				}
			}
		}
	}

	return t.getTopDocsResultsFromHeap(topCollector)
}

func (t *TopDocsCollector) getTopDocsResultsFromHeap(h *heap.Heap[*TopDocsResult]) []*TopDocsResult {
	topResults := make([]*TopDocsResult, 0, h.Len())
	for i := 0; i < t.limit; i++ {
		if h.Len() == 0 {
			break
		}
		result := h.Pop()
		if i < t.offset {
			continue
		}
		topResults = append(topResults, result)
	}
	return topResults
}
