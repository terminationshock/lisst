package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func LaunchProgram(match string) {
	args := prepareArguments(match)
	cmd := exec.Command(config.program, args...)

	// Try re-attaching stdin to /dev/tty because of pipe input
	stdin, err := os.Open("/dev/tty")
	if err == nil {
		cmd.Stdin = stdin
	} else {
		cmd.Stdin = os.Stdin
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	stdin.Close()

	if err != nil {
		// Try to forward the command's exit code
		exitError, ok := err.(*exec.ExitError)
		if ok {
			os.Exit(exitError.ExitCode())
		} else {
			os.Exit(1)
		}
	}
}

func PrintCommand(match string) string {
	if match == "" {
		return ""
	}
	args := prepareArguments(match)
	return fmt.Sprintf("%s %s", config.program, strings.Join(args, " "))
}

func prepareArguments(match string) []string {
	args := config.programArgs
	return append(args, match)
}
