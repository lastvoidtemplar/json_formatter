package lexer

import (
	"bufio"
	"go/token"
	"io"
	"iter"
)

func New(r io.Reader) iter.Seq[token.Token] {
	return func(yield func(token.Token) bool) {
		row := 1
		colm := 1

		scanner := bufio.NewScanner(r)
		scanner.Split(splitScannerFunc)

		_, _ = row, colm
	}
}

func splitScannerFunc(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF {
		return 0, nil, nil
	}

	quotes := 0
	ind := 0
	for ; ind < len(data); ind++ {
		b := data[ind]
		switch b {
		case ',':
			if quotes%2 == 0 {
				return ind + 1, data[:ind+1], nil
			}
		case '"':
			if ind == 0 || data[ind-1] != '\\' {
				quotes++
			}
		}
	}

	return ind, data, nil
}
