package enum

import (
	"iter"
	"reflect"
)

// func DefinitionOf[T any](namespace T) Namespace[T] {
// 	return registry[reflect.TypeFor[T]()]
// }

// ByName returns the enum member by name.
// Member zero value with false is returned if member name does not return initialized enum member
func ByName[T Namespace[T]](namespace T, name string) (Member, bool) {
	def, ok := registry[reflect.TypeFor[T]()]

	if !ok {
		return Member{}, false
	}

	return def.ByName(name)
}

// ByIndex returns the enum member by the index os its position in the enum list
func ByIndex[T Namespace[T]](namespace T, index int) (Member, bool) {
	def, ok := registry[reflect.TypeFor[T]()]

	if !ok {
		return Member{}, false
	}

	return def.ByIndex(index)
}

// Values returns a defensive copy of the definition's slice of Values
//
// an empty Member slice is returned for any internal error
func Values[T Namespace[T]](namespace T) []Member {
	def, ok := registry[reflect.TypeFor[T]()]

	if !ok {
		return []Member{}
	}

	return def.Values()
}

// Values returns a defensive copy of the definition's slice of Names
//
// an empty string slice is returned for any internal error
func Names[T Namespace[T]](namespace T) []string {
	def, ok := registry[reflect.TypeFor[T]()]

	if !ok {
		return []string{}
	}

	return def.Names()
}

// All provides allocation-free iteration over all enum members.
// Yields Member
//
// No yield for any internal error
func All[T Namespace[T]](namespace T) iter.Seq[Member] {
	def, ok := registry[reflect.TypeFor[T]()]

	if !ok {
		return func(yield func(Member) bool) {}
	}

	return def.All()
}

// Entries provides allocation-free iteration over all enum members.
// Yields Member name and associated Member
//
// No yield for any internal error
func Entries[T Namespace[T]](namespace T) iter.Seq2[string, Member] {
	def, ok := registry[reflect.TypeFor[T]()]

	if !ok {
		return func(yield func(string, Member) bool) {}
	}

	return def.Entries()
}
