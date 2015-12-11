package main

import (
	"fmt"
	"os"
)

func in(needle int, hay []int) bool {
	for _, j := range hay {
		if j == needle {
			return true
		}
	}
	return false
}

func info(args ...interface{}) {
	fmt.Fprintln(os.Stderr, args...)
}

func infof(pattern string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, pattern+"\n", args...)
}
