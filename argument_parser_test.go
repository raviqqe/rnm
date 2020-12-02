package main

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/stretchr/testify/assert"
)

func TestParseArguments(t *testing.T) {
	for _, ss := range [][]string{
		{"foo", "bar"},
		{"foo", "bar", "."},
		{"-b", "foo", "bar"},
		{"--bare", "foo", "bar"},
		{"-c", "camel", "foo", "bar"},
		{"--cases", "camel", "foo", "bar"},
		{"-c", "camel,kebab", "foo", "bar"},
		{"-i", "foo", "foo", "bar"},
		{"--include", "foo", "foo", "bar"},
		{"-e", "foo", "foo", "bar"},
		{"--exclude", "foo", "foo", "bar"},
		{"--ignore-git", "foo", "bar"},
		{"-v", "foo", "bar"},
		{"--verbose", "foo", "bar"},
		{"-h"},
		{"--help"},
		{"--version"},
	} {
		_, err := newArgumentParser(".").Parse(ss)
		assert.Nil(t, err)
	}
}

func TestParseArgumentsError(t *testing.T) {
	for _, ss := range [][]string{
		{},
		{"foo"},
		{"foo", "bar", "baz", "blah"},
		{"-c", "caml", "foo", "bar"},
		{"--exclude", "(", "foo", "bar"},
	} {
		_, err := newArgumentParser(".").Parse(ss)
		assert.NotNil(t, err)
	}
}

func TestParseArgumentsResolvingPath(t *testing.T) {
	args, err := newArgumentParser("foo").Parse([]string{"foo", "foo", "bar"})
	assert.Nil(t, err)

	assert.Equal(t, filepath.FromSlash("foo/bar"), args.Path)
}

func TestHelp(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip()
	}

	cupaloy.SnapshotT(t, newArgumentParser(".").Help())
}
