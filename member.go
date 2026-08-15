// Copyright 2026 Peter James Beard. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package enum

import (
	"database/sql/driver"
	"reflect"
)

type Entry[T any] struct {
	Value T
	Name  string
}

// Member is embedded in a struct to mark the struct type as an enum
//
// The zero value of Member is a nil definition signifying the enum is not initialized
type Member struct {
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

// identity returns the comparable for enum equality checks
func (e Member) identity() Member {
	return e
}

// Name returns the enum member name
func (e Member) Name() string {
	if !e.Valid() {
		return ""
	}
	return e.def.names[e.index]
}

// String returns the enum member name
func (e Member) String() string {
	if !e.Valid() {
		return ""
	}
	return e.def.names[e.index]
}

// Index returns the index of member's position in the enum list
func (e Member) Index() int {
	return e.index
}

// Valid reports whether or not the enum has been initialized
func (e Member) Valid() bool {
	return e.def != nil
}

// MarshalText marshals the enum member name
func (e Member) MarshalText() ([]byte, error) {
	if !e.Valid() {
		return nil, ErrUninitialized
	}
	return []byte(e.def.names[e.index]), nil
}

// Value allows the driver to handle the name of the enum member
func (e Member) Value() (driver.Value, error) {
	if !e.Valid() {
		return nil, ErrUninitialized
	}
	return e.def.names[e.index], nil
}

// enum marks Member as a valid enum implementation.
// It intentionally has no behavior; it seals the Enum interface.
func (e Member) enum() {}

// Namespace surfaces Enum Namespace with collection functions
func (e Member) Namespace() Namespace {
	if !e.Valid() {
		return Namespace{}
	}
	return Namespace{definition: e.def}
}

// Type returns the cached Type of Enum
func (e Member) Type() reflect.Type {
	if !e.Valid() {
		return nil
	}
	return e.def.memberType
}
