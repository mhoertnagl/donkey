package ast

import (
	"github.com/llir/llvm/ir/types"
	"github.com/mhoertnagl/donkey/cgen/llvm2/ctx"
)

type FunDefStmt struct {
	fun  *ctx.FuncContext
	body Stmts
}

func NewFunDefStmt(fun *ctx.FuncContext, body Stmts) *FunDefStmt {
	return &FunDefStmt{fun, body}
}

func (n *FunDefStmt) gen() {
	n.fun.CreateEntryBlock()
	n.allocArgs()
	n.body.gen()
}

func (n *FunDefStmt) allocArgs() {
	blk := n.fun.GetCurrentBlock()
	for _, arg := range n.fun.Params {
		ptr := blk.NewAlloca(types.I64)
		blk.NewStore(arg, ptr)
		n.fun.Set(arg.Name(), ptr)
	}
}
