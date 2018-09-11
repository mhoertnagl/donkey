package parser

import (
	"fmt"
	"github.com/mhoertnagl/donkey/lexer"
	"github.com/mhoertnagl/donkey/token"
)

const (
	_ int = iota
	LOWEST
	OR          // ||
	AND         // &&
	EQUALS      // ==, !=
	LESSGREATER // >, <, <=, >=
	BOR         // |
	XOR         // ^
	BAND        // &
	SHIFT       // <<, >>, >>>, <<>, <>>
	SUM         // +, -
	PRODUCT     // *, /
	PREFIX      // -, !, ~
	CALL        // foo()
)

type prefixParslet func() Expression
type infixParslet func(Expression) Expression

type Parser struct {
	lexer          *lexer.Lexer
	curToken       token.Token
	nxtToken       token.Token
	errors         []string
	prefixParslets map[token.TokenType]prefixParslet
	infixParslets  map[token.TokenType]infixParslet
}

func NewParser(lexer *lexer.Lexer) *Parser {
	p := &Parser{lexer: lexer, errors: []string{}}
	p.next()
	p.next()
	return p
}

func (p *Parser) registerPrefix(tok token.TokenType, f prefixParslet) {
	p.prefixParslets[tok] = f
}

func (p *Parser) registerInfix(tok token.TokenType, f infixParslet) {
	p.infixParslets[tok] = f
}

func (p *Parser) next() {
	p.curToken = p.nxtToken
	p.nxtToken = p.lexer.Next()
}

func (p *Parser) curTokenIs(exp token.TokenType) bool {
	return p.curToken.Typ == exp
}

func (p *Parser) nxtTokenIs(exp token.TokenType) bool {
	return p.nxtToken.Typ == exp
}

func (p *Parser) expectNext(exp token.TokenType) bool {
	if p.nxtTokenIs(exp) {
		p.next()
		return true
	}
	p.errorNext(exp)
	return false
}

func (p *Parser) errorNext(exp token.TokenType) {
	msg := fmt.Sprintf("Expected token [%s] but got [%s].", exp, p.curToken.Typ)
	p.errors = append(p.errors, msg)
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) Parse() *Program {
	prog := &Program{}
	prog.Statements = []Statement{}
	i := 10
	for !p.curTokenIs(token.EOF) && i > 0 {
		stmt := p.parseStatement()
		if stmt != nil {
			prog.Statements = append(prog.Statements, stmt)
		}
		p.next()
		i--
	}
	return prog
}

func (p *Parser) parseStatement() Statement {
	//fmt.Printf("STMT: %s\n", p.curToken.Typ)
	switch p.curToken.Typ {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	}
	return p.parseExpressionStatement()
}

func (p *Parser) parseLetStatement() *LetStatement {
	stmt := &LetStatement{Token: p.curToken}
	if !p.expectNext(token.ID) {
		return nil
	}
	//fmt.Printf("LET: %s\n", p.curToken.Typ)
	stmt.Name = &Identifier{Token: p.curToken, Value: p.curToken.Literal}
	if !p.expectNext(token.ASSIGN) {
		return nil
	}
	//fmt.Printf("LET: %s\n", p.curToken.Typ)
	//stmt.Value = nil
	for !p.curTokenIs(token.SCOLON) {
		p.next()
		//fmt.Printf("LET: %s\n", p.curToken.Typ)
	}
	return stmt
}

func (p *Parser) parseReturnStatement() *ReturnStatement {
	stmt := &ReturnStatement{Token: p.curToken}
	p.next()
	//stmt.Value = nil
	for !p.curTokenIs(token.SCOLON) {
		p.next()
	}
	return stmt
}

func (p *Parser) parseExpressionStatement() *ExpressionStatement {
	stmt := &ExpressionStatement{Token: p.curToken}
	p.next()
	//stmt.Value = nil
	for !p.curTokenIs(token.SCOLON) {
		p.next()
	}
	return stmt
}
