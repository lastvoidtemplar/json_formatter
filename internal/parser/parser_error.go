package parser

type ParserError struct {
	Message string
	Row     int
	Colm    int
}

func newParserErr(msg string, row int, colm int) *ParserError {
	return &ParserError{
		Message: msg,
		Row:     row,
		Colm:    colm,
	}
}
