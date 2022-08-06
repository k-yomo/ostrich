package analyzer

type Analyzer struct {
	charFilters  []CharFilter
	tokenizer    Tokenizer
	tokenFilters []TokenFilter
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
	return &Analyzer{tokenizer: tokenizer}
}

func (a *Analyzer) SetCharFilter(charFilters ...CharFilter) {
	a.charFilters = append(a.charFilters, charFilters...)
}

func (a *Analyzer) SetTokenFilter(tokenFilters ...TokenFilter) {
	a.tokenFilters = append(a.tokenFilters, tokenFilters...)
}

func (a *Analyzer) Analyze(text string) []string {
	for _, charFilter := range a.charFilters {
		text = charFilter.Filter(text)
	}
	tokens := a.tokenizer.Tokenize(text)
	for _, tokenFilter := range a.tokenFilters {
		tokens = tokenFilter.Filter(tokens)
	}
	return tokens
}
