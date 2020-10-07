package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"github.com/mattn/go-zglob"
	"github.com/raviqqe/rnm/rename"
)

func main() {
	err := command()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func command() error {
	if len(os.Args[1:]) != 2 {
		return errors.New("Usage: rnm <from> <to>")
	}

	r, err := rename.New(os.Args[1], os.Args[2])
	if err != nil {
		return err
	}

	ss, err := zglob.Glob("**/*")
	if err != nil {
		return err
	}

	g := &sync.WaitGroup{}
	ec := make(chan error, 1024)

	for _, s := range ss {
		g.Add(1)
		go func(s string) {
			defer g.Done()

			err := renameFile(r, s)
			if err != nil {
				ec <- err
			}
		}(s)
	}

	g.Wait()

	ok := false

	for err := range ec {
		ok = true

		fmt.Fprintln(os.Stderr, err)
	}

	if !ok {
		return errors.New("failed to rename some identifiers")
	}

	return nil
}

func renameFile(r *rename.Renamer, path string) error {
	p := r.Rename(path)

	err := os.Rename(path, p)
	if err != nil {
		return err
	}

	i, err := os.Lstat(p)
	if err != nil {
		return err
	} else if i.IsDir() {
		return nil
	}

	bs, err := ioutil.ReadFile(p)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(p, []byte(r.Rename(string(bs))), i.Mode())
}
