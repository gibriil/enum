// Copyright 2026 Peter James Beard. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package enum provides indexed, named, and iterable enum definitions.

There are two ways to define an enum:

  - DefineNamespace creates a namespace from a struct whose fields are enum values.
  - DefineType creates a namespace for comparable values such as integer-based constants.

Flags and flag-based values are still best handled with `iota`.

# Struct-backed enums

Embed Member in the enum value type and declare a namespace containing its members.

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

DefineNamespace initializes the enum fields and returns the initialized schema. Use
DefinitionFor when you need the namespace for lookups, iteration, scanning, or
metadata.

Enum values can carry additional data. The embedded Member is initialized along
with the other fields.

	type Vehicle struct {
		enum.Member[Vehicle]

		Tires              int
		Passengers         int
		CarbonPerKilometer int
	}

	type Vehicles struct {
		Car     Vehicle
		Bus     Vehicle
		Bicycle Vehicle
	}

	var vehicles = enum.DefineNamespace[Vehicle](Vehicles{
		Car:     Vehicle{Tires: 4, Passengers: 5, CarbonPerKilometer: 400},
		Bus:     Vehicle{Tires: 6, Passengers: 50, CarbonPerKilometer: 800},
		Bicycle: Vehicle{Tires: 2, Passengers: 1},
	})

# Type-backed enums

DefineType registers comparable values under names and returns their namespace.

	type State uint8

	const (
		StateIdle State = iota
		StateRunning
		StateStopped
	)

	var states = enum.DefineType(
		enum.As[State]{Name: "idle", Value: StateIdle},
		enum.As[State]{Name: "running", Value: StateRunning},
		enum.As[State]{Name: "stopped", Value: StateStopped},
	)

	current := enum.Of(StateRunning)
	_ = current.Name()

The Member returned by Of exposes the registered name, index, namespace, and
underlying value through Raw.

# Text and database decoding

Member implements text marshaling and database value conversion for initialized
members. Decoding is namespace-specific because the same enum type may be used
in more than one namespace. Use Namespace.UnmarshalText or Namespace.Scan from
the namespace that should resolve the name. A nil scan source resets the target
to its zero value.
*/
package enum
