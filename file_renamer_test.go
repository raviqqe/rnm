package main

import (
	"io/ioutil"
	"testing"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/stretchr/testify/assert"
)

func TestFileRenamerRenameFileShrinkingAfterRenaming(t *testing.T) {
	fs := memfs.New()
	f, err := fs.Create("foo")
	assert.Nil(t, err)

	_, err = f.Write([]byte("barBaz"))
	assert.Nil(t, err)

	f.Close()

	tr, err := newCaseTextRenamer("bar baz", "bar", nil)
	assert.Nil(t, err)

	err = newFileRenamer(fs).Rename(tr, "foo")
	assert.Nil(t, err)

	f, err = fs.Open("foo")
	assert.Nil(t, err)

	bs, err := ioutil.ReadAll(f)
	assert.Nil(t, err)

	assert.Equal(t, "bar", string(bs))
}
