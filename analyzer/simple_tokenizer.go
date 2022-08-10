package analyzer

import (
	"strings"
	"unicode"
)

type SimpleTokenizer struct{}

var _ Tokenizer = &SimpleTokenizer{}

func (s *SimpleTokenizer) Tokenize(text string) []string {
	return strings.FieldsFunc(text, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})
}
