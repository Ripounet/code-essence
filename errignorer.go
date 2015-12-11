package main

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"os"
	"strings"
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
				//infof("stmt = %T %v", stmt, stmt)
				switch child := stmt.(type) {
				case *ast.IfStmt:
					be, ok := child.Cond.(*ast.BinaryExpr)
					if !ok {
						break
					}
					identL, ok := be.X.(*ast.Ident)
					if !ok {
						break
					}
					// TODO: checking name only is a dirty approximation.
					// Try to check type instead.
					isErr := strings.HasPrefix(identL.Obj.Name, "err")
					isDiff := be.Op == token.NEQ
					identR, ok := be.Y.(*ast.Ident)
					if !ok {
						break
					}
					isNil := identR.Name == "nil"
					isTestErrDiffNil := isErr && isDiff && isNil
					if isTestErrDiffNil {
						fmt.Fprintln(os.Stderr, "Found err check", be)
						keep = false
						// But we want to keep the ELSE block!!
						// If exists, it is the actual main code.
						if child.Else != nil {
							kept = append(kept, child.Else)
						}
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
