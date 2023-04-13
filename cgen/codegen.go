package cgen

import (
	"github.com/mhoertnagl/donkey/parser"
)

type Codegen interface {
	Generate(node *parser.Program) string
}
