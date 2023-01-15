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
func TestGet(t *testing.T) {
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
	t.Run("none of many", func(t *testing.T) {
		var vv = []option.Option[any]{
			option.None[any](),
			option.None[any](),
			option.None[any](),
		}
		require.Len(t, option.Get(vv...), 0)
	})
	t.Run("all of many", func(t *testing.T) {
		var vv = []option.Option[any]{
			option.Some[any](1),
			option.Some[any](2),
			option.Some[any](3),
		}
		require.Len(t, option.Get(vv...), 3)
		require.Equal(t, []any{1, 2, 3}, option.Get(vv...))
	})
	t.Run("last of many", func(t *testing.T) {
		var vv = []option.Option[any]{
			option.None[any](),
			option.None[any](),
			option.Some[any](3),
		}
		require.Len(t, option.Get(vv...), 1)
		require.Equal(t, 3, option.Get(vv...)[0])
	})
}
func TestSet(t *testing.T) {
	var o = option.None[int]()
	o.Set(0)
	require.True(t, o.IsSome())
	require.False(t, o.IsNone())
	require.True(t, o.IsZero())
	require.Equal(t, 0, o.Get())

	t.Run("zero", func(t *testing.T) {
		var op option.Option[int]
		op.Set(0)
		require.True(t, op.IsSome())
		require.False(t, op.IsNone())
		require.True(t, op.IsZero())
		require.Equal(t, 0, op.Get())
	})
}
func TestGetFallback(t *testing.T) {
	t.Run("none", func(t *testing.T) {
		require.Equal(t, 1, option.None[int]().GetFallback(1))
	})
	t.Run("some", func(t *testing.T) {
		require.Equal(t, 1, option.Some(1).GetFallback(2))
	})
}
func TestSetFallback(t *testing.T) {
	t.Run("none", func(t *testing.T) {
		var o = option.None[int]()
		o.SetFallback(1)
		require.Equal(t, 1, o.Get())
	})
	t.Run("some", func(t *testing.T) {
		var o = option.Some(1)
		o.SetFallback(2)
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
		require.Equal(t, x[i].some, o.IsSome(), "%d Wrap(%#v).IsSome()", i, x[i].value)
		require.Equal(t, x[i].none, o.IsNone(), "%d Wrap(%#v)..IsNone()", i, x[i].value)
		require.Equal(t, x[i].zero, o.IsZero(), "%d Wrap(%#v)..IsZero()", i, x[i].value)

		require.Equal(t, x[i].some, option.IsSome(o), "%d option.IsSome(%#v)", i, o)
		require.Equal(t, x[i].none, option.IsNone(o), "%d option.IsNone(%#v)", i, o)
		require.Equal(t, x[i].zero, option.IsZero(o), "%d option.IsZero(%#v)", i, o)

		if x[i].none {
			require.Empty(t, option.Unwrap(o), "%d %#v == option.Unwrap(%#v)", i, x[i].value, o)
		} else {
			require.EqualValues(t, x[i].value, option.Unwrap(o), "%d %#v == option.Unwrap(%#v)", i, x[i].value, o)
		}
	}
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
	t.Run("masrshal nested some", func(t *testing.T) {
		var a = option.Some(1)
		var b = option.Some(a)
		var c = option.Some(b)
		var x, err = c.Value()
		require.NoError(t, err)
		require.Equal(t, 1, x)
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
	t.Run("unmarshal nested some", func(t *testing.T) {
		var a = option.None[int]()
		var b = option.Some(a)
		var c = option.Some(b)
		var err = c.Scan(int64(1))
		require.NoError(t, err)
		require.True(t, c.IsSome())
		require.True(t, c.Get().IsSome())
		require.True(t, c.Get().Get().IsSome())
		require.Equal(t, 1, c.Get().Get().Get())
	})
}
