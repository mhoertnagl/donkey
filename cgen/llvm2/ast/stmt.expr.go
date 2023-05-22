package ast

import "github.com/mhoertnagl/donkey/cgen/llvm2/ctx"

type ExprStmt struct {
	fun  *ctx.FuncContext
	expr Expr
}

func NewExprStmt(fun *ctx.FuncContext, expr Expr) *ExprStmt {
	return &ExprStmt{fun, expr}
}

func (n *ExprStmt) gen() {
	n.expr.gen()
}
