package main

import (
	"os"
)

func main() {
	err := newCommand().Run(os.Args[1:])
	if err != nil {
		printError(err)
		os.Exit(1)
	}
}
