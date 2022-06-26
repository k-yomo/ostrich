package schema

type FieldEntry struct {
	ID   FieldID `json:"id"`
	Name string  `json:"name"`
	FieldType FieldType `json:"fieldType"`
}

func newFieldEntry(id FieldID, name string, fieldType FieldType) *FieldEntry {
	return &FieldEntry{
		ID:        id,
		Name:      name,
		FieldType: fieldType,
	}
}
