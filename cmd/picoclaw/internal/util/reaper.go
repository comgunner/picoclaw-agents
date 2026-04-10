package util

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/comgunner/picoclaw/pkg/utils"
)

func newReaperCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reaper",
		Short: "Find and kill orphaned picoclaw-agents processes",
		RunE: func(cmd *cobra.Command, args []string) error {
			dryRun, _ := cmd.Flags().GetBool("dry-run")
			orphans, err := utils.FindOrphans()
			if err != nil {
				return fmt.Errorf("find orphans: %w", err)
			}
			if len(orphans) == 0 {
				fmt.Println("✅ No orphan processes found.")
				return nil
			}
			fmt.Printf("Found %d orphan(s):\n", len(orphans))
			for _, o := range orphans {
				fmt.Printf("  PID %d: %s\n", o.PID, o.Cmd)
			}
			if dryRun {
				fmt.Println("\nDry run — no processes killed.")
				return nil
			}
			killed, err := utils.KillOrphans()
			if err != nil {
				return fmt.Errorf("kill orphans: %w", err)
			}
			fmt.Printf("\n✅ Killed %d orphan(s).\n", len(killed))
			return nil
		},
	}
	cmd.Flags().Bool("dry-run", false, "Show orphans without killing them")
	return cmd
}
