# github.com/typomaker/option
Optional type for the **Golang**.
A generic value wrapper that adds two additional value representations:
 - __Zero__ - Same as `undefined`. Basic value state, persists until None or Some value is assigned.
 - __None__ - Same as `null`. A logically explicit absence of a value.
 - __Some__ - Usual value representation.


## Usage

### Value initialization
```go
import "github.com/typomaker/option"
var value option.Option[string]{}
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

### Value getting
```go
import "github.com/typomaker/option"
var value option.Option[string]{}

value.Get() // returns value if defined not null, otherwise panics.
value.GetOrZero() // always returns value, if undefined or null then zero value.
value.GetOrFunc(func() string {return "fallback"}) // returns a value if defined not null. otherwise, returns the result of the passed function.
```