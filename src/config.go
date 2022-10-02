package main

import (
	"fmt"
	"os"
	"regexp"
)

type Config struct {
	pattern *regexp.Regexp
	patternFunc func(string) bool
	program string
	programArgs []string
	showProgramOutput bool
	test bool
}

func NewConfig() *Config {
	config = &Config {
		pattern: nil,
		patternFunc: func(_ string) bool {
			return true
		},
		program: "",
		programArgs: []string{},
		showProgramOutput: false,
		test: false,
	}

	if len(os.Args) > 1 {
		inputPattern := ""
		remainingArgs := []string{}
		for _, arg := range os.Args[1:] {
			// Read switches in any order
			switch arg {
			case "--help":
				PrintHelp()
				os.Exit(0)
			case "--show-output":
				config.showProgramOutput = true
			case "--line":
				inputPattern = "^.*$"
			case "--git-commit-hash":
				inputPattern = "\\b[0-9a-f]{7,40}\\b"
			case "--filename":
				inputPattern = "[^\\s:]+"
				config.patternFunc = func(p string) bool {
					stat, err := os.Stat(p)
					return err == nil && !stat.IsDir()
				}
			case "--dirname":
				inputPattern = "[^\\s:]+"
				config.patternFunc = func(p string) bool {
					stat, err := os.Stat(p)
					return err == nil && stat.IsDir()
				}
			default:
				remainingArgs = append(remainingArgs, arg)
			}
		}

		offset := 0
		if inputPattern == "" && len(remainingArgs) > 0 {
			// No pattern has been set by a keyword, so assume the first argument is the pattern
			inputPattern = remainingArgs[0]
			offset++
		}

		if inputPattern != "" {
			// Regex pattern
			pattern, err := regexp.Compile(inputPattern)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Invalid regular expression")
				os.Exit(1)
			}
			config.pattern = pattern
		}

		if len(remainingArgs) > offset {
			// Program name
			config.program = remainingArgs[offset]
		}

		if len(remainingArgs) > offset + 1 {
			// Program arguments
			config.programArgs = remainingArgs[offset+1:]
		}
	}

	if os.Getenv("LISST_TEST") != "" {
		config.test = true
	}

	return config
}

