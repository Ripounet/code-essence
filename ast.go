package main

import (
	"fmt"
	"go/parser"
	"go/token"
)

func main() {

	// THIS IS JUST EXPERIMENTATION

	// Open sample program, get it AST
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "sample/program1/program1.go", nil, parser.AllErrors)
	if err != nil {
		panic(err)
	}
	fmt.Println(f)

	fmt.Println("End.")
}
