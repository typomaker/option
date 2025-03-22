package option

func Equal[T comparable](l, r Option[T]) bool {
	switch {
	case
		l.IsZero() && !r.IsZero(),
		l.IsNone() && !r.IsNone(),
		l.IsSome() && !r.IsSome(),
		l.IsSome() && !(l.Get() == r.Get()):
		return false
	default:
		return true
	}
}
func EqualFunc[L, R any](l Option[L], r Option[R], fn func(l L, r R) bool) bool {
	switch {
	case
		l.IsZero() && !r.IsZero(),
		l.IsNone() && !r.IsNone(),
		l.IsSome() && !r.IsSome(),
		l.IsSome() && !fn(l.Get(), r.Get()):
		return false
	default:
		return true
	}
}
