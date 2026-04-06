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
	var fbAppID string
	var fbAppSecret string
	var fbUserToken string
	var xAPIKey string
	var xAPISecret string
	var xAccessToken string
	var xAccessTokenSecret string

	cmd := &cobra.Command{
		Use:   "social-media-mcp-server",
		Short: "Start Social Media MCP server over stdio (Facebook + X)",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			cfg := &utils.SocialMediaMCPConfig{
				FacebookPageID:    orEnv("PICOCLAW_TOOLS_SOCIAL_FACEBOOK_DEFAULT_PAGE_ID", fbPageID),
				FacebookPageToken: orEnv("PICOCLAW_TOOLS_SOCIAL_FACEBOOK_DEFAULT_PAGE_TOKEN", fbPageToken),
				FacebookAppID:     orEnv("PICOCLAW_TOOLS_SOCIAL_FACEBOOK_APP_ID", fbAppID),
				FacebookAppSecret: orEnv("PICOCLAW_TOOLS_SOCIAL_FACEBOOK_APP_SECRET", fbAppSecret),
				FacebookUserToken: orEnv("PICOCLAW_TOOLS_SOCIAL_FACEBOOK_USER_TOKEN", fbUserToken),
				XAPIKey:           orEnv("PICOCLAW_TOOLS_SOCIAL_X_API_KEY", xAPIKey),
				XAPISecret:        orEnv("PICOCLAW_TOOLS_SOCIAL_X_API_SECRET", xAPISecret),
				XAccessToken:      orEnv("PICOCLAW_TOOLS_SOCIAL_X_ACCESS_TOKEN", xAccessToken),
				XAccessSecret:     orEnv("PICOCLAW_TOOLS_SOCIAL_X_ACCESS_TOKEN_SECRET", xAccessTokenSecret),
			}
			return utils.ServeSocialMediaMCPStdio(cfg)
		},
	}

	cmd.Flags().StringVar(&fbPageID, "fb-page-id", "", "Facebook Page ID")
	cmd.Flags().StringVar(&fbPageToken, "fb-page-token", "", "Facebook Page Token")
	cmd.Flags().StringVar(&fbAppID, "fb-app-id", "", "Facebook App ID")
	cmd.Flags().StringVar(&fbAppSecret, "fb-app-secret", "", "Facebook App Secret")
	cmd.Flags().StringVar(&fbUserToken, "fb-user-token", "", "Facebook User Token")
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
			cfg := &utils.NotionMCPConfig{
				APIKey: orEnv("NOTION_API_KEY", apiKey),
			}
			return utils.ServeNotionMCPStdio(cfg)
		},
	}

	cmd.Flags().StringVar(&apiKey, "api-key", "", "Notion API key")

	return cmd
}

// orEnv returns the env var value if non-empty, otherwise the fallback.
func orEnv(envKey, fallback string) string {
	if v := os.Getenv(envKey); v != "" {
		return strings.TrimSpace(v)
	}
	return strings.TrimSpace(fallback)
}

func newGoogleWorkspaceMCPServerCommand() *cobra.Command {
	var credentialsPath string

	cmd := &cobra.Command{
		Use:   "google-workspace-mcp-server",
		Short: "Start Google Workspace MCP server over stdio (Gmail + Calendar)",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			cfg := &utils.GoogleWorkspaceMCPConfig{
				CredentialsJSON: orEnv("GOOGLE_WORKSPACE_CREDENTIALS", credentialsPath),
			}
			return utils.ServeGoogleWorkspaceMCPStdio(cfg)
		},
	}

	cmd.Flags().StringVar(&credentialsPath, "credentials", "", "Path to Google credentials JSON")
	return cmd
}
