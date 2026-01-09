package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-git/go-billy/v5"
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
	wd, rd, err := f.findWorktreeDirectory(d)
	if err != nil {
		return nil, err
	} else if wd == "" {
		return nil, nil
	}

	rfs, err := f.fileSystem.Chroot(rd)
	if err != nil {
		return nil, err
	}

	wfs, err := f.fileSystem.Chroot(wd)
	if err != nil {
		return nil, err
	}

	rfs = billy.Filesystem(rfs)
	commonFs, err := f.findCommonGitDirectory(rd)
	if err != nil {
		return nil, err
	} else if commonFs != nil {
		rfs = dotgit.NewRepositoryFilesystem(rfs, commonFs)
	}

	r, err := git.Open(
		filesystem.NewStorage(rfs, cache.NewObjectLRUDefault()),
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

func (f *repositoryFileFinder) findWorktreeDirectory(d string) (string, string, error) {
	for {
		p := f.fileSystem.Join(d, ".git")
		i, err := f.fileSystem.Lstat(p)
		if err == nil && i.IsDir() {
			return d, p, nil
		} else if err == nil && !i.IsDir() {
			gitDir, err := f.readGitDirFromDotGitFile(p, d)
			return d, gitDir, err
		} else if err == billy.ErrCrossedBoundary || d == filepath.Dir(d) {
			return "", "", nil
		}

		d = filepath.Dir(d)
	}
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

func (f *repositoryFileFinder) findCommonGitDirectory(gitDir string) (billy.Filesystem, error) {
	file, err := f.fileSystem.Open(f.fileSystem.Join(gitDir, "commondir"))
	if os.IsNotExist(err) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	b, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	p := strings.TrimSpace(string(b))
	if p == "" {
		return nil, nil
	}

	if !filepath.IsAbs(p) {
		p = filepath.Clean(filepath.Join(gitDir, p))
	}

	if _, err := f.fileSystem.Stat(p); err != nil {
		return nil, err
	}

	return f.fileSystem.Chroot(p)
}
