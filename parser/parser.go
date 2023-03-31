package parser

import (
	"fmt"
	"strconv"

	"github.com/mhoertnagl/donkey/lexer"
	"github.com/mhoertnagl/donkey/token"
)

// TODO: for loop
// TODO: general iteration conditions
// TODO: switch case?
// TODO: strings
// TODO: arrays
// TODO: pointers?
// TODO: structs
// TODO: list support
// TODO: tuples?
// TODO: dictionaries?
// TODO: type inference

const (
	_       int = iota
	LOWEST      // LOWEST precedence.
	OR          // ||
	AND         // &&
	EQUALS      // ==, !=
	COMPARE     // >, <, <=, >=
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
	p.registerPrecedence(token.LT, COMPARE)
	p.registerPrecedence(token.LE, COMPARE)
	p.registerPrecedence(token.GT, COMPARE)
	p.registerPrecedence(token.GE, COMPARE)
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
	p.registerPrecedence(token.LPAR, CALL)

	p.registerPrefix(token.ID, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseInteger)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.MINUS, p.parsePrefix)
	p.registerPrefix(token.INV, p.parsePrefix)
	p.registerPrefix(token.NOT, p.parsePrefix)
	p.registerPrefix(token.LPAR, p.parseExpressionGroup)
	// p.registerPrefix(token.FUN, p.parseFunctionLiteral)

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
	p.registerInfix(token.LPAR, p.parseFunCall)

	// Sets the parsers current and next tokens.
	p.next()
	p.next()
	return p
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

// func (p *Parser) debug(prefix string) {
// 	fmt.Printf("%s: Current: %s | Next: %s\n", prefix, p.curToken, p.nxtToken)
// }

func (p *Parser) curTokenPrecedence() int {
	if pre, ok := p.precedences[p.curToken.Typ]; ok {
		return pre
	}
	return LOWEST
}

func (p *Parser) next() {
	p.curToken = p.nxtToken
	p.nxtToken = p.lexer.Next()
}

func (p *Parser) curTokenIs(exp token.TokenType) bool {
	return p.curToken.Typ == exp
}

func (p *Parser) curTokenIsNot(exp token.TokenType) bool {
	return p.curToken.Typ != exp
}

func (p *Parser) curTokenIsNone(exp ...token.TokenType) bool {
	for _, e := range exp {
		if p.curTokenIs(e) {
			return false
		}
	}
	return true
}

func (p *Parser) nxtTokenIs(exp token.TokenType) bool {
	return p.nxtToken.Typ == exp
}

func (p *Parser) consume(tok token.TokenType) {
	if p.curTokenIs(tok) {
		p.next()
		return
	}
	p.error("Expecting [%s] but got [%v].", tok, p.curToken)
}

func (p *Parser) HasNoErrors() bool {
	return len(p.errors) == 0
}

func (p *Parser) HasErrors() bool {
	return len(p.errors) > 0
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) error(format string, a ...any) {
	p.errors = append(p.errors, fmt.Sprintf(format, a...))
}

func (p *Parser) Parse() *Program {
	prog := &Program{Statements: []Statement{}}
	for p.curTokenIsNot(token.EOF) {
		stmt := p.parseStatement()
		prog.Statements = append(prog.Statements, stmt)
		p.consume(token.SCOLON)
	}
	return prog
}

func (p *Parser) parseStatement() Statement {
	switch p.curToken.Typ {
	case token.LET:
		return p.parseLetStatement()
	case token.FUN:
		return p.parseFunDefStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	case token.IF:
		return p.parseIfStatement()
	case token.LBRA:
		return p.parseBlockStatement()
	}
	return p.parseExpressionStatement()
}

// let <Identifier> = <Expression>
func (p *Parser) parseLetStatement() *LetStatement {
	stmt := &LetStatement{Token: p.curToken}
	p.consume(token.LET) // [let]
	stmt.Name = p.identifier()
	p.consume(token.ASSIGN) // [=]
	stmt.Value = p.parseExpression(LOWEST)
	return stmt
}

// fn <Identifier> <FunctionParams> <BlockStatement>
func (p *Parser) parseFunDefStatement() Statement {
	stmt := &FunDefStatement{Token: p.curToken}
	p.consume(token.FUN) // [fn]
	stmt.Name = p.identifier()
	stmt.Params = p.parseFunctionParams()
	stmt.Body = p.parseBlockStatement()
	return stmt
}

// ( <Identifier>* )
func (p *Parser) parseFunctionParams() []*Identifier {
	params := []*Identifier{}
	p.consume(token.LPAR) // [(]
	if p.curTokenIs(token.RPAR) {
		p.next() // [)]
		return params
	}
	param := p.identifier()
	params = append(params, param)
	for p.curTokenIs(token.COMMA) {
		p.next() // [,]
		param := p.identifier()
		params = append(params, param)
	}
	p.consume(token.RPAR) // [)]
	return params
}

// TODO: return <nil>
// return <Expression>
func (p *Parser) parseReturnStatement() *ReturnStatement {
	stmt := &ReturnStatement{Token: p.curToken}
	p.consume(token.RETURN) // [return]
	stmt.Value = p.parseExpression(LOWEST)
	return stmt
}

// if <Expression> <BlockStatement> <ElseStatement>
func (p *Parser) parseIfStatement() *IfStatement {
	stmt := &IfStatement{Token: p.curToken}
	p.consume(token.IF) // [if]
	stmt.Condition = p.parseExpression(LOWEST)
	stmt.Consequence = p.parseBlockStatement()
	stmt.Alternative = p.parseElseStatement()
	return stmt
}

// <nil>
// else <IfStatement>
// else <BlockStatement>
func (p *Parser) parseElseStatement() Statement {
	if p.curTokenIs(token.ELSE) {
		p.consume(token.ELSE)
		if p.curTokenIs(token.IF) {
			return p.parseIfStatement()
		}
		return p.parseBlockStatement()
	}
	return nil
}

// { <Statement>* }
func (p *Parser) parseBlockStatement() *BlockStatement {
	block := &BlockStatement{Token: p.curToken, Statements: []Statement{}}
	p.consume(token.LBRA)
	// for p.curTokenIsNot(token.RBRA) && p.curTokenIsNot(token.EOF) {
	for p.curTokenIsNone(token.RBRA, token.EOF) {
		stmt := p.parseStatement()
		block.Statements = append(block.Statements, stmt)
		p.consume(token.SCOLON)
	}
	p.consume(token.RBRA)
	return block
}

// <Expression>
func (p *Parser) parseExpressionStatement() *ExpressionStatement {
	stmt := &ExpressionStatement{Token: p.curToken}
	stmt.Value = p.parseExpression(LOWEST)
	return stmt
}

func (p *Parser) parseExpression(pre int) Expression {
	prefix := p.prefixParslets[p.curToken.Typ]
	if prefix == nil {
		p.error("No prefix parslet found for token [%s].", p.curToken.Literal)
		return nil
	}
	left := prefix()
	// for p.curTokenIsNot(token.SCOLON) && p.curTokenIsNot(token.EOF) && pre < p.curTokenPrecedence() {
	for p.curTokenIsNone(token.SCOLON, token.EOF) && pre < p.curTokenPrecedence() {
		infix := p.infixParslets[p.curToken.Typ]
		if infix == nil {
			return left
		}
		left = infix(left)
	}
	return left
}

func (p *Parser) parseIdentifier() Expression {
	return p.identifier()
}

func (p *Parser) identifier() *Identifier {
	id := &Identifier{Token: p.curToken, Value: p.curToken.Literal}
	p.consume(token.ID)
	return id
}

func (p *Parser) parseInteger() Expression {
	expr := &Integer{Token: p.curToken}
	n, err := strconv.ParseInt(p.curToken.Literal, 10, 64)
	if err != nil {
		p.error("Invaild number [%s].", p.curToken.Literal)
	}
	expr.Value = n
	p.next() // Consume integer.
	return expr
}

func (p *Parser) parseBoolean() Expression {
	expr := &Boolean{Token: p.curToken, Value: p.curToken.Typ == token.TRUE}
	p.next() // Consume boolean.
	return expr
}

func (p *Parser) parsePrefix() Expression {
	expr := &PrefixExpression{Token: p.curToken, Operator: p.curToken.Typ}
	p.next() // Consume operator.
	expr.Value = p.parseExpression(PREFIX)
	return expr
}

func (p *Parser) parseBinary(left Expression) Expression {
	expr := &BinaryExpression{Token: p.curToken, Operator: p.curToken.Typ}
	expr.Left = left
	precedence := p.curTokenPrecedence()
	p.next() // Consume operator.
	expr.Right = p.parseExpression(precedence)
	return expr
}

// ( <Expression> )
func (p *Parser) parseExpressionGroup() Expression {
	p.consume(token.LPAR) // [(]
	expr := p.parseExpression(LOWEST)
	p.consume(token.RPAR) // [)]
	return expr
}

// // fun <FunctionParams> <BlockStatement>
// func (p *Parser) parseFunctionLiteral() Expression {
// 	expr := &FunctionLiteral{Token: p.curToken}
// 	p.consume(token.FUN)
// 	expr.Params = p.parseFunctionParams()
// 	expr.Body = p.parseBlockStatement()
// 	return expr
// }

// <Expression> ( <Expression>* )
func (p *Parser) parseFunCall(left Expression) Expression {
	expr := &CallExpression{Token: p.curToken}
	expr.Function = left
	expr.Args = p.parseExprSeq(token.LPAR, token.COMMA, token.RPAR)
	return expr
}

func (p *Parser) parseExprSeq(start, delim, end token.TokenType) []Expression {
	exprs := []Expression{}
	p.consume(start)
	if p.curTokenIs(end) {
		p.next() // <end>
		return exprs
	}
	expr := p.parseExpression(LOWEST)
	exprs = append(exprs, expr)
	for p.curTokenIs(delim) {
		p.next() // <delim>
		expr := p.parseExpression(LOWEST)
		exprs = append(exprs, expr)
	}
	p.consume(end)
	return exprs
}
