package main

import (
	"github.com/stretchr/testify/assert"
	"go/ast"
	"go/token"
	"testing"
)

func TestParseAllDirs(t *testing.T) {
	parseStructsFunc = func(fset *token.FileSet, node *ast.File) []genStruct {
		return nil
	}
	defer func() {
		parseStructsFunc = parseStructs
	}()

	t.Run("invalid directory", func(t *testing.T) {
		assert.Panics(t, func() {
			parseAllDirs("./invaliddirname")
		})
	})

	t.Run("parse testdata dirname", func(t *testing.T) {
		results := parseAllDirs("./testdata")
		assert.Len(t, results, 3)
		assert.Equal(t, "./testdata", results[0].dirname)
		assert.Equal(t, "testdata", results[0].pkg)
		assert.Equal(t, "./testdata/pkg", results[1].dirname)
		assert.Equal(t, "empty", results[1].pkg)
		assert.Equal(t, "./testdata/recurse", results[2].dirname)
		assert.Equal(t, "recurse", results[2].pkg)
	})
}

func TestParseDir(t *testing.T) {
	parseStructsFunc = func(fset *token.FileSet, node *ast.File) []genStruct {
		return nil
	}
	defer func() {
		parseStructsFunc = parseStructs
	}()

	t.Run("invalid directory", func(t *testing.T) {
		assert.Panics(t, func() {
			parseDir("./invaliddirname")
		})
	})

	t.Run("skip test files", func(t *testing.T) {
		results := parseDir("./testdata")
		assert.Len(t, results, 1)
		for _, s := range results[0].structs {
			assert.NotEqual(t, "ShouldSkip", s.name)
		}
	})

	t.Run("parse testdata dirname", func(t *testing.T) {
		results := parseDir("./testdata")
		assert.Len(t, results, 1)
		assert.Equal(t, "testdata", results[0].pkg)
		assert.Equal(t, "./testdata", results[0].dirname)
	})
}

func TestParseFile(t *testing.T) {
	parseStructsFunc = func(fset *token.FileSet, node *ast.File) []genStruct {
		return nil
	}
	defer func() {
		parseStructsFunc = parseStructs
	}()

	t.Run("invalid file dirname", func(t *testing.T) {
		assert.Panics(t, func() {
			parseFile("testdata/invalidfile.go")
		})
	})

	t.Run("parse testdata/simple.go file", func(t *testing.T) {
		results := parseFile("testdata/simple.go")
		assert.Equal(t, "simple.go", results.filename)
		assert.Equal(t, "testdata/", results.dirname)
		assert.Equal(t, "testdata", results.pkg)
	})
}
