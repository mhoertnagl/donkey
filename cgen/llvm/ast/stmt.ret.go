package ast

import "github.com/mhoertnagl/donkey/cgen/llvm/ctx"

type ReturnStmt struct {
	fun *ctx.FuncContext
	val Expr
}

func NewReturnStmt(fun *ctx.FuncContext, val Expr) *ReturnStmt {
	return &ReturnStmt{fun, val}
}

func (n *ReturnStmt) gen() {
	blk := n.fun.GetCurrentBlock()
	blk.NewRet(n.val.gen())
}
