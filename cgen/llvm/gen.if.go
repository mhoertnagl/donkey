package llvm

import (
	"github.com/llir/llvm/ir/value"
	"github.com/mhoertnagl/donkey/parser"
)

func (c *LlvmCodegen) ifStmt(n *parser.IfStatement) value.Value {
	if n.Alternative != nil {
		return c.ifWithAlt(n)
	}
	return c.ifWithoutAlt(n)
}

func (c *LlvmCodegen) ifWithAlt(n *parser.IfStatement) value.Value {
	then_block := c.fun.NewBlock("if.then")
	else_block := c.fun.NewBlock("if.else")

	// Generate the condition and then a conditional branch.
	cond := c.expr(n.Condition)
	c.block.NewCondBr(cond, then_block, else_block)

	// Set the current block to then_block then generate the
	// consequence statements.
	c.setCurrentBlock(then_block)
	// TODO: push scope
	c.genStmt(n.Consequence)
	// TODO: pop scope
	// Finally set the then_block to the current block. The
	// current block may not be the same as then_block because
	// stmts could have changed it because of a
	// nested if statement for instance.
	then_block = c.getCurrentBlock()

	// Set the current block to else_block then generate the
	// alternative statements.
	c.setCurrentBlock(else_block)
	// TODO: push scope
	c.genStmt(n.Alternative)
	// TODO: pop scope
	// Finally set the else_block to the current block. The
	// current block may not be the same as else_block because
	// stmts could have changed it because of a
	// nested if statement for instance.
	else_block = c.getCurrentBlock()

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
		merge_block := c.fun.NewBlock("if.merge")
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
		c.setCurrentBlock(merge_block)
	}

	return nil
}

func (c *LlvmCodegen) ifWithoutAlt(n *parser.IfStatement) value.Value {
	then_block := c.fun.NewBlock("if.then")
	merge_block := c.fun.NewBlock("if.merge")

	// Generate the condition and then a conditional branch.
	cond := c.expr(n.Condition)
	c.block.NewCondBr(cond, then_block, merge_block)

	// Set the current block to then_block then generate the
	// consequence statements.
	c.setCurrentBlock(then_block)
	// TODO: push scope
	c.genStmt(n.Consequence)
	// TODO: pop scope
	// Finally set the then_block to the current block. The
	// current block may not be the same as then_block because
	// stmts could have changed it because of a
	// nested if statement for instance.
	then_block = c.getCurrentBlock()
	// If no terminator has been set, complete the block with
	// an unconditional jump to the merge_block.
	if then_block.Term == nil {
		then_block.NewBr(merge_block)
	}

	// Continue with merge_block as the new current block.
	c.setCurrentBlock(merge_block)

	return nil
}
