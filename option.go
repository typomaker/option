//	Example:
//      // Some value defintion.
//		var some = option.Some("")
//		fmt.Println(some.IsSome()) // true
//		fmt.Println(some.IsNone()) // false
//		fmt.Println(some.IsZero()) // false
//
//      // None value definition.
//      var none = option.None[string]()
//		fmt.Println(some.IsSome()) // false
//		fmt.Println(some.IsNone()) // true
//		fmt.Println(some.IsZero()) // false
//
//      // Zero value definition.
//      var zero = option.Option{}
//		fmt.Println(some.IsSome()) // false
//		fmt.Println(some.IsNone()) // false
//		fmt.Println(some.IsZero()) // true
//

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
)

type (
	Option[T any] struct {
		element T
		valued  bool
		defined bool
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
	return Option[T]{element: v, valued: true, defined: true}
}
func None[T any]() Option[T] {
	return Option[T]{defined: true}
}
func Nilable[T any](v *T) Option[T] {
	if v == nil {
		return Option[T]{}
	}
	return Some(*v)
}

// SomeOrZero returns some if the value is non zero. Otherwise returns Zero.
func SomeOrZero[T any](value T) Option[T] {
	if isZero(value) {
		return Option[T]{}
	}
	return Some(value)
}

// SomeOrNone returns Some if the value is non zero. Otherwise returns None.
func SomeOrNone[T any](value T) Option[T] {
	if isZero(value) {
		return None[T]()
	}
	return Some(value)
}

// Get returns a value if it some, in other case panics.
func (o Option[T]) Get() T {
	if o.IsZero() || o.IsNone() {
		var caller string
		if _, file, line, ok := runtime.Caller(1); ok {
			file = strings.Replace(file, basepath, "", 1)
			caller = file + ":" + strconv.Itoa(line)
		}
		panic(fmt.Errorf("option: %T is none in %s", o, caller))
	}
	return o.element
}

// GetNilable returns the nil value if the option is none.
// Pointer is refers to a copy of the origin value,
// so that means any changes to the pointer don't affect the value of the option.
func (o Option[T]) GetNilable() *T {
	if !o.valued {
		return nil
	}
	var cp = o.element
	return &cp
}

// GetOrZero returns the zero value if the option is none.
func (o Option[T]) GetOrZero() T {
	var zero T
	return o.GetOr(zero)
}

// GetOr returns the value if the option is none.
func (o Option[T]) GetOr(value T) T {
	if !o.valued {
		return value
	}
	return o.element
}

// GetOrFunc retunrs value from getter if the option is none
func (o Option[T]) GetOrFunc(getter func() T) T {
	if !o.valued {
		return getter()
	}
	return o.element
}
func (o Option[T]) String() string {
	if o.IsZero() {
		return ""
	}
	if s, ok := any(o.element).(fmt.Stringer); ok {
		return s.String()
	}
	rv := reflect.ValueOf(o.element)
	for rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
	}
	return fmt.Sprintf("%v", rv.Interface())
}
func (o Option[T]) GoString() string {
	if o.IsZero() {
		return fmt.Sprintf("option.Option[%T]{}", o.element)
	}
	if o.IsNone() {
		return fmt.Sprintf("option.None[%T]()", o.element)
	}
	return fmt.Sprintf("option.Some[%T](%#v)", o.element, o.element)
}

// MarshalJSON is a implementation of the json.Marshaler.
func (o Option[T]) MarshalJSON() (b []byte, err error) {
	if o.IsZero() || o.IsNone() {
		return []byte("null"), nil
	}
	return json.Marshal(o.element)
}

// UnmarshalJSON is a implementation of the json.Unmarshaler.
func (o *Option[T]) UnmarshalJSON(b []byte) (err error) {
	if b == nil {
		return nil
	}
	if bytes.Equal(b, []byte("null")) {
		o.defined = true
		return nil
	}
	if err = json.Unmarshal(b, &o.element); err != nil {
		return err
	}
	o.valued = true
	o.defined = true
	return nil
}

// IsNone returns a true if value is some.
func (o Option[T]) IsSome() bool {
	return o.defined && o.valued
}

// IsNone returns a true if value is none.
func (o Option[T]) IsNone() bool {
	return o.defined && !o.valued
}

// IsZero returns a true if value is zero.
func (o Option[T]) IsZero() bool {
	return !o.defined
}
func isZero(value any) bool {
	switch v := value.(type) {
	case Zeroable:
		return v.IsZero()
	case bool:
		return !v
	case string:
		return v == ""
	case
		int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64, complex64, complex128:
		return isZeroNumber(v)
	case
		[]bool, []string,
		[]int, []int8, []int16, []int32, []int64,
		[]uint, []uint8, []uint16, []uint32, []uint64,
		[]float32, []float64, []complex64, []complex128:
		return isZeroSlice(v)
	default:
		rv := reflect.ValueOf(value)
		return !rv.IsValid() || rv.IsZero()
	}
}
func isZeroNumber(value any) bool {
	switch v := value.(type) {
	case int:
		return v == 0
	case int8:
		return v == 0
	case int16:
		return v == 0
	case int32:
		return v == 0
	case int64:
		return v == 0
	case uint:
		return v == 0
	case uint8:
		return v == 0
	case uint16:
		return v == 0
	case uint32:
		return v == 0
	case uint64:
		return v == 0
	case float32:
		return v == 0
	case float64:
		return v == 0
	case complex64:
		return v == 0
	case complex128:
		return v == 0
	default:
		return false
	}
}
func isZeroSlice(value any) bool {
	switch v := value.(type) {
	case []bool:
		return v == nil
	case []string:
		return v == nil
	case []int:
		return v == nil
	case []int8:
		return v == nil
	case []int16:
		return v == nil
	case []int32:
		return v == nil
	case []int64:
		return v == nil
	case []uint:
		return v == nil
	case []uint8:
		return v == nil
	case []uint16:
		return v == nil
	case []uint32:
		return v == nil
	case []uint64:
		return v == nil
	case []float32:
		return v == nil
	case []float64:
		return v == nil
	case []complex64:
		return v == nil
	case []complex128:
		return v == nil
	default:
		return false
	}
}
