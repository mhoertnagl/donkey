package pass

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/mhoertnagl/donkey/cgen/llvm2/ast"
	"github.com/mhoertnagl/donkey/cgen/llvm2/ctx"
	"github.com/mhoertnagl/donkey/parser"
	"github.com/mhoertnagl/donkey/utils"
)

type AstPass struct {
	module *ctx.ModuleContext
	fun    *ctx.FuncContext
}

func NewAstPass() *AstPass {
	return &AstPass{
		module: ctx.NewModuleContext(ir.NewModule()),
	}
}

func (p *AstPass) Run(n parser.Program) *ast.Module {
	return p.program(n)
}

func (p *AstPass) program(n parser.Program) *ast.Module {
	p.stmts(n.Statements)
	return nil
}

func (p *AstPass) stmts(ns []parser.Statement) ast.Stmts {
	return utils.Map(ns, p.stmt)
}

func (p *AstPass) stmt(n parser.Statement) ast.Stmt {
	switch n := n.(type) {
	case *parser.LetStatement:
		return p.letStmt(n)
	case *parser.FunDefStatement:
		return p.funDefStmt(n)
	case *parser.BlockStatement:
		return p.blockStmt(n)
	case *parser.IfStatement:
		p.ifStmt(n)
	case *parser.ReturnStatement:
		return p.returnStmt(n)
	case *parser.ExpressionStatement:
		p.exprStmt(n)
	}
}

func (p *AstPass) letStmt(n *parser.LetStatement) *ast.LetStmt {
	return ast.NewLetStmt(p.fun, n.Name.Literal(), p.expr(n.Value))
}

func (p *AstPass) funDefStmt(n *parser.FunDefStatement) *ast.FunDefStmt {
	name := n.Name.Value
	params := utils.Map(n.Params, param)
	p.fun = p.module.NewFuncContext(name, types.I64, params...)
	body := p.blockStmt(n.Body)
	return ast.NewFunDefStmt(p.fun, body)
}

func param(p *parser.Identifier) *ir.Param {
	return ir.NewParam(p.Value, types.I64)
}

func (p *AstPass) blockStmt(n *parser.BlockStatement) ast.Stmts {
	return p.stmts(n.Statements)
}

func (p *AstPass) ifStmt(n *parser.IfStatement) {

}

func (p *AstPass) returnStmt(n *parser.ReturnStatement) *ast.ReturnStmt {
	return ast.NewReturnStmt(p.fun, p.expr(n.Value))
}

func (p *AstPass) exprStmt(n *parser.ExpressionStatement) {

}

func (p *AstPass) expr(n parser.Expression) ast.Expr {

}
