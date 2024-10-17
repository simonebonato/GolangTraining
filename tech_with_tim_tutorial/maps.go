// used to store key-value pairs (Python dictionaries)
package main

import "fmt"

func main() {
	var mappa map[string]int = map[string]int{
		"apple":1,
		"pear":2,
	}

	fmt.Println(mappa)

	// most common way to make an empty map
	mp := make(map[string]int)

	mp["the_best"] = 42
	mp["the_worst"] = 40

	fmt.Println(mp)

	// to delete a value from a map
	delete(mp, "the_best")
	
	fmt.Println(mp)

	// to try to access an element from a map and check if it exists
	val, ok := mp["simo"]
	fmt.Println(val, ok)

	val2, ok2 := mp["the_worst"]
	fmt.Println(val2, ok2)

}