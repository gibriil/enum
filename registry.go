package enum

import "reflect"

// Registry is a key:value store for reflection caching
var registry = map[reflect.Type]*definition{}

// Define registers the struct enum namespace and uses reflection over
// the struct fields to initialize each enum member
func Define[T any](schema T) T {

	class := reflect.TypeOf(schema)

	if class.Kind() != reflect.Struct {
		panic("enum.Define requires a struct")
	}

	def := definition{
		identity: class,
		name:     class.Name(),
		length:   class.NumField(),
		values:   make([]Member, class.NumField()),
		names:    make([]string, class.NumField()),
		lookup:   make(map[string]int),
		metadata: make([]metadata, class.NumField()),
	}

	registry[reflect.TypeFor[T]()] = &def

	index := 0

	enum := reflect.ValueOf(&schema).Elem()

	for field, data := range enum.Fields() {
		member := data.FieldByName("Member")

		embedded := member.Addr().Interface().(initializer)
		embedded.initialize(&def, index)

		def.values[index] = member.Interface().(Member)
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

func definitionOf(member Member) definition {
	return *member.def
}
