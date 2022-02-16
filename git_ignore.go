package main

import (
	"strings"

	ignore "github.com/sabhiram/go-gitignore"
)

type GitIgnore struct {
	ignore *ignore.GitIgnore
}

func NewGitIgnore(source string) *GitIgnore {
	return &GitIgnore{ignore.CompileIgnoreLines(strings.Split(source, "\n")...)}
}

func (*GitIgnore) Ignore(path string) bool {
	return true
}
