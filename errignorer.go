package main

import (
	"fmt"
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
		switch parent := node.(type) {
		case *ast.BlockStmt:
			kept := make([]ast.Stmt, 0, len(parent.List))
			for _, stmt := range parent.List {
				keep := true
				switch child := stmt.(type) {
				case *ast.IfStmt:
					be, ok := child.Cond.(*ast.BinaryExpr)
					if !ok {
						break
					}
					ident, ok := be.X.(*ast.Ident)
					if !ok {
						break
					}
					// TODO: fix all this mess
					fmt.Printf("ident.Obj.Decl = %T %v \n", ident.Obj.Decl, ident.Obj.Decl)
					isErr := true //be.X.(*ast.Object).Kind == ast.Var
					isDiff := be.Op == token.NEQ
					bl, ok := be.Y.(*ast.BasicLit)
					if !ok {
						break
					}
					isNil := bl.Value == "nil"
					isTestErrDiffNil := isErr && isDiff && isNil
					if isTestErrDiffNil {
						fmt.Fprintln(os.Stderr, "Found err check", be)
					}
				}
				if keep {
					kept = append(kept, stmt)
				} else {
					fmt.Fprintf(os.Stderr, "DESTROYING %T \n", stmt)
				}
			}
			parent.List = kept
		}
	}
	// For recursion
	return ei
}
