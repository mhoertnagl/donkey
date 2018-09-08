package lexer

import (
	"strings"

	"github.com/mhoertnagl/donkey/token"
)

type Lexer struct {
	input  string
	curPos int
	nxtPos int
	ch     byte
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.read()
	return l
}

func (l *Lexer) Next() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '(':
		tok = l.newToken(token.LPAR)
	case ')':
		tok = l.newToken(token.RPAR)
	case '{':
		tok = l.newToken(token.LBRA)
	case '}':
		tok = l.newToken(token.RBRA)
	case '+':
		tok = l.newToken(token.PLUS)
	case '=':
		tok = l.newToken(token.ASSIGN)
	case ',':
		tok = l.newToken(token.COMMA)
	case ';':
		tok = l.newToken(token.SCOLON)
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readID()
			tok.Typ = token.LookupId(tok.Literal)
			return tok
		} else if isDec(l.ch) {
			tok.Literal = l.readNum()
			tok.Typ = token.INT
			return tok
		} else {
			tok = l.newToken(token.ILLEGAL)
			return tok
		}
	}

	l.read()
	return tok
}

func (l *Lexer) read() {
	if l.nxtPos >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.nxtPos]
	}
	l.curPos = l.nxtPos
	l.nxtPos++
}

func (l *Lexer) newToken(typ token.TokenType) token.Token {
	return token.Token{Typ: typ, Literal: string(l.ch)}
}

// func (l *Lexer) newToken(type token.TokenType, literal string) {

// }

func (l *Lexer) skipWhitespace() {
	for isWhitespace(l.ch) {
		l.read()
	}
}

func (l *Lexer) readID() string {
	start := l.curPos
	for isLetter(l.ch) {
		l.read()
	}
	return l.input[start:l.curPos]
}

func (l *Lexer) readNum() string {
	start := l.curPos
	for isDec(l.ch) {
		l.read()
	}
	return l.input[start:l.curPos]
}

// isWhitespace returns true iff the character is one of [ \t\r\n].
func isWhitespace(c byte) bool {
	return strings.Contains(" \t\r\n", string(c))
}

// isNewline returns true iff the character is '\n'.
func isNewline(c byte) bool {
	return c == '\n'
}

// isDec returns true iff the character is a decimal digit.
func isDec(c byte) bool {
	return '0' <= c && c <= '9'
	//return strings.Contains("0123456789", string(c))
}

// isBin returns true iff the character is either '0' or '1'.
func isBin(c byte) bool {
	return c == '0' || c == '1'
	//return strings.Contains("01", string(c))
}

// isHex returns true iff the character is a hexadecimal digit. Note however,
// that the lower-case hexadecimal digits [a-f] are not supported.
func isHex(c byte) bool {
	return isDec(c) || ('A' <= c && c <= 'F')
	//return strings.Contains("0123456789ABCDEF", string(c))
}

// isLetter returns true iff the character is one of [a-zA-Z].
func isLetter(c byte) bool {
	return ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z')
	//return strings.Contains("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ", string(c))
}
