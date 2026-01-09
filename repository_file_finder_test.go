package main

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-billy/v5/util"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/filesystem"
	"github.com/stretchr/testify/assert"
)

func commitFiles(t *testing.T, fs billy.Filesystem, paths []string) {
	g, err := fs.Chroot(".git")
	assert.Nil(t, err)

	r, err := git.Init(filesystem.NewStorage(g, cache.NewObjectLRUDefault()), fs)
	assert.Nil(t, err)

	w, err := r.Worktree()
	assert.Nil(t, err)

	for _, p := range paths {
		_, err = fs.Create(p)
		assert.Nil(t, err)

		_, err = w.Add(p)
		assert.Nil(t, err)
	}

	_, err = w.Commit("foo", &git.CommitOptions{
		AllowEmptyCommits: true,
		Author:            &object.Signature{},
	})
	assert.Nil(t, err)
}

func TestRepositoryFileFinderFindNoPath(t *testing.T) {
	fs := memfs.New()
	commitFiles(t, fs, nil)

	ss, err := newRepositoryFileFinder(fs).Find(".")
	assert.Nil(t, err)
	assert.Equal(t, []string{}, ss)
}

func TestRepositoryFileFinderFindCommittedPath(t *testing.T) {
	fs := memfs.New()
	commitFiles(t, fs, []string{"foo"})

	ss, err := newRepositoryFileFinder(fs).Find(".")
	assert.Nil(t, err)
	assert.Equal(t, []string{"foo"}, ss)
}

func TestRepositoryFileFinderDoNotFindUncommittedPath(t *testing.T) {
	fs := memfs.New()
	commitFiles(t, fs, nil)

	_, err := fs.Create("foo")
	assert.Nil(t, err)

	ss, err := newRepositoryFileFinder(fs).Find(".")
	assert.Nil(t, err)
	assert.Equal(t, []string{}, ss)
}

func TestRepositoryFileFinderDoNotFindIgnoredUncommittedPath(t *testing.T) {
	fs := memfs.New()
	commitFiles(t, fs, nil)

	err := util.WriteFile(fs, ".gitignore", []byte("foo\n"), 0o400)
	assert.Nil(t, err)

	_, err = fs.Create("foo")
	assert.Nil(t, err)

	ss, err := newRepositoryFileFinder(fs).Find(".")
	assert.Nil(t, err)
	assert.Equal(t, []string{}, ss)
}

func TestRepositoryFileFinderFindPathInsideDirectory(t *testing.T) {
	fs := memfs.New()

	err := fs.MkdirAll("bar", 0o700)
	assert.Nil(t, err)

	commitFiles(t, fs, []string{"bar/foo"})

	ss, err := newRepositoryFileFinder(fs).Find("bar")
	assert.Nil(t, err)
	assert.Equal(t, []string{"bar/foo"}, normalizePaths(ss))
}

func TestRepositoryFileFinderDoNotFindPathOutsideDirectory(t *testing.T) {
	fs := memfs.New()

	err := fs.MkdirAll("bar", 0o700)
	assert.Nil(t, err)

	commitFiles(t, fs, []string{"foo"})

	ss, err := newRepositoryFileFinder(fs).Find("bar")
	assert.Nil(t, err)
	assert.Equal(t, []string{}, normalizePaths(ss))
}

func TestRepositoryFileFinderDoNotFindUncommittedPathInsideDirectory(t *testing.T) {
	fs := memfs.New()

	commitFiles(t, fs, nil)

	err := fs.MkdirAll("bar", 0o700)
	assert.Nil(t, err)

	_, err = fs.Create("bar/foo")
	assert.Nil(t, err)

	ss, err := newRepositoryFileFinder(fs).Find("bar")
	assert.Nil(t, err)
	assert.Equal(t, []string{}, normalizePaths(ss))
}

func TestRepositoryFileFinderFindPathInDirectory(t *testing.T) {
	fs := memfs.New()

	err := fs.MkdirAll("foo", 0o700)
	assert.Nil(t, err)

	commitFiles(t, fs, []string{"foo/foo", "bar"})

	ss, err := newRepositoryFileFinder(fs).Find("foo")
	assert.Nil(t, err)
	assert.Equal(t, []string{"foo/foo"}, normalizePaths(ss))
}

// TODO Test this when go-git v6 is released.
// func TestRepositoryFileFinderFindPathInLinkedWorktree(t *testing.T) {
// 	fs := memfs.New()
//
// 	err := fs.MkdirAll("repository", 0o700)
// 	assert.Nil(t, err)
//
// 	rfs, err := fs.Chroot("repository")
// 	assert.Nil(t, err)
// 	commitFiles(t, rfs, []string{"foo"})
//
// 	err = fs.MkdirAll("worktree", 0o700)
// 	assert.Nil(t, err)
//
// 	err = fs.MkdirAll("repository/.git/worktrees/0123", 0o700)
// 	assert.Nil(t, err)
//
// 	err = util.WriteFile(
// 		fs,
// 		"worktree/.git",
// 		[]byte("gitdir: ../repository/.git/worktrees/0123\n"),
// 		0o400,
// 	)
// 	assert.Nil(t, err)
//
// 	_, err = fs.Create("worktree/foo")
// 	assert.Nil(t, err)
//
// 	ss, err := newRepositoryFileFinder(fs).Find("worktree")
// 	assert.Nil(t, err)
// 	assert.Equal(t, []string{"worktree/foo"}, normalizePaths(ss))
// }

// TODO Test this when go-git v6 is released.
// func TestRepositoryFileFinderFindPathInLinkedWorktreeWithAbsoluteGitPaths(t *testing.T) {
// 	if runtime.GOOS == "windows" {
// 		t.Skip()
// 	}
//
// 	fs := memfs.New()
//
// 	err := fs.MkdirAll("repository", 0o700)
// 	assert.Nil(t, err)
//
// 	rfs, err := fs.Chroot("repository")
// 	assert.Nil(t, err)
// 	commitFiles(t, rfs, []string{"foo"})
//
// 	err = fs.MkdirAll("worktree", 0o700)
// 	assert.Nil(t, err)
//
// 	err = fs.MkdirAll("repository/.git/worktrees/0123", 0o700)
// 	assert.Nil(t, err)
//
// 	err = util.WriteFile(
// 		fs,
// 		"worktree/.git",
// 		// An absolute `gitdir` path
// 		fmt.Appendf(nil, "gitdir: %v\n", filepath.FromSlash("/repository/.git/worktrees/0123")),
// 		0o400,
// 	)
// 	assert.Nil(t, err)
//
// 	_, err = fs.Create("worktree/foo")
// 	assert.Nil(t, err)
//
// 	ss, err := newRepositoryFileFinder(fs).Find("worktree")
// 	assert.Nil(t, err)
// 	assert.Equal(t, []string{"worktree/foo"}, normalizePaths(ss))
// }

func TestRepositoryFileFinderFindPathInLinkedWorktreeWithInvalidGitdirEntry(t *testing.T) {
	fs := memfs.New()

	err := fs.MkdirAll("repository", 0o700)
	assert.Nil(t, err)

	rfs, err := fs.Chroot("repository")
	assert.Nil(t, err)
	commitFiles(t, rfs, []string{"foo"})

	err = fs.MkdirAll("worktree", 0o700)
	assert.Nil(t, err)

	err = fs.MkdirAll("repository/.git/worktrees/0123", 0o700)
	assert.Nil(t, err)

	err = util.WriteFile(
		fs,
		"worktree/.git",
		[]byte("git_dir: ../repository/.git/worktrees/0123\n"),
		0o400,
	)
	assert.Nil(t, err)

	_, err = newRepositoryFileFinder(fs).Find("worktree")
	assert.EqualError(t, err, fmt.Sprintf("no gitdir entry in .git file: %v", filepath.FromSlash("worktree/.git")))
}
