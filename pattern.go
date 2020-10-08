package main

import (
	"regexp"
)

type pattern struct {
	From *regexp.Regexp
	To   string
}

func compilePatterns(from string, to string, cs map[caseName]struct{}) ([]*pattern, error) {
	ps := make([]*pattern, 0, len(caseConfigurations))

	for _, c := range caseConfigurations {
		if _, ok := cs[c.name]; ok {
			p, err := compilePattern(from, to, c)
			if err != nil {
				return nil, err
			}

			ps = append(ps, p)
		}
	}

	return ps, nil
}

func compilePattern(from string, to string, c *caseConfiguration) (*pattern, error) {
	r, err := regexp.Compile(
		compileDelimiter(c.head, true) +
			c.convert(from) +
			compileDelimiter(c.tail, false),
	)
	if err != nil {
		return nil, err
	}

	return &pattern{r, "${1}" + c.convert(to) + "${2}"}, nil
}
