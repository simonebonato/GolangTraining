package main

import "fmt"

func main() {
	// it always starts with the "for" keyword
	// and can also work by itself
	// this for example is a while loop, but it needs the for keyword
	fmt.Println("First style of for loop")
	x := 3
	for x < 5 {
		fmt.Println(x)
		x++
	}


	fmt.Println("\nSecond style of for loop")
	for i:=0 ; i <=5; i++ {
		fmt.Println(i)
	}

	// this is instead a switch statement
	fmt.Println("Switch statement value ")
	value := 10

	switch value {
	case 10:
		fmt.Println(10)
	case 20:
		fmt.Println(20)
		
	}
}