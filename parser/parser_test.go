package parser

import (
	"fmt"
	"go/ast"
	"reflect"
	"testing"
)

func TestFindPackageReferences_function(t *testing.T) {
	testContent := `package example

import (
	"fmt"
	"os"
)

func main() {
	fmt.Printf("something")
}
`
	f, err := parseFile("src.go", testContent)
	if err != nil {
		t.Fatal(err)
	}

	refs, err := FindPackageReferences(f, "fmt")
	if err != nil {
		t.Fatal(err)
	}

	expectedRefs := []string{"Printf"}
	if !reflect.DeepEqual(stringifyAstNodes(refs), expectedRefs) {
		t.Fatalf("Expected: %q\nGiven: %q\n", expectedRefs, refs)
	}
}

func TestFindPackageReferences_ident(t *testing.T) {
	testContent := `package example

import (
	"fmt"
	"os"
)

func main() {
	fmt.Something
}
`
	f, err := parseFile("src.go", testContent)
	if err != nil {
		t.Fatal(err)
	}

	refs, err := FindPackageReferences(f, "fmt")
	if err != nil {
		t.Fatal(err)
	}

	expectedRefs := []string{"Something"}
	if !reflect.DeepEqual(stringifyAstNodes(refs), expectedRefs) {
		t.Fatalf("Expected: %q\nGiven: %q\n", expectedRefs, refs)
	}
}

func TestFindPackageReferences_funcsAndIdents(t *testing.T) {
	testContent := `package example

import (
	"fmt"
	"os"
)

func main() {
	fmt.Something
	fmt.ExampleFunc()
}
`
	f, err := parseFile("src.go", testContent)
	if err != nil {
		t.Fatal(err)
	}

	refs, err := FindPackageReferences(f, "fmt")
	if err != nil {
		t.Fatal(err)
	}

	expectedRefs := []string{"Something", "ExampleFunc"}
	if !reflect.DeepEqual(stringifyAstNodes(refs), expectedRefs) {
		t.Fatalf("Expected: %q\nGiven: %q\n", expectedRefs, refs)
	}

	if v, ok := refs[0].(*ast.Ident); !ok {
		t.Fatalf("Expected 1st reference to be of type *ast.Ident, given %T", v)
	}

	if v, ok := refs[1].(*CallExpr); !ok {
		t.Fatalf("Expected 2nd reference to be of type *CallExpr, given %T", v)
	}
}

func TestFindPackageReferences_matchingVariable(t *testing.T) {
	testContent := `package example

import (
	"fmt"
	"os"
)

type Formatter struct {
}

func (f *Formatter) Printf() string {
	return "hello"
}

func main() {
	fmt := &Formatter{}
	fmt.Printf("something")
}
`
	f, err := parseFile("src.go", testContent)
	if err != nil {
		t.Fatal(err)
	}

	refs, err := FindPackageReferences(f, "fmt")
	if err != nil {
		t.Fatal(err)
	}

	expectedRefs := []string{}
	if !reflect.DeepEqual(stringifyAstNodes(refs), expectedRefs) {
		t.Fatalf("Expected: %q\nGiven: %q\n", expectedRefs, refs)
	}
}

func stringifyAstNodes(nodes []ast.Node) []string {
	result := make([]string, len(nodes), len(nodes))
	for i, node := range nodes {
		result[i] = fmt.Sprint(node)
	}
	return result
}

func TestFindImportName(t *testing.T) {
	testContent := `package example

import (
	xyz "fmt"
	"os"
)

func main() {
	xyz.Printf("something")
}
`
	f, err := parseFile("src.go", testContent)
	if err != nil {
		t.Fatal(err)
	}

	aliasName, err := FindImportName(f.Imports, "fmt")
	if err != nil {
		t.Fatal(err)
	}
	expectedAliasName := "xyz"
	if aliasName != expectedAliasName {
		t.Fatalf("Expected: %q, given: %q", expectedAliasName, aliasName)
	}

	osName, err := FindImportName(f.Imports, "os")
	if err != nil {
		t.Fatal(err)
	}
	expecteOsName := "os"
	if osName != expecteOsName {
		t.Fatalf("Expected: %q, given: %q", expecteOsName, osName)
	}
}
