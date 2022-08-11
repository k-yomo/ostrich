package query

import (
	"github.com/k-yomo/ostrich/reader"
	"github.com/k-yomo/ostrich/schema"
	"sort"
)

type InterSectionScorer struct {
	left   reader.Scorer
	right  reader.Scorer
	others []reader.Scorer
}

func NewIntersectionScorer(scorers []reader.Scorer) reader.Scorer {
	if len(scorers) == 0 {
		return NewEmptyScorer()
	}
	if len(scorers) == 1 {
		return scorers[0]
	}

	//  start from the smallest scorer to reduce seek count
	sort.Slice(scorers, func(i, j int) bool {
		return scorers[i].SizeHint() < scorers[j].SizeHint()
	})

	advanceToFirstDoc(scorers)

	return &InterSectionScorer{
		left:   scorers[0],
		right:  scorers[1],
		others: scorers[2:],
	}
}

func (i *InterSectionScorer) Advance() schema.DocID {
	candidate := i.left.Advance()
outer:
	for {
		rightDoc := i.right.Seek(candidate)
		candidate = i.left.Seek(rightDoc)
		if candidate == rightDoc {
			break
		}
	}
	for _, scorer := range i.others {
		seekDoc := scorer.Seek(candidate)
		if seekDoc > candidate {
			candidate = i.left.Seek(seekDoc)
			goto outer
		}
	}
	return candidate
}

func (i *InterSectionScorer) Doc() schema.DocID {
	return i.left.Doc()
}

func (i *InterSectionScorer) Seek(target schema.DocID) schema.DocID {
	i.left.Seek(target)
	return advanceToFirstDoc(append([]reader.Scorer{i.left, i.right}, i.others...))
}

func (i *InterSectionScorer) SizeHint() uint32 {
	return i.left.SizeHint()
}

func (i *InterSectionScorer) Score() float64 {
	score := i.left.Score() + i.right.Score()
	for _, other := range i.others {
		score += other.Score()
	}
	return score
}

func advanceToFirstDoc(scorers []reader.Scorer) schema.DocID {
	var curDocID schema.DocID
outer:
	for _, docSet := range scorers {
		docID := docSet.Seek(curDocID)
		if docID > curDocID {
			curDocID = docID
			goto outer
		}
	}
	return curDocID
}
