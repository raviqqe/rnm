package main

import (
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/filesystem"
)

type repositoryFileFinder struct {
	fileSystem billy.Filesystem
}

func newRepositoryFileFinder(fs billy.Filesystem) *repositoryFileFinder {
	return &repositoryFileFinder{fs}
}

func (f *repositoryFileFinder) Find(d string, ignoreUntracked bool) ([]string, error) {
	d = f.findWorktreeDirectory(d)
	if d == "" {
		return nil, nil
	}

	gfs, err := f.fileSystem.Chroot(f.fileSystem.Join(d, ".git"))
	if err != nil {
		return nil, err
	}

	wfs, err := f.fileSystem.Chroot(d)
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
		ps = append(ps, f.fileSystem.Join(d, file.Name))

		return nil
	})
	if err != nil {
		return nil, err
	}

	if !ignoreUntracked {
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
	}

	return ps, nil
}

func (f *repositoryFileFinder) findWorktreeDirectory(d string) string {
	for {
		i, err := f.fileSystem.Lstat(f.fileSystem.Join(d, ".git"))
		if err == nil && i.IsDir() {
			return d
		} else if err == billy.ErrCrossedBoundary {
			return ""
		}

		d = f.fileSystem.Join(d, "..")
	}
}
