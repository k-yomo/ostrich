package opstamp

import (
	"sync"
)

type OpStamp uint64

type Stamper struct {
	mu      sync.Mutex
	opStamp OpStamp
}

func NewStamper(firstOpstamp OpStamp) *Stamper {
	return &Stamper{opStamp: firstOpstamp}
}

func (s *Stamper) Stamp() OpStamp {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.opStamp += 1
	return s.opStamp
}
