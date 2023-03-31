package llvmgen

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"github.com/mhoertnagl/donkey/parser"
	"github.com/mhoertnagl/donkey/token"
	"github.com/mhoertnagl/donkey/utils"
)

var zeroI64 = constant.NewInt(types.I64, 0)
var minusOneI64 = constant.NewInt(types.I64, -1)
var wordSizeI64 = constant.NewInt(types.I64, 32)

var _false = constant.NewInt(types.I1, 0)
var _true = constant.NewInt(types.I1, 1)

type llvmCodegen struct {
	module *ir.Module
	block  *ir.Block
	fun    *ir.Func
}

func (c *llvmCodegen) Generate(n parser.Program) {
	c.generateStatements(n.Statements)
}

func (c *llvmCodegen) generateStatements(ns []parser.Statement) value.Value {
	// m := len(n.Statements)
	// for i := 0; i < m-1; i++ {
	// 	c.generateStatement(n.Statements[i])
	// }
	// return c.generateStatement(n.Statements[m-1])
	var res value.Value = nil
	for _, s := range ns {
		res = c.generateStatement(s)
	}
	return res
}

func (c *llvmCodegen) generateStatement(n parser.Statement) value.Value {
	switch n := n.(type) {
	case *parser.LetStatement:
		return c.generateLet(n)
	case *parser.FunDefStatement:
		return c.generateFunDef(n)
	case *parser.BlockStatement:
		return c.generateBlock(n)
	case *parser.ReturnStatement:
		return c.generateReturn(n)
	case *parser.IfStatement:
		return c.generateIf(n)
	case *parser.ExpressionStatement:
		return c.generateExprStmt(n)
	}
	return nil
}

func (c *llvmCodegen) generateLet(n *parser.LetStatement) value.Value {
	val := c.generateExpression(n.Value)
	loc := c.block.NewAlloca(types.I64)
	c.block.NewStore(val, loc)
	return loc
}

func (c *llvmCodegen) generateFunDef(n *parser.FunDefStatement) value.Value {
	c.fun = c.generateFunDecl(n)
	c.generateBlock(n.Body)
	return c.fun
}

func (c *llvmCodegen) generateFunDecl(n *parser.FunDefStatement) *ir.Func {
	name := n.Name.Value
	params := utils.Map(n.Params, c.generateParam)
	// params := []*ir.Param{}
	// for _, param := range n.Params {
	// 	v := ir.NewParam(param.Value, types.I64)
	// 	params = append(params, v)
	// }
	return c.module.NewFunc(name, types.I64, params...)
}

func (c *llvmCodegen) generateParam(p *parser.Identifier) *ir.Param {
	return ir.NewParam(p.Value, types.I64)
}

func (c *llvmCodegen) generateBlock(n *parser.BlockStatement) value.Value {
	return c.generateStatements(n.Statements)
}

func (c *llvmCodegen) generateIf(n *parser.IfStatement) value.Value {
	then_block := c.fun.NewBlock("if.then")
	else_block := c.fun.NewBlock("if.else")
	merge_block := c.fun.NewBlock("if.merge")
	cond := c.generateExpression(n.Condition)
	c.block.NewCondBr(cond, then_block, else_block)
	// TODO: with then_block
	c.generateStatement(n.Consequence)
	then_block.NewBr(merge_block)
	if n.Alternative != nil {
		// TODO: with else_block
		c.generateStatement(n.Alternative)
		else_block.NewBr(merge_block)
	}
	// TODO: continue with merge_block
	return nil
}

func (c *llvmCodegen) generateReturn(n *parser.ReturnStatement) value.Value {
	v := c.generateExpression(n.Value)
	c.block.NewRet(v)
	return v
}

func (c *llvmCodegen) generateExprStmt(n *parser.ExpressionStatement) value.Value {
	return c.generateExpression(n.Value)
}

func (c *llvmCodegen) generateExpression(n parser.Expression) value.Value {
	switch n := n.(type) {
	case *parser.Boolean:
		return c.generateBoolean(n)
	case *parser.Integer:
		return c.generateInteger(n)
	case *parser.CallExpression:
		return c.generateCall(n)
	case *parser.BinaryExpression:
		return c.generateBinary(n)
	case *parser.PrefixExpression:
		return c.generatePrefix(n)
	}
	return nil
}

func (c *llvmCodegen) generateBoolean(n *parser.Boolean) value.Value {
	return constant.NewBool(n.Value)
}

func (c *llvmCodegen) generateInteger(n *parser.Integer) value.Value {
	return constant.NewInt(types.I64, n.Value)
}

// func (c *llvmCodegen) generateIdVal(n *parser.Identifier) value.Value {
// 	return c.block.NewLoad(src)
// }

func (c *llvmCodegen) generateCall(n *parser.CallExpression) value.Value {
	name := c.generateExpression(n.Function)
	args := utils.Map(n.Args, c.generateExpression)
	// args := []value.Value{}
	// for _, arg := range n.Args {
	// 	v := c.generateExpression(arg)
	// 	args = append(args, v)
	// }
	return c.block.NewCall(name, args...)
}

func (c *llvmCodegen) generateBinary(n *parser.BinaryExpression) value.Value {
	l := c.generateExpression(n.Left)
	r := c.generateExpression(n.Right)
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

func (c *llvmCodegen) generatePrefix(n *parser.PrefixExpression) value.Value {
	v := c.generateExpression(n.Value)
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
