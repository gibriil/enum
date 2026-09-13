// Copyright 2026 Peter James Beard. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package enum

import (
	"errors"
	"fmt"
)

var (
	// ErrUninitialized reports that an enum or member has its zero value.
	ErrUninitialized   = errors.New("enum is Zero Value")
	// ErrNotDefined reports that a requested enum definition is not registered.
	ErrNotDefined      = errors.New("enum has no registered definition")
	// ErrEnumNotFound reports that a name or index does not identify a member.
	ErrEnumNotFound    = errors.New("enum not found")
	// ErrInvalidEnumType reports that a registered definition has the wrong type.
	ErrInvalidEnumType = errors.New("enum is not expected type")
)

// Enum is a package sealed interface to identify the enum type
type Enum[T any] interface {
	enum()

	fmt.Stringer

	Name() string
	Index() int
	Valid() bool
	Namespace() Namespace[T]
}

// As associates a name with a comparable enum value for DefineType.
type As[T comparable] struct {
	// Name is the name used to look up the value.
	Name  string
	// Value is the comparable value associated with Name.
	Value T
}
