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
	executed string
	color string
	test bool
}

func NewConfig() *Config {
	if len(os.Args) > 1 && os.Args[1] == "--help" {
		printHelp()
		os.Exit(0)
	}

	config = &Config {
		pattern: nil,
		program: "",
		programArgs: []string{},
		executed: "",
		color: "red",
		test: false,
	}

	if len(os.Args) > 1 {
		// Regex pattern
		pattern, err := regexp.Compile(os.Args[1])
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

