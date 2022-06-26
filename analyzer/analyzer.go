package analyzer

type Analyzer interface {
	Analyze(text string) []string
}
