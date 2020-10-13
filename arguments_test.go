package main

import (
	"runtime"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/stretchr/testify/assert"
)

func TestGetArguments(t *testing.T) {
	for _, ss := range [][]string{
		{"foo", "bar"},
		{"foo", "bar", "."},
		{"-b", "foo", "bar"},
		{"--bare", "foo", "bar"},
		{"-c", "camel", "foo", "bar"},
		{"--cases", "camel", "foo", "bar"},
		{"-c", "camel,kebab", "foo", "bar"},
		{"-v", "foo", "bar"},
		{"--verbose", "foo", "bar"},
		{"-h"},
		{"--help"},
		{"--version"},
	} {
		_, err := getArguments(ss)
		assert.Nil(t, err)
	}
}

func TestGetArgumentsError(t *testing.T) {
	for _, ss := range [][]string{
		{},
		{"foo"},
		{"foo", "bar", "baz", "blah"},
		{"-c", "caml", "foo", "bar"},
	} {
		_, err := getArguments(ss)
		assert.NotNil(t, err)
	}
}

func TestHelp(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip()
	}

	cupaloy.SnapshotT(t, help())
}
