package main

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/stretchr/testify/assert"
)

func newTestFileFinder(fs billy.Filesystem) *fileFinder {
	return newFileFinder(newRepositoryPathFinder(fs), fs)
}

func normalizePaths(ss []string) []string {
	if runtime.GOOS == "windows" {
		for i, s := range ss {
			ss[i] = filepath.ToSlash(s)
		}
	}
	return ss
}

func TestFileFinderFindFile(t *testing.T) {
	fs := memfs.New()
	_, err := fs.Create("foo")
	assert.Nil(t, err)

	ss, err := newTestFileFinder(fs).Find(".", false)
	assert.Nil(t, err)
	assert.Equal(t, []string{"foo"}, ss)
}

func TestFileFinderDoNotFindDirectory(t *testing.T) {
	fs := memfs.New()
	err := fs.MkdirAll("foo", 0o755)
	assert.Nil(t, err)

	ss, err := newTestFileFinder(fs).Find(".", false)
	assert.Nil(t, err)
	assert.Equal(t, []string{}, ss)
}

func TestFileFinderFindRecursively(t *testing.T) {
	fs := memfs.New()
	_, err := fs.Create("foo/foo")
	assert.Nil(t, err)

	ss, err := newTestFileFinder(fs).Find(".", false)
	assert.Nil(t, err)
	assert.Equal(t, []string{"foo/foo"}, normalizePaths(ss))
}

func TestFileFinderIncludePathsNotIncludedInRepository(t *testing.T) {
	fs := memfs.New()
	_, err := fs.Create("foo")
	assert.Nil(t, err)

	commitFiles(t, fs, []string{"bar"})

	ss, err := newTestFileFinder(fs).Find(".", false)
	assert.Nil(t, err)
	assert.Equal(t, []string{"bar", "foo"}, normalizePaths(ss))
}

func TestFileFinderDoNotFindHiddenFile(t *testing.T) {
	fs := memfs.New()
	_, err := fs.Create(".foo")
	assert.Nil(t, err)

	ss, err := newTestFileFinder(fs).Find(".", false)
	assert.Nil(t, err)
	assert.Equal(t, []string{}, ss)
}

func TestFileFinderFindFileInHiddenDirectory(t *testing.T) {
	fs := memfs.New()
	_, err := fs.Create(".foo/foo")
	assert.Nil(t, err)

	ss, err := newTestFileFinder(fs).Find(".foo", false)
	assert.Nil(t, err)
	assert.Equal(t, []string{".foo/foo"}, ss)
}
