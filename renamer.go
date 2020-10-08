package main

type renamer struct {
	patterns []*pattern
}

func newRenamer(from string, to string) (*renamer, error) {
	ps, err := compilePatterns(from, to)
	if err != nil {
		return nil, err
	}

	return &renamer{ps}, nil
}

func (r *renamer) Rename(s string) string {
	for _, p := range r.patterns {
		s = p.From.ReplaceAllString(s, p.To)
	}

	return s
}
