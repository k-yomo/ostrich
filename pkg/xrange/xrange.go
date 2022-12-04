package xrange

type Range struct {
	From int
	To   int
}

func (r Range) Len() int {
	return r.To - r.From + 1
}
