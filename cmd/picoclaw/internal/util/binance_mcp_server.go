package util

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/comgunner/picoclaw/pkg/utils"
)

func newBinanceMCPServerCommand() *cobra.Command {
	var apiKey string
	var secretKey string

	cmd := &cobra.Command{
		Use:   "binance-mcp-server",
		Short: "Start Binance MCP server over stdio",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			resolvedAPI := strings.TrimSpace(apiKey)
			resolvedSecret := strings.TrimSpace(secretKey)
			if resolvedAPI == "" {
				resolvedAPI = strings.TrimSpace(os.Getenv(utils.EnvBinanceAPIKey))
			}
			if resolvedSecret == "" {
				resolvedSecret = strings.TrimSpace(os.Getenv(utils.EnvBinanceSecretKey))
			}

			if (resolvedAPI == "") != (resolvedSecret == "") {
				fmt.Fprintln(
					os.Stderr,
					"warning: set both BINANCE_API_KEY and BINANCE_SECRET_KEY to enable private balance tool",
				)
				resolvedAPI = ""
				resolvedSecret = ""
			}

			return utils.ServeBinanceMCPStdio(resolvedAPI, resolvedSecret)
		},
	}

	cmd.Flags().StringVar(&apiKey, "api-key", "", "Binance API key (overrides BINANCE_API_KEY)")
	cmd.Flags().StringVar(&secretKey, "secret-key", "", "Binance secret key (overrides BINANCE_SECRET_KEY)")

	return cmd
}

func newSocialMediaMCPServerCommand() *cobra.Command {
	var fbPageID string
	var fbPageToken string
	var xAPIKey string
	var xAPISecret string
	var xAccessToken string
	var xAccessTokenSecret string

	cmd := &cobra.Command{
		Use:   "social-media-mcp-server",
		Short: "Start Social Media MCP server over stdio (Facebook + X)",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			// Resolve Facebook credentials
			resolvedFBPageID := strings.TrimSpace(fbPageID)
			resolvedFBPageToken := strings.TrimSpace(fbPageToken)
			if resolvedFBPageID == "" {
				resolvedFBPageID = strings.TrimSpace(os.Getenv(utils.EnvSocialFacebookPageID))
			}
			if resolvedFBPageToken == "" {
				resolvedFBPageToken = strings.TrimSpace(os.Getenv(utils.EnvSocialFacebookPageToken))
			}

			// Resolve X credentials
			resolvedXAPIKey := strings.TrimSpace(xAPIKey)
			resolvedXAPISecret := strings.TrimSpace(xAPISecret)
			resolvedXAccessToken := strings.TrimSpace(xAccessToken)
			resolvedXAccessTokenSecret := strings.TrimSpace(xAccessTokenSecret)

			if resolvedXAPIKey == "" {
				resolvedXAPIKey = strings.TrimSpace(os.Getenv(utils.EnvSocialXAPIKey))
			}
			if resolvedXAPISecret == "" {
				resolvedXAPISecret = strings.TrimSpace(os.Getenv(utils.EnvSocialXAPISecret))
			}
			if resolvedXAccessToken == "" {
				resolvedXAccessToken = strings.TrimSpace(os.Getenv(utils.EnvSocialXAccessToken))
			}
			if resolvedXAccessTokenSecret == "" {
				resolvedXAccessTokenSecret = strings.TrimSpace(os.Getenv(utils.EnvSocialXAccessTokenSecret))
			}

			// TODO: Implement ServeSocialMediaMCPStdio in utils
			fmt.Printf("Social Media MCP Server (placeholder)\n")
			fmt.Printf("Facebook Page ID: %s\n", resolvedFBPageID)
			fmt.Printf("X API Key configured: %v\n", resolvedXAPIKey != "")
			return nil
		},
	}

	cmd.Flags().StringVar(&fbPageID, "fb-page-id", "", "Facebook Page ID")
	cmd.Flags().StringVar(&fbPageToken, "fb-page-token", "", "Facebook Page Token")
	cmd.Flags().StringVar(&xAPIKey, "x-api-key", "", "X API Key")
	cmd.Flags().StringVar(&xAPISecret, "x-api-secret", "", "X API Secret")
	cmd.Flags().StringVar(&xAccessToken, "x-access-token", "", "X Access Token")
	cmd.Flags().StringVar(&xAccessTokenSecret, "x-access-token-secret", "", "X Access Token Secret")

	return cmd
}

func newNotionMCPServerCommand() *cobra.Command {
	var apiKey string

	cmd := &cobra.Command{
		Use:   "notion-mcp-server",
		Short: "Start Notion MCP server over stdio",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			resolvedAPIKey := strings.TrimSpace(apiKey)
			if resolvedAPIKey == "" {
				resolvedAPIKey = strings.TrimSpace(os.Getenv(utils.EnvNotionAPIKey))
			}

			// TODO: Implement ServeNotionMCPStdio in utils
			fmt.Printf("Notion MCP Server (placeholder)\n")
			fmt.Printf("API Key configured: %v\n", resolvedAPIKey != "")
			return nil
		},
	}

	cmd.Flags().StringVar(&apiKey, "api-key", "", "Notion API key (overrides PICOCLAW_TOOLS_NOTION_API_KEY)")

	return cmd
}
