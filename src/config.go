package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

type Config struct {
	pattern *regexp.Regexp
	program string
	programArgs []string
	showProgramOutput bool
	color string
	test bool
}

func NewConfig() *Config {
	config = &Config {
		pattern: nil,
		program: "",
		programArgs: []string{},
		showProgramOutput: false,
		color: "red",
		test: false,
	}

	if len(os.Args) > 1 {
		argOffset := 1
		loop:
		for i, arg := range os.Args[1:] {
			// Read switches in any order
			switch arg {
			case "--help":
				PrintHelp()
				os.Exit(0)
			case "--pipe-less":
				config.showProgramOutput = true
			default:
				// No more switch found
				argOffset = i + 1
				break loop
			}
		}

		// Assume that the next argument is the pattern
		inputPattern := os.Args[argOffset]

		// Define human-readable keywords for frequently used patterns
		var shortcuts = map[string]string{
			"--line": "^.*$",
			"--grep-filename": "^(.*?):",
			"--git-commit-hash": "\\b[0-9a-f]{7,40}\\b",
		}
		value, ok := shortcuts[inputPattern]
		if ok {
			inputPattern = value
		}

		// Regex pattern
		pattern, err := regexp.Compile(inputPattern)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Invalid regular expression")
			os.Exit(1)
		}
		config.pattern = pattern

		if len(os.Args) > argOffset + 1 {
			// Program name
			config.program = os.Args[argOffset + 1]
		}

		if len(os.Args) > argOffset + 2 {
			// Program arguments
			config.programArgs = os.Args[argOffset + 2:]
		}
	}

	color := strings.TrimSpace(os.Getenv("LISST_COLOR"))
	if color != "" {
		// Highlighting color for matches
		config.color = color
	}

	if os.Getenv("LISST_TEST") != "" {
		config.test = true
	}

	return config
}

