package analyzer

import (
	"strings"
)

type LowerCaseCharFilter struct{}

var _ CharFilter = &LowerCaseCharFilter{}

func (s *LowerCaseCharFilter) Filter(text string) string {
	return strings.ToLower(text)
}
