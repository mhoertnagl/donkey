package lexer

import (
	"testing"

	"github.com/mhoertnagl/donkey/token"
)

func TestNext(t *testing.T) {
	input := `
		let five = 5;
		let ten = 10;

		let add = fun(x, y) {
			x + y;
		};

		let result = add(five, ten);
	`

	tokens := []token.Token{
		{token.LET, "let"},
		{token.ID, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SCOLON, ";"},

		{token.LET, "let"},
		{token.ID, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SCOLON, ";"},

		{token.LET, "let"},
		{token.ID, "add"},
		{token.ASSIGN, "="},
		{token.FUN, "fun"},
		{token.LPAR, "("},
		{token.ID, "x"},
		{token.COMMA, ","},
		{token.ID, "y"},
		{token.ID, ")"},
		{token.LBRA, "{"},
		{token.ID, "x"},
		{token.PLUS, "+"},
		{token.ID, "y"},
		{token.RBRA, "}"},
		{token.SCOLON, ";"},

		{token.LET, "let"},
		{token.ID, "result"},
		{token.ASSIGN, "="},
		{token.ID, "add"},
		{token.LPAR, "("},
		{token.ID, "five"},
		{token.COMMA, ","},
		{token.ID, "ten"},
		{token.ID, ")"},
		{token.SCOLON, ";"},
		{token.EOF, ""},
	}
	test(t, input, tokens)
}

func test(t *testing.T, input string, tokens []token.Token) {
	l := NewLexer(input)
	for _, e := range tokens {
		a := l.Next()
		if a.Typ != e.Typ {
			t.Errorf("Unexpected token type [%s]. Expecting [%s].", a.Typ, e.Typ)
		}
		if a.Literal != e.Literal {
			t.Errorf("Unexpected token literal [%s]. Expecting [%s].", a.Literal, e.Literal)
		}
	}
}
