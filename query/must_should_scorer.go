package query

import (
	"github.com/k-yomo/ostrich/reader"
	"github.com/k-yomo/ostrich/schema"
)

type MustShouldScorer struct {
	mustScorer    reader.Scorer
	shouldScorer  reader.Scorer
	scoreCombiner ScoreCombiner
	scoreCache    *float64
}

func NewMustShouldScorer(
	mustScorer reader.Scorer,
	shouldScorer reader.Scorer,
	scoreCombiner ScoreCombiner,
) reader.Scorer {
	return &MustShouldScorer{
		mustScorer:    mustScorer,
		shouldScorer:  shouldScorer,
		scoreCombiner: scoreCombiner,
	}
}

func (m *MustShouldScorer) Advance() schema.DocID {
	m.scoreCache = nil
	return m.mustScorer.Advance()
}

func (m *MustShouldScorer) Doc() schema.DocID {
	return m.mustScorer.Doc()
}

func (m *MustShouldScorer) Seek(target schema.DocID) schema.DocID {
	return m.mustScorer.Seek(target)
}

func (m *MustShouldScorer) SizeHint() uint32 {
	return m.mustScorer.SizeHint()
}

func (m *MustShouldScorer) Score() float64 {
	if m.scoreCache != nil {
		return *m.scoreCache
	}

	m.scoreCombiner.Update(m.mustScorer)
	doc := m.mustScorer.Doc()
	if m.shouldScorer.Doc() <= doc && m.shouldScorer.Seek(doc) == doc {
		m.scoreCombiner.Update(m.shouldScorer)
	}
	score := m.scoreCombiner.Score()
	m.scoreCache = &score
	m.scoreCombiner.Clear()
	return score
}
