package parser

import (
	"bytes"
	"fmt"
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
		buf.WriteString(";")
	}
	return buf.String()
}

type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
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
	//buf.WriteString(";")
	return buf.String()
}

type ReturnStatement struct {
	Token token.Token
	Value Expression
}

func (s *ReturnStatement) statement()      {}
func (s *ReturnStatement) Literal() string { return s.Token.Literal }

func (s *ReturnStatement) String() string {
	var buf bytes.Buffer
	buf.WriteString("return ")
	// TODO: Remove when expression parsing is in place.
	if s.Value != nil {
		buf.WriteString(s.Value.String())
	}
	//buf.WriteString(";")
	return buf.String()
}

type ExpressionStatement struct {
	Token token.Token
	Value Expression
}

func (s *ExpressionStatement) statement()      {}
func (s *ExpressionStatement) Literal() string { return s.Token.Literal }

func (s *ExpressionStatement) String() string {
	var buf bytes.Buffer
	// TODO: Remove when expression parsing is in place.
	if s.Value != nil {
		buf.WriteString(s.Value.String())
	}
	//buf.WriteString(";")
	return buf.String()
}

type Identifier struct {
	Token token.Token
	Value string
}

func (e *Identifier) expression()     {}
func (e *Identifier) Literal() string { return e.Token.Literal }
func (e *Identifier) String() string  { return e.Value }

type Integer struct {
	Token token.Token
	Value int64
}

func (e *Integer) expression()     {}
func (e *Integer) Literal() string { return e.Token.Literal }
func (e *Integer) String() string  { return fmt.Sprintf("%d", e.Value) }

type Boolean struct {
	Token token.Token
	Value bool
}

func (e *Boolean) expression()     {}
func (e *Boolean) Literal() string { return e.Token.Literal }
func (e *Boolean) String() string  { return fmt.Sprintf("%t", e.Value) }

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Value    Expression
}

func (e *PrefixExpression) expression()     {}
func (e *PrefixExpression) Literal() string { return e.Token.Literal }
func (e *PrefixExpression) String() string {
	var buf bytes.Buffer
	buf.WriteString("(")
	buf.WriteString(e.Operator)
	buf.WriteString(e.Value.String())
	buf.WriteString(")")
	return buf.String()
}

type BinaryExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (e *BinaryExpression) expression()     {}
func (e *BinaryExpression) Literal() string { return e.Token.Literal }
func (e *BinaryExpression) String() string {
	var buf bytes.Buffer
	buf.WriteString("(")
	buf.WriteString(e.Left.String())
	buf.WriteString(" ")
	buf.WriteString(e.Operator)
	buf.WriteString(" ")
	buf.WriteString(e.Right.String())
	buf.WriteString(")")
	return buf.String()
}
