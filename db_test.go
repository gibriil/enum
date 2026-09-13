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

// shipping is a basic enum value used for database integration tests.
type shipping struct {
	Member[shipping]
}

// Compile-time interface checks for shipping's enum and database integrations.
var (
	_ Enum[shipping]         = shipping{}
	_ encoding.TextMarshaler = shipping{}
	_ driver.Valuer          = (*shipping)(nil)
)

// carriers is the namespace used by the database value tests.
type carriers struct {
	UPS     shipping
	USPS    shipping
	FedEx   shipping
	DHL     shipping
	Digital shipping
}

// TestCarriers_Value verifies that initialized members return their names to database drivers.
func TestCarriers_Value(t *testing.T) {
	registry.DeleteDefinition(reflect.TypeFor[carriers]())
	carrier := DefineNamespace[shipping](carriers{
		UPS:     shipping{},
		USPS:    shipping{},
		FedEx:   shipping{},
		DHL:     shipping{},
		Digital: shipping{},
	})

	tests := []struct {
		Name    string
		Carrier shipping
		Index   int
		Want    driver.Value
		wantErr bool
	}{
		{"carrier.UPS", carrier.UPS, 0, "UPS", false},
		{"carrier.USPS", carrier.USPS, 1, "USPS", false},
		{"carrier.FedEx", carrier.FedEx, 2, "FedEx", false},
		{"carrier.DHL", carrier.DHL, 0, "DHL", false},
		{"carrier.Digital", carrier.Digital, 0, "Digital", false},
		{"carrier.Amazon", shipping{}, 0, nil, true},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			got, err := test.Carrier.Value()
			if (err != nil) != test.wantErr {
				t.Errorf("Value() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if got != test.Want {
				t.Errorf("%s: got %s, want %s", test.Name, got, test.Want)
			}
		})
	}
}
