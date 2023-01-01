package option_test

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/typomaker/option"
)

func TestCompatible(t *testing.T) {
	require.Implements(t, (*json.Unmarshaler)(nil), &option.Option[any]{})
	require.Implements(t, (*json.Marshaler)(nil), &option.Option[any]{})
	require.Implements(t, (*driver.Valuer)(nil), &option.Option[any]{})
	require.Implements(t, (*sql.Scanner)(nil), &option.Option[any]{})
}
func TestMust(t *testing.T) {
	t.Run("none", func(t *testing.T) {
		var o = option.None[int]()
		require.PanicsWithError(t, "option: option.Option[int] is none in /option_test.go:24", func() {
			o.Get()
		})
	})
	t.Run("some", func(t *testing.T) {
		var o = option.Some(1)
		require.NotPanics(t, func() {
			o.Get()
		})
		require.Equal(t, 1, o.Get())
	})
}
func TestIsNone(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		var vv = []option.Noneable{
			option.None[any](),
			option.None[any](),
			option.None[any](),
		}
		require.True(t, option.IsNone(vv...))
	})
	t.Run("false", func(t *testing.T) {
		var vv = []option.Noneable{
			option.None[any](),
			option.None[any](),
			option.Some[any](1),
		}
		require.False(t, option.IsNone(vv...))
	})
}
func TestIsSome(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		var vv = []option.Someable{
			option.Some[any](1),
			option.Some[any](""),
			option.Some[any](false),
		}
		require.True(t, option.IsSome(vv...))
	})
	t.Run("false", func(t *testing.T) {
		var vv = []option.Someable{
			option.Some[any](0),
			option.Some[any](""),
			option.None[any](),
		}
		require.False(t, option.IsSome(vv...))
	})
}
func TestWrap(t *testing.T) {
	var ptr int = 1
	var x = [...]struct {
		value any
		zero  bool
		some  bool
		none  bool
	}{
		{value: nil, zero: true, some: false, none: true},
		{value: []int(nil), zero: true, some: true, none: false},
		{value: []int{}, zero: false, some: true, none: false},
		{value: [3]int{}, zero: true, some: true, none: false},
		{value: [3]int{1}, zero: false, some: true, none: false},
		{value: map[int]int(nil), zero: true, some: true, none: false},
		{value: map[int]int{}, zero: false, some: true, none: false},
		{value: 0, zero: true, some: true, none: false},
		{value: 1, zero: false, some: true, none: false},
		{value: (*int)(nil), zero: true, some: false, none: true},
		{value: &ptr, zero: false, some: true, none: false},
		{value: "", zero: true, some: true, none: false},
		{value: "1", zero: false, some: true, none: false},
		{value: struct{ A int }{}, zero: true, some: true, none: false},
		{value: struct{ A int }{A: 1}, zero: false, some: true, none: false},
		{value: time.Time{}, zero: true, some: true, none: false},
		{value: time.Now(), zero: false, some: true, none: false},
	}
	for i := range x {
		var o = option.Wrap(x[i].value)
		require.Equal(t, x[i].some, o.IsSome(), "%d IsSome(%#v)", i, x[i].value)
		require.Equal(t, x[i].none, o.IsNone(), "%d IsNone(%#v)", i, x[i].value)
		require.Equal(t, x[i].zero, o.IsZero(), "%d IsZero(%#v)", i, x[i].value)

		require.Equal(t, x[i].some, option.IsSome(o), "%d IsSome(%#v)", i, x[i].value)
		require.Equal(t, x[i].none, option.IsNone(o), "%d IsNone(%#v)", i, x[i].value)
		require.Equal(t, x[i].zero, option.IsZero(o), "%d IsZero(%#v)", i, x[i].value)

		if x[i].none {
			require.Empty(t, option.Unwrap(o), "%d %#v == Unwrap(%#v)", i, x[i].value, o)
		} else {
			require.EqualValues(t, x[i].value, option.Unwrap(o), "%d %#v == Unwrap(%#v)", i, x[i].value, o)
		}
	}
}
func TestSomeOf(t *testing.T) {
	t.Run("none", func(t *testing.T) {
		var vv = []option.Option[any]{
			option.None[any](),
			option.None[any](),
			option.None[any](),
		}
		require.True(t, option.SomeOf(vv...).IsNone())
	})
	t.Run("first", func(t *testing.T) {
		var vv = []option.Option[any]{
			option.Some[any](1),
			option.Some[any](2),
			option.Some[any](3),
		}
		require.True(t, option.SomeOf(vv...).IsSome())
		require.Equal(t, 1, option.SomeOf(vv...).Get())
	})
	t.Run("last", func(t *testing.T) {
		var vv = []option.Option[any]{
			option.None[any](),
			option.None[any](),
			option.Some[any](3),
		}
		require.True(t, option.SomeOf(vv...).IsSome())
		require.Equal(t, 3, option.SomeOf(vv...).Get())
	})
}

// func TestDDD(t *testing.T) {
// 	var o = option.Some(0)
// 	o.UnmarshalJSON(nil)
// 	require.Equal(t, true, o.IsNone(), fmt.Sprintf("%T", o))
// }

// func TestScan(t *testing.T) {
// 	t.Run("string-string", func(t *testing.T) {
// 		expected, actual := option.Some(""), option.Option[string]{}
// 		require.NoError(t, actual.Scan(""))
// 		require.Equal(t, expected, actual)
// 	})
// 	t.Run("int64-int", func(t *testing.T) {
// 		expected, actual := option.Some(0), option.Option[int]{}
// 		require.NoError(t, actual.Scan(int64(0)))
// 		require.Equal(t, expected, actual)
// 	})
// 	t.Run("int64-int8", func(t *testing.T) {
// 		expected, actual := option.Some[int8](0), option.Option[int8]{}
// 		require.NoError(t, actual.Scan(int64(0)))
// 		require.Equal(t, expected, actual)
// 	})
// 	t.Run("int64-int16", func(t *testing.T) {
// 		expected, actual := option.Some[int16](0), option.Option[int16]{}
// 		require.NoError(t, actual.Scan(int64(0)))
// 		require.Equal(t, expected, actual)
// 	})
// 	t.Run("int64-int32", func(t *testing.T) {
// 		expected, actual := option.Some[int32](0), option.Option[int32]{}
// 		require.NoError(t, actual.Scan(int64(0)))
// 		require.Equal(t, expected, actual)
// 	})
// 	t.Run("int64-int64", func(t *testing.T) {
// 		expected, actual := option.Some[int64](0), option.Option[int64]{}
// 		require.NoError(t, actual.Scan(int64(0)))
// 		require.Equal(t, expected, actual)
// 	})
// 	t.Run("string-int", func(t *testing.T) {
// 		expected, actual := option.Some(0), option.Option[int]{}
// 		require.NoError(t, actual.Scan("0"))
// 		require.Equal(t, expected, actual)
// 	})
// 	t.Run("string-int8", func(t *testing.T) {
// 		expected, actual := option.Some[int8](0), option.Option[int8]{}
// 		require.NoError(t, actual.Scan("0"))
// 		require.Equal(t, expected, actual)
// 	})
// 	t.Run("string-int16", func(t *testing.T) {
// 		expected, actual := option.Some[int16](0), option.Option[int16]{}
// 		require.NoError(t, actual.Scan("0"))
// 		require.Equal(t, expected, actual)
// 	})
// 	t.Run("string-int32", func(t *testing.T) {
// 		expected, actual := option.Some[int32](0), option.Option[int32]{}
// 		require.NoError(t, actual.Scan("0"))
// 		require.Equal(t, expected, actual)
// 	})
// 	t.Run("string-int64", func(t *testing.T) {
// 		expected, actual := option.Some[int64](0), option.Option[int64]{}
// 		require.NoError(t, actual.Scan("0"))
// 		require.Equal(t, expected, actual)
// 	})
// 	t.Run("string-float32", func(t *testing.T) {
// 		expected, actual := option.Some[float32](0), option.Option[float32]{}
// 		require.NoError(t, actual.Scan("0"))
// 		require.Equal(t, expected, actual)
// 	})
// 	t.Run("string-float64", func(t *testing.T) {
// 		expected, actual := option.Some[float64](0), option.Option[float64]{}
// 		require.NoError(t, actual.Scan("0"))
// 		require.Equal(t, expected, actual)
// 	})
// 	t.Run("string-duration", func(t *testing.T) {
// 		expected, actual := option.Some(time.Second), option.Option[time.Duration]{}
// 		require.NoError(t, actual.Scan(time.Second.String()))
// 		require.Equal(t, expected, actual)
// 	})
// 	t.Run("bool-bool", func(t *testing.T) {
// 		expected, actual := option.Some(true), option.Option[bool]{}
// 		require.NoError(t, actual.Scan(true))
// 		require.Equal(t, expected, actual)
// 	})
// 	t.Run("bool-string", func(t *testing.T) {
// 		expected, actual := option.Some("true"), option.Option[string]{}
// 		require.NoError(t, actual.Scan(true))
// 		require.Equal(t, expected, actual)
// 	})
// 	t.Run("bool-int", func(t *testing.T) {
// 		expected, actual := option.Some(1), option.Option[int]{}
// 		require.NoError(t, actual.Scan(true))
// 		require.Equal(t, expected, actual)
// 	})
// 	t.Run("bool-int8", func(t *testing.T) {
// 		expected, actual := option.Some[int8](1), option.Option[int8]{}
// 		require.NoError(t, actual.Scan(true))
// 		require.Equal(t, expected, actual)
// 	})
// 	t.Run("bool-int16", func(t *testing.T) {
// 		expected, actual := option.Some[int16](1), option.Option[int16]{}
// 		require.NoError(t, actual.Scan(true))
// 		require.Equal(t, expected, actual)
// 	})
// 	t.Run("bool-int32", func(t *testing.T) {
// 		expected, actual := option.Some[int32](1), option.Option[int32]{}
// 		require.NoError(t, actual.Scan(true))
// 		require.Equal(t, expected, actual)
// 	})
// 	t.Run("bool-int64", func(t *testing.T) {
// 		expected, actual := option.Some[int64](1), option.Option[int64]{}
// 		require.NoError(t, actual.Scan(true))
// 		require.Equal(t, expected, actual)
// 	})
// 	t.Run("bool-float32", func(t *testing.T) {
// 		expected, actual := option.Some[float32](1), option.Option[float32]{}
// 		require.NoError(t, actual.Scan(true))
// 		require.Equal(t, expected, actual)
// 	})
// 	t.Run("bool-float64", func(t *testing.T) {
// 		expected, actual := option.Some[float64](1), option.Option[float64]{}
// 		require.NoError(t, actual.Scan(true))
// 		require.Equal(t, expected, actual)
// 	})
// 	t.Run("time-time", func(t *testing.T) {
// 		expected, actual := option.Some(time.Time{}), option.Option[time.Time]{}
// 		require.NoError(t, actual.Scan(time.Time{}))
// 		require.Equal(t, expected, actual)
// 	})
// 	t.Run("time-string", func(t *testing.T) {
// 		expected, actual := option.Some(time.Time{}.Format(time.RFC3339)), option.Option[string]{}
// 		require.NoError(t, actual.Scan(time.Time{}))
// 		require.Equal(t, expected, actual)
// 	})
// 	t.Run("time-string", func(t *testing.T) {
// 		expected, actual := option.Some(time.Time{}.Format(time.RFC3339)), option.Option[string]{}
// 		require.NoError(t, actual.Scan(time.Time{}))
// 		require.Equal(t, expected, actual)
// 	})
// }

// func TestGlobal(t *testing.T) {
// 	type optional struct {
// 		v    any
// 		zero bool
// 		some bool
// 		none bool
// 	}

// 	var x = []optional{
// 		{v: false, zero: true, none: true},
// 		{v: 0, zero: true, none: true},
// 		{v: 0.0, zero: true, none: true},
// 		{v: "", zero: true, none: true},
// 		{v: []any(nil), zero: true, none: true},
// 		{v: struct{ A int }{}, zero: true, none: true},
// 		{v: (*struct{ A int })(nil), zero: true, none: true},
// 		{v: map[any]any(nil), zero: true, none: true},
// 		{v: time.Time{}, zero: true, none: true},

// 		{v: true, zero: false, some: true},
// 		{v: 1, zero: false, some: true},
// 		{v: 1.0, zero: false, some: true},
// 		{v: "1", zero: false, some: true},
// 		{v: []any{}, zero: false, some: true},
// 		{v: struct{ A int }{A: 1}, zero: false, some: true},
// 		{v: map[any]any{}, zero: false, some: true},
// 		{v: time.Now(), zero: false, some: true},

// 		{v: option.Some(false), zero: true, some: true},
// 		{v: option.Some(0), zero: true, some: true},
// 		{v: option.Some(0.0), zero: true, some: true},
// 		{v: option.Some(""), zero: true, some: true},
// 		{v: option.Some([]any(nil)), zero: true, some: true},
// 		{v: option.Some((*struct{ A int })(nil)), zero: true, some: true},
// 		{v: option.Some(map[any]any(nil)), zero: true, some: true},
// 		{v: option.Some(time.Time{}), zero: true, some: true},

// 		{v: option.Some(true), zero: false, some: true},
// 		{v: option.Some(1), zero: false, some: true},
// 		{v: option.Some(1.1), zero: false, some: true},
// 		{v: option.Some("1"), zero: false, some: true},
// 		{v: option.Some([]any{}), zero: false, some: true},
// 		{v: option.Some(struct{ A int }{A: 1}), zero: false, some: true},
// 		{v: option.Some(map[any]any{}), zero: false, some: true},
// 		{v: option.Some(time.Now()), zero: false, some: true},

// 		{v: option.Wrap(false), zero: true, none: true},
// 		{v: option.Wrap(0), zero: true, none: true},
// 		{v: option.Wrap(0.0), zero: true, none: true},
// 		{v: option.Wrap(""), zero: true, none: true},
// 		{v: option.Wrap([]any(nil)), zero: true, none: true},
// 		{v: option.Wrap((*struct{ A int })(nil)), zero: true, none: true},
// 		{v: option.Wrap(map[any]any(nil)), zero: true, none: true},
// 		{v: option.Wrap(time.Time{}), zero: true, none: true},

// 		{v: option.Wrap(true), zero: false, some: true},
// 		{v: option.Wrap(1), zero: false, some: true},
// 		{v: option.Wrap(1.1), zero: false, some: true},
// 		{v: option.Wrap("1"), zero: false, some: true},
// 		{v: option.Wrap([]any{}), zero: false, some: true},
// 		{v: option.Wrap(struct{ A int }{A: 1}), zero: false, some: true},
// 		{v: option.Wrap(map[any]any{}), zero: false, some: true},
// 		{v: option.Wrap(time.Now()), zero: false, some: true},

// 		{v: option.Unwrap(option.Wrap(false)), zero: true, none: true},
// 		{v: option.Unwrap(option.Wrap(0)), zero: true, none: true},
// 		{v: option.Unwrap(option.Wrap(0.0)), zero: true, none: true},
// 		{v: option.Unwrap(option.Wrap("")), zero: true, none: true},
// 		{v: option.Unwrap(option.Wrap([]any(nil))), zero: true, none: true},
// 		{v: option.Unwrap(option.Wrap((*struct{ A int })(nil))), zero: true, none: true},
// 		{v: option.Unwrap(option.Wrap(map[any]any(nil))), zero: true, none: true},
// 		{v: option.Unwrap(option.Wrap(time.Time{})), zero: true, none: true},

// 		{v: option.Unwrap(option.Wrap(true)), zero: false, some: true},
// 		{v: option.Unwrap(option.Wrap(1)), zero: false, some: true},
// 		{v: option.Unwrap(option.Wrap(1.1)), zero: false, some: true},
// 		{v: option.Unwrap(option.Wrap("1")), zero: false, some: true},
// 		{v: option.Unwrap(option.Wrap([]any{})), zero: false, some: true},
// 		{v: option.Unwrap(option.Wrap(struct{ A int }{A: 1})), zero: false, some: true},
// 		{v: option.Unwrap(option.Wrap(map[any]any{})), zero: false, some: true},
// 		{v: option.Unwrap(option.Wrap(time.Now())), zero: false, some: true},
// 	}
// 	for i := range x {
// 		require.Equal(t, x[i].zero, option.IsZero(x[i].v), "iszero: case=%d value=%#v", i, x[i].v)
// 		require.Equal(t, x[i].some, option.IsSome(x[i].v), "issome: case=%d value=%#v", i, x[i].v)
// 		require.Equal(t, x[i].none, option.IsNone(x[i].v), "isnone: case=%d value=%#v", i, x[i].v)
// 	}
// }
// func TestGetUndefined(t *testing.T) {
// 	var o option.Option[int]
// 	require.Panics(t, func() {
// 		o.Get()
// 	})
// }
// func TestGet(t *testing.T) {
// 	var o = option.Some(1)
// 	require.Equal(t, 1, o.Get())
// }
// func TestSet(t *testing.T) {
// 	var o option.Option[int]
// 	o.Set(1)
// 	require.Equal(t, 1, option.Unwrap(o))
// }
// func TestGetDefault(t *testing.T) {
// 	var o option.Option[int]
// 	require.Equal(t, 1, o.GetDefault(1))
// 	o.Set(2)
// 	require.Equal(t, 2, o.GetDefault(1))
// }
// func TestSetDefault(t *testing.T) {
// 	var o option.Option[int]
// 	o.SetDefault(1)
// 	require.Equal(t, 1, option.Unwrap(o))
// 	o.Set(2)
// 	o.SetDefault(1)
// 	require.Equal(t, 2, option.Unwrap(o))
// }
// func TestGetDefaultFunc(t *testing.T) {
// 	var o option.Option[int]
// 	var fn = func() int {
// 		return 1
// 	}
// 	require.Equal(t, 1, o.GetDefaultFunc(fn))
// 	o.Set(2)
// 	require.Equal(t, 2, o.GetDefaultFunc(fn))
// }
// func TestSetDefaultFunc(t *testing.T) {
// 	var o option.Option[int]
// 	var fn = func() int {
// 		return 1
// 	}
// 	o.SetDefaultFunc(fn)
// 	require.Equal(t, 1, option.Unwrap(o))
// 	o.Set(2)
// 	o.SetDefaultFunc(fn)
// 	require.Equal(t, 2, option.Unwrap(o))
// }
// func TestGetFallback(t *testing.T) {
// 	var o option.Option[int]
// 	require.Equal(t, option.Some(1), o.GetFallback(option.Some(1)))
// 	o.Set(2)
// 	require.Equal(t, option.Some(2), o.GetFallback(option.Some(1)))
// }
// func TestSetFallback(t *testing.T) {
// 	var o option.Option[int]
// 	o.SetFallback(option.Some(1))
// 	require.Equal(t, 1, option.Unwrap(o))
// 	o.Set(2)
// 	o.SetFallback(option.Some(1))
// 	require.Equal(t, 2, option.Unwrap(o))
// }
// func TestGetFallbackFunc(t *testing.T) {
// 	var o option.Option[int]
// 	var fn = func() option.Option[int] {
// 		return option.Some(1)
// 	}
// 	require.Equal(t, option.Some(1), o.GetFallbackFunc(fn))
// 	o.Set(2)
// 	require.Equal(t, option.Some(2), o.GetFallbackFunc(fn))
// }
// func TestSetFallbackFunc(t *testing.T) {
// 	var o option.Option[int]
// 	var fn = func() option.Option[int] {
// 		return option.Some(1)
// 	}
// 	o.SetFallbackFunc(fn)
// 	require.Equal(t, 1, option.Unwrap(o))
// 	o.Set(2)
// 	o.SetFallbackFunc(fn)
// 	require.Equal(t, 2, option.Unwrap(o))
// }
// func TestMarshalJSON(t *testing.T) {
// 	t.Run("null", func(t *testing.T) {
// 		var o option.Option[int]
// 		var b, err = o.MarshalJSON()
// 		require.NoError(t, err)
// 		require.Equal(t, []byte("null"), b)
// 	})
// 	t.Run("notnull", func(t *testing.T) {
// 		var o = option.Some(0)
// 		var b, err = o.MarshalJSON()
// 		require.NoError(t, err)
// 		require.Equal(t, []byte("0"), b)
// 	})
// }
// func TestUnmarshalJSON(t *testing.T) {
// 	t.Run("null", func(t *testing.T) {
// 		var o option.Option[int]
// 		var err = o.UnmarshalJSON([]byte("null"))
// 		require.NoError(t, err)
// 		require.False(t, o.IsSome())
// 	})
// 	t.Run("notnull", func(t *testing.T) {
// 		var o option.Option[int]
// 		var err = o.UnmarshalJSON([]byte("1"))
// 		require.NoError(t, err)
// 		require.True(t, o.IsSome())
// 	})
// }
