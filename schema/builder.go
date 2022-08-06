package schema

import "github.com/k-yomo/ostrich/analyzer"

type Builder struct {
	Fields   []*FieldEntry `json:"fields"`
	FieldMap map[string]*FieldEntry
}

func NewBuilder() *Builder {
	return &Builder{
		FieldMap: make(map[string]*FieldEntry),
	}
}

func (b *Builder) AddTextField(fieldName string, analyzer *analyzer.Analyzer) FieldID {
	field := FieldID(len(b.Fields))
	fieldEntry := newFieldEntry(field, fieldName, FieldTypeText, analyzer)
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
