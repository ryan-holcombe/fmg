package main

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"path"
)

const (
	generatedFileName = "fm_gen.go"
)

// allow overriding to simplify testing
var parseStructsFunc = parseStructs
var parsedImportsFunc = parseImports

func parseAllDirs(dir string) []genPackage {
	fileInfos, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Panicf("unable to recursively traverse directory [%s] - %v", dir, errors.WithStack(err))
	}

	result := make([]genPackage, 0)
	result = append(result, parseDir(dir)...)
	for _, fi := range fileInfos {
		// skip vendor directory
		if fi.IsDir() && fi.Name() != "vendor" {
			// process directories recursively
			innerDir := fmt.Sprintf("%s/%s", dir, fi.Name())
			result = append(result, parseAllDirs(innerDir)...)
		}
	}

	return result
}

func parseDir(dir string) []genPackage {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, nil, parser.ParseComments)
	if err != nil {
		log.Panicf("unable to parse directory [%s] - %v", dir, errors.WithStack(err))
	}

	result := make([]genPackage, 0)

	for _, p := range pkgs {
		parsedStructs := make([]genStruct, 0)
		parsedImports := make([]string, 0)
		for _, file := range p.Files {
			parsedStructs = append(parsedStructs, parseStructsFunc(fset, file)...)
			parsedImports = append(parsedImports, parsedImportsFunc(file)...)
		}

		result = append(result, genPackage{
			dirname: dir,
			pkg:     p.Name,
			fset:    fset,
			structs: parsedStructs,
			imports: parsedImports,
		})
	}

	return result

}

func parseFile(filename string) genFile {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		log.Panicf("unable to parse file [%s] - %v", filename, errors.WithStack(err))
	}

	d, f := path.Split(filename)

	return genFile{
		dirname:  d,
		filename: f,
		pkg:      file.Name.Name,
		structs:  parseStructsFunc(fset, file),
		imports:  parsedImportsFunc(file),
	}
}

func createGeneratedFile(dirname, pkg string, imports []string, structs []genStruct) {
	var data bytes.Buffer
	writePackageFile(&data, pkg, imports, structs)
	if err := os.WriteFile(fmt.Sprintf("%s/%s", dirname, generatedFileName), data.Bytes(), 0644); err != nil {
		log.Panicf("unable to write fm_gen.go file to %s - %v", dirname, err)
	}
}
