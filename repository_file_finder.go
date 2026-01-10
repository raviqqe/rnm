package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-billy/v6"
	"github.com/go-git/go-billy/v6/util"
	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing/cache"
	"github.com/go-git/go-git/v6/storage/filesystem"
)

type repositoryFileFinder struct {
	fileSystem billy.Filesystem
}

func newRepositoryFileFinder(fs billy.Filesystem) *repositoryFileFinder {
	return &repositoryFileFinder{fs}
}

func (f *repositoryFileFinder) Find(d string) ([]string, error) {
	r, wd, err := f.openGitRepository(d)
	if err != nil {
		return nil, err
	} else if r == nil {
		return nil, nil
	}

	i, err := r.Storer.Index()
	if err != nil {
		return nil, err
	}

	ps := []string{}

	for _, e := range i.Entries {
		p := f.fileSystem.Join(wd, e.Name)

		i, err := f.fileSystem.Lstat(p)
		if err != nil {
			return nil, err
		} else if i.IsDir() {
			// Directories are meaningful only if they contain files in Git repositories.
			// Hence, file renaming handles directories implicitly.
			continue
		}

		b, err := filepath.Rel(d, p)

		if err == nil && !strings.HasPrefix(filepath.ToSlash(b), "../") {
			ps = append(ps, p)
		}
	}

	return ps, nil
}

func (f *repositoryFileFinder) openGitRepository(d string) (*git.Repository, string, error) {
	rd, i := f.findWorktreeDirectory(d)
	if rd == "" {
		return nil, "", nil
	}

	wd := filepath.Dir(rd)
	wfs, err := f.fileSystem.Chroot(wd)
	if err != nil {
		return nil, "", err
	}

	if !i.IsDir() {
		rd, err = f.findWorktreeDataDirectory(rd)
		if err != nil {
			return nil, "", err
		}
	}

	rfs, err := f.fileSystem.Chroot(rd)
	if err != nil {
		return nil, "", err
	}

	r, err := git.Open(
		filesystem.NewStorage(rfs, cache.NewObjectLRUDefault()),
		wfs,
	)
	if err != nil {
		return nil, "", err
	}

	return r, wd, nil
}

func (f *repositoryFileFinder) findWorktreeDirectory(d string) (string, os.FileInfo) {
	for {
		p := f.fileSystem.Join(d, ".git")

		if i, err := f.fileSystem.Lstat(p); err == nil {
			return p, i
		} else if err == billy.ErrCrossedBoundary || d == filepath.Dir(d) {
			return "", nil
		}

		d = filepath.Dir(d)
	}
}

func (f *repositoryFileFinder) findWorktreeDataDirectory(p string) (string, error) {
	bs, err := util.ReadFile(f.fileSystem, p)
	if err != nil {
		return "", err
	}

	s := strings.Split(string(bs), "\n")[0]
	const prefix = "gitdir:"

	if !strings.HasPrefix(s, prefix) {
		return "", fmt.Errorf("no gitdir entry in .git file: %v", p)
	}

	return f.resolvePath(filepath.Dir(p), strings.TrimSpace(s[len(prefix):])), nil
}

func (f *repositoryFileFinder) resolvePath(d, p string) string {
	if filepath.IsAbs(p) {
		return p
	}

	return filepath.Clean(f.fileSystem.Join(d, p))
}
