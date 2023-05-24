package ast

import (
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"github.com/mhoertnagl/donkey/cgen/llvm/ctx"
	"github.com/mhoertnagl/donkey/token"
)

var wordSizeI64 = constant.NewInt(types.I64, 64)

type BinaryExpr struct {
	fun   *ctx.FuncContext
	left  Expr
	op    token.TokenType
	right Expr
}

func NewBinaryExpr(fun *ctx.FuncContext, left Expr, op token.TokenType, right Expr) *BinaryExpr {
	return &BinaryExpr{fun, left, op, right}
}

func (e *BinaryExpr) gen() value.Value {
	blk := e.fun.GetCurrentBlock()
	l := e.left.gen()
	r := e.right.gen()
	switch e.op {
	case token.PLUS:
		return blk.NewAdd(l, r)
	case token.MINUS:
		return blk.NewSub(l, r)
	case token.TIMES:
		return blk.NewMul(l, r)
	case token.DIV:
		return blk.NewSDiv(l, r)

	case token.AND:
		return blk.NewAnd(l, r)
	case token.OR:
		return blk.NewOr(l, r)
	case token.XOR:
		return blk.NewXor(l, r)

	// TODO: Short circuiting?
	case token.CONJ:
		return blk.NewAnd(l, r)
	case token.DISJ:
		return blk.NewOr(l, r)

	case token.SLL:
		return blk.NewShl(l, r)
	case token.SRL:
		return blk.NewLShr(l, r)
	case token.SRA:
		return blk.NewAShr(l, r)
	case token.ROL:
		x := blk.NewShl(l, r)
		d := blk.NewSub(wordSizeI64, r)
		y := blk.NewLShr(l, d)
		return blk.NewOr(x, y)
	case token.ROR:
		x := blk.NewLShr(l, r)
		d := blk.NewSub(wordSizeI64, r)
		y := blk.NewShl(l, d)
		return blk.NewOr(x, y)

	case token.EQU:
		return blk.NewICmp(enum.IPredEQ, l, r)
	case token.NEQ:
		return blk.NewICmp(enum.IPredNE, l, r)
	case token.LT:
		return blk.NewICmp(enum.IPredSLT, l, r)
	case token.LE:
		return blk.NewICmp(enum.IPredSLE, l, r)
	case token.GT:
		return blk.NewICmp(enum.IPredSGT, l, r)
	case token.GE:
		return blk.NewICmp(enum.IPredSGE, l, r)
	}
	return nil
}
