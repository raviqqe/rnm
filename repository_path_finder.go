package main

import (
	"path/filepath"
	"regexp"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/filesystem"
)

var parentDirectoryRegexp *regexp.Regexp = regexp.MustCompile(`^\.\./`)

type repositoryPathFinder struct {
	fileSystem billy.Filesystem
}

func newRepositoryPathFinder(fs billy.Filesystem) *repositoryPathFinder {
	return &repositoryPathFinder{fs}
}

func (f *repositoryPathFinder) Find(d string) ([]string, error) {
	rd := f.findRepositoryRoot(d)
	if rd == "" {
		return nil, nil
	}

	gd, err := f.fileSystem.Chroot(f.fileSystem.Join(rd, ".git"))
	if err != nil {
		return nil, err
	}

	wd, err := f.fileSystem.Chroot(rd)
	if err != nil {
		return nil, err
	}

	r, err := git.Open(
		filesystem.NewStorage(gd, cache.NewObjectLRUDefault()),
		wd,
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

	ps := []string{}

	err = i.ForEach(func(file *object.File) error {
		ps = append(ps, f.fileSystem.Join(rd, file.Name))

		return nil
	})
	if err != nil {
		return nil, err
	}

	w, err := r.Worktree()
	if err != nil {
		return nil, err
	}

	st, err := w.Status()
	if err != nil {
		return nil, err
	}

	for p := range st {
		ps = append(ps, p)
	}

	pps := []string{}

	for _, p := range ps {
		p, err := filepath.Rel(d, p)
		if err == nil && !parentDirectoryRegexp.MatchString(filepath.ToSlash(p)) {
			pps = append(pps, p)
		}
	}

	return pps, nil
}

func (f *repositoryPathFinder) findRepositoryRoot(d string) string {
	for {
		_, err := f.fileSystem.Lstat(f.fileSystem.Join(d, ".git"))
		if err == nil {
			return d
		} else if err == billy.ErrCrossedBoundary {
			return ""
		}

		d = f.fileSystem.Join(d, "..")
	}
}
