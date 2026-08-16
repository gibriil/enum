// Copyright 2026 Peter James Beard. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package enum

import (
	"database/sql/driver"
	"reflect"
)

// Member is embedded in a struct to mark the struct type as an enum
//
// The zero value of Member is a nil definition signifying the enum is not initialized
type Member struct {
	def   *definition
	index int
}

// MemberAs is the internally created member for an Enum with an underlying type
//
// The zero value of MemberAs is a nil definition signifying the enum is not initialized
type MemberAs[T comparable] struct {
	Member
	raw T
}

// As is a package struct required for registering an underlying type as an Enum
//
// See DefineType
type As[T comparable] struct {
	Name  string
	Value T
}

// initializer is a non-exported interface for reflection type safety
type initializer interface {
	initialize(*definition, int)
}

// initialize initializes the enum member with its namespace definition
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

// UnmarshalText un-marshals the data to an Enum with an underlying type
func (e *MemberAs[T]) UnmarshalText(text []byte) error {
	registry.RLock()
	namespace, ok := registry.data[reflect.TypeFor[T]()]
	registry.RUnlock()

	if !ok {
		return ErrNotDefined
	}

	enum, ok := namespace.ByName(string(text))

	if !ok {
		return ErrEnumNotFound
	}

	constEnum, ok := enum.(EnumAs[T])

	*e = constEnum.(MemberAs[T])
	return nil
}

// Value allows the driver to handle the name of the enum member
func (e Member) Value() (driver.Value, error) {
	if !e.Valid() {
		return nil, ErrUninitialized
	}
	return e.def.names[e.index], nil
}

// Value allows the driver to handle the value of the enum member with an underlying type
func (e MemberAs[T]) Value() (driver.Value, error) {
	if !e.Valid() {
		return nil, ErrUninitialized
	}
	return e.raw, nil
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

// Raw returns the enums MemberAs raw value
func (e MemberAs[T]) Raw() T {
	return e.raw
}
