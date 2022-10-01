package main

import (
	"testing"
)

func TestPrintCommand(t *testing.T) {
	config = &Config{}
	config.program = "program"
	config.programArgs = []string{"arg1", "arg2"}

	cmd := PrintCommand("match")
	if cmd != "program arg1 arg2 match" {
		t.Error("Incorrect command string")
	}

	cmd = PrintCommand("")
	if cmd != "" {
		t.Error("Incorrect empty command string")
	}

	config.programArgs = []string{"arg1", "{}", "{}{}", "arg{}3", "arg4"}

	cmd = PrintCommand("match")
	if cmd != "program arg1 match matchmatch argmatch3 arg4" {
		t.Error("Incorrect inserted command string")
	}
}
