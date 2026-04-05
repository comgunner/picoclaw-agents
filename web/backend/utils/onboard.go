package utils

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/comgunner/picoclaw/pkg/config"
)

var execCommand = exec.Command

func EnsureOnboarded(configPath string) error {
	_, err := os.Stat(configPath)
	if err == nil {
		return nil // Config already exists, skip automatically
	}
	if !os.IsNotExist(err) {
		return fmt.Errorf("stat config: %w", err)
	}

	// Config doesn't exist, run onboard
	cmd := execCommand(FindPicoclawBinary(), "onboard")
	cmd.Env = append(os.Environ(), config.EnvConfig+"="+configPath)
	// Don't send stdin - let onboard handle skip automatically

	output, err := cmd.CombinedOutput()
	if err != nil {
		// Check if it's a skip message (config exists)
		outputStr := string(output)
		if strings.Contains(outputStr, "already exists") ||
			strings.Contains(outputStr, "Skipping onboard") {
			return nil // Skip was successful
		}

		trimmed := strings.TrimSpace(outputStr)
		if trimmed == "" {
			return fmt.Errorf("run onboard: %w", err)
		}
		return fmt.Errorf("run onboard: %w: %s", err, trimmed)
	}

	if _, err := os.Stat(configPath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("onboard completed but did not create config %s", configPath)
		}
		return fmt.Errorf("verify config after onboard: %w", err)
	}

	return nil
}
