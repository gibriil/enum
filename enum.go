// Copyright 2026 Peter James Beard. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package enum

import (
	"errors"
	"fmt"
)

var (
	ErrUninitialized   = errors.New("enum is Zero Value")
	ErrNotDefined      = errors.New("enum has no registered definition")
	ErrEnumNotFound    = errors.New("enum not found")
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

type As[T comparable] struct {
	Name  string
	Value T
}
