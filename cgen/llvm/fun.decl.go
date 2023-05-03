package llvm

import (
	"github.com/llir/llvm/ir"
	"github.com/mhoertnagl/donkey/parser"
	"github.com/mhoertnagl/donkey/utils"
)

func (c *LlvmCodegen) collectFunctionDefinitions(n *parser.Program) {
	stmts(c, n.Statements)
}

func stmts(c *LlvmCodegen, ns []parser.Statement) {
	for _, s := range ns {
		stmt(c, s)
	}
}

func stmt(c *LlvmCodegen, n parser.Statement) {
	switch n := n.(type) {
	case *parser.FunDefStatement:
		funDefStmt(c, n)
	}
}

func funDefStmt(c *LlvmCodegen, n *parser.FunDefStatement) {
	name := n.Name.Value
	params := utils.Map(n.Params, param)
	fun := c.module.NewFunc(name, i64, params...)
	c.ctx.SetFunction(n.Name.Value, fun)
}

func param(p *parser.Identifier) *ir.Param {
	return ir.NewParam(p.Value, i64)
}
