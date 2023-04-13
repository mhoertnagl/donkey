package llvm

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"github.com/mhoertnagl/donkey/cgen"
	"github.com/mhoertnagl/donkey/parser"
	"github.com/mhoertnagl/donkey/token"
	"github.com/mhoertnagl/donkey/utils"
)

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
	block  *ir.Block
	fun    *ir.Func
}

func NewLlvmCodegen() cgen.Codegen {
	ctx := NewContext()
	ctx.PushScope()
	return &LlvmCodegen{ctx: ctx}
}

func (c *LlvmCodegen) Generate(n parser.Program) {
	c.genStmts(n.Statements)
}

func (c *LlvmCodegen) genStmts(ns []parser.Statement) value.Value {
	// m := len(n.Statements)
	// for i := 0; i < m-1; i++ {
	// 	c.genStmt(n.Statements[i])
	// }
	// return c.genStmt(n.Statements[m-1])
	var res value.Value = nil
	for _, s := range ns {
		res = c.genStmt(s)
	}
	return res
}

func (c *LlvmCodegen) genStmt(n parser.Statement) value.Value {
	switch n := n.(type) {
	case *parser.LetStatement:
		return c.genLet(n)
	case *parser.FunDefStatement:
		return c.genFunDef(n)
	case *parser.BlockStatement:
		return c.genBlock(n)
	case *parser.ReturnStatement:
		return c.genReturn(n)
	case *parser.IfStatement:
		return c.genIf(n)
	case *parser.ExpressionStatement:
		return c.genExprStmt(n)
	}
	return nil
}

func (c *LlvmCodegen) genLet(n *parser.LetStatement) value.Value {
	val := c.genExpr(n.Value)
	loc := c.block.NewAlloca(i64)
	c.block.NewStore(val, loc)
	c.ctx.Set(n.Name.Value, loc)
	return loc
}

func (c *LlvmCodegen) genFunDef(n *parser.FunDefStatement) value.Value {
	c.fun = c.genFunDecl(n)
	// TODO: Push scope
	// TODO: add parameters
	c.genBlock(n.Body)
	// TODO: Pop scope
	return c.fun
}

func (c *LlvmCodegen) genFunDecl(n *parser.FunDefStatement) *ir.Func {
	name := n.Name.Value
	params := utils.Map(n.Params, c.genParam)
	// params := []*ir.Param{}
	// for _, param := range n.Params {
	// 	v := ir.NewParam(param.Value, i64)
	// 	params = append(params, v)
	// }
	return c.module.NewFunc(name, i64, params...)
}

func (c *LlvmCodegen) genParam(p *parser.Identifier) *ir.Param {
	return ir.NewParam(p.Value, i64)
}

func (c *LlvmCodegen) genBlock(n *parser.BlockStatement) value.Value {
	return c.genStmts(n.Statements)
}

func (c *LlvmCodegen) genIf(n *parser.IfStatement) value.Value {
	if n.Alternative != nil {
		return c.genIfWithAlt(n)
	}
	return c.genIfWithoutAlt(n)
}

func (c *LlvmCodegen) genIfWithAlt(n *parser.IfStatement) value.Value {
	then_block := c.fun.NewBlock("if.then")
	else_block := c.fun.NewBlock("if.else")
	merge_block := c.fun.NewBlock("if.merge")

	// Generate the condition and then a conditional branch.
	cond := c.genExpr(n.Condition)
	c.block.NewCondBr(cond, then_block, else_block)

	// Set the current block to then_block then generate the
	// consequence statements.
	c.setCurrentBlock(then_block)
	// TODO: push scope
	c.genStmt(n.Consequence)
	// TODO: pop scope
	// Finally set the then_block to the current block. The
	// current block may not be the same as then_block because
	// genStmts could have changed it because of a
	// nested if statement for instance.
	then_block = c.getCurrentBlock()
	// If no terminator has been set, complete the block with
	// an unconditional jump to the merge_block.
	if then_block.Term == nil {
		then_block.NewBr(merge_block)
	}

	// Set the current block to else_block then generate the
	// alternative statements.
	c.setCurrentBlock(else_block)
	// TODO: push scope
	c.genStmt(n.Alternative)
	// TODO: pop scope
	// Finally set the else_block to the current block. The
	// current block may not be the same as else_block because
	// genStmts could have changed it because of a
	// nested if statement for instance.
	else_block = c.getCurrentBlock()
	// If no terminator has been set, complete the block with
	// an unconditional jump to the merge_block.
	if else_block.Term == nil {
		else_block.NewBr(merge_block)
	}

	// Continue with merge_block as the new current block.
	c.setCurrentBlock(merge_block)

	return nil
}

func (c *LlvmCodegen) genIfWithoutAlt(n *parser.IfStatement) value.Value {
	then_block := c.fun.NewBlock("if.then")
	merge_block := c.fun.NewBlock("if.merge")

	// Generate the condition and then a conditional branch.
	cond := c.genExpr(n.Condition)
	c.block.NewCondBr(cond, then_block, merge_block)

	// Set the current block to then_block then generate the
	// consequence statements.
	c.setCurrentBlock(then_block)
	// TODO: push scope
	c.genStmt(n.Consequence)
	// TODO: pop scope
	// Finally set the then_block to the current block. The
	// current block may not be the same as then_block because
	// genStmts could have changed it because of a
	// nested if statement for instance.
	then_block = c.getCurrentBlock()
	// If no terminator has been set, complete the block with
	// an unconditional jump to the merge_block.
	if then_block.Term == nil {
		then_block.NewBr(merge_block)
	}

	// Continue with merge_block as the new current block.
	c.setCurrentBlock(merge_block)

	return nil
}

func (c *LlvmCodegen) genReturn(n *parser.ReturnStatement) value.Value {
	v := c.genExpr(n.Value)
	c.block.NewRet(v)
	return v
}

func (c *LlvmCodegen) genExprStmt(n *parser.ExpressionStatement) value.Value {
	return c.genExpr(n.Value)
}

func (c *LlvmCodegen) genExpr(n parser.Expression) value.Value {
	switch n := n.(type) {
	case *parser.Boolean:
		return c.genBoolean(n)
	case *parser.Integer:
		return c.genInteger(n)
	case *parser.Identifier:
		return c.genIdentifier(n)
	case *parser.CallExpression:
		return c.genCall(n)
	case *parser.BinaryExpression:
		return c.genBinary(n)
	case *parser.PrefixExpression:
		return c.genPrefix(n)
	}
	return nil
}

func (c *LlvmCodegen) genBoolean(n *parser.Boolean) value.Value {
	return constant.NewBool(n.Value)
}

func (c *LlvmCodegen) genInteger(n *parser.Integer) value.Value {
	return constant.NewInt(i64, n.Value)
}

func (c *LlvmCodegen) genIdentifier(n *parser.Identifier) value.Value {
	loc := c.ctx.Get(n.Value)
	// TODO: error if it not exists.
	return c.block.NewLoad(i64, loc)
}

func (c *LlvmCodegen) genCall(n *parser.CallExpression) value.Value {
	name := c.genExpr(n.Function)
	args := utils.Map(n.Args, c.genExpr)
	// args := []value.Value{}
	// for _, arg := range n.Args {
	// 	v := c.genExpr(arg)
	// 	args = append(args, v)
	// }
	return c.block.NewCall(name, args...)
}

func (c *LlvmCodegen) genBinary(n *parser.BinaryExpression) value.Value {
	l := c.genExpr(n.Left)
	r := c.genExpr(n.Right)
	switch n.Operator {
	// (int, int) -> int
	case token.PLUS:
		return c.block.NewAdd(l, r)
	// (int, int) -> int
	case token.MINUS:
		return c.block.NewSub(l, r)
	// (int, int) -> int
	case token.TIMES:
		return c.block.NewMul(l, r)
	// (int, int) -> int
	case token.DIV:
		return c.block.NewSDiv(l, r)

	// (int, int) -> int
	case token.AND:
		return c.block.NewAnd(l, r)
	// (int, int) -> int
	case token.OR:
		return c.block.NewOr(l, r)
	// (int, int) -> int
	case token.XOR:
		return c.block.NewXor(l, r)

	// TODO: Short circuiting?
	// (bool, bool) -> bool
	case token.CONJ:
		return c.block.NewAnd(l, r)
	// (bool, bool) -> bool
	case token.DISJ:
		return c.block.NewOr(l, r)

	// (int, int) -> int
	case token.SLL:
		return c.block.NewShl(l, r)
	// (int, int) -> int
	case token.SRL:
		return c.block.NewLShr(l, r)
	// (int, int) -> int
	case token.SRA:
		return c.block.NewAShr(l, r)
		// (int, int) -> int
	case token.ROL:
		x := c.block.NewShl(l, r)
		d := c.block.NewSub(wordSizeI64, r)
		y := c.block.NewLShr(l, d)
		return c.block.NewOr(x, y)
	// (int, int) -> int
	case token.ROR:
		x := c.block.NewLShr(l, r)
		d := c.block.NewSub(wordSizeI64, r)
		y := c.block.NewShl(l, d)
		return c.block.NewOr(x, y)

	// (int, int) -> boolean
	case token.EQU:
		return c.block.NewICmp(enum.IPredEQ, l, r)
		// (int, int) -> boolean
	case token.NEQ:
		return c.block.NewICmp(enum.IPredNE, l, r)
		// (int, int) -> boolean
	case token.LT:
		return c.block.NewICmp(enum.IPredSLT, l, r)
		// (int, int) -> boolean
	case token.LE:
		return c.block.NewICmp(enum.IPredSLE, l, r)
		// (int, int) -> boolean
	case token.GT:
		return c.block.NewICmp(enum.IPredSGT, l, r)
		// (int, int) -> boolean
	case token.GE:
		return c.block.NewICmp(enum.IPredSGE, l, r)
	}
	return nil
}

func (c *LlvmCodegen) genPrefix(n *parser.PrefixExpression) value.Value {
	v := c.genExpr(n.Value)
	switch n.Operator {
	// int -> int
	case token.MINUS:
		return c.block.NewSub(zeroI64, v)
	// int -> int
	case token.INV:
		return c.block.NewXor(minusOneI64, v)
	// bool -> bool
	case token.NOT:
		return c.block.NewXor(_true, v)
	}
	return nil
}

func (c *LlvmCodegen) setCurrentBlock(block *ir.Block) {
	c.block = block
}

func (c *LlvmCodegen) getCurrentBlock() *ir.Block {
	return c.block
}
