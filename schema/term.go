package schema

type Term struct {
	fieldID   FieldID
	fieldType FieldType
	data      []byte
}

func NewTermFromText(fieldID FieldID, text string) *Term {
	return &Term{
		fieldID:   fieldID,
		fieldType: FieldTypeText,
		data:      []byte(text),
	}
}

func (t *Term) FieldID() FieldID {
	return t.fieldID
}

func (t *Term) Text() string {
	return string(t.data)
}
