package query

import (
	"github.com/k-yomo/ostrich/index"
	"github.com/k-yomo/ostrich/schema"
)

type Parser struct {
	schema        *schema.Schema
	defaultFields []schema.FieldID
}

func NewParser(idx *index.Index, defaultFields []schema.FieldID) *Parser {
	return &Parser{
		schema:        idx.Schema(),
		defaultFields: defaultFields,
	}
}

func (p *Parser) Parse(query string) index.Query {
	// TODO: implement parser
	return &AllQuery{}
}
