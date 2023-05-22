package ast

import (
	"github.com/llir/llvm/ir/types"
	"github.com/mhoertnagl/donkey/cgen/llvm2/ctx"
)

type FunDefStmt struct {
	fun    *ctx.FuncContext
	params []*ParamExpr
	body   Stmts
}

func NewFunDefStmt(fun *ctx.FuncContext, params []*ParamExpr, body Stmts) *FunDefStmt {
	return &FunDefStmt{fun, params, body}
}

func (n *FunDefStmt) gen() {
	n.fun.CreateEntryBlock()
	n.allocArgs()
	n.body.gen()
}

func (n *FunDefStmt) allocArgs() {
	blk := n.fun.GetCurrentBlock()
	for _, arg := range n.params {
		ptr := blk.NewAlloca(types.I64)
		val := arg.gen()
		blk.NewStore(val, ptr)
		n.fun.Set(arg.name, ptr)
	}
}
