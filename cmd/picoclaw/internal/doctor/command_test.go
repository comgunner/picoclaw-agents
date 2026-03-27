// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package doctor_test

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/comgunner/picoclaw/cmd/picoclaw/internal/doctor"
	"github.com/comgunner/picoclaw/pkg/setup"
)

// TestDoctorCommand_Runs verifies the doctor command executes without panic.
func TestDoctorCommand_Runs(t *testing.T) {
	cmd := doctor.NewDoctorCommand()

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := cmd.Execute()

	// Restore stdout
	w.Close()
	os.Stdout = oldStdout
	_, _ = io.ReadAll(r)

	assert.NoError(t, err, "Doctor command should execute without error")
}

// TestDoctorCommand_JSONFlag verifies --json produces valid JSON.
func TestDoctorCommand_JSONFlag(t *testing.T) {
	cmd := doctor.NewDoctorCommand()
	cmd.SetArgs([]string{"--json"})

	// Capture stdout
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	err := cmd.Execute()
	require.NoError(t, err)

	output := buf.String()

	// Verify it's valid JSON by unmarshaling into EnvironmentReport
	var report setup.EnvironmentReport
	err = json.Unmarshal([]byte(output), &report)
	assert.NoError(t, err, "Output should be valid JSON")

	// Verify key fields are populated
	assert.NotEmpty(t, report.OS)
	assert.NotEmpty(t, report.Arch)
}

// TestDoctorCommand_DefaultOutput verifies default tabular output format.
func TestDoctorCommand_DefaultOutput(t *testing.T) {
	cmd := doctor.NewDoctorCommand()
	cmd.SetArgs([]string{})

	// Capture stdout
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	err := cmd.Execute()
	require.NoError(t, err)

	output := buf.String()

	// Verify expected sections in output
	assert.Contains(t, output, "=== PicoClaw-Agents Doctor ===")
	assert.Contains(t, output, "System:")
	assert.Contains(t, output, "Requirements:")
	assert.Contains(t, output, "Security:")
	assert.Contains(t, output, "OS/Arch:")
	assert.Contains(t, output, "WSL:")
	assert.Contains(t, output, "Go:")
	assert.Contains(t, output, "Docker:")
	assert.Contains(t, output, "Workspace:")
	assert.Contains(t, output, "Root:")
	assert.Contains(t, output, "Dangerous:")
	assert.Contains(t, output, "Open ports:")
}

// TestDoctorCommand_ExitCode verifies exit code is 0 on success.
func TestDoctorCommand_ExitCode(t *testing.T) {
	cmd := doctor.NewDoctorCommand()
	cmd.SetArgs([]string{})

	// Capture stdout
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	err := cmd.Execute()

	// Should return nil (exit code 0) in normal environment
	assert.NoError(t, err)
}

// TestRunDoctor_JSONOutput verifies JSON output contains all fields.
func TestRunDoctor_JSONOutput(t *testing.T) {
	cmd := doctor.NewDoctorCommand()
	cmd.SetArgs([]string{"--json"})

	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := cmd.Execute()
	require.NoError(t, err)

	var report setup.EnvironmentReport
	err = json.Unmarshal(buf.Bytes(), &report)
	require.NoError(t, err)

	// Verify all expected fields are present
	assert.NotEmpty(t, report.OS)
	assert.NotEmpty(t, report.Arch)
	assert.NotEmpty(t, report.GoVersion)
	assert.NotEmpty(t, report.WorkspacePath)
	assert.IsType(t, true, report.WSL)
	assert.IsType(t, setup.SecurityReport{}, report.SecurityChecks)
}

// TestRunDoctor_TabularOutput verifies tabular output formatting.
func TestRunDoctor_TabularOutput(t *testing.T) {
	cmd := doctor.NewDoctorCommand()

	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := cmd.Execute()
	require.NoError(t, err)

	output := buf.String()

	// Verify formatting
	lines := strings.Split(output, "\n")
	assert.Greater(t, len(lines), 5, "Output should have multiple lines")

	// Verify sections are in expected order
	systemIdx := strings.Index(output, "System:")
	requirementsIdx := strings.Index(output, "Requirements:")
	securityIdx := strings.Index(output, "Security:")

	assert.Less(t, systemIdx, requirementsIdx, "System should come before Requirements")
	assert.Less(t, requirementsIdx, securityIdx, "Requirements should come before Security")
}

// TestStatusIcon verifies status icon helper function.
func TestStatusIcon(t *testing.T) {
	// This tests the internal statusIcon function indirectly through output
	cmd := doctor.NewDoctorCommand()
	cmd.SetArgs([]string{})

	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := cmd.Execute()
	require.NoError(t, err)

	output := buf.String()

	// Should contain either [OK] or [MISSING]
	assert.True(t, strings.Contains(output, "[OK]") || strings.Contains(output, "[MISSING]"))
}

// TestDoctorCommand_Help verifies help text is available.
func TestDoctorCommand_Help(t *testing.T) {
	cmd := doctor.NewDoctorCommand()

	help := cmd.Help()
	assert.NoError(t, help)

	assert.Equal(t, "doctor", cmd.Use)
	assert.Contains(t, cmd.Short, "Check environment")
	assert.Contains(t, cmd.Long, "Diagnose environment")
}

// TestDoctorCommand_NoArgs verifies command rejects arguments.
func TestDoctorCommand_NoArgs(t *testing.T) {
	cmd := doctor.NewDoctorCommand()
	cmd.SetArgs([]string{"unexpected-arg"})

	err := cmd.Execute()

	// Should error because command expects no arguments
	assert.Error(t, err)
}
