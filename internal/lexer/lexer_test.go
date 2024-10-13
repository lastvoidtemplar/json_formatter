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
		t.Fatalf("tok.Literal was expected to be %s, but got %s", "false", tok.Literal)
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
