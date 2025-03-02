package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegexpRenamerRename(t *testing.T) {
	r, err := newRegexpTextRenamer("foo", "bar")
	assert.Nil(t, err)

	assert.Equal(t, "bar", r.Rename("foo"))
}

func TestNewRegexpTextRenamerFail(t *testing.T) {
	_, err := newRegexpTextRenamer("foo(", "bar(")
	assert.NotNil(t, err)
}

func TestRegexpRenamerRenameWithMetaCharacter(t *testing.T) {
	r, err := newRegexpTextRenamer("f.o", "bar")
	assert.Nil(t, err)

	assert.Equal(t, "bar", r.Rename("foo"))
}

func TestRegexpRenamerRenameWithPlaceholder(t *testing.T) {
	r, err := newRegexpTextRenamer("foo", "($0)")
	assert.Nil(t, err)

	assert.Equal(t, "(foo)", r.Rename("foo"))
}

func TestRegexpRenamerRenameNewline(t *testing.T) {
	r, err := newRegexpTextRenamer("a", "\n")
	assert.Nil(t, err)

	assert.Equal(t, "\n", r.Rename("a"))
}
