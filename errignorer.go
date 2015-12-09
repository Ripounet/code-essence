package main

import (
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		info("Usage: errignorer file.go")
		os.Exit(1)
	}
	in := os.Args[1]
	info("Processing", in)

	err := errorIgnoreFilter(in, os.Stdout)
	if err != nil {
		panic(err)
	}
}

func errorIgnoreFilter(in string, out io.Writer) error {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, in, nil, parser.AllErrors)
	if err != nil {
		return err
	}

	killer := &errignorer{fset}
	ast.Walk(killer, file)

	// Print altered program
	format.Node(out, fset, file)

	return nil
}

type errignorer struct {
	fset *token.FileSet
}

func (ei *errignorer) Visit(node ast.Node) (w ast.Visitor) {
	if node != nil {
		// TODO
	}
	// For recursion
	return ei
}
