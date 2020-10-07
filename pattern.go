package main

import "regexp"

type pattern struct {
	From *regexp.Regexp
	To   string
}
