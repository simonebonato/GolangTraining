package main

import "fmt"

type Simo struct {
	name string
	age int
}

func main() {
	simo := Simo{"Simone", 25}
	fmt.Println(simo)
}