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

var zeroI32 = constant.NewInt(types.I32, 0)
var minusOneI32 = constant.NewInt(types.I32, -1)
var wordSizeI32 = constant.NewInt(types.I32, 32)

type llvmCodegen struct {
  module *ir.Module
  block  *ir.Block
  fun    *ir.Func
}

func (c *llvmCodegen) Generate(node parser.Program) {

}

func (c *llvmCodegen) generate(node parser.Node) value.Value {

}

func (c *llvmCodegen) generateLet(node parser.LetStatement) value.Value {

}

func (c *llvmCodegen) generateBlock(node parser.BlockStatement) value.Value {

}

func (c *llvmCodegen) generateIf(node parser.IfStatement) value.Value {

}

func (c *llvmCodegen) generateReturn(node parser.ReturnStatement) value.Value {
  v := c.generate(node.Value)
  c.block.NewRet(v)
  // TODO: return the value that will be returned?
  return v
}

func (c *llvmCodegen) generateBool(n *parser.Boolean) value.Value {
	return constant.NewBool(n.Value)
}

func (c *llvmCodegen) generateInt(n *parser.Integer) value.Value {
	return constant.NewInt(types.I32, n.Value)
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
		d := c.block.NewSub(wordSizeI32, r)
		y := c.block.NewLShr(l, d)
		return c.block.NewOr(x, y)
	// (int, int) -> int
	case token.ROR:
		x := c.block.NewLShr(l, r)
		d := c.block.NewSub(wordSizeI32, r)
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
		return c.block.NewSub(zeroI32, v)
	// int -> int
	case token.INV:
		return c.block.NewXor(minusOneI32, v)
	// bool -> bool
	case token.NOT:
		return c.block.NewXor(one, v)
	}
	return nil
}
