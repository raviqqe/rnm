package main

import "regexp"

type compiledPattern struct {
	From *regexp.Regexp
	To   string
}
