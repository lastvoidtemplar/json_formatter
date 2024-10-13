package lexer

import (
	"bufio"
	"io"
	"iter"
	"json_formatter/internal/token"
)

func New(r io.Reader) iter.Seq[token.Token] {
	return func(yield func(token.Token) bool) {
		row := 1
		colm := 1

		scanner := bufio.NewScanner(r)
		scanner.Split(splitScannerFunc)

		for scanner.Scan() {
			input := scanner.Text()
			ind := 0
			var token token.Token

			n := len(input)
			for ind < n {
				token, ind, row, colm = getToken(input, ind, row, colm)
				if !yield(token) {
					return
				}
			}

		}

		if err := scanner.Err(); err != nil {
			yield(token.New(token.ERR, err.Error(), row, colm))
			return
		}
	}
}

func splitScannerFunc(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF {
		return 0, nil, nil
	}

	quotes := 0
	ind := 0
	n := len(data)
	for ; ind < n; ind++ {
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

func skipWhiteSpace(input string, ind int, row int, colm int) (int, int, int) {
	n := len(input)
	for i := ind; i < n; i++ {
		switch input[i] {
		case ' ':
			fallthrough
		case '\t':
			colm++
		case '\r':
			continue
		case '\n':
			colm = 1
			row++
		default:
			return i, row, colm
		}
	}
	return len(input), row, colm
}

func getToken(input string, ind int, row int, colm int) (token.Token, int, int, int) {
	ind, row, colm = skipWhiteSpace(input, ind, row, colm)

	var tok token.Token
	switch input[ind] {
	case 0:
		return token.New(token.EOF, "", row, colm), ind + 1, row, colm + 1
	case ':':
		return token.New(token.COLON, ":", row, colm), ind + 1, row, colm + 1
	case ',':
		return token.New(token.SEMICOLON, ",", row, colm), ind + 1, row, colm + 1
	case '[':
		return token.New(token.LEFT_SQUARE, "[", row, colm), ind + 1, row, colm + 1
	case ']':
		return token.New(token.RIGHT_SQUARE, "]", row, colm), ind + 1, row, colm + 1
	case '{':
		return token.New(token.LEFT_CURLY, "{", row, colm), ind + 1, row, colm + 1
	case '}':
		return token.New(token.RIGHT_CURLY, "}", row, colm), ind + 1, row, colm + 1
	case '"':
		var ok bool
		tok, ok, ind, row, colm = tryGetString(input, ind, row, colm)
		if ok {
			return tok, ind, row, colm
		} else {
			//TODO get undefined
		}
	default:
		var ok bool
		tok, ok, ind, row, colm = tryGetString(input, ind, row, colm)
		if ok {
			return tok, ind, row, colm
		}
		//TODO try get number
		//TODO get undefined
	}

	return token.New(token.UNDEFINED, "", row, colm), ind, row, colm
}

func tryGetString(input string, ind int, row int, colm int) (token.Token, bool, int, int, int) {
	if input[ind] != '"' {
		return token.Token{}, false, ind, row, colm
	}

	n := len(input)
	skip := 0
	for i, b := range input[ind+1:] {
		if skip > 0 {
			skip--
			continue
		}
		switch b {
		case '"':
			return token.New(token.STRING_LITERAL, input[ind+1:ind+i+1], row, colm), true, ind + i + 2, row, colm + i + 2
		case '\\':
			offest, ok := tryGetEscape(input, ind+i+1)
			if ok {
				skip += offest - 1
			} else {
				return token.Token{}, false, ind, row, colm
			}
		default:
			if isControlChar(b) {
				return token.Token{}, false, ind, row, colm
			}
		}

	}
	return token.New(token.ERR, "the buffer was to small to close a string", row, colm), true, n, row, colm + n - ind
}

func tryGetEscape(input string, ind int) (int, bool) {
	escapeU := false
	for i, b := range input[ind:] {
		if i == 0 {
			if b != '\\' {
				return 0, false
			}
		} else if i == 1 {
			switch b {
			case '"', '\\', '/', 'b', 'f', 'n', 'r', 't':
				return 2, true
			case 'u':
				if len(input[ind:]) < 6 {
					return 0, false
				}
				escapeU = true
			default:
				return 0, false
			}
		} else if 2 <= i && i <= 5 {
			if !escapeU || !isHex(b) {
				return 0, false
			}
		} else {
			return 6, true
		}
	}

	return 6, true
}

func isHex(b rune) bool {
	return '0' <= b && b <= '9' || 'a' <= b && b <= 'f' || 'A' <= b && b <= 'F'
}

func isControlChar(b rune) bool {
	return 0 <= b && b <= 31
}

func tryGetKeyword(input string, ind int, row int, colm int) (token.Token, bool, int, int, int) {
	n := len(input[ind:])
	if n >= 5 {
		keyword := input[ind : ind+5]
		if keyword == "false" {
			return token.New(token.FALSE, "false", row, colm), true, ind + 5, row, colm + 5
		}
	} else if n >= 4 {
		keyword := input[ind : ind+4]
		if keyword == "null" {
			return token.New(token.NULL, "null", row, colm), true, ind + 4, row, colm + 4
		} else if keyword == "true" {
			return token.New(token.TRUE, "true", row, colm), true, ind + 4, row, colm + 4
		}
	}
	return token.Token{}, false, ind, row, colm
}

func tryGetNumber(input string, ind int, row int, colm int) (token.Token, bool, int, int, int) {
	i := 0
	n := len(input)
	if input[ind] == '-' {
		i++
	}

	b := input[ind+i]
	if !isDigit(rune(b)) {
		return token.Token{}, false, ind, row, colm
	}

	if isDigitBiggerThanZero(b) {
		for _, v := range input[ind:] {
			if !isDigit(v) {
				break
			}
			i++
		}
		if ind+i == n {
			return token.New(token.NUMBER_LITERAL, input[ind:ind+i], row, colm), true, n, row, colm + n - ind
		}
	}

	return token.Token{}, false, ind, row, colm
}

func isDigitBiggerThanZero(b byte) bool {
	return '1' <= b && b <= '9'
}

func isDigit(b rune) bool {
	return '0' <= b && b <= '9'
}
