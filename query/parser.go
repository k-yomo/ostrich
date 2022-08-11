package query

import (
	"errors"
	"fmt"
	"strings"

	"github.com/k-yomo/ostrich/reader"
	"github.com/k-yomo/ostrich/schema"
)

type Parser struct {
	schema        *schema.Schema
	fieldNameMap  map[string]schema.FieldID
	defaultFields []schema.FieldID
}

func NewParser(indexSchema *schema.Schema, defaultFields []schema.FieldID) *Parser {
	fieldNameMap := map[string]schema.FieldID{}
	for _, field := range indexSchema.Fields {
		fieldNameMap[field.Name] = field.ID
	}
	return &Parser{
		schema:        indexSchema,
		fieldNameMap:  fieldNameMap,
		defaultFields: defaultFields,
	}
}

func (p *Parser) Parse(query string) (reader.Query, error) {
	tokens := tokenize(query)
	if len(tokens) == 0 {
		return &AllQuery{}, nil
	}
	astNode, rest, err := p.expr(tokens)
	if err != nil {
		return nil, err
	}
	if len(rest) > 0 {
		return nil, fmt.Errorf("there are unparsed tokens: %v", rest)
	}
	return astToQuery(astNode), nil
}

// expr = primary (primary | "AND" primary | "OR" primary)*
func (p *Parser) expr(tokens []string) (*ASTNode, []string, error) {
	if tokens[0] == ")" {
		return nil, tokens, nil
	}
	node, next, err := p.primary(tokens)
	if err != nil {
		return nil, nil, err
	}
	for {
		if len(next) == 0 {
			return node, nil, nil
		}
		switch strings.ToLower(next[0]) {
		case ")":
			return node, next, nil
		case "and":
			right, nextTokens, err := p.primary(next[1:])
			if err != nil {
				return nil, nil, err
			}
			node = NewLogicalOperationNode(NodeKindAnd, node, right)
			next = nextTokens
		case "or":
			right, nextTokens, err := p.primary(next[1:])
			if err != nil {
				return nil, nil, err
			}
			node = NewLogicalOperationNode(NodeKindOr, node, right)
			next = nextTokens
		default:
			right, nextTokens, err := p.primary(next)
			if err != nil {
				return nil, nil, err
			}
			node = NewLogicalOperationNode(NodeKindOr, node, right)
			next = nextTokens
		}
	}
}

// primary = value | "(" expr ")"
func (p *Parser) primary(tokens []string) (*ASTNode, []string, error) {
	if tokens[0] == "(" {
		node, rest, err := p.expr(tokens[1:])
		if err != nil {
			return nil, nil, err
		}
		if len(rest) == 0 || rest[0] != ")" {
			return nil, nil, errors.New("no closing parenthesis")
		}
		return node, rest[1:], nil
	}
	if fieldTerm := strings.SplitN(tokens[0], ":", 2); len(fieldTerm) == 2 {
		fieldName := fieldTerm[0]
		term := fieldTerm[1]
		if fieldID, ok := p.fieldNameMap[fieldName]; ok {
			return NewTermNode(schema.NewTermFromText(fieldID, term)), tokens[1:], nil
		}
	}
	terms := make([]*schema.Term, 0, len(p.defaultFields))
	for _, fieldID := range p.defaultFields {
		terms = append(terms, schema.NewTermFromText(fieldID, tokens[0]))
	}
	return NewTermsNode(terms), tokens[1:], nil
}

func tokenize(query string) []string {
	var tokens []string
	var curToken []rune
	for _, c := range query {
		switch c {
		case ' ':
			if len(curToken) > 0 {
				tokens = append(tokens, string(curToken))
				curToken = nil
			}
		case '(', ')':
			if len(curToken) > 0 {
				tokens = append(tokens, string(curToken))
				curToken = nil
			}
			tokens = append(tokens, string(c))
		default:
			curToken = append(curToken, c)
		}
	}
	if len(curToken) > 0 {
		tokens = append(tokens, string(curToken))
	}
	return tokens
}

func astToQuery(node *ASTNode) reader.Query {
	switch node.Kind {
	case NodeKindAnd:
		return NewBooleanIntersectionQuery([]reader.Query{astToQuery(node.Left), astToQuery(node.Right)})
	case NodeKindOr:
		return NewBooleanUnionQuery([]reader.Query{astToQuery(node.Left), astToQuery(node.Right)})
	default:
		switch v := node.Value.(type) {
		case *schema.Term:
			return NewTermQuery(v)
		case []*schema.Term:
			return NewMultiTermsQuery(v)
		default:
			panic(fmt.Sprintf("unexpected value: %+v", v))
		}
	}
}
