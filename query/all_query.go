package query

import (
	"github.com/k-yomo/ostrich/reader"
	"github.com/k-yomo/ostrich/schema"
	"io"
)

type AllQuery struct {
}

func NewAllQuery() reader.Query {
	return &AllQuery{}
}

func (a *AllQuery) Weight(_ *reader.Searcher, _ bool) (reader.Weight, error) {
	return &AllWeight{}, nil
}

type AllWeight struct {
}

func (a *AllWeight) Scorer(segmentReader *reader.SegmentReader) (reader.Scorer, error) {
	return &AllScorer{
		doc:    0,
		maxDoc: segmentReader.MaxDoc,
	}, nil
}

func (a *AllWeight) ForEachPruning(threshold float64, segmentReader *reader.SegmentReader, callback func(docID schema.DocID, score float64) float64) error {
	scorer, err := a.Scorer(segmentReader)
	doc, err := scorer.Doc()
	for err == nil {
		if score := scorer.Score(); score > threshold {
			threshold = callback(doc, threshold)
		}
		doc, err = scorer.Advance()
	}
	if err != io.EOF {
		return err
	}
	return nil
}

type AllScorer struct {
	doc    schema.DocID
	maxDoc schema.DocID
}

func (a *AllScorer) Advance() (schema.DocID, error) {
	if a.doc < a.maxDoc {
		a.doc += 1
	}
	return a.Doc()
}

func (a *AllScorer) Doc() (schema.DocID, error) {
	if a.doc >= a.maxDoc {
		return 0, io.EOF
	}
	return a.doc, nil
}

func (a *AllScorer) Score() float64 {
	return 1.0
}
