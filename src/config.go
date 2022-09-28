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
	color string
	test bool
}

func NewConfig() *Config {
	config = &Config {
		pattern: nil,
		program: "",
		programArgs: []string{},
		color: "red",
		test: false,
	}

	if len(os.Args) > 1 {
		if os.Args[1] == "--help" {
			PrintHelp()
			os.Exit(0)
		}

		inputPattern := os.Args[1]

		// Define human-readable keywords for frequently used patterns
		var shortcuts = map[string]string{
			"--line": "^.*$",
			"--grep-filename": "^(.*?):",
			"--git-commit-hash": "\\b[0-9a-f]{7,40}\\b",
		}
		value, ok := shortcuts[os.Args[1]]
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
	}

	if len(os.Args) > 2 {
		// Program name
		config.program = os.Args[2]
	}

	if len(os.Args) > 3 {
		// Program arguments
		config.programArgs = os.Args[3:]
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

