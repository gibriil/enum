// Copyright 2026 Peter James Beard. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package enum

import (
	"database/sql/driver"
	"errors"
)

// Member is embedded in a struct to mark the struct type as an enum
//
// The zero value of Member is a nil definition signifying the enum is not initialized
type Member struct {
	// Has unexported fields.
	def   *definition
	index int
}

// Initializer is a non-exported interface for reflection type safety
type initializer interface {
	initialize(*definition, int)
}

// Initialize initializes the enum member with its namespace definition
// and sets sets its position index in the list
func (e *Member) initialize(def *definition, index int) {
	e.def = def
	e.index = index
}

// Name returns the enum member name
func (e Member) Name() string {
	return e.def.names[e.index]
}

// String returns the enum member name
func (e Member) String() string {
	return e.def.names[e.index]
}

// Index returns member position in the enum list
func (e Member) Index() int {
	return e.index
}

// IsZero reports whether the member is its zero value
func (e Member) IsZero() bool {
	return e.def == nil
}

// Valid reports whether or not the enum has been initialized
func (e Member) Valid() bool {
	return !e.IsZero()
}

// MarshalText marshals the enum member name
func (e Member) MarshalText() ([]byte, error) {
	if e.IsZero() {
		return []byte("nil"), errors.New("Enum is Zero Value")
	}
	return []byte(e.def.names[e.index]), nil
}

func (e Member) Value() (driver.Value, error) {
	if e.def == nil {
		return nil, errors.New("Enum not found")
	}
	return e.def.names[e.index], nil
}

// enum marks Member as a valid enum implementation.
// It intentionally has no behavior; it seals the Enum interface.
func (e Member) enum() {}
