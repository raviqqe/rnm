package main

import (
	"path/filepath"
	"regexp"
	"runtime"
	"testing"

	"github.com/go-git/go-billy/v6"
	"github.com/go-git/go-billy/v6/memfs"
	"github.com/stretchr/testify/assert"
)

func newTestFileFinder(fs billy.Filesystem) *fileFinder {
	return newFileFinder(newRepositoryFileFinder(fs), fs)
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

	ss, err := newTestFileFinder(fs).Find(".", nil, nil, false)
	assert.Nil(t, err)
	assert.Equal(t, []string{"foo"}, ss)
}

func TestFileFinderDoNotFindDirectory(t *testing.T) {
	fs := memfs.New()
	err := fs.MkdirAll("foo", 0o700)
	assert.Nil(t, err)

	ss, err := newTestFileFinder(fs).Find(".", nil, nil, false)
	assert.Nil(t, err)
	assert.Equal(t, []string{}, ss)
}

func TestFileFinderFindRecursively(t *testing.T) {
	fs := memfs.New()
	_, err := fs.Create("foo/foo")
	assert.Nil(t, err)

	ss, err := newTestFileFinder(fs).Find(".", nil, nil, false)
	assert.Nil(t, err)
	assert.Equal(t, []string{"foo/foo"}, normalizePaths(ss))
}

func TestFileFinderDoNotIncludePathsNotIncludedInRepository(t *testing.T) {
	fs := memfs.New()
	_, err := fs.Create("foo")
	assert.Nil(t, err)

	commitFiles(t, fs, []string{"bar"})

	ss, err := newTestFileFinder(fs).Find(".", nil, nil, false)
	assert.Nil(t, err)
	assert.Equal(t, []string{"bar"}, normalizePaths(ss))
}

func TestFileFinderIgnoreGitRepositoryInformation(t *testing.T) {
	fs := memfs.New()

	commitFiles(t, fs, []string{".foo"})

	ss, err := newTestFileFinder(fs).Find(".", nil, nil, true)
	assert.Nil(t, err)
	assert.Equal(t, []string{}, normalizePaths(ss))
}

func TestFileFinderDoNotFindHiddenFile(t *testing.T) {
	fs := memfs.New()
	_, err := fs.Create(".foo")
	assert.Nil(t, err)

	ss, err := newTestFileFinder(fs).Find(".", nil, nil, false)
	assert.Nil(t, err)
	assert.Equal(t, []string{}, ss)
}

func TestFileFinderFindFileInHiddenDirectory(t *testing.T) {
	fs := memfs.New()
	_, err := fs.Create(".foo/foo")
	assert.Nil(t, err)

	ss, err := newTestFileFinder(fs).Find(".foo", nil, nil, false)
	assert.Nil(t, err)
	assert.Equal(t, []string{".foo/foo"}, normalizePaths(ss))
}

func TestFileFinderDoNotFindIncludedFile(t *testing.T) {
	fs := memfs.New()

	_, err := fs.Create("foo")
	assert.Nil(t, err)

	_, err = fs.Create("bar")
	assert.Nil(t, err)

	ss, err := newTestFileFinder(fs).Find(".", regexp.MustCompile("foo"), nil, false)
	assert.Nil(t, err)
	assert.Equal(t, []string{"foo"}, normalizePaths(ss))
}

func TestFileFinderDoNotFindExcludedFile(t *testing.T) {
	fs := memfs.New()

	_, err := fs.Create("foo")
	assert.Nil(t, err)

	_, err = fs.Create("bar")
	assert.Nil(t, err)

	ss, err := newTestFileFinder(fs).Find(".", nil, regexp.MustCompile("foo"), false)
	assert.Nil(t, err)
	assert.Equal(t, []string{"bar"}, normalizePaths(ss))
}
