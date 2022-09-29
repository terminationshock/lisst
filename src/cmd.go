package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func LaunchProgram(match string) (string, string) {
	args := prepareArguments(match)
	cmd := exec.Command(config.program, args...)

	// Try re-attaching stdin to /dev/tty because of pipe input
	stdin, err := os.Open("/dev/tty")
	if err == nil {
		cmd.Stdin = stdin
		defer stdin.Close()
	} else {
		cmd.Stdin = os.Stdin
	}

	var buffer bytes.Buffer
	if config.showProgramOutput {
		cmd.Stdout = &buffer
	} else {
		cmd.Stdout = os.Stdout
	}
	cmd.Stderr = os.Stderr

	err = cmd.Run()

	if err != nil {
		// Try to forward the command's exit code
		exitError, ok := err.(*exec.ExitError)
		if ok {
			os.Exit(exitError.ExitCode())
		} else {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	return PrintCommand(match), buffer.String()
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
