package option_test

import (
	"testing"
	"time"

	. "github.com/typomaker/option"
)

func BenchmarkNilableInt(b *testing.B) {
	var some Option[int]
	var zero Option[int]
	var ref = int(1)
	for j := 0; j < b.N; j++ {
		some = Nilable(&ref)
		zero = Nilable((*int)(nil))
	}
	_ = some
	_ = zero
}
func BenchmarkSomeOrZeroZeroable(b *testing.B) {
	var some Option[Time]
	var zero Option[Time]
	for j := 0; j < b.N; j++ {
		some = SomeOrZero(None[time.Time]())
		zero = SomeOrZero(Option[time.Time]{})
	}
	_ = some
	_ = zero
}
func BenchmarkSomeOrZeroBoolSlice(b *testing.B) {
	var some Option[[]bool]
	var zero Option[[]bool]
	for j := 0; j < b.N; j++ {
		some = SomeOrZero([]bool{})
		zero = SomeOrZero([]bool(nil))
	}
	_ = some
	_ = zero
}
func BenchmarkSomeOrZeroBool(b *testing.B) {
	var some Option[bool]
	var zero Option[bool]
	for j := 0; j < b.N; j++ {
		some = SomeOrZero(true)
		zero = SomeOrZero(false)
	}
	_ = some
	_ = zero
}
func BenchmarkSomeOrZeroInt(b *testing.B) {
	var some Option[int]
	var zero Option[int]
	for j := 0; j < b.N; j++ {
		some = SomeOrZero(1)
		zero = SomeOrZero(0)
	}
	_ = some
	_ = zero
}
func BenchmarkSomeOrZeroInt8(b *testing.B) {
	var some Option[int8]
	var zero Option[int8]
	for j := 0; j < b.N; j++ {
		some = SomeOrZero[int8](1)
		zero = SomeOrZero[int8](0)
	}
	_ = some
	_ = zero
}
func BenchmarkSomeOrZeroInt16(b *testing.B) {
	var some Option[int16]
	var zero Option[int16]
	for j := 0; j < b.N; j++ {
		some = SomeOrZero[int16](1)
		zero = SomeOrZero[int16](0)
	}
	_ = some
	_ = zero
}
func BenchmarkIsSome(b *testing.B) {
	var some = Some(1)
	var d1, d2 bool
	for i := 0; i < b.N; i++ {
		d1 = IsSome(some)
	}
	_ = d1
	_ = d2
}
