package query

import (
	"fmt"

	"github.com/k-yomo/ostrich/reader"
	"github.com/k-yomo/ostrich/schema"
)

type Occur int

const (
	OccurShould Occur = iota
	OccurMust
)

type BooleanQuery struct {
	subQueries []*BooleanSubQuery
}

type BooleanSubQuery struct {
	occur Occur
	query reader.Query
}

func NewBooleanSubQuery(occur Occur, query reader.Query) *BooleanSubQuery {
	return &BooleanSubQuery{
		occur: occur,
		query: query,
	}
}

func (b *BooleanSubQuery) Weight(searcher *reader.Searcher, scoringEnabled bool) (*BooleanSubWeight, error) {
	weight, err := b.query.Weight(searcher, scoringEnabled)
	if err != nil {
		return nil, err
	}
	return &BooleanSubWeight{
		occur:  b.occur,
		weight: weight,
	}, nil
}

func NewBooleanQuery(subQueries []*BooleanSubQuery) reader.Query {
	return &BooleanQuery{
		subQueries: subQueries,
	}
}

func NewBooleanIntersectionQuery(subQueries []reader.Query) reader.Query {
	booleanSubQueries := make([]*BooleanSubQuery, 0, len(subQueries))
	for _, subQuery := range subQueries {
		booleanSubQueries = append(booleanSubQueries, NewBooleanSubQuery(OccurMust, subQuery))
	}
	return &BooleanQuery{
		booleanSubQueries,
	}
}

func NewBooleanUnionQuery(subQueries []reader.Query) reader.Query {
	booleanSubQueries := make([]*BooleanSubQuery, 0, len(subQueries))
	for _, subQuery := range subQueries {
		booleanSubQueries = append(booleanSubQueries, NewBooleanSubQuery(OccurShould, subQuery))
	}
	return &BooleanQuery{
		booleanSubQueries,
	}
}

func NewMultiTermsQuery(terms []*schema.Term) reader.Query {
	termQueries := make([]reader.Query, 0, len(terms))
	for _, term := range terms {
		termQueries = append(termQueries, NewTermQuery(term))
	}
	return NewBooleanUnionQuery(termQueries)
}

func (a *BooleanQuery) Weight(searcher *reader.Searcher, scoringEnabled bool) (reader.Weight, error) {
	subWeights := make([]*BooleanSubWeight, 0, len(a.subQueries))
	for _, subQuery := range a.subQueries {
		subWeight, err := subQuery.Weight(searcher, scoringEnabled)
		if err != nil {
			return nil, err
		}
		subWeights = append(subWeights, subWeight)
	}
	return &BooleanWeight{subWeights: subWeights, scoringEnabled: scoringEnabled}, nil
}

type BooleanWeight struct {
	subWeights     []*BooleanSubWeight
	scoringEnabled bool
}

type BooleanSubWeight struct {
	occur  Occur
	weight reader.Weight
}

func (b *BooleanWeight) Scorer(segmentReader *reader.SegmentReader) (reader.Scorer, error) {
	if len(b.subWeights) == 1 {
		scorer, err := b.subWeights[0].weight.Scorer(segmentReader)
		if err != nil {
			return nil, err
		}
		return scorer, nil
	}
	booleanScorer, err := b.booleanScorerWrapper(segmentReader)
	if err != nil {
		return nil, err
	}
	return booleanScorer.Scorer(), nil
}

func (b *BooleanWeight) ForEachPruning(threshold float64, segmentReader *reader.SegmentReader, callback func(docID schema.DocID, score float64) float64) error {
	scorerWrapper, err := b.booleanScorerWrapper(segmentReader)
	if err != nil {
		return fmt.Errorf("initialize boolean score wrapper: %w", err)
	}
	if scorerWrapper.IsTermUnion() {
		// TODO: use specialized method to do WAND
		return ForEachPruning(scorerWrapper.Scorer(), threshold, callback)
	} else {
		return ForEachPruning(scorerWrapper.other, threshold, callback)
	}
}

func (b *BooleanWeight) perOccurScorers(segmentReader *reader.SegmentReader) (map[Occur][]reader.Scorer, error) {
	perOccurScorer := make(map[Occur][]reader.Scorer, len(b.subWeights))
	for _, subWeight := range b.subWeights {
		scorer, err := subWeight.weight.Scorer(segmentReader)
		if err != nil {
			return nil, err
		}
		perOccurScorer[subWeight.occur] = append(perOccurScorer[subWeight.occur], scorer)
	}
	return perOccurScorer, nil
}

func (b *BooleanWeight) booleanScorerWrapper(segmentReader *reader.SegmentReader) (*BooleanScorerWrapper, error) {
	occurScorersMap, err := b.perOccurScorers(segmentReader)
	if err != nil {
		return nil, err
	}

	var shouldScorer *BooleanScorerWrapper
	if scorers := occurScorersMap[OccurShould]; len(scorers) > 0 {
		if len(scorers) == 1 {
			shouldScorer = &BooleanScorerWrapper{other: scorers[0]}
		} else {
			termScorers := make([]*TermScorer, 0, len(scorers))
			isAllTermScorers := true
			for _, scorer := range scorers {
				if termScorer, ok := scorer.(*TermScorer); ok {
					termScorers = append(termScorers, termScorer)
				} else {
					isAllTermScorers = false
					break
				}
			}
			if isAllTermScorers {
				shouldScorer = &BooleanScorerWrapper{termUnion: termScorers}
			} else {
				shouldScorer = &BooleanScorerWrapper{other: NewUnionScorer(scorers)}
			}
		}
	}
	var mustScorer reader.Scorer
	if len(occurScorersMap[OccurMust]) > 0 {
		mustScorer = NewIntersectionScorer(occurScorersMap[OccurMust])
	}

	if mustScorer != nil && shouldScorer != nil {
		mustShouldScorer := NewMustShouldScorer(mustScorer, shouldScorer.Scorer(), NewSumWithCombiner())
		return &BooleanScorerWrapper{other: mustShouldScorer}, nil
	} else if mustScorer != nil {
		return &BooleanScorerWrapper{other: mustScorer}, nil
	} else if shouldScorer != nil {
		return shouldScorer, nil
	}
	return &BooleanScorerWrapper{other: NewEmptyScorer()}, nil
}

type BooleanScorerWrapper struct {
	termUnion []*TermScorer
	other     reader.Scorer
}

func (b *BooleanScorerWrapper) IsTermUnion() bool {
	return len(b.termUnion) > 0
}

func (b *BooleanScorerWrapper) Scorer() reader.Scorer {
	if b.IsTermUnion() {
		termScorers := make([]reader.Scorer, 0, len(b.termUnion))
		for _, termScorer := range b.termUnion {
			termScorers = append(termScorers, termScorer)
		}
		return NewUnionScorer(termScorers)
	} else {
		return b.other
	}
}
