package enum

import "fmt"

type Fruit struct {
	Member[Fruit]
}

var fruits = struct {
	Apple   Fruit
	Banana  Fruit
	Coconut Fruit
	Grapes  Fruit
}{}

func Example() {

	fruits = DefineNamespace[Fruit](fruits)

	fmt.Printf("I like %ss, %ss, and %s", fruits.Banana, fruits.Coconut, fruits.Grapes)

}
