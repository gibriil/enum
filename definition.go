// Copyright 2026 Peter James Beard. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package enum

import (
	"reflect"
	"sync"

	"github.com/gibriil/enum/internal"
)

// Registry is a key:value store for reflection caching
var registry = struct {
	sync.RWMutex
	data map[reflect.Type]any
}{
	data: map[reflect.Type]any{},
}

type Namespace[T any] struct {
	*internal.Definition[T]
}

func (def Namespace[T]) identity() *internal.Definition[T] {
	return def.Definition
}

// UnmarshalText is a helper for creating UnmarshalText on Enum type
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

// Scan is a helper for creating Scan on Enum type
func (def *Namespace[T]) Scan(enum *T, src any) error {
	if def == nil || src == nil {
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
