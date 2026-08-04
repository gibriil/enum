package enum

import "fmt"

type Fruit struct {
	Member
}

var fruits = struct {
	Apple   Fruit
	Banana  Fruit
	Coconut Fruit
	Grapes  Fruit
}{}

func Example() {

	fruits = Define(fruits)

	fmt.Printf("I like %ss, %ss, and %s", fruits.Banana, fruits.Coconut, fruits.Grapes)

}
