package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

type Config struct {
	pattern *regexp.Regexp
	patternFunc func(string) bool
	program string
	programArgs []string
	filter bool
	sort int
	showProgramOutput bool
	ignoreProgramError bool
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
		filter: false,
		sort: 0,
		showProgramOutput: false,
		ignoreProgramError: false,
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
			case "--filter":
				config.filter = true
			case "--sort":
				config.sort = 1
			case "--sort-rev":
				config.sort = -1
			case "--show-output":
				config.showProgramOutput = true
			case "--ignore-error":
				config.ignoreProgramError = true
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
			case "--filename-lineno":
				inputPattern = "[^\\s:]+:[1-9][0-9]*"
				config.patternFunc = func(p string) bool {
					splitted := strings.Split(p, ":")
					if len(splitted) == 0 {
						return false
					}
					filename := strings.Join(splitted[:len(splitted) - 1], ":")
					stat, err := os.Stat(filename)
					return err == nil && !stat.IsDir()
				}
			case "--dirname":
				inputPattern = "[^\\s:]+"
				config.patternFunc = func(p string) bool {
					stat, err := os.Stat(p)
					return err == nil && stat.IsDir()
				}
			default:
				if strings.HasPrefix(arg, "--") {
					fmt.Fprintln(os.Stderr, "Invalid command-line option " + arg)
					os.Exit(1)
				}
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

