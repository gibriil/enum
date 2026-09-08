// Copyright 2026 Peter James Beard. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package internal

import (
	"fmt"
	"reflect"
	"slices"
)

func SetEntry[T any](e *Entry[T], name string, value T) {
	e.entryIdentity = entryIdentity[T]{
		name: name,
	}
	e.value = value
}

func InitializeDefinition[T any](id reflect.Type, name string) Definition[T] {
	return Definition[T]{
		identity:  id,
		name:      name,
		entryType: reflect.TypeFor[T](),
	}
}

func RegisterDefinition[T any](def *Definition[T], entries ...Entry[T]) {
	def.length = len(entries)
	def.values = make([]Entry[T], def.length)
	def.names = make([]string, def.length)
	def.lookup = make(map[string]int, def.length)
	def.metadata = make([]Metadata, def.length)

	for _, entry := range entries {
		if slices.Contains(def.names, entry.name) {
			panic(fmt.Sprintf("%v already has entry for %s", def.identity, entry.name))
		}
	}

}

func AttachMetadata[T any](def *Definition[T], data ...Metadata) {
	def.metadata = data
}

func EntryAt[T any](def *Definition[T], index int) Entry[T] {
	return def.values[index]
}

func IdentityOf[T any](entry Entry[T]) entryIdentity[T] {
	return entry.entryIdentity
}
