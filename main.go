package main

import (
	"flag"
	"log"
)

var (
	verbose bool
)

// to allow for testing
var createGeneratedFileFunc = createGeneratedFile

func run(directory string, recurse bool, file string) {
	if file != "" {
		parsed := parseFile(file)
		if len(parsed.structs) > 0 {
			createGeneratedFileFunc(parsed.dirname, parsed.pkg, parsed.imports, parsed.structs)
		} else {
			log.Fatalf("no structs found in file %s, aborting", file)
		}
	} else {
		var pkgs []genPackage
		if recurse {
			pkgs = parseAllDirs(directory)
		} else {
			pkgs = parseDir(directory)
		}
		for _, pkg := range pkgs {
			if len(pkg.structs) > 0 {
				createGeneratedFileFunc(pkg.dirname, pkg.pkg, pkg.imports, pkg.structs)
			}
		}
	}
}

func main() {
	directory := flag.String("d", "./", "generate factory methods for all structs within the directory")
	recurse := flag.Bool("r", true, "recursively generate factory methods for all packages")
	file := flag.String("f", "", "generate factory methods only for file specific")
	flag.BoolVar(&verbose, "v", false, "verbose output")
	flag.Parse()
	run(*directory, *recurse, *file)
}
