package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRun(t *testing.T) {
	defer func() {
		createGeneratedFileFunc = createGeneratedFile
	}()

	t.Run("file", func(t *testing.T) {
		var executed bool
		createGeneratedFileFunc = func(dirname, pkg string, imports []string, structs []genStruct) {
			assert.Equal(t, "testdata", pkg)
			assert.Equal(t, "testdata/", dirname)
			executed = true
		}
		run("", "testdata/test.go")
		assert.True(t, executed)
	})

	t.Run("directory", func(t *testing.T) {
		var runCnt int
		createGeneratedFileFunc = func(dirname, pkg string, imports []string, structs []genStruct) {
			runCnt++
		}
		run("testdata/", "")
		assert.Equal(t, 3, runCnt)
	})
}