// PicoClaw - Google Workspace MCP Server (Gmail + Calendar)
// Exposes Gmail and Google Calendar tools via Model Context Protocol (stdio).
//
// Usage:
//
//	picoclaw-agents util google-workspace-mcp-server --credentials /path/to/credentials.json
package utils

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

// GoogleWorkspaceMCPConfig holds Google credentials.
type GoogleWorkspaceMCPConfig struct {
	CredentialsJSON string // Path to service account or OAuth credentials JSON
}

// GoogleWorkspaceConfigFromEnv resolves config from environment.
func GoogleWorkspaceConfigFromEnv() *GoogleWorkspaceMCPConfig {
	return &GoogleWorkspaceMCPConfig{
		CredentialsJSON: strings.TrimSpace(os.Getenv("GOOGLE_WORKSPACE_CREDENTIALS")),
	}
}

// NewGoogleWorkspaceMCPServer builds an MCP server with Gmail and Calendar tools.
func NewGoogleWorkspaceMCPServer(cfg *GoogleWorkspaceMCPConfig) (*server.MCPServer, error) {
	s := server.NewMCPServer("google-workspace-mcp", "1.0.0")

	if cfg.CredentialsJSON == "" {
		return s, nil // Return server without tools if not configured
	}

	ctx := context.Background()
	gmailSvc, err := gmail.NewService(ctx, option.WithCredentialsFile(cfg.CredentialsJSON))
	if err != nil {
		return nil, fmt.Errorf("gmail auth: %w", err)
	}
	calSvc, err := calendar.NewService(ctx, option.WithCredentialsFile(cfg.CredentialsJSON))
	if err != nil {
		return nil, fmt.Errorf("calendar auth: %w", err)
	}

	// ─── Gmail Tools ──────────────────────────────────────────────

	s.AddTool(
		mcp.NewTool("gmail_list_emails",
			mcp.WithDescription("List recent emails in Gmail inbox"),
			mcp.WithNumber("max_results", mcp.Description("Max results (default 10, max 50)")),
			mcp.WithString("query", mcp.Description("Optional Gmail search query")),
		),
		gmailListEmailsHandler(gmailSvc),
	)

	s.AddTool(
		mcp.NewTool("gmail_get_email",
			mcp.WithDescription("Get a specific email by message ID"),
			mcp.WithString("message_id", mcp.Description("Gmail message ID"), mcp.Required()),
		),
		gmailGetEmailHandler(gmailSvc),
	)

	s.AddTool(
		mcp.NewTool("gmail_send_email",
			mcp.WithDescription("Send an email via Gmail"),
			mcp.WithString("to", mcp.Description("Recipient email"), mcp.Required()),
			mcp.WithString("subject", mcp.Description("Email subject"), mcp.Required()),
			mcp.WithString("body", mcp.Description("Email body"), mcp.Required()),
		),
		gmailSendEmailHandler(gmailSvc),
	)

	s.AddTool(
		mcp.NewTool("gmail_search",
			mcp.WithDescription("Search emails using Gmail search syntax"),
			mcp.WithString("query", mcp.Description("Gmail search query"), mcp.Required()),
			mcp.WithNumber("max_results", mcp.Description("Max results (default 10, max 50)")),
		),
		gmailSearchHandler(gmailSvc),
	)

	// ─── Calendar Tools ───────────────────────────────────────────

	s.AddTool(
		mcp.NewTool("calendar_list_events",
			mcp.WithDescription("List upcoming calendar events"),
			mcp.WithNumber("max_results", mcp.Description("Max results (default 10)")),
			mcp.WithString("calendar_id", mcp.Description("Calendar ID (default: primary)")),
		),
		calendarListEventsHandler(calSvc),
	)

	s.AddTool(
		mcp.NewTool("calendar_create_event",
			mcp.WithDescription("Create a new calendar event"),
			mcp.WithString("summary", mcp.Description("Event title"), mcp.Required()),
			mcp.WithString("start_time", mcp.Description("Start time RFC3339"), mcp.Required()),
			mcp.WithString("end_time", mcp.Description("End time RFC3339"), mcp.Required()),
			mcp.WithString("description", mcp.Description("Event description")),
			mcp.WithString("calendar_id", mcp.Description("Calendar ID (default: primary)")),
		),
		calendarCreateEventHandler(calSvc),
	)

	s.AddTool(
		mcp.NewTool("calendar_delete_event",
			mcp.WithDescription("Delete a calendar event"),
			mcp.WithString("event_id", mcp.Description("Event ID"), mcp.Required()),
			mcp.WithString("calendar_id", mcp.Description("Calendar ID (default: primary)")),
		),
		calendarDeleteEventHandler(calSvc),
	)

	s.AddTool(
		mcp.NewTool("calendar_free_busy",
			mcp.WithDescription("Check free/busy status for a time range"),
			mcp.WithString("start_time", mcp.Description("Start time RFC3339"), mcp.Required()),
			mcp.WithString("end_time", mcp.Description("End time RFC3339"), mcp.Required()),
			mcp.WithString("calendar_id", mcp.Description("Calendar ID (default: primary)")),
		),
		calendarFreeBusyHandler(calSvc),
	)

	return s, nil
}

// ServeGoogleWorkspaceMCPStdio starts the Google Workspace MCP server over stdio.
func ServeGoogleWorkspaceMCPStdio(cfg *GoogleWorkspaceMCPConfig) error {
	s, err := NewGoogleWorkspaceMCPServer(cfg)
	if err != nil {
		return err
	}
	return server.ServeStdio(s)
}

// ─── Gmail Handlers ──────────────────────────────────────────────

func gmailListEmailsHandler(svc *gmail.Service) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		maxResults := int64(getMCPInt(req, "max_results", 10))
		if maxResults > 50 {
			maxResults = 50
		}
		query := getMCPStr(req, "query", "")

		resp, err := svc.Users.Messages.List("me").Q(query).MaxResults(maxResults).Do()
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Gmail list failed: %v", err)), nil
		}
		if len(resp.Messages) == 0 {
			return mcp.NewToolResultText("📬 No emails found"), nil
		}

		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("📬 %d emails:\n\n", len(resp.Messages)))
		for _, msg := range resp.Messages {
			full, _ := svc.Users.Messages.Get("me", msg.Id).Do()
			if full == nil {
				continue
			}
			subject := getGmailHeader(full.Payload.Headers, "Subject")
			from := getGmailHeader(full.Payload.Headers, "From")
			date := getGmailHeader(full.Payload.Headers, "Date")
			sb.WriteString(fmt.Sprintf("📧 %s\nFrom: %s\nDate: %s\nID: %s\n\n", subject, from, date, msg.Id))
		}
		return mcp.NewToolResultText(sb.String()), nil
	}
}

func gmailGetEmailHandler(svc *gmail.Service) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		msgID := getMCPStr(req, "message_id", "")
		if msgID == "" {
			return mcp.NewToolResultError("message_id is required"), nil
		}
		msg, err := svc.Users.Messages.Get("me", msgID).Format("full").Do()
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Gmail get failed: %v", err)), nil
		}
		headers := msg.Payload.Headers
		subject := getGmailHeader(headers, "Subject")
		from := getGmailHeader(headers, "From")
		to := getGmailHeader(headers, "To")
		date := getGmailHeader(headers, "Date")
		body := extractEmailBody(msg.Payload)
		result := fmt.Sprintf("📧 %s\nFrom: %s\nTo: %s\nDate: %s\n\n%s", subject, from, to, date, body)
		return mcp.NewToolResultText(result), nil
	}
}

func gmailSendEmailHandler(svc *gmail.Service) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		to := getMCPStr(req, "to", "")
		subject := getMCPStr(req, "subject", "")
		body := getMCPStr(req, "body", "")
		if to == "" || subject == "" || body == "" {
			return mcp.NewToolResultError("to, subject, and body are required"), nil
		}
		var msgStr strings.Builder
		msgStr.WriteString(fmt.Sprintf("From: me\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", to, subject, body))
		gmailMsg := &gmail.Message{Raw: base64.URLEncoding.EncodeToString([]byte(msgStr.String()))}
		sent, err := svc.Users.Messages.Send("me", gmailMsg).Do()
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Gmail send failed: %v", err)), nil
		}
		return mcp.NewToolResultText(
			fmt.Sprintf("✅ Email sent! ID: %s\nTo: %s\nSubject: %s", sent.Id, to, subject),
		), nil
	}
}

func gmailSearchHandler(svc *gmail.Service) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		query := getMCPStr(req, "query", "")
		if query == "" {
			return mcp.NewToolResultError("query is required"), nil
		}
		maxResults := int64(getMCPInt(req, "max_results", 10))
		if maxResults > 50 {
			maxResults = 50
		}
		resp, err := svc.Users.Messages.List("me").Q(query).MaxResults(maxResults).Do()
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Gmail search failed: %v", err)), nil
		}
		if len(resp.Messages) == 0 {
			return mcp.NewToolResultText(fmt.Sprintf("🔍 No results for: %s", query)), nil
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("🔍 %d results for \"%s\":\n\n", len(resp.Messages), query))
		for _, msg := range resp.Messages {
			full, _ := svc.Users.Messages.Get("me", msg.Id).Do()
			if full == nil {
				continue
			}
			subject := getGmailHeader(full.Payload.Headers, "Subject")
			from := getGmailHeader(full.Payload.Headers, "From")
			sb.WriteString(fmt.Sprintf("📧 %s\nFrom: %s\nID: %s\n\n", subject, from, msg.Id))
		}
		return mcp.NewToolResultText(sb.String()), nil
	}
}

// ─── Calendar Handlers ───────────────────────────────────────────

func calendarListEventsHandler(svc *calendar.Service) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		calendarID := getMCPStr(req, "calendar_id", "primary")
		maxResults := int64(getMCPInt(req, "max_results", 10))
		if maxResults > 50 {
			maxResults = 50
		}
		now := time.Now().Format(time.RFC3339)
		resp, err := svc.Events.List(calendarID).
			TimeMin(now).
			MaxResults(maxResults).
			SingleEvents(true).
			OrderBy("startTime").
			Do()
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Calendar list failed: %v", err)), nil
		}
		if len(resp.Items) == 0 {
			return mcp.NewToolResultText("📅 No upcoming events"), nil
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("📅 %d upcoming events:\n\n", len(resp.Items)))
		for _, ev := range resp.Items {
			start := "All day"
			if ev.Start.DateTime != "" {
				start = ev.Start.DateTime
			} else if ev.Start.Date != "" {
				start = ev.Start.Date
			}
			sb.WriteString(
				fmt.Sprintf("📌 %s\nStart: %s\nLocation: %s\nID: %s\n\n", ev.Summary, start, ev.Location, ev.Id),
			)
		}
		return mcp.NewToolResultText(sb.String()), nil
	}
}

func calendarCreateEventHandler(svc *calendar.Service) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		calendarID := getMCPStr(req, "calendar_id", "primary")
		summary := getMCPStr(req, "summary", "")
		startTime := getMCPStr(req, "start_time", "")
		endTime := getMCPStr(req, "end_time", "")
		description := getMCPStr(req, "description", "")
		if summary == "" || startTime == "" || endTime == "" {
			return mcp.NewToolResultError("summary, start_time, and end_time are required"), nil
		}
		event := &calendar.Event{
			Summary:     summary,
			Description: description,
			Start:       &calendar.EventDateTime{DateTime: startTime},
			End:         &calendar.EventDateTime{DateTime: endTime},
		}
		created, err := svc.Events.Insert(calendarID, event).Do()
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Calendar create failed: %v", err)), nil
		}
		return mcp.NewToolResultText(
			fmt.Sprintf(
				"✅ Event created!\nTitle: %s\nStart: %s\nLink: %s",
				created.Summary,
				created.Start.DateTime,
				created.HtmlLink,
			),
		), nil
	}
}

func calendarDeleteEventHandler(svc *calendar.Service) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		calendarID := getMCPStr(req, "calendar_id", "primary")
		eventID := getMCPStr(req, "event_id", "")
		if eventID == "" {
			return mcp.NewToolResultError("event_id is required"), nil
		}
		if err := svc.Events.Delete(calendarID, eventID).Do(); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Calendar delete failed: %v", err)), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("✅ Event deleted: %s", eventID)), nil
	}
}

func calendarFreeBusyHandler(svc *calendar.Service) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		calendarID := getMCPStr(req, "calendar_id", "primary")
		startTime := getMCPStr(req, "start_time", "")
		endTime := getMCPStr(req, "end_time", "")
		if startTime == "" || endTime == "" {
			return mcp.NewToolResultError("start_time and end_time are required"), nil
		}
		req2 := &calendar.FreeBusyRequest{
			TimeMin: startTime,
			TimeMax: endTime,
			Items:   []*calendar.FreeBusyRequestItem{{Id: calendarID}},
		}
		resp, err := svc.Freebusy.Query(req2).Do()
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Calendar free/busy failed: %v", err)), nil
		}
		cals, ok := resp.Calendars[calendarID]
		if !ok {
			return mcp.NewToolResultText("📅 No calendar data found"), nil
		}
		if len(cals.Busy) == 0 {
			return mcp.NewToolResultText(fmt.Sprintf("✅ Free from %s to %s", startTime, endTime)), nil
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("📅 Busy from %s to %s:\n\n", startTime, endTime))
		for _, b := range cals.Busy {
			sb.WriteString(fmt.Sprintf("⛔ %s → %s\n", b.Start, b.End))
		}
		return mcp.NewToolResultText(sb.String()), nil
	}
}

// ─── Helpers ─────────────────────────────────────────────────────

func getGmailHeader(headers []*gmail.MessagePartHeader, name string) string {
	for _, h := range headers {
		if strings.EqualFold(h.Name, name) {
			return h.Value
		}
	}
	return ""
}

func extractEmailBody(payload *gmail.MessagePart) string {
	if payload == nil {
		return ""
	}
	if payload.Body.Data != "" {
		decoded, err := base64.URLEncoding.DecodeString(payload.Body.Data)
		if err == nil {
			return string(decoded)
		}
	}
	for _, part := range payload.Parts {
		if part.MimeType == "text/plain" && part.Body.Data != "" {
			decoded, err := base64.URLEncoding.DecodeString(part.Body.Data)
			if err == nil {
				return string(decoded)
			}
		}
	}
	return "(unable to extract body)"
}

func getMCPInt(req mcp.CallToolRequest, key string, fallback int) int {
	v, ok := req.GetArguments()[key]
	if !ok {
		return fallback
	}
	if f, ok := v.(float64); ok {
		return int(f)
	}
	return fallback
}

func getMCPStr(req mcp.CallToolRequest, key string, fallback string) string {
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

// envOr returns the env var value if non-empty, otherwise the fallback.
func envOr(envKey, fallback string) string {
	if v := os.Getenv(envKey); v != "" {
		return strings.TrimSpace(v)
	}
	return strings.TrimSpace(fallback)
}
