package main

import "regexp"

type regexpTextRenamer struct {
	from *regexp.Regexp
	to   string
}

func newRegexpTextRenamer(from, to string) (textRenamer, error) {
	r, err := regexp.Compile(from)
	if err != nil {
		return nil, err
	}

	return &regexpTextRenamer{r, to}, nil
}

func (r *regexpTextRenamer) Rename(s string) string {
	return r.from.ReplaceAllString(s, r.to)
}
