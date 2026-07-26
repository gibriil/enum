package enum

import (
	"database/sql/driver"
	"errors"
)

type Member struct {
	def   *definition
	index int
}

type initializer interface {
	initialize(*definition, int)
}

func (e *Member) initialize(def *definition, index int) {
	e.def = def
	e.index = index
}

func (e Member) Name() string {
	return e.def.names[e.index]
}

func (e Member) String() string {
	return e.def.names[e.index]
}

func (e Member) Index() int {
	return e.index
}

func (e Member) IsZero() bool {
	return e.def == nil
}

func (e Member) Valid() bool {
	return !e.IsZero()
}

func (e Member) MarshalText() ([]byte, error) {
	if e.IsZero() {
		return []byte("nil"), errors.New("Enum is Zero Value")
	}
	return []byte(e.def.names[e.index]), nil
}

func (e *Member) UnmarshalText(data []byte) error {

	def, ok := registry[e.def.identity]

	if !ok {
		return errors.New("Enum not found")
	}

	member, ok := def.ByName(string(data))

	if !ok {
		return errors.New("Enum not found")
	}

	*e = member
	return nil
}

func (e *Member) Scan(src any) error {
	switch data := src.(type) {
	case []byte:
		member, ok := e.def.ByName(string(data))

		if !ok {
			return nil
		}

		*e = member
	case string:
		member, ok := e.def.ByName(data)

		if ok {
			return nil
		}

		*e = member
	case int:
		member, ok := e.def.ByIndex(data)

		if !ok {
			return nil
		}

		*e = member
	default:
	}
	return errors.New("Could not identify enum value in scan")
}

func (e Member) Value() (driver.Value, error) {
	if e.def == nil {
		return nil, errors.New("Enum not found")
	}
	return e.def.names[e.index], nil
}

// enum marks Value as a valid enum implementation.
// It intentionally has no behavior; it seals the Enum interface.
func (e Member) enum() {}
