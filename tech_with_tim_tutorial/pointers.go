package main

import "fmt"

func takeAPointer(pointX *int) {
	*pointX = 1999
}

func main() {
	x := 7
	fmt.Printf("Using & we can print the memory reference of the variable x, %v, called the pointer or reference to x\n", &x)

	// we can also assign the location of x to another variable
	y := &x

	*y = 10
	fmt.Println(x)
	fmt.Println("We can see that the value of x changed :) \n")

	takeAPointer(&x)
	fmt.Println(x)
}