package main

import "fmt"

func main() {

	var sl = []int{11,2,3,42,5,6,3,8,9,10}
	
	for index, element := range sl {
		fmt.Printf("\nValue: %d and index %d", element, index)
	}

	for index1, element1 := range sl {
		for index2, element2 := range sl[index1 + 1:] {
			if element1 == element2 {
				fmt.Printf("The element %d is duplicated in the list, in positions %d and %d", element1, index1, index2)
			}
		}
	}
}