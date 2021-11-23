package main

import (
	"go/token"
	"strings"
)

var (
	skipStructComment = []string{"fmgen:-", "fmgen:skip", "fmgen:exclude"}
)

type genField struct {
	name     string
	typ      string
	optional bool
	skip     bool
	ptr      bool
	array    bool
}

type genComment struct {
	lineNum int
	value   string
}

type genStruct struct {
	name    string
	lineNum int
	fields  []genField
	comment *genComment
}

func (g genStruct) Skip() bool {
	if g.comment == nil {
		return false
	}

	for _, c := range skipStructComment {
		if strings.Contains(strings.ToLower(g.comment.value), c) {
			return true
		}
	}

	return false
}

type genPackage struct {
	dirname string
	pkg     string
	fset    *token.FileSet
	structs []genStruct
	imports []string
}

type genFile struct {
	dirname  string
	filename string
	pkg      string
	fset     *token.FileSet
	structs  []genStruct
	imports  []string
}
