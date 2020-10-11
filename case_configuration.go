package main

import "github.com/iancoleman/strcase"

type caseConfiguration struct {
	name    caseName
	convert func(string) string
	head    delimiter
	tail    delimiter
}

var caseConfigurations = []*caseConfiguration{
	{
		upperCamel,
		strcase.ToCamel,
		lowerCase,
		upperCase,
	},
	{
		camel,
		strcase.ToLowerCamel,
		nonAlphabet,
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
