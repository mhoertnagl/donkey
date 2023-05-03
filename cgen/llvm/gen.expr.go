package llvm

import (
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/value"
	"github.com/mhoertnagl/donkey/parser"
	"github.com/mhoertnagl/donkey/token"
	"github.com/mhoertnagl/donkey/utils"
)

func (c *LlvmCodegen) expr(n parser.Expression) value.Value {
	switch n := n.(type) {
	case *parser.Boolean:
		return c.boolLit(n)
	case *parser.Integer:
		return c.intLit(n)
	case *parser.Identifier:
		return c.identifier(n)
	case *parser.CallExpression:
		return c.callExpr(n)
	case *parser.BinaryExpression:
		return c.binaryExpr(n)
	case *parser.PrefixExpression:
		return c.prefixExpr(n)
	}
	return nil
}

func (c *LlvmCodegen) boolLit(n *parser.Boolean) value.Value {
	return constant.NewBool(n.Value)
}

func (c *LlvmCodegen) intLit(n *parser.Integer) value.Value {
	return constant.NewInt(i64, n.Value)
}

func (c *LlvmCodegen) identifier(n *parser.Identifier) value.Value {
	sym := c.ctx.Get((n.Value))
	switch sym := sym.(type) {
	case *ValueSymbol:
		return c.block.NewLoad(i64, sym.GetValue())
	case *FuncSymbol:
		return sym.GetValue()
	}
	return nil
}

func (c *LlvmCodegen) callExpr(n *parser.CallExpression) value.Value {
	name := c.expr(n.Function)
	args := utils.Map(n.Args, c.expr)
	return c.block.NewCall(name, args...)
}

func (c *LlvmCodegen) binaryExpr(n *parser.BinaryExpression) value.Value {
	l := c.expr(n.Left)
	r := c.expr(n.Right)
	switch n.Operator {
	case token.PLUS:
		return c.block.NewAdd(l, r)
	case token.MINUS:
		return c.block.NewSub(l, r)
	case token.TIMES:
		return c.block.NewMul(l, r)
	case token.DIV:
		return c.block.NewSDiv(l, r)

	case token.AND:
		return c.block.NewAnd(l, r)
	case token.OR:
		return c.block.NewOr(l, r)
	case token.XOR:
		return c.block.NewXor(l, r)

	// TODO: Short circuiting?
	case token.CONJ:
		return c.block.NewAnd(l, r)
	case token.DISJ:
		return c.block.NewOr(l, r)

	case token.SLL:
		return c.block.NewShl(l, r)
	case token.SRL:
		return c.block.NewLShr(l, r)
	case token.SRA:
		return c.block.NewAShr(l, r)
	case token.ROL:
		x := c.block.NewShl(l, r)
		d := c.block.NewSub(wordSizeI64, r)
		y := c.block.NewLShr(l, d)
		return c.block.NewOr(x, y)
	case token.ROR:
		x := c.block.NewLShr(l, r)
		d := c.block.NewSub(wordSizeI64, r)
		y := c.block.NewShl(l, d)
		return c.block.NewOr(x, y)

	case token.EQU:
		return c.block.NewICmp(enum.IPredEQ, l, r)
	case token.NEQ:
		return c.block.NewICmp(enum.IPredNE, l, r)
	case token.LT:
		return c.block.NewICmp(enum.IPredSLT, l, r)
	case token.LE:
		return c.block.NewICmp(enum.IPredSLE, l, r)
	case token.GT:
		return c.block.NewICmp(enum.IPredSGT, l, r)
	case token.GE:
		return c.block.NewICmp(enum.IPredSGE, l, r)
	}
	return nil
}

func (c *LlvmCodegen) prefixExpr(n *parser.PrefixExpression) value.Value {
	v := c.expr(n.Value)
	switch n.Operator {
	case token.MINUS:
		return c.block.NewSub(zeroI64, v)
	case token.INV:
		return c.block.NewXor(minusOneI64, v)
	case token.NOT:
		return c.block.NewXor(_true, v)
	}
	return nil
}
