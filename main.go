package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sync"

	"github.com/logrusorgru/aurora/v3"
	"github.com/mattn/go-zglob"
)

func main() {
	err := run()
	if err != nil {
		printError(err)
		os.Exit(1)
	}
}

func run() error {
	args, err := getArguments()
	if err != nil {
		return err
	}

	r, err := newRenamer(args.From, args.To, args.CaseNames)
	if err != nil {
		return err
	}

	ss, err := getPaths()
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

			err = renameFile(r, s)
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

		printError(err)
	}

	if !ok {
		return errors.New("failed to rename some identifiers")
	}

	return nil
}

func renameFile(r *renamer, path string) error {
	ok, err := validatePath(path)
	if err != nil {
		return err
	} else if !ok {
		return nil
	}

	p := r.Rename(path)

	if p != path {
		err := os.Rename(path, p)
		if err != nil {
			return err
		}
	}

	i, err := os.Lstat(p)
	if err != nil {
		return err
	} else if i.IsDir() {
		return nil
	}

	ok, err = isTextFile(p)
	if err != nil {
		return err
	} else if !ok {
		return nil
	}

	bs, err := ioutil.ReadFile(p)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(p, []byte(r.Rename(string(bs))), i.Mode())
}

func validatePath(s string) (bool, error) {
	ok, err := regexp.MatchString("(^|/)\\.", s)
	if err != nil {
		return false, err
	}

	return !ok, nil
}

func isTextFile(path string) (bool, error) {
	f, err := os.Open(path)
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

func getPaths() ([]string, error) {
	ss, err := glob()
	if err != nil {
		return nil, err
	}

	sss, err := listGitFiles()
	if err != nil {
		return nil, err
	} else if len(sss) == 0 {
		return ss, nil
	}

	return intersectStringSets(ss, sss), nil
}

func glob() ([]string, error) {
	ss, err := zglob.Glob("**/*")
	if err != nil {
		return nil, err
	}

	sss := make([]string, 0, len(ss))

	for _, s := range ss {
		s, err = filepath.Abs(s)
		if err != nil {
			return nil, err
		}

		sss = append(sss, s)
	}

	return sss, nil
}

func intersectStringSets(ss, sss []string) []string {
	sm := map[string]struct{}{}

	for _, s := range ss {
		sm[s] = struct{}{}
	}

	ss = make([]string, 0, len(sm))

	for _, s := range sss {
		if _, ok := sm[s]; ok {
			ss = append(ss, s)
		}
	}

	return ss
}

func printError(err error) {
	fmt.Fprintln(os.Stderr, aurora.Red(err))
}
