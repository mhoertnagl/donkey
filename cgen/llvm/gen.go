package llvm

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"github.com/mhoertnagl/donkey/cgen"
	"github.com/mhoertnagl/donkey/parser"
)

// TODO: Scopes for if
// TODO: for loop
// TODO: Assignment

var i1 = types.I1
var i64 = types.I64

var zeroI64 = constant.NewInt(i64, 0)
var minusOneI64 = constant.NewInt(i64, -1)
var wordSizeI64 = constant.NewInt(i64, 64)

var _false = constant.NewInt(i1, 0)
var _true = constant.NewInt(i1, 1)

type LlvmCodegen struct {
	ctx    *Context
	module *ir.Module
	fun    *ir.Func
	block  *ir.Block
}

func NewLlvmCodegen() cgen.Codegen {
	return &LlvmCodegen{
		ctx:    NewContext(),
		module: ir.NewModule(),
	}
}

func (c *LlvmCodegen) Generate(n *parser.Program) string {
	c.collectFunctionDefinitions(n)
	c.stmts(n.Statements)
	return c.module.String()
}

func (c *LlvmCodegen) stmts(ns []parser.Statement) value.Value {
	var res value.Value = nil
	for _, s := range ns {
		res = c.genStmt(s)
	}
	return res
}

func (c *LlvmCodegen) genStmt(n parser.Statement) value.Value {
	switch n := n.(type) {
	case *parser.LetStatement:
		return c.letStmt(n)
	case *parser.FunDefStatement:
		return c.funDefStmt(n)
	case *parser.BlockStatement:
		return c.blockStmt(n)
	case *parser.IfStatement:
		return c.ifStmt(n)
	case *parser.ReturnStatement:
		return c.returnStmt(n)
	case *parser.ExpressionStatement:
		return c.exprStmt(n)
	}
	return nil
}

func (c *LlvmCodegen) letStmt(n *parser.LetStatement) value.Value {
	name := n.Name.Value
	val := c.expr(n.Value)
	ptr := c.block.NewAlloca(i64)
	c.block.NewStore(val, ptr)
	c.ctx.SetValue(name, ptr)
	return ptr
}

func (c *LlvmCodegen) blockStmt(n *parser.BlockStatement) value.Value {
	return c.stmts(n.Statements)
}

func (c *LlvmCodegen) returnStmt(n *parser.ReturnStatement) value.Value {
	v := c.expr(n.Value)
	c.block.NewRet(v)
	return v
}

func (c *LlvmCodegen) exprStmt(n *parser.ExpressionStatement) value.Value {
	return c.expr(n.Value)
}

func (c *LlvmCodegen) setCurrentBlock(block *ir.Block) {
	c.block = block
}

func (c *LlvmCodegen) getCurrentBlock() *ir.Block {
	return c.block
}
