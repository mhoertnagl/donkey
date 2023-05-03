package llvm

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/value"
	"github.com/mhoertnagl/donkey/parser"
)

func (c *LlvmCodegen) funDefStmt(n *parser.FunDefStatement) value.Value {
	name := n.Name.Value
	// Load the appropriate function declaration.
	// Function declarations have been collected already in fun.decl.go.
	sym := c.ctx.Get(n.Name.Value)
	// TODO: map functions and identifiers separately
	switch fun := sym.GetValue().(type) {
	case *ir.Func:
		c.fun = fun
		// Create a new function scoped context. Arguments
		// and local variables are only visible in the
		// function body.
		c.ctx.PushScope()

		// Create the function entry block.
		// TODO: Wrapper for ir.Func that also holds function AST and provides convenient methods.
		//       c.fun.CreateEntryBlock()
		c.setCurrentBlock(c.fun.NewBlock(name + ".entry"))

		// Allocate and store function arguments.
		for _, arg := range c.fun.Params {
			// For each argument allocate space and store to
			// it the argument value.
			ptr := c.block.NewAlloca(i64)
			c.block.NewStore(arg, ptr)
			// Add the storage location of the argument to
			// the function scoped context.
			c.ctx.SetValue(arg.Name(), ptr)
		}

		// Compile the function body.
		c.blockStmt(n.Body)

		// Pop the function scoped context.
		c.ctx.PopScope()
	}
	return c.fun
}
