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
	fileSystem       billy.Filesystem
	workingDirectory string
}

func newRepositoryPathFinder(fs billy.Filesystem, workingDirectory string) *repositoryPathFinder {
	return &repositoryPathFinder{fs, workingDirectory}
}

func (f *repositoryPathFinder) Find() ([]string, error) {
	d := f.findRepositoryRoot()
	if d == "" {
		return nil, nil
	}

	gd, err := f.fileSystem.Chroot(f.fileSystem.Join(d, ".git"))
	if err != nil {
		return nil, err
	}

	wd, err := f.fileSystem.Chroot(d)
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
		ps = append(ps, f.fileSystem.Join(d, file.Name))

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
		p, err := filepath.Rel(f.workingDirectory, p)
		if err == nil && !parentDirectoryRegexp.MatchString(filepath.ToSlash(p)) {
			pps = append(pps, p)
		}
	}

	return pps, nil
}

func (f *repositoryPathFinder) findRepositoryRoot() string {
	s := f.workingDirectory

	for {
		_, err := f.fileSystem.Lstat(f.fileSystem.Join(s, ".git"))
		if err == nil {
			return s
		} else if err == billy.ErrCrossedBoundary {
			return ""
		}

		s = f.fileSystem.Join(s, "..")
	}
}
