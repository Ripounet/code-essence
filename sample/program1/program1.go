package main

import "fmt"

func main() {
	if 1 == 2 {
		fmt.Println(1)
	} else {
		fmt.Println(2)
	}

	if "a"+"b" == "ab" {
		fmt.Println(3)
	} else {
		fmt.Println(4)
		fmt.Println(5)
	}
}
