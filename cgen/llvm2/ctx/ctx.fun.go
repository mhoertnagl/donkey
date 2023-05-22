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
	scopes Scopes
}

func NewFuncContext(module *ModuleContext, fun *ir.Func, name string) *FuncContext {
	ctx := &FuncContext{
		module: module,
		fun:    fun,
		name:   name,
		scopes: make(Scopes, 0),
	}
	ctx.PushScope()
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

func (c *FuncContext) PushScope() {
	c.scopes = append(c.scopes, make(SymbolTable))
}

func (c *FuncContext) PopScope() {
	c.scopes = c.scopes[:len(c.scopes)-1]
}

func (c *FuncContext) Set(name string, value value.Value) {
	c.scopes[len(c.scopes)-1][name] = NewValueSymbol(value)
}

func (c *FuncContext) Get(name string) (Symbol, bool) {
	for i := len(c.scopes) - 1; i >= 0; i-- {
		if sym, ok := c.scopes[i][name]; ok {
			return sym, true
		}
	}
	if fun, ok := c.module.Get(name); ok {
		return fun, ok
	}
	return nil, false
}
