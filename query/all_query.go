package query

import (
	"fmt"

	"github.com/k-yomo/ostrich/index"
	"github.com/k-yomo/ostrich/reader"
	"github.com/k-yomo/ostrich/schema"
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
	if err != nil {
		return fmt.Errorf("get scorer: %v", err)
	}
	doc := scorer.Doc()
	for !doc.IsTerminated() {
		if score := scorer.Score(); score > threshold {
			threshold = callback(doc, threshold)
		}
		doc = scorer.Advance()
	}
	return nil
}

type AllScorer struct {
	doc    schema.DocID
	maxDoc schema.DocID
}

func (a *AllScorer) Advance() schema.DocID {
	if a.doc < a.maxDoc {
		a.doc += 1
	}
	return a.Doc()
}

func (a *AllScorer) Doc() schema.DocID {
	if a.doc >= a.maxDoc {
		return schema.DocIDTerminated
	}
	return a.doc
}

func (a *AllScorer) Seek(target schema.DocID) schema.DocID {
	return index.SeekDocSet(a, target)
}

func (a *AllScorer) SizeHint() uint32 {
	return uint32(a.maxDoc)
}

func (a *AllScorer) Score() float64 {
	return 1.0
}
