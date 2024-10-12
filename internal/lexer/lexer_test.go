package lexer

import (
	"bufio"
	"strings"
	"testing"
)

func TestScannerSplitFunc(t *testing.T) {
	input := `"hello" : "world\"ddd",
			"age": 11111`
	scanner := bufio.NewScanner(strings.NewReader(input))

	expected := []string{
		`"hello" : "world\"ddd",`,
		`			"age": 11111`,
	}

	ind := 0
	for scanner.Scan() {
		text := scanner.Text()
		if text != expected[ind] {
			t.Errorf("Text[%d] was expected to be %s, but got %s", ind, expected[ind], text)
		}
		ind++
	}
}
