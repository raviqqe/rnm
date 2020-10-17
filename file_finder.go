package main

import (
	"regexp"

	"github.com/go-git/go-billy/v5"
)

var hiddenPathRegexp = regexp.MustCompile(`^\.`)

type fileFinder struct {
	repositoryPathFinder *repositoryPathFinder
	fileSystem           billy.Filesystem
}

func newFileFinder(f *repositoryPathFinder, fs billy.Filesystem) *fileFinder {
	return &fileFinder{f, fs}
}

func (g *fileFinder) Find(d string, ignoreUntracked bool) ([]string, error) {
	fs, err := g.repositoryPathFinder.Find(d, ignoreUntracked)
	if err != nil {
		return nil, err
	} else if len(fs) != 0 {
		return fs, nil
	}

	fs = []string{}
	ds := []string{d}

	for len(ds) != 0 {
		d := ds[0]
		ds = ds[1:]

		is, err := g.fileSystem.ReadDir(d)
		if err != nil {
			return nil, err
		}

		for _, i := range is {
			if hiddenPathRegexp.MatchString(i.Name()) {
				continue
			}

			p := g.fileSystem.Join(d, i.Name())

			i, err := g.fileSystem.Lstat(p)
			if err != nil {
				return nil, err
			} else if i.IsDir() {
				ds = append(ds, p)
			} else {
				fs = append(fs, p)
			}
		}
	}

	return fs, nil
}
