package main

import (
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/go-git/go-billy/v5"
	"github.com/logrusorgru/aurora/v3"
)

type command struct {
	pathGlobber *pathGlobber
	fileRenamer *fileRenamer
	fileSystem  billy.Filesystem
	stdout      io.Writer
	stderr      io.Writer
}

func newCommand(g *pathGlobber, r *fileRenamer, fs billy.Filesystem, stdout, stderr io.Writer) *command {
	return &command{g, r, fs, stdout, stderr}
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

	r, err := newTextRenamer(args.From, args.To, args.CaseNames)
	if err != nil {
		return err
	}

	i, err := c.fileSystem.Lstat(args.Path)
	if err != nil {
		return err
	} else if !i.IsDir() {
		return c.fileRenamer.Rename(r, args.Path)
	}

	ss, err = c.pathGlobber.Glob(args.Path)
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

			err = c.fileRenamer.Rename(r, s)
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

func (c *command) printError(err error) {
	fmt.Fprintln(c.stderr, aurora.Red(err))
}
