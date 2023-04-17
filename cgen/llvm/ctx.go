package llvm

import "github.com/llir/llvm/ir/value"

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
	value value.Value
}

func (sym *FuncSymbol) GetValue() value.Value {
	return sym.value
}

type Scope map[string]Symbol

type Context struct {
	// parent  *Context
	scopes []Scope
}

func NewContext() *Context {
	return &Context{
		// parent:  nil,
		scopes: make([]Scope, 0),
	}
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

func (ctx *Context) SetFunction(name string, value value.Value) {
	ctx.setSymbol(name, &FuncSymbol{value})
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
