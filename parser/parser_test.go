package parser

import (
	"testing"

	"github.com/mhoertnagl/donkey/lexer"
)

func TestStatements(t *testing.T) {
	test(t, "let a = 0;", "let a = 0;", 1)
	test(t, "return 42;", "return 42;", 1)
	test(t, "return a;", "return a;", 1)
	test(t, "0;", "0;", 1)
	test(t, "x;", "x;", 1)
	test(t, "-15;", "(-15);", 1)
	test(t, "!true;", "(!true);", 1)
	test(t, "~0;", "(~0);", 1)
	test(t, "~~0;", "(~(~0));", 1)
	test(t, "0; 1; 2;", "0;1;2;", 3)
	test(t, "-0; --1; !false;", "(-0);(-(-1));(!false);", 3)
	test(t, "false || false;", "(false || false);", 1)
	test(t, "false && false;", "(false && false);", 1)
	test(t, "a == 5;", "(a == 5);", 1)
	test(t, "a != 5;", "(a != 5);", 1)
	test(t, "a < 5;", "(a < 5);", 1)
	test(t, "a > 5;", "(a > 5);", 1)
	test(t, "a <= 5;", "(a <= 5);", 1)
	test(t, "a >= 5;", "(a >= 5);", 1)
	test(t, "5 | 5;", "(5 | 5);", 1)
	test(t, "5 ^ 5;", "(5 ^ 5);", 1)
	test(t, "5 & 5;", "(5 & 5);", 1)
	test(t, "5 << 5;", "(5 << 5);", 1)
	test(t, "5 >> 5;", "(5 >> 5);", 1)
	test(t, "5 >>> 5;", "(5 >>> 5);", 1)
	test(t, "5 <>> 5;", "(5 <>> 5);", 1)
	test(t, "5 <<> 5;", "(5 <<> 5);", 1)
	test(t, "5 + 5;", "(5 + 5);", 1)
	test(t, "5 - 5;", "(5 - 5);", 1)
	test(t, "5 * 5;", "(5 * 5);", 1)
	test(t, "5 / 5;", "(5 / 5);", 1)
	// TODO: Test operator precedence.
}

func test(t *testing.T, input string, expected string, n int) {
	lexer := lexer.NewLexer(input)
	parser := NewParser(lexer)
	root := parser.Parse()
	actual := root.String()
	m := len(root.Statements)
	if m != n {
		t.Errorf("Expected [%d] statements but got [%d].", n, m)
	}
	if actual != expected {
		t.Errorf("Expected [%s] but got [%s].", expected, actual)
	}
}
