package parser

import (
	"testing"

	"github.com/mhoertnagl/donkey/lexer"
)

func Test(t *testing.T) {
	test(t, "let a = 0;", "let a = 0;")
	test(t, "return a;", "return a;")
}

func test(t *testing.T, input string, expected string) {
	lexer := lexer.NewLexer(input)
	parser := NewParser(lexer)
	root := parser.Parse()
	actual := root.String()
	//t.Logf("Debug: %v", root)
	if actual != expected {
		t.Errorf("Expected [%s] but got [%s].", expected, actual)
	}
}
