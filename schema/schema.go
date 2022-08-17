package schema

import "sync"

type Schema struct {
	mu sync.RWMutex

	Fields []*FieldEntry `json:"fields"`
}

func NewSchema() *Schema {
	return &Schema{}
}

func (s *Schema) AddTextField(fieldName string, analyzerName string) FieldID {
	field := FieldID(len(s.Fields))
	fieldEntry := newFieldEntry(field, fieldName, FieldTypeText, analyzerName)
	s.Fields = append(s.Fields, fieldEntry)
	return field
}

func (s *Schema) FieldEntry(fieldID FieldID) *FieldEntry {
	return s.Fields[fieldID]
}

func (s *Schema) FieldByName(name string) *FieldEntry {
	for _, field := range s.Fields {
		if field.Name == name {
			return field
		}
	}
	return nil
}

func (s *Schema) FieldIDs() []FieldID {
	fieldIDs := make([]FieldID, 0, len(s.Fields))
	for _, field := range s.Fields {
		fieldIDs = append(fieldIDs, field.ID)
	}
	return fieldIDs
}
