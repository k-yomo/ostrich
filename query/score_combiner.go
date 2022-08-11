package query

import "github.com/k-yomo/ostrich/reader"

type ScoreCombiner interface {
	Update(scorer reader.Scorer)
	Clear()
	Score() float64
}

type noopCombiner struct {
	score float64
}

func NewNoopCombiner() ScoreCombiner {
	return &noopCombiner{}
}

func (n *noopCombiner) Update(scorer reader.Scorer) {
	n.score += scorer.Score()
}

func (n *noopCombiner) Clear() {
	n.score = 0
}

func (n *noopCombiner) Score() float64 {
	return n.score
}

type sumWithCombiner struct {
	score float64
}

func NewSumWithCombiner() ScoreCombiner {
	return &sumWithCombiner{}
}

func (n *sumWithCombiner) Update(scorer reader.Scorer) {
	n.score += scorer.Score()
}

func (n *sumWithCombiner) Clear() {
	n.score = 0
}

func (n *sumWithCombiner) Score() float64 {
	return n.score
}
