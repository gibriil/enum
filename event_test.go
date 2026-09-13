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

// event is an enhanced enum carrying arbitrary event details.
type event struct {
	Member[event]

	Detail any
}

// Compile-time interface checks for event's enum and database integrations.
var (
	_ Enum[event]            = event{}
	_ encoding.TextMarshaler = event{}
	_ driver.Valuer          = (*event)(nil)
)

// keyboardEvents is one namespace for event values.
type keyboardEvents struct {
	Enter event
	Focus event
}

// keyboardEvent contains the initial keyboard event definitions.
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

// windowEvents is a separate namespace for the same event type.
type windowEvents struct {
	Move  event
	Focus event
}

// windowEvent contains the initial window event definitions.
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

// TestEventNamespacesMembersAreDistinct verifies that equal positions in
// different namespaces still represent different members.
func TestEventNamespacesMembersAreDistinct(t *testing.T) {
	registry.DeleteDefinition(reflect.TypeFor[windowEvents]())
	windowEvent = DefineNamespace[event](windowEvent)
	registry.DeleteDefinition(reflect.TypeFor[keyboardEvents]())
	keyboardEvent = DefineNamespace[event](keyboardEvent)

	if windowEvent.Move.Member == keyboardEvent.Enter.Member {
		t.Error("windowEvent should not equal keyboardEvent")
	}

	if windowEvent.Focus.Member == keyboardEvent.Focus.Member {
		t.Error("The same member names and index should not have equality across namespaces")
	}
}

// TestEventNamespacesIdentitiesAreDistinct verifies that namespace definitions
// retain distinct identities even when they use the same enum value type.
func TestEventNamespacesIdentitiesAreDistinct(t *testing.T) {
	registry.DeleteDefinition(reflect.TypeFor[windowEvents]())
	windowEvent = DefineNamespace[event](windowEvent)
	registry.DeleteDefinition(reflect.TypeFor[keyboardEvents]())
	keyboardEvent = DefineNamespace[event](keyboardEvent)

	if windowEvent.Move.Namespace().identity() == keyboardEvent.Enter.Namespace().identity() {
		t.Errorf("windowEvent and keyboard should not have the same type: got %s", windowEvent.Move.Namespace().identity().Type())
	}

	if windowEvent.Move.Namespace().Name() == keyboardEvent.Enter.Namespace().Name() {
		t.Errorf("windowEvent and keyboard should not have the same name: got %s", windowEvent.Move.Namespace().Name())
	}
}
