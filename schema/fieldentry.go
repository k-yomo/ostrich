package schema

import "github.com/k-yomo/ostrich/analyzer"

type FieldEntry struct {
	ID        FieldID   `json:"id"`
	Name      string    `json:"name"`
	FieldType FieldType `json:"fieldType"`
	Analyzer  *analyzer.Analyzer
}

func newFieldEntry(id FieldID, name string, fieldType FieldType, analyzer *analyzer.Analyzer) *FieldEntry {
	return &FieldEntry{
		ID:        id,
		Name:      name,
		FieldType: fieldType,
		Analyzer:  analyzer,
	}
}
