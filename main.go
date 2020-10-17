package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/go-git/go-billy/v5/osfs"
	"github.com/logrusorgru/aurora/v3"
	"github.com/mattn/go-colorable"
)

func main() {
	stdout := colorable.NewColorableStdout()
	stderr := colorable.NewColorableStderr()

	fs := osfs.New(filepath.FromSlash("/"))

	d, err := os.Getwd()
	if err != nil {
		fail(stderr, err)
	}

	err = newCommand(
		newArgumentParser(d),
		newFileFinder(newRepositoryFileFinder(fs), fs),
		newFileRenamer(fs, os.Stderr),
		fs,
		stdout,
		stderr,
	).Run(os.Args[1:])
	if err != nil {
		fail(stderr, err)
	}
}

func fail(stderr io.Writer, err error) {
	fmt.Fprintln(stderr, aurora.Red(err))
	os.Exit(1)
}
