package parser

import (
	"fmt"
	"github.com/mhoertnagl/donkey/lexer"
	"github.com/mhoertnagl/donkey/token"
	"strconv"
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
	p := &Parser{
		lexer:          lexer,
		errors:         []string{},
		prefixParslets: make(map[token.TokenType]prefixParslet),
		infixParslets:  make(map[token.TokenType]infixParslet),
	}

	p.registerPrefix(token.ID, p.parseIdentifer)
	p.registerPrefix(token.INT, p.parseInteger)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.MINUS, p.parsePrefix)
	p.registerPrefix(token.INV, p.parsePrefix)
	p.registerPrefix(token.NOT, p.parsePrefix)

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
	msg := fmt.Sprintf("Expected token [%s] but got [%s].", exp, p.nxtToken.Typ)
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
		//fmt.Printf("STMT: %s :: %s\n", p.curToken.Typ, p.nxtToken.Typ)
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
	//fmt.Printf("LET: %s :: %s\n", p.curToken.Typ, p.nxtToken.Typ)
	stmt.Name = &Identifier{Token: p.curToken, Value: p.curToken.Literal}
	if !p.expectNext(token.ASSIGN) {
		return nil
	}
	p.next() // Skip =
	//fmt.Printf("LET: %s\n", p.curToken.Typ)
	stmt.Value = p.parseExpression(LOWEST)
	// if p.nxtTokenIs(token.SCOLON) {
	// 	p.next()
	// }
	return stmt
}

func (p *Parser) parseReturnStatement() *ReturnStatement {
	stmt := &ReturnStatement{Token: p.curToken}
	p.next()
	stmt.Value = p.parseExpression(LOWEST)
	// if p.nxtTokenIs(token.SCOLON) {
	// 	p.next()
	// }
	return stmt
}

func (p *Parser) parseExpressionStatement() *ExpressionStatement {
	stmt := &ExpressionStatement{Token: p.curToken}
	stmt.Value = p.parseExpression(LOWEST)
	//fmt.Printf("STMT: %s :: %s\n", p.curToken.Typ, p.nxtToken.Typ)
	// if p.nxtTokenIs(token.SCOLON) {
	// 	p.next()
	// }
	return stmt
}

func (p *Parser) parseExpression(pre int) Expression {
	prefix := p.prefixParslets[p.curToken.Typ]
	if prefix == nil {
		msg := fmt.Sprintf("No prefix parslet found for token [%s].", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	left := prefix()
	return left
}

func (p *Parser) parseIdentifer() Expression {
	expr := &Identifier{Token: p.curToken, Value: p.curToken.Literal}
	p.next()
	return expr
}

func (p *Parser) parseInteger() Expression {
	n, err := strconv.ParseInt(p.curToken.Literal, 10, 64)
	if err != nil {
		p.errorNext(token.INT)
	}
	p.next()
	return &Integer{Token: p.curToken, Value: n}
}

func (p *Parser) parseBoolean() Expression {
	expr := &Boolean{Token: p.curToken}
	if p.curToken.Typ == token.TRUE {
		expr.Value = true
	} else {
		expr.Value = false
	}
	p.next()
	return expr
}

func (p *Parser) parsePrefix() Expression {
	//fmt.Printf("PRE1: %s :: %s\n", p.curToken.Typ, p.nxtToken.Typ)
	expr := &PrefixExpression{Token: p.curToken}
	expr.Operator = p.curToken.Literal
	p.next()
	//fmt.Printf("PRE2: %s :: %s\n", p.curToken.Typ, p.nxtToken.Typ)
	expr.Value = p.parseExpression(PREFIX)
	//fmt.Printf("PRE3: %s :: %s\n", p.curToken.Typ, p.nxtToken.Typ)
	p.next()
	//fmt.Printf("PRE4: %s :: %s\n", p.curToken.Typ, p.nxtToken.Typ)
	return expr
}

func (p *Parser) parseBinary(left Expression) Expression {
	//fmt.Printf("PRE1: %s :: %s\n", p.curToken.Typ, p.nxtToken.Typ)
	expr := &BinaryExpression{Token: p.curToken}
	expr.Left = left
	expr.Operator = p.curToken.Literal
	p.next()
	//fmt.Printf("PRE2: %s :: %s\n", p.curToken.Typ, p.nxtToken.Typ)
	expr.Right = p.parseExpression(PREFIX)
	//fmt.Printf("PRE3: %s :: %s\n", p.curToken.Typ, p.nxtToken.Typ)
	p.next()
	//fmt.Printf("PRE4: %s :: %s\n", p.curToken.Typ, p.nxtToken.Typ)
	return expr
}
