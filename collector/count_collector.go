package collector

import (
	"fmt"

	"github.com/k-yomo/ostrich/reader"
	"github.com/k-yomo/ostrich/schema"
)

type CountCollector struct{}

func NewCountCollector() reader.Collector[int] {
	return &CountCollector{}
}

func (c *CountCollector) CollectSegment(w reader.Weight, segmentOrd int, segmentReader *reader.SegmentReader) (int, error) {
	segmentCollector := c.SegmentCollector(segmentOrd)
	err := w.ForEach(segmentReader, func(docID schema.DocID, score float64) {
		segmentCollector.Collect(docID, score)
	})
	if err != nil {
		return 0, fmt.Errorf("collect segment: %w", err)
	}
	return segmentCollector.Result(), nil
}

func (c *CountCollector) MergeResults(results []int) int {
	var sum int
	for _, result := range results {
		sum += result
	}
	return sum
}

func (c *CountCollector) SegmentCollector(segmentOrd int) reader.SegmentCollector[int] {
	return &countSegmentCollector{}
}

type countSegmentCollector struct {
	count int
}

func (c *countSegmentCollector) Collect(docID schema.DocID, score float64) {
	c.count += 1
}

func (c *countSegmentCollector) Result() int {
	return c.count
}
