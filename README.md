# Enum [![Go Reference](https://pkg.go.dev/badge/github.com/gibriil/enum.svg)](https://pkg.go.dev/github.com/gibriil/enum)

<!-- The enum package aims to become the canonical way to build typed registries in Go. Enums are simply the first and most common registry. -->

The enum package builds typed, named, indexed, and iterable registries. Enums are simply the first and most common registry.

Unlike other Go enum packages, this is not a code generation tool. The package works at runtime and requires zero additional dependencies.

<!-- This package is being proposed to the Go community for consideration to be added in to the Go standard library as the idiomatic way to create Enums. -->

**Flags or flag based enums are still best handled by [Iota](https://go.dev/ref/spec#Iota)**

## Requirements

- Go 1.23 or later

## Installation and usage

The import path for the package is *github.com/gibriil/enum*.

To install it, run:

```bash
go get github.com/gibriil/enum@latest
```

Define a struct-backed enum by embedding a typed `Member`:

```go
type Color struct {
	enum.Member[Color]
}

type Colors struct {
	Red   Color
	Green Color
	Blue  Color
}

var colors = enum.DefineNamespace[Color](Colors{})
var colorNamespace = enum.DefinitionFor[Color, Colors]()
```

For comparable values such as integer constants, use `DefineType`:

```go
type State uint8

const (
	StateIdle State = iota
	StateRunning
)

var states = enum.DefineType(
	enum.As[State]{Name: "idle", Value: StateIdle},
	enum.As[State]{Name: "running", Value: StateRunning},
)
```

Text decoding is namespace-specific. Use the relevant `Namespace`'s
`UnmarshalText` or `Scan` helper; this avoids ambiguity when one enum type is
used by multiple namespaces.

## API Documentation

See: <https://pkg.go.dev/github.com/gibriil/enum>
