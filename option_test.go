package option

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"testing"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	jsoniter.ConfigDefault = jsoniter.ConfigCompatibleWithStandardLibrary
	os.Exit(m.Run())
}
func ExampleOption_states() {
	var value Option[string]
	fmt.Println("value is zero, same as undefined:", value.IsZero())

	value = None[string]()
	fmt.Println("value is none, defined and same as null:", value.IsNone())

	value = Some[string]("hello world")
	fmt.Println("value is some, defined and not null:", value.IsSome())

	// Output:
	// value is zero, same as undefined: true
	// value is none, defined and same as null: true
	// value is some, defined and not null: true
}
func ExampleSomeOrNone() {
	fmt.Printf("%#v\n", SomeOrNone(0))
	fmt.Printf("%#v\n", SomeOrNone(1))
	fmt.Printf("%#v\n", SomeOrNone((*string)(nil)))
	fmt.Printf("%#v\n", SomeOrNone("123123"))
	// Output:
	// option.None[int]()
	// option.Some[int](1)
	// option.None[*string]()
	// option.Some[string]("123123")
}
func ExampleSomeOrZero() {
	fmt.Printf("%#v\n", SomeOrZero(0))
	fmt.Printf("%#v\n", SomeOrZero(1))
	fmt.Printf("%#v\n", SomeOrZero((*string)(nil)))
	fmt.Printf("%#v\n", SomeOrZero("123123"))
	// Output:
	// option.Option[int]{}
	// option.Some[int](1)
	// option.Option[*string]{}
	// option.Some[string]("123123")
}
func ExampleOption_GetOrFunc() {
	fmt.Println("none", None[int]().GetOrFunc(func() int { return 1 }))
	fmt.Println("some", Some(2).GetOrFunc(func() int { return 1 }))
	fmt.Println("zero", Option[int]{}.GetOrFunc(func() int { return 3 }))
	// Output:
	// none 1
	// some 2
	// zero 3
}
func ExampleOption_GetOr() {
	fmt.Println("none", None[int]().GetOr(1))
	fmt.Println("some", Some(2).GetOr(1))
	fmt.Println("zero", Option[int]{}.GetOr(3))
	// Output:
	// none 1
	// some 2
	// zero 3
}
func ExampleOption_GetOrZero() {
	fmt.Println("none", None[int]().GetOrZero())
	fmt.Println("some", Some(1).GetOrZero())
	fmt.Println("zero", Option[int]{}.GetOrZero())
	// Output:
	// none 0
	// some 1
	// zero 0
}
func ExampleOption_GetNilable() {
	fmt.Println("none", None[int]().GetNilable())
	fmt.Printf("some %v\n", *Some("1").GetNilable())
	fmt.Println("zero", Option[int]{}.GetNilable())
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
func TestNilable(t *testing.T) {
	require.Equal(t, Option[string]{}, Nilable((*string)(nil)))
	var value = "123"
	require.Equal(t, Some[string]("123"), Nilable(&value))
}

func TestSomeOrZero_zeroable(t *testing.T) {
	var zeroable = []any{
		nil,
		0,
		[]int(nil),
		(*bool)(nil),
		time.Time{},
		Option[string]{},
		bool(false),
		int(0),
		int8(0),
		int16(0),
		int32(0),
		int64(0),
		uint(0),
		uint8(0),
		uint16(0),
		uint32(0),
		uint64(0),
		string(""),
		float32(0),
		float64(0),
		[]bool(nil),
		[]int(nil),
		[]int8(nil),
		[]int16(nil),
		[]int32(nil),
		[]int64(nil),
		[]uint(nil),
		[]uint8(nil),
		[]uint16(nil),
		[]uint32(nil),
		[]uint64(nil),
		[]string(nil),
		[]float32(nil),
		[]float64(nil),
	}
	for i := range zeroable {
		t.Run(fmt.Sprintf("zero from %T %v\n", zeroable[i], zeroable[i]), func(t *testing.T) {
			actual := SomeOrZero(zeroable[i])
			require.True(t, actual.IsZero())
		})
	}
}
func TestSomeOrZero_someable(t *testing.T) {
	var value = 0
	var someable = []any{
		&value,
		1,
		time.Now(),
		true,
		Some(123),
	}
	for i := range someable {
		t.Run(fmt.Sprintf("some from %T %v\n", someable[i], someable[i]), func(t *testing.T) {
			require.True(t, SomeOrZero(someable[i]).IsSome())
		})
	}
}
func TestSomeOrNone_someable(t *testing.T) {
	var refsome = 0
	var someable = []any{
		&refsome,
		1,
		time.Now(),
		true,
		Some(123),
	}
	for i := range someable {
		t.Run(fmt.Sprintf("some from %T %v\n", someable[i], someable[i]), func(t *testing.T) {
			require.True(t, SomeOrNone(someable[i]).IsSome())
		})
	}
}
func TestSomeOrZero_noneable(t *testing.T) {
	var nilsubtype json.RawMessage
	var noneable = []any{
		nil,
		0,
		false,
		[]int(nil),
		nilsubtype,
		(*bool)(nil),
		time.Time{},
		Option[int]{},
	}

	for i := range noneable {
		t.Run(fmt.Sprintf("none from %T %v\n", noneable[i], noneable[i]), func(t *testing.T) {
			require.True(t, SomeOrNone(noneable[i]).IsNone())
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
func TestGoString(t *testing.T) {
	cases := []struct {
		want  string
		value interface{ GoString() string }
	}{
		{`option.Some[int](123)`, Some(123)},
		{`option.Some[int32](123)`, Some(int32(123))},
		{`option.Some[int64](123)`, Some(int64(123))},
		{`option.Some[float32](123.1)`, Some[float32](123.1)},
		{`option.Some[float64](123.1)`, Some(123.1)},
		{`option.Some[string]("foo")`, Some("foo")},
		{`option.None[string]()`, None[string]()},
		{`option.Option[string]{}`, Option[string]{}},
	}
	for i := range cases {
		require.Equal(t, cases[i].want, fmt.Sprintf("%#v", cases[i].value))
	}
}
func TestEqual(t *testing.T) {
	cases := [...]struct {
		hint        string
		expected    bool
		left, right Option[int]
	}{
		{"zero == zero", true,
			Option[int]{}, Option[int]{}},
		{"zero != none", false,
			Option[int]{}, None[int]()},
		{"zero != some", false,
			Option[int]{}, Some[int](1)},

		{"none != zero", false,
			None[int](), Option[int]{}},
		{"none == none", true,
			None[int](), None[int]()},
		{"none != some", false,
			None[int](), Some[int](1)},

		{"some != zero", false,
			Some[int](1), Option[int]{}},
		{"some != none", false,
			Some[int](1), None[int]()},
		{"some == some", true,
			Some[int](1), Some[int](1)},
		{"some != some", false,
			Some[int](1), Some[int](2)},
	}
	for i := range cases {
		t.Run(cases[i].hint, func(t *testing.T) {
			actual := Equal(cases[i].left, cases[i].right)
			require.Equal(t, cases[i].expected, actual)
		})
	}
}
func TestEqualFunc(t *testing.T) {
	var fn = func(l []string, r string) bool {
		for _, v := range l {
			if v == r {
				return true
			}
		}
		return false
	}
	t.Run("==", func(t *testing.T) {
		left := Some([]string{"foo", "bar"})
		right := Some("bar")

		actual := EqualFunc(left, right, fn)
		require.True(t, actual)
	})
	t.Run("!=", func(t *testing.T) {
		left := Some([]string{"foo", "bar"})
		right := Some("buz")

		actual := EqualFunc(left, right, fn)
		require.False(t, actual)
	})
}
