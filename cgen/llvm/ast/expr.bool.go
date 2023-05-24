package ast

import (
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/value"
	"github.com/mhoertnagl/donkey/cgen/llvm/ctx"
)

type BoolLiteralExpr struct {
	fun *ctx.FuncContext
	val bool
}

func NewBoolLiteralExpr(fun *ctx.FuncContext, val bool) *BoolLiteralExpr {
	return &BoolLiteralExpr{fun, val}
}

func (e *BoolLiteralExpr) gen() value.Value {
	return constant.NewBool(e.val)
}
