// Copyright 2026 Peter James Beard. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package enum_test

import (
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
		enum.As[serverState]{Name: "idle", Value: stateIdle},
		enum.As[serverState]{Name: "connected", Value: stateConnected},
		enum.As[serverState]{Name: "error", Value: stateError},
		enum.As[serverState]{Name: "retrying", Value: stateRetrying},
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

// TestEnumEqualsConst verifies lookup of registered underlying values.
func TestEnumEqualsConst(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Unexpected panic message: %v", r)
		}
	}()

	if e := enum.Of(stateIdle); e.Entry.Value() != stateIdle {
		t.Errorf("Enum %s: got %d, want %d", e.Name(), e.Entry.Value(), stateIdle)
	}
	if e := enum.Of(stateConnected); e.Entry.Value() != stateConnected {
		t.Errorf("Enum %s: got %d, want %d", e.Name(), e.Entry.Value(), stateConnected)
	}
	if e := enum.Of(stateError); e.Entry.Value() != stateError {
		t.Errorf("Enum %s: got %d, want %d", e.Name(), e.Entry.Value(), stateError)
	}
	if e := enum.Of(stateRetrying); e.Entry.Value() != stateRetrying {
		t.Errorf("Enum %s: got %d, want %d", e.Name(), e.Entry.Value(), stateRetrying)
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
		enum.As[state]{Name: "idle", Value: idle},
		enum.As[state]{Name: "running", Value: running},
		enum.As[state]{Name: "stopped", Value: stopped},
	)

	got := enum.Of(running)

	if !got.Valid() {
		t.Fatal("Of(running) returned invalid enum")
	}

	if got.Entry.Value() != running {
		t.Fatalf("Value() = %v, want %v", got.Entry.Value(), running)
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
