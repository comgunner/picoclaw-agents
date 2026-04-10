package util

import (
	"github.com/spf13/cobra"

	"github.com/comgunner/picoclaw/cmd/picoclaw/internal/mcp"
)

func NewUtilCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "util",
		Short: "Utility commands",
	}

	cmd.AddCommand(
		newBinanceMCPServerCommand(),
		newSocialMediaMCPServerCommand(),
		newNotionMCPServerCommand(),
		newGoogleWorkspaceMCPServerCommand(),
		newCodegenCommand(),
		mcp.NewMCPCommand(),
		// Productivity tools
		newBenchCommand(),
		newReaperCommand(),
		newArchLintCommand(),
		newMdAuditCommand(),
	)
	return cmd
}
