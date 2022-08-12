package reader

import (
	"github.com/k-yomo/ostrich/index"
	"github.com/k-yomo/ostrich/schema"
)

type Collector[T any] interface {
	CollectSegment(w Weight, segmentOrd int, segmentReader *SegmentReader) (T, error)
	MergeResults(results []T) T
	SegmentCollector(segmentOrd int) SegmentCollector[T]
}

type SegmentCollector[T any] interface {
	Collect(docID schema.DocID, score float64)
	Result() T
}

type Query interface {
	Weight(searcher *Searcher, scoringEnabled bool) (Weight, error)
}

type Weight interface {
	Scorer(segmentReader *SegmentReader) (Scorer, error)
	ForEach(reader *SegmentReader, callback func(docID schema.DocID, score float64)) error
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

func (s *Searcher) SegmentReaders() []*SegmentReader {
	return s.segmentReaders
}

func (s *Searcher) DocFreq(fieldID schema.FieldID, term string) (int, error) {
	totalDocFreq := 0
	for _, segmentReader := range s.segmentReaders {
		postingsReader := segmentReader.InvertedIndex(fieldID)
		totalDocFreq += postingsReader.DocFreq(term)
	}
	return totalDocFreq, nil
}

func (s *Searcher) Close() error {
	for _, segmentReader := range s.segmentReaders {
		// if err := segmentReader.storeFile.Close(); err != nil {
		// 	return err
		// }
		if err := segmentReader.postingsFile.Close(); err != nil {
			return err
		}
	}
	return nil
}
