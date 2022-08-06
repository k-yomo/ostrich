package reader

import (
	"github.com/k-yomo/ostrich/index"
	"github.com/k-yomo/ostrich/schema"
)

type Collector[T any] interface {
	CollectSegment(w Weight, segmentOrd int, segmentReader *SegmentReader) (T, error)
	MergeResults(results []T) T
}

type Query interface {
	Weight(searcher *Searcher, scoringEnabled bool) Weight
}

type Weight interface {
	Scorer(segmentReader *SegmentReader) (Scorer, error)
	ForEachPruning(threshold float64, reader *SegmentReader, callback func(docID schema.DocID, score float64) float64) error
}

type Scorer interface {
	index.DocSet
	Score() float64
}

type Searcher struct {
	schema         *schema.Schema
	index          *index.Index
	segmentReaders []*SegmentReader
}

func NewSearcher(idx *index.Index, segmentReaders []*SegmentReader) *Searcher {
	return &Searcher{
		schema:         idx.Schema(),
		index:          idx,
		segmentReaders: segmentReaders,
	}
}

func Search[T any](searcher *Searcher, q Query, c Collector[T]) (T, error) {
	results := make([]T, 0, len(searcher.segmentReaders))
	weight := q.Weight(searcher, false)
	for i, segmentReader := range searcher.segmentReaders {
		result, err := c.CollectSegment(weight, i, segmentReader)
		if err != nil {
			var temp T
			return temp, err
		}
		results = append(results, result)
	}
	return c.MergeResults(results), nil
}

func (s *Searcher) Close() error {
	for _, segmentReader := range s.segmentReaders {
		if err := segmentReader.termdictFile.Close(); err != nil {
			return err
		}
		if err := segmentReader.storeFile.Close(); err != nil {
			return err
		}
		if err := segmentReader.postingsFile.Close(); err != nil {
			return err
		}
	}
	return nil
}
