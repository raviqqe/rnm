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
	pathFinder     *pathFinder
	fileRenamer    *fileRenamer
	fileSystem     billy.Filesystem
	stdout         io.Writer
	stderr         io.Writer
}

func newCommand(p *argumentParser, g *pathFinder, r *fileRenamer, fs billy.Filesystem, stdout, stderr io.Writer) *command {
	return &command{p, g, r, fs, stdout, stderr}
}

func (c *command) Run(ss []string) error {
	args, err := c.argumentParser.Parse(ss)
	if err != nil {
		return err
	} else if args.Help {
		fmt.Fprint(c.stdout, c.argumentParser.Help())
		return nil
	} else if args.Version {
		fmt.Fprintln(c.stdout, version)
		return nil
	}

	r := newBareTextRenamer(args.From, args.To)

	if !args.Bare {
		r, err = newCaseTextRenamer(args.From, args.To, args.CaseNames)
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

	ss, err = c.pathFinder.Find(args.Path, args.IgnoreUntracked)
	if err != nil {
		return err
	}

	fs := make([]string, 0, len(ss))

	// Rename all directories first to avoid concurrency bugs.
	// It assumes that a number of directories is quite small compared to a number of files.
	for _, s := range ss {
		i, err := c.fileSystem.Stat(s)
		if err != nil {
			return err
		} else if i.IsDir() {
			err = c.renameFile(r, s, args.Path, args.Verbose)
			if err != nil {
				return err
			}
		} else {
			fs = append(fs, s)
		}
	}

	g := &sync.WaitGroup{}
	sm := newSemaphore(maxOpenFiles)
	ec := make(chan error, errorChannelCapacity)

	for _, s := range fs {
		g.Add(1)
		go func(s string) {
			defer g.Done()

			sm.Request()
			defer sm.Release()

			err = c.renameFile(r, s, args.Path, args.Verbose)
			if err != nil {
				ec <- err
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

func (c *command) renameFile(tr textRenamer, path string, baseDir string, verbose bool) error {
	err := c.fileRenamer.Rename(tr, path, baseDir, verbose)
	if err != nil {
		return fmt.Errorf("%v: %v", path, err)
	}

	return nil
}

func (c *command) printError(err error) {
	fmt.Fprintln(c.stderr, aurora.Red(err))
}
