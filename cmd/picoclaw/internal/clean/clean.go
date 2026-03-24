package clean

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

func NewCleanCommand() *cobra.Command {
	var (
		olderThan string
		dryRun    bool
		all       bool
	)

	cmd := &cobra.Command{
		Use:   "clean",
		Short: "Clean old or corrupt session files",
		Long: `Clean session files from the workspace.

Examples:
  picoclaw clean --all              # Remove all sessions
  picoclaw clean --older-than 7d    # Remove sessions older than 7 days
  picoclaw clean --dry-run          # Show what would be deleted`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runClean(olderThan, dryRun, all)
		},
	}

	cmd.Flags().StringVarP(&olderThan, "older-than", "o", "", "Remove sessions older than (e.g., 7d, 24h)")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what would be deleted without actually deleting")
	cmd.Flags().BoolVar(&all, "all", false, "Remove all sessions")

	return cmd
}

func runClean(olderThan string, dryRun bool, all bool) error {
	// Get workspace path
	home, _ := os.UserHomeDir()
	sessionsDir := filepath.Join(home, ".picoclaw", "workspace", "sessions")

	// Check if directory exists
	if _, err := os.Stat(sessionsDir); os.IsNotExist(err) {
		fmt.Println("No sessions directory found")
		return nil
	}

	// Parse older-than duration
	var cutoffTime time.Time
	if olderThan != "" {
		duration, err := time.ParseDuration(olderThan)
		if err != nil {
			return fmt.Errorf("invalid duration format: %w", err)
		}
		cutoffTime = time.Now().Add(-duration)
	}

	// Walk through sessions
	var deletedCount int
	var totalSize int64

	err := filepath.Walk(sessionsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// Check if should delete
		shouldDelete := false
		if all {
			shouldDelete = true
		} else if olderThan != "" && info.ModTime().Before(cutoffTime) {
			shouldDelete = true
		}

		if shouldDelete {
			if dryRun {
				fmt.Printf("Would delete: %s (%d bytes)\n", path, info.Size())
			} else {
				if err := os.Remove(path); err != nil {
					fmt.Printf("Error deleting %s: %v\n", path, err)
				} else {
					fmt.Printf("Deleted: %s\n", path)
					deletedCount++
				}
			}
			totalSize += info.Size()
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("error walking sessions: %w", err)
	}

	if dryRun {
		fmt.Printf("\nWould delete %d files (%.2f KB)\n", deletedCount, float64(totalSize)/1024)
	} else {
		fmt.Printf("\nDeleted %d files (%.2f KB freed)\n", deletedCount, float64(totalSize)/1024)
	}

	return nil
}
