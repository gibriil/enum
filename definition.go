package enum

import (
	"iter"
	"reflect"
)

type definition struct {
	identity reflect.Type
	name     string
	length   int
	values   []Member
	names    []string
	lookup   map[string]int
	metadata []metadata
}

type metadata struct {
	Name  string
	Field reflect.StructField
	Type  reflect.Type
}

func (def definition) Len() int {
	return def.length
}

func (def definition) ByName(name string) (Member, bool) {
	index, ok := def.lookup[name]

	if !ok {
		return Member{}, false
	}

	return def.values[index], true
}

func (def definition) ByIndex(index int) (Member, bool) {
	return def.values[index], true
}

// Returns a defensive copy of the slice of Values
func (def definition) Values() []Member {
	out := make([]Member, def.length)
	copy(out, def.values)
	return out
}

// Returns a defensive copy of the slice of Names
func (def definition) Names() []string {
	out := make([]string, def.length)
	copy(out, def.names)
	return out
}

// Allocation-free iteration
func (def definition) All() iter.Seq[Member] {
	return func(yield func(Member) bool) {
		for _, value := range def.values {
			if !yield(value) {
				return
			}
		}
	}
}

func (def definition) Entries() iter.Seq2[string, Member] {
	return func(yield func(string, Member) bool) {
		for i := 0; i < def.length; i++ {
			if !yield(def.names[i], def.values[i]) {
				return
			}
		}
	}
}
