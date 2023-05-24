package pass

import (
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/mhoertnagl/donkey/cgen/llvm/ast"
	"github.com/mhoertnagl/donkey/cgen/llvm/ctx"
	"github.com/mhoertnagl/donkey/parser"
	"github.com/mhoertnagl/donkey/utils"
)

type AstPass struct {
	module *ctx.ModuleContext
	fun    *ctx.FuncContext
}

func NewAstPass() *AstPass {
	return &AstPass{
		module: ctx.NewModuleContext(),
	}
}

func (p *AstPass) Run(n *parser.Program) *ast.Module {
	return p.program(n)
}

func (p *AstPass) program(n *parser.Program) *ast.Module {
	stmts := p.stmts(n.Statements)
	return ast.NewModule(p.module, stmts)
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
		return p.ifStmt(n)
	case *parser.ReturnStatement:
		return p.returnStmt(n)
	case *parser.ExpressionStatement:
		return p.exprStmt(n)
	case nil:
		// n can be nil if it is the missing alternative
		// of an if statement.
		return nil
	}
	panic(fmt.Sprintf("Unsupported statement type: %v", n))
}

func (p *AstPass) letStmt(n *parser.LetStatement) *ast.LetStmt {
	name := n.Name.Literal()
	expr := p.expr(n.Value)
	return ast.NewLetStmt(p.fun, name, expr)
}

func (p *AstPass) funDefStmt(n *parser.FunDefStatement) *ast.FunDefStmt {
	name := n.Name.Value
	irParams := utils.Map(n.Params, irParam)
	p.fun = p.module.NewFuncContext(name, types.I64, irParams...)
	params := p.mapParams(n.Params, irParams)
	body := p.blockStmt(n.Body)
	return ast.NewFunDefStmt(p.fun, params, body)
}

func irParam(p *parser.Identifier) *ir.Param {
	return ir.NewParam(p.Value, types.I64)
}

// TODO: Can this be done better?
func (p *AstPass) mapParams(is []*parser.Identifier, ps []*ir.Param) []*ast.ParamExpr {
	res := make([]*ast.ParamExpr, len(is))
	for i := 0; i < len(is); i++ {
		res[i] = ast.NewParamExpr(p.fun, is[i].Value, ps[i])
	}
	return res
}

func (p *AstPass) blockStmt(n *parser.BlockStatement) ast.Stmts {
	return p.stmts(n.Statements)
}

func (p *AstPass) ifStmt(n *parser.IfStatement) *ast.IfStmt {
	cond := p.expr(n.Condition)
	cons := p.stmt(n.Consequence)
	alt := p.stmt(n.Alternative)
	return ast.NewIfStmt(p.fun, cond, cons, alt)
}

func (p *AstPass) returnStmt(n *parser.ReturnStatement) *ast.ReturnStmt {
	expr := p.expr(n.Value)
	return ast.NewReturnStmt(p.fun, expr)
}

func (p *AstPass) exprStmt(n *parser.ExpressionStatement) *ast.ExprStmt {
	expr := p.expr(n.Value)
	return ast.NewExprStmt(p.fun, expr)
}

func (p *AstPass) expr(n parser.Expression) ast.Expr {
	switch n := n.(type) {
	case *parser.Boolean:
		return p.boolLit(n)
	case *parser.Integer:
		return p.intLit(n)
	case *parser.Identifier:
		return p.identifier(n)
	case *parser.CallExpression:
		return p.callExpr(n)
	case *parser.BinaryExpression:
		return p.binaryExpr(n)
	case *parser.PrefixExpression:
		return p.prefixExpr(n)
	case *parser.AssignmentExpression:
		return p.assignExpr(n)
	}
	panic("Unsupported expression type.")
}

func (p *AstPass) boolLit(n *parser.Boolean) ast.Expr {
	return ast.NewBoolLiteralExpr(p.fun, n.Value)
}

func (p *AstPass) intLit(n *parser.Integer) ast.Expr {
	return ast.NewIntLiteralExpr(p.fun, n.Value)
}

func (p *AstPass) identifier(n *parser.Identifier) *ast.IdentifierExpr {
	return ast.NewIdentifierExpr(p.fun, n.Value)
}

func (p *AstPass) callExpr(n *parser.CallExpression) ast.Expr {
	callee := p.expr(n.Function)
	args := utils.Map(n.Args, p.expr)
	return ast.NewFunCallExpr(p.fun, callee, args)
}

func (p *AstPass) binaryExpr(n *parser.BinaryExpression) ast.Expr {
	left := p.expr(n.Left)
	right := p.expr(n.Right)
	return ast.NewBinaryExpr(p.fun, left, n.Operator, right)
}

func (p *AstPass) prefixExpr(n *parser.PrefixExpression) ast.Expr {
	val := p.expr(n.Value)
	return ast.NewPrefixExpr(p.fun, n.Operator, val)
}

func (p *AstPass) assignExpr(n *parser.AssignmentExpression) ast.Expr {
	id := p.identifier(n.Name)
	value := p.expr(n.Value)
	return ast.NewAssignExpr(p.fun, id, value)
}
