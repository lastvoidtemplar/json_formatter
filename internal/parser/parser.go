package parser

import (
	"errors"
	"iter"

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

// if EOF token is not found, there is a bug in the lexer
var ErrEofNotFound = errors.New("expected EOF")

// undefined token as a leaf
var ErrUndefinedToken = errors.New("expected STRING, NUMBER, TRUE, FALSE or NULL")

// did not found closing square bracket
var ErrMissingArrayClosingBracket = errors.New("expected ']'")

// missing semicolon seperator in array element or missing closing square bracket
var ErrMissingArraySeperator = errors.New("expected SEMICOLON or ']'")

// did not found closing curly bracket
var ErrMissingObjectClosingBracket = errors.New("expected '}'")

// missing semicolon seperator in object element or missing closing curly bracket
var ErrMissingObjectSeperator = errors.New("expected SEMICOLON or '}'")

// found duplacate keys
var ErrDuplicateKeys = errors.New("duplicate keys")

// key must be string
var ErrKeyNotString = errors.New("key must be string")

// missing colon in keyval
var ErrMissingKeyvalSeperator = errors.New("expectted COLON")

func (p *Parser) Parse() (ast.Node, *ParserError) {
	defer p.stopLexer()

	root := p.parseNode()

	if p.parserErr != nil {
		return nil, p.parserErr
	}

	if root == nil {
		panic("Couldn`t parse, but there was no parsing error")
	}

	if p.currToken.Type != token.EOF {
		p.parserErr = newParserErr(ErrEofNotFound, p.currToken.Row, p.currToken.Colm, p.currToken.Literal)
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
	if p.currToken.Type == token.UNDEFINED {
		p.parserErr = newParserErr(ErrUndefinedToken, p.currToken.Row, p.currToken.Colm, p.currToken.Literal)
		return nil
	}

	return ast.NewLeafNode(p.currToken)
}

func (p *Parser) parseArray() *ast.ArrayNode {
	arrNode := ast.NewArrayNode()

	p.NextToken()
	for p.currToken.Type != token.RIGHT_SQUARE {
		if p.currToken.Type == token.EOF {
			p.parserErr = newParserErr(ErrMissingArrayClosingBracket, p.currToken.Row, p.currToken.Colm, "EOF")
			return nil
		}

		node := p.parseNode()

		if p.parserErr != nil {
			return nil
		}

		if node == nil {
			panic("Couldn`t parse array element, but there was not parsing error")
		}
		arrNode.Add(node)

		if p.currToken.Type != token.SEMICOLON && p.currToken.Type != token.RIGHT_SQUARE {
			p.parserErr = newParserErr(ErrMissingArraySeperator, p.currToken.Row, p.currToken.Colm, p.currToken.Literal)
			return nil
		}

		if p.currToken.Type == token.SEMICOLON {
			p.NextToken()
		}
	}

	return arrNode
}

func (p *Parser) parseObject() *ast.ObjectNode {
	objNode := ast.NewObjectNode()

	p.NextToken()

	for p.currToken.Type != token.RIGHT_CURLY {
		if p.currToken.Type == token.EOF {
			p.parserErr = newParserErr(ErrMissingObjectClosingBracket, p.currToken.Row, p.currToken.Colm, "EOF")
			return nil
		}

		node := p.parseKeyVal()

		if p.parserErr != nil {
			return nil
		}

		if node == nil {
			panic("Couldn`t parse object element, but there was no parsing error")
		}

		ok := objNode.Add(node)

		if !ok {
			p.parserErr = newParserErr(ErrDuplicateKeys, node.Key.Row, node.Key.Colm, "")
			return nil
		}

		if p.currToken.Type != token.SEMICOLON && p.currToken.Type != token.RIGHT_CURLY {
			p.parserErr = newParserErr(ErrMissingObjectSeperator, p.currToken.Row, p.currToken.Colm, p.currToken.Literal)
			return nil
		}

		if p.currToken.Type == token.SEMICOLON {
			p.NextToken()
		}
	}

	return objNode
}

func (p *Parser) parseKeyVal() *ast.KeyValNode {
	if p.currToken.Type != token.STRING_LITERAL {
		p.parserErr = newParserErr(ErrKeyNotString, p.currToken.Row, p.currToken.Colm, p.currToken.Literal)
		return nil
	}

	key := p.currToken
	p.NextToken()

	if p.currToken.Type != token.COLON {
		p.parserErr = newParserErr(ErrMissingKeyvalSeperator, p.currToken.Row, p.currToken.Colm, p.currToken.Literal)
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
