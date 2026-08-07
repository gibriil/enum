// Copyright 2026 Peter James Beard. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package enum

import (
	"fmt"
	"iter"
	"reflect"
)

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

	if enumType != enumType {
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

	if enumType != enumType {
		return func(yield func(*E, error) bool) {
			err := fmt.Errorf("failed to cast %s to %s", def.memberType, enumType)
			if !yield(nil, err) {
				return
			}
		}
	}

	arr := def.Values()

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
