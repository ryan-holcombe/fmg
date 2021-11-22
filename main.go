package main

import (
	"flag"
	"log"
)

var (
	directory = flag.String("dirname", "./", "starting directory to process")
	file      = flag.String("file", "", "optional filename to process")
)

// to allow for testing
var createGeneratedFileFunc = createGeneratedFile

func run(directory, file string) {
	if file != "" {
		parsed := parseFile(file)
		if len(parsed.structs) > 0 {
			createGeneratedFileFunc(parsed.dirname, parsed.pkg, parsed.imports, parsed.structs)
		} else {
			log.Fatalf("no structs found in file %s, aborting", file)
		}
	} else {
		pkgs := parseAllDirs(directory)
		for _, pkg := range pkgs {
			if len(pkg.structs) > 0 {
				createGeneratedFileFunc(pkg.dirname, pkg.pkg, pkg.imports, pkg.structs)
			}
		}
	}
}

func main() {
	run(*directory, *file)
}
