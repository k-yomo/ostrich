package analyzer

import "sync"

const DefaultAnalyzerName = "default"

var (
	analyzerMu  sync.RWMutex
	analyzerMap = map[string]*Analyzer{
		DefaultAnalyzerName: {
			CharFilters: []CharFilter{&LowerCaseCharFilter{}},
			Tokenizer:   &SimpleTokenizer{},
		},
	}
)

func Register(name string, analyzer *Analyzer) {
	analyzerMu.Lock()
	defer analyzerMu.Unlock()
	if analyzer == nil {
		panic("register analyzer is nil")
	}
	if _, dup := analyzerMap[name]; dup {
		panic("Register called twice for analyzer " + name)
	}
	analyzerMap[name] = analyzer
}

func Get(name string) (*Analyzer, bool) {
	a, ok := analyzerMap[name]
	return a, ok
}

type Analyzer struct {
	CharFilters  []CharFilter
	Tokenizer    Tokenizer
	TokenFilters []TokenFilter
}

type CharFilter interface {
	Filter(text string) string
}

type Tokenizer interface {
	Tokenize(text string) []string
}

type TokenFilter interface {
	Filter(tokens []string) []string
}

func NewAnalyzer(tokenizer Tokenizer) *Analyzer {
	return &Analyzer{Tokenizer: tokenizer}
}

func (a *Analyzer) SetCharFilter(charFilters ...CharFilter) {
	a.CharFilters = append(a.CharFilters, charFilters...)
}

func (a *Analyzer) SetTokenFilter(tokenFilters ...TokenFilter) {
	a.TokenFilters = append(a.TokenFilters, tokenFilters...)
}

func (a *Analyzer) Analyze(text string) []string {
	for _, charFilter := range a.CharFilters {
		text = charFilter.Filter(text)
	}
	tokens := a.Tokenizer.Tokenize(text)
	for _, tokenFilter := range a.TokenFilters {
		tokens = tokenFilter.Filter(tokens)
	}
	return tokens
}
