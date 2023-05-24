package ast

import (
	"github.com/llir/llvm/ir/value"
	"github.com/mhoertnagl/donkey/cgen/llvm/ctx"
)

type AssignExpr struct {
	fun *ctx.FuncContext
	id  *IdentifierExpr
	// op    token.TokenType
	value Expr
}

func NewAssignExpr(fun *ctx.FuncContext, id *IdentifierExpr, value Expr) *AssignExpr {
	return &AssignExpr{fun, id, value}
}

func (e *AssignExpr) gen() value.Value {
	blk := e.fun.GetCurrentBlock()
	val := e.value.gen()
	if sym, ok := e.fun.Get(e.id.name); ok {
		blk.NewStore(val, sym.GetValue())
	}
	return val
	// 	l := e.left.gen()
	// 	r := e.right.gen()
	// case token.ASSIGN:
	// 	// TODO: l is not correct - we need the pointer not the value.
	// 	blk.NewStore(r, l)
	// 	// Return the assigned value. This way
	// 	// we can assign the same value to multiple
	// 	// variables like in a = b = c = 0;
	// 	return r

}
