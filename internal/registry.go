// Copyright 2026 Peter James Beard. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package internal

import (
	"errors"
	"reflect"
	"sync"
)

var ErrNoRegistry = errors.New("registry is not initialized")

type Registry struct {
	sync.RWMutex
	data map[reflect.Type]any
}

func (r *Registry) Lookup(key reflect.Type) any {
	if r == nil {
		panic(ErrNoRegistry)
	}
	r.RLock()
	defer r.RUnlock()
	return r.data[key]
}

func (r *Registry) DeleteDefinition(key reflect.Type) {
	if r == nil {
		panic(ErrNoRegistry)
	}
	r.Lock()
	defer r.Unlock()
	delete(r.data, key)
}
