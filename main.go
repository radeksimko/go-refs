package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"strings"
)

func main() {
	fset := token.NewFileSet()

	var packagePath string

	flag.StringVar(&packagePath, "pkg", "", "full path to package")
	flag.Parse()

	if packagePath == "" {
		log.Fatal("Please provide package path via '-pkg'")
	}

	if len(flag.Args()) != 1 {
		log.Fatal("Please provider path to a file")
	}

	filePath := flag.Arg(0)

	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	src := string(b)

	f, err := parser.ParseFile(fset, filePath, src, 0)
	if err != nil {
		log.Fatal(err)
	}

	refs, err := findPkgRefs(f, packagePath)
	if err != nil {
		log.Fatal(err)
	}

	for _, ref := range refs {
		fmt.Println(ref)
	}
}

func findPkgRefs(f *ast.File, importPath string) ([]ast.Node, error) {
	importName, err := findImportName(f.Imports, importPath)
	if err != nil {
		return nil, err
	}

	var refs = make([]ast.Node, 0)

	ast.Inspect(f, func(node ast.Node) bool {
		if node == nil {
			return false
		}

		switch node.(type) {
		case *ast.File:
			return true
		case *ast.ImportSpec:
			// Imports are processed separately (above)
			return false
		case *ast.SelectorExpr:
			selector := node.(*ast.SelectorExpr)
			x := selector.X

			ident, ok := x.(*ast.Ident)
			if !ok {
				return true
			}

			if ident.Name != importName {
				return true
			}

			refs = append(refs, selector.Sel)
		}

		return true
	})

	return refs, nil
}

// This implementation does not handle case where package
// is imported under different name without explicit alias
// e.g. github.com/hashicorp/go-discover (imported as "discover")
// This would require downloading the module and adding more complexity.
func findImportName(imports []*ast.ImportSpec, importPath string) (string, error) {
	if len(importPath) == 0 {
		return "", fmt.Errorf("Unknown import path")
	}

	parts := strings.Split(importPath, "/")
	importName := parts[len(parts)-1]
	for _, imp := range imports {
		path := strings.Trim(imp.Path.Value, "\"")
		if path == importPath {
			if imp.Name != nil {
				importName = imp.Name.Name
			}

			return importName, nil
		}
	}

	return "", fmt.Errorf("Import %q not found", importPath)
}
