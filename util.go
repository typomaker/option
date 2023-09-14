package option

import (
	"time"
)

type (
	Time       = Option[time.Time]
	Duration   = Option[time.Duration]
	Bool       = Option[bool]
	String     = Option[string]
	Int        = Option[int]
	Int8       = Option[int8]
	Int16      = Option[int16]
	Int32      = Option[int32]
	Int64      = Option[int64]
	Uint       = Option[uint]
	Uint8      = Option[uint8]
	Uint16     = Option[uint16]
	Uint32     = Option[uint32]
	Uint64     = Option[uint64]
	Float32    = Option[float32]
	Float64    = Option[float64]
	Complex64  = Option[complex64]
	Complex128 = Option[complex128]
)

// OneOf returns first not zero option.
// If there are no some options, then returns none value.
func OneOf[T any](op ...Option[T]) Option[T] {
	for i := range op {
		if !op[i].IsZero() {
			return op[i]
		}
	}
	return Option[T]{}
}

// GetOf returns value of first some value.
// If there are no some value, then return zero value.
func GetOf[T any](op ...Option[T]) T {
	for i := range op {
		if op[i].IsSome() {
			return op[i].Get()
		}
	}
	var zero T
	return zero
}

// PickOf returns all some values.
func PickOf[T any](op ...Option[T]) []T {
	var v []T
	for i := range op {
		if op[i].IsSome() {
			v = append(v, op[i].Get())
		}
	}
	return v
}

// IsSome returns a true if all passed values is some.
func IsSome(ss ...Someable) (ok bool) {
	for i := range ss {
		if !ss[i].IsSome() {
			return false
		}
	}
	return true
}

// IsNone returns a true if all passed values is none.
func IsNone(nn ...Noneable) (ok bool) {
	for i := range nn {
		if !nn[i].IsNone() {
			return false
		}
	}
	return true
}

// IsZero returns a true if all passed values is zero.
func IsZero(zz ...Zeroable) (ok bool) {
	for i := range zz {
		if !zz[i].IsZero() {
			return false
		}
	}
	return true
}
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
