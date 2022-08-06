package termdict

type Range struct {
	From uint64
	To   uint64
}

func (r Range) Len() uint64 {
	return r.To - r.From
}

type TermInfo struct {
	Term          string
	DocFreq       int
	PostingsRange Range
}
