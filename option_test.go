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
func TestJSON(t *testing.T) {
	t.Run("marshal none", func(t *testing.T) {
		var o = option.None[int]()
		var b, err = o.MarshalJSON()
		require.NoError(t, err)
		require.Equal(t, []byte("null"), b)
	})
	t.Run("marshal some", func(t *testing.T) {
		var o = option.Some(1)
		var b, err = o.MarshalJSON()
		require.NoError(t, err)
		require.Equal(t, []byte("1"), b)
	})
	t.Run("unmarshal none", func(t *testing.T) {
		var o option.Option[int]
		var err = o.UnmarshalJSON([]byte("null"))
		require.NoError(t, err)
		require.True(t, o.IsNone())
		require.True(t, o.IsZero())
	})
	t.Run("unmarshal some", func(t *testing.T) {
		var o option.Option[int]
		var err = o.UnmarshalJSON([]byte("1"))
		require.NoError(t, err)
		require.True(t, o.IsSome())
		require.Equal(t, 1, o.Get())
	})
}
func TestSQL(t *testing.T) {
	t.Run("marshal none", func(t *testing.T) {
		var o = option.None[int]()
		var b, err = o.Value()
		require.NoError(t, err)
		require.Equal(t, nil, b)
	})
	t.Run("marshal some", func(t *testing.T) {
		var o = option.Some(1)
		var b, err = o.Value()
		require.NoError(t, err)
		require.Equal(t, 1, b)
	})
	t.Run("unmarshal none", func(t *testing.T) {
		var o option.Option[int]
		var err = o.Scan(nil)
		require.NoError(t, err)
		require.True(t, o.IsNone())
		require.True(t, o.IsZero())
	})
	t.Run("unmarshal some", func(t *testing.T) {
		var o option.Option[int]
		var err = o.Scan(int64(1))
		require.NoError(t, err)
		require.True(t, o.IsSome())
		require.Equal(t, 1, o.Get())
	})
}
