package lexer

import (
  //"fmt"
	"github.com/mhoertnagl/donkey/token"
)

// TODO: runes support.
// TODO: skip multi line comments.
// TODO: track positional information.
// TODO: turn into a library.

type Lexer struct {
	input string
  len   int
	pos   int
	ch    byte
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input, len: len(input), pos: -1}
	l.read()
	return l
}

func (l *Lexer) Next() token.Token {
	var tok token.Token

  l.skipWhitespace()
  l.skipSingleLineComments("//")
	
	switch {
	case l.ch == 0:
		tok = l.emit2(token.EOF, "")
	case l.peeks(2) == "==":
    l.read()
    tok = l.emit2(token.EQU, "==")
  case l.ch == '=':
		tok = l.emit(token.ASSIGN)
	case l.ch == '+':
		tok = l.emit(token.PLUS)
	case l.ch == '-':
		tok = l.emit(token.MINUS)
	case l.ch == '*':
		tok = l.emit(token.TIMES)
	case l.ch == '/':
		tok = l.emit(token.DIV)
	case l.ch == '~':
		tok = l.emit(token.INV)
	case l.peeks(2) == "&&":
    l.read()
    tok = l.emit2(token.CONJ, "&&")
  case l.ch == '&':
    tok = l.emit(token.AND)
	case l.peeks(2) == "||":
    l.read()
    tok = l.emit2(token.DISJ, "||")
  case l.ch == '|':
		tok = l.emit(token.OR)
	case l.ch == '^':
		tok = l.emit(token.XOR)
	case l.peeks(2) == "!=":
    l.read()
    tok = l.emit2(token.NEQ, "!=")
  case l.ch == '!':
		tok = l.emit(token.NOT)
  case l.peeks(3) == "<<>":
    l.read()
    l.read()
    tok = l.emit2(token.ROL, "<<>")  
  case l.peeks(3) == "<>>":
    l.read()
    l.read()
    tok = l.emit2(token.ROR, "<>>")
  case l.peeks(2) == "<<":
    l.read()
    tok = l.emit2(token.SLL, "<<")   
  case l.peeks(2) == "<=":
    l.read()
    tok = l.emit2(token.LE, "<=")
  case l.ch == '<':
    tok = l.emit2(token.LT, "<")     
  case l.peeks(3) == ">>>":
    l.read()
    l.read()
    tok = l.emit2(token.SRA, ">>>")
  case l.peeks(2) == ">>":
    l.read()
    tok = l.emit2(token.SRL, ">>")
  case l.peeks(2) == ">=":
    l.read()
    tok = l.emit2(token.GE, ">=")
	case l.ch == '>':
    tok = l.emit(token.GT)
	case l.ch == '(':
		tok = l.emit(token.LPAR)
	case l.ch == ')':
		tok = l.emit(token.RPAR)
	case l.ch == '{':
		tok = l.emit(token.LBRA)
	case l.ch == '}':
		tok = l.emit(token.RBRA)
	case l.ch == ',':
		tok = l.emit(token.COMMA)
	case l.ch == ';':
		tok = l.emit(token.SCOLON)
	case isLetter(l.ch):
    // TODO: return token?
		tok.Literal = l.readID()
		tok.Typ = token.LookupId(tok.Literal)
		return tok
	case isDec(l.ch):
    // TODO: return token?
		tok.Literal = l.readNum()
		tok.Typ = token.INT
		return tok
	default:
		tok = l.emit(token.ILLEGAL)
		return tok
	}

	l.read()
	return tok
}

func (l *Lexer) read() {
	l.ch = l.peek()
  l.pos++
}

func (l *Lexer) peek() byte {
  return l.peekAt(1)
}

func (l *Lexer) peekAt(n uint) byte {
  posAt := l.pos + int(n)
  if posAt >= l.len {
    return 0
  }
  return l.input[posAt]
}

func (l *Lexer) peeks(n uint) string {
  posAt := l.pos + int(n)
  if posAt > l.len {
    return ""
  }
  return l.input[l.pos:posAt]
}

func (l *Lexer) emit(typ token.TokenType) token.Token {
	return l.emit2(typ, string(l.ch))
}

func (l *Lexer) emit2(typ token.TokenType, literal string) token.Token {
	return token.Token{Typ: typ, Literal: literal}
}

func (l *Lexer) skipWhitespace() {
	for isWhitespace(l.ch) {
		l.read()
	}
}

func (l *Lexer) skipSingleLineComments(marker string) {
  if l.peeks(uint(len(marker))) == marker {
    for l.ch != '\n' && l.ch != 0 {
      l.read()
    }
    l.skipWhitespace()
  }
}

func (l *Lexer) readID() string {
	start := l.pos
	for isLetter(l.ch) {
		l.read()
	}
	return l.input[start:l.pos]
}

func (l *Lexer) readNum() string {
	start := l.pos
	for isDec(l.ch) {
		l.read()
	}
	return l.input[start:l.pos]
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
