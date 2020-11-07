package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCaseTextRenamer(t *testing.T) {
	_, err := newCaseTextRenamer("foo", "bar", nil)
	assert.Nil(t, err)
}

func TestCaseTextRenamerRenameDifferentCases(t *testing.T) {
	r, err := newCaseTextRenamer("foo bar", "baz qux", nil)
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

func TestCaseTextRenamerDoNotRenameDifferentCases(t *testing.T) {
	r, err := newCaseTextRenamer("foo", "bar", nil)
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

func TestCaseTextRenamerRenameNameWithSpecificCase(t *testing.T) {
	r, err := newCaseTextRenamer("bar", "bar baz", map[caseName]struct{}{kebab: {}})
	assert.Nil(t, err)

	assert.Equal(t, "foo-bar-baz-baz", r.Rename("foo-bar-baz"))
}

func TestCaseTextRenamerRenameAcronym(t *testing.T) {
	r, err := newCaseTextRenamer("u s a", "u k", nil)
	assert.Nil(t, err)

	assert.Equal(t, "UK", r.Rename("USA"))
}

func TestCaseTextRenamerRenamePlurals(t *testing.T) {
	r, err := newCaseTextRenamer("bad apple", "nice orange", nil)
	assert.Nil(t, err)

	assert.Equal(t, "NiceOranges", r.Rename("BadApples"))
}

func TestCaseTextRenamerRenameNameInText(t *testing.T) {
	r, err := newCaseTextRenamer("foo", "bar", nil)
	assert.Nil(t, err)

	assert.Equal(t, "ab bar cd", r.Rename("ab foo cd"))
}

func TestCaseTextRenamerRenameCamelCaseMatchingBothLowerAndUpperCamelPatterns(t *testing.T) {
	r, err := newCaseTextRenamer("foo bar", "baz foo bar", nil)
	assert.Nil(t, err)

	assert.Equal(t, "bazFooBar", r.Rename("fooBar"))
}

func TestCaseTextRenamerRenameNameWithRegexpCharacters(t *testing.T) {
	r, err := newCaseTextRenamer("foo(", "bar(", nil)
	assert.Nil(t, err)

	assert.Equal(t, "bar()", r.Rename("foo()"))
}

func TestCaseTextRenamerFallbackToDefaultText(t *testing.T) {
	r, err := newCaseTextRenamer("foo", "ふー", nil)
	assert.Nil(t, err)

	assert.Equal(t, "ふー", r.Rename("foo"))
}

func TestCaseTextRenamerReplaceWithPatternIncludingRegexpCharacter(t *testing.T) {
	r, err := newCaseTextRenamer("foo", "\\", nil)
	assert.Nil(t, err)

	assert.Equal(t, "\\", r.Rename("foo"))
}

func TestCaseTextRenamerReplaceWithPluralPattern(t *testing.T) {
	r, err := newCaseTextRenamer("apples", "orange", nil)
	assert.Nil(t, err)

	assert.Equal(t, "orange", r.Rename("apples"))
}
