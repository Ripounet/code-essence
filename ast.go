package main

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
)

func main() {

	// THIS IS JUST EXPERIMENTATION

	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: ast file.go")
		os.Exit(1)
	}
	in := os.Args[1]
	fmt.Fprintln(os.Stderr, "Processing", in)

	// Open program, get it AST
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, in, nil, parser.AllErrors)
	if err != nil {
		panic(err)
	}

	// Apply some changes
	ast.Walk(alterer{fset}, f)

	// Print altered program
	format.Node(os.Stdout, fset, f)
}

type alterer struct {
	fset *token.FileSet
}

func (a alterer) Visit(node ast.Node) (w ast.Visitor) {

	// Replace 4 with 666
	if bl, ok := node.(*ast.BasicLit); ok {
		if bl.Kind == token.INT && bl.Value == "4" {
			bl.Value = "666"
		}
	}

	// Replace if body with else body
	if is, ok := node.(*ast.IfStmt); ok {
		//		is. = is.Else
		if elseBlock, ok := is.Else.(*ast.BlockStmt); ok {
			is.Body.List = elseBlock.List
		} else {
			is.Body.List = []ast.Stmt{is.Else}
		}
		is.Else = nil
	}

	// Use original line
	if node != nil {
		p := node.Pos()
		f := a.fset.File(p)
		//		if f.Line(p) == 13 {
		fmt.Fprintf(os.Stderr, "Found <<<%v>>> at line %v \n", node, f.Line(p))
		//		}
	}

	// Return self, for recursive calls on children
	return a
}
