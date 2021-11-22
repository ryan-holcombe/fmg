package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseTag(t *testing.T) {
	t.Run("no fmgen tags", func(t *testing.T) {
		_, found := parseTag(`db:"dirname" json:"foo"`)
		assert.False(t, found)
	})

	t.Run("fmgen skip tag", func(t *testing.T) {
		results, found := parseTag(`fmgen:"-"`)
		assert.True(t, found)
		assert.Contains(t, results.values, "-")
	})

	t.Run("fmgen multiple tags", func(t *testing.T) {
		results, found := parseTag(`db:"dirname" json:"foo" fmgen:"optional,-"`)
		assert.True(t, found)
		assert.True(t, results.skip())
		assert.True(t, results.optional())
	})
}
