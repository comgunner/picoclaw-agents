package util

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/comgunner/picoclaw/pkg/utils"
)

func newArchLintCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "arch-lint [root]",
		Short: "Check for forbidden import patterns between packages",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			root := "."
			if len(args) > 0 {
				root = args[0]
			}
			root, err := filepath.Abs(root)
			if err != nil {
				return fmt.Errorf("resolve root: %w", err)
			}
			violations, err := utils.CheckImports(root, nil)
			if err != nil {
				return fmt.Errorf("check imports: %w", err)
			}
			utils.PrintViolations(violations)
			if len(violations) > 0 {
				os.Exit(1)
			}
			return nil
		},
	}
	return cmd
}
