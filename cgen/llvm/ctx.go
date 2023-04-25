package llvm

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/value"
)

type Symbol interface {
	GetValue() value.Value
}

type ValueSymbol struct {
	value value.Value
}

func (sym *ValueSymbol) GetValue() value.Value {
	return sym.value
}

type FuncSymbol struct {
	fun *ir.Func
}

func (sym *FuncSymbol) GetValue() value.Value {
	return sym.fun
}

type Scope map[string]Symbol

type Context struct {
	scopes []Scope
}

func NewContext() *Context {
	ctx := &Context{
		scopes: make([]Scope, 0),
	}
	ctx.PushScope()
	return ctx
}

func (ctx *Context) PushScope() {
	ctx.scopes = append(ctx.scopes, Scope{})
}

func (ctx *Context) PopScope() {
	ctx.scopes = ctx.scopes[:len(ctx.scopes)-1]
}

func (ctx *Context) SetValue(name string, value value.Value) {
	ctx.setSymbol(name, &ValueSymbol{value})
}

func (ctx *Context) SetFunction(name string, fun *ir.Func) {
	ctx.setSymbol(name, &FuncSymbol{fun})
}

func (ctx *Context) setSymbol(name string, symbol Symbol) {
	ctx.scopes[len(ctx.scopes)-1][name] = symbol
}

func (ctx *Context) Get(name string) Symbol {
	for i := len(ctx.scopes) - 1; i >= 0; i-- {
		if v, ok := ctx.scopes[i][name]; ok {
			return v
		}
	}
	return nil
}

// https://github.com/zegl/tre/blob/master/compiler/compiler/compiler.go#L84
