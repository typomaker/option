package option_test

import (
	"encoding/json"
	"fmt"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/typomaker/option"
	. "github.com/typomaker/option"
)

func ExampleOption_states() {
	var value option.Option[string]
	fmt.Println("value is zero, same as undefined:", value.IsZero())

	value = option.None[string]()
	fmt.Println("value is none, defined and same as null:", value.IsNone())

	value = option.Some[string]("hello world")
	fmt.Println("value is some, defined and not null:", value.IsSome())

	// Output:
	// value is zero, same as undefined: true
	// value is none, defined and same as null: true
	// value is some, defined and not null: true
}
func ExampleOption_GetOrFunc() {
	fmt.Println("none", option.None[int]().GetOrFunc(func() int { return 1 }))
	fmt.Println("some", option.Some(2).GetOrFunc(func() int { return 1 }))
	fmt.Println("zero", option.Option[int]{}.GetOrFunc(func() int { return 3 }))
	// Output:
	// none 1
	// some 2
	// zero 3
}
func ExampleOption_GetOr() {
	fmt.Println("none", option.None[int]().GetOr(1))
	fmt.Println("some", option.Some(2).GetOr(1))
	fmt.Println("zero", option.Option[int]{}.GetOr(3))
	// Output:
	// none 1
	// some 2
	// zero 3
}
func ExampleOption_GetOrZero() {
	fmt.Println("none", option.None[int]().GetOrZero())
	fmt.Println("some", option.Some(1).GetOrZero())
	fmt.Println("zero", option.Option[int]{}.GetOrZero())
	// Output:
	// none 0
	// some 1
	// zero 0
}
func ExampleOption_GetOrNil() {
	fmt.Println("none", option.None[int]().GetOrNil())
	fmt.Printf("some %v\n", *option.Some("1").GetOrNil())
	fmt.Println("zero", option.Option[int]{}.GetOrNil())
	// Output:
	// none <nil>
	// some 1
	// zero <nil>
}
func TestCompatible(t *testing.T) {
	require.Implements(t, (*json.Unmarshaler)(nil), &Option[any]{})
	require.Implements(t, (*json.Marshaler)(nil), &Option[any]{})
}
func TestGet(t *testing.T) {
	t.Run("zero", func(t *testing.T) {
		_, _, line, _ := runtime.Caller(0)
		require.PanicsWithError(t, "option: option.Option[int] is none in /option_test.go:"+fmt.Sprintf("%d", line+2), func() {
			Option[int]{}.Get()
		})
	})
	t.Run("none", func(t *testing.T) {
		_, _, line, _ := runtime.Caller(0)
		require.PanicsWithError(t, "option: option.Option[int] is none in /option_test.go:"+fmt.Sprintf("%d", line+2), func() {
			None[int]().Get()
		})
	})
	t.Run("some", func(t *testing.T) {
		var o = Some(1)
		require.NotPanics(t, func() {
			o.Get()
		})
		require.Equal(t, 1, o.Get())
	})
}
func TestIsNone(t *testing.T) {
	t.Run("zero", func(t *testing.T) {
		var vv = []Noneable{
			None[any](),
			None[any](),
			Option[any]{},
		}
		require.False(t, IsNone(vv...))
	})
	t.Run("none", func(t *testing.T) {
		var vv = []Noneable{
			None[any](),
			None[any](),
			None[any](),
		}
		require.True(t, IsNone(vv...))
	})
	t.Run("some", func(t *testing.T) {
		var vv = []Noneable{
			None[any](),
			None[any](),
			Some[any](1),
		}
		require.False(t, IsNone(vv...))
	})
}
func TestIsSome(t *testing.T) {
	t.Run("zero", func(t *testing.T) {
		var vv = []Someable{
			Some[any](""),
			Some[any](0),
			Option[any]{},
		}
		require.False(t, IsSome(vv...))
	})
	t.Run("none", func(t *testing.T) {
		var vv = []Someable{
			Some[any](""),
			Some[any](0),
			None[any](),
		}
		require.False(t, IsSome(vv...))
	})
	t.Run("some", func(t *testing.T) {
		var vv = []Someable{
			Some[any](""),
			Some[any](0),
			Some[any](1),
		}
		require.True(t, IsSome(vv...))
	})
}
func TestIsZero(t *testing.T) {
	t.Run("zero", func(t *testing.T) {
		var vv = []Zeroable{
			Option[any]{},
			Option[any]{},
			Option[any]{},
		}
		require.True(t, IsZero(vv...))
	})
	t.Run("none", func(t *testing.T) {
		var vv = []Zeroable{
			Option[any]{},
			Option[any]{},
			None[any](),
		}
		require.False(t, IsZero(vv...))
	})
	t.Run("some", func(t *testing.T) {
		var vv = []Zeroable{
			Option[any]{},
			Option[any]{},
			Some[any](0),
		}
		require.False(t, IsZero(vv...))
	})
}
func TestNil(t *testing.T) {
	require.Equal(t, None[string](), Nil((*string)(nil)))
	var value = "123"
	require.Equal(t, Some[string]("123"), Nil(&value))
}
func TestMaybe(t *testing.T) {
	var refsome = 0
	var someable = []any{
		&refsome,
		1,
		time.Now(),
		true,
	}
	for i := range someable {
		t.Run(fmt.Sprintln(someable[i]), func(t *testing.T) {
			require.True(t, Maybe(someable[i]).IsSome())
		})
	}
	var noneable = []any{
		nil,
		0,
		time.Time{},
		false,
		(*bool)(nil),
	}

	for i := range noneable {
		t.Run(fmt.Sprintln(noneable[i]), func(t *testing.T) {
			require.True(t, Maybe(noneable[i]).IsNone())
		})
	}
}
func TestJSON(t *testing.T) {
	t.Run("marshal none", func(t *testing.T) {
		var o = None[int]()
		var b, err = o.MarshalJSON()
		require.NoError(t, err)
		require.Equal(t, []byte("null"), b)
	})
	t.Run("marshal some", func(t *testing.T) {
		var o = Some(1)
		var b, err = o.MarshalJSON()
		require.NoError(t, err)
		require.Equal(t, []byte("1"), b)
	})
	t.Run("unmarshal zero", func(t *testing.T) {
		var o Option[int]
		var err = o.UnmarshalJSON(nil)
		require.NoError(t, err)
		require.True(t, o.IsZero())
		require.False(t, o.IsNone())
	})
	t.Run("unmarshal none", func(t *testing.T) {
		var o Option[int]
		var err = o.UnmarshalJSON([]byte("null"))
		require.NoError(t, err)
		require.False(t, o.IsZero())
		require.True(t, o.IsNone())
	})
	t.Run("unmarshal some", func(t *testing.T) {
		var o Option[int]
		var err = o.UnmarshalJSON([]byte("1"))
		require.NoError(t, err)
		require.True(t, o.IsSome())
		require.Equal(t, 1, o.Get())
	})
}
func TestSomeOne(t *testing.T) {
	require.Equal(t, None[int](), SomeOf[int]())
	require.Equal(t, None[int](), SomeOf(None[int]()))
	require.Equal(t, Some(1), SomeOf(Some(1)))
	require.Equal(t, Some(1), SomeOf(Some(1), None[int](), Some(2)))
	require.Equal(t, Some(1), SomeOf(None[int](), Some(1), None[int]()))
}
func TestPickOf(t *testing.T) {
	require.Equal(t, ([]int)(nil), PickOf[int]())
	require.Equal(t, ([]int)(nil), PickOf(None[int]()))
	require.Equal(t, []int{1}, PickOf(Some(1)))
	require.Equal(t, []int{1, 2}, PickOf(Some(1), None[int](), Some(2)))
	require.Equal(t, []int{1}, PickOf(None[int](), Some(1), None[int]()))
}
func TestGetOne(t *testing.T) {
	require.Equal(t, 0, GetOf[int]())
	require.Equal(t, 0, GetOf(None[int]()))
	require.Equal(t, 1, GetOf(Some(1)))
	require.Equal(t, 1, GetOf(None[int](), Some(1), None[int](), Some(2)))
	require.Equal(t, 1, GetOf(None[int](), Some(1), None[int]()))
}
func TestGoString(t *testing.T) {
	cases := []struct {
		want  string
		value interface{ GoString() string }
	}{
		{`Some[int](123)`, Some(123)},
		{`Some[int32](123)`, Some(int32(123))},
		{`Some[int64](123)`, Some(int64(123))},
		{`Some[float32](123.1)`, Some[float32](123.1)},
		{`Some[float64](123.1)`, Some(123.1)},
		{`Some[string]("foo")`, Some("foo")},
		{`None[string]()`, None[string]()},
	}
	for i := range cases {
		require.Equal(t, cases[i].want, cases[i].value.GoString())
	}
}
