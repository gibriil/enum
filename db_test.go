package enum

import (
	"database/sql/driver"
	"encoding"
	"testing"
)

var (
	_ Enum                     = color{}
	_ encoding.TextMarshaler   = color{}
	_ encoding.TextUnmarshaler = (*color)(nil)
	_ driver.Valuer            = (*color)(nil)
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

// func TestCarriers_Scan(t *testing.T) {
// 	carrier := Define(carriers{
// 		UPS:     shipping{},
// 		USPS:    shipping{},
// 		FedEx:   shipping{},
// 		DHL:     shipping{},
// 		Digital: shipping{},
// 	})

// 	tests := []struct {
// 		Name    string
// 		Input   any
// 		Want    shipping
// 		wantErr bool
// 	}{
// 		{"scan bytes", []byte("UPS"), carrier.UPS, false},
// 		{"scan int", 3, carrier.DHL, false},
// 		{"scan nil", nil, shipping{}, true},
// 		{"scan invalid type", "Amazon", shipping{}, true},
// 	}

// 	for _, test := range tests {
// 		t.Run(test.Name, func(t *testing.T) {
// 			var carrier shipping
// 			err := carrier.Scan(test.Input)
// 			if (err != nil) != test.wantErr {
// 				t.Errorf("Scan() error = %v, wantErr %v", err, test.wantErr)
// 				return
// 			}
// 			if carrier != test.Want {
// 				t.Errorf("Scan(): got %v, want %v", carrier, test.Want)
// 			}
// 		})
// 	}
// }
