package main

import (
	"regexp"

	"github.com/iancoleman/strcase"
)

type renamer struct {
	patterns []compiledPattern
}

func newRenamer(from string, to string) (*renamer, error) {
	ps := []compiledPattern{}

	for _, f := range [](func(string) string){
		func(s string) string {
			return strcase.ToDelimited(s, ' ')
		},
		func(s string) string {
			return strcase.ToScreamingDelimited(s, ' ', 0, true)
		},
		strcase.ToCamel,
		strcase.ToKebab,
		strcase.ToLowerCamel,
		strcase.ToScreamingKebab,
		strcase.ToScreamingSnake,
		strcase.ToSnake,
	} {
		r, err := regexp.Compile(f(from))
		if err != nil {
			return nil, err
		}

		ps = append(ps, compiledPattern{r, f(to)})
	}

	return &renamer{ps}, nil
}

func (r *renamer) Rename(s string) string {
	for _, p := range r.patterns {
		s = p.From.ReplaceAllLiteralString(s, p.To)
	}

	return s
}
