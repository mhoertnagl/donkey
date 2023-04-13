package llvm_test

import (
	"math/big"
	"testing"

	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/mhoertnagl/donkey/cgen/llvm"
)

func TestScopesLevel1(t *testing.T) {
	ctx := llvm.NewContext()
	ctx.PushScope()
	ctx.Set("x", constant.NewInt(types.I64, 1))

	expectExisting(t, ctx, "x", 1)
	expectNonExisting(t, ctx, "y")
}

func TestScopesLevel2(t *testing.T) {
	ctx := llvm.NewContext()
	ctx.PushScope()
	ctx.Set("x", constant.NewInt(types.I64, 1))
	ctx.PushScope()
	ctx.Set("z", constant.NewInt(types.I64, 2))

	expectExisting(t, ctx, "x", 1)
	expectNonExisting(t, ctx, "y")
	expectExisting(t, ctx, "z", 2)

	ctx.PopScope()

	expectExisting(t, ctx, "x", 1)
	expectNonExisting(t, ctx, "y")
	expectNonExisting(t, ctx, "z")
}

func expectExisting(t *testing.T, ctx *llvm.Context, name string, exp int64) {
	t.Helper()
	switch c := ctx.Get(name).(type) {
	case *constant.Int:
		if c.X.Cmp(big.NewInt(exp)) != 0 {
			t.Errorf("Expected value is [%d] but got [%v]", exp, c.X)
		}
	default:
		t.Errorf("Unexpected type")
	}
}

func expectNonExisting(t *testing.T, ctx *llvm.Context, name string) {
	t.Helper()
	if ctx.Get(name) != nil {
		t.Errorf("[%s] should be undefined.", name)
	}
}
