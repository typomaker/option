package option

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"

	"github.com/typomaker/option/internal/sql"
)

type (
	Option[T any] struct {
		vl T
		ok bool
	}
	multiple[T any] []Option[T]

	Zeroable interface{ IsZero() bool }
	Someable interface{ IsSome() bool }
	Noneable interface{ IsNone() bool }
)

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)

func Some[T any](v T) Option[T] {
	return Option[T]{vl: v, ok: true}
}
func None[T any]() Option[T] {
	return Option[T]{}
}
func Each[T any](v ...Option[T]) multiple[T] {
	return multiple[T](v)
}
func SomeOf[T any](v ...Option[T]) Option[T] {
	return Each(v...).Some().First()
}
func GetOf[T any](v ...Option[T]) T {
	return Each(v...).Some().First().GetZero()
}

// Wrap returns none for nilled value, otherwise some.
func Wrap[T any](v T) Option[T] {
	switch r := reflect.ValueOf(v); r.Kind() {
	case reflect.Invalid:
		return None[T]()
	case reflect.Pointer:
		if r.IsNil() {
			return None[T]()
		}
	}
	return Some(v)
}

// Unwrap returns value of some, otherwise zero value.
func Unwrap[T any](o Option[T]) (v T) {
	if o.ok {
		return o.vl
	}
	return v
}

// IsSome returns a true for some value
func IsSome(ss ...Someable) (ok bool) {
	for i := range ss {
		if !ss[i].IsSome() {
			return false
		}
	}
	return true
}
func (ls multiple[T]) IsEmpty() bool {
	return len(ls) == 0
}

// IsNone returns a true if value is some
func (o Option[T]) IsSome() (ok bool) {
	return o.ok
}
func (ls multiple[T]) IsSome() bool {
	for i := range ls {
		if !ls[i].IsSome() {
			return false
		}
	}
	return !ls.IsEmpty()
}

// IsNone returns a true for none value
func IsNone(nn ...Noneable) (ok bool) {
	for i := range nn {
		if !nn[i].IsNone() {
			return false
		}
	}
	return true
}

// IsNone returns a true if value is none
func (o Option[T]) IsNone() (ok bool) {
	return !o.ok
}
func (ls multiple[T]) IsNone() bool {
	for i := range ls {
		if !ls[i].IsNone() {
			return false
		}
	}
	return true
}

// IsZero returns a true for none and zero
func IsZero(zz ...Zeroable) (ok bool) {
	for i := range zz {
		if !zz[i].IsZero() {
			return false
		}
	}
	return true
}

// IsZero returns a true if value is zero
func (o Option[T]) IsZero() (ok bool) {
	if !o.ok {
		return true
	}
	if r := reflect.ValueOf(o.vl); !r.IsValid() || r.IsZero() {
		return true
	}
	return false
}
func (ls multiple[T]) IsZero() bool {
	for i := range ls {
		if !ls[i].IsZero() {
			return false
		}
	}
	return true
}

// GetZero returns a value, if it none or zero then returns zero value
func (o Option[T]) GetZero() T {
	if !o.ok {
		var zero T
		return zero
	}
	return o.Get()
}
func (ls multiple[T]) GetZero() []T {
	var some = make([]T, len(ls))
	for i := range ls {
		some[i] = ls[i].GetZero()
	}
	return some
}

// Get returns a value if it some, in other case panics
func (o Option[T]) Get() T {
	if !o.ok {
		var caller string
		if _, file, line, ok := runtime.Caller(1); ok {
			file = strings.Replace(file, basepath, "", 1)
			caller = file + ":" + strconv.Itoa(line)
		}
		panic(fmt.Errorf("option: %T is none in %s", o, caller))
	}
	return o.vl
}
func (ls multiple[T]) Get() []T {
	var some = make([]T, 0, len(ls))
	for i := range ls {
		if ls[i].IsSome() {
			some = append(some, ls[i].Get())
		}
	}
	return some
}
func (ls multiple[T]) Some() multiple[T] {
	var some = make(multiple[T], 0, len(ls))
	for i := range ls {
		if ls[i].IsSome() {
			some = append(some, ls[i])
		}
	}
	return some
}
func (ls multiple[T]) None() multiple[T] {
	var none = make(multiple[T], 0, len(ls))
	for i := range ls {
		if ls[i].IsNone() {
			none = append(none, ls[i])
		}
	}
	return none
}
func (ls multiple[T]) Zero() multiple[T] {
	var zero = make(multiple[T], 0, len(ls))
	for i := range ls {
		if ls[i].IsZero() {
			zero = append(zero, ls[i])
		}
	}
	return zero
}
func (ls multiple[T]) First() Option[T] {
	if len(ls) == 0 {
		return None[T]()
	}
	return ls[0]
}
func (ls multiple[T]) Last() Option[T] {
	if len(ls) == 0 {
		return None[T]()
	}
	return ls[len(ls)-1]
}
func (o Option[T]) MarshalJSON() (b []byte, err error) {
	if !o.ok {
		return json.Marshal(nil)
	}
	return json.Marshal(o.vl)
}
func (ls multiple[T]) MarshalJSON() (b []byte, err error) {
	return json.Marshal([]Option[T](ls.Some()))
}
func (o *Option[T]) UnmarshalJSON(b []byte) (err error) {
	if b == nil || bytes.Equal(b, []byte("null")) {
		return
	}
	if err = json.Unmarshal(b, &o.vl); err != nil {
		return
	}
	o.ok = true
	return
}
func (o Option[T]) Value() (val sql.Value, err error) {
	if !o.ok {
		return nil, nil
	}
	if val, err = sql.Marshal(o.vl); err != nil {
		return val, fmt.Errorf("option: value from %T: %w", o.vl, err)
	}
	return val, nil
}
func (o *Option[T]) Scan(src any) (err error) {
	if src == nil {
		*o = Option[T]{}
		return nil
	}
	if err = sql.Unmarshal(src, &o.vl); err != nil {
		return fmt.Errorf("option: scan from %T to %T: %w", src, o.vl, err)
	}
	o.ok = true
	return
}
