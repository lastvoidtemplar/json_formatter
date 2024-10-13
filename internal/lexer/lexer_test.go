package lexer

import (
	"bufio"
	"json_formatter/internal/token"
	"strings"
	"testing"
)

func TestScannerSplitFunc(t *testing.T) {
	input := `"hello" : "world\"ddd",
			"age": 11111`
	scanner := bufio.NewScanner(strings.NewReader(input))
	scanner.Split(splitScannerFunc)

	expected := []string{
		`"hello" : "world\"ddd",`,
		"\n			\"age\": 11111",
	}

	ind := 0
	for scanner.Scan() {
		text := scanner.Text()
		if text != expected[ind] {
			t.Errorf("Text[%d] was expected to be %s(len = %d), but got %s(len = %d)",
				ind, expected[ind], len(expected[ind]), text, len(text))
		}
		ind++
	}
}

func TestTryGetEscapeWithout(t *testing.T) {
	input := `//`

	offset, ok := tryGetEscape(input, 0)

	if ok {
		t.Fatal("Escape invaid escape seq")
	}

	if offset != 0 {
		t.Fatalf("The offset was expected to be %d, but got %d", 0, offset)
	}

}

func TestTryGetEscapeOffest2(t *testing.T) {
	input := `\"\\\/\b\f\n\r\t`
	n := len(input)

	ind := 0
	for ind < n {
		offset, ok := tryGetEscape(input, 0)

		if !ok {
			t.Fatalf("Couldn't escape [%d]", ind/2)
		}

		if offset != 2 {
			t.Fatalf("The offset was expected to be %d, but got %d", 2, offset)
		}
		ind += offset
	}

}
func TestTryGetEscapeOffset6(t *testing.T) {
	input := `\u0055`

	offset, ok := tryGetEscape(input, 0)

	if !ok {
		t.Fatal("Couldn't escape unicode")
	}

	if offset != 6 {
		t.Fatalf("The offset was expected to be %d, but got %d", 6, offset)
	}
}
func TestTryGetEscapeOffset6InvalidLen(t *testing.T) {
	input := `\u005`

	offset, ok := tryGetEscape(input, 0)

	if ok {
		t.Fatal("Escape invaid unicode")
	}

	if offset != 0 {
		t.Fatalf("The offset was expected to be %d, but got %d", 0, offset)
	}
}

func TestTryGetEscapeOffset6InvalidChar(t *testing.T) {
	input := `\u005t`

	offset, ok := tryGetEscape(input, 0)

	if ok {
		t.Fatal("Escape invaid unicode")
	}

	if offset != 0 {
		t.Fatalf("The offset was expected to be %d, but got %d", 0, offset)
	}
}

func TestTryGetStringValid(t *testing.T) {
	input := `junk"hello Деян\nHow are you?😀\uD83D\uDE00"junk`
	ind, row, colm := 4, 1, 5
	var tok token.Token
	var ok bool

	tok, ok, ind, row, colm = tryGetString(input, ind, row, colm)

	if !ok {
		t.Fatal("Failed to get the string token")
	}

	if tok.Type != token.STRING_LITERAL {
		t.Fatalf("tok.Type was expected to be %d, but got %d", token.STRING_LITERAL, tok.Type)
	}

	expectedLiteral := `hello Деян\nHow are you?😀\uD83D\uDE00`
	if tok.Literal != expectedLiteral {
		t.Fatalf("tok.Literal was expected to be %s, but got %s", expectedLiteral, tok.Literal)
	}

	if ind == 0 || input[ind] != 'j' {
		t.Fatalf("Wrong value for ind, expected input[ind]='j', but got %c", input[ind])
	}

	if ind+1 != colm {
		t.Fatalf("Colm was expected to be 1 more then ind (%d), but got %d", ind+1, colm)
	}

	if row != 1 {
		t.Fatalf("Row was expected to remain 1, but got %d", row)
	}
}

func TestTryGetStringValidSmallBuffer(t *testing.T) {
	input := `"hello world`
	ind, row, colm := 0, 1, 1

	var tok token.Token
	var ok bool

	tok, ok, ind, row, colm = tryGetString(input, ind, row, colm)

	if !ok {
		t.Fatal("Failed to get the error token")
	}

	if tok.Type != token.ERR {
		t.Fatalf("tok.Type was expected to be %d, but got %d", token.STRING_LITERAL, tok.Type)
	}

	if ind != len(input) {
		t.Fatalf("Wrong value for ind, expected %d, but got %d", len(input), ind)
	}

	if ind+1 != colm {
		t.Fatalf("Colm was expected to be 1 more then ind (%d), but got %d", ind+1, colm)
	}

	if row != 1 {
		t.Fatalf("Row was expected to remain 1, but got %d", row)
	}
}

func TestTryGetStringInvalidStaringWithNoQuote(t *testing.T) {
	input := `hello world"`
	ind, row, colm := 0, 1, 1

	var ok bool

	_, ok, ind, row, colm = tryGetString(input, ind, row, colm)

	if ok {
		t.Fatal("Failed to get discard invalid string")
	}

	if ind != 0 || row != 1 || colm != 1 {
		t.Fatal("Pointers have moved")
	}
}

func TestTryGetStringInvalidEscape(t *testing.T) {
	input := `"Hello \' world"`
	ind, row, colm := 0, 1, 1

	var ok bool

	_, ok, ind, row, colm = tryGetString(input, ind, row, colm)

	if ok {
		t.Fatal("Failed to get discard invalid string")
	}

	if ind != 0 || row != 1 || colm != 1 {
		t.Fatal("Pointers have moved")
	}
}

func TestTryGetStringInvalidControlChar(t *testing.T) {
	input := `"Hello
	 world"`
	ind, row, colm := 0, 1, 1

	var ok bool

	_, ok, ind, row, colm = tryGetString(input, ind, row, colm)

	if ok {
		t.Fatal("Failed to get discard invalid string")
	}

	if ind != 0 || row != 1 || colm != 1 {
		t.Fatal("Pointers have moved")
	}
}

func TestTryGetKeywordValidNull(t *testing.T) {
	input := "null"
	ind, row, colm := 0, 1, 1

	var tok token.Token
	var ok bool
	tok, ok, ind, row, colm = tryGetKeyword(input, ind, row, colm)

	if !ok {
		t.Fatal("Couldn`t get null token")
	}

	if tok.Type != token.NULL {
		t.Fatalf("tok.Type was expected to be %d, but got %d", token.NULL, tok.Type)
	}

	if tok.Literal != "null" {
		t.Fatalf("tok.Literal was expected to be %s, but got %s", "null", tok.Literal)
	}

	if ind != len(input) {
		t.Fatalf("Wrong value for ind, expected %d, but got %d", len(input), ind)
	}

	if ind+1 != colm {
		t.Fatalf("Colm was expected to be 1 more then ind (%d), but got %d", ind+1, colm)
	}

	if row != 1 {
		t.Fatalf("Row was expected to remain 1, but got %d", row)
	}
}

func TestTryGetKeywordInvalidNull(t *testing.T) {
	input := "nulle"
	ind, row, colm := 0, 1, 1

	var ok bool
	_, ok, ind, row, colm = tryGetKeyword(input, ind, row, colm)

	if ok {
		t.Fatal("Failed to discard invalid null token")
	}

	if ind != 0 || row != 1 || colm != 1 {
		t.Fatal("Pointers have moved")
	}
}

func TestTryGetKeywordValidTrue(t *testing.T) {
	input := "true"
	ind, row, colm := 0, 1, 1

	var tok token.Token
	var ok bool
	tok, ok, ind, row, colm = tryGetKeyword(input, ind, row, colm)

	if !ok {
		t.Fatal("Couldn`t get true token")
	}

	if tok.Type != token.TRUE {
		t.Fatalf("tok.Type was expected to be %d, but got %d", token.TRUE, tok.Type)
	}

	if tok.Literal != "true" {
		t.Fatalf("tok.Literal was expected to be %s, but got %s", "true", tok.Literal)
	}

	if ind != len(input) {
		t.Fatalf("Wrong value for ind, expected %d, but got %d", len(input), ind)
	}

	if ind+1 != colm {
		t.Fatalf("Colm was expected to be 1 more then ind (%d), but got %d", ind+1, colm)
	}

	if row != 1 {
		t.Fatalf("Row was expected to remain 1, but got %d", row)
	}
}

func TestTryGetKeywordInvalidTrue(t *testing.T) {
	input := "truee"
	ind, row, colm := 0, 1, 1

	var ok bool
	_, ok, ind, row, colm = tryGetKeyword(input, ind, row, colm)

	if ok {
		t.Fatal("Failed to discard invalid true token")
	}

	if ind != 0 || row != 1 || colm != 1 {
		t.Fatal("Pointers have moved")
	}
}
func TestTryGetKeywordValidFalse(t *testing.T) {
	input := "false"
	ind, row, colm := 0, 1, 1

	var tok token.Token
	var ok bool
	tok, ok, ind, row, colm = tryGetKeyword(input, ind, row, colm)

	if !ok {
		t.Fatal("Couldn`t get false token")
	}

	if tok.Type != token.FALSE {
		t.Fatalf("tok.Type was expected to be %d, but got %d", token.FALSE, tok.Type)
	}

	if tok.Literal != "false" {
		t.Fatalf("tok.Literal was expected to be %s, but got %s", "false", tok.Literal)
	}

	if ind != len(input) {
		t.Fatalf("Wrong value for ind, expected %d, but got %d", len(input), ind)
	}

	if ind+1 != colm {
		t.Fatalf("Colm was expected to be 1 more then ind (%d), but got %d", ind+1, colm)
	}

	if row != 1 {
		t.Fatalf("Row was expected to remain 1, but got %d", row)
	}
}
func TestTryGetKeywordInvalidFalse(t *testing.T) {
	input := "falsee"
	ind, row, colm := 0, 1, 1

	var ok bool
	_, ok, ind, row, colm = tryGetKeyword(input, ind, row, colm)

	if ok {
		t.Fatal("Failed to discard invalid false token")
	}

	if ind != 0 || row != 1 || colm != 1 {
		t.Fatal("Pointers have moved")
	}
}

func TestTryGetKeywordInvalid(t *testing.T) {
	input := "dfdf"
	ind, row, colm := 0, 1, 1

	var ok bool
	_, ok, ind, row, colm = tryGetKeyword(input, ind, row, colm)

	if ok {
		t.Fatal("Failed to discard invalid invalid keyword")
	}

	if ind != 0 || row != 1 || colm != 1 {
		t.Fatal("Pointers have moved")
	}
}

func TestTryGetNumberWhole(t *testing.T) {
	input := "erer -1234"
	ind, row, colm := 5, 1, 6

	var tok token.Token
	var ok bool
	tok, ok, ind, row, colm = tryGetNumber(input, ind, row, colm)

	if !ok {
		t.Fatal("Couldn`t get number token")
	}

	if tok.Type != token.NUMBER_LITERAL {
		t.Fatalf("tok.Type was expected to be %d, but got %d", token.NUMBER_LITERAL, tok.Type)
	}

	if tok.Literal != "-1234" {
		t.Fatalf("tok.Literal was expected to be %s, but got %s", "-1234", tok.Literal)
	}

	if ind != len(input) {
		t.Fatalf("Wrong value for ind, expected %d, but got %d", len(input), ind)
	}

	if ind+1 != colm {
		t.Fatalf("Colm was expected to be 1 more then ind (%d), but got %d", ind+1, colm)
	}

	if row != 1 {
		t.Fatalf("Row was expected to remain 1, but got %d", row)
	}
}
func TestTryGetNumberFraction(t *testing.T) {
	input := "erer 1234.25"
	ind, row, colm := 5, 1, 6

	var tok token.Token
	var ok bool
	tok, ok, ind, row, colm = tryGetNumber(input, ind, row, colm)

	if !ok {
		t.Fatal("Couldn`t get number token")
	}

	if tok.Type != token.NUMBER_LITERAL {
		t.Fatalf("tok.Type was expected to be %d, but got %d", token.NUMBER_LITERAL, tok.Type)
	}

	if tok.Literal != "1234.25" {
		t.Fatalf("tok.Literal was expected to be %s, but got %s", "1234.25", tok.Literal)
	}

	if ind != len(input) {
		t.Fatalf("Wrong value for ind, expected %d, but got %d", len(input), ind)
	}

	if ind+1 != colm {
		t.Fatalf("Colm was expected to be 1 more then ind (%d), but got %d", ind+1, colm)
	}

	if row != 1 {
		t.Fatalf("Row was expected to remain 1, but got %d", row)
	}
}
func TestTryGetNumberSciNotion(t *testing.T) {
	input := "erer 1234e-3"
	ind, row, colm := 5, 1, 6

	var tok token.Token
	var ok bool
	tok, ok, ind, row, colm = tryGetNumber(input, ind, row, colm)

	if !ok {
		t.Fatal("Couldn`t get number token")
	}

	if tok.Type != token.NUMBER_LITERAL {
		t.Fatalf("tok.Type was expected to be %d, but got %d", token.NUMBER_LITERAL, tok.Type)
	}

	if tok.Literal != "1234e-3" {
		t.Fatalf("tok.Literal was expected to be %s, but got %s", "1234e-3", tok.Literal)
	}

	if ind != len(input) {
		t.Fatalf("Wrong value for ind, expected %d, but got %d", len(input), ind)
	}

	if ind+1 != colm {
		t.Fatalf("Colm was expected to be 1 more then ind (%d), but got %d", ind+1, colm)
	}

	if row != 1 {
		t.Fatalf("Row was expected to remain 1, but got %d", row)
	}
}
func TestTryGetNumberFractionStartingWithZeroAndHasSciNotion(t *testing.T) {
	input := "erer 0.1234e3 ere"
	ind, row, colm := 5, 1, 6

	var tok token.Token
	var ok bool
	tok, ok, ind, row, colm = tryGetNumber(input, ind, row, colm)

	if !ok {
		t.Fatal("Couldn`t get number token")
	}

	if tok.Type != token.NUMBER_LITERAL {
		t.Fatalf("tok.Type was expected to be %d, but got %d", token.NUMBER_LITERAL, tok.Type)
	}

	if tok.Literal != "0.1234e3" {
		t.Fatalf("tok.Literal was expected to be %s, but got %s", "0.1234e3", tok.Literal)
	}

	if ind != 13 {
		t.Fatalf("Wrong value for ind, expected %d, but got %d", 13, ind)
	}

	if ind+1 != colm {
		t.Fatalf("Colm was expected to be 1 more then ind (%d), but got %d", ind+1, colm)
	}

	if row != 1 {
		t.Fatalf("Row was expected to remain 1, but got %d", row)
	}
}

func TestTryGetNumberInvalidFirstCharNotDigitOrSign(t *testing.T) {
	input := "f4444"
	ind, row, colm := 0, 1, 1

	var ok bool
	_, ok, ind, row, colm = tryGetNumber(input, ind, row, colm)

	if ok {
		t.Fatal("Failed to discard invalid number token")
	}

	if ind != 0 || row != 1 || colm != 1 {
		t.Fatal("Pointers have moved")
	}
}

func TestTryNumberInvalidOctal(t *testing.T) {
	input := "0640"

	ind, row, colm := 0, 1, 1

	var ok bool
	_, ok, ind, row, colm = tryGetNumber(input, ind, row, colm)

	if ok {
		t.Fatal("Failed to discard invalid number token")
	}

	if ind != 0 || row != 1 || colm != 1 {
		t.Fatal("Pointers have moved")
	}
}
func TestTryNumberInvalidOHex(t *testing.T) {
	input := "0x640"

	ind, row, colm := 0, 1, 1

	var ok bool
	_, ok, ind, row, colm = tryGetNumber(input, ind, row, colm)

	if ok {
		t.Fatal("Failed to discard invalid number token")
	}

	if ind != 0 || row != 1 || colm != 1 {
		t.Fatal("Pointers have moved")
	}
}
func TestTryNumberInvalidNonDigitAferPoint(t *testing.T) {
	input := "0.e3"

	ind, row, colm := 0, 1, 1

	var ok bool
	_, ok, ind, row, colm = tryGetNumber(input, ind, row, colm)

	if ok {
		t.Fatal("Failed to discard invalid number token")
	}

	if ind != 0 || row != 1 || colm != 1 {
		t.Fatal("Pointers have moved")
	}
}
func TestTryNumberInvalidNonDigitAferPoint2(t *testing.T) {
	input := "0."

	ind, row, colm := 0, 1, 1

	var ok bool
	_, ok, ind, row, colm = tryGetNumber(input, ind, row, colm)

	if ok {
		t.Fatal("Failed to discard invalid number token")
	}

	if ind != 0 || row != 1 || colm != 1 {
		t.Fatal("Pointers have moved")
	}
}
func TestTryNumberInvalidNonDigitAfterExponent(t *testing.T) {
	input := "10e"

	ind, row, colm := 0, 1, 1

	var ok bool
	_, ok, ind, row, colm = tryGetNumber(input, ind, row, colm)

	if ok {
		t.Fatal("Failed to discard invalid number token")
	}

	if ind != 0 || row != 1 || colm != 1 {
		t.Fatal("Pointers have moved")
	}
}
func TestTryNumberInvalidNonDigitAfterExponent2(t *testing.T) {
	input := "10e+"

	ind, row, colm := 0, 1, 1

	var ok bool
	_, ok, ind, row, colm = tryGetNumber(input, ind, row, colm)

	if ok {
		t.Fatal("Failed to discard invalid number token")
	}

	if ind != 0 || row != 1 || colm != 1 {
		t.Fatal("Pointers have moved")
	}
}

func TestGetUndefined(t *testing.T) {
	input := "dfsfsf3333"
	ind, row, colm := 0, 1, 1
	var tok token.Token

	tok, ind, row, colm = getUndefined(input, ind, row, colm)

	if tok.Type != token.UNDEFINED {
		t.Fatalf("tok.Type was expected to be %d, but got %d", token.UNDEFINED, tok.Type)
	}

	if tok.Literal != input {
		t.Fatalf("tok.Literal was expected to be %s, but got %s", input, tok.Literal)
	}

	if ind != len(input) {
		t.Fatalf("Wrong value for ind, expected %d, but got %d", len(input), ind)
	}

	if ind+1 != colm {
		t.Fatalf("Colm was expected to be 1 more then ind (%d), but got %d", ind+1, colm)
	}

	if row != 1 {
		t.Fatalf("Row was expected to remain 1, but got %d", row)
	}
}
func TestGetUndefined1(t *testing.T) {
	input := "dfsfsf3333\n"
	ind, row, colm := 0, 1, 1
	var tok token.Token

	tok, ind, row, colm = getUndefined(input, ind, row, colm)

	if tok.Type != token.UNDEFINED {
		t.Fatalf("tok.Type was expected to be %d, but got %d", token.UNDEFINED, tok.Type)
	}

	if tok.Literal != input[:len(input)-1] {
		t.Fatalf("tok.Literal was expected to be %s, but got %s", input[:len(input)-1], tok.Literal)
	}

	if ind != len(input)-1 {
		t.Fatalf("Wrong value for ind, expected %d, but got %d", len(input), ind)
	}

	if ind+1 != colm {
		t.Fatalf("Colm was expected to be 1 more then ind (%d), but got %d", ind+1, colm)
	}

	if row != 1 {
		t.Fatalf("Row was expected to remain 1, but got %d", row)
	}
}

func TestLexer(t *testing.T) {
	input := `{
	"name": "Name",
	"age": 21,
	"human": true,
	"hobbies": [
		"Programming",
		false,
		42.69,
		null
	],
	"indefi\'ned": 2.
}`
	lex := New(strings.NewReader(input))

	expected := []token.Token{
		{Type: token.LEFT_CURLY, Literal: "{", Row: 1, Colm: 1},
		{Type: token.STRING_LITERAL, Literal: "name", Row: 2, Colm: 2},
		{Type: token.COLON, Literal: ":", Row: 2, Colm: 8},
		{Type: token.STRING_LITERAL, Literal: "Name", Row: 2, Colm: 10},
		{Type: token.SEMICOLON, Literal: ",", Row: 2, Colm: 16},
		{Type: token.STRING_LITERAL, Literal: "age", Row: 3, Colm: 2},
		{Type: token.COLON, Literal: ":", Row: 3, Colm: 7},
		{Type: token.NUMBER_LITERAL, Literal: "21", Row: 3, Colm: 9},
		{Type: token.SEMICOLON, Literal: ",", Row: 3, Colm: 11},
		{Type: token.STRING_LITERAL, Literal: "human", Row: 4, Colm: 2},
		{Type: token.COLON, Literal: ":", Row: 4, Colm: 9},
		{Type: token.TRUE, Literal: "true", Row: 4, Colm: 11},
		{Type: token.SEMICOLON, Literal: ",", Row: 4, Colm: 15},
		{Type: token.STRING_LITERAL, Literal: "hobbies", Row: 5, Colm: 2},
		{Type: token.COLON, Literal: ":", Row: 5, Colm: 11},
		{Type: token.LEFT_SQUARE, Literal: "[", Row: 5, Colm: 13},
		{Type: token.STRING_LITERAL, Literal: "Programming", Row: 6, Colm: 3},
		{Type: token.SEMICOLON, Literal: ",", Row: 6, Colm: 16},
		{Type: token.FALSE, Literal: "false", Row: 7, Colm: 3},
		{Type: token.SEMICOLON, Literal: ",", Row: 7, Colm: 8},
		{Type: token.NUMBER_LITERAL, Literal: "42.69", Row: 8, Colm: 3},
		{Type: token.SEMICOLON, Literal: ",", Row: 8, Colm: 8},
		{Type: token.NULL, Literal: "null", Row: 9, Colm: 3},
		{Type: token.RIGHT_SQUARE, Literal: "]", Row: 10, Colm: 2},
		{Type: token.SEMICOLON, Literal: ",", Row: 10, Colm: 3},
		{Type: token.UNDEFINED, Literal: `"indefi\'ned"`, Row: 11, Colm: 2},
		{Type: token.COLON, Literal: ":", Row: 11, Colm: 15},
		{Type: token.UNDEFINED, Literal: `2.`, Row: 11, Colm: 17},
		{Type: token.RIGHT_CURLY, Literal: "}", Row: 12, Colm: 1},
	}
	ind := 0
	for tok := range lex {
		exp := expected[ind]
		if tok.Type != exp.Type {
			t.Errorf("tok[%d].Type was expected to be %d, but got %d", ind, exp.Type, tok.Type)
		}
		if tok.Literal != exp.Literal {
			t.Errorf("tok[%d].Literal was expected to be %s, but got %s", ind, exp.Literal, tok.Literal)
		}
		if tok.Row != exp.Row {
			t.Errorf("tok[%d].Row was expected to be %d, but got %d", ind, exp.Row, tok.Row)
		}
		if tok.Colm != exp.Colm {
			t.Errorf("tok[%d].Colm was expected to be %d, but got %d", ind, exp.Colm, tok.Colm)
		}
		ind++
	}

	if ind != len(expected) {
		t.Errorf("Expected len was %d, but got %d", len(expected), ind)
	}
}
