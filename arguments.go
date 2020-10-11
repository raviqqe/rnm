package main

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/jessevdk/go-flags"
)

const usage = "[options] <from> <to>"

type arguments struct {
	Bare         bool   `short:"b" long:"bare" description:"Use given patterns as they are"`
	RawCaseNames string `short:"c" long:"cases" description:"Comma-separated names of enabled cases (options: camel, upper-camel, kebab, upper-kebab, snake, upper-snake, space, upper-space)"`
	Help         bool   `short:"h" long:"help" description:"Show this help"`
	Version      bool   `long:"version" description:"Show version"`
	From         string
	To           string
	Path         string
	CaseNames    map[caseName]struct{}
}

func getArguments(ss []string) (*arguments, error) {
	args := arguments{}
	p := flags.NewParser(&args, flags.PassDoubleDash)
	p.Usage = usage

	ss, err := p.ParseArgs(ss)
	if err != nil {
		return nil, err
	} else if args.Help || args.Version {
		return &args, nil
	} else if len(ss) < 2 || len(ss) > 3 {
		return nil, errors.New("invalid number of arguments")
	}

	args.From, args.To = ss[0], ss[1]

	if len(ss) == 3 {
		args.Path = ss[2]
	} else {
		args.Path = "."
	}

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

func help() string {
	p := flags.NewParser(&arguments{}, flags.PassDoubleDash)
	p.Usage = usage

	// Parse() is run here to show default values in help.
	// This seems to be a bug in go-flags.
	p.Parse() // nolint:errcheck

	b := &bytes.Buffer{}
	p.WriteHelp(b)
	return b.String()
}
