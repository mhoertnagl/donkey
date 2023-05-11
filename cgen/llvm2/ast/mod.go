package ast

import "github.com/mhoertnagl/donkey/cgen/llvm2/ctx"

type Module struct {
	module *ctx.ModuleContext
	stmts  Stmts
}

func NewModule(module *ctx.ModuleContext, stmts Stmts) *Module {
	return &Module{module, stmts}
}

func (m *Module) Gen() string {
	for _, stmt := range m.stmts {
		stmt.gen()
	}
	return m.module.String()
}
