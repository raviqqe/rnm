package rename

import (
	"regexp"

	"github.com/iancoleman/strcase"
)

// Renamer renames all identifiers similar to a given identifier.
type Renamer struct {
	patterns []pattern
}

// New creates a renamer.
func New(from string, to string) (*Renamer, error) {
	ps := []pattern{}

	for _, f := range [](func(string) string){strcase.ToCamel} {
		r, err := regexp.Compile(f(from))
		if err != nil {
			return nil, err
		}

		ps = append(ps, pattern{r, f(to)})
	}

	return &Renamer{ps}, nil
}

// Rename renames all identifiers in a string.
func (r *Renamer) Rename(s string) string {
	for _, p := range r.patterns {
		s = p.From.ReplaceAllLiteralString(s, p.To)
	}

	return s
}
