package main

import "fmt"

func abs(x int)int {
	if x < 0 {return -x}
	return x
}

func func_with_defer() {
	// this will be executed after the function returns or exits
	defer fmt.Println("The function is over!")

	fmt.Println("Hello")
}

func main()  {
	fmt.Println(abs(-10))
	fmt.Println(abs(10))

	func_with_defer()
}