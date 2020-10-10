package main

import (
	"testing"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-billy/v5/util"
	"github.com/stretchr/testify/assert"
)

func newTestPathGlobber(fs billy.Filesystem) *pathGlobber {
	return newPathGlobber(newRepositoryPathFinder(fs, "."), fs)
}

func TestPathGlobberGlobFile(t *testing.T) {
	fs := memfs.New()
	err := util.WriteFile(fs, "foo", []byte("foo"), 0o222)
	assert.Nil(t, err)

	ss, err := newTestPathGlobber(fs).Glob(".")
	assert.Nil(t, err)
	assert.Equal(t, []string{"foo"}, ss)
}

func TestPathGlobberGlobRecursively(t *testing.T) {
	fs := memfs.New()
	err := util.WriteFile(fs, "foo/foo", []byte("foo"), 0o222)
	assert.Nil(t, err)

	ss, err := newTestPathGlobber(fs).Glob(".")
	assert.Nil(t, err)
	assert.Equal(t, []string{"foo", "foo/foo"}, ss)
}

func TestPathGlobberDoNotIncludePathsNotIncludedInRepository(t *testing.T) {
	fs := memfs.New()
	err := util.WriteFile(fs, "foo", []byte("foo"), 0o222)
	assert.Nil(t, err)

	commitFiles(t, fs, []string{"bar"})

	ss, err := newTestPathGlobber(fs).Glob(".")
	assert.Nil(t, err)
	assert.Equal(t, []string{"bar"}, ss)
}
