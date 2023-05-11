package parser

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/mhoertnagl/donkey/token"
)

type Node interface {
	Literal() string
	String() string
}

type Statement interface {
	Node
	statement()
}

type Expression interface {
	Node
	expression()
}

type Program struct {
	Statements []Statement
}

func NewProgram() *Program {
	return &Program{Statements: []Statement{}}
}

func (p *Program) Literal() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].Literal()
	}
	return ""
}

func (p *Program) String() string {
	var buf bytes.Buffer
	for _, s := range p.Statements {
		buf.WriteString(s.String())
		//buf.WriteString(";")
	}
	return buf.String()
}

type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func NewLetStmt(token token.Token) *LetStatement {
	return &LetStatement{Token: token}
}

func (s *LetStatement) statement()      {}
func (s *LetStatement) Literal() string { return s.Token.Literal }
func (s *LetStatement) String() string {
	var buf bytes.Buffer
	buf.WriteString("let ")
	buf.WriteString(s.Name.String())
	buf.WriteString(" = ")
	// TODO: Remove when expression parsing is in place.
	if s.Value != nil {
		buf.WriteString(s.Value.String())
	}
	buf.WriteString(";")
	return buf.String()
}

type FunDefStatement struct {
	Token  token.Token
	Name   *Identifier
	Params []*Identifier
	Body   *BlockStatement
}

func NewFunDefStmt(token token.Token) *FunDefStatement {
	return &FunDefStatement{Token: token}
}

func (e *FunDefStatement) statement()      {}
func (e *FunDefStatement) Literal() string { return e.Token.Literal }
func (e *FunDefStatement) String() string {
	params := []string{}
	for _, id := range e.Params {
		params = append(params, id.String())
	}

	var buf bytes.Buffer
	buf.WriteString("fn")
	buf.WriteString(" ")
	buf.WriteString(e.Name.String())
	buf.WriteString("(")
	buf.WriteString(strings.Join(params, ", "))
	buf.WriteString(")")
	buf.WriteString(" ")
	buf.WriteString(e.Body.String())
	return buf.String()
}

type ReturnStatement struct {
	Token token.Token
	Value Expression
}

func NewReturnStmt(token token.Token) *ReturnStatement {
	return &ReturnStatement{Token: token}
}

func (s *ReturnStatement) statement()      {}
func (s *ReturnStatement) Literal() string { return s.Token.Literal }
func (s *ReturnStatement) String() string {
	var buf bytes.Buffer
	buf.WriteString("return ")
	buf.WriteString(s.Value.String())
	buf.WriteString(";")
	return buf.String()
}

type IfStatement struct {
	Token       token.Token
	Condition   Expression
	Consequence Statement
	Alternative Statement
}

func NewIfStmt(token token.Token) *IfStatement {
	return &IfStatement{Token: token}
}

func (s *IfStatement) statement()      {}
func (s *IfStatement) Literal() string { return s.Token.Literal }
func (s *IfStatement) String() string {
	var buf bytes.Buffer
	buf.WriteString("if")
	buf.WriteString(" ")
	buf.WriteString(s.Condition.String())
	buf.WriteString(" ")
	buf.WriteString(s.Consequence.String())
	if s.Alternative != nil {
		buf.WriteString(" ")
		buf.WriteString("else")
		buf.WriteString(" ")
		buf.WriteString(s.Alternative.String())
	}
	return buf.String()
}

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func NewBlockStmt(token token.Token) *BlockStatement {
	return &BlockStatement{Token: token, Statements: []Statement{}}
}

func (s *BlockStatement) statement()      {}
func (s *BlockStatement) Literal() string { return s.Token.Literal }
func (s *BlockStatement) String() string {
	var buf bytes.Buffer
	buf.WriteString("{")
	buf.WriteString(" ")
	for _, stmt := range s.Statements {
		buf.WriteString(stmt.String())
	}
	buf.WriteString(" ")
	buf.WriteString("}")
	return buf.String()
}

type ExpressionStatement struct {
	Token token.Token
	Value Expression
}

func NewExprStmt(token token.Token) *ExpressionStatement {
	return &ExpressionStatement{Token: token}
}

func (s *ExpressionStatement) statement()      {}
func (s *ExpressionStatement) Literal() string { return s.Token.Literal }
func (s *ExpressionStatement) String() string {
	var buf bytes.Buffer
	// TODO: Remove when expression parsing is in place.
	if s.Value != nil {
		buf.WriteString(s.Value.String())
	}
	buf.WriteString(";")
	return buf.String()
}

type Identifier struct {
	Token token.Token
	Value string
}

func NewIdentifier(token token.Token) *Identifier {
	return &Identifier{Token: token, Value: token.Literal}
}

func (e *Identifier) expression()     {}
func (e *Identifier) Literal() string { return e.Token.Literal }
func (e *Identifier) String() string  { return e.Value }

type Integer struct {
	Token token.Token
	Value int64
}

func NewIntLiteral(token token.Token) *Integer {
	return &Integer{Token: token}
}

func (e *Integer) expression()     {}
func (e *Integer) Literal() string { return e.Token.Literal }
func (e *Integer) String() string  { return fmt.Sprintf("%d", e.Value) }

type Boolean struct {
	Token token.Token
	Value bool
}

func NewBoolLiteral(_token token.Token) *Boolean {
	return &Boolean{Token: _token, Value: _token.Typ == token.TRUE}
}

func (e *Boolean) expression()     {}
func (e *Boolean) Literal() string { return e.Token.Literal }
func (e *Boolean) String() string  { return fmt.Sprintf("%t", e.Value) }

type PrefixExpression struct {
	Token    token.Token
	Operator token.TokenType
	Value    Expression
}

func NewPrefixExpr(token token.Token) *PrefixExpression {
	return &PrefixExpression{Token: token, Operator: token.Typ}
}

func (e *PrefixExpression) expression()     {}
func (e *PrefixExpression) Literal() string { return e.Token.Literal }
func (e *PrefixExpression) String() string {
	var buf bytes.Buffer
	buf.WriteString("(")
	buf.WriteString(string(e.Operator))
	buf.WriteString(e.Value.String())
	buf.WriteString(")")
	return buf.String()
}

type BinaryExpression struct {
	Token    token.Token
	Left     Expression
	Operator token.TokenType
	Right    Expression
}

func NewBinaryExpr(token token.Token) *BinaryExpression {
	return &BinaryExpression{Token: token, Operator: token.Typ}
}

func (e *BinaryExpression) expression()     {}
func (e *BinaryExpression) Literal() string { return e.Token.Literal }
func (e *BinaryExpression) String() string {
	var buf bytes.Buffer
	buf.WriteString("(")
	buf.WriteString(e.Left.String())
	buf.WriteString(" ")
	buf.WriteString(string(e.Operator))
	buf.WriteString(" ")
	buf.WriteString(e.Right.String())
	buf.WriteString(")")
	return buf.String()
}

// type FunctionLiteral struct {
// 	Token  token.Token
// 	Params []*Identifier
// 	Body   *BlockStatement
// }

// func (e *FunctionLiteral) expression()     {}
// func (e *FunctionLiteral) Literal() string { return e.Token.Literal }
// func (e *FunctionLiteral) String() string {
// 	params := []string{}
// 	for _, id := range e.Params {
// 		params = append(params, id.String())
// 	}

// 	var buf bytes.Buffer
// 	buf.WriteString("fun")
// 	buf.WriteString(" ")
// 	buf.WriteString("(")
// 	buf.WriteString(strings.Join(params, ", "))
// 	buf.WriteString(")")
// 	buf.WriteString(" ")
// 	buf.WriteString(e.Body.String())
// 	return buf.String()
// }

type CallExpression struct {
	Token    token.Token
	Function Expression
	Args     []Expression
}

func NewCallExpr(token token.Token) *CallExpression {
	return &CallExpression{Token: token}
}

func (e *CallExpression) expression()     {}
func (e *CallExpression) Literal() string { return e.Token.Literal }
func (e *CallExpression) String() string {
	args := []string{}
	for _, arg := range e.Args {
		args = append(args, arg.String())
	}

	var buf bytes.Buffer
	buf.WriteString(e.Function.String())
	// buf.WriteString(" ")
	buf.WriteString("(")
	buf.WriteString(strings.Join(args, ", "))
	buf.WriteString(")")
	return buf.String()
}
