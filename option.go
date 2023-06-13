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
//      var zero = option.Value{}
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
func Nil[T any](v *T) Option[T] {
	if v == nil {
		return None[T]()
	}
	return Some(*v)
}

// Maybe returns some if the value is non zero and nil.
//
//	Example:
//		log.Println(Maybe(0)) // None[int]()
//		log.Println(Maybe(1)) // Some[int](1)
//		log.Println(Maybe((*string)(nil))) // None[*string]()
//		ptr := "foo"
//		log.Println(Maybe(&ptr)) // Some[*string](foo)
func Maybe[T any](value T) Option[T] {
	if vv, ok := any(value).(Noneable); ok && vv.IsNone() {
		return None[T]()
	} else if vv, ok := any(value).(Zeroable); ok && vv.IsZero() {
		return None[T]()
	} else if rv := reflect.ValueOf(value); !rv.IsValid() || rv.IsZero() {
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

// GetOrNil returns the nil value if the option is none.
// Pointer is refers to a copy of the origin value,
// so that means any changes to the pointer don't affect the value of the option.
func (o Option[T]) GetOrNil() *T {
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
		return fmt.Sprintf("Value[%T]{}", o.element)
	}
	if o.IsNone() {
		return fmt.Sprintf("None[%T]()", o.element)
	}
	return fmt.Sprintf("Some[%T](%#v)", o.element, o.element)
}

// MarshalJSON is a implementation of the json.Marshaler.
func (o Option[T]) MarshalJSON() (b []byte, err error) {
	if o.IsZero() || o.IsNone() {
		return json.Marshal(nil)
	}
	return json.Marshal(o.element)
}

// UnmarshalJSON is a implementation of the json.Unmarshaler.
func (o *Option[T]) UnmarshalJSON(b []byte) (err error) {
	if b == nil {
		return
	}
	if bytes.Equal(b, []byte("null")) {
		o.defined = true
		return
	}
	if err = json.Unmarshal(b, &o.element); err != nil {
		return
	}
	o.valued = true
	o.defined = true
	return
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
