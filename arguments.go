package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/jessevdk/go-flags"
)

type arguments struct {
	RawCaseNames string `short:"c" long:"cases" description:"Comma-separated names of enabled cases (options: camel, upper-camel, kebab, upper-kebab, snake, upper-snake, space, upper-space)"`
	Help         bool   `short:"h" long:"help" description:"Show this help"`
	Version      bool   `long:"version" description:"Show version"`
	From         string
	To           string
	CaseNames    map[caseName]struct{}
}

func getArguments() (*arguments, error) {
	args := arguments{}
	p := flags.NewParser(&args, flags.PassDoubleDash)
	p.Usage = "[options] <from> <to>"

	ss, err := p.Parse()
	if err != nil {
		return nil, err
	} else if args.Help {
		p.WriteHelp(os.Stderr)
		os.Exit(0)
	} else if args.Version {
		fmt.Println(version)
		os.Exit(0)
	} else if len(ss) != 2 {
		return nil, errors.New("invalid number of arguments")
	}

	args.From, args.To = ss[0], ss[1]

	if args.RawCaseNames != "" {
		args.CaseNames = map[caseName]struct{}{}

		for _, n := range strings.Split(args.RawCaseNames, ",") {
			n := caseName(n)

			if _, ok := allCaseNames[n]; !ok {
				return nil, fmt.Errorf("invalid case name: %v", n)
			}

			args.CaseNames[n] = struct{}{}
		}
	}

	return &args, nil
}
