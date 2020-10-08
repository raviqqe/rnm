package main

type delimiter int

const (
	nonAlphabet delimiter = iota
	upperCase
	lowerCase
)

func compileDelimiter(d delimiter, head bool) string {
	s := "[^[:alpha:]]"

	switch d {
	case upperCase:
		s = "(?:[[:upper:]]|" + s + ")"
	case lowerCase:
		s = "(?:[[:lower:]]|" + s + ")"
	}

	return compilePosition(s, head)
}

func compilePosition(pattern string, head bool) string {
	s := "$"

	if head {
		s = "^"
	}

	return "(" + pattern + "|" + s + ")"
}
