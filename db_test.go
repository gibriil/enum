// Copyright 2026 Peter James Beard. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package enum

import (
	"database/sql/driver"
	"encoding"
	"testing"
)

var (
	_ Enum                   = color{}
	_ encoding.TextMarshaler = color{}
	_ driver.Valuer          = (*color)(nil)
)

type shipping struct {
	Member
}

type carriers struct {
	UPS     shipping
	USPS    shipping
	FedEx   shipping
	DHL     shipping
	Digital shipping
}

func TestCarriers_Value(t *testing.T) {
	carrier := Define(carriers{
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
