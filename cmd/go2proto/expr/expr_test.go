package expr_test

import (
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/block/ftl/cmd/go2proto/expr"
)

func TestBasicConversions(t *testing.T) {
	a := expr.Ref[expr.Basic]("a")
	c1 := expr.Cast(a, "int32")
	c2 := expr.Cast(c1, "int")
	ptr := expr.PointerOf(c2)
	ref := expr.ReferenceFrom(ptr)

	assert.Equal(t, "*(&(int(int32(a))))", ref.Code)

	ptr2 := expr.PointerOf(ptr)
	assert.Equal(t, "&(&(int(int32(a))))", ptr2.Code)

	assert.Equal(t, "*(&(&(int(int32(a)))))", expr.ReferenceFrom(ptr2).Code)
}

func TestFunctions(t *testing.T) {
	t.Run("supports function calls", func(t *testing.T) {
		f := expr.Ref[expr.Func[expr.Pointer[expr.Basic], expr.OrError[expr.Basic]]]("FooFromProto")
		input := expr.Ref[expr.Pointer[expr.Basic]]("p")
		result := expr.FunctionCall(input, f)

		assert.Equal(t, "FooFromProto(p)", result.Code)
	})
	t.Run("supports method calls", func(t *testing.T) {
		i := expr.Ref[expr.Basic]("p")
		p := expr.PointerOf(i)
		pr := expr.ReferenceFrom(p)
		method := expr.Ref[expr.Func[expr.Basic, expr.Pointer[expr.Basic]]]("ToProto")
		r := expr.MethodCall(pr, method)

		assert.Equal(t, "*(&(p)).ToProto()", r.Code)
	})
}
