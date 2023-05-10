package ast

import "github.com/llir/llvm/ir/value"

type Expr interface {
	gen() value.Value
}
