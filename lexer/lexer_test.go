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
			if y > 5 && x < 2 || y != 0 {
				return x + y / 5 - 2;
			} else {
				return 0;
			}			
		};

		let result = add(five, ten);
	`

	tokens := []token.Token{
		// let five = 5;
		{token.LET, "let"},
		{token.ID, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SCOLON, ";"},
		// let ten = 10;
		{token.LET, "let"},
		{token.ID, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SCOLON, ";"},
		// let add = fun(x, y) {
		{token.LET, "let"},
		{token.ID, "add"},
		{token.ASSIGN, "="},
		{token.FUN, "fun"},
		{token.LPAR, "("},
		{token.ID, "x"},
		{token.COMMA, ","},
		{token.ID, "y"},
		{token.RPAR, ")"},
		{token.LBRA, "{"},
		// 	if y > 5 && x < 2 || y != 0 {
		{token.IF, "if"},
		{token.ID, "y"},
		{token.GT, ">"},
		{token.INT, "5"},
		{token.CONJ, "&&"},
		{token.ID, "x"},
		{token.LT, "<"},
		{token.INT, "2"},
		{token.DISJ, "||"},
		{token.ID, "y"},
		{token.NEQ, "!="},
		{token.INT, "0"},
		{token.LBRA, "{"},
		// return x + y / 5 - 2;
		{token.RETURN, "return"},
		{token.ID, "x"},
		{token.PLUS, "+"},
		{token.ID, "y"},
		{token.DIV, "/"},
		{token.INT, "5"},
		{token.MINUS, "-"},
		{token.INT, "2"},
		{token.SCOLON, ";"},
		// } else {
		{token.RBRA, "}"},
		{token.ELSE, "else"},
		{token.LBRA, "{"},
		// return 0;
		{token.RETURN, "return"},
		{token.INT, "0"},
		{token.SCOLON, ";"},
		// }
		{token.RBRA, "}"},
		// };
		{token.RBRA, "}"},
		{token.SCOLON, ";"},
		// let result = add(five, ten);
		{token.LET, "let"},
		{token.ID, "result"},
		{token.ASSIGN, "="},
		{token.ID, "add"},
		{token.LPAR, "("},
		{token.ID, "five"},
		{token.COMMA, ","},
		{token.ID, "ten"},
		{token.RPAR, ")"},
		{token.SCOLON, ";"},
		{token.EOF, string(0)},
	}
	test(t, input, tokens)
}

const msgErrUnexpType = "%d: Unexpected token type [%s]. Expecting [%s]."
const msgErrUnexpLiteral = "%d: Unexpected token literal [%s]. Expecting [%s]."

func test(t *testing.T, input string, tokens []token.Token) {
	l := NewLexer(input)
	for i, e := range tokens {
		a := l.Next()
		if a.Typ != e.Typ {
			t.Errorf(msgErrUnexpType, i, a.Typ, e.Typ)
		}
		if a.Literal != e.Literal {
			t.Errorf(msgErrUnexpLiteral, i, a.Literal, e.Literal)
		}
		t.Logf("%d: %s", i, a.Literal)
	}
}
