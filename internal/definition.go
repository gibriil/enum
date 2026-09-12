// Copyright 2026 Peter James Beard. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package internal

import (
	"iter"
	"reflect"
)

// Definition holds the encapsulation information for the entry list.
type Definition[T any] struct {
	identity  reflect.Type   // Namespace type
	name      string         // Namespace Name
	entryType reflect.Type   // Entry Type
	length    int            // Number of Entries
	values    []Entry[T]     // Slice of all initialized entries
	names     []string       // Slice of names for each initialized Entry
	lookup    map[string]int //Lookup map for identifying the values index of an Entry by name
	metadata  []Metadata     // Slice of the reflection details of each Entry
}

// Len returns the number of initialized entries
//
// Len returns -1 if there is not a valid definition
func (def *Definition[T]) Len() int {
	if def == nil {
		return -1
	}
	return def.length
}

// ByName returns the Entry by name.
// Entry zero value with false is returned if Entry name does not return initialized Entry
func (def *Definition[T]) ByName(name string) (T, bool) {
	if def == nil {
		return *new(T), false
	}

	index, ok := def.lookup[name]

	if !ok {
		return *new(T), false
	}

	return def.values[index].Value(), true
}

// ByIndex returns the Entry by the index of its position in the list
func (def *Definition[T]) ByIndex(index int) (T, bool) {
	if def == nil {
		return *new(T), false
	}

	if index < 0 || index >= def.Len() {
		return *new(T), false
	}
	return def.values[index].Value(), true
}

// Values returns a defensive copy of the definition's slice of Values
func (def *Definition[T]) Values() []T {
	if def == nil {
		return []T{}
	}
	out := make([]T, def.length)
	for i, val := range def.values {
		out[i] = val.Value()
	}
	return out
}

// Names returns a defensive copy of the definition's slice of Names
func (def *Definition[T]) Names() []string {
	if def == nil {
		return []string{}
	}
	out := make([]string, def.length)
	copy(out, def.names)
	return out
}

// All provides iteration over all entries.
// Yields Entry
func (def *Definition[T]) All() iter.Seq[T] {
	if def == nil {
		return func(yield func(T) bool) {}
	}

	return func(yield func(T) bool) {
		for _, entry := range def.values {
			if !yield(entry.Value()) {
				return
			}
		}
	}
}

// Entries provides iteration over all entries.
// Yields Entry name and associated value
func (def *Definition[T]) Entries() iter.Seq2[string, T] {
	if def == nil {
		return func(yield func(string, T) bool) {}
	}
	return func(yield func(string, T) bool) {
		for i := 0; i < def.length; i++ {
			if !yield(def.names[i], def.values[i].Value()) {
				return
			}
		}
	}
}

func (def *Definition[T]) Metadata() []Metadata {
	if def == nil {
		return []Metadata{}
	}

	out := make([]Metadata, def.length)
	copy(out, def.metadata)
	return out
}

// EntryType return the type for the Entry set
func (def *Definition[T]) EntryType() reflect.Type {
	if def == nil {
		return nil
	}
	return def.entryType
}

// Type returns the Namespace type
func (def *Definition[T]) Type() reflect.Type {
	if def == nil {
		return nil
	}
	return def.identity
}

// Name returns the Namespace Name
func (def *Definition[T]) Name() string {
	if def == nil {
		return ""
	}
	return def.name
}
