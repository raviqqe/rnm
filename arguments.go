package main

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
)

type arguments struct {
	Patterns struct {
		From string
		To   string
	} `positional-args:"true"`
	Version bool `long:"version" description:"Show version"`
}

func getArguments() (*arguments, error) {
	args := arguments{}
	_, err := flags.Parse(&args)
	if err != nil {
		return nil, err
	} else if args.Version {
		fmt.Println(version)
		os.Exit(0)
	}

	return &args, nil
}
