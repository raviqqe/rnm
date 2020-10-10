package main

import (
	"path/filepath"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/filesystem"
)

type gitPathFinder struct {
	fileSystem       billy.Filesystem
	workingDirectory string
}

func newGitPathFinder(fs billy.Filesystem, workingDirectory string) *gitPathFinder {
	return &gitPathFinder{fs, workingDirectory}
}

func (f *gitPathFinder) Find() ([]string, error) {
	d := f.findGitRoot()
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
		if err == nil && !filepath.HasPrefix(s, "../") {
			ss = append(ss, s)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return ss, nil
}

func (f *gitPathFinder) findGitRoot() string {
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
