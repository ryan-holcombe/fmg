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

	// check comment for skip directive
	if g.comment != nil {
		for _, c := range skipStructComment {
			if strings.Contains(strings.ToLower(g.comment.value), c) {
				return true
			}
		}
	}

	// check flags for struct includes
	if *flagStructs == "" {
		return false
	}

	skip := true
	structSplit := strings.Split(*flagStructs, ",")
	for _, s := range structSplit {
		if strings.TrimSpace(s) == g.name {
			skip = false
		}
	}

	return skip
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

type structIncludes struct {
	structs map[string]struct{}
}

func NewStructIncludes(structStr string) *structIncludes {
	// create a map of the structs to include flag
	splitStructs := strings.Split(structStr, ",")
	structMap := make(map[string]struct{}, len(splitStructs))
	for _, s := range splitStructs {
		structMap[strings.TrimSpace(s)] = struct{}{}
	}
	return &structIncludes{structMap}
}

func (s structIncludes) Skip(structName string) bool {
	// if no structs are passed in, assume all structs should be included
	if len(s.structs) == 0 {
		return false
	}

	_, ok := s.structs[structName]
	return ok
}
