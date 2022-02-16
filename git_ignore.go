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

func (i *GitIgnore) Ignore(path string) bool {
	return i.ignore.MatchesPath(path)
}
