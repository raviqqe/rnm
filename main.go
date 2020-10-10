package main

import (
	"fmt"
	"os"

	"github.com/logrusorgru/aurora/v3"
	"gopkg.in/src-d/go-billy.v4/osfs"
)

func main() {
	err := newCommand(osfs.New("."), os.Stderr).Run(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, aurora.Red(err))
		os.Exit(1)
	}
}
