package enum

import (
	"reflect"
)

type collection interface {
	Values() []Value
	Names() []string
	Len() int
	ByName(string) (Value, bool)
	ByIndex(int) (Value, bool)
}

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
	member := def.lookup[name]
	return def.values[member], true
}

func (def definition) ByIndex(index int) (Value, bool) {
	return def.values[index], true
}
