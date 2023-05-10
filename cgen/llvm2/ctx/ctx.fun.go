package ctx

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/value"
)

type FuncContext struct {
	module *ModuleContext
	fun    *ir.Func
	block  *ir.Block
	name   string
	locals Symbols
}

func (c *FuncContext) SetLocal(name string, value value.Value) {
	c.locals[name] = NewValueSymbol(value)
}

func (c *FuncContext) CreateEntryBlock() *ir.Block {
	c.block = c.CreateBlock(c.name + ".entry")
	return c.block
}

func (c *FuncContext) CreateBlock(name string) *ir.Block {
	return c.fun.NewBlock(name)
}

func (c *FuncContext) SetCurrentBlock(block *ir.Block) {
	c.block = block
}

func (c *FuncContext) GetCurrentBlock() *ir.Block {
	return c.block
}
