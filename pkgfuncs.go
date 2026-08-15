// Copyright 2026 Peter James Beard. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package enum

import (
	"fmt"
	"iter"
	"reflect"
)

// Define registers the struct enum namespace and uses reflection over
// the struct fields to initialize each enum member
func Define[T any](schema T) T {

	class := reflect.TypeFor[T]()

	if class.Kind() != reflect.Struct {
		panic("enum.Define requires a struct")
	}

	registry.RLock()
	if _, exists := registry.data[class]; exists {
		panic(fmt.Sprintf("enum has already been defined for %s", class))
	}
	registry.RUnlock()

	def := definition{
		identity: class,
		name:     class.Name(),
		values:   []Enum{},
		names:    []string{},
		lookup:   map[string]int{},
		metadata: []metadata{},
	}

	memberIndex := 0

	for i := 0; i < class.NumField(); i++ {
		field := class.Field(i)

		if !field.Type.Implements(reflect.TypeFor[Enum]()) {
			continue
		}

		if def.memberType == nil {
			def.memberType = field.Type
		} else if def.memberType != field.Type {
			panic(fmt.Sprintf("the %s Namespace must contain only one enum type: set as %s, tried to add %s", def.identity.Name(), def.memberType.Name(), field.Type.Name()))
		}

		member := reflect.ValueOf(&schema).Elem().Field(i)

		embedded := member.Addr().Interface().(initializer)
		embedded.initialize(&def, memberIndex)

		def.values = append(def.values, member.Interface().(Enum))
		def.names = append(def.names, field.Name)

		def.lookup[field.Name] = memberIndex

		def.metadata = append(def.metadata, metadata{
			Name:  field.Name,
			Field: field,
			Type:  field.Type,
		})

		memberIndex++
	}

	def.length = memberIndex

	registry.Lock()
	registry.data[class] = &def
	registry.Unlock()

	return schema
}

// Register registers a comparable type into a Namespace and initializes each enum member
func Register[T comparable](entries ...Entry[T]) Namespace {

	class := reflect.TypeFor[T]()

	registry.RLock()
	if _, exists := registry.data[class]; exists {
		panic(fmt.Sprintf("enum has already been defined for %s", class))
	}
	registry.RUnlock()

	def := definition{
		memberType: class,
		name:       class.Name(),
		values:     make([]Enum, len(entries)),
		names:      make([]string, len(entries)),
		lookup:     make(map[string]int, len(entries)),
		metadata:   make([]metadata, len(entries)),
	}

	for memberIndex, entry := range entries {
		member := Member{
			def:   &def,
			index: memberIndex,
		}

		def.values = append(def.values, member)
		def.names = append(def.names, entry.Name)

		def.lookup[entry.Name] = memberIndex

		def.metadata = append(def.metadata, metadata{
			Name: entry.Name,
			Type: class,
		})
	}

	registry.Lock()
	registry.data[class] = &def
	registry.Unlock()

	return Namespace{
		definition: &def,
	}
}

func clearRegisteredNamespace[T any]() {
	registry.Lock()
	defer registry.Unlock()
	delete(registry.data, reflect.TypeFor[T]())
}

// Equal evaluates to Enums and returns if their identity is equal
//
// This is needed for when Enums contain non-comparable members
func Equal(a, b Enum) bool {
	if !a.Valid() || !b.Valid() {
		return false
	}

	return a.identity() == b.identity()
}

// ByName returns the enum member by name.
// Member zero value with false is returned if member name does not return initialized enum member
func ByName[T any](namespace T, name string) (Enum, bool) {
	registry.RLock()
	defer registry.RUnlock()
	def, ok := registry.data[reflect.TypeFor[T]()]

	if !ok {
		return Member{}, false
	}

	return def.ByName(name)
}

// ByNameAs returns the enum member, cast to enum type, by name.
func ByNameAs[E Enum, T any](namespace T, name string) (*E, error) {

	namespaceType := reflect.TypeFor[T]()

	registry.RLock()
	def, ok := registry.data[namespaceType]
	registry.RUnlock()

	if !ok {
		return nil, fmt.Errorf("could not locate Namespace[%v] in registry", namespaceType.Name())
	}

	enum, ok := def.ByName(name)

	if !ok {
		return nil, ErrEnumNotFound
	}

	castEnum, ok := enum.(E)

	if !ok {
		return nil, fmt.Errorf("failed to cast %s to %s", def.memberType, reflect.TypeFor[E]())
	}

	return &castEnum, nil
}

// ByIndex returns the enum member by the index of its position in the enum list
func ByIndex[T any](namespace T, index int) (Enum, bool) {
	registry.RLock()
	defer registry.RUnlock()
	def, ok := registry.data[reflect.TypeFor[T]()]

	if !ok {
		return Member{}, false
	}

	return def.ByIndex(index)
}

// ByIndexAs returns the enum member, cast to enum type, by the index of its position in the enum list
func ByIndexAs[E Enum, T any](namespace T, index int) (*E, error) {

	namespaceType := reflect.TypeFor[T]()

	registry.RLock()
	def, ok := registry.data[namespaceType]
	registry.RUnlock()

	if !ok {
		return nil, fmt.Errorf("could not locate Namespace[%v] in registry", namespaceType.Name())
	}

	enum, ok := def.ByIndex(index)

	if !ok {
		return nil, ErrEnumNotFound
	}

	castEnum, ok := enum.(E)

	if !ok {
		return nil, fmt.Errorf("failed to cast %s to %s", def.memberType, reflect.TypeFor[E]())
	}

	return &castEnum, nil
}

// Values returns a defensive copy of the definition's slice of Values
//
// an empty Member slice is returned for any internal error
func Values[T any](namespace T) []Enum {
	registry.RLock()
	defer registry.RUnlock()
	def, ok := registry.data[reflect.TypeFor[T]()]

	if !ok {
		return []Enum{}
	}

	return def.Values()
}

// ValuesAs returns a defensive copy of the definition's slice of Values each cast to enum type
func ValuesAs[E Enum, T any](namespace T) ([]E, error) {

	namespaceType := reflect.TypeFor[T]()
	enumType := reflect.TypeFor[E]()

	registry.RLock()
	def, ok := registry.data[namespaceType]
	registry.RUnlock()

	if !ok {
		return nil, fmt.Errorf("could not locate Namespace[%v] in registry", namespaceType.Name())
	}

	if def.memberType != enumType {
		return nil, fmt.Errorf("failed to cast %s to %s", def.memberType, enumType)
	}

	arr := def.Values()

	values := make([]E, len(arr))

	for i, enum := range arr {
		castEnum, _ := enum.(E)

		values[i] = castEnum
	}

	return values, nil
}

// Values returns a defensive copy of the definition's slice of Names
//
// an empty string slice is returned for any internal error
func Names[T any](namespace T) []string {
	registry.RLock()
	defer registry.RUnlock()
	def, ok := registry.data[reflect.TypeFor[T]()]

	if !ok {
		return []string{}
	}

	return def.Names()
}

// All provides iteration over all enum members.
// Yields Member
//
// No yield for any internal error
func All[T any](namespace T) iter.Seq[Enum] {
	registry.RLock()
	defer registry.RUnlock()
	def, ok := registry.data[reflect.TypeFor[T]()]

	if !ok {
		return func(yield func(Enum) bool) {}
	}

	return def.All()
}

// AllAs provides iteration over all enum members.
// Yields pointer for Member cast to enum type
//
// Yield error for any internal error
func AllAs[E Enum, T any](namespace T) iter.Seq2[*E, error] {

	namespaceType := reflect.TypeFor[T]()
	enumType := reflect.TypeFor[E]()

	registry.RLock()
	def, ok := registry.data[namespaceType]
	registry.RUnlock()

	if !ok {
		return func(yield func(*E, error) bool) {
			err := fmt.Errorf("could not locate Namespace[%v] in registry", namespaceType.Name())
			if !yield(nil, err) {
				return
			}
		}
	}

	if def.memberType != enumType {
		return func(yield func(*E, error) bool) {
			err := fmt.Errorf("failed to cast %s to %s", def.memberType, enumType)
			if !yield(nil, err) {
				return
			}
		}
	}

	arr := make([]Enum, def.length)

	copy(arr, def.Values())

	return func(yield func(*E, error) bool) {
		for _, enum := range arr {
			castEnum, ok := enum.(E)

			if !ok {
				if !yield(nil, fmt.Errorf("failed to cast %s to %s", def.memberType, reflect.TypeFor[E]())) {
					return
				}
				continue
			}

			if !yield(&castEnum, nil) {
				return
			}
		}
	}
}

// Entries provides iteration over all enum members.
// Yields Member name and associated Member
//
// No yield for any internal error
func Entries[T any](namespace T) iter.Seq2[string, Enum] {
	registry.RLock()
	defer registry.RUnlock()
	def, ok := registry.data[reflect.TypeFor[T]()]

	if !ok {
		return func(yield func(string, Enum) bool) {}
	}

	return def.Entries()
}

func Decode[T Enum](namespace Namespace, enum *T, src any) error {
	if namespace.EnumType() != reflect.TypeFor[T]() {
		return ErrInvalidEnumType
	}

	if src == nil {
		return nil
	}

	var e Enum

	err := namespace.Scan(&e, src)

	if err != nil {
		return err
	}

	*enum = e.(T)
	return nil
}
