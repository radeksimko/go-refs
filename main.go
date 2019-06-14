package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/radeksimko/go-refs/parser"
)

func main() {
	var packagePath string
	var printFile bool

	flag.StringVar(&packagePath, "pkg", "", "full path to package")
	flag.BoolVar(&printFile, "printfile", false, "whether to print location alongside package")
	flag.Parse()

	if packagePath == "" {
		log.Fatal("Please provide package path via '-pkg'")
	}

	if len(flag.Args()) != 1 {
		log.Fatal("Please provider path to a file")
	}

	filePath := flag.Arg(0)

	log.Printf("Parsing %s ...", filePath)
	f, err := parser.ParseFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	refs, err := parser.FindPackageReferences(f, packagePath)
	if err != nil {
		log.Fatal(err)
	}

	for _, ref := range refs {
		if printFile {
			fmt.Printf("%s:%s\n", filePath, ref)
		} else {
			fmt.Println(ref)
		}
	}
}
