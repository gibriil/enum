// Copyright 2026 Peter James Beard. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package enum

import (
	"fmt"
	"iter"
	"reflect"
	"unicode"
	"unicode/utf8"

	"github.com/gibriil/enum/internal"
)

// DefineNamespace registers a struct-backed enum namespace and initializes
// each field whose type implements Enum[T]. It returns the initialized schema.
func DefineNamespace[T any, S any](schema S) S {

	class := reflect.TypeFor[S]()

	if class.Kind() != reflect.Struct {
		panic("enum.DefineNamespace requires a struct")
	}

	def := internal.NewDefinition[T](class, class.Name())

	entries := []internal.Entry[T]{}
	metadata := []internal.Metadata{}

	entryIndex := 0

	for i := 0; i < class.NumField(); i++ {
		field := class.Field(i)

		if !field.Type.Implements(reflect.TypeFor[Enum[T]]()) {
			continue
		}

		if def.EntryType() != field.Type {
			panic(fmt.Sprintf("the %s Namespace must contain only one enum type: set as %s, tried to add %s", def.Type().Name(), def.EntryType().Name(), field.Type.Name()))
		}

		member := reflect.ValueOf(&schema).Elem().Field(i)

		embedded := member.FieldByName("Entry").Interface().(internal.Entry[T])

		internal.InitializeEntry(&embedded, field.Name, &def, entryIndex)

		member.FieldByName("Entry").Set(reflect.ValueOf(embedded))

		internal.InitializeEntryValue(&embedded, member.Interface().(T))

		member.FieldByName("Entry").Set(reflect.ValueOf(embedded))

		entries = append(entries, embedded)

		metadata = append(metadata, internal.Metadata{
			Name:  field.Name,
			Field: field,
			Type:  field.Type,
			Value: member,
		})

		entryIndex++
	}

	if len(entries) == 0 {
		panic(fmt.Sprintf("attempted to register %v enums with 0 members", class))
	}

	internal.PopulateDefinition(&def, entries...)
	internal.AttachMetadata(&def, metadata...)

	internal.RegisterDefinition(registry, &def)

	return schema
}

// DefineType registers named comparable values and returns their Namespace.
func DefineType[T comparable](members ...As[T]) Namespace[T] {

	class := reflect.TypeFor[T]()

	if len(members) == 0 {
		panic(fmt.Sprintf("attempted to register %v enums with 0 members", class))
	}

	def := internal.NewDefinition[T](class, class.Name())

	entries := []internal.Entry[T]{}
	metadata := []internal.Metadata{}

	fields := []reflect.StructField{}

	for index, entry := range members {
		embedded := reflect.ValueOf(Member[T]{}).FieldByName("Entry").Interface().(internal.Entry[T])

		internal.InitializeEntry(&embedded, entry.Name, &def, index)

		firstLetter, size := utf8.DecodeRuneInString(entry.Name)
		if firstLetter == utf8.RuneError {
			panic(fmt.Sprintf("invalid entry name at position %d", index))
		}

		internal.InitializeEntryValue(&embedded, entry.Value)

		field := reflect.StructField{
			Name: string(unicode.ToUpper(firstLetter)) + entry.Name[size:],
			Type: class,
		}

		fields = append(fields, field)
		entries = append(entries, embedded)
		metadata = append(metadata, internal.Metadata{
			Name:  entry.Name,
			Field: field,
			Type:  field.Type,
			Value: reflect.ValueOf(entry.Value),
		})
	}

	internal.PopulateDefinition(&def, entries...)
	internal.AttachMetadata(&def, metadata...)

	internal.RegisterDefinition(registry, &def)

	return Namespace[T]{
		Definition: &def,
	}
}

// DefinitionFor returns the registered definition for an enum type
//
// Panics if type is not registered
func DefinitionFor[E, S any]() Namespace[E] {
	def := registry.Lookup(reflect.TypeFor[S]())

	if def == nil {
		panic(ErrNotDefined)
	}

	namespace, ok := def.(*internal.Definition[E])

	if !ok {
		panic(ErrInvalidEnumType)
	}

	return Namespace[E]{
		Definition: namespace,
	}
}

// Equal reports whether a and b identify the same registered enum member.
// It returns false for uninitialized values, missing members, or members from
// different definitions.
func Equal[T any](a, b Enum[T]) bool {
	if !a.Valid() || !b.Valid() {
		return false
	}

	aEntry, ok := internal.Lookup(a.Namespace().identity(), a.Name())
	if !ok {
		return false
	}

	bEntry, ok := internal.Lookup(b.Namespace().identity(), b.Name())
	if !ok {
		return false
	}

	return internal.IdentityOf(aEntry) == internal.IdentityOf(bEntry)
}

// ByName returns the enum member by name.
// Member zero value with false is returned if member name does not return initialized enum member
func ByName[T any](namespace Namespace[T], name string) (T, bool) {
	def, ok := registry.Lookup(namespace.Type()).(*internal.Definition[T])

	if !ok {
		return *new(T), false
	}

	return def.ByName(name)
}

// ByIndex returns the enum member by the index of its position in the enum list
func ByIndex[T any](namespace Namespace[T], index int) (T, bool) {
	def, ok := registry.Lookup(namespace.Type()).(*internal.Definition[T])

	if !ok {
		return *new(T), false
	}

	return def.ByIndex(index)
}

// Values returns a defensive copy of the definition's slice of Values
//
// an empty Member slice is returned for any internal error
func Values[T any](namespace Namespace[T]) []T {
	def, ok := registry.Lookup(namespace.Type()).(*internal.Definition[T])

	if !ok {
		return make([]T, 0)
	}

	return def.Values()
}

// Names returns a defensive copy of the definition's slice of Names
//
// an empty string slice is returned for any internal error
func Names[T any](namespace Namespace[T]) []string {
	def, ok := registry.Lookup(namespace.Type()).(*internal.Definition[T])

	if !ok {
		return []string{}
	}

	return def.Names()
}

// All provides iteration over all enum members in definition order.
//
// No yield for any internal error
func All[T any](namespace Namespace[T]) iter.Seq[T] {
	def, ok := registry.Lookup(namespace.Type()).(*internal.Definition[T])

	if !ok {
		return func(yield func(T) bool) {}
	}

	return def.All()
}

// Entries provides iteration over member names and values in definition order.
//
// No yield for any internal error
func Entries[T any](namespace Namespace[T]) iter.Seq2[string, T] {
	def, ok := registry.Lookup(namespace.Type()).(*internal.Definition[T])

	if !ok {
		return func(yield func(string, T) bool) {}
	}

	return def.Entries()
}

// Decode decodes a source value into an enum using namespace. It accepts the
// same source forms as Namespace.Scan.
func Decode[T Enum[T]](namespace Namespace[T], enum *T, src any) error {
	if enum == nil {
		return ErrUninitialized
	}
	if namespace.EntryType() != reflect.TypeFor[T]() {
		return ErrInvalidEnumType
	}

	if src == nil {
		*enum = *new(T)
		return nil
	}

	return namespace.Scan(enum, src)
}

// Of returns the registered member whose underlying value equals enum. It
// returns an uninitialized Member when no matching value is registered.
func Of[T comparable](enum T) Member[T] {
	registration := registry.Lookup(reflect.TypeFor[T]())

	if registration == nil {
		panic(ErrNotDefined)
	}

	def, ok := registration.(*internal.Definition[T])

	if !ok {
		panic(ErrInvalidEnumType)
	}

	for name, v := range def.Entries() {
		if v == enum {
			cnst, ok := internal.Lookup(def, name)
			if ok {
				return Member[T]{
					Entry: cnst,
				}
			}

		}
	}

	return Member[T]{}
}
