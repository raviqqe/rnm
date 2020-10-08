package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRenamer(t *testing.T) {
	_, err := newRenamer("foo", "bar", nil)
	assert.Nil(t, err)
}

func TestRenameDifferentCases(t *testing.T) {
	r, err := newRenamer("foo bar", "baz qux", nil)
	assert.Nil(t, err)

	for _, ss := range [][2]string{
		{"foo_bar", "baz_qux"},
		{"FOO_BAR", "BAZ_QUX"},
		{"fooBar", "bazQux"},
		{"FooBar", "BazQux"},
		{"foo bar", "baz qux"},
		{"FOO BAR", "BAZ QUX"},
		{"AfooBar", "AfooBar"},
		{" FooBar ", " BazQux "},
		{"aFooBar", "aBazQux"},
	} {
		assert.Equal(t, ss[1], r.Rename(ss[0]))
	}
}

func TestDoNotRenameDifferentCases(t *testing.T) {
	r, err := newRenamer("foo", "bar", nil)
	assert.Nil(t, err)

	for _, s := range []string{
		"ffoo",
		"fooo",
		"FFOO",
		"FOOO",
	} {
		assert.Equal(t, s, r.Rename(s))
	}
}

func TestRenameNameWithSpecificCase(t *testing.T) {
	r, err := newRenamer("bar", "bar baz", map[caseName]struct{}{kebab: {}})
	assert.Nil(t, err)

	assert.Equal(t, "foo-bar-baz-baz", r.Rename("foo-bar-baz"))
}

func TestRenameAcronym(t *testing.T) {
	r, err := newRenamer("u s a", "u k", nil)
	assert.Nil(t, err)

	assert.Equal(t, "UK", r.Rename("USA"))
}

func TestRenameNameInText(t *testing.T) {
	r, err := newRenamer("foo", "bar", nil)
	assert.Nil(t, err)

	assert.Equal(t, "ab bar cd", r.Rename("ab foo cd"))
}
