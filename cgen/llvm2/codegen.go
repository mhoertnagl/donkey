package llvm2

import (
	"github.com/mhoertnagl/donkey/cgen"
	"github.com/mhoertnagl/donkey/parser"
)

type LlvmCodegen struct {
}

func NewLlvmCodegen() cgen.Codegen {
	return &LlvmCodegen{}
}

func (c *LlvmCodegen) Generate(n *parser.Program) string {
	return ""
}
