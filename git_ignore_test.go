package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGitIgnore(t *testing.T) {
	NewGitIgnore("")
}

func TestNewGitIgnoreWithSource(t *testing.T) {
	NewGitIgnore("foo\nbar")
}

func TestNewGitIgnoreIgnorePath(t *testing.T) {
	g := NewGitIgnore("foo")

	assert.True(t, g.Ignore("foo"))
	assert.False(t, g.Ignore("bar"))
}

func TestNewGitIgnoreIgnoreWildcard(t *testing.T) {
	assert.True(t, NewGitIgnore("*").Ignore("foo"))
}

func TestNewGitIgnoreIgnoreNegatingWildcard(t *testing.T) {
	g := NewGitIgnore("*\n!foo")

	assert.False(t, g.Ignore("foo"))
	assert.True(t, g.Ignore("bar"))
}

func TestNewGitIgnoreIgnoreRecursiveWildcard(t *testing.T) {
	g := NewGitIgnore("**/bar")

	assert.True(t, g.Ignore("bar"))
	assert.True(t, g.Ignore("foo/bar"))
	assert.False(t, g.Ignore("foo"))
	assert.False(t, g.Ignore("foo/foo"))
}
