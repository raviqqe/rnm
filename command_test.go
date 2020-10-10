package main

import (
	"bytes"
	"io/ioutil"
	"regexp"
	"strings"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/stretchr/testify/assert"
	"gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-billy.v4/memfs"
	"gopkg.in/src-d/go-billy.v4/util"
)

func newTestCommand(fs billy.Filesystem) *command {
	return newCommand(newPathGlobber(fs), fs, ioutil.Discard, ioutil.Discard)
}

func TestCommandHelp(t *testing.T) {
	b := &bytes.Buffer{}
	fs := memfs.New()

	err := newCommand(newPathGlobber(fs), fs, b, ioutil.Discard).Run([]string{"--help"})
	assert.Nil(t, err)

	cupaloy.SnapshotT(t, b.String())
}

func TestCommandVersion(t *testing.T) {
	b := &bytes.Buffer{}
	fs := memfs.New()

	err := newCommand(newPathGlobber(fs), fs, b, ioutil.Discard).Run([]string{"--version"})
	assert.Nil(t, err)

	assert.True(
		t,
		regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+$`).MatchString(
			strings.TrimSpace(b.String()),
		),
	)
}

func TestCommandRun(t *testing.T) {
	fs := memfs.New()
	err := util.WriteFile(fs, "foo", []byte("foo"), 0o222)
	assert.Nil(t, err)

	err = newTestCommand(fs).Run([]string{"foo", "bar"})
	assert.Nil(t, err)

	f, err := fs.Open("bar")
	assert.Nil(t, err)

	bs, err := ioutil.ReadAll(f)
	assert.Nil(t, err)

	assert.Equal(t, "bar", string(bs))
}
