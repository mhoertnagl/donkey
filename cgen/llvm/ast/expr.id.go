package ast

import (
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"github.com/mhoertnagl/donkey/cgen/llvm/ctx"
)

type IdentifierExpr struct {
	fun  *ctx.FuncContext
	name string
}

func NewIdentifierExpr(fun *ctx.FuncContext, name string) *IdentifierExpr {
	return &IdentifierExpr{fun, name}
}

func (e *IdentifierExpr) gen() value.Value {
	sym, _ := e.fun.Get(e.name)
	switch sym := sym.(type) {
	case *ctx.ValueSymbol:
		blk := e.fun.GetCurrentBlock()
		return blk.NewLoad(types.I64, sym.GetValue())
	case *ctx.FuncSymbol:
		return sym.GetValue()
	}
	return nil
}
