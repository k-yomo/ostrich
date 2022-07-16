package collector

import (
	"github.com/k-yomo/ostrich/index"
	"github.com/k-yomo/ostrich/pkg/heap"
	"math"
)

type TopDocsCollector struct {
	limit  int
	offset int
}

type TopDocsResult struct {
	DocAddress index.DocAddress
	Score      float64
}

func NewTopDocsCollector(limit int, offset int) index.Collector[[]*TopDocsResult] {
	return &TopDocsCollector{
		limit:  limit,
		offset: offset,
	}
}

func (t *TopDocsCollector) CollectSegment(w index.Weight, segmentOrd int, segmentReader *index.SegmentReader) []*TopDocsResult {
	topCollector := heap.NewHeap[*TopDocsResult](func(a, b *TopDocsResult) bool {
		return a.Score < b.Score
	})

	heapLen := t.limit + t.offset
	threshold := math.SmallestNonzeroFloat64
	w.ForEachPruning(threshold, segmentReader, func(docID index.DocID, score float64) float64 {
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

	return t.getTopDocsResultsFromHeap(topCollector)
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
