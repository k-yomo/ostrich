package termdict

type Range struct {
	From int
	To   int
}

func (r Range) Len() int {
	return r.To - r.From
}

type TermInfo struct {
	Term          string
	DocFreq       int
	PostingsRange Range
}
