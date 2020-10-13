package main

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/stretchr/testify/assert"
)

func newTestPathFinder(fs billy.Filesystem) *pathFinder {
	return newPathFinder(newRepositoryPathFinder(fs, "."), fs)
}

func TestPathFinderFindFile(t *testing.T) {
	fs := memfs.New()
	_, err := fs.Create("foo")
	assert.Nil(t, err)

	ss, err := newTestPathFinder(fs).Find(".")
	assert.Nil(t, err)
	assert.Equal(t, []string{"foo"}, ss)
}

func normalize(ss []string) []string {
	if runtime.GOOS == "windows" {
		for i, s := range ss {
			ss[i] = filepath.ToSlash(s)
		}
	}
	return ss
}

func TestPathFinderFindRecursively(t *testing.T) {
	fs := memfs.New()
	_, err := fs.Create("foo/foo")
	assert.Nil(t, err)

	ss, err := newTestPathFinder(fs).Find(".")
	assert.Nil(t, err)
	assert.Equal(t, []string{"foo", "foo/foo"}, normalize(ss))
}

func TestPathFinderIncludePathsNotIncludedInRepository(t *testing.T) {
	fs := memfs.New()
	_, err := fs.Create("foo")
	assert.Nil(t, err)

	commitFiles(t, fs, []string{"bar"})

	ss, err := newTestPathFinder(fs).Find(".")
	assert.Nil(t, err)
	assert.Equal(t, []string{"bar", "foo"}, normalize(ss))
}
