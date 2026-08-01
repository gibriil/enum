# Enum

<!-- The enum package aims to become the canonical way to build typed registries in Go. Enums are simply the first and most common registry. -->

The enum package is a way to build typed registries in Go. Enums are simply the first and most common registry.

Unlike other Go enum packages, this is not a code generation tool. The package works at runtime and requires zero additional dependencies.

<!-- This package is being proposed to the Go community for consideration to be added in to the Go standard library as the idiomatic way to create Enums. -->

**Flags or flag based enums are still best handled by [Iota](https://go.dev/ref/spec#Iota)**

## Requirements

- Go 1.26.0 or later (recommended for latest capabilities)
- Go 1.18.0 or later (minimum for Go Generics dependency)

##  Installation and Usage

The import path for the package is *github.com/gibriil/enum*.

To install it, run:

```bash
go get github.com/gibriil/enum@latest
```

## API Documentation

See: <https://pkg.go.dev/github.com/gibriil/enum>