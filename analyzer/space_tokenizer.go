package analyzer

import "strings"

type SpaceTokenizer struct{}

var _ Tokenizer = &SpaceTokenizer{}

func (s *SpaceTokenizer) Tokenize(text string) []string {
	return strings.Fields(text)
}
