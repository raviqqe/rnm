package main

import (
	"fmt"
	"io"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/util"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/filesystem"
	"github.com/go-git/go-git/v5/storage/filesystem/dotgit"
)

var parentDirectoryRegexp = regexp.MustCompile(`^\.\./`)

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

func (f repositoryFileFinder) openGitRepository(d string) (*git.Repository, string, error) {
	rd, i := f.findWorktreeDirectory(d)
	if rd == "" {
		return nil, "", nil
	}

	wfs, err := f.fileSystem.Chroot(filepath.Dir(rd))
	if err != nil {
		return nil, "", err
	}

	if !i.IsDir() {
		rd, err := f.findRepositoryDirectory(rd)
		if err != nil {
			return nil, err
		}
	}

	rfs, err := f.fileSystem.Chroot(rd)
	if err != nil {
		return nil, "", err
	}

	return git.Open(
		filesystem.NewStorage(rfs, cache.NewObjectLRUDefault()),
		wfs,
	)
}

func (f *repositoryFileFinder) readGitDirFromDotGitFile(dotGitPath, worktreeDirectory string) (string, error) {
	file, err := f.fileSystem.Open(dotGitPath)
	if err != nil {
		return "", err
	}
	defer func() { _ = file.Close() }()

	b, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	line := strings.Split(string(b), "\n")[0]
	const prefix = "gitdir:"

	if !strings.HasPrefix(line, prefix) {
		return "", fmt.Errorf(".git file has no %s prefix: %v", prefix, dotGitPath)
	}

	d := strings.TrimSpace(line[len(prefix):])

	if filepath.IsAbs(d) {
		return filepath.Clean(d), nil
	}

	return filepath.Clean(filepath.Join(worktreeDirectory, d)), nil
}

func (f *repositoryFileFinder) findWorktreeDirectory(d string) (string, os.FileInfo) {
	for {
		p := f.fileSystem.Join(d, ".git")

		if i, err := f.fileSystem.Lstat(p); err == nil {
			return p, i
		} else if err == billy.ErrCrossedBoundary || d == filepath.Dir(d) {
			return "", os.FileInfo{}
		}

		d = filepath.Dir(d)
	}
}

func (f *repositoryFileFinder) findCommonDirectory(d string) (string, error) {
	p := f.fileSystem.Join(d, "commondir")
	bs, err := util.ReadFile(f.fileSystem, p)
	if err != nil {
		return "", err
	}

	cd := strings.TrimSpace(string(bs))
	if cd == "" {
		return "", fmt.Errorf("empty commondir file: %v", p)
	}

	if filepath.IsAbs(cd) {
		return cd, nil
	}

	return filepath.Clean(filepath.Join(d, cd)), nil
}
