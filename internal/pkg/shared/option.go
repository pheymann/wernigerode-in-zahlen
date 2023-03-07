package shared

type Option[T any] struct {
	IsSome bool
	Value  T
}

func Some[T any](value T) Option[T] {
	return Option[T]{IsSome: true, Value: value}
}

func None[T any]() Option[T] {
	return Option[T]{IsSome: false}
}

func (opt *Option[T]) ToSome(value T) {
	opt.IsSome = true
	opt.Value = value
}

func (opt *Option[T]) ForEach(f func(T)) {
	if opt.IsSome {
		f(opt.Value)
	}
}

func (opt Option[T]) GetOrElse(alt T) T {
	if opt.IsSome {
		return opt.Value
	}

	return alt
}

func Map[A any, B any](opt Option[A], f func(A) B) Option[B] {
	if opt.IsSome {
		return Some(f(opt.Value))
	}

	return None[B]()
}
