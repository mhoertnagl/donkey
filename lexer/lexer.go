package lexer

import (
	"github.com/mhoertnagl/donkey/token"
)

type Lexer struct {
}

func NewLexer() *Lexer {
	return &Lexer{}
}

func (l *Lexer) Next() token.Token {
	return token.Token{}
}
