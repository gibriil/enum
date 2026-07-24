package enum

import (
	"iter"
	"reflect"
)

type definition struct {
	identity reflect.Type
	name     string
	length   int
	values   []Value
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

func (def definition) ByName(name string) (Value, bool) {
	index, ok := def.lookup[name]

	if !ok {
		return Value{}, false
	}

	return def.values[index], true
}

func (def definition) ByIndex(index int) (Value, bool) {
	return def.values[index], true
}

// Returns a defensive copy of the slice of Values
func (def definition) Values() []Value {
	out := make([]Value, def.length)
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
func (def definition) All() iter.Seq[Value] {
	return func(yield func(Value) bool) {
		for _, value := range def.values {
			if !yield(value) {
				return
			}
		}
	}
}

func (def definition) Entries() iter.Seq2[string, Value] {
	return func(yield func(string, Value) bool) {
		for i := 0; i < def.length; i++ {
			if !yield(def.names[i], def.values[i]) {
				return
			}
		}
	}
}
