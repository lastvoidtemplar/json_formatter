package parser

import (
	"errors"
	"fmt"
	"iter"
	"log"

	"github.com/lastvoidtemplar/json_formatter/internal/ast"
	"github.com/lastvoidtemplar/json_formatter/internal/token"
)

type Parser struct {
	currToken token.Token
	peekToken token.Token
	nextToken func() (token.Token, bool)
	stopLexer func()
	parserErr *ParserError
}

var ErrEmptyLexer = errors.New("the lexer is empty")

func New(lex iter.Seq[token.Token]) (*Parser, error) {
	next, stop := iter.Pull(lex)

	currToken, ok := next()

	if !ok || currToken.Type == token.EOF {
		return nil, ErrEmptyLexer
	}

	peekToken, ok := next()

	if !ok {
		peekToken = newEOF()
	}

	return &Parser{
		nextToken: next,
		stopLexer: stop,
		currToken: currToken,
		peekToken: peekToken,
		parserErr: nil,
	}, nil
}

func (p *Parser) Parse() (ast.Node, *ParserError) {
	defer p.stopLexer()

	root := p.parseNode()

	if root == nil {
		panic("root is nil")
	}

	if p.currToken.Type != token.EOF {
		p.parserErr = newParserErr(fmt.Sprintf("expected EOF on row %d, colm %d but got %s",
			p.currToken.Row, p.currToken.Colm, p.currToken.Literal), p.currToken.Row, p.currToken.Colm)
		return nil, p.parserErr
	}

	if p != nil {
		return nil, p.parserErr
	}

	return root, nil
}

func (p *Parser) parseNode() ast.Node {
	if isCurrLeaf(p.currToken) {
		tok := p.parserLeaf()
		p.NextToken()
		return tok
	}

	if p.currToken.Type == token.LEFT_SQUARE {
		tok := p.parseArray()
		p.NextToken()
		return tok

	}

	if p.currToken.Type == token.LEFT_CURLY {
		tok := p.parseObject()
		p.NextToken()
		return tok
	}

	return nil
}

func (p *Parser) parserLeaf() ast.LeafNode {
	if !isCurrLeaf(p.currToken) {
		log.Println("Expecred leaf token")
		p.parserErr = newParserErr(fmt.Sprintf("expected leaf on row %d, colm %d, but got %s",
			p.currToken.Row, p.currToken.Colm, p.currToken.Literal), p.currToken.Row, p.currToken.Colm)
		return nil
	}

	if p.currToken.Type == token.UNDEFINED {
		p.parserErr = newParserErr(fmt.Sprintf("undefined token on row %d, colm %d - %s",
			p.currToken.Row, p.currToken.Colm, p.currToken.Literal), p.currToken.Row, p.currToken.Colm)
		return nil
	}

	return ast.NewLeafNode(p.currToken)
}

func (p *Parser) parseArray() *ast.ArrayNode {
	if p.currToken.Type != token.LEFT_SQUARE {
		log.Println("The starting token of array is not [")
		p.parserErr = newParserErr(fmt.Sprintf("expected [ on row %d, colm %d, but got %s",
			p.currToken.Row, p.currToken.Colm, p.currToken.Literal), p.currToken.Row, p.currToken.Colm)
		return nil
	}

	arrNode := ast.NewArrayNode()

	p.NextToken()
	for p.currToken.Type != token.RIGHT_SQUARE {
		if p.currToken.Type == token.EOF {
			p.parserErr = newParserErr(fmt.Sprintf("Expected ] on row %d, colm %d, but got EOF",
				p.currToken.Row, p.currToken.Colm), p.currToken.Row, p.currToken.Colm)
			return nil
		}

		node := p.parseNode()

		if p.parserErr != nil {
			return nil
		}

		if node == nil {
			panic("Couldn`t parse array element")
		}
		arrNode.Add(node)

		p.NextToken()

		if p.currToken.Type != token.SEMICOLON && p.currToken.Type != token.RIGHT_SQUARE {
			p.parserErr = newParserErr(fmt.Sprintf("expected SEMICOLON or ] on row %d, colm%d got %s",
				p.currToken.Row, p.currToken.Colm, p.currToken.Literal), p.currToken.Row, p.currToken.Colm)
			return nil
		}

		if p.currToken.Type == token.SEMICOLON {
			p.NextToken()
		}
	}

	return arrNode
}

func (p *Parser) parseObject() *ast.ObjectNode {
	if p.currToken.Type == token.EOF {
		p.parserErr = newParserErr(fmt.Sprintf("Expected ] on row %d, colm %d, but got EOF",
			p.currToken.Row, p.currToken.Colm), p.currToken.Row, p.currToken.Colm)
		return nil
	}

	objNode := ast.NewObjectNode()

	p.NextToken()

	for p.currToken.Type != token.RIGHT_CURLY {
		if p.currToken.Type == token.EOF {
			p.parserErr = newParserErr(fmt.Sprintf("Expected ] on row %d, colm %d, but got EOF",
				p.currToken.Row, p.currToken.Colm), p.currToken.Row, p.currToken.Colm)
			return nil
		}

		node := p.parseKeyVal()

		if p.parserErr != nil {
			return nil
		}

		if node == nil {
			panic("Couldn`t parse object element")
		}

		ok := objNode.Add(node)

		if !ok {
			p.parserErr = newParserErr(fmt.Sprintf("duplicate key on row %d, colm%d",
				node.Key.Row, node.Key.Colm), node.Key.Row, node.Key.Colm)
		}

		p.NextToken()

		if p.currToken.Type != token.SEMICOLON && p.currToken.Type != token.RIGHT_CURLY {
			p.parserErr = newParserErr(fmt.Sprintf("expected SEMICOLON or ] on row %d, colm%d got %s",
				p.currToken.Row, p.currToken.Colm, p.currToken.Literal), p.currToken.Row, p.currToken.Colm)
			return nil
		}

		if p.currToken.Type != token.SEMICOLON {
			p.NextToken()
		}
	}

	return objNode
}

func (p *Parser) parseKeyVal() *ast.KeyValNode {
	if p.currToken.Type != token.STRING_LITERAL {
		p.parserErr = newParserErr(fmt.Sprintf("the key must be string on row %d, colm %d, but got %s",
			p.currToken.Row, p.currToken.Colm, p.currToken.Literal), p.currToken.Row, p.currToken.Colm)
		return nil
	}

	key := p.currToken
	p.NextToken()

	if p.currToken.Type != token.COLON {
		p.parserErr = newParserErr(fmt.Sprintf("expected COLON on row %d, colm %d, but got %s",
			p.currToken.Row, p.currToken.Colm, p.currToken.Literal), p.currToken.Row, p.currToken.Colm)
		return nil
	}

	p.NextToken()

	val := p.parseNode()
	if p.parserErr != nil {
		return nil
	}

	return ast.NewKeyVal(key, val)
}

func (p *Parser) NextToken() {
	p.currToken = p.peekToken
	tok, ok := p.nextToken()

	if !ok {
		p.peekToken = newEOF()
		return
	}

	p.peekToken = tok
}

func isCurrLeaf(tok token.Token) bool {
	switch tok.Type {
	case token.UNDEFINED, token.NULL, token.TRUE, token.FALSE, token.NUMBER_LITERAL, token.STRING_LITERAL:
		return true
	default:
		return false
	}
}

func newEOF() token.Token {
	return token.New(token.EOF, "", -1, -1)
}
