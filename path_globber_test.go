package main

import (
	"testing"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/stretchr/testify/assert"
)

func newTestPathGlobber(fs billy.Filesystem) *pathGlobber {
	return newPathGlobber(newRepositoryPathFinder(fs, "."), fs)
}

func TestPathGlobberGlobFile(t *testing.T) {
	fs := memfs.New()
	_, err := fs.Create("foo")
	assert.Nil(t, err)

	ss, err := newTestPathGlobber(fs).Glob(".")
	assert.Nil(t, err)
	assert.Equal(t, []string{"foo"}, ss)
}

func TestPathGlobberGlobRecursively(t *testing.T) {
	fs := memfs.New()
	_, err := fs.Create("foo/foo")
	assert.Nil(t, err)

	ss, err := newTestPathGlobber(fs).Glob(".")
	assert.Nil(t, err)
	assert.Equal(t, []string{"foo", "foo/foo"}, ss)
}

func TestPathGlobberIncludePathsNotIncludedInRepository(t *testing.T) {
	fs := memfs.New()
	_, err := fs.Create("foo")
	assert.Nil(t, err)

	commitFiles(t, fs, []string{"bar"})

	ss, err := newTestPathGlobber(fs).Glob(".")
	assert.Nil(t, err)
	assert.Equal(t, []string{"bar", "foo"}, ss)
}
