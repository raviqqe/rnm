package main

import (
	"github.com/jinzhu/inflection"
)

type renamer struct {
	patterns []*pattern
}

func newRenamer(from string, to string, enabled map[caseName]struct{}) (*renamer, error) {
	if enabled == nil {
		enabled = allCaseNames
	}

	ps, err := compilePatterns(from, to, enabled)
	if err != nil {
		return nil, err
	}

	pps, err := compilePatterns(inflection.Plural(from), inflection.Plural(to), enabled)
	if err != nil {
		return nil, err
	}

	return &renamer{append(ps, pps...)}, nil
}

func (r *renamer) Rename(s string) string {
	for _, p := range r.patterns {
		s = p.From.ReplaceAllString(s, p.To)
	}

	return s
}
