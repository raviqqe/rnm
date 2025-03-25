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
	argumentParser *argumentParser
	fileFinder     *fileFinder
	fileRenamer    *fileRenamer
	fileSystem     billy.Filesystem
	stdout         io.Writer
	stderr         io.Writer
}

func newCommand(p *argumentParser, g *fileFinder, r *fileRenamer, fs billy.Filesystem, stdout, stderr io.Writer) *command {
	return &command{p, g, r, fs, stdout, stderr}
}

func (c *command) Run(ss []string) error {
	args, err := c.argumentParser.Parse(ss)
	if err != nil {
		return err
	} else if args.Help {
		_, err := fmt.Fprint(c.stdout, c.argumentParser.Help())
		return err
	} else if args.Version {
		_, err := fmt.Fprintln(c.stdout, version)
		return err
	}

	r, err := newCaseTextRenamer(args.From, args.To, args.CaseNames)
	if err != nil {
		return err
	}

	if args.Bare {
		r = newBareTextRenamer(args.From, args.To)
	} else if args.Regexp {
		r, err = newRegexpTextRenamer(args.From, args.To)
		if err != nil {
			return err
		}
	}

	i, err := c.fileSystem.Stat(args.Path)
	if err != nil {
		return err
	} else if !i.IsDir() {
		// Rename only filenames but not their directories.
		return c.fileRenamer.Rename(
			r,
			args.Path,
			filepath.Dir(args.Path),
			args.Verbose,
		)
	}

	ss, err = c.fileFinder.Find(
		args.Path,
		args.Include,
		args.Exclude,
		args.NoGit,
	)
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

			err = c.fileRenamer.Rename(r, s, args.Path, args.Verbose)
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

		if err := c.printError(err); err != nil {
			return err
		}
	}

	if !ok {
		return errors.New("failed to rename some identifiers")
	}

	return nil
}

func (c *command) printError(err error) error {
	_, err = fmt.Fprintln(c.stderr, aurora.Red(err))
	return err
}
