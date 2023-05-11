package ast

import (
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"github.com/mhoertnagl/donkey/cgen/llvm2/ctx"
)

type IntLiteralExpr struct {
	fun *ctx.FuncContext
	val int64
}

func NewIntLiteralExpr(fun *ctx.FuncContext, val int64) *IntLiteralExpr {
	return &IntLiteralExpr{fun, val}
}

func (e *IntLiteralExpr) gen() value.Value {
	return constant.NewInt(types.I64, e.val)
}
