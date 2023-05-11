package ast

import (
	"github.com/llir/llvm/ir/types"
	"github.com/mhoertnagl/donkey/cgen/llvm2/ctx"
)

// TODO: BlockContext
type LetStmt struct {
	fun  *ctx.FuncContext
	name string
	val  Expr
}

func NewLetStmt(fun *ctx.FuncContext, name string, val Expr) *LetStmt {
	return &LetStmt{fun, name, val}
}

// TODO: Allocate appropriate type.
func (n *LetStmt) gen() {
	blk := n.fun.GetCurrentBlock()
	val := n.val.gen()
	ptr := blk.NewAlloca(types.I64)
	blk.NewStore(val, ptr)
	n.fun.Set(n.name, ptr)
}
