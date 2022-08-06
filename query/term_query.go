package query

import (
	"fmt"
	"github.com/k-yomo/ostrich/reader"
	"github.com/k-yomo/ostrich/schema"
	"io"
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

func (a *TermQuery) Weight(_ *reader.Searcher, _ bool) reader.Weight {
	return &TermWeight{
		fieldID: a.fieldID,
		term:    a.term,
	}
}

type TermWeight struct {
	fieldID schema.FieldID
	term    string
}

func (a *TermWeight) Scorer(segmentReader *reader.SegmentReader) (reader.Scorer, error) {
	postingsReader, err := segmentReader.InvertedIndex(a.fieldID)
	if err != nil {
		return nil, fmt.Errorf("initialize inverted index: %w", err)
	}
	postingList, err := postingsReader.ReadPostings(a.fieldID, a.term)
	if err != nil {
		return nil, fmt.Errorf("read postings: %w", err)
	}
	return &TermScorer{
		postingList: postingList,
		curIdx:      0,
	}, nil
}

func (a *TermWeight) ForEachPruning(threshold float64, segmentReader *reader.SegmentReader, callback func(docID schema.DocID, score float64) float64) error {
	scorer, err := a.Scorer(segmentReader)
	if err != nil {
		return fmt.Errorf("open scorer: %w", err)
	}
	doc, err := scorer.Doc()
	for err != io.EOF {
		if score := scorer.Score(); score > threshold {
			threshold = callback(doc, threshold)
		}
		doc, err = scorer.Advance()
	}

	return nil
}

type TermScorer struct {
	postingList []schema.DocID
	curIdx      int
}

func (a *TermScorer) Advance() (schema.DocID, error) {
	if a.curIdx < len(a.postingList) {
		a.curIdx += 1
	}
	return a.Doc()
}

func (a *TermScorer) Doc() (schema.DocID, error) {
	if a.curIdx >= len(a.postingList) {
		return 0, io.EOF
	}
	return a.postingList[a.curIdx], nil
}

func (a *TermScorer) Score() float64 {
	return 1.0
}
