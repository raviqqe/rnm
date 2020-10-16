package main

import "github.com/go-git/go-billy/v5"

type fileFinder struct {
	repositoryPathFinder *repositoryPathFinder
	fileSystem           billy.Filesystem
}

func newFileFinder(f *repositoryPathFinder, fs billy.Filesystem) *fileFinder {
	return &fileFinder{f, fs}
}

func (g *fileFinder) Find(d string, ignoreUntracked bool) ([]string, error) {
	fs := []string{}
	ds := []string{d}

	for len(ds) != 0 {
		d := ds[0]
		ds = ds[1:]

		is, err := g.fileSystem.ReadDir(d)
		if err != nil {
			return nil, err
		}

		for _, i := range is {
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

	rps, err := g.repositoryPathFinder.Find(d, ignoreUntracked)
	if err != nil {
		return nil, err
	} else if len(rps) == 0 {
		return fs, nil
	}

	return intersectStringSets(fs, rps), nil
}

func intersectStringSets(ss, sss []string) []string {
	sm := make(map[string]struct{}, len(ss))

	for _, s := range ss {
		sm[s] = struct{}{}
	}

	ss = make([]string, 0, len(sm))

	for _, s := range sss {
		if _, ok := sm[s]; ok {
			ss = append(ss, s)
		}
	}

	return ss
}
