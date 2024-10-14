package ast

import (
	"github.com/lastvoidtemplar/json_formatter/internal/token"
)

type Node interface {
	String(tabs string) string
}

type LeafNode interface {
	Node
	leatNode()
}

type NullNode struct {
	Token token.Token
}

func (node *NullNode) String(tabs string) string {
	return tabs + node.Token.Literal
}

func (node *NullNode) leatNode() {}

type BoolNode struct {
	Token token.Token
}

func (node *BoolNode) String(tabs string) string {
	return tabs + node.Token.Literal
}

func (node *BoolNode) leatNode() {}

type NumberNode struct {
	Token token.Token
}

func (node *NumberNode) String(tabs string) string {
	return tabs + node.Token.Literal
}

func (node *NumberNode) leatNode() {}

type StringNode struct {
	Token token.Token
}

func (node *StringNode) String(tabs string) string {
	return tabs + node.Token.Literal
}

func (node *StringNode) leatNode() {}

type UndefinedNode struct {
	Token token.Token
}

func (node *UndefinedNode) String(tabs string) string {
	return tabs + node.Token.Literal
}

func (node *UndefinedNode) leatNode() {}

func NewLeafNode(tok token.Token) LeafNode {
	switch tok.Type {
	case token.NULL:
		return &NullNode{Token: tok}
	case token.TRUE, token.FALSE:
		return &BoolNode{Token: tok}
	case token.NUMBER_LITERAL:
		return &NumberNode{Token: tok}
	case token.STRING_LITERAL:
		return &StringNode{Token: tok}
	case token.UNDEFINED:
		return &UndefinedNode{Token: tok}
	default:
		return nil
	}
}

type ArrayNode struct {
	Nodes []Node
}

func (array *ArrayNode) String(tabs string) string {
	return "[]"
}

type KeyValNode struct {
	Key token.Token
	Val Node
}

func (node *KeyValNode) String(tabs string) string {
	return ":"
}

func NewKeyVal(key token.Token, val Node) *KeyValNode {
	if val == nil {
		return nil
	}

	if key.Type != token.STRING_LITERAL && key.Type != token.UNDEFINED {
		return nil
	}

	return &KeyValNode{
		Key: key,
		Val: val,
	}
}

type ObjectNode struct {
	Nodes []Node
}

func (array *ObjectNode) String(tabs string) string {
	return "{}"
}
