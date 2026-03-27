// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package cli_test

import (
	"bufio"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/comgunner/picoclaw/pkg/cli"
)

// TestReadLine verifies basic line reading.
func TestReadLine(t *testing.T) {
	input := "test input\n"
	scanner := bufio.NewScanner(strings.NewReader(input))

	result := cli.ReadLine(scanner)
	assert.Equal(t, "test input", result)
}

// TestReadLine_Empty verifies empty input handling.
func TestReadLine_Empty(t *testing.T) {
	scanner := bufio.NewScanner(strings.NewReader(""))

	result := cli.ReadLine(scanner)
	assert.Equal(t, "", result)
}

// TestConfirm_Yes verifies yes confirmation.
func TestConfirm_Yes(t *testing.T) {
	input := "y\n"
	scanner := bufio.NewScanner(strings.NewReader(input))

	result := cli.Confirm("Are you sure?", scanner)
	assert.True(t, result)
}

// TestConfirm_No verifies no confirmation.
func TestConfirm_No(t *testing.T) {
	input := "n\n"
	scanner := bufio.NewScanner(strings.NewReader(input))

	result := cli.Confirm("Are you sure?", scanner)
	assert.False(t, result)
}

// TestConfirm_DefaultNo verifies default is no.
func TestConfirm_DefaultNo(t *testing.T) {
	input := "\n"
	scanner := bufio.NewScanner(strings.NewReader(input))

	result := cli.Confirm("Are you sure?", scanner)
	assert.False(t, result)
}

// TestConfirm_CapitalY verifies capital Y works.
func TestConfirm_CapitalY(t *testing.T) {
	input := "Y\n"
	scanner := bufio.NewScanner(strings.NewReader(input))

	result := cli.Confirm("Are you sure?", scanner)
	assert.True(t, result)
}

// TestReadMaskedWithFallback_NonTTY verifies fallback works in non-TTY.
func TestReadMaskedWithFallback_NonTTY(t *testing.T) {
	// In tests, we're non-TTY so it should use fallback
	input := "secret-key-123\n"
	scanner := bufio.NewScanner(strings.NewReader(input))

	result, err := cli.ReadMaskedWithFallback("Enter API key: ", scanner)
	assert.NoError(t, err)
	assert.Equal(t, "secret-key-123", result)
}

// TestReadMaskedWithFallback_Empty verifies empty input handling.
func TestReadMaskedWithFallback_Empty(t *testing.T) {
	scanner := bufio.NewScanner(strings.NewReader(""))

	result, err := cli.ReadMaskedWithFallback("Enter API key: ", scanner)
	assert.Error(t, err)
	assert.Empty(t, result)
}

// TestReadLine_Whitespace verifies whitespace trimming.
func TestReadLine_Whitespace(t *testing.T) {
	input := "  trimmed  \n"
	scanner := bufio.NewScanner(strings.NewReader(input))

	result := cli.ReadLine(scanner)
	assert.Equal(t, "trimmed", result)
}

// TestConfirm_Whitespace verifies whitespace handling in confirm.
func TestConfirm_Whitespace(t *testing.T) {
	input := "  y  \n"
	scanner := bufio.NewScanner(strings.NewReader(input))

	result := cli.Confirm("Are you sure?", scanner)
	// Should trim and match
	assert.True(t, result) // "y" after trimming
}
