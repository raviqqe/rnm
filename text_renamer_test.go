package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTextRenamer(t *testing.T) {
	_, err := newTextRenamer("foo", "bar", nil)
	assert.Nil(t, err)
}

func TestRenameDifferentCases(t *testing.T) {
	r, err := newTextRenamer("foo bar", "baz qux", nil)
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
	r, err := newTextRenamer("foo", "bar", nil)
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
	r, err := newTextRenamer("bar", "bar baz", map[caseName]struct{}{kebab: {}})
	assert.Nil(t, err)

	assert.Equal(t, "foo-bar-baz-baz", r.Rename("foo-bar-baz"))
}

func TestRenameAcronym(t *testing.T) {
	r, err := newTextRenamer("u s a", "u k", nil)
	assert.Nil(t, err)

	assert.Equal(t, "UK", r.Rename("USA"))
}

func TestRenamePlurals(t *testing.T) {
	r, err := newTextRenamer("bad apple", "nice orange", nil)
	assert.Nil(t, err)

	assert.Equal(t, "NiceOranges", r.Rename("BadApples"))
}

func TestRenameNameInText(t *testing.T) {
	r, err := newTextRenamer("foo", "bar", nil)
	assert.Nil(t, err)

	assert.Equal(t, "ab bar cd", r.Rename("ab foo cd"))
}

func TestRenameBarePattern(t *testing.T) {
	r, err := newTextRenamer("foo/v1", "foo/v2", nil)
	assert.Nil(t, err)

	assert.Equal(t, "foo/v2", r.Rename("foo/v1"))
}
