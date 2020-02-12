package parser

import (
	"fmt"
	"go/ast"
	"reflect"
	"testing"
)

func TestFindPackageReferences(t *testing.T) {
	testCases := []struct {
		name         string
		code         string
		pkgName      string
		expectedRefs []string
	}{
		{
			"function",
			`package example

import (
	"fmt"
	"os"
)

func main() {
	fmt.Printf("something")
}
`,
			"fmt",
			[]string{"Printf"},
		},
		{
			"ident",
			`package example

import (
	"fmt"
	"os"
)

func main() {
	fmt.Something
}
`,
			"fmt",
			[]string{"Something"},
		},
		{
			"funcs and idents",
			`package example

import (
	"fmt"
	"os"
)

func main() {
	fmt.Something
	fmt.ExampleFunc()
}
`,
			"fmt",
			[]string{"Something", "ExampleFunc"},
		},
		{
			"matching variable",
			`package example

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
`,
			"fmt",
			[]string{},
		},
		{
			"anonymous function",
			`package terraform

import (
	"fmt"
)

func Run() error {
	defer func() {
		fmt.Println("at last")
	}()
	return nil
}
`,
			"fmt",
			[]string{"Println"},
		},
		{
			"type declaration",
			`package terraform

import (
	"fmt"
)

type MyFunc func(f fmt.Formatter ) error

`,
			"fmt",
			[]string{"Formatter"},
		},
		{
			"variable declaration",
			`package terraform

import (
	"fmt"
)

var formatter fmt.Formatter
var ptr *fmt.Ptr

`,
			"fmt",
			[]string{"Formatter", "Ptr"},
		},
		{
			"function arguments",
			`package terraform

import (
	"fmt"
)

func run(f fmt.CustomFunc) error {
	return f()
}

`,
			"fmt",
			[]string{"CustomFunc"},
		},
		{
			"function results",
			`package terraform

import (
	"fmt"
)

func run() *fmt.Error {
	return nil
}

`,
			"fmt",
			[]string{"Error"},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d-%s", i, tc.name), func(t *testing.T) {
			f, err := parseFile("src.go", tc.code)
			if err != nil {
				t.Fatal(err)
			}

			refs, err := FindPackageReferences(f, tc.pkgName)
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(stringifyAstNodes(refs), tc.expectedRefs) {
				t.Fatalf("Expected: %q\nGiven: %q\n", tc.expectedRefs, refs)
			}
		})
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
