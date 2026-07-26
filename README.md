# Enum

The enum package aims to become the canonical way to build typed registries in Go. Enums are simply the first and most common registry.

Unlike other Go enum packages, this is not a code generation tool. The package works at runtime and requires zero additional dependencies.

This package is being proposed to the Go community for consideration to be added in to the Go standard library as the idiomatic way to create Enums.

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

<!-- ## API Documentation

See: <https://pkg.go.dev/gibriil/enum> -->

## Basic Enum
Declare an Enum type by embedding Member in a struct.

```go
type Color struct {
    enum.Member
}
```

This enum type can now be used to create an Enum list or namespace.

```go
type Colors struct {
    Red Color
    Green Color
    Blue Color
}
```

The enums can now be initialized by passing your struct to the Define function.

```go
Colors := enum.Define(Colors{})
```

## Enhanced Enum

Because an enum is just a struct, we can build enhanced enums that carry additional data.

```go
type Vehicle struct {
    enum.Member

    Tires int
    Passengers int
    CarbonPerKilometer int
}
```

These are identical to our basic enum in every other way. Unlike the name and index of enums, the additional data is technically mutable. By convention enhanced enum data should be treated as immutable, though mutation may be desirable in some cases.

Create an Enum list or namespace.

```go
type Vehicles struct {
    Car Vehicle
    Bus Vehicle
    Bicycle Vehicle
}
```

The additional data can then be populated with their non-zero values when passing the struct declaration to Define.

```go
Vehicle := enum.Define(Vehicles{
    Car: Vehicle{
        Tires: 4,
        Passengers: 5,
        CarbonPerKilometer: 400,
    },
    Bus: Vehicle{
        Tires: 6,
        Passengers: 50,
        CarbonPerKilometer: 800,
    },
    Bicycle: Vehicle{
        Tires: 2,
        Passengers: 1,
        CarbonPerKilometer: 0,
    },
})
```

Because package enum uses reflection to initialize, it may be advisable to declare your list globally and pass it to Define in the init function