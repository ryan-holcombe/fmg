package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"io"
	"log"
	"os"
)

func lineNum(fset *token.FileSet, pos token.Pos) int {
	return fset.File(pos).Line(pos)
}

// if the comment is 1 line before the struct definition, then consider it a struct comment
func findComment(lineNum int, comments []genComment) *genComment {
	for _, c := range comments {
		if c.lineNum+1 == lineNum {
			return &c
		}
	}
	return nil
}

func logComments(w io.Writer, comments []genComment) {
	if len(comments) > 0 {
		fmt.Fprintln(w, "Comments")
		fmt.Fprintln(w, "----------------------------------------------------")
		for _, c := range comments {
			fmt.Fprintf(w, "Comment:\n")
			fmt.Fprintf(w, "    value=%s\n", c.value)
			fmt.Fprintf(w, "    loc=%d\n", c.lineNum)
		}
	}
}

func logStructs(w io.Writer, structs []genStruct) {
	if len(structs) > 0 {
		fmt.Fprintln(w, "Structs")
		fmt.Fprintln(w, "----------------------------------------------------")
		for _, s := range structs {
			fmt.Fprintf(w, "Struct:\n")
			fmt.Fprintf(w, "    dirname=%s\n", s.name)
			fmt.Fprintf(w, "    loc=%d\n", s.lineNum)
			fmt.Fprintf(w, "    Fields:\n")
			for _, f := range s.fields {
				fmt.Fprintf(w, "        dirname=%s type=%s ptr=%t optional=%t skip=%t\n", f.name, f.typ, f.ptr, f.optional, f.skip)
			}
		}
	}
}

func logImports(w io.Writer, imports []string) {
	if len(imports) > 0 {
		fmt.Fprintln(w, "Imports")
		fmt.Fprintln(w, "----------------------------------------------------")
		for _, i := range imports {
			fmt.Fprintf(w, "    %s\n", i)
		}
	}
}

func parseImports(node *ast.File) []string {
	imports := make([]string, 0)
	for _, i := range node.Imports {
		imports = append(imports, i.Path.Value)
	}

	if *flagVerbose {
		logImports(os.Stdout, imports)
	}

	return imports
}

func parseFieldName(field *ast.Field) string {
	return field.Names[0].Name
}

func buildField(field *genField, expr ast.Expr, fieldName string, fieldTag *ast.BasicLit) *genField {
	if field == nil {
		var tags tag
		if fieldTag != nil {
			tags, _ = parseTag(fieldTag.Value)
		}

		field = &genField{
			name:     fieldName,
			optional: tags.optional(),
			skip:     tags.skip(),
		}
	}

	var typ string
	switch fieldType := expr.(type) {
	case *ast.Ident:
		typ = fieldType.Name
	case *ast.SelectorExpr:
		typ = fmt.Sprintf("%s.%s", fieldType.X.(*ast.Ident).Name, fieldType.Sel.Name)
	case *ast.StarExpr:
		field.ptr = true
		return buildField(field, fieldType.X, fieldName, fieldTag)
	case *ast.ArrayType:
		field.array = true
		return buildField(field, fieldType.Elt, fieldName, fieldTag)
	default:
		log.Panicf("skipping field type - %v\n", fieldType)
	}

	field.typ = typ
	return field
}

func parseStructs(fset *token.FileSet, node *ast.File) []genStruct {

	// process all comments in the file to match with structs later
	var comments []genComment
	for _, c := range node.Comments {
		comments = append(comments, genComment{
			lineNum: lineNum(fset, c.Pos()),
			value:   c.Text(),
		})
	}

	if *flagVerbose {
		logComments(os.Stdout, comments)
	}

	var structs []genStruct

	// look for structs and within the file, parse out the fields and tags into a genStruct
	for _, decl := range node.Decls {
		switch decl.(type) {

		case *ast.GenDecl:
			genDecl := decl.(*ast.GenDecl)
			for _, spec := range genDecl.Specs {
				switch spec.(type) {
				case *ast.TypeSpec:
					typeSpec := spec.(*ast.TypeSpec)

					structName := typeSpec.Name.Name
					structLineNum := lineNum(fset, typeSpec.Pos())
					structFields := make([]genField, 0)

					switch typeSpec.Type.(type) {
					case *ast.StructType:
						structType := typeSpec.Type.(*ast.StructType)
						for _, field := range structType.Fields.List {
							fieldStruct := buildField(nil, field.Type, parseFieldName(field), field.Tag)
							structFields = append(structFields, *fieldStruct)
						}
						structs = append(structs, genStruct{
							name:    structName,
							lineNum: structLineNum,
							fields:  structFields,
							comment: findComment(structLineNum, comments),
						})
					default:
						log.Printf("skipping spec type in [%s], struct [%s] - %v\n", node.Name.Name, structName, typeSpec.Type)
					}
				}
			}
		}
	}

	if *flagVerbose {
		logStructs(os.Stdout, structs)
	}

	return structs
}
