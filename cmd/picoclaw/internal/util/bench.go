package util

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/comgunner/picoclaw/pkg/utils"
)

func newBenchCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "bench [command]",
		Short:   "Benchmark command startup time and memory usage",
		Example: "  picoclaw-agents util bench ./build/picoclaw-agents -- --help\n  picoclaw-agents util bench --self",
		RunE: func(cmd *cobra.Command, args []string) error {
			self, _ := cmd.Flags().GetBool("self")
			if self {
				exe, _ := os.Executable()
				elapsed, snap, err := utils.BenchmarkStartup(exe, []string{"--version"})
				if err != nil {
					return fmt.Errorf("benchmark failed: %w", err)
				}
				fmt.Printf("\nStartup time: %v\n", elapsed)
				fmt.Printf("Peak RSS: %d MB\n", snap.SysMB)
				fmt.Printf("Alloc: %d MB\n", snap.AllocMB)
				return nil
			}
			if len(args) == 0 {
				return fmt.Errorf("provide a command to benchmark, or use --self")
			}
			elapsed, snap, err := utils.BenchmarkStartup(args[0], args[1:])
			if err != nil {
				return err
			}
			fmt.Printf("\nStartup time: %v\n", elapsed)
			fmt.Printf("Peak RSS: %d MB\n", snap.SysMB)
			return nil
		},
	}
	cmd.Flags().Bool("self", false, "Benchmark the picoclaw-agents binary itself")
	return cmd
}
