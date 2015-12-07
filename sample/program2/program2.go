package main

import "fmt"

func main() {
	f(2)
	f(3)
}

func f(i int) {
	switch i {
	case 1:
		fmt.Println("Un")
	case 2:
		fmt.Println("Dos")
	case 3:
		fmt.Println("Tres")
	default:
		panic(i)
	}
}
