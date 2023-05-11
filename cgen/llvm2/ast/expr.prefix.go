package ast

import (
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"github.com/mhoertnagl/donkey/cgen/llvm2/ctx"
	"github.com/mhoertnagl/donkey/token"
)

var zeroI64 = constant.NewInt(types.I64, 0)
var minusOneI64 = constant.NewInt(types.I64, -1)
var _true = constant.NewInt(types.I1, 1)

type PrefixExpr struct {
	fun *ctx.FuncContext
	op  token.TokenType
	val Expr
}

func NewPrefixExpr(fun *ctx.FuncContext, op token.TokenType, val Expr) *PrefixExpr {
	return &PrefixExpr{fun, op, val}
}

func (e *PrefixExpr) gen() value.Value {
	blk := e.fun.GetCurrentBlock()
	v := e.val.gen()
	switch e.op {
	case token.MINUS:
		return blk.NewSub(zeroI64, v)
	case token.INV:
		return blk.NewXor(minusOneI64, v)
	case token.NOT:
		return blk.NewXor(_true, v)
	}
	return nil
}
