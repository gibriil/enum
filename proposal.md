# Proposal: Enum interface

## Abstract

Enums are a highly requested feature. Many packages exist for code generation that just use the iota convention underneath. Many request enums as a type, which is still a valid discussion for future Go versions, but there could be an intermediate step that fits Go's philosophy.

By providing a package that creates a standard interface for enums, structs could be converted to an enum through composition without introducing new syntax or creating language changing types.

## Background

The idiomatic way to create an enum would be using iota:

```go
const (
    Undefined int = iota    // 0
    Red int                 // 1
    Green int               // 2
    Blue int                // 3
)
```

For type protection you create a type and follow the same approach.

```go
type Color int

const (
    Undefined Color = iota    // 0
    Red Color                 // 1
    Green Color               // 2
    Blue Color                // 3
)
```

There is then a good deal of setup. See [Go by Example](https://gobyexample.com/enums) to start.

Some approaches would skip iota all together and declare each const individually for a non-int type

```go
type Color string

const (
    Undefined Color = ""
    Red Color = "Red"
    Green Color = "Green"
    Blue Color = "Blue"
)
```

This approach still requires the setup.

There is nothing preventing someone from simply type casting their own matching case to now be a member of your enum list. See [Rational](#rationale)

Enums should consist of at least the following key characteristics
1. consist of named values which are usually related and represent a finite set of possible options
2. ensure that variables can only take one of the predefined values, providing a layer of type safety
3. makes the code more readable by replacing numeric or string literals with meaningful names.
4. Self documenting the allowable values a variable can take

## Proposal

Provide a stdlib package that will declare enums in a typed namespace and provide much of the setup currently being duplicated by hand.

For example:

```go

type Color struct {
    enum // some exported package type
}

type Colors struct {
    Undefined Color
    Red Color
    Green Color
    Blue Color
}
```

You have the same amount of initial setup up, but then you simply initialize it and the majority of the extra boilerplate is handled by the package. For example: `enum.Initialize(Colors)`

## Rationale

Here are just a few of the many way the community tries to solve the enum problem.

### Iota

See [background](#background) for examples

- **Strengths**
    - Bitwise operations for creating bitmasks and flags
    - Blank identifier (_) to skip values that shouldn't be assigned to an active enum member.
    - Arithmetic operations for creating sequential yet not necessarily consecutive
    - Immutable

- **Weaknesses**
    - Lack of Type Safety
    - No Identity
    - Options or Values are not iterable
    - No useful string representation
    - Type is not finite

### Other typed const

See [background](#background) for examples

- **Strengths**
    - Immutable
    - Useful string representation

- **Weaknesses**
    - Lack of Type Safety
    - No Identity
    - Options or Values are not iterable
    - Type is not finite

### Namespacing trick

```go
type Color int

var Colors = struct {
    Undefined Color
    Red Color
    Green Color
    Blue Color
}{
    Undefined: 0,
    Red: 1,
    Green: 2,
    Blue: 3
}
```

- **Strengths**
    - Type is namespaced and finite

- **Weaknesses**
    - Lack of Type Safety
    - Mutable
    - Options or Values are not iterable


## Compatibility

[A discussion of the change with regard to the
[compatibility guidelines](https://go.dev/doc/go1compat).]

## Implementation

[A description of the steps in the implementation, who will do them, and when.
This should include a discussion of how the work fits into [Go's release cycle](https://go.dev/wiki/Go-Release-Cycle).]

## Open issues

[A discussion of issues relating to this proposal for which the author does not
know the solution. This section may be omitted if there are none.]



```go
type Fruit struct {
    enum.Member
}

type Fruits struct {
    Apple Fruit
    Banana Fruit
    Coconut Fruit
    Grapes Fruit
}

func main() {
    fmt.Println(PostType.Unknown.Name())

    fmt.Println(PostType.Post.CommentEnabled)

    fmt.Println(PostType.ValueOf("Page").Ordinal())

    fmt.Println(PostType.ValueOf("Note") == PostType.Note)

    for _, postType := range PostType.Values() {
        fmt.Println(postType.CommentEnabled())
    }
}
```