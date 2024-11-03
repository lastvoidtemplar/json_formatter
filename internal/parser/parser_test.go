package parser

import (
	"errors"
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
		t.Fatal(parserErr.Error())
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

func TestParserValid2(t *testing.T) {
	input := `{
		"id" : 11,
		"name": "cabal",
		"available" : true,
		"orders": null
	}`
	lex := lexer.New(strings.NewReader(input))
	parser, err := New(lex)

	if err != nil {
		t.Fatal(err.Error())
	}

	root, parserErr := parser.Parse()

	if parserErr != nil {
		t.Fatal(parserErr.Error())
	}

	if root == nil {
		t.Fatal("Root is nil")
	}

	expeced := &ast.ObjectNode{
		Type: ast.OBJECT,
		Nodes: []*ast.KeyValNode{
			{
				Type: ast.KEYVAL,
				Key:  token.Token{Type: token.STRING_LITERAL, Literal: "id"},
				Val: &ast.NumberNode{
					Type:  ast.NUMBER,
					Token: token.Token{Type: token.NUMBER_LITERAL, Literal: "11"},
				},
			},
			{
				Type: ast.KEYVAL,
				Key:  token.Token{Type: token.STRING_LITERAL, Literal: "name"},
				Val: &ast.NumberNode{
					Type:  ast.STRING,
					Token: token.Token{Type: token.STRING_LITERAL, Literal: "cabal"},
				},
			},
			{
				Type: ast.KEYVAL,
				Key:  token.Token{Type: token.STRING_LITERAL, Literal: "available"},
				Val: &ast.NumberNode{
					Type:  ast.BOOL,
					Token: token.Token{Type: token.TRUE, Literal: "true"},
				},
			},
			{
				Type: ast.KEYVAL,
				Key:  token.Token{Type: token.STRING_LITERAL, Literal: "orders"},
				Val: &ast.NumberNode{
					Type:  ast.NULL,
					Token: token.Token{Type: token.NULL, Literal: "null"},
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
		t.Fatal(parserErr.Error())
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

func TestParserEmptyLexer(t *testing.T) {
	input := ""
	lex := lexer.New(strings.NewReader(input))
	_, err := New(lex)

	if err == nil {
		t.Fatal("Expected error, but got nil")
	}

	if err != ErrEmptyLexer {
		t.Fatalf("Expected ErrEmptyLexer, but got %s", err.Error())
	}
}

func TestParserExtraTokens(t *testing.T) {
	input := `{
	
	},`
	lex := lexer.New(strings.NewReader(input))
	parser, err := New(lex)

	if err != nil {
		t.Fatal(err.Error())
	}

	_, err = parser.Parse()

	if err == nil {
		t.Fatal("Expected error, but got nil")
	}

	if !errors.Is(err, ErrExtraTokens) {
		t.Fatalf("Expected ErrorExtraTokens, but got %s", err.Error())
	}
}

func TestParserUndefinedToken(t *testing.T) {
	input := `A3`
	lex := lexer.New(strings.NewReader(input))
	parser, err := New(lex)

	if err != nil {
		t.Fatal(err.Error())
	}

	_, err = parser.Parse()

	if err == nil {
		t.Fatal("Expected error, but got nil")
	}

	if !errors.Is(err, ErrUndefinedToken) {
		t.Fatalf("Expected ErrUndefinedToken, but got %s", err.Error())
	}
}

func TestParserMissingArrayClosingBracket1(t *testing.T) {
	input := `[
		true,
		"hello",
		null
		`
	lex := lexer.New(strings.NewReader(input))
	parser, err := New(lex)

	if err != nil {
		t.Fatal(err.Error())
	}

	_, err = parser.Parse()

	if err == nil {
		t.Fatal("Expected error, but got nil")
	}

	if !errors.Is(err, ErrMissingArrayClosingBracket) {
		t.Fatalf("Expected ErrMissingArrayClosingBracket, but got %s", err.Error())
	}
}

func TestParserMissingArrayClosingBracket2(t *testing.T) {
	input := `[
		`
	lex := lexer.New(strings.NewReader(input))
	parser, err := New(lex)

	if err != nil {
		t.Fatal(err.Error())
	}

	_, err = parser.Parse()

	if err == nil {
		t.Fatal("Expected error, but got nil")
	}

	if !errors.Is(err, ErrMissingArrayClosingBracket) {
		t.Fatalf("Expected ErrMissingArrayClosingBracket, but got %s", err.Error())
	}
}
func TestParserMissingArraySeparator(t *testing.T) {
	input := `[
		true,
		"hello"
		null
		]`
	lex := lexer.New(strings.NewReader(input))
	parser, err := New(lex)

	if err != nil {
		t.Fatal(err.Error())
	}

	_, err = parser.Parse()

	if err == nil {
		t.Fatal("Expected error, but got nil")
	}

	if !errors.Is(err, ErrMissingArraySeparator) {
		t.Fatalf("Expected ErrMissingArrayClosingBracket, but got %s", err.Error())
	}
}

func TestParserKeyvalKeysNotString(t *testing.T) {
	input := `{
		123: true
	}`
	lex := lexer.New(strings.NewReader(input))
	parser, err := New(lex)

	if err != nil {
		t.Fatal(err.Error())
	}

	_, err = parser.Parse()

	if err == nil {
		t.Fatal("Expected error, but got nil")
	}

	if !errors.Is(err, ErrKeyNotString) {
		t.Fatalf("Expected ErrKeyNotString, but got %s", err.Error())
	}
}

func TestParserKeyvalMissingSeparator(t *testing.T) {
	input := `{
		"key" "val"
	}`
	lex := lexer.New(strings.NewReader(input))
	parser, err := New(lex)

	if err != nil {
		t.Fatal(err.Error())
	}

	_, err = parser.Parse()

	if err == nil {
		t.Fatal("Expected error, but got nil")
	}

	if !errors.Is(err, ErrMissingKeyvalSeparator) {
		t.Fatalf("Expected ErrMissingKeyvalSeparator, but got %s", err.Error())
	}
}
func TestParserMissingObjectClosingBracket1(t *testing.T) {
	input := `{
		"key1": "val1",
		"key2": "val2"
		`
	lex := lexer.New(strings.NewReader(input))
	parser, err := New(lex)

	if err != nil {
		t.Fatal(err.Error())
	}

	_, err = parser.Parse()

	if err == nil {
		t.Fatal("Expected error, but got nil")
	}

	if !errors.Is(err, ErrMissingObjectClosingBracket) {
		t.Fatalf("Expected ErrMissingObjectClosingBracket, but got %s", err.Error())
	}
}

func TestParserMissingObjectClosingBracket2(t *testing.T) {
	input := `{
		`
	lex := lexer.New(strings.NewReader(input))
	parser, err := New(lex)

	if err != nil {
		t.Fatal(err.Error())
	}

	_, err = parser.Parse()

	if err == nil {
		t.Fatal("Expected error, but got nil")
	}

	if !errors.Is(err, ErrMissingObjectClosingBracket) {
		t.Fatalf("Expected ErrMissingObjectClosingBracket, but got %s", err.Error())
	}
}
func TestParserMissingObjectSeparator(t *testing.T) {
	input := `{
		"key1": "val1"
		"key2": "val2"
		}`
	lex := lexer.New(strings.NewReader(input))
	parser, err := New(lex)

	if err != nil {
		t.Fatal(err.Error())
	}

	_, err = parser.Parse()

	if err == nil {
		t.Fatal("Expected error, but got nil")
	}

	if !errors.Is(err, ErrMissingObjectSeparator) {
		t.Fatalf("Expected ErrMissingObjectSeparator, but got %s", err.Error())
	}
}

func TestParserDuplicateKeys(t *testing.T) {
	input := `{
		"key1": "val1",
		"key1": "val2"
		}`
	lex := lexer.New(strings.NewReader(input))
	parser, err := New(lex)

	if err != nil {
		t.Fatal(err.Error())
	}

	_, err = parser.Parse()

	if err == nil {
		t.Fatal("Expected error, but got nil")
	}

	if !errors.Is(err, ErrDuplicateKeys) {
		t.Fatalf("Expected ErrDuplicateKeys, but got %s", err.Error())
	}
}

func TestParserErrPropragation(t *testing.T) {
	input := `{
		"key1": [
			A3
		]
	}`

	lex := lexer.New(strings.NewReader(input))
	parser, err := New(lex)

	if err != nil {
		t.Fatal(err.Error())
	}

	_, err = parser.Parse()

	if err == nil {
		t.Fatal("Expected error, but got nil")
	}

	if !errors.Is(err, ErrUndefinedToken) {
		t.Fatalf("Expected ErrUndefinedToken, but got %s", err.Error())
	}
}
