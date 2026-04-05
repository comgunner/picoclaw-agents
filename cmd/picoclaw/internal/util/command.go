package util

import "github.com/spf13/cobra"

func NewUtilCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "util",
		Short: "Utility commands",
	}

	cmd.AddCommand(
		newBinanceMCPServerCommand(),
		newSocialMediaMCPServerCommand(),
		newNotionMCPServerCommand(),
		newCodegenCommand(),
	)
	return cmd
}
