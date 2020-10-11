package main

import (
	"fmt"
	"os"

	"github.com/go-git/go-billy/v5/osfs"
	"github.com/logrusorgru/aurora/v3"
)

func main() {
	fs := osfs.New(".")

	err := newCommand(
		newPathGlobber(newRepositoryPathFinder(fs, "."), fs),
		newFileRenamer(fs),
		os.Stdout,
		os.Stderr,
	).Run(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, aurora.Red(err))
		os.Exit(1)
	}
}
