package query

import (
	"github.com/k-yomo/ostrich/index"
	"io"
)

type AllQuery struct {
}

func NewAllQuery() index.Query {
	return &AllQuery{}
}

func (a *AllQuery) Weight(_ *index.Searcher, _ bool) index.Weight {
	return &AllWeight{}
}

type AllWeight struct {
}

func (a *AllWeight) Scorer(segmentReader *index.SegmentReader) index.Scorer {
	return &AllScorer{
		doc:    0,
		maxDoc: segmentReader.MaxDoc(),
	}
}

func (a *AllWeight) ForEachPruning(threshold float64, segmentReader *index.SegmentReader, callback func(docID index.DocID, score float64) float64) {
	scorer := a.Scorer(segmentReader)
	doc, err := scorer.Doc()
	for err != io.EOF {
		if score := scorer.Score(); score > threshold {
			threshold = callback(doc, threshold)
		}
		doc, err = scorer.Advance()
	}
}

type AllScorer struct {
	doc    index.DocID
	maxDoc index.DocID
}

func (a *AllScorer) Advance() (index.DocID, error) {
	if a.doc <= a.maxDoc {
		a.doc += 1
	}
	return a.Doc()
}

func (a *AllScorer) Doc() (index.DocID, error) {
	if a.doc > a.maxDoc {
		return 0, io.EOF
	}
	return a.doc, nil
}

func (a *AllScorer) Score() float64 {
	return 1.0
}
