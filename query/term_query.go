package query

import (
	"fmt"

	"github.com/k-yomo/ostrich/postings"
	"github.com/k-yomo/ostrich/reader"
	"github.com/k-yomo/ostrich/schema"
)

type TermQuery struct {
	FieldID schema.FieldID
	Term    *schema.Term
}

func NewTermQuery(term *schema.Term) reader.Query {
	return &TermQuery{
		Term: term,
	}
}

func (a *TermQuery) Weight(searcher *reader.Searcher, _ bool) (reader.Weight, error) {
	var totalDocNum uint64 = 0
	for _, segmentReader := range searcher.SegmentReaders() {
		totalDocNum += uint64(segmentReader.MaxDoc)
	}
	docFrequency, err := searcher.DocFreq(a.FieldID, a.Term.Text())
	if err != nil {
		return nil, fmt.Errorf("get doc frequency: %w", err)
	}

	return &TermWeight{
		term:             a.Term,
		similarityWeight: NewTFIDFWeight(totalDocNum, docFrequency),
	}, nil
}

type TermWeight struct {
	term             *schema.Term
	similarityWeight *TfIDFWeight
}

func (t *TermWeight) Scorer(segmentReader *reader.SegmentReader) (reader.Scorer, error) {
	invertedIndexReader := segmentReader.InvertedIndex(t.term.FieldID())
	postingsReader, err := invertedIndexReader.ReadPostings(t.term.Text())
	if err != nil {
		return nil, fmt.Errorf("read postings: %w", err)
	}
	return &TermScorer{
		postingsReader:   postingsReader,
		similarityWeight: t.similarityWeight,
	}, nil
}

func (t *TermWeight) ForEach(segmentReader *reader.SegmentReader, callback func(docID schema.DocID, score float64)) error {
	scorer, err := t.Scorer(segmentReader)
	if err != nil {
		return fmt.Errorf("initialize scorer: %w", err)
	}
	return ForEach(scorer, callback)
}

func (t *TermWeight) ForEachPruning(threshold float64, segmentReader *reader.SegmentReader, callback func(docID schema.DocID, score float64) float64) error {
	scorer, err := t.Scorer(segmentReader)
	if err != nil {
		return fmt.Errorf("open scorer: %w", err)
	}
	doc := scorer.Doc()
	for !doc.IsTerminated() {
		if score := scorer.Score(); score > threshold {
			threshold = callback(doc, score)
		}
		doc = scorer.Advance()
	}

	return nil
}

type TermScorer struct {
	postingsReader   *postings.PostingsReader
	similarityWeight *TfIDFWeight
}

func (a *TermScorer) Advance() schema.DocID {
	return a.postingsReader.Advance()
}

func (a *TermScorer) Doc() schema.DocID {
	return a.postingsReader.Doc()
}

func (a *TermScorer) Seek(target schema.DocID) schema.DocID {
	return a.postingsReader.Seek(target)
}

func (a *TermScorer) SizeHint() uint32 {
	return a.postingsReader.SizeHint()
}

func (a *TermScorer) TermFreq() uint64 {
	return a.postingsReader.TermFreq()
}

func (a *TermScorer) Score() float64 {
	return a.similarityWeight.Score(float64(a.TermFreq()))
}
