package parser

import (
	"fmt"
	"github.com/mhoertnagl/donkey/lexer"
	"github.com/mhoertnagl/donkey/token"
)

type Parser struct {
	lexer    *lexer.Lexer
	curToken token.Token
	nxtToken token.Token
	errors   []string
}

func NewParser(lexer *lexer.Lexer) *Parser {
	p := &Parser{lexer: lexer, errors: []string{}}
	p.next()
	p.next()
	return p
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
	for p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		prog.Statements = append(prog.Statements, stmt)
	}
	return prog
}

func (p *Parser) parseStatement() Statement {
	switch p.curToken.Typ {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	}
	return nil
}

func (p *Parser) parseLetStatement() *LetStatement {
	stmt := &LetStatement{Token: p.curToken}
	if !p.expectNext(token.ID) {
		return nil
	}
	stmt.Name = &Identifier{Token: p.curToken, Value: p.curToken.Literal}
	if !p.expectNext(token.ASSIGN) {
		return nil
	}
	//stmt.Value = nil
	for !p.curTokenIs(token.SCOLON) {
		p.next()
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
