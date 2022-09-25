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
	debug bool
}

func NewConfig(input []string) *Config {
	if len(os.Args) > 1 && os.Args[1] == "--help" {
		fmt.Println("Usage: COMMAND_IN | " + os.Args[0] + " PATTERN [COMMAND]")
		fmt.Println("\nThis program displays the output of COMMAND_IN as interactive list")
		fmt.Println("and opens COMMAND after a certain match has been selected.")
		fmt.Println("   PATTERN    - Accepts regular expressions")
		fmt.Println("              - Case insensitive (can be toggled with key [c])")
		fmt.Println("   DIRECTORY  - Recursive search within this directory")
		fmt.Println("              - Defaults to the current working directory")
		fmt.Println("\nKey bindings:")
		fmt.Println("   [q] or [Esc]      Quit")
		fmt.Println("   [Up] and [Down]   Select a match")
		fmt.Println("   [Enter]           Open the selected file with the editor (see variable EDITOR)")
		fmt.Println("\nEnvironment variables:")
		fmt.Println("   LISST_COLOR        - Highlight color for matches")
		fmt.Println("                       - Assign \"-\" to disable highlighting")
		fmt.Println("                       - Default: \"red\"")
		os.Exit(0)
	}

	config = &Config {
		pattern: nil,
		program: "",
		programArgs: []string{},
		color: "#ff0000",
		debug: false,
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

	if os.Getenv("LISST_DEBUG") != "" {
		config.debug = true
	}

	return config
}

