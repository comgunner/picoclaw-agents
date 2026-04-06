package util

import (
	"github.com/comgunner/picoclaw/cmd/picoclaw/internal/mcp"
	"github.com/spf13/cobra"
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
