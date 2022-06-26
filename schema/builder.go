package schema

type Builder struct {
	Fields   []*FieldEntry `json:"fields"`
	FieldMap map[string]*FieldEntry
}

func NewBuilder() *Builder {
	return &Builder{
		FieldMap: make(map[string]*FieldEntry),
	}
}

func (b *Builder) AddTextField(fieldName string) FieldID {
	field := FieldID(len(b.Fields))
	fieldEntry := newFieldEntry(field, fieldName, FieldTypeText)
	b.Fields = append(b.Fields, fieldEntry)
	b.FieldMap[fieldName] = fieldEntry
	return field
}

func (b *Builder) Build() *Schema {
	return &Schema{
		Fields:   b.Fields,
		fieldMap: b.FieldMap,
	}
}
