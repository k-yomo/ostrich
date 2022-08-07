package query

import "math"

type TfIDFWeight struct {
	idf float64
}

func NewTFIDFWeight(totalDocNum uint64, documentFrequency int) *TfIDFWeight {
	df := 1.0 + float64(documentFrequency)
	return &TfIDFWeight{
		idf: 1.0 + math.Log2(float64(totalDocNum)/df),
	}
}

func (t *TfIDFWeight) Score(termFrequency float64) float64 {
	return termFrequency * t.idf
}
