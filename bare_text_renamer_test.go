package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBareRenamerRename(t *testing.T) {
	assert.Equal(t, "bar", newBareTextRenamer("foo", "bar").Rename("foo"))
}

func TestBareRenamerRenameWithUnclosedParenthesis(t *testing.T) {
	assert.Equal(t, "bar(", newBareTextRenamer("foo(", "bar(").Rename("foo("))
}

func TestBareRenamerReplaceWithRegexpCharacter(t *testing.T) {
	assert.Equal(t, "\\", newBareTextRenamer("foo", "\\").Rename("foo"))
}
