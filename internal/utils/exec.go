package utils

import (
	"os/exec"
	"strings"
)

// runCommand executes a command and returns its combined output.
// This is a simple wrapper around os/exec for use in utility functions.
func runCommand(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), err
	}
	return string(out), nil
}

// splitLines splits a string into lines, filtering empty ones.
func splitLines(s string) []string {
	var lines []string
	for _, line := range strings.Split(s, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			lines = append(lines, trimmed)
		}
	}
	return lines
}

// split splits a string by whitespace.
func split(s string) []string {
	return strings.Fields(s)
}
