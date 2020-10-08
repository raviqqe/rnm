package main

import (
	"regexp"

	"github.com/iancoleman/strcase"
)

type patternOptions struct {
	convert func(string) string
	head    delimiter
	tail    delimiter
}

type pattern struct {
	From *regexp.Regexp
	To   string
}

var patternOptionsList = []*patternOptions{
	{
		func(s string) string {
			return strcase.ToDelimited(s, ' ')
		},
		nonAlphabet,
		nonAlphabet,
	},
	{
		func(s string) string {
			return strcase.ToScreamingDelimited(s, ' ', 0, true)
		},
		nonAlphabet,
		nonAlphabet,
	},
	{
		strcase.ToKebab,
		nonAlphabet,
		nonAlphabet,
	},
	{
		strcase.ToScreamingKebab,
		nonAlphabet,
		nonAlphabet,
	},
	{
		strcase.ToSnake,
		nonAlphabet,
		nonAlphabet,
	},
	{
		strcase.ToScreamingSnake,
		nonAlphabet,
		nonAlphabet,
	},
	{
		strcase.ToCamel,
		lowerCase,
		upperCase,
	},
	{
		strcase.ToLowerCamel,
		upperCase,
		upperCase,
	},
}

func compilePatterns(from string, to string) ([]*pattern, error) {
	ps := make([]*pattern, 0, len(patternOptionsList))

	for _, o := range patternOptionsList {
		p, err := compilePattern(from, to, o)
		if err != nil {
			return nil, err
		}

		ps = append(ps, p)
	}

	return ps, nil
}

func compilePattern(from string, to string, o *patternOptions) (*pattern, error) {
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
