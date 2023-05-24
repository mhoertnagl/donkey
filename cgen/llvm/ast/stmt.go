package ast

import "github.com/mhoertnagl/donkey/utils"

type Stmt interface {
	gen()
}

type Stmts []Stmt

func (ss Stmts) gen() {
	// for _, s := range ss {
	// 	s.gen()
	// }
	utils.For(ss, Stmt.gen)
}
