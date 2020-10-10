package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/src-d/go-billy.v4/memfs"
	"gopkg.in/src-d/go-billy.v4/util"
)

func TestFileGlobberGlobFile(t *testing.T) {
	fs := memfs.New()
	err := util.WriteFile(fs, "foo", []byte("foo"), 0222)
	assert.Nil(t, err)

	ss, err := newFileGlobber(fs).Glob(".")
	assert.Nil(t, err)
	assert.Equal(t, []string{"foo"}, ss)
}

func TestFileGlobberGlobRecursively(t *testing.T) {
	fs := memfs.New()
	err := util.WriteFile(fs, "foo/foo", []byte("foo"), 0222)
	assert.Nil(t, err)

	ss, err := newFileGlobber(fs).Glob(".")
	assert.Nil(t, err)
	assert.Equal(t, []string{"foo", "foo/foo"}, ss)
}
