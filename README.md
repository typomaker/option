# github.com/typomaker/option
Optional type for the **Golang**.
A generic value wrapper that adds two additional value representations:
 - __Zero__ - Same as `undefined`. Basic value state, persists until None or Some value is assigned.
 - __None__ - Same as `null`. A logically explicit absence of a value.
 - __Some__ - Usual value representation.


## Usage
```go
import "github.com/typomaker/option"
```

### Value initialization and checks.
```go
import "github.com/typomaker/option"
var value option.Option[string]
if value.IsZero() { // if value is undefined.
    value = option.Some("foo") // then define it using Some[string](...).
}
if value.IsSome() { // if value is defined and not null
    value = option.None[string]() // then set it null using None[string]() 
}
if value.IsNone() { // if value is defined and null
    value = option.Option[string]{} // then undefine it with the zero Option struct.
}

```

### Value getting.
```go
var value option.Option[string]

value.Get() // returns value if defined not null, otherwise panics.
value.GetOr("fallback") // returns value if defined not null, otherwise passed value.
value.GetOrZero() // always returns value, if undefined or null then zero value.
value.GetOrFunc(func() string {return "fallback"}) // returns a value if defined not null. otherwise, returns the result of the passed function.
```

### Convert value to optional.
SomeOrNone function returns `Some[T](...)` for non zero and not nil value, otherwise returns `None[T]()`
```go
fmt.Printf("%#v\n", SomeOrNone(0))
fmt.Printf("%#v\n", SomeOrZero(0))
fmt.Printf("%#v\n", SomeOrNone(1))
fmt.Printf("%#v\n", SomeOrNone((*string)(nil)))
fmt.Printf("%#v\n", SomeOrNone("123123"))
// Output:
// option.None[int]()
// option.Option[int]{}
// option.Some[int](1)
// option.None[*string]()
// option.Some[string]("123123")
```

### OneOf helper
Returns __first__ some option from passed.
If passed doesn't contain some value, then returns `None[T]()`.
```go
var value1 = option.None[int]()
var value2 = option.Some[int](1)
var value3 = option.Some[int](2)

fmt.Printf("%#v", option.OneOf(value1, value2, value3))

// Output:
// Some[int](1)
```

### GetOf helper
Returns value of __first__ some option from passed.
```go
var value1 = option.None[int]()
var value2 = option.Some[int](1)
var value3 = option.Some[int](2)

fmt.Printf("%#v", option.GetOf(value1, value2, value3))

// Output:
// 1
```

### PickOf helper
Returns all values of passed options, exclude the zero and none option.
```go
var value1 = option.None[int]()
var value2 = option.Some[int](1)
var value3 = option.None[int]()
var value4 = option.Option[int]{}

fmt.Printf("%#v", option.PickOf(value1, value2, value3, value4))

// Output:
// []int{1, 2}
```