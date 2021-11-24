package main

import (
	"flag"
)

var (
	flagDirectory = flag.String("d", "./", "generate factory methods for all structs within the directory")
	flagRecurse   = flag.Bool("r", true, "recursively generate factory methods for all packages")
	flagFile      = flag.String("f", "", "generate factory methods only for file specific")
	flagStructs   = flag.String("s", "", "comma separated list of structs to generate factory methods for")
	flagVerbose   = flag.Bool("v", false, "verbose output")
)

// to allow for testing
var createGeneratedFileFunc = createGeneratedFile

func run(directory string, recurse bool, file string) {
	if file != "" {
		parsed := parseFile(file)
		createGeneratedFileFunc(parsed.dirname, parsed.pkg, parsed.imports, parsed.structs)
	} else {
		var pkgs []genPackage
		if recurse {
			pkgs = parseAllDirs(directory)
		} else {
			pkgs = parseDir(directory)
		}
		for _, pkg := range pkgs {
			createGeneratedFileFunc(pkg.dirname, pkg.pkg, pkg.imports, pkg.structs)
		}
	}
}

func main() {
	flag.Parse()
	run(*flagDirectory, *flagRecurse, *flagFile)
}
