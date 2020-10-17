package main

import (
	"testing"
	"time"

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
		Author: &object.Signature{
			Name:  "",
			Email: "",
			When:  time.Now(),
		},
	})
	assert.Nil(t, err)
}

func TestRepositoryFileFinderFindNoPath(t *testing.T) {
	fs := memfs.New()
	commitFiles(t, fs, nil)

	ss, err := newRepositoryFileFinder(fs).Find(".", false)
	assert.Nil(t, err)
	assert.Equal(t, []string{}, ss)
}

func TestRepositoryFileFinderFindCommittedPath(t *testing.T) {
	fs := memfs.New()
	commitFiles(t, fs, []string{"foo"})

	ss, err := newRepositoryFileFinder(fs).Find(".", false)
	assert.Nil(t, err)
	assert.Equal(t, []string{"foo"}, ss)
}

func TestRepositoryFileFinderFindUncommittedPath(t *testing.T) {
	fs := memfs.New()
	commitFiles(t, fs, nil)

	_, err := fs.Create("foo")
	assert.Nil(t, err)

	ss, err := newRepositoryFileFinder(fs).Find(".", false)
	assert.Nil(t, err)
	assert.Equal(t, []string{"foo"}, ss)
}

func TestRepositoryFileFinderDoNotFindUncommittedPath(t *testing.T) {
	fs := memfs.New()
	commitFiles(t, fs, nil)

	_, err := fs.Create("foo")
	assert.Nil(t, err)

	ss, err := newRepositoryFileFinder(fs).Find(".", true)
	assert.Nil(t, err)
	assert.Equal(t, []string{}, ss)
}

func TestRepositoryFileFinderDoNotFindIgnoredUncommittedPath(t *testing.T) {
	fs := memfs.New()
	commitFiles(t, fs, nil)

	err := util.WriteFile(fs, ".gitignore", []byte("foo\n"), 0o444)
	assert.Nil(t, err)

	_, err = fs.Create("foo")
	assert.Nil(t, err)

	ss, err := newRepositoryFileFinder(fs).Find(".", false)
	assert.Nil(t, err)
	assert.Equal(t, []string{".gitignore"}, ss)
}

func TestRepositoryFileFinderFindPathInsideDirectory(t *testing.T) {
	fs := memfs.New()

	err := fs.MkdirAll("bar", 0o755)
	assert.Nil(t, err)

	commitFiles(t, fs, []string{"bar/foo"})

	ss, err := newRepositoryFileFinder(fs).Find("bar", false)
	assert.Nil(t, err)
	assert.Equal(t, []string{"bar/foo"}, normalizePaths(ss))
}

func TestRepositoryFileFinderFindPathOutsideDirectory(t *testing.T) {
	fs := memfs.New()

	err := fs.MkdirAll("bar", 0o755)
	assert.Nil(t, err)

	commitFiles(t, fs, []string{"foo"})

	ss, err := newRepositoryFileFinder(fs).Find("bar", false)
	assert.Nil(t, err)
	assert.Equal(t, []string{}, normalizePaths(ss))
}

func TestRepositoryFileFinderFindUncommittedPathInsideDirectory(t *testing.T) {
	fs := memfs.New()

	err := fs.MkdirAll("bar", 0o755)
	assert.Nil(t, err)

	commitFiles(t, fs, nil)

	_, err = fs.Create("bar/foo")
	assert.Nil(t, err)

	ss, err := newRepositoryFileFinder(fs).Find("bar", false)
	assert.Nil(t, err)
	assert.Equal(t, []string{"bar/foo"}, normalizePaths(ss))
}

func TestRepositoryFileFinderFindPathInDirectory(t *testing.T) {
	fs := memfs.New()

	commitFiles(t, fs, nil)

	err := fs.MkdirAll("foo", 0o755)
	assert.Nil(t, err)

	_, err = fs.Create("foo/foo")
	assert.Nil(t, err)

	_, err = fs.Create("bar")
	assert.Nil(t, err)

	ss, err := newRepositoryFileFinder(fs).Find("foo", false)
	assert.Nil(t, err)
	assert.Equal(t, []string{"foo/foo"}, normalizePaths(ss))
}

// TODO Support multiple worktrees of the same repositories.
func TestRepositoryFileFinderDoNotFindPathInDifferentWorktree(t *testing.T) {
	fs := memfs.New()

	err := fs.MkdirAll("foo", 0o755)
	assert.Nil(t, err)

	subfs, err := fs.Chroot("foo")
	assert.Nil(t, err)
	commitFiles(t, subfs, nil)

	err = fs.MkdirAll("bar", 0o755)
	assert.Nil(t, err)

	err = util.WriteFile(
		fs,
		"bar/.git",
		[]byte("gitdir: /foo/.git"),
		0o444,
	)
	assert.Nil(t, err)

	_, err = fs.Create("bar/foo")
	assert.Nil(t, err)

	ss, err := newRepositoryFileFinder(fs).Find("bar", false)
	assert.Nil(t, err)
	assert.Equal(t, []string(nil), normalizePaths(ss))
}
