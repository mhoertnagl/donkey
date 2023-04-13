package parser_test

import (
	"testing"

	"github.com/mhoertnagl/donkey/lexer"
	"github.com/mhoertnagl/donkey/parser"
)

func TestEmpty(t *testing.T) {
	test(t, "", "", 0)
}

func TestLetStatements(t *testing.T) {
	test(t, "let a = 0;", "let a = 0;", 1)
}

func TestReturnStatements(t *testing.T) {
	test(t, "return 42;", "return 42;", 1)
	test(t, "return a;", "return a;", 1)
}

func TestExpressionStatements(t *testing.T) {
	test(t, "0;", "0;", 1)
	test(t, "x;", "x;", 1)

	test(t, "-15;", "(-15);", 1)
	test(t, "!true;", "(!true);", 1)
	test(t, "~0;", "(~0);", 1)
	test(t, "~~0;", "(~(~0));", 1)
	test(t, "--99;", "(-(-99));", 1)
	test(t, "!~-a;", "(!(~(-a)));", 1)

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

	test(t, "1 - -2;", "(1 - (-2));", 1)
	test(t, "-1 - 2;", "((-1) - 2);", 1)
}

func TestOperatorPrecedence(t *testing.T) {
	test(t, "a + b + c;", "((a + b) + c);", 1)
	test(t, "a + b * c + d;", "((a + (b * c)) + d);", 1)
	test(t, "a * b + c;", "((a * b) + c);", 1)
	test(t, "a + b * c;", "(a + (b * c));", 1)
	test(t, "a + (b + c) + d;", "((a + (b + c)) + d);", 1)
	test(t, "(a + b) * c;", "((a + b) * c);", 1)
}

func TestExpressionGroupStatements(t *testing.T) {
	test(t, "(0);", "0;", 1)
	test(t, "((0));", "0;", 1)
}

func TestBlockStatements(t *testing.T) {
	test(t, "{ }", "{  }", 1)
	test(t, "{ 1; }", "{ 1; }", 1)
	test(t, "{ -1; -2; }", "{ (-1);(-2); }", 1)
	test(t, "{ 0; 1; 2; }", "{ 0;1;2; }", 1)
	test(t, "{ -0; --1; !false; }", "{ (-0);(-(-1));(!false); }", 1)
}

func TestIfStatements(t *testing.T) {
	test(t, "if 0 < 1 { return 1; }", "if (0 < 1) { return 1; }", 1)
	test(t, "if 0 < 1 && a > b { return 1; } else { return 0; }", "if ((0 < 1) && (a > b)) { return 1; } else { return 0; }", 1)
	test(t, "if 0 < 1 || x == y { return 1; }", "if ((0 < 1) || (x == y)) { return 1; }", 1)
	test(t, "if 0 < 1 { let a = b + c; return a; }", "if (0 < 1) { let a = (b + c);return a; }", 1)
	test(t, "if 0 < 1 { return 1; } else { return 0; }", "if (0 < 1) { return 1; } else { return 0; }", 1)
	test(t, "if a { return b; } else if b { return c; } else { return d; }", "if a { return b; } else if b { return c; } else { return d; }", 1)
}

// :>, <:, =>, -> >>=, =<<, |>, <|, ~>, <~, +>, <+, ::, :, #, ?:, (| |), {| |}, |{  }|, <>, ><, <|>, <+>, <->, <=>, ?, ++, --, @,

func TestFunDefn(t *testing.T) {
	test(t, "fn foo() {}", "fn foo() {  }", 1)
	test(t, "fn bar(a) {}", "fn bar(a) {  }", 1)
	test(t, "fn baz(a, b) {}", "fn baz(a, b) {  }", 1)
}

// func TestFunLiterals(t *testing.T) {
// 	test(t, "fn () {};", "fun () {  };", 1)
// 	test(t, "fn (a) {};", "fun (a) {  };", 1)
// 	test(t, "fn (a, b) {};", "fun (a, b) {  };", 1)
// }

// func TestFunCall(t *testing.T) {
// 	test(t, "foo();", "foo();", 1)
// 	test(t, "foo(a);", "foo(a);", 1)
// 	test(t, "foo(a, b);", "foo(a, b);", 1)

// 	test(t, "fn () { return 0; }();", "fun () { return 0; }();", 1)
// 	test(t, "fn (a) { return a + 1; }(1);", "fun (a) { return (a + 1); }(1);", 1)
// 	test(t, "fn (a, b) { return a + b; }(1, 2);", "fun (a, b) { return (a + b); }(1, 2);", 1)
// }

// TODO: Test error cases.

func test(t *testing.T, input string, expected string, n int) {
	t.Helper()
	lexer := lexer.NewLexer(input)
	parser := parser.NewParser(lexer)
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
