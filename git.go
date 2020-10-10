package main

import (
	"os"
	"path"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func listGitFiles() ([]string, error) {
	r, err := git.PlainOpenWithOptions(
		".",
		&git.PlainOpenOptions{DetectDotGit: true},
	)
	if err != nil {
		return nil, nil
	}

	ref, err := r.Head()
	if err != nil {
		return nil, err
	}

	c, err := r.CommitObject(ref.Hash())
	if err != nil {
		return nil, err
	}

	i, err := c.Files()
	if err != nil {
		return nil, err
	}

	d, err := findGitRoot()
	if err != nil {
		return nil, err
	}

	ss := []string{}

	err = i.ForEach(func(f *object.File) error {
		ss = append(ss, path.Join(d, f.Name))
		return nil
	})
	if err != nil {
		return nil, err
	}

	return ss, nil
}

func findGitRoot() (string, error) {
	s, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		i, err := os.Lstat(path.Join(s, ".git"))
		if err != nil {
			return "", err
		} else if i.IsDir() || s == "" {
			return s, nil
		}

		s = path.Dir(s)

	}

}
