package main

import (
	"regexp"

	"github.com/go-git/go-billy/v5"
)

var hiddenPathRegexp = regexp.MustCompile(`^\.`)

type fileFinder struct {
	repositoryFileFinder *repositoryFileFinder
	fileSystem           billy.Filesystem
}

func newFileFinder(f *repositoryFileFinder, fs billy.Filesystem) *fileFinder {
	return &fileFinder{f, fs}
}

func (f *fileFinder) Find(d string, excludedPattern *regexp.Regexp, ignoreUntracked bool) ([]string, error) {
	fs, err := f.findFiles(d, ignoreUntracked)
	if err != nil {
		return nil, err
	} else if excludedPattern == nil {
		return fs, nil
	}

	ffs := make([]string, 0, len(fs))

	for _, f := range fs {
		if !excludedPattern.MatchString(f) {
			ffs = append(ffs, f)
		}
	}

	return ffs, nil
}

func (f *fileFinder) findFiles(d string, ignoreUntracked bool) ([]string, error) {
	fs, err := f.repositoryFileFinder.Find(d, ignoreUntracked)
	if err != nil {
		return nil, err
	} else if len(fs) != 0 {
		return fs, nil
	}

	return f.findFilesOutsideRepository(d)
}

func (f *fileFinder) findFilesOutsideRepository(d string) ([]string, error) {
	fs := []string{}
	ds := []string{d}

	for len(ds) != 0 {
		d := ds[0]
		ds = ds[1:]

		is, err := f.fileSystem.ReadDir(d)
		if err != nil {
			return nil, err
		}

		for _, i := range is {
			if hiddenPathRegexp.MatchString(i.Name()) {
				continue
			}

			p := f.fileSystem.Join(d, i.Name())

			i, err := f.fileSystem.Lstat(p)
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
