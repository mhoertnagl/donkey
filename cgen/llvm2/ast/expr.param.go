package ast

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"github.com/mhoertnagl/donkey/cgen/llvm2/ctx"
)

type ParamExpr struct {
	fun   *ctx.FuncContext
	name  string
	param *ir.Param
}

func NewParamExpr(fun *ctx.FuncContext, name string, param *ir.Param) *ParamExpr {
	return &ParamExpr{fun, name, param}
}

func (e *ParamExpr) gen() value.Value {
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
