// PicoClaw - Social Media MCP Server
// Exposes Facebook and X/Twitter tools via Model Context Protocol (stdio).
//
// Usage:
//
//	picoclaw-agents util social-media-mcp-server \
//	  --fb-page-id YOUR_PAGE_ID \
//	  --fb-page-token YOUR_PAGE_TOKEN \
//	  --x-api-key YOUR_KEY \
//	  --x-api-secret YOUR_SECRET \
//	  --x-access-token YOUR_TOKEN \
//	  --x-access-token-secret YOUR_TOKEN_SECRET
package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// SocialMediaMCPConfig holds credentials for Facebook and X.
type SocialMediaMCPConfig struct {
	FacebookPageID    string
	FacebookPageToken string
	FacebookAppID     string
	FacebookAppSecret string
	FacebookUserToken string
	XAPIKey           string
	XAPISecret        string
	XAccessToken      string
	XAccessSecret     string
}

// SocialMediaConfigFromEnv resolves credentials from environment variables.
func SocialMediaConfigFromEnv() *SocialMediaMCPConfig {
	return &SocialMediaMCPConfig{
		FacebookPageID:    strings.TrimSpace(os.Getenv("PICOCLAW_TOOLS_SOCIAL_FACEBOOK_DEFAULT_PAGE_ID")),
		FacebookPageToken: strings.TrimSpace(os.Getenv("PICOCLAW_TOOLS_SOCIAL_FACEBOOK_DEFAULT_PAGE_TOKEN")),
		FacebookAppID:     strings.TrimSpace(os.Getenv("PICOCLAW_TOOLS_SOCIAL_FACEBOOK_APP_ID")),
		FacebookAppSecret: strings.TrimSpace(os.Getenv("PICOCLAW_TOOLS_SOCIAL_FACEBOOK_APP_SECRET")),
		FacebookUserToken: strings.TrimSpace(os.Getenv("PICOCLAW_TOOLS_SOCIAL_FACEBOOK_USER_TOKEN")),
		XAPIKey:           strings.TrimSpace(os.Getenv("PICOCLAW_TOOLS_SOCIAL_X_API_KEY")),
		XAPISecret:        strings.TrimSpace(os.Getenv("PICOCLAW_TOOLS_SOCIAL_X_API_SECRET")),
		XAccessToken:      strings.TrimSpace(os.Getenv("PICOCLAW_TOOLS_SOCIAL_X_ACCESS_TOKEN")),
		XAccessSecret:     strings.TrimSpace(os.Getenv("PICOCLAW_TOOLS_SOCIAL_X_ACCESS_TOKEN_SECRET")),
	}
}

// NewSocialMediaMCPServer builds an MCP server with Facebook and X tools.
func NewSocialMediaMCPServer(cfg *SocialMediaMCPConfig) *server.MCPServer {
	s := server.NewMCPServer(
		"social-media-mcp",
		"1.0.0",
	)

	// ─── Facebook Tools ──────────────────────────────────────────

	s.AddTool(
		mcp.NewTool(
			"facebook_post_text",
			mcp.WithDescription("Post text-only message to a Facebook Page"),
			mcp.WithString("page_id", mcp.Description("Facebook Page ID"), mcp.Required()),
			mcp.WithString("page_token", mcp.Description("Page Access Token")),
			mcp.WithString("message", mcp.Description("Post message"), mcp.Required()),
			mcp.WithString("comment", mcp.Description("Optional comment to add after posting")),
		),
		fbPostTextHandler(cfg),
	)

	s.AddTool(
		mcp.NewTool(
			"facebook_post_image",
			mcp.WithDescription("Post image with caption to a Facebook Page"),
			mcp.WithString("page_id", mcp.Description("Facebook Page ID"), mcp.Required()),
			mcp.WithString("page_token", mcp.Description("Page Access Token")),
			mcp.WithString("message", mcp.Description("Image caption"), mcp.Required()),
			mcp.WithString("image_path", mcp.Description("Absolute path to image file"), mcp.Required()),
			mcp.WithString("comment", mcp.Description("Optional comment to add after posting")),
		),
		fbPostImageHandler(cfg),
	)

	s.AddTool(
		mcp.NewTool(
			"facebook_post",
			mcp.WithDescription("Post to Facebook (text or image). Full control."),
			mcp.WithString("page_id", mcp.Description("Facebook Page ID"), mcp.Required()),
			mcp.WithString("page_token", mcp.Description("Page Access Token")),
			mcp.WithString("message", mcp.Description("Post message"), mcp.Required()),
			mcp.WithString("image_path", mcp.Description("Optional absolute image path")),
			mcp.WithString("comment", mcp.Description("Optional comment")),
		),
		fbPostHandler(cfg),
	)

	// ─── X/Twitter Tools ─────────────────────────────────────────

	s.AddTool(
		mcp.NewTool(
			"x_post_tweet",
			mcp.WithDescription("Post a single tweet on X/Twitter"),
			mcp.WithString("message", mcp.Description("Tweet text (max 280 chars)"), mcp.Required()),
			mcp.WithString("image_path", mcp.Description("Optional absolute image path")),
			mcp.WithString("reply_to_tweet_id", mcp.Description("Optional tweet ID to reply to")),
		),
		xPostTweetHandler(cfg),
	)

	s.AddTool(
		mcp.NewTool(
			"x_post_thread",
			mcp.WithDescription("Post a thread of tweets on X/Twitter"),
			mcp.WithString(
				"tweets",
				mcp.Description("JSON array of tweet texts, e.g. [\"tweet1\",\"tweet2\"]"),
				mcp.Required(),
			),
		),
		xPostThreadHandler(cfg),
	)

	return s
}

// ServeSocialMediaMCPStdio starts the Social Media MCP server over stdio.
func ServeSocialMediaMCPStdio(cfg *SocialMediaMCPConfig) error {
	return server.ServeStdio(NewSocialMediaMCPServer(cfg))
}

// ─── Handlers: Facebook ──────────────────────────────────────────

func fbPostTextHandler(cfg *SocialMediaMCPConfig) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		pageID := getStrArg(req, "page_id", cfg.FacebookPageID)
		pageToken := getStrArg(req, "page_token", cfg.FacebookPageToken)
		message := getStrArg(req, "message", "")
		comment := getStrArg(req, "comment", "")

		if pageID == "" {
			return mcp.NewToolResultError("page_id is required"), nil
		}
		if message == "" {
			return mcp.NewToolResultError("message is required"), nil
		}

		// appID/appSecret/userToken come from config env, not MCP args (for token refresh flow)
		postID, err := FacebookPostTextOnly(
			ctx,
			pageID,
			pageToken,
			cfg.FacebookAppID,
			cfg.FacebookAppSecret,
			cfg.FacebookUserToken,
			message,
			comment,
		)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Facebook post failed: %v", err)), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("✅ Facebook post published. Post ID: %s", postID)), nil
	}
}

func fbPostImageHandler(cfg *SocialMediaMCPConfig) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		pageID := getStrArg(req, "page_id", cfg.FacebookPageID)
		pageToken := getStrArg(req, "page_token", cfg.FacebookPageToken)
		message := getStrArg(req, "message", "")
		imagePath := getStrArg(req, "image_path", "")
		comment := getStrArg(req, "comment", "")

		if pageID == "" {
			return mcp.NewToolResultError("page_id is required"), nil
		}
		if message == "" {
			return mcp.NewToolResultError("message is required"), nil
		}
		if imagePath == "" {
			return mcp.NewToolResultError("image_path is required"), nil
		}

		req2 := FBPostRequest{
			PageID:    pageID,
			PageToken: pageToken,
			Message:   message,
			ImagePath: imagePath,
			Comment:   comment,
		}
		postID, err := FacebookPost(ctx, req2)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Facebook image post failed: %v", err)), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("✅ Facebook image post published. Post ID: %s", postID)), nil
	}
}

func fbPostHandler(cfg *SocialMediaMCPConfig) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		pageID := getStrArg(req, "page_id", cfg.FacebookPageID)
		pageToken := getStrArg(req, "page_token", cfg.FacebookPageToken)
		message := getStrArg(req, "message", "")
		imagePath := getStrArg(req, "image_path", "")
		comment := getStrArg(req, "comment", "")

		if pageID == "" {
			return mcp.NewToolResultError("page_id is required"), nil
		}
		if message == "" {
			return mcp.NewToolResultError("message is required"), nil
		}

		var postID string
		var err error
		if imagePath == "" {
			postID, err = FacebookPostTextOnly(
				ctx,
				pageID,
				pageToken,
				cfg.FacebookAppID,
				cfg.FacebookAppSecret,
				cfg.FacebookUserToken,
				message,
				comment,
			)
		} else {
			req2 := FBPostRequest{
				PageID:    pageID,
				PageToken: pageToken,
				Message:   message,
				ImagePath: imagePath,
				Comment:   comment,
			}
			postID, err = FacebookPost(ctx, req2)
		}
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Facebook post failed: %v", err)), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("✅ Facebook post published. Post ID: %s", postID)), nil
	}
}

// ─── Handlers: X/Twitter ─────────────────────────────────────────

func xPostTweetHandler(cfg *SocialMediaMCPConfig) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		apiKey := cfg.XAPIKey
		apiSecret := cfg.XAPISecret
		accessToken := cfg.XAccessToken
		accessSecret := cfg.XAccessSecret

		if apiKey == "" || apiSecret == "" || accessToken == "" || accessSecret == "" {
			return mcp.NewToolResultError("X/Twitter credentials not configured"), nil
		}

		message := getStrArg(req, "message", "")
		imagePath := getStrArg(req, "image_path", "")
		replyToTweetID := getStrArg(req, "reply_to_tweet_id", "")

		if message == "" {
			return mcp.NewToolResultError("message is required"), nil
		}

		tReq := XPostRequest{
			APIKey:            apiKey,
			APISecret:         apiSecret,
			AccessToken:       accessToken,
			AccessTokenSecret: accessSecret,
			Message:           message,
			ImagePath:         imagePath,
			ReplyToTweetID:    replyToTweetID,
		}

		tweetID, _, err := XPostTweet(ctx, tReq)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("X tweet failed: %v", err)), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("✅ Tweet published. Tweet ID: %s", tweetID)), nil
	}
}

func xPostThreadHandler(cfg *SocialMediaMCPConfig) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		apiKey := cfg.XAPIKey
		apiSecret := cfg.XAPISecret
		accessToken := cfg.XAccessToken
		accessSecret := cfg.XAccessSecret

		if apiKey == "" || apiSecret == "" || accessToken == "" || accessSecret == "" {
			return mcp.NewToolResultError("X/Twitter credentials not configured"), nil
		}

		tweetsStr := getStrArg(req, "tweets", "")
		if tweetsStr == "" {
			return mcp.NewToolResultError("tweets JSON array is required"), nil
		}

		tweets, err := parseTweetArray(tweetsStr)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid tweets JSON: %v", err)), nil
		}
		if len(tweets) == 0 {
			return mcp.NewToolResultError("tweets array must not be empty"), nil
		}
		if len(tweets) > 25 {
			return mcp.NewToolResultError("Thread limited to 25 tweets"), nil
		}

		tReq := XPostRequest{
			APIKey:            apiKey,
			APISecret:         apiSecret,
			AccessToken:       accessToken,
			AccessTokenSecret: accessSecret,
		}

		ids := make([]string, 0, len(tweets))
		for i, msg := range tweets {
			tReq.Message = msg
			if i > 0 && len(ids) > 0 {
				tReq.ReplyToTweetID = ids[len(ids)-1]
			}
			tweetID, _, err := XPostTweet(ctx, tReq)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Thread failed at tweet %d: %v", i+1, err)), nil
			}
			ids = append(ids, tweetID)
		}

		return mcp.NewToolResultText(
			fmt.Sprintf("✅ Thread published: %d tweets. IDs: %s", len(ids), strings.Join(ids, ", ")),
		), nil
	}
}

// ─── Helpers ─────────────────────────────────────────────────────

func getStrArg(req mcp.CallToolRequest, key string, fallback string) string {
	v, ok := req.GetArguments()[key]
	if !ok {
		return fallback
	}
	s, ok := v.(string)
	if !ok {
		return fallback
	}
	return strings.TrimSpace(s)
}

// parseTweetArray parses a JSON array of tweet texts.
func parseTweetArray(s string) ([]string, error) {
	var tweets []string
	if err := json.Unmarshal([]byte(s), &tweets); err != nil {
		return nil, err
	}
	// Trim each
	for i, t := range tweets {
		tweets[i] = strings.TrimSpace(t)
	}
	return tweets, nil
}
