package ctx

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/value"
)

type Symbol interface {
	GetValue() value.Value
}

type ValueSymbol struct {
	value value.Value
}

func NewValueSymbol(value value.Value) *ValueSymbol {
	return &ValueSymbol{value}
}

func (sym *ValueSymbol) GetValue() value.Value {
	return sym.value
}

type FuncSymbol struct {
	fun *ir.Func
}

func (sym *FuncSymbol) GetValue() value.Value {
	return sym.fun
}
