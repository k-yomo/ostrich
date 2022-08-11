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

func (s *Schema) FieldIDs() []FieldID {
	fieldIDs := make([]FieldID, 0, len(s.Fields))
	for _, field := range s.Fields {
		fieldIDs = append(fieldIDs, field.ID)
	}
	return fieldIDs
}
