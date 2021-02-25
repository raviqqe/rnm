package main

import (
	"bytes"
	"io"
	"os"
	"runtime"
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

	err = newFileRenamer(fs, io.Discard).Rename(tr, "foo", ".", false)
	assert.Nil(t, err)

	f, err := fs.Open("foo")
	assert.Nil(t, err)

	bs, err := io.ReadAll(f)
	assert.Nil(t, err)

	assert.Equal(t, "bar", string(bs))
}

func TestFileRenamerRenameFileWithVerboseOption(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip()
	}

	fs := memfs.New()

	err := util.WriteFile(fs, "foo", []byte("foo"), 0o444)
	assert.Nil(t, err)

	tr, err := newCaseTextRenamer("foo", "bar", nil)
	assert.Nil(t, err)

	b := &bytes.Buffer{}
	err = newFileRenamer(fs, b).Rename(tr, "foo", ".", true)
	assert.Nil(t, err)

	cupaloy.SnapshotT(t, b.String())
}

func TestFileRenamerRenameSymlinkToFile(t *testing.T) {
	fs := memfs.New()

	err := util.WriteFile(fs, "baz", []byte("foo"), 0o444)
	assert.Nil(t, err)

	err = fs.Symlink("baz", "foo")
	assert.Nil(t, err)

	tr, err := newCaseTextRenamer("foo", "bar", nil)
	assert.Nil(t, err)

	err = newFileRenamer(fs, io.Discard).Rename(tr, "foo", ".", true)
	assert.Nil(t, err)

	i, err := fs.Lstat("bar")
	assert.Nil(t, err)

	assert.True(t, i.Mode()&os.ModeSymlink > 0)

	f, err := fs.Open("bar")
	assert.Nil(t, err)

	bs, err := io.ReadAll(f)
	assert.Nil(t, err)

	assert.Equal(t, "bar", string(bs))
}

func TestFileRenamerFailToRenameDirectory(t *testing.T) {
	fs := memfs.New()

	err := fs.MkdirAll("foo", 0o755)
	assert.Nil(t, err)

	tr, err := newCaseTextRenamer("foo", "bar", nil)
	assert.Nil(t, err)

	err = newFileRenamer(fs, io.Discard).Rename(tr, "foo", ".", true)
	assert.Error(t, err)
}

func TestFileRenamerFailToRenameSymlinkToDirectory(t *testing.T) {
	fs := memfs.New()

	err := fs.MkdirAll("baz", 0o755)
	assert.Nil(t, err)

	err = fs.Symlink("baz", "foo")
	assert.Nil(t, err)

	tr, err := newCaseTextRenamer("foo", "bar", nil)
	assert.Nil(t, err)

	err = newFileRenamer(fs, io.Discard).Rename(tr, "foo", ".", true)
	assert.Error(t, err)
}

func TestFileRenamerFileInDirectory(t *testing.T) {
	fs := memfs.New()

	err := fs.MkdirAll("foo", 0o755)
	assert.Nil(t, err)

	_, err = fs.Create("foo/foo")
	assert.Nil(t, err)

	tr, err := newCaseTextRenamer("foo", "bar", nil)
	assert.Nil(t, err)

	err = newFileRenamer(fs, io.Discard).Rename(tr, "foo/foo", "foo", true)
	assert.Nil(t, err)

	_, err = fs.Lstat("foo/bar")
	assert.Nil(t, err)
}
