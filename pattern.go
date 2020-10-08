package main

import (
	"regexp"
)

type pattern struct {
	From *regexp.Regexp
	To   string
}

func compilePatterns(from string, to string, enabled map[caseName]struct{}) ([]*pattern, error) {
	ps := make([]*pattern, 0, len(caseConfigurations))

	for _, m := range caseConfigurations {
		if _, ok := enabled[m.name]; ok {
			p, err := compilePattern(from, to, m)
			if err != nil {
				return nil, err
			}

			ps = append(ps, p)
		}
	}

	return ps, nil
}

func compilePattern(from string, to string, o *caseConfiguration) (*pattern, error) {
	r, err := regexp.Compile(
		compileDelimiter(o.head, true) +
			o.convert(from) +
			compileDelimiter(o.tail, false),
	)
	if err != nil {
		return nil, err
	}

	return &pattern{r, "${1}" + o.convert(to) + "${2}"}, nil
}
