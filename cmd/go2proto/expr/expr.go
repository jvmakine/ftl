package expr

import (
	"fmt"
)

type Expr[T Type] struct {
	Code string
}

func Ref[T Type](code string) Expr[T] {
	return Expr[T]{
		Code: code,
	}
}

func Cast(from Expr[Basic], typ string) Expr[Basic] {
	return Expr[Basic]{
		Code: fmt.Sprintf("%s(%s)", typ, from.Code),
	}
}

func FunctionCall[I Value, O Type](from Expr[I], function Expr[Func[I, O]]) Expr[O] {
	return Expr[O]{
		Code: fmt.Sprintf("%s(%s)", function.Code, from.Code),
	}
}

func PointerOf[U Simple](from Expr[U]) Expr[Pointer[U]] {
	return Expr[Pointer[U]]{
		Code: fmt.Sprintf("&(%s)", from.Code),
	}
}

func ReferenceFrom[U Pointer[T], T Simple](from Expr[U]) Expr[T] {
	return Expr[T]{
		Code: fmt.Sprintf("*(%s)", from.Code),
	}
}

func MethodCall[I Simple, O Type](from Expr[I], function Expr[Func[I, O]]) Expr[O] {
	return Expr[O]{
		Code: fmt.Sprintf("%s.%s()", from.Code, function.Code),
	}
}
