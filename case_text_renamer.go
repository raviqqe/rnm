package main

import (
	"github.com/jinzhu/inflection"
)

type caseTextRenamer struct {
	patterns []*pattern
}

func newCaseTextRenamer(from string, to string, cs map[caseName]struct{}) (*caseTextRenamer, error) {
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

	return &caseTextRenamer{append(ps, pps...)}, nil
}

func (r *caseTextRenamer) Rename(s string) string {
	for _, p := range r.patterns {
		s = p.From.ReplaceAllString(s, p.To)
	}

	return s
}
