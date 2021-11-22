package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenStruct(t *testing.T) {
	t.Run("skip struct based on comment", func(t *testing.T) {
		s := genStruct{
			lineNum: 9,
			comment: &genComment{
				lineNum: 8,
				value:   "Simple struct fmgen:skip",
			},
		}

		assert.True(t, s.Skip())
	})

	t.Run("comment with no skip", func(t *testing.T) {
		s := genStruct{
			lineNum: 9,
			comment: &genComment{
				lineNum: 8,
				value:   "Simple struct",
			},
		}

		assert.False(t, s.Skip())
	})

	t.Run("no comment", func(t *testing.T) {
		s := genStruct{
			lineNum: 9,
			comment: nil,
		}

		assert.False(t, s.Skip())
	})
}
