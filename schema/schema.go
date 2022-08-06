package schema

import "sync"

type Schema struct {
	mu sync.RWMutex

	Fields   []*FieldEntry `json:"fields"`
	fieldMap map[string]*FieldEntry
}

func (s *Schema) FieldEntry(fieldID FieldID) *FieldEntry {
	return s.Fields[fieldID]
}
