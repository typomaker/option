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

// IsNone returns a true if value is some
func (o Option[T]) IsSome() (ok bool) {
	return o.ok
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

// GetZero returns a values in any case, even if it none or zero
func GetZero[T any](oo ...Option[T]) (values []T) {
	values = make([]T, len(oo))
	for i := range oo {
		values[i] = oo[i].GetZero()
	}
	return values
}

// Get returns a value if it some, in other case panics
func (o Option[T]) GetZero() T {
	var zero T
	return o.GetFallback(zero)
}

// GetFallback returns vl if current value is not some
func (o Option[T]) GetFallback(vl T) T {
	if o.ok {
		return o.vl
	}
	return vl
}

// GetZero returns a some values
func Get[T any](oo ...Option[T]) (values []T) {
	values = make([]T, 0, len(oo))
	for i := range oo {
		if oo[i].IsSome() {
			values = append(values, oo[i].Get())
		}
	}
	return values
}

// GetZero returns a value if it some, in other case panics
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

// SetZero value if current value is none
func (o *Option[T]) SetZero() {
	var zero T
	o.SetFallback(zero)
}

// SetFallback value if current value is none
func (o *Option[T]) SetFallback(vl T) {
	if o.ok {
		return
	}
	o.Set(vl)
}

// Set value
func (o *Option[T]) Set(v T) {
	o.vl, o.ok = v, true
}
func (o Option[T]) MarshalJSON() (b []byte, err error) {
	if !o.ok {
		return json.Marshal(nil)
	}
	return json.Marshal(o.vl)
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
