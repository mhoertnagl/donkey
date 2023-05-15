package ast

import (
	"github.com/llir/llvm/ir/value"
	"github.com/mhoertnagl/donkey/cgen/llvm2/ctx"
)

type FunCallExpr struct {
	fun    *ctx.FuncContext
	callee Expr
	args   Exprs
}

func NewFunCallExpr(fun *ctx.FuncContext, callee Expr, args Exprs) *FunCallExpr {
	return &FunCallExpr{fun, callee, args}
}

func (e *FunCallExpr) gen() value.Value {
	callee := e.callee.gen()
	args := e.args.gen()
	blk := e.fun.GetCurrentBlock()
	return blk.NewCall(callee, args...)
}
