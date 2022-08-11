package query

import (
	"github.com/k-yomo/ostrich/reader"
	"github.com/k-yomo/ostrich/schema"
)

type EmptyScorer struct{}

func NewEmptyScorer() reader.Scorer {
	return &EmptyScorer{}
}

func (e *EmptyScorer) Advance() schema.DocID {
	return schema.DocIDTerminated
}

func (e *EmptyScorer) Doc() schema.DocID {
	return schema.DocIDTerminated
}

func (e *EmptyScorer) Seek(_ schema.DocID) schema.DocID {
	return schema.DocIDTerminated
}

func (e *EmptyScorer) SizeHint() uint32 {
	return 0
}

func (e *EmptyScorer) Score() float64 {
	return 0
}
