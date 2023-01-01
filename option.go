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
		value T
		ok    bool
	}

	Zeroable interface{ IsZero() bool }
	Someable interface{ IsSome() bool }
	Noneable interface{ IsNone() bool }
)

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)

func Some[T any](v T) Option[T] {
	return Option[T]{value: v, ok: true}
}
func None[T any]() Option[T] {
	return Option[T]{}
}
func Wrap[T any](v T) Option[T] {
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
func Unwrap[T any](o Option[T]) (v T) {
	if o.ok {
		return o.value
	}
	return v
}
func SomeOf[T any](oo ...Option[T]) (o Option[T]) {
	for _, o = range oo {
		if o.ok {
			break
		}
	}
	return o
}

// IsNone returns a true if value is some
func (o Option[T]) IsSome() bool {
	return o.ok
}

// IsSome returns a true for some value
func IsSome(v ...Someable) bool {
	for i := range v {
		if !v[i].IsSome() {
			return false
		}
	}
	return true
}

// IsNone returns a true if value is none
func (o Option[T]) IsNone() bool {
	return !o.ok
}

// IsNone returns a true for none value
func IsNone(v ...Noneable) bool {
	for i := range v {
		if !v[i].IsNone() {
			return false
		}
	}
	return true
}

// IsZero returns a true if value is zero
func (o Option[T]) IsZero() bool {
	return IsZero(o.value)
}

// IsZero returns a true for none and zero
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

// Get returns a value for some and panics for none
func (o Option[T]) Get() T {
	if !o.ok {
		var caller string
		if _, file, line, ok := runtime.Caller(1); ok {
			file = strings.Replace(file, basepath, "", 1)
			caller = file + ":" + strconv.Itoa(line)
		}
		panic(fmt.Errorf("option: %T is none in %s", o, caller))
	}
	return o.value
}
func (o *Option[T]) Set(v T) {
	o.ok = true
	o.value = v
}

func (o Option[T]) MarshalJSON() (b []byte, err error) {
	if !o.ok {
		return json.Marshal(nil)
	}
	return json.Marshal(o.value)
}
func (o *Option[T]) UnmarshalJSON(b []byte) (err error) {
	if b == nil || bytes.Equal(b, []byte("null")) {
		return
	}
	if err = json.Unmarshal(b, &o.value); err != nil {
		return
	}
	o.ok = true
	return
}
func (o Option[T]) Value() (val sql.Value, err error) {
	if !o.ok {
		return nil, nil
	}
	if val, err = sql.Marshal(o.value); err != nil {
		return val, fmt.Errorf("option: sql value from %T: %w", o.value, err)
	}
	return val, nil
}
func (o *Option[T]) Scan(src any) (err error) {
	if src == nil {
		*o = Option[T]{}
		return nil
	}
	if err = sql.Unmarshal(src, &o.value); err != nil {
		return fmt.Errorf("option: scan from %T to %T: %w", src, o.value, err)
	}
	o.ok = true
	return
}
