package parser

import (
	"fmt"
)

type ParserError struct {
	WrapError error
	Row       int
	Colm      int
	Actual    string
}

func (err *ParserError) Error() string {
	if err.Actual == "" {
		return fmt.Sprintf("%s on row %d colm %d", err.WrapError.Error(), err.Row, err.Colm)
	}
	return fmt.Sprintf("%s on row %d colm %d, but got %s", err.WrapError.Error(), err.Row, err.Colm, err.Actual)
}

// for errors.Unwrap
func (err *ParserError) Unwrap() error {
	return err.WrapError
}

func newParserErr(err error, row int, colm int, actual string) *ParserError {
	return &ParserError{
		WrapError: err,
		Row:       row,
		Colm:      colm,
		Actual:    actual,
	}
}
