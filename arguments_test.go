package main

import (
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/stretchr/testify/assert"
)

func TestGetArguments(t *testing.T) {
	for _, ss := range [][]string{
		{"foo", "bar"},
		{"-c", "camel", "foo", "bar"},
		{"-c", "camel,kebab", "foo", "bar"},
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
		{"foo", "bar", "baz"},
		{"-c", "caml", "foo", "bar"},
	} {
		_, err := getArguments(ss)
		assert.NotNil(t, err)
	}
}

func TestHelp(t *testing.T) {
	cupaloy.SnapshotT(t, help())
}
