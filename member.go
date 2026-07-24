package enum

import (
	"database/sql/driver"
	"errors"
)

type Value struct {
	def   *definition
	index int
}

type initializer interface {
	initialize(*definition, int)
}

func (v *Value) initialize(def *definition, index int) {
	v.def = def
	v.index = index
}

func (v Value) Name() string {
	return v.def.names[v.index]
}

func (v Value) String() string {
	return v.def.names[v.index]
}

func (v Value) Index() int {
	return v.index
}

func (v Value) MarshalText() ([]byte, error) {
	return []byte(v.def.names[v.index]), nil
}

func (v *Value) UnmarshalText(data []byte) error {

	def, ok := registry[v.def.identity]

	if !ok {
		return errors.New("Enum not found")
	}

	member, ok := def.ByName(string(data))

	if !ok {
		return errors.New("Enum not found")
	}

	*v = member
	return nil
}

func (v *Value) Scan(src any) error {
	switch data := src.(type) {
	case []byte:
		member, ok := v.def.ByName(string(data))

		if !ok {
			return nil
		}

		v = &member
	case string:
		member, ok := v.def.ByName(data)

		if ok {
			return nil
		}

		v = &member
	case int:
		member, ok := v.def.ByIndex(data)

		if !ok {
			return nil
		}

		v = &member
	default:
	}
	return errors.New("Could not identify enum value in scan")
}

func (v Value) Value() (driver.Value, error) {
	if v.def == nil {
		return nil, nil
	}
	return v.def.names[v.index], nil
}

// enum marks Value as a valid enum implementation.
// It intentionally has no behavior; it seals the Enum interface.
func (v Value) enum() {}
