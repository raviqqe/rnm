package main

import (
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/filesystem"
)

type repositoryPathFinder struct {
	fileSystem billy.Filesystem
}

func newRepositoryPathFinder(fs billy.Filesystem) *repositoryPathFinder {
	return &repositoryPathFinder{fs}
}

func (f *repositoryPathFinder) Find(d string) ([]string, error) {
	d = f.findRepositoryRoot(d)
	if d == "" {
		return nil, nil
	}

	gd, err := f.fileSystem.Chroot(f.fileSystem.Join(d, ".git"))
	if err != nil {
		return nil, err
	}

	fs, err := f.fileSystem.Chroot(d)
	if err != nil {
		return nil, err
	}

	r, err := git.Open(
		filesystem.NewStorage(gd, cache.NewObjectLRUDefault()),
		fs,
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
		ps = append(ps, f.fileSystem.Join(d, p))
	}

	return ps, nil
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
