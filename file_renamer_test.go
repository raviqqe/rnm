package main

import (
	"io/ioutil"
	"testing"

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

	err = newFileRenamer(fs).Rename(tr, "foo")
	assert.Nil(t, err)

	f, err := fs.Open("foo")
	assert.Nil(t, err)

	bs, err := ioutil.ReadAll(f)
	assert.Nil(t, err)

	assert.Equal(t, "bar", string(bs))
}
