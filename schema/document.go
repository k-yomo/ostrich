package schema

import "sort"

type Document struct {
	FieldValues []*FieldValue
}

func (d *Document) SortedFieldValues() []*FieldAndFieldValues {
	if len(d.FieldValues) == 0 {
		return nil
	}
	fieldValues := append(d.FieldValues)
	sort.Slice(fieldValues, func(i, j int) bool {
		return fieldValues[i].FieldID < fieldValues[j].FieldID
	})

	var fieldAndFieldValues []*FieldAndFieldValues
	fieldAndFieldValues = append(fieldAndFieldValues, &FieldAndFieldValues{
		Field:       fieldValues[0].FieldID,
		FieldValues: []*FieldValue{fieldValues[0]},
	})
	for _, fieldValue := range fieldValues[1:] {
		lastField := fieldAndFieldValues[len(fieldAndFieldValues)-1]
		if fieldValue.FieldID == lastField.Field {
			lastField.FieldValues = append(lastField.FieldValues, fieldValue)
		} else {
			fieldAndFieldValues = append(fieldAndFieldValues, &FieldAndFieldValues{
				Field:       fieldValue.FieldID,
				FieldValues: []*FieldValue{fieldValue},
			})
		}
	}
	return fieldAndFieldValues
}
