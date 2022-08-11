package query

import (
	"github.com/k-yomo/ostrich/reader"
	"github.com/k-yomo/ostrich/schema"
)

type UnionScorer struct {
	scorers []reader.Scorer
	doc     schema.DocID
	score   float64
}

func NewUnionScorer(scorers []reader.Scorer) reader.Scorer {
	nonEmptyScorers := make([]reader.Scorer, 0, len(scorers))
	for _, scorer := range scorers {
		if !scorer.Doc().IsTerminated() {
			nonEmptyScorers = append(nonEmptyScorers, scorer)
		}
	}
	if len(nonEmptyScorers) == 0 {
		return NewEmptyScorer()
	}

	unionScorer := &UnionScorer{
		scorers: nonEmptyScorers,
	}
	unionScorer.storeMinDocAndScore()
	return unionScorer
}

func (u *UnionScorer) Advance() schema.DocID {
	if u.doc.IsTerminated() {
		return u.doc
	}

	for _, scorer := range u.scorers {
		if scorer.Doc() == u.doc {
			scorer.Advance()
		}
	}
	u.storeMinDocAndScore()
	return u.doc
}

func (u *UnionScorer) Doc() schema.DocID {
	return u.doc
}

func (u *UnionScorer) Seek(target schema.DocID) schema.DocID {
	if u.doc >= target {
		return u.doc
	}
	for _, scorer := range u.scorers {
		scorer.Seek(target)
	}
	u.storeMinDocAndScore()
	return u.doc
}

func (u *UnionScorer) SizeHint() uint32 {
	var maxSizeHint uint32
	for _, scorer := range u.scorers {
		sizeHint := scorer.SizeHint()
		if sizeHint > maxSizeHint {
			maxSizeHint = sizeHint
		}
	}
	return maxSizeHint
}

func (u *UnionScorer) Score() float64 {
	return u.score
}

func (u *UnionScorer) storeMinDocAndScore() {
	minDocID := schema.DocIDTerminated
	scoreCombiner := &sumWithCombiner{}
	for _, scorer := range u.scorers {
		doc := scorer.Doc()
		if doc.IsTerminated() {
			continue
		}
		if doc == minDocID {
			scoreCombiner.Update(scorer)
		} else if doc < minDocID {
			scoreCombiner.Clear()
			scoreCombiner.Update(scorer)
			minDocID = doc
		}
	}
	u.doc = minDocID
	u.score = scoreCombiner.Score()
}
