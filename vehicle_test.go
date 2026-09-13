// Copyright 2026 Peter James Beard. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package enum

import (
	"database/sql/driver"
	"encoding"
	"reflect"
	"testing"
)

// vehicle is an enhanced enum carrying transport metadata.
type vehicle struct {
	Member[vehicle]

	Tires              int
	Passengers         int
	CarbonPerKilometer int
}

// Compile-time interface checks for vehicle's enum and database integrations.
var (
	_ Enum[vehicle]          = vehicle{}
	_ encoding.TextMarshaler = vehicle{}
	_ driver.Valuer          = (*vehicle)(nil)
)

// vehicles is the namespace used by the enhanced enum tests.
type vehicles struct {
	Car     vehicle
	Bus     vehicle
	Bicycle vehicle
}

// transportation is the enhanced enum namespace schema.
var transportation = vehicles{
	Car: vehicle{
		Tires:              4,
		Passengers:         5,
		CarbonPerKilometer: 400,
	},
	Bus: vehicle{
		Tires:              6,
		Passengers:         50,
		CarbonPerKilometer: 800,
	},
	Bicycle: vehicle{
		Tires:              2,
		Passengers:         1,
		CarbonPerKilometer: 0,
	},
}

// TestEnhancedEnumVehicleTires verifies the Tires field survives initialization.
func TestEnhancedEnumVehicleTires(t *testing.T) {
	registry.DeleteDefinition(reflect.TypeFor[vehicles]())
	Vehicle := DefineNamespace[vehicle](transportation)

	tireTests := []struct {
		Name    string
		Vehicle vehicle
		Index   int
		Want    int
		wantErr bool
	}{
		{"Vehicle.Car", Vehicle.Car, 0, 4, false},
		{"Vehicle.Bus", Vehicle.Bus, 1, 6, false},
		{"Vehicle.Bicycle", Vehicle.Bicycle, 2, 2, false},
		{"Vehicle.Trike", vehicle{}, 0, 0, true},
	}

	for _, test := range tireTests {
		t.Run(test.Name, func(t *testing.T) {
			got := test.Vehicle.Tires
			if got != test.Want {
				t.Errorf("%s: got %d, want %d", test.Name, got, test.Want)
			}
		})
	}
}

// TestEnhancedEnumVehiclePassengers verifies the Passengers field survives initialization.
func TestEnhancedEnumVehiclePassengers(t *testing.T) {
	registry.DeleteDefinition(reflect.TypeFor[vehicles]())
	Vehicle := DefineNamespace[vehicle](transportation)

	tireTests := []struct {
		Name    string
		Vehicle vehicle
		Index   int
		Want    int
		wantErr bool
	}{
		{"Vehicle.Car", Vehicle.Car, 0, 5, false},
		{"Vehicle.Bus", Vehicle.Bus, 1, 50, false},
		{"Vehicle.Bicycle", Vehicle.Bicycle, 2, 1, false},
		{"Vehicle.Trike", vehicle{}, 0, 0, true},
	}

	for _, test := range tireTests {
		t.Run(test.Name, func(t *testing.T) {
			got := test.Vehicle.Passengers
			if got != test.Want {
				t.Errorf("%s: got %d, want %d", test.Name, got, test.Want)
			}
		})
	}
}

// TestEnhancedEnumVehicleCarbonPerKilometer verifies the CarbonPerKilometer field survives initialization.
func TestEnhancedEnumVehicleCarbonPerKilometer(t *testing.T) {
	registry.DeleteDefinition(reflect.TypeFor[vehicles]())
	Vehicle := DefineNamespace[vehicle](transportation)

	tireTests := []struct {
		Name    string
		Vehicle vehicle
		Index   int
		Want    int
		wantErr bool
	}{
		{"Vehicle.Car", Vehicle.Car, 0, 400, false},
		{"Vehicle.Bus", Vehicle.Bus, 1, 800, false},
		{"Vehicle.Bicycle", Vehicle.Bicycle, 2, 0, false},
		{"Vehicle.Trike", vehicle{}, 0, 0, true},
	}

	for _, test := range tireTests {
		t.Run(test.Name, func(t *testing.T) {
			got := test.Vehicle.CarbonPerKilometer
			if got != test.Want {
				t.Errorf("%s: got %d, want %d", test.Name, got, test.Want)
			}
		})
	}
}
