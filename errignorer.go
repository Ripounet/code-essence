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
	info("")

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
					isLeftErr := strings.HasPrefix(identL.Obj.Name, "err")
					identR, ok := be.Y.(*ast.Ident)
					if !ok {
						break
					}
					isRightNil := identR.Name == "nil"
					if isLeftErr && isRightNil {
						fmt.Fprintln(os.Stderr, "Found err check", be)
						fmt.Fprintf(os.Stderr, "body: %T \n", child.Body)
						fmt.Fprintf(os.Stderr, "else: %T \n", child.Else)
						switch be.Op {
						case token.NEQ:
							keep = false
							// But we want to keep the ELSE block!!
							// If exists, it is the actual main code.
							if child.Else != nil {
								switch strictness {
								case Strict:
									kept = append(kept, child.Else)
								case Loose, Exotic:
									kept = appendNodes(kept, child.Else)
								}
							}
						case token.EQL:
							keep = false
							// But we want to keep the BODY block!!
							// It is the actual main code.
							switch strictness {
							case Strict:
								kept = append(kept, child.Body)
							case Loose, Exotic:
								kept = appendNodes(kept, child.Body)
							}
						}
					}
				case *ast.ExprStmt:
					if expr, ok := child.X.(*ast.CallExpr); ok {
						if selexpr, ok := expr.Fun.(*ast.SelectorExpr); ok {
							if selexpr.Sel.Name == "Exit" {
								if left, ok := selexpr.X.(*ast.Ident); ok && left.Name == "os" {
									infof("Found an Exit call : %v \n", expr.Fun)
								}
							}
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

func appendNodes(list []ast.Stmt, stmt ast.Stmt) []ast.Stmt {
	switch x := stmt.(type) {
	case *ast.BlockStmt:
		return append(list, x.List...)
	default:
		return append(list, x)
	}
}
