package main

import (
	"fmt"
	"path/filepath"
	"regexp"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/filesystem"
)

var parentDirectoryRegexp = regexp.MustCompile(`^\.\./`)

type repositoryFileFinder struct {
	fileSystem billy.Filesystem
}

func newRepositoryFileFinder(fs billy.Filesystem) *repositoryFileFinder {
	return &repositoryFileFinder{fs}
}

func (f *repositoryFileFinder) Find(d string) ([]string, error) {
	wd, err := f.findWorktreeDirectory(d)
	if err != nil {
		return nil, err
	} else if wd == "" {
		return nil, nil
	}

	gfs, err := f.fileSystem.Chroot(f.fileSystem.Join(wd, ".git"))
	if err != nil {
		return nil, err
	}

	wfs, err := f.fileSystem.Chroot(wd)
	if err != nil {
		return nil, err
	}

	r, err := git.Open(
		filesystem.NewStorage(gfs, cache.NewObjectLRUDefault()),
		wfs,
	)
	if err != nil {
		return nil, err
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
		ps = append(ps, f.fileSystem.Join(wd, file.Name))

		return nil
	})
	if err != nil {
		return nil, err
	}

	pps := make([]string, 0, len(ps))

	for _, p := range ps {
		b, err := filepath.Rel(d, p)
		if err == nil && !parentDirectoryRegexp.MatchString(filepath.ToSlash(b)) {
			pps = append(pps, p)
		}
	}

	return pps, nil
}

func (f *repositoryFileFinder) findWorktreeDirectory(d string) (string, error) {
	for {
		p := f.fileSystem.Join(d, ".git")
		i, err := f.fileSystem.Lstat(p)
		if err == nil && i.IsDir() {
			return d, nil
		} else if err == nil && !i.IsDir() {
			return "", fmt.Errorf("multiple worktrees not supported: %v", p)
		} else if err == billy.ErrCrossedBoundary || d == filepath.Dir(d) {
			return "", nil
		}

		d = filepath.Dir(d)
	}
}
