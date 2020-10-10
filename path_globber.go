package main

import "gopkg.in/src-d/go-billy.v4"

type pathGlobber struct{ fileSystem billy.Filesystem }

func newPathGlobber(fs billy.Filesystem) *pathGlobber {
	return &pathGlobber{fs}
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

	return ps, nil
}
