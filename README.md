# github.com/typomaker/option
Optional type for the **Golang**.
A generic value wrapper that adds two additional value representations:
 - __Zero__ - Same as `undefined`. Basic value state, persists until None or Some value is assigned.
 - __None__ - Same as `null`. A logically explicit absence of a value.
 - __Some__ - Defined value.


## Usage
```go
import "github.com/typomaker/option"
// Some value defintion.
var some = option.Some("foo")
fmt.Println(some.IsSome()) // true
fmt.Println(some.GetOrZero()) // foo
fmt.Println(some.GetOr("bar")) // foo
fmt.Println(some.Get()) // foo

// None value definition.
var none = option.None[string]()
fmt.Println(none.IsNone()) // true
fmt.Println(none.GetOrZero()) // ""
fmt.Println(none.GetOr("bar")) // bar
fmt.Println(none.Get()) // panic

// Zero value definition.
var zero = option.Option[string]{}
fmt.Println(zero.IsZero()) // true
fmt.Println(zero.GetOrZero()) // ""
fmt.Println(zero.GetOr("bar")) // bar
fmt.Println(zero.Get()) // panic
```
