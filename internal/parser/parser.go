package parser

import (
	"errors"
	"fmt"
	"iter"

	"github.com/lastvoidtemplar/json_formatter/internal/ast"
	"github.com/lastvoidtemplar/json_formatter/internal/token"
)

type Parser struct {
	currToken token.Token
	nextToken func() (token.Token, bool)
	stopLexer func()
	errs      []error
}

func New(lex iter.Seq[token.Token]) *Parser {
	next, stop := iter.Pull(lex)

	tok, ok := next()

	if !ok || tok.Type == token.EOF {
		return nil
	}

	return &Parser{
		nextToken: next,
		stopLexer: stop,
		currToken: tok,
		errs:      make([]error, 0),
	}
}

func (p *Parser) Parse() (ast.Node, []error) {
	defer p.stopLexer()

	if isCurrLeaf(p.currToken) {
		peekToken, ok := p.nextToken()

		if !ok || peekToken.Type == token.EOF {
			return ast.NewLeafNode(p.currToken), nil
		}

		if peekToken.Type == token.ERR {
			p.errs = append(p.errs, errors.New(peekToken.Literal))
			return ast.NewLeafNode(p.currToken), p.errs
		}

		p.errs = append(p.errs,
			fmt.Errorf("row %d colm %d: the expected token was EOF, but got %s", peekToken.Row, peekToken.Colm, peekToken.Literal))
		return ast.NewLeafNode(p.currToken), p.errs
	}

	if p.currToken.Type == token.LEFT_CURLY {
		return nil, nil
	}

	if p.currToken.Type == token.LEFT_SQUARE {
		return nil, nil
	}

	p.errs = append(
		p.errs, fmt.Errorf(
			"row %d colm %d: expected first token was STRING, NUMBER, BOOL, NULL, { or [, but got %s",
			p.currToken.Row, p.currToken.Colm, p.currToken.Literal))
	return nil, p.errs
}

func isCurrLeaf(tok token.Token) bool {
	switch tok.Type {
	case token.NULL, token.TRUE, token.FALSE, token.NUMBER_LITERAL, token.STRING_LITERAL:
		return true
	default:
		return false
	}
}
