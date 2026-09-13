// Copyright 2026 Peter James Beard. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package enum

import (
	"github.com/gibriil/enum/internal"
)

// Registry is a key:value store for reflection caching
var registry = internal.NewRegistry()

// Namespace describes the registered members of an enum definition.
type Namespace[T any] struct {
	// Definition contains the registered members and lookup indexes.
	*internal.Definition[T]
}

func (def Namespace[T]) identity() *internal.Definition[T] {
	return def.Definition
}

// UnmarshalText looks up text in this namespace and stores the matching enum
// value in enum. It is intended to be called by a concrete enum type's
// UnmarshalText method, because decoding is namespace-specific.
func (def *Namespace[T]) UnmarshalText(enum *T, text []byte) error {
	if def == nil {
		return nil
	}

	e, ok := def.ByName(string(text))
	if !ok {
		return ErrEnumNotFound
	}

	*enum = e
	return nil
}

// Scan decodes a database-style source into enum using this namespace. Strings
// and byte slices are interpreted as member names; ints are interpreted as
// member indexes. A nil source resets enum to its zero value.
func (def *Namespace[T]) Scan(enum *T, src any) error {
	if enum == nil {
		return ErrUninitialized
	}

	if def == nil {
		*enum = *new(T)
		return ErrNotDefined
	}

	if src == nil {
		*enum = *new(T)
		return nil
	}

	switch data := src.(type) {
	case []byte:
		e, ok := def.ByName(string(data))
		if !ok {
			return ErrEnumNotFound
		}
		*enum = e
		return nil
	case string:
		e, ok := def.ByName(data)
		if !ok {
			return ErrEnumNotFound
		}

		*enum = e
		return nil
	case int:
		e, ok := def.ByIndex(data)
		if !ok {
			return ErrEnumNotFound
		}

		*enum = e
		return nil
	default:
		return ErrEnumNotFound
	}
}
