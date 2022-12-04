package termdict

import "github.com/k-yomo/ostrich/pkg/xrange"

type TermInfo struct {
	Term          string
	DocFreq       int
	PostingsRange xrange.Range
}
