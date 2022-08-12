package collector

import (
	"fmt"

	"github.com/k-yomo/ostrich/reader"
	"github.com/k-yomo/ostrich/schema"
)

type TupleCollector[L, R any] struct {
	left  reader.Collector[L]
	right reader.Collector[R]
}

type TupleResult[L, R any] struct {
	Left  L
	Right R
}

func NewTupleCollector[L, R any](l reader.Collector[L], r reader.Collector[R]) reader.Collector[*TupleResult[L, R]] {
	return &TupleCollector[L, R]{
		left:  l,
		right: r,
	}
}

func (t *TupleCollector[L, R]) CollectSegment(w reader.Weight, segmentOrd int, segmentReader *reader.SegmentReader) (*TupleResult[L, R], error) {
	leftSegmentCollector := t.left.SegmentCollector(segmentOrd)
	rightSegmentCollector := t.right.SegmentCollector(segmentOrd)
	err := w.ForEach(segmentReader, func(docID schema.DocID, score float64) {
		leftSegmentCollector.Collect(docID, score)
		rightSegmentCollector.Collect(docID, score)
	})
	if err != nil {
		return nil, fmt.Errorf("collect segment: %w", err)
	}

	return &TupleResult[L, R]{
		Left:  leftSegmentCollector.Result(),
		Right: rightSegmentCollector.Result(),
	}, nil
}

func (t *TupleCollector[L, R]) MergeResults(results []*TupleResult[L, R]) *TupleResult[L, R] {
	leftResults := make([]L, 0, len(results))
	rightResults := make([]R, 0, len(results))
	for _, result := range results {
		leftResults = append(leftResults, result.Left)
		rightResults = append(rightResults, result.Right)
	}
	return &TupleResult[L, R]{
		Left:  t.left.MergeResults(leftResults),
		Right: t.right.MergeResults(rightResults),
	}
}

func (t *TupleCollector[L, R]) SegmentCollector(segmentOrd int) reader.SegmentCollector[*TupleResult[L, R]] {
	return &tupleSegmentCollector[L, R]{left: t.left.SegmentCollector(segmentOrd), right: t.right.SegmentCollector(segmentOrd)}
}

type tupleSegmentCollector[L, R any] struct {
	left  reader.SegmentCollector[L]
	right reader.SegmentCollector[R]
}

func (t *tupleSegmentCollector[L, R]) Collect(docID schema.DocID, score float64) {
	t.left.Collect(docID, score)
	t.right.Collect(docID, score)
}

func (t *tupleSegmentCollector[L, R]) Result() *TupleResult[L, R] {
	return &TupleResult[L, R]{Left: t.left.Result(), Right: t.right.Result()}
}
