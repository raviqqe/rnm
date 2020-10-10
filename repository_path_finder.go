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

	ss := []string{}

	err = i.ForEach(func(file *object.File) error {
		s, err := filepath.Rel(f.workingDirectory, f.fileSystem.Join(d, file.Name))
		if err != nil {
			// Ignore errors assuming that they are not in the current directory.
			return nil
		}

		ok, err := regexp.MatchString(`^\.\./`, s)
		if err != nil {
			return err
		} else if !ok {
			ss = append(ss, s)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return ss, nil
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
