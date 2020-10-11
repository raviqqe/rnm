package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBareRenamerRename(t *testing.T) {
	assert.Equal(t, "bar", newBareTextRenamer("foo", "bar").Rename("foo"))
}

func TestBareRenamerRenameWithRegexpCharacter(t *testing.T) {
	assert.Equal(t, "bar(", newBareTextRenamer("foo(", "bar(").Rename("foo("))
}
