// Copyright 2026 Peter James Beard. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package internal

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

// ErrNoRegistry reports that a registry receiver has not been initialized.
var ErrNoRegistry = errors.New("registry is not initialized")

// Registry stores definitions keyed by reflection type.
type Registry struct {
	sync.RWMutex
	data map[reflect.Type]any
}

// Lookup returns the value registered for key, or nil when key is absent.
func (r *Registry) Lookup(key reflect.Type) any {
	if r == nil {
		panic(ErrNoRegistry)
	}
	r.RLock()
	defer r.RUnlock()
	return r.data[key]
}

// DeleteDefinition removes the definition registered for key.
func (r *Registry) DeleteDefinition(key reflect.Type) {
	if r == nil {
		panic(ErrNoRegistry)
	}
	r.Lock()
	defer r.Unlock()
	delete(r.data, key)
}

// RegisterDefinition registers def under its definition identity.
func RegisterDefinition[T any](r *Registry, def *Definition[T]) {
	r.Lock()
	defer r.Unlock()

	if _, exists := r.data[def.identity]; exists {
		panic(fmt.Sprintf("enum has already been defined for %s", def.identity))
	}

	r.data[def.identity] = def
}
