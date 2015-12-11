package main

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
)

func main0() {

	// THIS IS JUST EXPERIMENTATION

	if len(os.Args) < 2 {
		info("Usage: ast file.go")
		os.Exit(1)
	}
	in := os.Args[1]
	info("Processing", in)

	// Open program, get it AST
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, in, nil, parser.AllErrors)
	if err != nil {
		panic(err)
	}

	// Apply some changes
	r := &remover{fset, []int{7, 6, 10}, 0}
	ast.Walk(r, f)
	info("remove made", r.visits, "visits")
	//ast.Walk(alterer{fset}, f)

	// Print altered program
	format.Node(os.Stdout, fset, f)
}

type alterer struct {
	fset *token.FileSet
}

func (a *alterer) Visit(node ast.Node) (w ast.Visitor) {

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
		_, ok1 := node.(*ast.ExprStmt)
		_, ok2 := node.(*ast.BasicLit)
		if ok1 || ok2 {
			fmt.Fprintf(os.Stderr, "Found <<<%T: %v>>> at line %v \n", node, node, f.Line(p))
		}
		//		}
	}

	// Return self, for recursive calls on children
	return a
}

// remover deletes nodes starting at specific lines
type remover struct {
	fset        *token.FileSet
	lineNumbers []int
	visits      int
}

func (r *remover) Visit(node ast.Node) (w ast.Visitor) {
	r.visits++
	if node != nil {
		if bs, ok := node.(*ast.BlockStmt); ok {
			kept := make([]ast.Stmt, 0, len(bs.List))
			for _, stmt := range bs.List {
				node := stmt
				p := node.Pos()
				f := r.fset.File(p)
				if in(f.Line(p), r.lineNumbers) {
					fmt.Fprintf(os.Stderr, "DESTROYING %T from line %d\n", node, f.Line(p))
				} else {
					kept = append(kept, stmt)
				}
			}
			bs.List = kept
		}
	}
	return r
}
