//	Example:
// 		import "github.com/typomaker/option"
// 		// Some value defintion.
// 		var some = option.Some("foo")
// 		fmt.Println(some.IsSome()) // true
// 		fmt.Println(some.GetOrZero()) // foo
// 		fmt.Println(some.GetOr("bar")) // foo
// 		fmt.Println(some.Get()) // foo
//
// 		// None value definition.
// 		var none = option.None[string]()
// 		fmt.Println(none.IsNone()) // true
// 		fmt.Println(none.GetOrZero()) // ""
// 		fmt.Println(none.GetOr("bar")) // bar
// 		fmt.Println(none.Get()) // panic
//
// 		// Zero value definition.
// 		var zero = option.Option[string]{}
//		fmt.Println(zero.IsZero()) // true
// 		fmt.Println(zero.GetOrZero()) // ""
// 		fmt.Println(zero.GetOr("bar")) // bar
// 		fmt.Println(zero.Get()) // panic

package option

import (
	"bytes"
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"

	"log/slog"

	jsoniter "github.com/json-iterator/go"
)

type (
	Option[T any] struct {
		value T
		some  bool
		none  bool
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
	return Option[T]{some: true, value: v}
}
func None[T any]() Option[T] {
	return Option[T]{none: true}
}

// Get returns a value if it some, in other case panics.
func (o Option[T]) Get() T {
	if o.IsSome() {
		return o.value
	}
	var caller string
	if _, file, line, ok := runtime.Caller(1); ok {
		file = strings.Replace(file, basepath, "", 1)
		caller = file + ":" + strconv.Itoa(line)
	}
	if o.IsNone() {
		panic(fmt.Errorf("option: %T is none in %s", o, caller))
	}
	panic(fmt.Errorf("option: %T is zero in %s", o, caller))
}

// GetOrZero returns the zero value if the option is none.
func (o Option[T]) GetOrZero() T {
	var zero T
	return o.GetOr(zero)
}

// GetOr returns the value if the option is none.
func (o Option[T]) GetOr(value T) T {
	if !o.IsSome() {
		return value
	}
	return o.value
}

// GetOrFunc retunrs value from getter if the option is none
func (o Option[T]) GetOrFunc(getter func() T) T {
	if !o.IsSome() {
		return getter()
	}
	return o.value
}
func (o Option[T]) LogValue() slog.Value {
	if o.IsZero() {
		return slog.GroupValue()
	}
	if o.IsNone() {
		return slog.Value{}
	}
	if v, ok := any(o.value).(slog.LogValuer); ok {
		return v.LogValue()
	}
	return slog.AnyValue(o.value)
}
func (o Option[T]) String() string {
	if o.IsZero() {
		return ""
	}
	if s, ok := any(o.value).(fmt.Stringer); ok {
		return s.String()
	}
	rv := reflect.ValueOf(o.value)
	for rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
	}
	return fmt.Sprintf("%v", rv.Interface())
}
func (o Option[T]) GoString() string {
	if o.IsZero() {
		return fmt.Sprintf("option.Option[%T]{}", o.value)
	}
	if o.IsNone() {
		return fmt.Sprintf("option.None[%T]()", o.value)
	}
	return fmt.Sprintf("option.Some[%T](%#v)", o.value, o.value)
}

// MarshalJSON is a implementation of the json.Marshaler.
func (o Option[T]) MarshalJSON() (b []byte, err error) {
	if o.IsZero() {
		return nil, nil
	}
	if o.IsNone() {
		return []byte("null"), nil
	}
	return jsoniter.Marshal(o.value)
}

// UnmarshalJSON is a implementation of the json.Unmarshaler.
func (o *Option[T]) UnmarshalJSON(b []byte) (err error) {
	if b == nil {
		return nil
	}
	if bytes.Equal(b, []byte("null")) {
		o.none = true
		return nil
	}
	if err = jsoniter.Unmarshal(b, &o.value); err != nil {
		return err
	}
	o.some = true
	return nil
}

// IsNone returns a true if value is some.
func (o Option[T]) IsSome() bool {
	return o.some
}

// IsNone returns a true if value is none.
func (o Option[T]) IsNone() bool {
	return o.none
}

// IsZero returns a true if value is zero.
func (o Option[T]) IsZero() bool {
	return !o.none && !o.some
}
