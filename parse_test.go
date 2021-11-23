package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestLineNum(t *testing.T) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "testdata/linenum.go", nil, parser.ParseComments)
	assert.NoError(t, err)
	assert.Len(t, file.Comments, 1)

	result := lineNum(fset, file.Comments[0].Pos())
	assert.Equal(t, 3, result)
}

func TestFindComment(t *testing.T) {
	t.Run("comment match", func(t *testing.T) {
		structLineNum := 10
		comments := []genComment{
			{
				lineNum: 9,
			},
		}
		result := findComment(structLineNum, comments)
		assert.NotNil(t, result)
		assert.Equal(t, comments[0], *result)
	})

	t.Run("comment not above struct", func(t *testing.T) {
		structLineNum := 10
		comments := []genComment{
			{
				lineNum: 1000,
			},
		}
		result := findComment(structLineNum, comments)
		assert.Nil(t, result)
	})
}

func TestLogComments(t *testing.T) {
	w := &bytes.Buffer{}
	comments := []genComment{
		{
			lineNum: 9,
			value:   "Simple struct comment",
		},
		{
			lineNum: 132,
			value:   "random comment",
		},
	}
	logComments(w, comments)

	expected := `Comments
----------------------------------------------------
Comment:
    value=Simple struct comment
    loc=9
Comment:
    value=random comment
    loc=132
`
	assert.Equal(t, expected, w.String())
}

func TestLogStructs(t *testing.T) {
	w := &bytes.Buffer{}
	structs := []genStruct{
		{
			name:    "Simple",
			lineNum: 9,
			fields: []genField{
				{
					name:     "Name",
					typ:      "string",
					ptr:      false,
					optional: true,
					skip:     false,
				},
				{
					name:     "Age",
					typ:      "int",
					ptr:      true,
					optional: false,
					skip:     false,
				},
			},
			comment: &genComment{
				lineNum: 8,
				value:   "struct comment",
			},
		},
	}

	logStructs(w, structs)

	expected := `Structs
----------------------------------------------------
Struct:
    dirname=Simple
    loc=9
    Fields:
        dirname=Name type=string ptr=false optional=true skip=false
        dirname=Age type=int ptr=true optional=false skip=false
`

	assert.Equal(t, expected, w.String())
}

func TestParseASTFile(t *testing.T) {
	t.Run("simple.go", func(t *testing.T) {
		fset := token.NewFileSet()
		astFile, err := parser.ParseFile(fset, "testdata/simple.go", nil, parser.ParseComments)
		assert.NoError(t, err)
		structs := parseStructs(fset, astFile)
		assert.Len(t, structs, 2)

		expected1 := genStruct{
			name:    "Sample",
			lineNum: 6,
			fields: []genField{
				{name: "ID", typ: "int64", optional: false, skip: true},
				{name: "Name", typ: "string", optional: false, skip: false},
				{name: "Age", typ: "int64", optional: true, skip: false},
				{name: "LastUpdated", typ: "time.Time", optional: false, skip: false},
			},
			comment: &genComment{
				lineNum: 5,
				value:   "Sample simple struct fmgen:omit\n",
			},
		}
		assert.Equal(t, expected1, structs[0])
	})

	t.Run("pointer.go", func(t *testing.T) {
		fset := token.NewFileSet()
		astFile, err := parser.ParseFile(fset, "testdata/pointer.go", nil, parser.ParseComments)
		assert.NoError(t, err)
		structs := parseStructs(fset, astFile)
		assert.Len(t, structs, 1)
		expected := genStruct{
			name:    "Pointer",
			lineNum: 6,
			fields: []genField{
				{name: "ID", typ: "int64", optional: false, skip: true},
				{name: "Name", typ: "string", optional: false, skip: false},
				{name: "Age", typ: "int64", optional: true, skip: false},
				{name: "PtrS", typ: "string", ptr: true, optional: false, skip: false},
				{name: "PtrOpt", typ: "string", ptr: true, optional: true, skip: false},
				{name: "PtrI", typ: "int", ptr: true, optional: false, skip: false},
				{name: "LastUpdated", typ: "time.Time", ptr: true, optional: false, skip: false},
			},
			comment: &genComment{
				lineNum: 5,
				value:   "Pointer struct to help test pointers\n",
			},
		}
		assert.Equal(t, expected, structs[0])
	})

	t.Run("interface.go", func(t *testing.T) {
		fset := token.NewFileSet()
		astFile, err := parser.ParseFile(fset, "testdata/interface.go", nil, parser.ParseComments)
		assert.NoError(t, err)
		structs := parseStructs(fset, astFile)
		assert.Len(t, structs, 1)
		expected := genStruct{
			name:    "impl",
			lineNum: 7,
			fields: []genField{
				{name: "i", typ: "iface", optional: false, skip: false},
			},
		}
		assert.Equal(t, expected, structs[0])
	})
}

func TestParseImports(t *testing.T) {
	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, "testdata/imports.go", nil, parser.ParseComments)
	assert.NoError(t, err)
	results := parseImports(astFile)

	expected := []string{`"net/url"`, `"time"`}
	assert.Equal(t, expected, results)
}

func TestParseFieldName(t *testing.T) {
	t.Run("object", func(t *testing.T) {
		astData := `package parse
type s struct {
Name string
}
`
		parsed, err := parser.ParseFile(token.NewFileSet(), "", []byte(astData), parser.ParseComments)
		assert.NoError(t, err)

		field := parsed.Decls[0].(*ast.GenDecl).Specs[0].(*ast.TypeSpec).Type.(*ast.StructType).Fields.List[0]
		result := parseFieldName(field)
		assert.Equal(t, "Name", result)
	})

	t.Run("array", func(t *testing.T) {
		astData := `package parse
type s struct {
Name []string
}
`
		parsed, err := parser.ParseFile(token.NewFileSet(), "", []byte(astData), parser.ParseComments)
		assert.NoError(t, err)

		field := parsed.Decls[0].(*ast.GenDecl).Specs[0].(*ast.TypeSpec).Type.(*ast.StructType).Fields.List[0]
		result := parseFieldName(field)
		assert.Equal(t, "Name", result)
	})

	t.Run("interface", func(t *testing.T) {
		astData := `package parse
type s struct {
Name interface{}
}
`
		parsed, err := parser.ParseFile(token.NewFileSet(), "", []byte(astData), parser.ParseComments)
		assert.NoError(t, err)

		field := parsed.Decls[0].(*ast.GenDecl).Specs[0].(*ast.TypeSpec).Type.(*ast.StructType).Fields.List[0]
		result := parseFieldName(field)
		assert.Equal(t, "Name", result)
	})
}

func TestBuildField(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		astData := `package parse
type s struct {
Name string
}
`
		parsed, err := parser.ParseFile(token.NewFileSet(), "", []byte(astData), parser.ParseComments)
		assert.NoError(t, err)

		field := parsed.Decls[0].(*ast.GenDecl).Specs[0].(*ast.TypeSpec).Type.(*ast.StructType).Fields.List[0]
		result := buildField(nil, field.Type, "Name", nil)
		assert.Equal(t, "Name", result.name)
		assert.Equal(t, "string", result.typ)
		assert.False(t, result.ptr)
		assert.False(t, result.array)
	})

	t.Run("string pointer", func(t *testing.T) {
		astData := `package parse
type s struct {
Name *string
}
`
		parsed, err := parser.ParseFile(token.NewFileSet(), "", []byte(astData), parser.ParseComments)
		assert.NoError(t, err)

		field := parsed.Decls[0].(*ast.GenDecl).Specs[0].(*ast.TypeSpec).Type.(*ast.StructType).Fields.List[0]
		result := buildField(nil, field.Type, "Name", nil)
		assert.Equal(t, "Name", result.name)
		assert.Equal(t, "string", result.typ)
		assert.True(t, result.ptr)
		assert.False(t, result.array)
	})

	t.Run("array", func(t *testing.T) {
		astData := `package parse
type s struct {
Name []string
}
`
		parsed, err := parser.ParseFile(token.NewFileSet(), "", []byte(astData), parser.ParseComments)
		assert.NoError(t, err)

		field := parsed.Decls[0].(*ast.GenDecl).Specs[0].(*ast.TypeSpec).Type.(*ast.StructType).Fields.List[0]
		result := buildField(nil, field.Type, "Name", nil)
		assert.Equal(t, "Name", result.name)
		assert.Equal(t, "string", result.typ)
		assert.False(t, result.ptr)
		assert.True(t, result.array)
	})

	t.Run("time.Time", func(t *testing.T) {
		astData := `package parse
type s struct {
Name time.Time
}
`
		parsed, err := parser.ParseFile(token.NewFileSet(), "", []byte(astData), parser.ParseComments)
		assert.NoError(t, err)

		field := parsed.Decls[0].(*ast.GenDecl).Specs[0].(*ast.TypeSpec).Type.(*ast.StructType).Fields.List[0]
		result := buildField(nil, field.Type, "Name", nil)
		assert.Equal(t, "Name", result.name)
		assert.Equal(t, "time.Time", result.typ)
		assert.False(t, result.ptr)
		assert.False(t, result.array)
	})
}
