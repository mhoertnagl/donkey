package ast

import (
	"github.com/mhoertnagl/donkey/cgen/llvm/ctx"
)

type IfStmt struct {
	fun  *ctx.FuncContext
	cond Expr
	cons Stmt
	alt  Stmt
}

func NewIfStmt(fun *ctx.FuncContext, cond Expr, cons Stmt, alt Stmt) *IfStmt {
	return &IfStmt{fun, cond, cons, alt}
}

func (n *IfStmt) gen() {
	if n.alt != nil {
		n.genIfElse()
	} else {
		n.genIf()
	}
}

func (n *IfStmt) genIfElse() {
	then_block := n.fun.CreateBlock("if.then")
	else_block := n.fun.CreateBlock("if.else")

	// Generate the condition and then a conditional branch.
	cond := n.cond.gen()
	blk := n.fun.GetCurrentBlock()
	blk.NewCondBr(cond, then_block, else_block)

	// Set the current block to then_block then generate the
	// consequence statements.
	n.fun.SetCurrentBlock(then_block)
	n.fun.PushScope()
	n.cons.gen()
	n.fun.PopScope()
	// Finally set the then_block to the current block. The
	// current block may not be the same as then_block because
	// the consequence statements could could have changed it.
	then_block = n.fun.GetCurrentBlock()

	// Set the current block to else_block then generate the
	// alternative statements.
	n.fun.SetCurrentBlock(else_block)
	n.fun.PushScope()
	n.alt.gen()
	n.fun.PopScope()
	// Finally set the else_block to the current block. The
	// current block may not be the same as else_block because
	// the alternative statements could have changed it.
	else_block = n.fun.GetCurrentBlock()

	// If either block is missing a terminator, create a merge
	// block. Consider the following program:
	//
	//   fn main() {
	//     let a = 1;
	//     let b = 2;
	//     if b < a {
	//       return a;
	//     } else {
	//       return b;
	//     }
	//   }
	//
	// If both blocks already terminate the merge block is
	// superfluous and ends without a terminator which results
	// in a compilation error.
	if then_block.Term == nil || else_block.Term == nil {
		merge_block := n.fun.CreateBlock("if.merge")
		// If the then_block has no terminator, complete the block
		// with an unconditional jump to the merge_block.
		if then_block.Term == nil {
			then_block.NewBr(merge_block)
		}
		// If the else_block has no terminator, complete the block
		// with an unconditional jump to the merge_block.
		if else_block.Term == nil {
			else_block.NewBr(merge_block)
		}

		// Continue with merge_block as the new current block.
		n.fun.SetCurrentBlock(merge_block)
	}
}

func (n *IfStmt) genIf() {
	then_block := n.fun.CreateBlock("if.then")
	merge_block := n.fun.CreateBlock("if.merge")

	// Generate the condition and then a conditional branch.
	cond := n.cond.gen()
	blk := n.fun.GetCurrentBlock()
	blk.NewCondBr(cond, then_block, merge_block)

	// Set the current block to then_block then generate the
	// consequence statements.
	n.fun.SetCurrentBlock(then_block)
	n.fun.PushScope()
	n.cons.gen()
	n.fun.PopScope()
	// Finally set the then_block to the current block. The
	// current block may not be the same as then_block because
	// the consequence statements could have changed.

	then_block = n.fun.GetCurrentBlock()
	// If no terminator has been set, complete the block with
	// an unconditional jump to the merge_block.
	if then_block.Term == nil {
		then_block.NewBr(merge_block)
	}

	// Continue with merge_block as the new current block.
	n.fun.SetCurrentBlock(merge_block)
}
