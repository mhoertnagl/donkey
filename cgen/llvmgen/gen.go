package llvmgen

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"github.com/mhoertnagl/donkey/parser"
	"github.com/mhoertnagl/donkey/token"
)

var zeroI64 = constant.NewInt(types.I64, 0)
var minusOneI64 = constant.NewInt(types.I64, -1)
var wordSizeI64 = constant.NewInt(types.I64, 32)

type llvmCodegen struct {
	module *ir.Module
	block  *ir.Block
	fun    *ir.Func
}

func (c *llvmCodegen) Generate(n parser.Program) {

}

func (c *llvmCodegen) generate(n parser.Node) value.Value {

}

func (c *llvmCodegen) generateLet(n parser.LetStatement) value.Value {

}

func (c *llvmCodegen) generateFunDef(n parser.FunctionDef) value.Value {
	fn := c.generateFunDecl(n)
}

func (c *llvmCodegen) generateFunDecl(n parser.FunctionDef) *ir.Func {
	name := n.Name.Value
	params := []*ir.Param{}
	for _, param := range n.Params {
		v := ir.NewParam(param.Value, types.I64)
		params = append(params, v)
	}
	return c.module.NewFunc(name, types.I64, params...)
}

func (c *llvmCodegen) generateBlock(n parser.BlockStatement) value.Value {

}

func (c *llvmCodegen) generateIf(n parser.IfStatement) value.Value {

}

func (c *llvmCodegen) generateReturn(n parser.ReturnStatement) value.Value {
	v := c.generate(n.Value)
	c.block.NewRet(v)
	// TODO: return the value that will be returned?
	return v
}

func (c *llvmCodegen) generateBool(n *parser.Boolean) value.Value {
	return constant.NewBool(n.Value)
}

func (c *llvmCodegen) generateInt(n *parser.Integer) value.Value {
	return constant.NewInt(types.I64, n.Value)
}

// func (c *llvmCodegen) generateLoadIdentifer(n *parser.Identifier) value.Value {
// 	return c.block.NewLoad(src) n.Value
// }

func (c *llvmCodegen) generateCall(n *parser.CallExpression) value.Value {
	name := c.generate(n.Function)
	args := []value.Value{}
	for _, arg := range n.Args {
		v := c.generate(arg)
		args = append(args, v)
	}
	return c.block.NewCall(name, args...)
}

func (c *llvmCodegen) generateBinary(n *parser.BinaryExpression) value.Value {
	l := c.generate(n.Left)
	r := c.generate(n.Right)
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
	one := constant.NewInt(types.I1, 1)
	v := c.generate(n)
	switch n.Operator {
	// int -> int
	case token.MINUS:
		return c.block.NewSub(zeroI64, v)
	// int -> int
	case token.INV:
		return c.block.NewXor(minusOneI64, v)
	// bool -> bool
	case token.NOT:
		return c.block.NewXor(one, v)
	}
	return nil
}
