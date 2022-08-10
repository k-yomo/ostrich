package query

import (
	"fmt"
	"github.com/k-yomo/ostrich/postings"
	"github.com/k-yomo/ostrich/reader"
	"github.com/k-yomo/ostrich/schema"
)

type TermQuery struct {
	fieldID schema.FieldID
	term    string
}

func NewTermQuery(fieldID schema.FieldID, term string) reader.Query {
	return &TermQuery{
		fieldID: fieldID,
		term:    term,
	}
}

func (a *TermQuery) Weight(searcher *reader.Searcher, _ bool) (reader.Weight, error) {
	var totalDocNum uint64 = 0
	for _, segmentReader := range searcher.SegmentReaders() {
		totalDocNum += uint64(segmentReader.MaxDoc)
	}
	docFrequency, err := searcher.DocFreq(a.fieldID, a.term)
	if err != nil {
		return nil, fmt.Errorf("get doc frequency: %w", err)
	}

	return &TermWeight{
		fieldID:          a.fieldID,
		term:             a.term,
		similarityWeight: NewTFIDFWeight(totalDocNum, docFrequency),
	}, nil
}

type TermWeight struct {
	fieldID          schema.FieldID
	term             string
	similarityWeight *TfIDFWeight
}

func (a *TermWeight) Scorer(segmentReader *reader.SegmentReader) (reader.Scorer, error) {
	invertedIndexReader := segmentReader.InvertedIndex(a.fieldID)
	postingsReader, err := invertedIndexReader.ReadPostings(a.term)
	if err != nil {
		return nil, fmt.Errorf("read postings: %w", err)
	}
	return &TermScorer{
		postingsReader:   postingsReader,
		similarityWeight: a.similarityWeight,
	}, nil
}

func (a *TermWeight) ForEachPruning(threshold float64, segmentReader *reader.SegmentReader, callback func(docID schema.DocID, score float64) float64) error {
	scorer, err := a.Scorer(segmentReader)
	if err != nil {
		return fmt.Errorf("open scorer: %w", err)
	}
	doc, err := scorer.Doc()
	for err == nil {
		if score := scorer.Score(); score > threshold {
			threshold = callback(doc, threshold)
		}
		doc, err = scorer.Advance()
	}

	return nil
}

type TermScorer struct {
	postingsReader   *postings.PostingsReader
	similarityWeight *TfIDFWeight
}

func (a *TermScorer) Advance() (schema.DocID, error) {
	return a.postingsReader.Advance()
}

func (a *TermScorer) Doc() (schema.DocID, error) {
	return a.postingsReader.Doc()
}

func (a *TermScorer) TermFreq() uint64 {
	return a.postingsReader.TermFreq()
}

func (a *TermScorer) Score() float64 {
	return a.similarityWeight.Score(float64(a.TermFreq()))
}
