package expr

type Type interface {
	exprType()
}

type Value interface {
	Type
	value()
}

type Simple interface {
	Value
	simple()
}

type Basic interface {
	Simple
	basic()
}

type Pointer[U Type] interface {
	Simple
	pointer() U
}

type Slice[U Simple] interface {
	Value
	slice() U
}

type OrError[U Simple] interface {
	Type
	orError() U
}

type Func[I Value, O Type] interface {
	Type
	function() (I, O)
}
