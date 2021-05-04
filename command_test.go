package main

import (
	"bytes"
	"regexp"
	"strings"
	"testing"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-billy/v5/util"
	"github.com/stretchr/testify/assert"
)

func newTestCommand(fs billy.Filesystem, d string) *command {
	return newCommand(
		newArgumentParser(d),
		newFileFinder(newRepositoryFileFinder(fs), fs),
		newFileRenamer(fs, io.Discard),
		fs,
		io.Discard,
		io.Discard,
	)
}

func TestCommandHelp(t *testing.T) {
	b := &bytes.Buffer{}
	fs := memfs.New()

	err := newCommand(
		newArgumentParser("."),
		newFileFinder(newRepositoryFileFinder(fs), fs),
		newFileRenamer(fs, io.Discard),
		fs,
		b,
		io.Discard,
	).Run([]string{"--help"})
	assert.Nil(t, err)

	assert.Greater(t, len(b.String()), 0)
}

func TestCommandVersion(t *testing.T) {
	b := &bytes.Buffer{}
	fs := memfs.New()

	err := newCommand(
		newArgumentParser("."),
		newFileFinder(newRepositoryFileFinder(fs), fs),
		newFileRenamer(fs, io.Discard),
		fs,
		b,
		io.Discard,
	).Run([]string{"--version"})
	assert.Nil(t, err)

	assert.True(
		t,
		regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+$`).MatchString(
			strings.TrimSpace(b.String()),
		),
	)
}

func TestCommandRenameWithoutPathOption(t *testing.T) {
	fs := memfs.New()
	err := util.WriteFile(fs, "foo", []byte("foo"), 0o444)
	assert.Nil(t, err)

	err = newTestCommand(fs, ".").Run([]string{"foo", "bar"})
	assert.Nil(t, err)

	f, err := fs.Open("bar")
	assert.Nil(t, err)

	bs, err := io.ReadAll(f)
	assert.Nil(t, err)

	assert.Equal(t, "bar", string(bs))
}

func TestCommandRenameOnlyFile(t *testing.T) {
	fs := memfs.New()
	err := util.WriteFile(fs, "foo", []byte("baz"), 0o444)
	assert.Nil(t, err)

	err = util.WriteFile(fs, "bar", []byte("baz"), 0o444)
	assert.Nil(t, err)

	err = newTestCommand(fs, ".").Run([]string{"baz", "blah", "foo"})
	assert.Nil(t, err)

	f, err := fs.Open("foo")
	assert.Nil(t, err)

	bs, err := io.ReadAll(f)
	assert.Nil(t, err)

	assert.Equal(t, "blah", string(bs))

	f, err = fs.Open("bar")
	assert.Nil(t, err)

	bs, err = io.ReadAll(f)
	assert.Nil(t, err)

	assert.Equal(t, "baz", string(bs))
}

func TestCommandRenameOnlyDirectory(t *testing.T) {
	fs := memfs.New()

	err := fs.MkdirAll("foo", 0o755)
	assert.Nil(t, err)

	err = util.WriteFile(fs, "foo/foo", []byte("baz"), 0o444)
	assert.Nil(t, err)

	err = util.WriteFile(fs, "bar", []byte("baz"), 0o444)
	assert.Nil(t, err)

	err = newTestCommand(fs, ".").Run([]string{"baz", "blah", "foo"})
	assert.Nil(t, err)

	f, err := fs.Open("foo/foo")
	assert.Nil(t, err)

	bs, err := io.ReadAll(f)
	assert.Nil(t, err)

	assert.Equal(t, "blah", string(bs))

	f, err = fs.Open("bar")
	assert.Nil(t, err)

	bs, err = io.ReadAll(f)
	assert.Nil(t, err)

	assert.Equal(t, "baz", string(bs))
}

func TestCommandRenameWithBarePattern(t *testing.T) {
	fs := memfs.New()
	err := util.WriteFile(fs, "foo", []byte("foo()"), 0o444)
	assert.Nil(t, err)

	err = newTestCommand(fs, ".").Run([]string{"--bare", "foo(", "bar("})
	assert.Nil(t, err)

	f, err := fs.Open("foo")
	assert.Nil(t, err)

	bs, err := io.ReadAll(f)
	assert.Nil(t, err)

	assert.Equal(t, "bar()", string(bs))
}

func TestCommandRenameInDirectory(t *testing.T) {
	fs := memfs.New()
	err := fs.MkdirAll("foo", 0o755)
	assert.Nil(t, err)

	err = util.WriteFile(fs, "foo/foo", []byte("foo"), 0o444)
	assert.Nil(t, err)

	err = newTestCommand(fs, "foo").Run([]string{"foo", "bar"})
	assert.Nil(t, err)

	f, err := fs.Open("foo/bar")
	assert.Nil(t, err)

	bs, err := io.ReadAll(f)
	assert.Nil(t, err)

	assert.Equal(t, "bar", string(bs))
}

func TestCommandRenameInDirectoryWithFileSpecified(t *testing.T) {
	fs := memfs.New()
	err := fs.MkdirAll("foo", 0o755)
	assert.Nil(t, err)

	err = util.WriteFile(fs, "foo/foo", []byte("foo"), 0o444)
	assert.Nil(t, err)

	err = newTestCommand(fs, "foo").Run([]string{"foo", "bar", "foo"})
	assert.Nil(t, err)

	f, err := fs.Open("foo/bar")
	assert.Nil(t, err)

	bs, err := io.ReadAll(f)
	assert.Nil(t, err)

	assert.Equal(t, "bar", string(bs))
}

func TestCommandRenameWithRegularExpression(t *testing.T) {
	fs := memfs.New()

	err := util.WriteFile(fs, "foo", []byte("foo"), 0o444)
	assert.Nil(t, err)

	err = newTestCommand(fs, ".").Run([]string{"-r", "(f.)o", "${1}e"})
	assert.Nil(t, err)

	f, err := fs.Open("foe")
	assert.Nil(t, err)

	bs, err := io.ReadAll(f)
	assert.Nil(t, err)

	assert.Equal(t, "foe", string(bs))
}
