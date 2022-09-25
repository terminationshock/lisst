package main

import (
	"fmt"
	"os"
	"os/exec"
)

func LaunchProgram(match string) {
	args := config.programArgs
	args = append(args, match)

	cmd := exec.Command(config.program, args...)

	// Re-attach stdin to /dev/tty because of pipe input
	stdin, err := os.Open("/dev/tty")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	cmd.Stdin = stdin

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
