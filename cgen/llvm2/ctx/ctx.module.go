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

func NewModuleContext(module *ir.Module) *ModuleContext {
	return &ModuleContext{
		module:    module,
		functions: make(Functions),
	}
}

func (c *ModuleContext) NewFuncContext(name string, retType types.Type, params ...*ir.Param) *FuncContext {
	c.functions[name] = &FuncContext{
		module: c,
		fun:    c.module.NewFunc(name, retType, params...),
		name:   name,
		locals: make(Symbols),
	}
	return c.functions[name]
}

// func (c *ModuleContext) GetFunction
