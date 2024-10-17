package main

import "fmt"

func main() {
	// arrays might be less powerful than slices 
	// same definition as a slice but the size MUST be defined
	var arr = [5]int{1,2,3,4,5}

	fmt.Println(arr)
	fmt.Printf("Lenght: %d, capacity %d\n", len(arr), cap(arr))

	// with slices is a bit easier because we do not need to know 
	// how big the slice is going to be from the beginning
	var s []int = arr[1:4]
	fmt.Println(s)
	fmt.Printf("Lenght: %d, capacity %d\n", len(s), cap(s))

	// now it is also possible to extend the array until it's full capacity
	s = s[:cap(s)]
	fmt.Println(s)
	fmt.Printf("Lenght: %d, capacity %d\n", len(s), cap(s))

	// we can append values to the slice
	s = append(s, 10)
	fmt.Println(s)
	fmt.Printf("Lenght: %d, capacity %d\n", len(s), cap(s))

	s = s[:cap(s)]
	fmt.Println(s)
	fmt.Printf("Lenght: %d, capacity %d\n", len(s), cap(s))
	
}