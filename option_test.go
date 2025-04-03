package option

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"testing"

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
func TestLogValue(t *testing.T) {
	b := bytes.Buffer{}
	h := slog.NewJSONHandler(&b, &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			switch a.Key {
			case "time", "level", "msg":
				return slog.Attr{}
			}
			return a
		},
	})
	t.Run("zero", func(t *testing.T) {
		l := slog.New(h)
		v := Option[string]{}
		l.Info("foo", slog.Any("v", v))
		require.JSONEq(t, `{}`, b.String())
		b.Reset()
	})
	t.Run("none", func(t *testing.T) {
		l := slog.New(h)
		v := None[string]()
		l.Info("foo", slog.Any("v", v))
		require.JSONEq(t, `{"v":null}`, b.String())
		b.Reset()
	})
	t.Run("some", func(t *testing.T) {
		l := slog.New(h)
		v := Some("bar")
		l.Info("foo", slog.Any("v", v))
		require.JSONEq(t, `{"v":"bar"}`, b.String())
		b.Reset()
	})
	t.Run("implement log valuer", func(t *testing.T) {
		l := slog.New(h)
		v := Some(SlogValue{})
		l.Info("foo", slog.Any("v", v))
		require.JSONEq(t, `{"v":{"foo":"bar"}}`, b.String())
		b.Reset()
	})
}

type SlogValue struct{}

func (SlogValue) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Any("foo", "bar"),
	)
}
