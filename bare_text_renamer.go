package main

import "regexp"

type bareTextRenamer struct {
	from *regexp.Regexp
	to   string
}

func newBareTextRenamer(from, to string) textRenamer {
	return &bareTextRenamer{regexp.MustCompile(regexp.QuoteMeta(from)), to}
}

func (r *bareTextRenamer) Rename(s string) string {
	return r.from.ReplaceAllString(s, r.to)
}
