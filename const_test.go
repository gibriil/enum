// Copyright 2026 Peter James Beard. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package enum_test

import (
	"reflect"
	"testing"

	"github.com/gibriil/enum"
)

// serverState is the underlying type used by the package-level constant tests.
type serverState int

// Server-state constants registered with DefineType.
const (
	stateIdle serverState = iota
	stateConnected
	stateError
	stateRetrying
)

// init registers the package-level server-state definition.
func init() {
	_ = enum.DefineType(
		enum.Entry("idle", stateIdle),
		enum.Entry("connected", stateConnected),
		enum.Entry("error", stateError),
		enum.Entry("retrying", stateRetrying),
	)
}

// TestConstHasEnum verifies that the constant definition can be retrieved.
func TestConstHasEnum(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Unexpected panic message: %v", r)
		}
	}()

	_ = enum.DefinitionFor[serverState, serverState]()
}

// TestNamespaceFacade verifies that the namespace forwards its public
// collection operations without exposing the internal definition type.
func TestNamespaceFacade(t *testing.T) {
	namespace := enum.DefinitionFor[serverState, serverState]()

	if namespace.Len() != 4 {
		t.Fatalf("Len() = %d, want 4", namespace.Len())
	}
	if got, ok := namespace.ByName("error"); !ok || got != stateError {
		t.Fatalf("ByName(error) = %v, %v; want %v, true", got, ok, stateError)
	}
	if got, ok := namespace.ByIndex(1); !ok || got != stateConnected {
		t.Fatalf("ByIndex(1) = %v, %v; want %v, true", got, ok, stateConnected)
	}
	if got := namespace.Values(); len(got) != 4 || got[0] != stateIdle {
		t.Fatalf("Values() = %v, want four values beginning with %v", got, stateIdle)
	}
	if got := namespace.Names(); len(got) != 4 || got[2] != "error" {
		t.Fatalf("Names() = %v, want four names with error at index 2", got)
	}

	all := 0
	for range namespace.All() {
		all++
	}
	if all != namespace.Len() {
		t.Fatalf("All() yielded %d values, want %d", all, namespace.Len())
	}

	entries := 0
	for name, value := range namespace.Entries() {
		if entries == 0 && (name != "idle" || value != stateIdle) {
			t.Fatalf("Entries() first result = %q, %v; want idle, %v", name, value, stateIdle)
		}
		entries++
	}
	if entries != namespace.Len() {
		t.Fatalf("Entries() yielded %d values, want %d", entries, namespace.Len())
	}

	if namespace.EntryType() != reflect.TypeFor[serverState]() {
		t.Fatalf("EntryType() = %v, want %v", namespace.EntryType(), reflect.TypeFor[serverState]())
	}
	if namespace.Type() != reflect.TypeFor[serverState]() {
		t.Fatalf("Type() = %v, want %v", namespace.Type(), reflect.TypeFor[serverState]())
	}
	if namespace.Name() != "serverState" {
		t.Fatalf("Name() = %q, want serverState", namespace.Name())
	}
}

// TestEnumEqualsConst verifies lookup of registered underlying values.
func TestEnumEqualsConst(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Unexpected panic message: %v", r)
		}
	}()

	if e := enum.Of(stateIdle); e.Raw() != stateIdle {
		t.Errorf("Enum %s: got %d, want %d", e.Name(), e.Raw(), stateIdle)
	}
	if e := enum.Of(stateConnected); e.Raw() != stateConnected {
		t.Errorf("Enum %s: got %d, want %d", e.Name(), e.Raw(), stateConnected)
	}
	if e := enum.Of(stateError); e.Raw() != stateError {
		t.Errorf("Enum %s: got %d, want %d", e.Name(), e.Raw(), stateError)
	}
	if e := enum.Of(stateRetrying); e.Raw() != stateRetrying {
		t.Errorf("Enum %s: got %d, want %d", e.Name(), e.Raw(), stateRetrying)
	}

}

// TestMemberAs verifies the Member wrapper returned for type-backed values.
func TestMemberAs(t *testing.T) {
	type state uint8

	const (
		idle state = iota
		running
		stopped
	)

	states := enum.DefineType(
		enum.Entry("idle", idle),
		enum.Entry("running", running),
		enum.Entry("stopped", stopped),
	)

	got := enum.Of(running)

	if !got.Valid() {
		t.Fatal("Of(running) returned invalid enum")
	}

	if got.Raw() != running {
		t.Fatalf("Value() = %v, want %v", got.Raw(), running)
	}

	if got.Name() != "running" {
		t.Fatalf("Name() = %q, want %q", got.Name(), "running")
	}

	if got.Index() != 1 {
		t.Fatalf("Index() = %d, want 1", got.Index())
	}

	try := enum.Of(states.Values()[1])

	if !enum.Equal(got, try) {
		t.Fatal("Of(running) did not identify the registered member")
	}
}
