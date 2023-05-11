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
	locals Scopes
}

func NewFuncContext(module *ModuleContext, fun *ir.Func, name string) *FuncContext {
	ctx := &FuncContext{
		module: module,
		fun:    fun,
		name:   name,
		locals: make(Scopes, 0),
	}
	// Initial function scope.
	ctx.PushScope()
	return ctx
}

func (c *FuncContext) PushScope() {
	c.locals = append(c.locals, make(Symbols))
}

func (c *FuncContext) PopScope() {
	c.locals = c.locals[:len(c.locals)-1]
}

func (c *FuncContext) SetLocal(name string, value value.Value) {
	c.locals[len(c.locals)-1][name] = NewValueSymbol(value)
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
