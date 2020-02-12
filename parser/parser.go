package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"strings"
)

type CallExpr struct {
	*ast.CallExpr
	importName string
}

func (ce *CallExpr) String() string {
	ident, _ := getIdentFromSelector(ce.Fun.(*ast.SelectorExpr), ce.importName)
	if ident == nil {
		return fmt.Sprintf("%s", ce.Fun)
	}
	return fmt.Sprintf("%s", ident)
}

func ParseFile(filePath string) (*ast.File, error) {
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return parseFile(filePath, string(b))
}

func parseFile(filePath, content string) (*ast.File, error) {
	fset := token.NewFileSet()
	return parser.ParseFile(fset, filePath, content, 0)
}

func getIdentFromSelector(selector *ast.SelectorExpr, importName string) (*ast.Ident, bool) {
	x := selector.X

	ident, ok := x.(*ast.Ident)
	if !ok {
		return nil, true
	}

	if ident.Name != importName {
		return nil, true
	}

	if ident.Obj != nil {
		// Avoid reporting references to local variables
		// which may have the same name as package
		return nil, true
	}

	return selector.Sel, true
}

func FindPackageReferences(f *ast.File, importPath string) ([]ast.Node, error) {
	importName, err := FindImportName(f.Imports, importPath)
	if err != nil {
		return nil, err
	}

	i := &inspector{
		importName: importName,
		foundRefs:  make([]ast.Node, 0),
	}

	ast.Inspect(f, func(node ast.Node) bool {
		return i.inspectNode(node)
	})

	return i.foundRefs, nil
}

type inspector struct {
	importName string
	foundRefs  []ast.Node
}

func (i *inspector) inspectNode(node ast.Node) bool {
	if node == nil {
		return false
	}

	switch n := node.(type) {
	case *ast.File:
		return true
	case *ast.ImportSpec:
		// Imports are processed separately (above)
		return false
	case *ast.CallExpr:
		if n.Fun == nil {
			return true
		}

		switch expr := n.Fun.(type) {
		case *ast.SelectorExpr:
			ident, keepTraversing := getIdentFromSelector(expr, i.importName)
			if ident == nil {
				return keepTraversing
			}

			i.foundRefs = append(i.foundRefs, &CallExpr{n, i.importName})
			return false
		default:
			return true
		}
	case *ast.SelectorExpr:
		ident, keepTraversing := getIdentFromSelector(n, i.importName)
		if ident == nil {
			return keepTraversing
		}

		i.foundRefs = append(i.foundRefs, ident)
	default:
	}

	return true
}

// This implementation does not handle case where package
// is imported under different name without explicit alias
// e.g. github.com/hashicorp/go-discover (imported as "discover")
// This would require downloading the module and adding more complexity.
func FindImportName(imports []*ast.ImportSpec, importPath string) (string, error) {
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
