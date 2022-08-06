package analyzer

func NewEnglishAnalyzer() *Analyzer {
	analyzer := NewAnalyzer(&SpaceTokenizer{})
	analyzer.SetCharFilter(&LowerCaseCharFilter{})
	analyzer.SetTokenFilter(&StemmingTokenFilter{}, &StopWordTokenFilter{})
	return analyzer
}
