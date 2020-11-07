package main

import (
	"testing"

	"github.com/iancoleman/strcase"
	"github.com/stretchr/testify/assert"
)

// TODO Use delimiters only for heads?
func TestPatternDoNotReplaceOverlappedPatterns(t *testing.T) {
	p, err := newPattern("ab", "cd", &caseConfiguration{
		camel,
		strcase.ToLowerCamel,
		nonAlphabet,
		upperCase,
	})
	assert.Nil(t, err)

	assert.Equal(t, "cd,ab", p.Replace("ab,ab"))
}
