package analyzer

func NewEnglishAnalyzer() *Analyzer {
	analyzer := NewAnalyzer(&SimpleTokenizer{})
	analyzer.SetCharFilter(&LowerCaseCharFilter{})
	analyzer.SetTokenFilter(&StemmingTokenFilter{}, &StopWordTokenFilter{})
	return analyzer
}
