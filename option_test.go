package option_test

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"runtime"
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
		_, _, line, _ := runtime.Caller(0)
		require.PanicsWithError(t, "option: option.Option[int] is none in /option_test.go:"+fmt.Sprintf("%d", line+2), func() {
			option.None[int]().Get()
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
func TestIsZero(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		var vv = []option.Zeroable{
			option.Some[any](0),
			option.Some[any](""),
			option.Some[any](false),
		}
		require.True(t, option.IsZero(vv...))
	})
	t.Run("false", func(t *testing.T) {
		var vv = []option.Zeroable{
			option.Some[any](1),
			option.Some[any]("1"),
			option.None[any](),
		}
		require.False(t, option.IsZero(vv...))
	})
}
func TestMaybe(t *testing.T) {
	var refsome = 0
	var someable = []any{
		&refsome,
		1,
		time.Now(),
		true,
	}
	var noneable = []any{
		nil,
		0,
		time.Time{},
		false,
		(*bool)(nil),
	}
	for i := range someable {
		t.Run(fmt.Sprintln(someable[i]), func(t *testing.T) {
			require.True(t, option.Maybe(someable[i]).IsSome())
		})
	}
	for i := range noneable {
		t.Run(fmt.Sprintln(noneable[i]), func(t *testing.T) {
			require.True(t, option.Maybe(noneable[i]).IsNone())
		})
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
func TestSomeAll(t *testing.T) {
	require.Equal(t, ([]option.Option[int])(nil), option.SomeAll[int]())
	require.Equal(t, ([]option.Option[int])(nil), option.SomeAll(option.None[int]()))
	require.Equal(t, []option.Option[int]{option.Some(1)}, option.SomeAll(option.Some(1)))
	require.Equal(t, []option.Option[int]{option.Some(1), option.Some(2)}, option.SomeAll(option.Some(1), option.None[int](), option.Some(2)))
	require.Equal(t, []option.Option[int]{option.Some(1)}, option.SomeAll(option.None[int](), option.Some(1), option.None[int]()))
}
func TestSomeOne(t *testing.T) {
	require.Equal(t, option.None[int](), option.SomeOne[int]())
	require.Equal(t, option.None[int](), option.SomeOne(option.None[int]()))
	require.Equal(t, option.Some(1), option.SomeOne(option.Some(1)))
	require.Equal(t, option.Some(1), option.SomeOne(option.Some(1), option.None[int](), option.Some(2)))
	require.Equal(t, option.Some(1), option.SomeOne(option.None[int](), option.Some(1), option.None[int]()))
}
func TestGetAll(t *testing.T) {
	require.Equal(t, ([]int)(nil), option.GetAll[int]())
	require.Equal(t, ([]int)(nil), option.GetAll(option.None[int]()))
	require.Equal(t, []int{1}, option.GetAll(option.Some(1)))
	require.Equal(t, []int{1, 2}, option.GetAll(option.Some(1), option.None[int](), option.Some(2)))
	require.Equal(t, []int{1}, option.GetAll(option.None[int](), option.Some(1), option.None[int]()))
}
func TestGetOne(t *testing.T) {
	require.Equal(t, 0, option.GetOne[int]())
	require.Equal(t, 0, option.GetOne(option.None[int]()))
	require.Equal(t, 1, option.GetOne(option.Some(1)))
	require.Equal(t, 1, option.GetOne(option.None[int](), option.Some(1), option.None[int](), option.Some(2)))
	require.Equal(t, 1, option.GetOne(option.None[int](), option.Some(1), option.None[int]()))
}
