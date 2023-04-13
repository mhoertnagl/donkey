package llvm

import "github.com/llir/llvm/ir/value"

type Scope map[string]value.Value

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

func (ctx *Context) Set(name string, value value.Value) {
	ctx.scopes[len(ctx.scopes)-1][name] = value
}

func (ctx *Context) Get(name string) value.Value {
	for i := len(ctx.scopes) - 1; i >= 0; i-- {
		if v, ok := ctx.scopes[i][name]; ok {
			return v
		}
	}
	return nil
}

// https://github.com/zegl/tre/blob/master/compiler/compiler/compiler.go#L84
