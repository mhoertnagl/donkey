package lexer

import (
	"testing"
	"github.com/mhoertnagl/donkey/token"
)

func TestNext(t *testing.T) {
	input := `// A comment.
    let five = 5;
		let ten = 0x10; // Another comment.

    // Commen line 1.
    // Comment line 2.

    /**
     * Adds two numbers.
     */
		let add = fun(x, y) {
			if y > 5 && x < 2 || y != 0 {
				return x + y / 5 - 2;
			} else {
				return 0;
			}			
		};

    /* Another comment */
    
    /**
     * The result.
     * Some more comments.
     */
		let result = add(five, ten);
	`

	tokens := []token.Token{
		// let five = 5;
		{Typ: token.LET, Literal: "let", Line: 2, Col: 6},
		{Typ: token.ID, Literal: "five", Line: 2, Col: 10},
		{Typ: token.ASSIGN, Literal: "=", Line: 2, Col: 15},
		{Typ: token.INT, Literal: "5"},
		{Typ: token.SCOLON, Literal: ";"},
		// let ten = 0x10;
		{Typ: token.LET, Literal: "let"},
		{Typ: token.ID, Literal: "ten"},
		{Typ: token.ASSIGN, Literal: "="},
		{Typ: token.INT, Literal: "0x10"},
		{Typ: token.SCOLON, Literal: ";"},
		// let add = fun(x, y) {
		{Typ: token.LET, Literal: "let"},
		{Typ: token.ID, Literal: "add"},
		{Typ: token.ASSIGN, Literal: "="},
		{Typ: token.FUN, Literal: "fun"},
		{Typ: token.LPAR, Literal: "("},
		{Typ: token.ID, Literal: "x"},
		{Typ: token.COMMA, Literal: ","},
		{Typ: token.ID, Literal: "y"},
		{Typ: token.RPAR, Literal: ")"},
		{Typ: token.LBRA, Literal: "{"},
		// 	if y > 5 && x < 2 || y != 0 {
		{Typ: token.IF, Literal: "if"},
		{Typ: token.ID, Literal: "y"},
		{Typ: token.GT, Literal: ">"},
		{Typ: token.INT, Literal: "5"},
		{Typ: token.CONJ, Literal: "&&"},
		{Typ: token.ID, Literal: "x"},
		{Typ: token.LT, Literal: "<"},
		{Typ: token.INT, Literal: "2"},
		{Typ: token.DISJ, Literal: "||"},
		{Typ: token.ID, Literal: "y"},
		{Typ: token.NEQ, Literal: "!="},
		{Typ: token.INT, Literal: "0"},
		{Typ: token.LBRA, Literal: "{"},
		// return x + y / 5 - 2;
		{Typ: token.RETURN, Literal: "return"},
		{Typ: token.ID, Literal: "x"},
		{Typ: token.PLUS, Literal: "+"},
		{Typ: token.ID, Literal: "y"},
		{Typ: token.DIV, Literal: "/"},
		{Typ: token.INT, Literal: "5"},
		{Typ: token.MINUS, Literal: "-"},
		{Typ: token.INT, Literal: "2"},
		{Typ: token.SCOLON, Literal: ";"},
		// } else {
		{Typ: token.RBRA, Literal: "}"},
		{Typ: token.ELSE, Literal: "else"},
		{Typ: token.LBRA, Literal: "{"},
		// return 0;
		{Typ: token.RETURN, Literal: "return"},
		{Typ: token.INT, Literal: "0"},
		{Typ: token.SCOLON, Literal: ";"},
		// }
		{Typ: token.RBRA, Literal: "}"},
		// };
		{Typ: token.RBRA, Literal: "}"},
		{Typ: token.SCOLON, Literal: ";"},
		// let result = add(five, ten);
		{Typ: token.LET, Literal: "let"},
		{Typ: token.ID, Literal: "result"},
		{Typ: token.ASSIGN, Literal: "="},
		{Typ: token.ID, Literal: "add"},
		{Typ: token.LPAR, Literal: "("},
		{Typ: token.ID, Literal: "five"},
		{Typ: token.COMMA, Literal: ","},
		{Typ: token.ID, Literal: "ten"},
		{Typ: token.RPAR, Literal: ")"},
		{Typ: token.SCOLON, Literal: ";"},
		{Typ: token.EOF, Literal: ""},
	}
	testBlock(t, input, tokens)
}

func TestSingleTokens(t *testing.T) {
  test(t, "", token.Token{Typ: token.EOF, Literal: ""})
  
  test(t, "=", token.Token{Typ: token.ASSIGN, Literal: "="})
  test(t, "==", token.Token{Typ: token.EQU, Literal: "=="})
  
  test(t, "+", token.Token{Typ: token.PLUS, Literal: "+"})
  test(t, "-", token.Token{Typ: token.MINUS, Literal: "-"})
  test(t, "*", token.Token{Typ: token.TIMES, Literal: "*"})
  test(t, "/", token.Token{Typ: token.DIV, Literal: "/"})
  
  test(t, "~", token.Token{Typ: token.INV, Literal: "~"})
  
  test(t, "&", token.Token{Typ: token.AND, Literal: "&"})
  test(t, "&&", token.Token{Typ: token.CONJ, Literal: "&&"})
  
  test(t, "|", token.Token{Typ: token.OR, Literal: "|"})
  test(t, "||", token.Token{Typ: token.DISJ, Literal: "||"})
  
  test(t, "^", token.Token{Typ: token.XOR, Literal: "^"})
  
  test(t, "!", token.Token{Typ: token.NOT, Literal: "!"})
  test(t, "!=", token.Token{Typ: token.NEQ, Literal: "!="})
  
  test(t, ">", token.Token{Typ: token.GT, Literal: ">"})
  test(t, ">=", token.Token{Typ: token.GE, Literal: ">="})
  test(t, ">>", token.Token{Typ: token.SRL, Literal: ">>"})
  test(t, ">>>", token.Token{Typ: token.SRA, Literal: ">>>", Line: 1, Col: 1})
  
  test(t, "<", token.Token{Typ: token.LT, Literal: "<"})
  test(t, "<=", token.Token{Typ: token.LE, Literal: "<="})
  test(t, "<<", token.Token{Typ: token.SLL, Literal: "<<"})
  test(t, "<>>", token.Token{Typ: token.ROR, Literal: "<>>"})
  test(t, "<<>", token.Token{Typ: token.ROL, Literal: "<<>"})
  
  test(t, "(", token.Token{Typ: token.LPAR, Literal: "("})
  test(t, ")", token.Token{Typ: token.RPAR, Literal: ")"})
  test(t, "{", token.Token{Typ: token.LBRA, Literal: "{"})
  test(t, "}", token.Token{Typ: token.RBRA, Literal: "}"})
  test(t, ",", token.Token{Typ: token.COMMA, Literal: ","})
  test(t, ";", token.Token{Typ: token.SCOLON, Literal: ";"})
  
  test(t, "xxx", token.Token{Typ: token.ID, Literal: "xxx"})
  test(t, "x1", token.Token{Typ: token.ID, Literal: "x1"})
  test(t, "0123456789", token.Token{Typ: token.INT, Literal: "0123456789"})
  test(t, "0x0123456789ABCDEF", token.Token{Typ: token.INT, Literal: "0x0123456789ABCDEF"})
  
  test(t, "#", token.Token{Typ: token.ILLEGAL, Literal: "#"})
}


const msgErrUnexpType = "%d: Unexpected token type [%s]. Expecting [%s]."
const msgErrUnexpLiteral = "%d: Unexpected token literal [%s]. Expecting [%s]."
const msgErrLineMismatch = "%d: Line mismatch [%d]. Expecting [%d]."
const msgErrColMismatch = "%d: Column mismatch [%d]. Expecting [%d]."

func testBlock(t *testing.T, input string, tokens []token.Token) {
	l := NewLexer(input)
	for i, e := range tokens {
		a := l.Next()
		if a.Typ != e.Typ {
			t.Errorf(msgErrUnexpType, i, a.Typ, e.Typ)
		}
		if a.Literal != e.Literal {
			t.Errorf(msgErrUnexpLiteral, i, a.Literal, e.Literal)
		}
    if e.Line > 0 && a.Line != e.Line {
      t.Errorf(msgErrLineMismatch, 0, a.Line, e.Line)
    }
    if e.Col > 0 && a.Col != e.Col {
      t.Errorf(msgErrColMismatch, 0, a.Col, e.Col)
    }
	}
}

func test(t *testing.T, input string, token token.Token) {
  l := NewLexer(input)
  a := l.Next()
  if a.Typ != token.Typ {
    t.Errorf(msgErrUnexpType, 0, a.Typ, token.Typ)
  }
  if a.Literal != token.Literal {
    t.Errorf(msgErrUnexpLiteral, 0, a.Literal, token.Literal)
  }
  if token.Line > 0 && a.Line != token.Line {
    t.Errorf(msgErrLineMismatch, 0, a.Line, token.Line)
  }
  if token.Col > 0 && a.Col != token.Col {
    t.Errorf(msgErrColMismatch, 0, a.Col, token.Col)
  }
}
