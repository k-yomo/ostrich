package query

import "github.com/k-yomo/ostrich/schema"

type NodeKind int

const (
	NodeKindLeaf NodeKind = iota
	NodeKindAnd
	NodeKindOr
)

type ASTNode struct {
	Kind  NodeKind
	Left  *ASTNode
	Right *ASTNode
	// Value is the literal value
	// For now, it's always []*schema.Term
	Value interface{}
}

func NewLogicalOperationNode(kind NodeKind, left *ASTNode, right *ASTNode) *ASTNode {
	return &ASTNode{
		Kind:  kind,
		Left:  left,
		Right: right,
	}
}

func NewTermsNode(terms []*schema.Term) *ASTNode {
	return &ASTNode{
		Kind:  NodeKindLeaf,
		Value: terms,
	}
}
