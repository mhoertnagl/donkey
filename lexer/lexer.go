package lexer

import (
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
	// TODO: skip comments.

	switch l.ch {
	case 0:
		tok = l.newToken(token.EOF)
	case '=':
		if l.peek() == '=' {
			l.read()
			tok = l.newToken2(token.EQU, "==")
		} else {
			tok = l.newToken(token.ASSIGN)
		}
	case '+':
		tok = l.newToken(token.PLUS)
	case '-':
		tok = l.newToken(token.MINUS)
	case '*':
		tok = l.newToken(token.TIMES)
	case '/':
		tok = l.newToken(token.DIV)
	case '~':
		tok = l.newToken(token.INV)
	case '&':
		if l.peek() == '&' {
			l.read()
			tok = l.newToken2(token.CONJ, "&&")
		} else {
			tok = l.newToken(token.AND)
		}
	case '|':
		if l.peek() == '|' {
			l.read()
			tok = l.newToken2(token.DISJ, "||")
		} else {
			tok = l.newToken(token.OR)
		}
	case '^':
		tok = l.newToken(token.XOR)
		// SLL     = "<<"
		// SRL     = ">>"
		// SRA     = ">>>"
		// ROL     = "<<>"
		// ROR     = "<>>"
	case '!':
		if l.peek() == '=' {
			l.read()
			tok = l.newToken2(token.NEQ, "!=")
		} else {
			tok = l.newToken(token.NOT)
		}
	case '<':
		switch l.peek() {
		case '=':
			l.read()
			tok = l.newToken2(token.LE, "<=")
		case '<':
			l.read()
			tok = l.newToken2(token.SLL, "<<")
			// TODO: <>> ROR
			// TODO: <<> ROL
		default:
			tok = l.newToken(token.LT)
		}
	case '>':
		switch l.peek() {
		case '=':
			l.read()
			tok = l.newToken2(token.GE, ">=")
		case '>':
			l.read()
			tok = l.newToken2(token.SRL, ">>")
			// TODO: >>> SRA
		default:
			tok = l.newToken(token.GT)
		}
	case '(':
		tok = l.newToken(token.LPAR)
	case ')':
		tok = l.newToken(token.RPAR)
	case '{':
		tok = l.newToken(token.LBRA)
	case '}':
		tok = l.newToken(token.RBRA)
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
	l.ch = l.peek()
	l.curPos = l.nxtPos
	l.nxtPos++
}

func (l *Lexer) peek() byte {
	if l.nxtPos >= len(l.input) {
		return 0
	}
	return l.input[l.nxtPos]
}

func (l *Lexer) newToken(typ token.TokenType) token.Token {
	return token.Token{Typ: typ, Literal: string(l.ch)}
}

func (l *Lexer) newToken2(typ token.TokenType, literal string) token.Token {
	return token.Token{Typ: typ, Literal: literal}
}

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
	return c == ' ' || c == '\t' || c == '\r' || c == '\n'
}

// isNewline returns true iff the character is '\n'.
func isNewline(c byte) bool {
	return c == '\n'
}

// isDec returns true iff the character is a decimal digit.
func isDec(c byte) bool {
	return '0' <= c && c <= '9'
}

// isBin returns true iff the character is either '0' or '1'.
func isBin(c byte) bool {
	return c == '0' || c == '1'
}

// isHex returns true iff the character is a hexadecimal digit. Note however,
// that the lower-case hexadecimal digits [a-f] are not supported.
func isHex(c byte) bool {
	return isDec(c) || ('A' <= c && c <= 'F')
}

// isLetter returns true iff the character is one of [a-zA-Z].
func isLetter(c byte) bool {
	return ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z')
}
