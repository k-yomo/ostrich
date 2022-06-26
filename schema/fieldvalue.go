package schema

type FieldID uint32

type FieldAndFieldValues struct {
	Field       FieldID
	FieldValues []*FieldValue
}

type FieldValue struct {
	FieldID FieldID
	Value   Value
}
