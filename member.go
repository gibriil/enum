package enum

import "errors"

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

func (e Value) Name() string {
	return e.def.names[e.index]
}

func (e Value) String() string {
	return e.def.names[e.index]
}

func (e Value) Index() int {
	return e.index
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

func (v Value) definition() *definition {
	return v.def
}
