package enum

import (
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

func (def definition) Values() []Value {
	return def.values
}

func (def definition) Names() []string {
	return def.names
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
