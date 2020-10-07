package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRenamer(t *testing.T) {
	_, err := newRenamer("foo", "bar")
	assert.Nil(t, err)
}

func TestRenameDifferentCases(t *testing.T) {
	r, err := newRenamer("foo bar", "baz qux")
	assert.Nil(t, err)

	for _, ss := range [][2]string{
		{"foo_bar", "baz_qux"},
		{"FOO_BAR", "BAZ_QUX"},
		{"fooBar", "bazQux"},
		{"FooBar", "BazQux"},
	} {
		assert.Equal(t, ss[1], r.Rename(ss[0]))
	}
}

func TestRenameAcronym(t *testing.T) {
	r, err := newRenamer("u s a", "u k")
	assert.Nil(t, err)

	assert.Equal(t, "UK", r.Rename("USA"))
}

func TestRenameNameInText(t *testing.T) {
	r, err := newRenamer("foo", "bar")
	assert.Nil(t, err)

	assert.Equal(t, "ab bar cd", r.Rename("ab foo cd"))
}
