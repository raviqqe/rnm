package main

import (
	"github.com/jinzhu/inflection"
)

type renamer struct {
	patterns []*pattern
}

func newRenamer(from string, to string, cs map[caseName]struct{}) (*renamer, error) {
	if cs == nil {
		cs = allCaseNames
	}

	ps, err := compilePatterns(from, to, cs)
	if err != nil {
		return nil, err
	}

	pps, err := compilePatterns(inflection.Plural(from), inflection.Plural(to), cs)
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
