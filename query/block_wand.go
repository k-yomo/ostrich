package query

import (
	"github.com/k-yomo/ostrich/schema"
	"sort"
)

type TermScorerWithMaxScore struct {
	*TermScorer
	maxScore float64
}

func newTermScorerWithMaxScore(termScorer *TermScorer) *TermScorerWithMaxScore {
	return &TermScorerWithMaxScore{
		TermScorer: termScorer,
		maxScore:   termScorer.MaxScore(),
	}
}

func blockWand(termScorers []*TermScorer, threshold float64, callback func(docID schema.DocID, score float64) float64) error {
	sort.Slice(termScorers, func(i, j int) bool {
		return termScorers[i].Doc() < termScorers[j].Doc()
	})
	scorers := make([]*TermScorerWithMaxScore, 0, len(termScorers))
	for _, termScorer := range termScorers {
		scorers = append(scorers, newTermScorerWithMaxScore(termScorer))
	}
	for beforePivotLen, pivotLen, pivotDoc, ok := findPivotDoc(scorers, threshold); ok {
		blockMaxScoreUpperbound := scorers[:pivotLen]
	}
}

func findPivotDoc(scorers []*TermScorerWithMaxScore, threshold float64) (int, int, schema.DocID, bool) {
	maxScore := 0.0
	beforePivotLen := 0
	pivotDoc := schema.DocIDTerminated
	for beforePivotLen < len(scorers) {
		scorer := scorers[beforePivotLen]
		maxScore += scorer.maxScore
		if maxScore > threshold {
			pivotDoc = scorer.Doc()
			break
		}
		beforePivotLen += 1
	}
	if pivotDoc == schema.DocIDTerminated {
		return 0, 0, 0, false
	}

	pivotLen := beforePivotLen + 1
	for _, scorer := range scorers[pivotLen:] {
		if scorer.Doc() == pivotDoc {
			pivotLen += 1
		}
	}
	return beforePivotLen, pivotLen, pivotDoc, true
}
