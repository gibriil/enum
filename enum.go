// Copyright 2026 Peter James Beard. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package enum provides a standard enum interface for creating a list of enums. Enums are indexed, iterable, and namespaced.

Flags or flag based enums are still best handled by [Iota]

# Basic Enum

Declare an Enum type by embedding Member in a struct.

	type Color struct {
		enum.Member
	}

This enum type can now be used to create an Enum list or namespace.

	type Colors struct {
		Red Color
		Green Color
		Blue Color
	}

The enums can now be initialized by passing your struct to the Define function.

	Colors := enum.Define(Colors{})

# Enhanced Enum

Because an enum is just a struct, we can build enhanced enums that carry additional data.

	type Vehicle struct {
		enum.Member

		Tires int
		Passengers int
		CarbonPerKilometer int
	}

These are identical to our basic enum in every other way. Unlike the name and index of enums, the additional data is technically mutable. By convention enhanced enum data should be treated as immutable, though mutation may be desirable in some cases.

Create an Enum list or namespace.

	type Vehicles struct {
		Car Vehicle
		Bus Vehicle
		Bicycle Vehicle
	}

The additional data can then be populated with their non-zero values when passing the struct declaration to Define.

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

Because package enum uses reflection to initialize, it may be advisable to declare your list globally and pass it to Define in the init function
*/
package enum

import (
	"errors"
)

var (
	ErrUninitialized = errors.New("Enum is Zero Value")
)

// Enum is a package sealed interface to identify the enum type
type Enum interface {
	enum()

	Index() int
	Name() string
	String() string
}

type Namespace struct {
	*definition
}
