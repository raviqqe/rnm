package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"sync"

	"github.com/logrusorgru/aurora/v3"
	"gopkg.in/src-d/go-billy.v4"
)

type command struct {
	fileSystem billy.Filesystem
	stdout     io.Writer
	stderr     io.Writer
}

func newCommand(fileSystem billy.Filesystem, stdout, stderr io.Writer) *command {
	return &command{fileSystem, stdout, stderr}
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

	r, err := newRenamer(args.From, args.To, args.CaseNames)
	if err != nil {
		return err
	}

	ss, err = c.glob()
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

			err = c.rename(r, s)
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

func (c *command) glob() ([]string, error) {
	fs := []string{}
	ds := []string{"."}

	for len(ds) != 0 {
		d := ds[0]
		ds = ds[1:]

		is, err := c.fileSystem.ReadDir(d)
		if err != nil {
			return nil, err
		}

		for _, i := range is {
			p := c.fileSystem.Join(d, i.Name())

			i, err := c.fileSystem.Lstat(p)
			if err != nil {
				return nil, err
			} else if i.IsDir() {
				ds = append(ds, p)
			} else {
				fs = append(fs, p)
			}
		}
	}

	return fs, nil
}

func (c *command) rename(r *renamer, path string) error {
	ok, err := c.validatePath(path)
	if err != nil {
		return err
	} else if !ok {
		return nil
	}

	p := r.Rename(path)

	if p != path {
		err := c.fileSystem.Rename(path, p)
		if err != nil {
			return err
		}
	}

	i, err := c.fileSystem.Lstat(p)
	if err != nil {
		return err
	} else if i.IsDir() {
		return nil
	}

	ok, err = c.isTextFile(p)
	if err != nil {
		return err
	} else if !ok {
		return nil
	}

	f, err := c.fileSystem.OpenFile(p, os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}

	defer f.Close()

	bs, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	bbs := []byte(r.Rename(string(bs)))
	if bytes.Equal(bs, bbs) {
		return nil
	}

	_, err = f.Seek(0, 0)
	if err != nil {
		return err
	}

	_, err = f.Write(bbs)
	return err
}

func (*command) validatePath(s string) (bool, error) {
	ok, err := regexp.MatchString("(^|/)\\.", s)
	if err != nil {
		return false, err
	}

	return !ok, nil
}

func (c *command) isTextFile(path string) (bool, error) {
	f, err := c.fileSystem.Open(path)
	if err != nil {
		return false, err
	}

	// Read only 512 bytes for file type detection.
	bs := make([]byte, 0, 512)
	_, err = f.Read(bs)
	if err != nil && err != io.EOF {
		return false, err
	}

	return regexp.MatchString("^text/", http.DetectContentType(bs))
}

func (c *command) printError(err error) {
	fmt.Fprintln(c.stderr, aurora.Red(err))
}
