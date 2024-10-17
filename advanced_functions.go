package main

import "fmt"

func test() {
	fmt.Println("Test")
}

func func_that_takes_a_func(x int, myFunc func(int) int) {
	myFunc(x)
}

func main() {
	// a function can be assigned to a variable
	x := test
	x()

	// functions can also be created within like this
	test2 := func(x int) int {
		fmt.Println(x)
		fmt.Println("This is so weird")
		return x * -1
	}

	fmt.Println(test2)

	func_that_takes_a_func(100, test2)
}