package ast

type Stmt interface {
	gen()
}

type Stmts []Stmt

func (ss Stmts) gen() {
	for _, s := range ss {
		s.gen()
	}
}
