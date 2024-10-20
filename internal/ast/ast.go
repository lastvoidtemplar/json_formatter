package ast

import (
	"log"

	"github.com/lastvoidtemplar/json_formatter/internal/token"
)

type NodeType byte

const (
	NULL NodeType = iota
	BOOL
	NUMBER
	STRING
	UNDEFINED

	ARRAY

	KEYVAL
	OBJECT
)

type Node interface {
	NodeType() NodeType
}

type LeafNode interface {
	Node
	leatNode()
}

type NullNode struct {
	Type  NodeType
	Token token.Token
}

func (node *NullNode) NodeType() NodeType {
	return node.Type
}

func (node *NullNode) leatNode() {}

type BoolNode struct {
	Type  NodeType
	Token token.Token
}

func (node *BoolNode) NodeType() NodeType {
	return node.Type
}

func (node *BoolNode) leatNode() {}

type NumberNode struct {
	Type  NodeType
	Token token.Token
}

func (node *NumberNode) NodeType() NodeType {
	return node.Type
}

func (node *NumberNode) leatNode() {}

type StringNode struct {
	Type  NodeType
	Token token.Token
}

func (node *StringNode) NodeType() NodeType {
	return node.Type
}

func (node *StringNode) leatNode() {}

type UndefinedNode struct {
	Type  NodeType
	Token token.Token
}

func (node *UndefinedNode) NodeType() NodeType {
	return node.Type
}

func (node *UndefinedNode) leatNode() {}

func NewLeafNode(tok token.Token) LeafNode {
	switch tok.Type {
	case token.NULL:
		return &NullNode{Type: NULL, Token: tok}
	case token.TRUE, token.FALSE:
		return &BoolNode{Type: BOOL, Token: tok}
	case token.NUMBER_LITERAL:
		return &NumberNode{Type: NUMBER, Token: tok}
	case token.STRING_LITERAL:
		return &StringNode{Type: STRING, Token: tok}
	default:
		return nil
	}
}

type ArrayNode struct {
	Type  NodeType
	Nodes []Node
}

func (array *ArrayNode) NodeType() NodeType {
	return array.Type
}

func (array *ArrayNode) Add(node Node) {
	if node == nil {
		log.Println("Tried to add nil node to array")
		return
	}

	array.Nodes = append(array.Nodes, node)
}

func NewArrayNode() *ArrayNode {
	return &ArrayNode{
		Type:  ARRAY,
		Nodes: make([]Node, 0),
	}
}

type KeyValNode struct {
	Type NodeType
	Key  token.Token
	Val  Node
}

func (node *KeyValNode) NodeType() NodeType {
	return node.Type
}

func NewKeyVal(key token.Token, val Node) *KeyValNode {
	if val == nil {
		return nil
	}

	return &KeyValNode{
		Type: KEYVAL,
		Key:  key,
		Val:  val,
	}
}

type ObjectNode struct {
	Type  NodeType
	Nodes []Node
	keys  map[string]struct{}
}

func (array *ObjectNode) NodeType() NodeType {
	return array.Type
}

func NewObjectNode() *ObjectNode {
	return &ObjectNode{
		Type:  OBJECT,
		Nodes: make([]Node, 0),
		keys:  make(map[string]struct{}),
	}
}

func (object *ObjectNode) Add(keyval *KeyValNode) bool {
	if keyval == nil {
		log.Println("Tried to add nil keylval")
		return false
	}

	if _, ok := object.keys[keyval.Key.Literal]; !ok {
		return false
	}

	object.Nodes = append(object.Nodes, keyval)
	object.keys[keyval.Key.Literal] = struct{}{}
	return true
}
