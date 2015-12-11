package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: program5 <date>")
		os.Exit(1)
	}
	y, m, d, err := parse(os.Args[1])
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Successfully parsed")
	}
	fmt.Println("Year =", y)
	fmt.Println("Month =", m)
	fmt.Println("Day =", d)
}

func parse(date string) (year, month, day int, err error) {
	parts := strings.Split(date, "-")
	if len(parts) != 3 {
		err = fmt.Errorf("Wrong input %s", date)
		return
	}
	year, err = strconv.Atoi(parts[0])
	if err != nil {
		return
	}
	month, err = strconv.Atoi(parts[1])
	if err != nil {
		return
	}
	day, err = strconv.Atoi(parts[2])
	if err == nil {
		fmt.Println("Cool.")
	} else {
		fmt.Println("Not cool.")
	}
	return
}
