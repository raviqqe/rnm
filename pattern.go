package main

import (
	"regexp"

	"github.com/iancoleman/strcase"
)

type caseConfig struct {
	name    caseName
	convert func(string) string
	head    delimiter
	tail    delimiter
}

type pattern struct {
	From *regexp.Regexp
	To   string
}

var caseConfigs = []*caseConfig{
	{
		camel,
		strcase.ToLowerCamel,
		nonAlphabet,
		upperCase,
	},
	{
		upperCamel,
		strcase.ToCamel,
		lowerCase,
		upperCase,
	},
	{
		kebab,
		strcase.ToKebab,
		nonAlphabet,
		nonAlphabet,
	},
	{
		upperKebab,
		strcase.ToScreamingKebab,
		nonAlphabet,
		nonAlphabet,
	},
	{
		snake,
		strcase.ToSnake,
		nonAlphabet,
		nonAlphabet,
	},
	{
		upperSnake,
		strcase.ToScreamingSnake,
		nonAlphabet,
		nonAlphabet,
	},
	{
		space,
		func(s string) string {
			return strcase.ToDelimited(s, ' ')
		},
		nonAlphabet,
		nonAlphabet,
	},
	{
		upperSpace,
		func(s string) string {
			return strcase.ToScreamingDelimited(s, ' ', 0, true)
		},
		nonAlphabet,
		nonAlphabet,
	},
}

func compilePatterns(from string, to string, enabled map[caseName]struct{}) ([]*pattern, error) {
	ps := make([]*pattern, 0, len(caseConfigs))

	for _, m := range caseConfigs {
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

func compilePattern(from string, to string, o *caseConfig) (*pattern, error) {
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
