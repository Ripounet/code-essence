package main

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"os"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) < 3 {
		info("Usage: coveragefilter file.go profile.out")
		os.Exit(1)
	}
	in := os.Args[1]
	info("Processing", in)

	profile, err := readProfile(os.Args[2])
	if err != nil {
		panic(err)
	}

	err = coverageFilter(in, os.Stdout, profile)
	if err != nil {
		panic(err)
	}
}

type profileItem struct {
	Line, Column int
	Count        int
}

func readProfile(profileName string) ([]profileItem, error) {
	profile := []profileItem{}

	inFile, _ := os.Open(profileName)
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)

	// Skip "mode: count"
	scanner.Scan()

	for scanner.Scan() {
		item := readProfileItem(scanner.Text())
		profile = append(profile, item)
	}

	return profile, nil
}

func readProfileItem(line string) profileItem {
	info(line)
	parts := strings.Split(line, " ")
	countStr := parts[len(parts)-1]
	count, _ := strconv.Atoi(countStr)
	line = parts[0]

	i := strings.Index(line, ":")
	line = line[i+1:]
	j := strings.Index(line, ",")
	line = line[:j]
	parts = strings.Split(line, ".")
	lineNum, _ := strconv.Atoi(parts[0])
	colNum, _ := strconv.Atoi(parts[1])

	item := profileItem{
		Line:   lineNum, // yuk :(
		Column: colNum,  // yuk :()
		Count:  count,
	}
	info("->", item)
	return item
}

func coverageFilter(in string, out io.Writer, profile []profileItem) error {
	info("Filtering", len(profile), "items.")
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, in, nil, parser.AllErrors)
	if err != nil {
		return err
	}

	killer := &filter{fset, profile}
	ast.Walk(killer, file)

	// Print altered program
	format.Node(out, fset, file)

	return nil
}

type filter struct {
	fset    *token.FileSet
	profile []profileItem
}

func (f *filter) Visit(node ast.Node) (w ast.Visitor) {
	if node != nil {
		position := func(p token.Pos) token.Position {
			file := f.fset.File(p)
			return file.Position(p)
		}
		nodePosition := func(node ast.Node) token.Position {
			return position(node.Pos())
		}
		unused := func(child ast.Node) bool {
			childPosition := nodePosition(child)
			return isUnused(childPosition, f.profile)
		}

		// Filter strategy:
		// FOR EACH node having direct children,
		// DO remove unused children.
		switch parent := node.(type) {
		case *ast.BlockStmt:
			kept := make([]ast.Stmt, 0, len(parent.List))
			for _, child := range parent.List {
				if unused(child) {
					fmt.Fprintf(os.Stderr, "DESTROYING %T at %v\n", child, nodePosition(child))
				} else {
					kept = append(kept, child)
				}
			}
			parent.List = kept
		case *ast.IfStmt:
			if child := parent.Body; unused(child) {
				fmt.Fprintf(os.Stderr, "DESTROYING %T at %v\n", child, nodePosition(child))
				// TODO
			}
			if child := parent.Else; isUnused(position(parent.Body.End()), f.profile) {
				fmt.Fprintf(os.Stderr, "DESTROYING %T at %v\n", child, nodePosition(child))
				parent.Else = nil
			}
		}
	}
	// For recursion
	return f
}

func isUnused(nodePosition token.Position, profile []profileItem) bool {
	info("Testing position", nodePosition)
	for _, item := range profile {
		if nodePosition.Line == item.Line && nodePosition.Column == item.Column && item.Count == 0 {
			return true
		}
	}
	return false
}
