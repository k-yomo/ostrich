package schema

type FieldEntry struct {
	ID           FieldID   `json:"id"`
	Name         string    `json:"name"`
	FieldType    FieldType `json:"fieldType"`
	AnalyzerName string    `json:"analyzer"`
}

func newFieldEntry(id FieldID, name string, fieldType FieldType, analyzerName string) *FieldEntry {
	return &FieldEntry{
		ID:           id,
		Name:         name,
		FieldType:    fieldType,
		AnalyzerName: analyzerName,
	}
}
