package main

import (
	"os"
)

func main() {
	err := run()
	if err != nil {
		printError(err)
		os.Exit(1)
	}
}
