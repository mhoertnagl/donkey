package ctx

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
)

type Functions map[string]*FuncContext

type ModuleContext struct {
	module    *ir.Module
	functions Functions
}

func NewModuleContext() *ModuleContext {
	return &ModuleContext{
		module:    ir.NewModule(),
		functions: make(Functions),
	}
}

func (c *ModuleContext) NewFuncContext(name string, retType types.Type, params ...*ir.Param) *FuncContext {
	fun := c.module.NewFunc(name, retType, params...)
	c.functions[name] = NewFuncContext(c, fun, name)
	return c.functions[name]
}

// func (c *ModuleContext) GetFunction

func (c *ModuleContext) String() string {
	return c.module.String()
}
