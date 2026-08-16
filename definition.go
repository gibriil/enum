// Copyright 2026 Peter James Beard. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package enum

import (
	"iter"
	"reflect"
	"sync"
)

// Registry is a key:value store for reflection caching
var registry = struct {
	sync.RWMutex
	data map[reflect.Type]*definition
}{
	data: map[reflect.Type]*definition{},
}

// Definition holds the namespace information for the enum list.
type definition struct {
	identity   reflect.Type   // Enum Namespace type
	name       string         // Enum Namespace Name
	memberType reflect.Type   // Enum Member Type
	length     int            // Number of struct members
	values     []Enum         // Slice of all initialized enum members
	names      []string       // Slice of names for each initialized enum member
	lookup     map[string]int //Lookup map for identifying the values index of an enum member by name
	metadata   []metadata     // Slice of the reflection details of each enum member
}

// Metadata holds the reflection information for an enum
type metadata struct {
	Name  string
	Field reflect.StructField
	Type  reflect.Type
}

// Len returns the number of initialized enum members
//
// Len returns -1 if there is not a valid definition
func (def *definition) Len() int {
	if def == nil {
		return -1
	}
	return def.length
}

// ByName returns the enum member by name.
// Member zero value with false is returned if member name does not return initialized enum member
func (def *definition) ByName(name string) (Enum, bool) {
	if def == nil {
		return Member{}, false
	}

	index, ok := def.lookup[name]

	if !ok {
		return Member{}, false
	}

	return def.values[index], true
}

// ByIndex returns the enum member by the index of its position in the enum list
func (def *definition) ByIndex(index int) (Enum, bool) {
	if def == nil {
		return Member{}, false
	}

	if index < 0 || index >= def.Len() {
		return Member{}, false
	}
	return def.values[index], true
}

// Values returns a defensive copy of the definition's slice of Values
func (def *definition) Values() []Enum {
	if def == nil {
		return []Enum{}
	}
	out := make([]Enum, def.length)
	copy(out, def.values)
	return out
}

// Values returns a defensive copy of the definition's slice of Names
func (def *definition) Names() []string {
	if def == nil {
		return []string{}
	}
	out := make([]string, def.length)
	copy(out, def.names)
	return out
}

// All provides iteration over all enum members.
// Yields Member
func (def *definition) All() iter.Seq[Enum] {
	if def == nil {
		return func(yield func(Enum) bool) {}
	}

	return func(yield func(Enum) bool) {
		for _, value := range def.values {
			if !yield(value) {
				return
			}
		}
	}
}

// Entries provides iteration over all enum members.
// Yields Member name and associated Member
func (def *definition) Entries() iter.Seq2[string, Enum] {
	if def == nil {
		return func(yield func(string, Enum) bool) {}
	}
	return func(yield func(string, Enum) bool) {
		for i := 0; i < def.length; i++ {
			if !yield(def.names[i], def.values[i]) {
				return
			}
		}
	}
}

// EnumType return the type for the Member set
func (def *definition) EnumType() reflect.Type {
	if def == nil {
		return nil
	}
	return def.memberType
}

// Type returns the Namespace type
func (def *definition) Type() reflect.Type {
	if def == nil {
		return nil
	}
	return def.identity
}

// UnmarshalText is a helper for creating UnmarshalText on Enum type
func (def *definition) UnmarshalText(enum *Enum, text []byte) error {
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
func (def *definition) Scan(enum *Enum, src any) error {
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
