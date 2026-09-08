// Copyright 2026 Peter James Beard. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package internal

import (
	"errors"
	"reflect"
)

var (
	ErrUninitialized   = errors.New("entry is zero Value")
	ErrNotDefined      = errors.New("entry has no registered definition")
	ErrEnumNotFound    = errors.New("entry not found")
	ErrInvalidEnumType = errors.New("entry is not expected type")
)

type Entry[T any] struct {
	entryIdentity[T]
	value T
}

type entryIdentity[T any] struct {
	name  string
	def   *Definition[T]
	index int
}

// Metadata holds the reflection information for an enum
type Metadata struct {
	Name  string
	Field reflect.StructField
	Type  reflect.Type
	Value reflect.Value
}

// Valid reports whether or not the Entry has been initialized
func (e Entry[T]) Valid() bool {
	return e.def != nil
}

// identity returns the comparable for Entry equality checks
func (e Entry[T]) identity() entryIdentity[T] {
	return e.entryIdentity
}

// Name returns the entry's name
func (e Entry[T]) Name() string {
	if !e.Valid() {
		return ""
	}
	return e.name
}

// String returns the entry's name
func (e Entry[T]) String() string {
	if !e.Valid() {
		return ""
	}
	return e.name
}

// Index returns the index of entry's position in the definition's list
//
// zero value will return -1
func (e Entry[T]) Index() int {
	if !e.Valid() {
		return -1
	}
	return e.index
}

// Value returns the concrete value of the entry's type
func (e Entry[T]) Value() T {
	return e.value
}

// Definition returns the internally registered definition for a collection of entries
func (e Entry[T]) Definition() *Definition[T] {
	if !e.Valid() {
		return &Definition[T]{}
	}
	return e.Definition()
}
