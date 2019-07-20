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
	p.registerPrecedence(token.LPAR, CALL)

	p.registerPrefix(token.ID, p.parseIdentifer)
	p.registerPrefix(token.INT, p.parseInteger)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.MINUS, p.parsePrefix)
	p.registerPrefix(token.INV, p.parsePrefix)
	p.registerPrefix(token.NOT, p.parsePrefix)
	p.registerPrefix(token.LPAR, p.parseExpressionGroup)
	p.registerPrefix(token.FUN, p.parseFunctionLiteral)

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

func (p *Parser) debug(prefix string) {
	fmt.Printf("%s: Current: %s | Next: %s\n", prefix, p.curToken, p.nxtToken)
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

func (p *Parser) HasErrors() bool {
	return len(p.errors) > 0
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
	//i := 10
	for !p.curTokenIs(token.EOF) { // && i > 0 {
		stmt := p.parseStatement()
		if stmt != nil {
			prog.Statements = append(prog.Statements, stmt)
		}
		p.next()
		//i--
	}
	return prog
}

func (p *Parser) parseStatement() Statement {
	switch p.curToken.Typ {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	case token.IF:
		return p.parseIfStatement()
	case token.LBRA:
		return p.parseBlockStatement()
	}
	return p.parseExpressionStatement()
}

func (p *Parser) parseLetStatement() *LetStatement {
	// stmt := &LetStatement{Token: p.curToken}
	// p.consume(token.LET)
	// stmt.Name = p.parseIdentifer()
	// p.consume(token.ASSIGN)
	// stmt.Value = p.parseExpression(LOWEST)
	// return stmt
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

// return <Expression>
func (p *Parser) parseReturnStatement() *ReturnStatement {
	// stmt := &ReturnStatement{Token: p.curToken}
	// p.consume(token.RETURN)
	// stmt.Value = p.parseExpression(LOWEST)
	// return stmt
	stmt := &ReturnStatement{Token: p.curToken}
	p.next() // Consume [return].
	stmt.Value = p.parseExpression(LOWEST)
	if p.nxtTokenIs(token.SCOLON) {
		p.next()
	}
	return stmt
}

// if <Expression> <Block> <ElseStatement>
func (p *Parser) parseIfStatement() *IfStatement {
	// stmt := &IfStatement{Token: p.curToken}
	// p.consume(token.IF)
	// stmt.Condition = p.parseExpression(LOWEST)
	// stmt.Consequence = p.parseBlockStatement()
	// stmt.Alternative = p.parseElseStatement()
	// return stmt

	stmt := &IfStatement{Token: p.curToken}
	p.next() // Consume [if].
	stmt.Condition = p.parseExpression(LOWEST)
	p.next() // Consume end of expression.
	stmt.Consequence = p.parseBlockStatement()
	p.next() // Consume [}].
	//p.debug("cons after")
	if p.curTokenIs(token.ELSE) {
		p.next() // Consume [else].
		//p.debug("alt before")
		if p.curTokenIs(token.IF) {
			stmt.Alternative = p.parseIfStatement()
		} else if p.curTokenIs(token.LBRA) {
			stmt.Alternative = p.parseBlockStatement()
		} else {
			p.error("Expecting block or if statement.")
			return nil
		}
		p.next() // Consume [}].
		//p.debug("alt after")
	}
	return stmt
}

// func (p *Parser) parseElseStatement() Statement {
//   if p.curTokenIs(token.ELSE) {
//     p.consume(token.ELSE)
//     if p.curTokenIs(token.IF) {
//       return p.parseIfStatement()
//     }
//     return p.parseBlockStatement()
//   }
//   return nil
// }

func (p *Parser) parseBlockStatement() *BlockStatement {
	// block := &BlockStatement{Token: p.curToken}
	// block.Statements = p.stmtSeq(token.LBRA, token.SCOLON, token.RBRA, p.parseStatement)
	// return block
	block := &BlockStatement{Token: p.curToken}
	block.Statements = []Statement{}
	p.next() // Consume [{].
	// i := 10
	//p.debug("block init")
	for !p.curTokenIs(token.RBRA) && !p.curTokenIs(token.EOF) { // && i > 0 {
		stmt := p.parseStatement()
		//p.debug("block stmt")
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.next() // Consume [;].
		//p.debug("block stmt 2")
		// i--
	}
	//p.next() // Consume [}].
	//p.debug("block exit")
	return block
}

func (p *Parser) parseExpressionStatement() *ExpressionStatement {
	// stmt := &ExpressionStatement{Token: p.curToken}
	// stmt.Value = p.parseExpression(LOWEST)
	// return stmt
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
	// id := &Identifier{Token: p.curToken, Value: p.curToken.Literal}
	// p.consume(token.ID)
	// return id
	expr := p.identifer() // &Identifier{Token: p.curToken, Value: p.curToken.Literal}
	return expr
}

func (p *Parser) identifer() *Identifier {
	// id := &Identifier{Token: p.curToken, Value: p.curToken.Literal}
	// p.consume(token.ID)
	// return id
	return &Identifier{Token: p.curToken, Value: p.curToken.Literal}
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
	// return &Boolean{Token: p.curToken, Value: p.curToken.Typ == token.TRUE}
	expr := &Boolean{Token: p.curToken}
	expr.Value = (p.curToken.Typ == token.TRUE)
	return expr
}

func (p *Parser) parsePrefix() Expression {
	expr := &PrefixExpression{Token: p.curToken, Operator: p.curToken.Literal}
	//expr.Operator = p.curToken.Literal
	p.next() // Consume operator.
	expr.Value = p.parseExpression(PREFIX)
	return expr
}

func (p *Parser) parseBinary(left Expression) Expression {
	expr := &BinaryExpression{Token: p.curToken, Operator: p.curToken.Literal}
	expr.Left = left
	//expr.Operator = p.curToken.Literal
	precedence := p.curTokenPrecedence()
	p.next() // Consume operator.
	expr.Right = p.parseExpression(precedence)
	return expr
}

func (p *Parser) parseExpressionGroup() Expression {
	// p.consume(token.LPAR)
	// expr := p.parseExpression(LOWEST)
	// p.consume(token.RPAR)
	// return expr
	p.next() // Consume [(].
	expr := p.parseExpression(LOWEST)
	if p.nxtTokenIs(token.RPAR) {
		p.next() // Consume [)].
		return expr
	}
	p.error("Missing closing parenthesis in [%s].", expr)
	return nil
}

// func (p *Parser) consume(tok token.TokenType) {
//   if p.ok {
//     if p.curTokenIs(tok) {
//       p.next()
//       return
//     }
//     p.error("Expecting [%s] but got [%s].", tok, p.curToken)
//     p.ok = false
//   }
// }

// Kontrakt: curToken zeigt vor Aufruf einer parser-subroutine auf das erste
// element in der seqeuenz der subroutine und wenn fertig dann zeigt curTokenIs
// auf das erste element für die nächste routine.

// func (p *Parser) zeroOrOne(parse func() Node) Node {
//
//   return nil
// }
//
// func (p *Parser) seq(start token.TokenType, delim token.TokenType, end token.TokenType, parse func() Node) []Node {
//   nodes := []Node{}
//   p.consume(start)
//
//   if p.curTokenIs(end) {
//     p.consume(end) // p.next()
//     return nodes
//   }
//
//   node := parse()
//   nodes = append(nodes, node)
//
//   for p.curTokenIs(delim) {
//     p.consume(delim)
//     node := parse()
//     nodes = append(nodes, node)
//   }
//
//   p.consume(end)
//   return nodes
// }
//
// func (p *Parser) stmtSeq(start token.TokenType, delim token.TokenType, end token.TokenType, parse func() Node) []Node {
//   nodes := []Node{}
//   p.consume(start)
//
//   if p.curTokenIs(end) {
//     p.consume(end)
//     return nodes
//   }
//
//   node := parse()
//   nodes = append(nodes, node)
//
//   // for p.curTokenIs(delim) {
//   //   p.consume(delim)
//   //   // The last item can have an optional delimiter.
//   //   if p.curTokenIs(end) {
//   //     p.consume(end)
//   //     return nodes
//   //   }
//   //   node := parse()
//   //   nodes = append(nodes, node)
//   // }
//
//   for p.curTokenIs(delim) && !p.nxtTokenIs(end) {
//     p.consume(delim)
//     node := parse()
//     nodes = append(nodes, node)
//   }
//
//   p.consume(delim)
//   p.consume(end)
//   return nodes
// }

func (p *Parser) parseFunctionLiteral() Expression {
	// expr := &FunctionLiteral{Token: p.curToken}
	// p.consume("fun")
	//// p.consume("(")
	// expr.Params = p.seq("(", ",", ")", p.parseParam)
	//// p.consume(")")
	// expr.Body = p.parseBlockStatement()
	// return expr
	//p.debug("fun stmt")
	expr := &FunctionLiteral{Token: p.curToken}
	if !p.expectNext(token.LPAR) {
		p.error("Missing opening parenthesis in [%s].", expr)
		return nil
	}
	//p.debug("fun args begin")
	//p.next() // Consume [(].
	expr.Params = p.parseFunctionParams()

	if !p.expectNext(token.RPAR) {
		p.error("Missing closing parenthesis in [%s].", expr)
		return nil
	}
	p.next() // Consume [)].
	//p.debug("fun args end")

	expr.Body = p.parseBlockStatement()
	return expr
}

func (p *Parser) parseFunctionParams() []*Identifier {
	params := []*Identifier{}

	if p.nxtTokenIs(token.RPAR) {
		return params
	}

	p.next() // Consume [(].

	param := p.identifer() // &Identifier{Token: p.curToken, Value: p.curToken.Literal}
	params = append(params, param)

	for p.nxtTokenIs(token.COMMA) {
		p.next()               // Consume [ID].
		p.next()               // Consume [,].
		param := p.identifer() // &Identifier{Token: p.curToken, Value: p.curToken.Literal}
		params = append(params, param)
	}

	return params
}

func (p *Parser) parseFunCall(left Expression) Expression {
	expr := &CallExpression{Token: p.curToken}
	expr.Function = left
	expr.Args = p.parseExprSeq(token.COMMA, token.RPAR)
	return expr
}

// func (p *Parser) parseExprSeq(start token.TokenType, delim token.TokenType, end token.TokenType) []Expression {
func (p *Parser) parseExprSeq(delim token.TokenType, end token.TokenType) []Expression {
	exprs := []Expression{}

	// if !p.expectNext(start) {
	// 	p.error("Missing opening [%s] in [%s].", start, "???")
	// 	return nil
	// }

	if p.nxtTokenIs(end) {
		p.next() // Consume [start].
		return exprs
	}

	p.next() // Consume [start].

	expr := p.parseExpression(LOWEST)
	exprs = append(exprs, expr)

	for p.nxtTokenIs(delim) {
		p.next() // Consume [?].
		p.next() // Consume [delim].
		expr := p.parseExpression(LOWEST)
		exprs = append(exprs, expr)
	}

	if !p.expectNext(end) {
		p.error("Missing closing [%s] in [%s].", end, "???")
		return nil
	}

	return exprs
}
