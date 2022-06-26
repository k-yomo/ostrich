package schema

import "sync"

type Schema struct {
	mu sync.RWMutex

	Fields   []*FieldEntry `json:"fields"`
	fieldMap map[string]*FieldEntry
}
