package option

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
)

type (
	Immutable[T any] interface {
		driver.Valuer
		json.Marshaler
		// IsSome returns a true for some value
		IsSome() (ok bool)
		// IsNone returns a true for none value
		IsNone() (ok bool)
		// IsZero returns a true for none and zero
		IsZero() (ok bool)
		// Must returns a value for some and panics for none
		Must() (value T)
		// Mutable returns a value that can be changed
		Mutable() Mutable[T]
	}
	Mutable[T any] interface {
		Immutable[T]
		// Some set as value
		Some(v T)
		// None set as value
		None()
		// Immutable returns a value that cannot be changed
		Immutable() Immutable[T]
	}

	Zeroable interface{ IsZero() bool }
	Someable interface{ IsSome() bool }
	Noneable interface{ IsNone() bool }

	none[T any]  struct{}
	some[T any]  struct{ v T }
	every[T any] struct{ o Immutable[T] }
)

var (
	_ Immutable[any] = (*none[any])(nil)
	_ Immutable[any] = (*some[any])(nil)
	_ Mutable[any]   = (*every[any])(nil)

	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)

func (e *every[T]) Some(v T) {
	e.o = Some(v)
}
func Some[T any](v T) Immutable[T] {
	return some[T]{v}
}
func (e *every[T]) None() {
	e.o = None[T]()
}
func None[T any]() Immutable[T] {
	return none[T]{}
}
func Every[T any](o Immutable[T]) Mutable[T] {
	return &every[T]{o}
}
func Wrap[T any](v T) Immutable[T] {
	r := reflect.ValueOf(v)
	if !r.IsValid() {
		return None[T]()
	}
	switch r.Kind() {
	case reflect.Pointer, reflect.UnsafePointer:
		if r.IsNil() {
			return None[T]()
		}
		return Some(v)
	default:
		return Some(v)
	}
}
func Unwrap[T any](o Immutable[T]) (v T) {
	if o.IsSome() {
		v = o.Must()
	}
	return v
}
func (e every[T]) IsSome() bool {
	return e.o.IsSome()
}
func (none[T]) IsSome() bool {
	return false
}
func (some[T]) IsSome() bool {
	return true
}
func IsSome(v ...any) bool {
	for i := range v {
		if someable, ok := v[i].(Someable); ok && !someable.IsSome() {
			return false
		}
	}
	return true
}
func (e every[T]) IsNone() bool {
	return e.o.IsNone()
}
func (none[T]) IsNone() bool {
	return true
}
func (some[T]) IsNone() bool {
	return false
}
func IsNone(v ...any) bool {
	for i := range v {
		if noneable, ok := v[i].(Noneable); ok && !noneable.IsNone() {
			return false
		}
	}
	return true
}
func (e every[T]) IsZero() bool {
	return e.o.IsZero()
}
func (none[T]) IsZero() bool {
	return true
}
func (o some[T]) IsZero() bool {
	return IsZero(o.v)
}
func IsZero(v ...any) bool {
	for i := range v {
		if zeroable, ok := v[i].(Zeroable); ok && zeroable.IsZero() {
			return true
		}
		r := reflect.ValueOf(v[i])
		if !r.IsValid() || r.IsZero() {
			return true
		}
	}
	return false
}
func (e every[T]) Must() T {
	return e.o.Must()
}
func (o none[T]) Must() T {
	var caller string
	if _, file, line, ok := runtime.Caller(1); ok {
		file = strings.Replace(file, basepath, "", 1)
		caller = file + ":" + strconv.Itoa(line)
	}
	panic(fmt.Errorf("option: %T is none in %s", o, caller))
}
func (o some[T]) Must() T {
	return o.v
}
func (e every[T]) MarshalJSON() (b []byte, err error) {
	return e.o.MarshalJSON()
}
func (none[T]) MarshalJSON() (b []byte, err error) {
	return []byte("null"), nil
}
func (o some[T]) MarshalJSON() (b []byte, err error) {
	return json.Marshal(o.v)
}
func (e every[T]) Value() (driver.Value, error) {
	return e.o.Value()
}
func (o none[T]) Value() (driver.Value, error) {
	return nil, nil
}
func (o some[T]) Value() (driver.Value, error) {
	if vo, ok := any(o.v).(driver.Valuer); ok {
		return vo.Value()
	}
	return o.v, nil
}
func (e every[T]) Immutable() Immutable[T] {
	return e.o
}
func (e *every[T]) Mutable() Mutable[T] {
	return e
}
func (e some[T]) Mutable() Mutable[T] {
	return Every[T](e)
}
func (e none[T]) Mutable() Mutable[T] {
	return Every[T](e)
}

func SomeOf[T any](oo ...Immutable[T]) (o Immutable[T]) {
	for _, o = range oo {
		if o.IsSome() {
			break
		}
	}
	return o
}
