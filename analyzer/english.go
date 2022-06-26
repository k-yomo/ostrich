package analyzer

import (
	"github.com/kljensen/snowball/english"
	"strings"
)

type EnglishAnalyzer struct{}

func (a *EnglishAnalyzer) Analyze(text string) []string {
	tokens := tokenize(text)
	tokens = lowercaseFilter(tokens)
	tokens = stopWordFilter(tokens)
	tokens = stemmerFilter(tokens)

	return tokens
}

func tokenize(text string) []string {
	return strings.Split(text, " ")
}

func lowercaseFilter(tokens []string) []string {
	r := make([]string, 0, len(tokens))
	for _, token := range tokens {
		r = append(r, strings.ToLower(token))
	}
	return r
}

var stopWords = map[string]struct{}{
	"a": {}, "and": {}, "be": {}, "have": {}, "i": {},
	"in": {}, "of": {}, "that": {}, "the": {}, "to": {},
}

func stopWordFilter(tokens []string) []string {
	r := make([]string, 0, len(tokens))
	for _, token := range tokens {
		if _, ok := stopWords[token]; !ok {
			r = append(r, token)
		}
	}
	return r
}

func stemmerFilter(tokens []string) []string {
	r := make([]string, 0, len(tokens))
	for _, token := range tokens {
		r = append(r, english.Stem(token, false))
	}
	return r
}
