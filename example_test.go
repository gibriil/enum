package enum

import "fmt"

// Fruit is the enum value type used by the package example.
type Fruit struct {
	Member[Fruit]
}

// fruits is the example namespace schema.
var fruits = struct {
	Apple   Fruit
	Banana  Fruit
	Coconut Fruit
	Grapes  Fruit
}{}

// Example demonstrates defining a namespace and formatting its members.
func Example() {

	fruits = DefineNamespace[Fruit](fruits)

	fmt.Printf("I like %ss, %ss, and %s", fruits.Banana, fruits.Coconut, fruits.Grapes)

}
