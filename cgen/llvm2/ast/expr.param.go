package ast

import (
	"github.com/llir/llvm/ir"
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
