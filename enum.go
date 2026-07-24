package enum

import (
	"reflect"
)

var registry = map[reflect.Type]*definition{}

type Enum interface {
	definition() *definition
	Index() int
	Name() string
	String() string
}

func Define[T any](schema T) T {

	class := reflect.TypeOf(schema)

	if class.Kind() != reflect.Struct {
		panic("enum.Define requires a struct")
	}

	def := definition{
		identity: class,
		name:     class.Name(),
		length:   class.NumField(),
		values:   make([]Value, class.NumField()),
		names:    make([]string, class.NumField()),
		lookup:   make(map[string]int),
		metadata: make([]metadata, class.NumField()),
	}

	registry[reflect.TypeFor[T]()] = &def

	index := 0

	enum := reflect.ValueOf(&schema).Elem()

	for field, data := range enum.Fields() {
		member := data.FieldByName("Value")

		embedded := member.Addr().Interface().(initializer)
		embedded.initialize(&def, index)

		def.values[index] = member.Interface().(Value)
		def.names[index] = field.Name

		def.lookup[field.Name] = index

		def.metadata[index] = metadata{
			Name:  field.Name,
			Field: field,
			Type:  field.Type,
		}

		index++
	}

	return schema
}

func DefinitionOf(member Enum) definition {
	return *member.definition()
}
