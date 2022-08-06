package index

import "github.com/k-yomo/ostrich/schema"

type Document struct {
	fields map[string]interface{}
}

type DocAddress struct {
	SegmentOrd int
	DocID      schema.DocID
}

type DocSet interface {
	Advance() (schema.DocID, error)
	Doc() (schema.DocID, error)
}
