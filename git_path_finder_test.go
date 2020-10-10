package main

import (
	"testing"
	"time"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
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
			Name:  "John Doe",
			Email: "john@doe.org",
			When:  time.Now(),
		},
	})
	assert.Nil(t, err)
}

func TestGitPathFinderFindNoPath(t *testing.T) {
	fs := memfs.New()
	commitFiles(t, fs, nil)

	ss, err := newGitPathFinder(fs, ".").Find()
	assert.Nil(t, err)
	assert.Equal(t, []string{}, ss)
}

func TestGitPathFinderFindPath(t *testing.T) {
	fs := memfs.New()
	commitFiles(t, fs, []string{"foo"})

	ss, err := newGitPathFinder(fs, ".").Find()
	assert.Nil(t, err)
	assert.Equal(t, []string{"foo"}, ss)
}

func TestGitPathFinderFindPathInsideDirectory(t *testing.T) {
	fs := memfs.New()

	err := fs.MkdirAll("bar", 0o755)
	assert.Nil(t, err)

	commitFiles(t, fs, []string{"bar/foo"})

	ss, err := newGitPathFinder(fs, "bar").Find()
	assert.Nil(t, err)
	assert.Equal(t, []string{"foo"}, ss)
}

func TestGitPathFinderFindPathOutsideDirectory(t *testing.T) {
	fs := memfs.New()

	err := fs.MkdirAll("bar", 0o755)
	assert.Nil(t, err)

	commitFiles(t, fs, []string{"foo"})

	ss, err := newGitPathFinder(fs, "bar").Find()
	assert.Nil(t, err)
	assert.Equal(t, []string{}, ss)
}
