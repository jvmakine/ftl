package main

import "fmt"

type ConverterExpr interface {
	asCode() string
}

// Pointer tags the type to contain a pointer
type Pointer interface{}

// Nested tags the type to contain a nested value
type Nested interface{}

type Tagged[T Pointer | Nested] interface {
	tagged() T
}

// ValueErr is a (v, err) typed exression
type ValueErr interface {
	ConverterExpr
	valueErr()
}

// SimpleExpr is an expression of a single simple underlying value.
type SimpleExpr interface {
	ConverterExpr
	simpleExpr()
}

// NoErrValue is (value, nil) pair
type NoErrValue[T Pointer | Nested] struct {
	Value SimpleExpr
}

func (e *NoErrValue[T]) valueErr() {}

func (e *NoErrValue[T]) tagged() T { var k T; return k }

func (e *NoErrValue[T]) asCode() string {
	return fmt.Sprintf("%s, nil", e.Value.asCode())
}

// BasicTypeConversionExpr is a int(v.Int) like conversion
type BasicTypeConversionExpr struct {
	Underlying interface {
		SimpleExpr
		Tagged[Nested]
	}
	Type string
}

func (e *BasicTypeConversionExpr) simpleExpr() {}

func (e *BasicTypeConversionExpr) tagged() Nested { var k Nested; return k }

func (e *BasicTypeConversionExpr) asCode() string {
	return fmt.Sprintf("%s(%s)", e.Type, e.Underlying.asCode())
}

// InputExpr is a single value as an input to type conversion
type InputExpr[T Pointer | Nested] struct {
	Value string
}

func (e *InputExpr[T]) simpleExpr() {}

func (e *InputExpr[T]) tagged() T { var k T; return k }

func (e *InputExpr[T]) asCode() string {
	return e.Value
}

// PtrToPtrErrFunctionExpr is a call to *FromProto like function
// using (pointer) input and (pointer, err) output
type PtrToPtrErrFunctionExpr struct {
	Underlying interface {
		SimpleExpr
		Tagged[Pointer]
	}
	Function string
}

func (e *PtrToPtrErrFunctionExpr) valueErr() {}

func (e *PtrToPtrErrFunctionExpr) tagged() Pointer { var k Pointer; return k }

func (e *PtrToPtrErrFunctionExpr) asCode() string {
	return fmt.Sprintf("%s(%s)", e.Function, e.Underlying.asCode())
}

// NestedToNestedErrFunctionExpr is a call to *FromProto like function
// using (nested) input and (nested, err) output
type NestedToNestedErrFunctionExpr struct {
	Underlying interface {
		SimpleExpr
		Tagged[Nested]
	}
	Function string
}

func (e *NestedToNestedErrFunctionExpr) valueErr() {}

func (e *NestedToNestedErrFunctionExpr) tagged() Nested { var k Nested; return k }

func (e *NestedToNestedErrFunctionExpr) asCode() string {
	return fmt.Sprintf("%s(%s)", e.Function, e.Underlying.asCode())
}

// PtrToNestedMethodExpr is a call to a method taking a pointer receiver
// and producing a nested value
type PtrToNestedMethodExpr struct {
	Underlying interface {
		SimpleExpr
		Tagged[Pointer]
	}
	Method string
}

func (e *PtrToNestedMethodExpr) simpleExpr() {}

func (e *PtrToNestedMethodExpr) tagged() Nested { var k Nested; return k }

func (e *PtrToNestedMethodExpr) asCode() string {
	return fmt.Sprintf("%s.%s()", e.Underlying.asCode(), e.Method)
}

// NestedToPtrErrTypedFunction converts a nested value to a pointer type
// and error using an external typed pointer expression
type NestedToPtrErrTypedFunction struct {
	Underlying interface {
		SimpleExpr
		Tagged[Nested]
	}
	TypeExr interface {
		SimpleExpr
		Tagged[Pointer]
	}
	Function string
}

func (e *NestedToPtrErrTypedFunction) valueErr() {}

func (e *NestedToPtrErrTypedFunction) tagged() Pointer { var k Pointer; return k }

func (e *NestedToPtrErrTypedFunction) asCode() string {
	return fmt.Sprintf("%s(%s, %s)", e.Function, e.Underlying.asCode(), e.TypeExr.asCode())
}

type PtrExpr struct {
	Underlying interface {
		SimpleExpr
		Tagged[Nested]
	}
}

func (e *PtrExpr) simpleExpr() {}

func (e *PtrExpr) asCode() string {
	return fmt.Sprintf("&(%s)", e.Underlying.asCode())
}

func (e *PtrExpr) tagged() Pointer { var k Pointer; return k }

type NestedExpr struct {
	Underlying interface {
		SimpleExpr
		Tagged[Pointer]
	}
}

func (e *NestedExpr) simpleExpr() {}

func (e *NestedExpr) asCode() string {
	return fmt.Sprintf("*(%s)", e.Underlying.asCode())
}

func (e *NestedExpr) tagged() Nested { var k Nested; return k }

func asPointer(exp SimpleExpr) interface {
	SimpleExpr
	Tagged[Pointer]
} {
	if e, ok := exp.(interface {
		SimpleExpr
		Tagged[Pointer]
	}); ok {
		return e
	} else if e, ok := exp.(interface {
		SimpleExpr
		Tagged[Nested]
	}); ok {
		return &PtrExpr{Underlying: e}
	}
	panic("should not happen")
}

func asNested(exp SimpleExpr) interface {
	SimpleExpr
	Tagged[Nested]
} {
	if e, ok := exp.(interface {
		SimpleExpr
		Tagged[Nested]
	}); ok {
		return e
	} else if e, ok := exp.(interface {
		SimpleExpr
		Tagged[Pointer]
	}); ok {
		return &NestedExpr{Underlying: e}
	}
	panic("should not happen")
}
