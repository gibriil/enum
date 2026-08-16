// Copyright 2026 Peter James Beard. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package enum_test

import (
	"testing"

	"github.com/gibriil/enum"
)

type serverState int

const (
	stateIdle serverState = iota
	stateConnected
	stateError
	stateRetrying
)

var serverStates enum.Namespace

func init() {
	serverStates = enum.DefineType(
		enum.As[serverState]{Name: "idle", Value: stateIdle},
		enum.As[serverState]{Name: "connected", Value: stateConnected},
		enum.As[serverState]{Name: "error", Value: stateError},
		enum.As[serverState]{Name: "retrying", Value: stateRetrying},
	)
}

func TestConstHasEnum(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Unexpected panic message: %v", r)
		}
	}()

	_ = enum.DefinitionFor[serverState]()
}

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
