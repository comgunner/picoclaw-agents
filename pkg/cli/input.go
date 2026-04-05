// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

// Package cli provides command-line input utilities.
package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

// ReadMasked reads input from terminal without echo (for API keys and passwords).
// Returns the entered string and an error if the terminal is non-interactive.
func ReadMasked(prompt string) (string, error) {
	fmt.Print(prompt + " ")
	b, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println() // newline after silent input
	if err != nil {
		return "", fmt.Errorf("reading masked input: %w", err)
	}
	return strings.TrimSpace(string(b)), nil
}

// ReadMaskedWithFallback uses ReadMasked if TTY is available, falls back to bufio.Scanner for tests/pipes.
// This ensures compatibility with both interactive terminals and non-TTY environments (CI, pipes).
func ReadMaskedWithFallback(prompt string, scanner *bufio.Scanner) (string, error) {
	if term.IsTerminal(int(os.Stdin.Fd())) {
		return ReadMasked(prompt)
	}
	// Fallback for non-TTY (tests, pipes)
	fmt.Print(prompt + " ")
	if scanner.Scan() {
		return strings.TrimSpace(scanner.Text()), nil
	}
	return "", fmt.Errorf("no input available")
}

// ReadLine reads a single line of input using the provided scanner.
// Returns the trimmed input string.
func ReadLine(scanner *bufio.Scanner) string {
	if scanner.Scan() {
		return strings.TrimSpace(scanner.Text())
	}
	return ""
}

// Confirm prompts the user for a yes/no confirmation.
// Returns true if the user enters 'y' or 'Y', false otherwise.
func Confirm(prompt string, scanner *bufio.Scanner) bool {
	fmt.Print(prompt + " [y/N]: ")
	if scanner.Scan() {
		input := strings.TrimSpace(scanner.Text())
		return input == "y" || input == "Y"
	}
	return false
}
