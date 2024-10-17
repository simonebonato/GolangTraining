package main

import "fmt"

func main() {
	// some data type are mutable and some others are immutable 
	// slices are mutable, because it can change
	var x []int = []int{1,2,3}

	y := x
	y[0] = 100

	fmt.Println(x, y)
	fmt.Println("\nWe can see that the values of both x and y have changed.")
	fmt.Println("The same works with maps. And it depends on the used data type.")


}