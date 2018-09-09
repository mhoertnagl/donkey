package ast

import "github.com/mhoertnagl/donkey/token"

type Node interface {
	Literal() string
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

type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (s *LetStatement) statement()      {}
func (s *LetStatement) Literal() string { return s.Token.Literal }

type Identifier struct {
	Token token.Token
	Value string
}

func (s *Identifier) statement()      {}
func (s *Identifier) Literal() string { return s.Token.Literal }
