package main

import "github.com/go-git/go-billy/v5"

type pathGlobber struct {
	repositoryPathFinder *repositoryPathFinder
	fileSystem           billy.Filesystem
}

func newPathGlobber(f *repositoryPathFinder, fs billy.Filesystem) *pathGlobber {
	return &pathGlobber{f, fs}
}

func (g *pathGlobber) Glob(d string) ([]string, error) {
	ps := []string{}
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
			ps = append(ps, p)

			i, err := g.fileSystem.Lstat(p)
			if err != nil {
				return nil, err
			} else if i.IsDir() {
				ds = append(ds, p)
			}
		}
	}

	rps, err := g.repositoryPathFinder.Find()
	if err != nil {
		return nil, err
	} else if len(rps) == 0 {
		return ps, nil
	}

	return intersectStringSets(ps, rps), nil
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
