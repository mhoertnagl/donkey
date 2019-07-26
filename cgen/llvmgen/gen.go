package llvmgen

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"github.com/mhoertnagl/donkey/parser"
)

func (c *llvmCodegen) Generate(node parser.Program) {

}

func (c *llvmCodegen) generateBool(n *parser.Boolean) value.Value {
	return constant.NewBool(n.Value)
}

func (c *llvmCodegen) generateInt(n *parser.Integer) value.Value {
	return constant.NewInt(types.I32, n.Value)
}
