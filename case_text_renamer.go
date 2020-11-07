package main

import (
	"regexp"

	"github.com/jinzhu/inflection"
)

type caseTextRenamer struct {
	patterns []*pattern
}

func newCaseTextRenamer(from string, to string, cs map[caseName]struct{}) (textRenamer, error) {
	if cs == nil {
		cs = allCaseNames
	}

	from = regexp.QuoteMeta(from)

	ps, err := compilePatterns(from, to, cs)
	if err != nil {
		return nil, err
	}

	froms := inflection.Plural(from)

	if froms != from {
		pps, err := compilePatterns(froms, inflection.Plural(to), cs)
		if err != nil {
			return nil, err
		}

		ps = append(ps, pps...)
	}

	return &caseTextRenamer{ps}, nil
}

func (r *caseTextRenamer) Rename(s string) string {
	for _, p := range r.patterns {
		s = p.Replace(s)
	}

	return s
}
