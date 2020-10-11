package main

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-billy/v5/util"
	"github.com/stretchr/testify/assert"
)

func TestFileRenamerRenameFileShrinkingAfterRenaming(t *testing.T) {
	fs := memfs.New()

	err := util.WriteFile(fs, "foo", []byte("barBaz"), 0o444)
	assert.Nil(t, err)

	tr, err := newCaseTextRenamer("bar baz", "bar", nil)
	assert.Nil(t, err)

	err = newFileRenamer(fs, ioutil.Discard).Rename(tr, "foo", false)
	assert.Nil(t, err)

	f, err := fs.Open("foo")
	assert.Nil(t, err)

	bs, err := ioutil.ReadAll(f)
	assert.Nil(t, err)

	assert.Equal(t, "bar", string(bs))
}

func TestFileRenamerRenameFileWithVerboseOption(t *testing.T) {
	fs := memfs.New()

	err := util.WriteFile(fs, "foo", []byte("foo"), 0o444)
	assert.Nil(t, err)

	tr, err := newCaseTextRenamer("foo", "bar", nil)
	assert.Nil(t, err)

	b := &bytes.Buffer{}
	err = newFileRenamer(fs, b).Rename(tr, "foo", true)
	assert.Nil(t, err)

	cupaloy.SnapshotT(t, b.String())
}
