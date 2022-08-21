package query

import "math"

type TfIDFWeight struct {
	idf float64
}

func NewTFIDFWeight(totalDocNum uint64, documentFrequency int) *TfIDFWeight {
	// using inverse document frequency smooth
	// https://en.wikipedia.org/wiki/Tf%E2%80%93idf#Inverse_document_frequency_2
	return &TfIDFWeight{
		idf: 1.0 + math.Log(float64(totalDocNum)/float64(1+documentFrequency)),
	}
}

func (t *TfIDFWeight) Score(termFrequency float64) float64 {
	return termFrequency * t.idf
}

func (t *TfIDFWeight) MaxScore() float64 {
	return t.Score(2_013_265_944)
}
