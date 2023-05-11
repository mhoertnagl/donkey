package llvm2

import (
	"github.com/mhoertnagl/donkey/cgen"
	"github.com/mhoertnagl/donkey/cgen/llvm2/pass"
	"github.com/mhoertnagl/donkey/parser"
)

type LlvmCodegen struct {
}

func NewLlvmCodegen() cgen.Codegen {
	return &LlvmCodegen{}
}

func (c *LlvmCodegen) Generate(n *parser.Program) string {
	astPass := pass.NewAstPass()
	mod := astPass.Run(n)
	return mod.Gen()
}
