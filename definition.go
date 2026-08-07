// Copyright 2026 Peter James Beard. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package enum

import (
	"fmt"
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

func clearRegisteredNamespace[T any]() {
	registry.Lock()
	defer registry.Unlock()
	delete(registry.data, reflect.TypeFor[T]())
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

// Define registers the struct enum namespace and uses reflection over
// the struct fields to initialize each enum member
func Define[T any](schema T) T {

	class := reflect.TypeFor[T]()

	if class.Kind() != reflect.Struct {
		panic("enum.Define requires a struct")
	}

	if _, exists := registry.data[class]; exists {
		panic(fmt.Sprintf("enum has already been defined for %s", class))
	}

	def := definition{
		identity: class,
		name:     class.Name(),
		values:   []Enum{},
		names:    []string{},
		lookup:   map[string]int{},
		metadata: []metadata{},
	}

	for i := 0; i < class.NumField(); i++ {
		field := class.Field(i)

		if def.memberType != field.Type {
			def.memberType = field.Type
		}

		if !field.Type.Implements(reflect.TypeFor[Enum]()) {
			continue
		}

		member := reflect.ValueOf(&schema).Elem().FieldByIndex(field.Index)

		embedded := member.Addr().Interface().(initializer)
		embedded.initialize(&def, i)

		def.values = append(def.values, member.Interface().(Enum))
		def.names = append(def.names, field.Name)

		def.lookup[field.Name] = i

		def.metadata = append(def.metadata, metadata{
			Name:  field.Name,
			Field: field,
			Type:  field.Type,
		})
	}

	def.length = len(def.values)

	registry.Lock()
	registry.data[class] = &def
	registry.Unlock()

	return schema
}

// Len returns the number of initialized enum members
func (def *definition) Len() int {
	return def.length
}

// ByName returns the enum member by name.
// Member zero value with false is returned if member name does not return initialized enum member
func (def *definition) ByName(name string) (Enum, bool) {
	index, ok := def.lookup[name]

	if !ok {
		return Member{}, false
	}

	return def.values[index], true
}

// ByIndex returns the enum member by the index of its position in the enum list
func (def *definition) ByIndex(index int) (Enum, bool) {
	if index < 0 || index >= def.Len() {
		return Member{}, false
	}
	return def.values[index], true
}

// Values returns a defensive copy of the definition's slice of Values
func (def *definition) Values() []Enum {
	out := make([]Enum, def.length)
	copy(out, def.values)
	return out
}

// Values returns a defensive copy of the definition's slice of Names
func (def *definition) Names() []string {
	out := make([]string, def.length)
	copy(out, def.names)
	return out
}

// All provides allocation-free iteration over all enum members.
// Yields Member
func (def *definition) All() iter.Seq[Enum] {
	return func(yield func(Enum) bool) {
		for _, value := range def.values {
			if !yield(value) {
				return
			}
		}
	}
}

// Entries provides allocation-free iteration over all enum members.
// Yields Member name and associated Member
func (def *definition) Entries() iter.Seq2[string, Enum] {
	return func(yield func(string, Enum) bool) {
		for i := 0; i < def.length; i++ {
			if !yield(def.names[i], def.values[i]) {
				return
			}
		}
	}
}

func (def *definition) EnumType() reflect.Type {
	return def.memberType
}

func (def *definition) Type() reflect.Type {
	return def.identity
}
