package parser

import (
	"strings"
	"testing"

	"github.com/lastvoidtemplar/json_formatter/internal/ast"
	"github.com/lastvoidtemplar/json_formatter/internal/lexer"
	"github.com/lastvoidtemplar/json_formatter/internal/token"
)

func TestParserValid1(t *testing.T) {
	input := `
[
    {
        "name": "Jason",
        "gender": "M",
        "age": 27
    },
    {
        "name": "Rosita",
        "gender": "F",
        "age": 23
    },
    {
        "name": "Leo",
        "gender": "M",
        "age": 19
    }
]
`
	lex := lexer.New(strings.NewReader(input))
	parser, err := New(lex)

	if err != nil {
		t.Fatal(err.Error())
	}

	root, parserErr := parser.Parse()

	if parserErr != nil {
		t.Fatal(parserErr.Message)
	}

	if root == nil {
		t.Fatal("Root is nil")
	}

	expeced := &ast.ArrayNode{
		Type: ast.ARRAY,
		Nodes: []ast.Node{
			&ast.ObjectNode{
				Type: ast.OBJECT,
				Nodes: []*ast.KeyValNode{
					{
						Type: ast.KEYVAL,
						Key:  token.Token{Type: token.STRING_LITERAL, Literal: "name"},
						Val: &ast.StringNode{
							Type:  ast.STRING,
							Token: token.Token{Type: token.STRING_LITERAL, Literal: "Jason"},
						},
					},
					{
						Type: ast.KEYVAL,
						Key:  token.Token{Type: token.STRING_LITERAL, Literal: "gender"},
						Val: &ast.StringNode{
							Type:  ast.STRING,
							Token: token.Token{Type: token.STRING_LITERAL, Literal: "M"},
						},
					},
					{
						Type: ast.KEYVAL,
						Key:  token.Token{Type: token.STRING_LITERAL, Literal: "age"},
						Val: &ast.NumberNode{
							Type:  ast.NUMBER,
							Token: token.Token{Type: token.NUMBER_LITERAL, Literal: "27"},
						},
					},
				},
			}, &ast.ObjectNode{
				Type: ast.OBJECT,
				Nodes: []*ast.KeyValNode{
					{
						Type: ast.KEYVAL,
						Key:  token.Token{Type: token.STRING_LITERAL, Literal: "name"},
						Val: &ast.StringNode{
							Type:  ast.STRING,
							Token: token.Token{Type: token.STRING_LITERAL, Literal: "Rosita"},
						},
					},
					{
						Type: ast.KEYVAL,
						Key:  token.Token{Type: token.STRING_LITERAL, Literal: "gender"},
						Val: &ast.StringNode{
							Type:  ast.STRING,
							Token: token.Token{Type: token.STRING_LITERAL, Literal: "F"},
						},
					},
					{
						Type: ast.KEYVAL,
						Key:  token.Token{Type: token.STRING_LITERAL, Literal: "age"},
						Val: &ast.NumberNode{
							Type:  ast.NUMBER,
							Token: token.Token{Type: token.NUMBER_LITERAL, Literal: "23"},
						},
					},
				},
			},
			&ast.ObjectNode{
				Type: ast.OBJECT,
				Nodes: []*ast.KeyValNode{
					{
						Type: ast.KEYVAL,
						Key:  token.Token{Type: token.STRING_LITERAL, Literal: "name"},
						Val: &ast.StringNode{
							Type:  ast.STRING,
							Token: token.Token{Type: token.STRING_LITERAL, Literal: "Leo"},
						},
					},
					{
						Type: ast.KEYVAL,
						Key:  token.Token{Type: token.STRING_LITERAL, Literal: "gender"},
						Val: &ast.StringNode{
							Type:  ast.STRING,
							Token: token.Token{Type: token.STRING_LITERAL, Literal: "M"},
						},
					},
					{
						Type: ast.KEYVAL,
						Key:  token.Token{Type: token.STRING_LITERAL, Literal: "age"},
						Val: &ast.NumberNode{
							Type:  ast.NUMBER,
							Token: token.Token{Type: token.NUMBER_LITERAL, Literal: "19"},
						},
					},
				},
			},
		},
	}

	compareNode(root, expeced, t)
}
func TestParserValidUTF8(t *testing.T) {
	input := `
[
    {
        "име": "Деян",
        "пол": "М",
        "възраст": 27
    },
    {
        "име": "Ивана",
        "пол": "Ж",
        "възраст": 23
    },
    {
        "име": "Димитър",
        "пол": "М",
        "възраст": 19
    }
]
`
	lex := lexer.New(strings.NewReader(input))
	parser, err := New(lex)

	if err != nil {
		t.Fatal(err.Error())
	}

	root, parserErr := parser.Parse()

	if parserErr != nil {
		t.Fatal(parserErr.Message)
	}

	if root == nil {
		t.Fatal("Root is nil")
	}

	expeced := &ast.ArrayNode{
		Type: ast.ARRAY,
		Nodes: []ast.Node{
			&ast.ObjectNode{
				Type: ast.OBJECT,
				Nodes: []*ast.KeyValNode{
					{
						Type: ast.KEYVAL,
						Key:  token.Token{Type: token.STRING_LITERAL, Literal: "име"},
						Val: &ast.StringNode{
							Type:  ast.STRING,
							Token: token.Token{Type: token.STRING_LITERAL, Literal: "Деян"},
						},
					},
					{
						Type: ast.KEYVAL,
						Key:  token.Token{Type: token.STRING_LITERAL, Literal: "пол"},
						Val: &ast.StringNode{
							Type:  ast.STRING,
							Token: token.Token{Type: token.STRING_LITERAL, Literal: "М"},
						},
					},
					{
						Type: ast.KEYVAL,
						Key:  token.Token{Type: token.STRING_LITERAL, Literal: "възраст"},
						Val: &ast.NumberNode{
							Type:  ast.NUMBER,
							Token: token.Token{Type: token.NUMBER_LITERAL, Literal: "27"},
						},
					},
				},
			}, &ast.ObjectNode{
				Type: ast.OBJECT,
				Nodes: []*ast.KeyValNode{
					{
						Type: ast.KEYVAL,
						Key:  token.Token{Type: token.STRING_LITERAL, Literal: "име"},
						Val: &ast.StringNode{
							Type:  ast.STRING,
							Token: token.Token{Type: token.STRING_LITERAL, Literal: "Ивана"},
						},
					},
					{
						Type: ast.KEYVAL,
						Key:  token.Token{Type: token.STRING_LITERAL, Literal: "пол"},
						Val: &ast.StringNode{
							Type:  ast.STRING,
							Token: token.Token{Type: token.STRING_LITERAL, Literal: "Ж"},
						},
					},
					{
						Type: ast.KEYVAL,
						Key:  token.Token{Type: token.STRING_LITERAL, Literal: "възраст"},
						Val: &ast.NumberNode{
							Type:  ast.NUMBER,
							Token: token.Token{Type: token.NUMBER_LITERAL, Literal: "23"},
						},
					},
				},
			},
			&ast.ObjectNode{
				Type: ast.OBJECT,
				Nodes: []*ast.KeyValNode{
					{
						Type: ast.KEYVAL,
						Key:  token.Token{Type: token.STRING_LITERAL, Literal: "име"},
						Val: &ast.StringNode{
							Type:  ast.STRING,
							Token: token.Token{Type: token.STRING_LITERAL, Literal: "Димитър"},
						},
					},
					{
						Type: ast.KEYVAL,
						Key:  token.Token{Type: token.STRING_LITERAL, Literal: "пол"},
						Val: &ast.StringNode{
							Type:  ast.STRING,
							Token: token.Token{Type: token.STRING_LITERAL, Literal: "М"},
						},
					},
					{
						Type: ast.KEYVAL,
						Key:  token.Token{Type: token.STRING_LITERAL, Literal: "възраст"},
						Val: &ast.NumberNode{
							Type:  ast.NUMBER,
							Token: token.Token{Type: token.NUMBER_LITERAL, Literal: "19"},
						},
					},
				},
			},
		},
	}

	compareNode(root, expeced, t)
}
func compareNode(actual ast.Node, expected ast.Node, t *testing.T) {
	if actual.NodeType() != expected.NodeType() {
		t.Errorf("Exptected type was %T, but got %T", expected, actual)
		return
	}

	switch node := actual.(type) {
	case ast.LeafNode:
		exp, _ := expected.(ast.LeafNode)
		if node.Literal() != exp.Literal() {
			t.Errorf("Expected literal was %s, but got %s", exp.Literal(), node.Literal())
			return
		}
	case *ast.KeyValNode:
		exp, _ := expected.(*ast.KeyValNode)
		if exp.Key.Literal != node.Key.Literal {
			t.Errorf("Expected key literal was %s, but got %s", exp.Key.Literal, node.Key.Literal)
		}
		compareNode(node.Val, exp.Val, t)
	case *ast.ArrayNode:
		exp, _ := expected.(*ast.ArrayNode)
		if len(node.Nodes) != len(exp.Nodes) {
			t.Errorf("Expected number of elements was %d, but got %d", len(exp.Nodes), len(node.Nodes))
			return
		}
		for i := 0; i < len(node.Nodes); i++ {
			compareNode(node.Nodes[i], exp.Nodes[i], t)
		}
	case *ast.ObjectNode:
		exp, _ := expected.(*ast.ObjectNode)
		if len(node.Nodes) != len(exp.Nodes) {
			t.Errorf("Expected number of elements was %d, but got %d", len(exp.Nodes), len(node.Nodes))
			return
		}
		for i := 0; i < len(node.Nodes); i++ {
			compareNode(node.Nodes[i], exp.Nodes[i], t)
		}
	}
}
