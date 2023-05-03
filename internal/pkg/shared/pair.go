package shared

type Pair[A any, B any] struct {
	First  A
	Second B
}

func NewPair[A any, B any](first A, second B) Pair[A, B] {
	return Pair[A, B]{First: first, Second: second}
}
