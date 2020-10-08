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
	RawCaseNames []patternName `long:"enable" description:"Enable only specified cases (options: camel, upper-camel, kebab, upper-kebab, snake, upper-snake, space, upper-space)"`
	Help         bool          `short:"h" long:"help" description:"Show this help"`
	Version      bool          `long:"version" description:"Show version"`
	CaseNames    map[patternName]struct{}
}

func getArguments() (*arguments, error) {
	args := arguments{}
	p := flags.NewParser(&args, flags.PassDoubleDash)

	_, err := p.Parse()
	if err != nil {
		return nil, err
	} else if args.Help {
		p.WriteHelp(os.Stderr)
		os.Exit(0)
	} else if args.Version {
		fmt.Println(version)
		os.Exit(0)
	} else if args.RawCaseNames != nil {
		args.CaseNames = map[patternName]struct{}{}

		for _, n := range args.RawCaseNames {
			args.CaseNames[n] = struct{}{}
		}
	}

	return &args, nil
}
