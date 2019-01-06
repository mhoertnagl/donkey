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

	test(t, "!~-a;", "(!(~(-a)));", 1)
	test(t, "a + b + c;", "((a + b) + c);", 1)
	test(t, "a + b * c + d;", "((a + (b * c)) + d);", 1)
	test(t, "1 * 2 + 3;", "((1 * 2) + 3);", 1)
	test(t, "1 + 2 * 3;", "(1 + (2 * 3));", 1)
	test(t, "1 - -2;", "(1 - (-2));", 1)
	test(t, "-1 - 2;", "((-1) - 2);", 1)

	test(t, "1 + (2 + 3) + 4", "((1 + (2 + 3)) + 4)", 1)
	test(t, "(2 + 3) * 4", "((2 + 3) * 4)", 1)
	// TODO: Test operator precedence.

	test(t, "if (0 < 1) return 1;", "if (0 < 1) return 1;", 1)
	test(t, "if (0 < 1) return 1; else return 0;", "if (0 < 1) return 1; else return 0;", 1)
	test(t, "if (0 < 1) { return 1; }", "if (0 < 1) { return 1; }", 1)
	test(t, "if (0 < 1) { let a = b + c; return a; }", "if (0 < 1) { let a = (b + c);return a; }", 1)
	test(t, "if (0 < 1) { return 1; } else { return 0; }", "if (0 < 1) { return 1; } else { return 0; }", 1)
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
