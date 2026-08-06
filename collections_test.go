package enum_test

import (
	"database/sql/driver"
	"encoding"
	"strings"
	"testing"

	"github.com/gibriil/enum"
)

// role represents an Okta-style role with embedded enum, a slice, and a map.
type role struct {
	enum.Member

	// Scope defines the allowed resources or endpoints.
	Scope []string

	// Permissions maps resource identifiers to allowed actions.
	Permissions map[string]string
}

// Ensures role type satisfies Enum interface
var (
	_ enum.Enum              = role{}
	_ encoding.TextMarshaler = role{}
	_ driver.Valuer          = (*role)(nil)
)

// AccessControl groups role assignments following Okta semantics.
var accessControl = struct {
	User    role
	Guest   role
	Admin   role
	Manager role
	Viewer  role
}{
	// Admin: full access
	Admin: role{
		Scope: []string{"read", "write", "delete", "admin"},
		Permissions: map[string]string{
			"users":        "full",
			"applications": "manage",
			"reports":      "export",
			"settings":     "configure",
		},
	},

	// User: standard access
	User: role{
		Scope: []string{"read", "write"},
		Permissions: map[string]string{
			"users":        "read",
			"applications": "access",
			"reports":      "view",
		},
	},

	// Manager: team oversight
	Manager: role{
		Scope: []string{"read", "write", "review"},
		Permissions: map[string]string{
			"users":        "read,write",
			"applications": "assign",
			"reports":      "view,export",
			"team":         "manage",
		},
	},

	// Viewer: read-only access
	Viewer: role{
		Scope: []string{"read"},
		Permissions: map[string]string{
			"users":        "read",
			"applications": "view",
			"reports":      "view",
		},
	},

	// Guest: minimal access
	Guest: role{
		Scope: []string{"read"},
		Permissions: map[string]string{
			"applications": "view",
		},
	},
}

func init() {
	accessControl = enum.Define(accessControl)
}

func TestCollectionEquality(t *testing.T) {

	if enum.Equal(accessControl.Admin, accessControl.Guest) {
		t.Errorf("accessControl.Admin should not equal accessControl.Guest")
	}

	if !enum.Equal(accessControl.Admin, accessControl.Admin) {
		t.Errorf("accessControl.Admin should be equal to itself")
	}
}

func TestCollectionsByName(t *testing.T) {
	tests := []struct {
		Name   string
		Lookup string
		Want   role
		wantOk bool
	}{
		{"User", "User", accessControl.User, true},
		{"Guest", "Guest", accessControl.Guest, true},
		{"Admin", "Admin", accessControl.Admin, true},
		{"Manager", "Manager", accessControl.Manager, true},
		{"Viewer", "Viewer", accessControl.Viewer, true},
		{"Missing", "Reporter", role{}, false},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			for i := range 3 {
				switch i {
				case 1:
					test.Lookup = strings.ToLower(test.Lookup)
					test.wantOk = false
				case 2:
					test.Lookup = strings.ToUpper(test.Lookup)
					test.wantOk = false
				}
				got, ok := enum.ByName(accessControl, test.Lookup)
				if ok != test.wantOk {
					t.Fatalf("ByName(%q): ok=%v, want %v", test.Lookup, ok, test.wantOk)
				}
				if !ok {
					return
				}
				if !enum.Equal(got, test.Want) {
					t.Errorf("ByName(%q): got %v, want %v", test.Lookup, got, test.Want)
				}
				if got.Name() != test.Lookup {
					t.Errorf("Name() = %q, want %q", got.Name(), test.Lookup)
				}
			}
		})
	}
}

func TestCollectionsByIndex(t *testing.T) {

}

func TestCollectionsByValues(t *testing.T) {

}

func TestCollectionsNames(t *testing.T) {

}

func TestCollectionsAll(t *testing.T) {

}

func TestCollectionsEntries(t *testing.T) {

}
