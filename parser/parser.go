package parser

import (
	"fmt"
	"strconv"

	"github.com/mhoertnagl/donkey/lexer"
	"github.com/mhoertnagl/donkey/token"
)

const (
	_           int = iota
	LOWEST          // LOWEST precedence.
	OR              // ||
	AND             // &&
	EQUALS          // ==, !=
	LESSGREATER     // >, <, <=, >=
	BOR             // |
	XOR             // ^
	BAND            // &
	SHIFT           // <<, >>, >>>, <<>, <>>
	SUM             // +, -
	PRODUCT         // *, /
	PREFIX          // -, !, ~
	CALL            // foo()
)

type prefixParslet func() Expression
type infixParslet func(Expression) Expression

type Parser struct {
	lexer          *lexer.Lexer
	curToken       token.Token
	nxtToken       token.Token
	errors         []string
	precedences    map[token.TokenType]int
	prefixParslets map[token.TokenType]prefixParslet
	infixParslets  map[token.TokenType]infixParslet
}

func NewParser(lexer *lexer.Lexer) *Parser {
	p := &Parser{
		lexer:          lexer,
		errors:         []string{},
		precedences:    make(map[token.TokenType]int),
		prefixParslets: make(map[token.TokenType]prefixParslet),
		infixParslets:  make(map[token.TokenType]infixParslet),
	}

	p.registerPrecedence(token.DISJ, OR)
	p.registerPrecedence(token.CONJ, AND)
	p.registerPrecedence(token.EQU, EQUALS)
	p.registerPrecedence(token.NEQ, EQUALS)
	p.registerPrecedence(token.LT, LESSGREATER)
	p.registerPrecedence(token.LE, LESSGREATER)
	p.registerPrecedence(token.GT, LESSGREATER)
	p.registerPrecedence(token.GE, LESSGREATER)
	p.registerPrecedence(token.OR, BOR)
	p.registerPrecedence(token.XOR, XOR)
	p.registerPrecedence(token.AND, BAND)
	p.registerPrecedence(token.SLL, SHIFT)
	p.registerPrecedence(token.SRA, SHIFT)
	p.registerPrecedence(token.SRL, SHIFT)
	p.registerPrecedence(token.ROR, SHIFT)
	p.registerPrecedence(token.ROL, SHIFT)
	p.registerPrecedence(token.PLUS, SUM)
	p.registerPrecedence(token.MINUS, SUM)
	p.registerPrecedence(token.TIMES, PRODUCT)
	p.registerPrecedence(token.DIV, PRODUCT)

	p.registerPrefix(token.ID, p.parseIdentifer)
	p.registerPrefix(token.INT, p.parseInteger)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.MINUS, p.parsePrefix)
	p.registerPrefix(token.INV, p.parsePrefix)
	p.registerPrefix(token.NOT, p.parsePrefix)

	p.registerInfix(token.DISJ, p.parseBinary)
	p.registerInfix(token.CONJ, p.parseBinary)
	p.registerInfix(token.EQU, p.parseBinary)
	p.registerInfix(token.NEQ, p.parseBinary)
	p.registerInfix(token.LT, p.parseBinary)
	p.registerInfix(token.LE, p.parseBinary)
	p.registerInfix(token.GT, p.parseBinary)
	p.registerInfix(token.GE, p.parseBinary)
	p.registerInfix(token.OR, p.parseBinary)
	p.registerInfix(token.XOR, p.parseBinary)
	p.registerInfix(token.AND, p.parseBinary)
	p.registerInfix(token.SLL, p.parseBinary)
	p.registerInfix(token.SRA, p.parseBinary)
	p.registerInfix(token.SRL, p.parseBinary)
	p.registerInfix(token.ROR, p.parseBinary)
	p.registerInfix(token.ROL, p.parseBinary)
	p.registerInfix(token.PLUS, p.parseBinary)
	p.registerInfix(token.MINUS, p.parseBinary)
	p.registerInfix(token.TIMES, p.parseBinary)
	p.registerInfix(token.DIV, p.parseBinary)

	// Sets the parsers current and next tokens.
	p.next()
	p.next()
	return p
}

func (p *Parser) curTokenPrecedence() int {
	if pre, ok := p.precedences[p.curToken.Typ]; ok {
		return pre
	}
	return LOWEST
}

func (p *Parser) nxtTokenPrecedence() int {
	if pre, ok := p.precedences[p.nxtToken.Typ]; ok {
		return pre
	}
	return LOWEST
}

func (p *Parser) registerPrecedence(tok token.TokenType, pre int) {
	p.precedences[tok] = pre
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

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) error(format string, a ...interface{}) {
	p.errors = append(p.errors, fmt.Sprintf(format, a...))
}

func (p *Parser) errorNext(exp token.TokenType) {
	p.error("Expected token [%s] but got [%s].", exp, p.nxtToken.Typ)
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
	stmt.Name = &Identifier{Token: p.curToken, Value: p.curToken.Literal}
	if !p.expectNext(token.ASSIGN) {
		return nil
	}
	p.next() // Consume [=].
	stmt.Value = p.parseExpression(LOWEST)
	if p.nxtTokenIs(token.SCOLON) {
		p.next()
	}
	return stmt
}

func (p *Parser) parseReturnStatement() *ReturnStatement {
	stmt := &ReturnStatement{Token: p.curToken}
	p.next() // Consume [return].
	stmt.Value = p.parseExpression(LOWEST)
	if p.nxtTokenIs(token.SCOLON) {
		p.next()
	}
	return stmt
}

func (p *Parser) parseExpressionStatement() *ExpressionStatement {
	stmt := &ExpressionStatement{Token: p.curToken}
	stmt.Value = p.parseExpression(LOWEST)
	if p.nxtTokenIs(token.SCOLON) {
		p.next()
	}
	return stmt
}

func (p *Parser) parseExpression(pre int) Expression {
	prefix := p.prefixParslets[p.curToken.Typ]
	if prefix == nil {
		p.error("No prefix parslet found for token [%s].", p.curToken.Literal)
		return nil
	}
	left := prefix()

	for !p.nxtTokenIs(token.SCOLON) && pre < p.nxtTokenPrecedence() {
		infix := p.infixParslets[p.nxtToken.Typ]
		if infix == nil {
			return left
		}
		p.next()
		left = infix(left)
	}

	return left
}

func (p *Parser) parseIdentifer() Expression {
	expr := &Identifier{Token: p.curToken, Value: p.curToken.Literal}
	return expr
}

func (p *Parser) parseInteger() Expression {
	expr := &Integer{Token: p.curToken}
	n, err := strconv.ParseInt(p.curToken.Literal, 10, 64)
	if err != nil {
		p.errorNext(token.INT)
	}
	expr.Value = n
	return expr
}

func (p *Parser) parseBoolean() Expression {
	expr := &Boolean{Token: p.curToken}
	expr.Value = (p.curToken.Typ == token.TRUE)
	return expr
}

func (p *Parser) parsePrefix() Expression {
	expr := &PrefixExpression{Token: p.curToken}
	expr.Operator = p.curToken.Literal
	p.next() // Consume operator.
	expr.Value = p.parseExpression(PREFIX)
	return expr
}

func (p *Parser) parseBinary(left Expression) Expression {
	expr := &BinaryExpression{Token: p.curToken}
	expr.Left = left
	expr.Operator = p.curToken.Literal
	precedence := p.curTokenPrecedence()
	p.next() // Consume operator.
	expr.Right = p.parseExpression(precedence)
	return expr
}
