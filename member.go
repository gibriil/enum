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

func (e Value) Is(other Value) bool {
	return e == other
}

func (v Value) MarshalText() ([]byte, error) {
	return []byte(v.def.names[v.index]), nil
}

func (v *Value) UnmarshalText(data []byte) error {

	def, ok := registry[v.def.identity]

	if !ok {
		return errors.New("Enum not found")
	}

	val, ok := def.ByName(string(data))

	if !ok {
		return errors.New("Enum not found")
	}

	*v = val
	return nil
}

func (v *Value) Scan(src any) error {
	switch data := src.(type) {
	case []byte:
		if _, ok := v.def.ByName(string(data)); ok {
			return nil
		}
	case string:
		if _, ok := v.def.ByName(data); ok {
			return nil
		}
	case int:
		if _, ok := v.def.ByIndex(data); ok {
			return nil
		}
	default:
	}
	return errors.New("Could not identify enum value in scan")
}
