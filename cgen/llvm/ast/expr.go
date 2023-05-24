package ast

import (
	"github.com/llir/llvm/ir/value"
	"github.com/mhoertnagl/donkey/utils"
)

type Expr interface {
	gen() value.Value
}

type Exprs []Expr

func (es Exprs) gen() []value.Value {
	// vs := make([]value.Value, len(es))
	// for i, e := range es {
	// 	vs[i] = e.gen()
	// }
	// return vs
	return utils.Map(es, Expr.gen)
}
