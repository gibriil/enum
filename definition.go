// Copyright 2026 Peter James Beard. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package enum

import (
	"iter"
	"reflect"
)

// Definition holds the namespace information for the enum list.
type definition struct {
	identity reflect.Type   // Enum Namespace type
	name     string         // Enum Namespace Name
	length   int            // Number of struct members
	values   []Member       // Slice of all initialized enum members
	names    []string       // Slice of names for each initialized enum member
	lookup   map[string]int //Lookup map for identifying the values index of an enum member by name
	metadata []metadata     // Slice of the reflection details of each enum member
}

// Metadata holds the reflection information for an enum
type metadata struct {
	Name  string
	Field reflect.StructField
	Type  reflect.Type
}

// Len returns the number of initialized enum members
func (def definition) Len() int {
	return def.length
}

// ByName returns the enum member by name.
// Member zero value with false is returned if member name does not return initialized enum member
func (def definition) ByName(name string) (Member, bool) {
	index, ok := def.lookup[name]

	if !ok {
		return Member{}, false
	}

	return def.values[index], true
}

// ByIndex returns the enum member by the index os its position in the enum list
func (def definition) ByIndex(index int) (Member, bool) {
	if index < 0 || index > def.Len() {
		return Member{}, false
	}
	return def.values[index], true
}

// Values returns a defensive copy of the definition's slice of Values
func (def definition) Values() []Member {
	out := make([]Member, def.length)
	copy(out, def.values)
	return out
}

// Values returns a defensive copy of the definition's slice of Names
func (def definition) Names() []string {
	out := make([]string, def.length)
	copy(out, def.names)
	return out
}

// All provides allocation-free iteration over all enum members.
// Yields Member
func (def definition) All() iter.Seq[Member] {
	return func(yield func(Member) bool) {
		for _, value := range def.values {
			if !yield(value) {
				return
			}
		}
	}
}

// Entries provides allocation-free iteration over all enum members.
// Yields Member name and associated Member
func (def definition) Entries() iter.Seq2[string, Member] {
	return func(yield func(string, Member) bool) {
		for i := 0; i < def.length; i++ {
			if !yield(def.names[i], def.values[i]) {
				return
			}
		}
	}
}

// namespace marks Namespace as a valid enum Namespace implementation.
// It intentionally has no behavior; it seals the Namespace interface.
func (e definition) namespace() {}
