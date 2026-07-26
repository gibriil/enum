// Copyright 2026 Peter James Beard. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package enum

import (
	"database/sql/driver"
	"encoding"
	"testing"
)

// Declares type event as an enum
type event struct {
	Member

	Detail any
}

// Ensures event type satisfies Enum interface
var (
	_ Enum                   = event{}
	_ encoding.TextMarshaler = event{}
	_ driver.Valuer          = (*event)(nil)
)

// KeyboardEvent namespace for event enums
type keyboardEvents struct {
	Enter event
	Focus event
}

var keyboardEvent = keyboardEvents{
	Enter: event{
		Detail: map[string]any{
			"id": 1,
		},
	},
	Focus: event{
		Detail: map[string]any{
			"id": 1,
		},
	},
}

// WindowEvent namespace for event enums
type windowEvents struct {
	Move  event
	Focus event
}

var windowEvent = windowEvents{
	Move: event{
		Detail: map[string]any{
			"id": 1,
		},
	},
	Focus: event{
		Detail: map[string]any{
			"id": 1,
		},
	},
}

func init() {
	windowEvent = Define(windowEvent)
	keyboardEvent = Define(keyboardEvent)
}

// Test to ensure that enum at position 0 of one namespace is not equal
// an enum at position 0 of another namespace
func TestEventNamespacesMembersAreDistinct(t *testing.T) {
	if windowEvent.Move == keyboardEvent.Enter {
		t.Error("windowEvent should not equal keyboardEvent")
	}

	if windowEvent.Focus == keyboardEvent.Focus {
		t.Error("The same member names and index should not have equality across namespaces")
	}
}

// Test to ensure that definitions are not equal between two namespaces.
// Test to ensure internal registry only contains namespaces defined in init
func TestEventNamespacesIdentitiesAreDistinct(t *testing.T) {

	if windowEvent.Move.def.identity == keyboardEvent.Enter.def.identity {
		t.Errorf("windowEvent and keyboard should not have the same type: got %s", windowEvent.Move.def.identity)
	}

	if windowEvent.Move.def.name == keyboardEvent.Enter.def.name {
		t.Errorf("windowEvent and keyboard should not have the same name: got %s", windowEvent.Move.def.name)
	}

	if len(registry) != 2 {
		t.Error("Registry length does not matched number of defined Enums")
	}
}
