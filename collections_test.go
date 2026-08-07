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

		variations := []struct {
			Name   string
			Lookup string
			wantOk bool
		}{
			{Name: "exact", Lookup: test.Lookup, wantOk: test.wantOk},
			{Name: "lowercase", Lookup: strings.ToLower(test.Lookup), wantOk: false},
			{Name: "uppercase", Lookup: strings.ToUpper(test.Lookup), wantOk: false},
		}
		for _, v := range variations {
			t.Run(test.Name+"/"+v.Name, func(t *testing.T) {
				got, ok := enum.ByName(accessControl, v.Lookup)
				if ok != v.wantOk {
					t.Fatalf("ByName(%q): ok=%v, want %v", v.Lookup, ok, v.wantOk)
				}
				if !ok {
					return
				}
				if !enum.Equal(got, test.Want) {
					t.Errorf("ByName(%q): got %v, want %v", v.Lookup, got, test.Want)
				}
				if got.Name() != v.Lookup {
					t.Errorf("Name() = %q, want %q", got.Name(), v.Lookup)
				}
			})
		}
	}
}

func TestCollectionsByIndex(t *testing.T) {
	tests := []struct {
		Name   string
		Index  int
		Want   role
		wantOk bool
	}{
		{"User", 0, accessControl.User, true},
		{"Guest", 1, accessControl.Guest, true},
		{"Admin", 2, accessControl.Admin, true},
		{"Manager", 3, accessControl.Manager, true},
		{"Viewer", 4, accessControl.Viewer, true},
		{"TooHigh", len(enum.Values(accessControl)), role{}, false},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			got, ok := enum.ByIndex(accessControl, test.Index)
			if ok != test.wantOk {
				t.Fatalf("ByIndex(%d): ok=%v, want %v", test.Index, ok, test.wantOk)
			}
			if !ok {
				return
			}
			if !enum.Equal(got, test.Want) {
				t.Errorf("ByIndex(%d): got %v, want %v", test.Index, got, test.Want)
			}
			if got.Index() != test.Index {
				t.Errorf("Index() = %q, want %d", got.Name(), test.Index)
			}
		})
	}
}

func TestCollectionsValues(t *testing.T) {
	tests := []struct {
		Name   string
		Want   role
		wantOk bool
	}{
		{"User", accessControl.User, true},
		{"Guest", accessControl.Guest, true},
		{"Admin", accessControl.Admin, true},
		{"Manager", accessControl.Manager, true},
		{"Viewer", accessControl.Viewer, true},
	}

	members := enum.Values(accessControl)

	if len(members) != len(tests) {
		t.Fatalf("unexpected Values length: got %d, want %d", len(tests), len(members))
	}

	for i, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			role := members[i]
			if !enum.Equal(role, test.Want) {
				t.Errorf("Value did not match expected enum: got %s, want %s", role, test.Want)
			}
		})
	}
}

func TestCollectionsNames(t *testing.T) {
	tests := []struct {
		Name   string
		Want   role
		wantOk bool
	}{
		{"User", accessControl.User, true},
		{"Guest", accessControl.Guest, true},
		{"Admin", accessControl.Admin, true},
		{"Manager", accessControl.Manager, true},
		{"Viewer", accessControl.Viewer, true},
	}

	members := enum.Names(accessControl)

	if len(members) != len(tests) {
		t.Fatalf("unexpected Names length: got %d, want %d", len(tests), len(members))
	}

	for i, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			name := members[i]
			if name != test.Want.Name() {
				t.Errorf("Name did not match expected enum: got %s, want %s", name, test.Want)
			}
		})
	}

}

func TestCollectionsAll(t *testing.T) {
	tests := []struct {
		Name   string
		Want   role
		wantOk bool
	}{
		{"User", accessControl.User, true},
		{"Guest", accessControl.Guest, true},
		{"Admin", accessControl.Admin, true},
		{"Manager", accessControl.Manager, true},
		{"Viewer", accessControl.Viewer, true},
	}

	for e := range enum.All(accessControl) {
		t.Run(e.Name(), func(t *testing.T) {
			if !enum.Equal(e, tests[e.Index()].Want) {
				t.Errorf("iteration did not yield expected enum: got %s, want %s", e, tests[e.Index()].Want)
			}
		})
	}
}

func TestCollectionsEntries(t *testing.T) {
	tests := []struct {
		Name   string
		Want   role
		wantOk bool
	}{
		{"User", accessControl.User, true},
		{"Guest", accessControl.Guest, true},
		{"Admin", accessControl.Admin, true},
		{"Manager", accessControl.Manager, true},
		{"Viewer", accessControl.Viewer, true},
	}

	for name, e := range enum.Entries(accessControl) {
		t.Run(name, func(t *testing.T) {
			if !enum.Equal(e, tests[e.Index()].Want) {
				t.Errorf("iteration did not yield expected enum: got %s, want %s", e, tests[e.Index()].Want)
			}
		})
	}
}
