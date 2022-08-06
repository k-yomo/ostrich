package analyzer

import "github.com/kljensen/snowball/english"

type StemmingTokenFilter struct{}

var _ TokenFilter = &StemmingTokenFilter{}

func (s *StemmingTokenFilter) Filter(tokens []string) []string {
	r := make([]string, 0, len(tokens))
	for _, token := range tokens {
		r = append(r, english.Stem(token, false))
	}
	return r
}
