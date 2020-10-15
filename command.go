package main

import (
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"sync"

	"github.com/go-git/go-billy/v5"
	"github.com/logrusorgru/aurora/v3"
)

type command struct {
	pathFinder       *pathFinder
	fileRenamer      *fileRenamer
	fileSystem       billy.Filesystem
	workingDirectory string
	stdout           io.Writer
	stderr           io.Writer
}

func newCommand(g *pathFinder, r *fileRenamer, fs billy.Filesystem, d string, stdout, stderr io.Writer) *command {
	return &command{g, r, fs, d, stdout, stderr}
}

func (c *command) Run(ss []string) error {
	args, err := getArguments(ss)
	if err != nil {
		return err
	} else if args.Help {
		fmt.Fprint(c.stdout, help())
		return nil
	} else if args.Version {
		fmt.Fprintln(c.stdout, version)
		return nil
	}

	p := c.resolvePath(args.Path)
	r := newBareTextRenamer(args.From, args.To)

	if !args.Bare {
		r, err = newCaseTextRenamer(args.From, args.To, args.CaseNames)
		if err != nil {
			return err
		}
	}

	i, err := c.fileSystem.Stat(p)
	if err != nil {
		return err
	} else if !i.IsDir() {
		// Rename only filenames but not their directories.
		return c.fileRenamer.Rename(r, p, filepath.Dir(p), args.Verbose)
	}

	ss, err = c.pathFinder.Find(p)
	if err != nil {
		return err
	}

	g := &sync.WaitGroup{}
	sm := newSemaphore(maxOpenFiles)
	ec := make(chan error, errorChannelCapacity)

	for _, s := range ss {
		g.Add(1)
		go func(s string) {
			defer g.Done()

			sm.Request()
			defer sm.Release()

			err = c.fileRenamer.Rename(r, s, p, args.Verbose)
			if err != nil {
				ec <- fmt.Errorf("%v: %v", s, err)
			}
		}(s)
	}

	go func() {
		g.Wait()

		close(ec)
	}()

	ok := true

	for err := range ec {
		ok = false

		c.printError(err)
	}

	if !ok {
		return errors.New("failed to rename some identifiers")
	}

	return nil
}

func (c *command) resolvePath(p string) string {
	if p == "" {
		return c.workingDirectory
	} else if filepath.IsAbs(p) {
		return p
	}

	return filepath.Join(c.workingDirectory, p)
}

func (c *command) printError(err error) {
	fmt.Fprintln(c.stderr, aurora.Red(err))
}
