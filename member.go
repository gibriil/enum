// Copyright 2026 Peter James Beard. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package enum

import (
	"database/sql/driver"
	"reflect"

	"github.com/gibriil/enum/internal"
)

// Member is embedded in a struct to mark the struct type as an enum or is the internally created member for an Enum with an underlying type
//
// The zero value of Member is a nil definition signifying the enum is not initialized
type Member[T any] struct {
	internal.Entry[T]
}

// MarshalText marshals the enum member name
func (e Member[T]) MarshalText() ([]byte, error) {
	if !e.Valid() {
		return nil, ErrUninitialized
	}
	return []byte(e.Name()), nil
}

// Value allows the driver to handle the name of the enum member
func (e Member[T]) Value() (driver.Value, error) {
	if !e.Valid() {
		return nil, ErrUninitialized
	}
	return e.Entry.Name(), nil
}

// enum marks Member as a valid enum implementation.
// It intentionally has no behavior; it seals the Enum interface.
func (e Member[T]) enum() {}

// Namespace surfaces Enum Namespace with collection functions
func (e Member[T]) Namespace() Namespace[T] {
	if !e.Valid() {
		return Namespace[T]{}
	}
	return Namespace[T]{
		Definition: e.Definition(),
	}
}

// Type returns the cached Type of Enum
func (e Member[T]) Type() reflect.Type {
	if !e.Valid() {
		return nil
	}
	return e.Definition().EntryType()
}

// Raw returns the enums MemberAs raw value
func (e Member[T]) Raw() T {
	if !e.Valid() {
		return *new(T)
	}
	return e.Entry.Value()
}
