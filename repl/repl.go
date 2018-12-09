package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/mhoertnagl/donkey/lexer"
	"github.com/mhoertnagl/donkey/token"
)

func Start(in io.Reader, out io.Writer) {
	s := bufio.NewScanner(in)
	for {
		fmt.Fprintf(out, ">> ")
		if ok := s.Scan(); !ok {
			return
		}
		input := s.Text()
		if input == ":exit" {
			fmt.Fprintf(out, "Bye.\n")
			return
		}
		lexer := lexer.NewLexer(input)
		tok := lexer.Next()
		for tok.Typ != token.EOF && tok.Typ != token.ILLEGAL {
			fmt.Fprintf(out, "%s [%s]\n", tok.Literal, tok.Typ)
			tok = lexer.Next()
		}
	}
}
