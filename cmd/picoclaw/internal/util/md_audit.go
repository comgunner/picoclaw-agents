package util

import (
	"fmt"

	"github.com/comgunner/picoclaw/pkg/utils"
	"github.com/spf13/cobra"
)

func newMdAuditCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "md-audit [docs-dir]",
		Short: "Audit Markdown files for broken internal links",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			root := "docs"
			if len(args) > 0 {
				root = args[0]
			}
			issues, err := utils.AuditMarkdown(root)
			if err != nil {
				return fmt.Errorf("audit: %w", err)
			}
			utils.PrintLinkIssues(issues)
			return nil
		},
	}
	return cmd
}
