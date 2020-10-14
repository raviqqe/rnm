package main

import (
	"fmt"
	"os"

	"github.com/go-git/go-billy/v5/osfs"
	"github.com/logrusorgru/aurora/v3"
	"github.com/mattn/go-colorable"
)

func main() {
	stdout := colorable.NewColorableStdout()
	stderr := colorable.NewColorableStderr()

	fs := osfs.New(".")

	err := newCommand(
		newPathFinder(newRepositoryPathFinder(fs), fs),
		newFileRenamer(fs, os.Stderr),
		fs,
		stdout,
		stderr,
	).Run(os.Args[1:])
	if err != nil {
		fmt.Fprintln(stderr, aurora.Red(err))
		os.Exit(1)
	}
}
