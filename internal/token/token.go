package token

type TokenType byte

const (
	UNDEFINED TokenType = iota

	EOF

	COLON
	SEMICOLON

	LEFT_SQUARE
	RIGHT_SQUARE
	LEFT_CURLY
	RIGHT_CURLY

	NULL
	BOOL_LITERAL
	NUMBER_LITERAL
	STRING_LITERAL
)

type Token struct {
	Type    TokenType
	Literal string
	Row     int
	Colm    int
}

func New(typ TokenType, literal string, row int, colm int) Token {
	return Token{
		Type:    typ,
		Literal: literal,
		Row:     row,
		Colm:    colm,
	}
}
