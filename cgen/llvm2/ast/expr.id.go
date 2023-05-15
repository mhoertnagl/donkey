package ast

import (
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"github.com/mhoertnagl/donkey/cgen/llvm2/ctx"
)

type IdentifierExpr struct {
	fun *ctx.FuncContext
	val string
}

func NewIdentifierExpr(fun *ctx.FuncContext, val string) *IdentifierExpr {
	return &IdentifierExpr{fun, val}
}

func (e *IdentifierExpr) gen() value.Value {
	sym, _ := e.fun.Get(e.val)
	switch sym := sym.(type) {
	case *ctx.ValueSymbol:
		blk := e.fun.GetCurrentBlock()
		return blk.NewLoad(types.I64, sym.GetValue())
	case *ctx.FuncSymbol:
		return sym.GetValue()
	}
	return nil
}
