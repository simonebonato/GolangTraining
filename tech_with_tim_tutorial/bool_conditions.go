package main

import "fmt"

func main() {
	fmt.Println("Hello")

	var num1 float32 = 9.5
	num2 := 10.0

	condition := float32(num2) > num1 && num1 == 9
	fmt.Println(condition)
}