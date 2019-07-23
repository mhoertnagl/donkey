package lexer

import (
  //"fmt"
	"github.com/mhoertnagl/donkey/token"
)

// TODO: runes support.
// TODO: track positional information.
// TODO: turn into a library.

type Lexer struct {
	input string
  len   int     // Input length.
	pos   int
  line  int     // Token line number.
  col   int     // Token column.
	ch    byte
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input, len: len(input), pos: -1, line: 1, col: 1}
	l.read()
	return l
}

func (l *Lexer) Next() token.Token {
	var tok token.Token
  
  // Skip over any sequence of whitespace, single- or multiline comments.
  for isWhitespace(l.ch) || l.peeksIs("//") || l.peeksIs("/*") {
    l.skipWhitespace()
    l.skipSingleLineComment("//")  
    l.skipMultiLineComment("/*", "*/")
  }
	
	switch {
	case l.ch == 0:
		tok = l.emit2(token.EOF, "")
	case l.peeksIs("=="):
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
	case l.peeksIs("&&"):
    l.read()
    tok = l.emit2(token.CONJ, "&&")
  case l.ch == '&':
    tok = l.emit(token.AND)
	case l.peeksIs("||"):
    l.read()
    tok = l.emit2(token.DISJ, "||")
  case l.ch == '|':
		tok = l.emit(token.OR)
	case l.ch == '^':
		tok = l.emit(token.XOR)
	case l.peeksIs("!="):
    l.read()
    tok = l.emit2(token.NEQ, "!=")
  case l.ch == '!':
		tok = l.emit(token.NOT)
  case l.peeksIs("<<>"):
    l.read()
    l.read()
    tok = l.emit2(token.ROL, "<<>")  
  case l.peeksIs("<>>"):
    l.read()
    l.read()
    tok = l.emit2(token.ROR, "<>>")
  case l.peeksIs("<<"):
    l.read()
    tok = l.emit2(token.SLL, "<<")   
  case l.peeksIs("<="):
    l.read()
    tok = l.emit2(token.LE, "<=")
  case l.ch == '<':
    tok = l.emit2(token.LT, "<")     
  case l.peeksIs(">>>"):
    l.read()
    l.read()
    tok = l.emit2(token.SRA, ">>>")
  case l.peeksIs(">>"):
    l.read()
    tok = l.emit2(token.SRL, ">>")
  case l.peeksIs(">="):
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
	case isAlpha(l.ch):
		return l.readID()
  case l.peeksIs("0x"):
    return l.readHex()
	case isDec(l.ch):
		return l.readDec()
	default:
		return l.emit(token.ILLEGAL)
	}

	l.read()
	return tok
}

func (l *Lexer) read() {
	l.ch = l.peek()
  l.pos++
  if l.ch == '\n' {
    l.col = 1
    l.line++
  } else {
    l.col++
  }
}

func (l *Lexer) readWhile(pred func(byte)bool) {
  for pred(l.ch) {
    l.read()
  }
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

func (l *Lexer) peeksIs(pattern string) bool {
  return l.peeks(uint(len(pattern))) == pattern
}

func (l *Lexer) peeksIsNot(pattern string) bool {
  return l.peeksIs(pattern) == false
}

func (l *Lexer) emit(typ token.TokenType) token.Token {
	return l.emit2(typ, string(l.ch))
}

func (l *Lexer) emit2(typ token.TokenType, literal string) token.Token {
	return token.Token{
    Typ: typ, 
    Literal: literal, 
    Line: l.line, 
    Col: l.col - len(literal),
  }
}

func (l *Lexer) skipWhitespace() {
  l.readWhile(isWhitespace)
}

func (l *Lexer) skipSingleLineComment(start string) {
  if l.peeksIs(start) {
    for l.ch != '\n' && l.ch != 0 {
      l.read()
    }
  }
}

func (l *Lexer) skipMultiLineComment(start string, end string) {
  if l.peeksIs(start) {
    for l.peeksIsNot(end) && l.ch != 0 {
      l.read()
    }    
    l.read() // [*]
    l.read() // [/]
  }
}

func (l *Lexer) readID() token.Token {
	start := l.pos
  l.read() // [a-zA-Z]
  l.readWhile(isAlphaNum)
  literal := l.input[start:l.pos]
  typ := token.LookupId(literal)
  return l.emit2(typ, literal)
}

func (l *Lexer) readDec() token.Token {
	start := l.pos
  l.readWhile(isDec)
	return l.emit2(token.INT, l.input[start:l.pos])
}

func (l *Lexer) readHex() token.Token {
  start := l.pos
  l.read() // [0]
  l.read() // [x]
  l.readWhile(isHex)
	return l.emit2(token.INT, l.input[start:l.pos])
}

// isWhitespace returns true iff the character is one of [ \t\r\n].
func isWhitespace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\r' || c == '\n'
}

// isDec returns true iff the character is a decimal digit.
func isDec(c byte) bool {
	return '0' <= c && c <= '9'
}

// isHex returns true iff the character is a hexadecimal digit.
func isHex(c byte) bool {
	return isDec(c) || ('a' <= c && c <= 'f') || ('A' <= c && c <= 'F')
}

// isAlpha returns true iff the character is one of [a-zA-Z].
func isAlpha(c byte) bool {
	return ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z')
}

// isAlphaNum returns true iff the character is one of [a-zA-Z0-9].
func isAlphaNum(c byte) bool {
	return isAlpha(c) || isDec(c) 
}
