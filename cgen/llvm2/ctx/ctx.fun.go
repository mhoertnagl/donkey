package ctx

import (
	"github.com/llir/llvm/ir"
)

type FuncContext struct {
	SymbolContainer
	module *ModuleContext
	fun    *ir.Func
	block  *ir.Block
	name   string
}

func NewFuncContext(module *ModuleContext, fun *ir.Func, name string) *FuncContext {
	ctx := &FuncContext{
		module: module,
		fun:    fun,
		name:   name,
	}
	ctx.InitSymbolContainer(nil)
	return ctx
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
