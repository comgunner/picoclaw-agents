// PicoClaw - Notion MCP Server
// Exposes Notion tools via Model Context Protocol (stdio).
//
// Usage:
//
//	picoclaw-agents util notion-mcp-server --api-key YOUR_NOTION_API_KEY
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

// NotionMCPConfig holds Notion API credentials.
type NotionMCPConfig struct {
	APIKey string
}

// NotionConfigFromEnv resolves Notion API key from environment.
func NotionConfigFromEnv() *NotionMCPConfig {
	return &NotionMCPConfig{
		APIKey: strings.TrimSpace(os.Getenv("NOTION_API_KEY")),
	}
}

// NewNotionMCPServer builds an MCP server with Notion tools.
func NewNotionMCPServer(cfg *NotionMCPConfig) *server.MCPServer {
	s := server.NewMCPServer("notion-mcp", "1.0.0")

	s.AddTool(
		mcp.NewTool(
			"notion_query_database",
			mcp.WithDescription("Query a Notion database and return results"),
			mcp.WithString("database_id", mcp.Description("Notion Database ID"), mcp.Required()),
			mcp.WithString("filter", mcp.Description("Optional JSON filter string")),
		),
		notionQueryDBHandler(cfg),
	)

	s.AddTool(
		mcp.NewTool(
			"notion_get_page",
			mcp.WithDescription("Get a Notion page by ID"),
			mcp.WithString("page_id", mcp.Description("Notion Page ID"), mcp.Required()),
		),
		notionGetPageHandler(cfg),
	)

	s.AddTool(
		mcp.NewTool(
			"notion_create_page",
			mcp.WithDescription("Create a new Notion page"),
			mcp.WithString("parent_id", mcp.Description("Parent page or database ID"), mcp.Required()),
			mcp.WithString("properties", mcp.Description("Page properties as JSON"), mcp.Required()),
			mcp.WithString("content", mcp.Description("Optional page content (JSON blocks array)")),
		),
		notionCreatePageHandler(cfg),
	)

	s.AddTool(
		mcp.NewTool(
			"notion_update_page",
			mcp.WithDescription("Update an existing Notion page"),
			mcp.WithString("page_id", mcp.Description("Notion Page ID"), mcp.Required()),
			mcp.WithString("properties", mcp.Description("Updated properties as JSON"), mcp.Required()),
		),
		notionUpdatePageHandler(cfg),
	)

	s.AddTool(
		mcp.NewTool(
			"notion_search",
			mcp.WithDescription("Search Notion pages and databases"),
			mcp.WithString("query", mcp.Description("Search query")),
		),
		notionSearchHandler(cfg),
	)

	s.AddTool(
		mcp.NewTool(
			"notion_append_block",
			mcp.WithDescription("Append content blocks to a Notion page"),
			mcp.WithString("page_id", mcp.Description("Page ID to append to"), mcp.Required()),
			mcp.WithString("blocks", mcp.Description("Blocks as JSON array"), mcp.Required()),
		),
		notionAppendBlockHandler(cfg),
	)

	return s
}

// ServeNotionMCPStdio starts the Notion MCP server over stdio.
func ServeNotionMCPStdio(cfg *NotionMCPConfig) error {
	if cfg.APIKey == "" {
		return fmt.Errorf("NOTION_API_KEY is required")
	}
	return server.ServeStdio(NewNotionMCPServer(cfg))
}

// ─── Helpers ─────────────────────────────────────────────────────

func notionClient(cfg *NotionMCPConfig) *NotionClient {
	return NewNotionClient(cfg.APIKey)
}

func notionQueryDBHandler(cfg *NotionMCPConfig) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if cfg.APIKey == "" {
			return mcp.NewToolResultError("NOTION_API_KEY not configured"), nil
		}
		dbID := getMCPStr(req, "database_id", "")
		if dbID == "" {
			return mcp.NewToolResultError("database_id is required"), nil
		}
		client := notionClient(cfg)
		resp, err := client.NotionQueryDatabase(ctx, dbID, nil)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Notion query failed: %v", err)), nil
		}
		data, _ := json.Marshal(resp)
		return mcp.NewToolResultText(string(data)), nil
	}
}

func notionGetPageHandler(cfg *NotionMCPConfig) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if cfg.APIKey == "" {
			return mcp.NewToolResultError("NOTION_API_KEY not configured"), nil
		}
		pageID := getMCPStr(req, "page_id", "")
		if pageID == "" {
			return mcp.NewToolResultError("page_id is required"), nil
		}
		client := notionClient(cfg)
		page, err := client.NotionGetPage(ctx, pageID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Notion get page failed: %v", err)), nil
		}
		data, _ := json.Marshal(page)
		return mcp.NewToolResultText(string(data)), nil
	}
}

func notionCreatePageHandler(cfg *NotionMCPConfig) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if cfg.APIKey == "" {
			return mcp.NewToolResultError("NOTION_API_KEY not configured"), nil
		}
		parentID := getMCPStr(req, "parent_id", "")
		propsStr := getMCPStr(req, "properties", "")
		contentStr := getMCPStr(req, "content", "")
		if parentID == "" {
			return mcp.NewToolResultError("parent_id is required"), nil
		}
		if propsStr == "" {
			return mcp.NewToolResultError("properties is required"), nil
		}
		var props map[string]any
		if err := json.Unmarshal([]byte(propsStr), &props); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid properties JSON: %v", err)), nil
		}
		client := notionClient(cfg)
		page, err := client.NotionCreatePage(ctx, parentID, props)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Notion create failed: %v", err)), nil
		}
		// Append blocks if content provided
		if contentStr != "" {
			var blocks []map[string]any
			if err := json.Unmarshal([]byte(contentStr), &blocks); err == nil {
				_ = client.NotionCreateBlock(ctx, page.ID, blocks)
			}
		}
		data, _ := json.Marshal(page)
		return mcp.NewToolResultText(string(data)), nil
	}
}

func notionUpdatePageHandler(cfg *NotionMCPConfig) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if cfg.APIKey == "" {
			return mcp.NewToolResultError("NOTION_API_KEY not configured"), nil
		}
		pageID := getMCPStr(req, "page_id", "")
		propsStr := getMCPStr(req, "properties", "")
		if pageID == "" {
			return mcp.NewToolResultError("page_id is required"), nil
		}
		if propsStr == "" {
			return mcp.NewToolResultError("properties is required"), nil
		}
		var props map[string]any
		if err := json.Unmarshal([]byte(propsStr), &props); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid properties JSON: %v", err)), nil
		}
		client := notionClient(cfg)
		page, err := client.NotionUpdatePage(ctx, pageID, props)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Notion update failed: %v", err)), nil
		}
		data, _ := json.Marshal(page)
		return mcp.NewToolResultText(string(data)), nil
	}
}

func notionSearchHandler(cfg *NotionMCPConfig) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if cfg.APIKey == "" {
			return mcp.NewToolResultError("NOTION_API_KEY not configured"), nil
		}
		query := getMCPStr(req, "query", "")
		client := notionClient(cfg)
		resp, err := client.NotionSearch(ctx, query)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Notion search failed: %v", err)), nil
		}
		data, _ := json.Marshal(resp)
		return mcp.NewToolResultText(string(data)), nil
	}
}

func notionAppendBlockHandler(cfg *NotionMCPConfig) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if cfg.APIKey == "" {
			return mcp.NewToolResultError("NOTION_API_KEY not configured"), nil
		}
		pageID := getMCPStr(req, "page_id", "")
		blocksStr := getMCPStr(req, "blocks", "")
		if pageID == "" {
			return mcp.NewToolResultError("page_id is required"), nil
		}
		if blocksStr == "" {
			return mcp.NewToolResultError("blocks is required"), nil
		}
		var blocks []map[string]any
		if err := json.Unmarshal([]byte(blocksStr), &blocks); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid blocks JSON: %v", err)), nil
		}
		client := notionClient(cfg)
		if err := client.NotionCreateBlock(ctx, pageID, blocks); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Notion append block failed: %v", err)), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("✅ Appended %d blocks to page %s", len(blocks), pageID)), nil
	}
}
