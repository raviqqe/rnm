package main

import "gopkg.in/src-d/go-billy.v4"

type fileGlobber struct{ fileSystem billy.Filesystem }

func newFileGlobber(fs billy.Filesystem) *fileGlobber {
	return &fileGlobber{fs}
}

func (g *fileGlobber) Glob(d string) ([]string, error) {
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
