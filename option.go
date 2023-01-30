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

// Some creates new some value
func Some[T any](v T) Option[T] {
	return Option[T]{vl: v, ok: true}
}

// None creates new none value
func None[T any]() Option[T] {
	return Option[T]{}
}

// Maybe returns some if the value is non zero and nil
func Maybe[T any](value T) Option[T] {
	rv := reflect.ValueOf(value)
	if !rv.IsValid() {
		return None[T]()
	} else if vv, ok := any(value).(Zeroable); ok && vv.IsZero() {
		return None[T]()
	} else if rv.IsZero() {
		return None[T]()
	}
	return Some(value)
}

// SomeAll returns all some options
func SomeAll[T any](op ...Option[T]) []Option[T] {
	for i := 0; i < len(op); i++ {
		if op[i].IsNone() {
			op = append(op[:i], op[i+1:]...)
			i--
		}
	}
	if len(op) == 0 {
		op = nil
	}
	return op
}

// SomeOne returns first some option. If there are no some options, then returns none option
func SomeOne[T any](op ...Option[T]) Option[T] {
	for i := range op {
		if op[i].IsSome() {
			return op[i]
		}
	}
	return None[T]()
}

// GetAll returns all some values
func GetAll[T any](op ...Option[T]) (some []T) {
	for i := range op {
		if op[i].IsSome() {
			some = append(some, op[i].Get())
		}
	}
	return some
}

// GetOne returns first some value. If there are no some value, then return zero value
func GetOne[T any](op ...Option[T]) (some T) {
	for i := range op {
		if op[i].IsSome() {
			return op[i].Get()
		}
	}
	return some
}

// IsSome returns a true if all passed values is some
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

// IsNone returns a true if all passed values is none
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

// IsZero returns a true if all passed values is zero
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
