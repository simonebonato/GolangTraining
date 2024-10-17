package main

import (
	"fmt"
	"bufio"
	"os"
	"strconv"
)

func main() {
	fmt.Println("Hellooooooo")

	// how to define variables
	var x int = 10
	var name string = "Simo"

	fmt.Println(x)
	fmt.Println(name)

	// to get user input we have to define a bufio scanner object
	// os.Stdin is to use the command line as input to the scanner
	scanner := bufio.NewScanner(os.Stdin)

	// we will store inside the scanner object what we collected from the CLI
	fmt.Println("Type something please: ")
	scanner.Scan() // what we type in will be interpreted as a string

	input := scanner.Text()
	fmt.Printf("You typed: %q", input)

	fmt.Println("Type the year in which you were born: ")
	scanner.Scan()
	age_year, _ := strconv.ParseInt(scanner.Text(), 10, 64)
	
	fmt.Printf(
		"You were born in the year %d, so at the end of 2025 you should be %d (old!)", age_year, 2025 - age_year)

}