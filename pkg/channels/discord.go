// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package channels

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/comgunner/picoclaw/pkg/bus"
	"github.com/comgunner/picoclaw/pkg/commands"
	"github.com/comgunner/picoclaw/pkg/config"
	"github.com/comgunner/picoclaw/pkg/logger"
	"github.com/comgunner/picoclaw/pkg/utils"
	"github.com/comgunner/picoclaw/pkg/voice"
)

const (
	transcriptionTimeout = 30 * time.Second
	sendTimeout          = 10 * time.Second
)

type DiscordChannel struct {
	*BaseChannel
	session      *discordgo.Session
	config       config.DiscordConfig
	transcriber  *voice.GroqTranscriber
	ctx          context.Context
	typingMu     sync.Mutex
	typingStop   map[string]chan struct{} // chatID → stop signal
	botUserID    string                   // stored for mention checking
	modelHandler *commands.ModelCommandHandler
}

func NewDiscordChannel(
	cfg config.DiscordConfig,
	bus *bus.MessageBus,
	fullConfig *config.Config,
) (*DiscordChannel, error) {
	session, err := discordgo.New("Bot " + cfg.Token)
	if err != nil {
		return nil, fmt.Errorf("failed to create discord session: %w", err)
	}

	base := NewBaseChannel("discord", cfg, bus, cfg.AllowFrom)

	// Initialize model command handler
	var modelHandler *commands.ModelCommandHandler
	if fullConfig != nil {
		configPath := filepath.Join(fullConfig.WorkspacePath(), "..", "config.json")
		modelHandler = commands.NewModelCommandHandler(configPath, nil)
	} else {
		modelHandler = commands.NewModelCommandHandler("config.json", nil)
	}

	return &DiscordChannel{
		BaseChannel:  base,
		session:      session,
		config:       cfg,
		transcriber:  nil,
		ctx:          context.Background(),
		typingStop:   make(map[string]chan struct{}),
		modelHandler: modelHandler,
	}, nil
}

func (c *DiscordChannel) SetTranscriber(transcriber *voice.GroqTranscriber) {
	c.transcriber = transcriber
}

func (c *DiscordChannel) getContext() context.Context {
	if c.ctx == nil {
		return context.Background()
	}
	return c.ctx
}

func (c *DiscordChannel) Start(ctx context.Context) error {
	logger.InfoC("discord", "Starting Discord bot")

	c.ctx = ctx

	// Get bot user ID before opening session to avoid race condition
	botUser, err := c.session.User("@me")
	if err != nil {
		return fmt.Errorf("failed to get bot user: %w", err)
	}
	c.botUserID = botUser.ID

	c.session.AddHandler(c.handleMessage)
	c.session.AddHandler(c.handleInteraction)

	if err := c.session.Open(); err != nil {
		return fmt.Errorf("failed to open discord session: %w", err)
	}

	// Register Slash Commands
	logger.InfoC("discord", "Registering slash commands...")
	_, err = c.session.ApplicationCommandBulkOverwrite(c.botUserID, "", DiscordCommands)
	if err != nil {
		logger.ErrorCF("discord", "Failed to register commands", map[string]any{"error": err.Error()})
	}

	c.setRunning(true)

	logger.InfoCF("discord", "Discord bot connected", map[string]any{
		"username": botUser.Username,
		"user_id":  botUser.ID,
	})

	return nil
}

func (c *DiscordChannel) Stop(ctx context.Context) error {
	logger.InfoC("discord", "Stopping Discord bot")
	c.setRunning(false)

	// Stop all typing goroutines before closing session
	c.typingMu.Lock()
	for chatID, stop := range c.typingStop {
		close(stop)
		delete(c.typingStop, chatID)
	}
	c.typingMu.Unlock()

	if err := c.session.Close(); err != nil {
		return fmt.Errorf("failed to close discord session: %w", err)
	}

	return nil
}

func (c *DiscordChannel) Send(ctx context.Context, msg bus.OutboundMessage) error {
	c.stopTyping(msg.ChatID)

	if !c.IsRunning() {
		return fmt.Errorf("discord bot not running")
	}

	channelID := msg.ChatID
	if channelID == "" {
		return fmt.Errorf("channel ID is empty")
	}

	// Prepare files if any
	var files []*discordgo.File
	for _, m := range msg.Media {
		f, err := os.Open(m)
		if err != nil {
			logger.ErrorCF("discord", "Failed to open media file", map[string]any{"path": m, "error": err.Error()})
			continue
		}
		defer f.Close()
		files = append(files, &discordgo.File{
			Name:        filepath.Base(m),
			ContentType: "", // Discord will detect
			Reader:      f,
		})
	}

	chunks := utils.SplitMessage(msg.Content, 2000)
	if len(chunks) == 0 && len(files) > 0 {
		chunks = []string{""} // Send files with empty content
	}

	for i, chunk := range chunks {
		var currentFiles []*discordgo.File
		var currentComponents []discordgo.MessageComponent
		if i == 0 {
			currentFiles = files // Attach files only to the first chunk
		}
		if i == len(chunks)-1 && len(msg.Buttons) > 0 {
			// Add buttons to the last chunk
			var actionRow discordgo.ActionsRow
			for _, btn := range msg.Buttons {
				actionRow.Components = append(actionRow.Components, &discordgo.Button{
					Label:    btn.Text,
					Style:    discordgo.PrimaryButton,
					CustomID: btn.Data,
				})
			}
			currentComponents = append(currentComponents, actionRow)
		}

		if err := c.sendChunk(ctx, channelID, chunk, currentFiles, currentComponents); err != nil {
			return err
		}
	}

	return nil
}

func (c *DiscordChannel) sendChunk(
	ctx context.Context,
	channelID, content string,
	files []*discordgo.File,
	components []discordgo.MessageComponent,
) error {
	sendCtx, cancel := context.WithTimeout(ctx, sendTimeout)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		var err error
		if len(files) > 0 || len(components) > 0 {
			_, err = c.session.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
				Content:    content,
				Files:      files,
				Components: components,
			})
		} else {
			_, err = c.session.ChannelMessageSend(channelID, content)
		}
		done <- err
	}()

	select {
	case err := <-done:
		if err != nil {
			return fmt.Errorf("failed to send discord message: %w", err)
		}
		return nil
	case <-sendCtx.Done():
		return fmt.Errorf("send message timeout: %w", sendCtx.Err())
	}
}

// appendContent safely appends content to existing text
func appendContent(content, suffix string) string {
	if content == "" {
		return suffix
	}
	return content + "\n" + suffix
}

func (c *DiscordChannel) handleMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m == nil || m.Author == nil {
		return
	}

	if m.Author.ID == s.State.User.ID {
		return
	}

	// Check allowlist first to avoid downloading attachments and transcribing for rejected users
	if !c.IsAllowed(m.Author.ID) {
		logger.DebugCF("discord", "Message rejected by allowlist", map[string]any{
			"user_id": m.Author.ID,
		})
		return
	}

	// If configured to only respond to mentions, check if bot is mentioned
	// Skip this check for DMs (GuildID is empty) - DMs should always be responded to
	if c.config.MentionOnly && m.GuildID != "" {
		isMentioned := false
		for _, mention := range m.Mentions {
			if mention.ID == c.botUserID {
				isMentioned = true
				break
			}
		}
		if !isMentioned {
			logger.DebugCF("discord", "Message ignored - bot not mentioned", map[string]any{
				"user_id": m.Author.ID,
			})
			return
		}
	}

	senderID := m.Author.ID
	senderName := m.Author.Username
	if m.Author.Discriminator != "" && m.Author.Discriminator != "0" {
		senderName += "#" + m.Author.Discriminator
	}

	content := m.Content
	content = c.stripBotMention(content)
	mediaPaths := make([]string, 0, len(m.Attachments))
	localFiles := make([]string, 0, len(m.Attachments))

	// Ensure temp files are cleaned up when function returns
	defer func() {
		for _, file := range localFiles {
			if err := os.Remove(file); err != nil {
				logger.DebugCF("discord", "Failed to cleanup temp file", map[string]any{
					"file":  file,
					"error": err.Error(),
				})
			}
		}
	}()

	for _, attachment := range m.Attachments {
		isAudio := utils.IsAudioFile(attachment.Filename, attachment.ContentType)

		if isAudio {
			localPath := c.downloadAttachment(attachment.URL, attachment.Filename)
			if localPath != "" {
				localFiles = append(localFiles, localPath)

				var transcribedText string
				if c.transcriber != nil && c.transcriber.IsAvailable() {
					ctx, cancel := context.WithTimeout(c.getContext(), transcriptionTimeout)
					result, err := c.transcriber.Transcribe(ctx, localPath)
					cancel() // Release context resources immediately to avoid leaks in for loop

					if err != nil {
						logger.ErrorCF("discord", "Voice transcription failed", map[string]any{
							"error": err.Error(),
						})
						transcribedText = fmt.Sprintf("[audio: %s (transcription failed)]", attachment.Filename)
					} else {
						transcribedText = fmt.Sprintf("[audio transcription: %s]", result.Text)
						logger.DebugCF("discord", "Audio transcribed successfully", map[string]any{
							"text": result.Text,
						})
					}
				} else {
					transcribedText = fmt.Sprintf("[audio: %s]", attachment.Filename)
				}

				content = appendContent(content, transcribedText)
			} else {
				logger.WarnCF("discord", "Failed to download audio attachment", map[string]any{
					"url":      attachment.URL,
					"filename": attachment.Filename,
				})
				mediaPaths = append(mediaPaths, attachment.URL)
				content = appendContent(content, fmt.Sprintf("[attachment: %s]", attachment.URL))
			}
		} else {
			mediaPaths = append(mediaPaths, attachment.URL)
			content = appendContent(content, fmt.Sprintf("[attachment: %s]", attachment.URL))
		}
	}

	if content == "" && len(mediaPaths) == 0 {
		return
	}

	if content == "" {
		content = "[media only]"
	}

	// Start typing after all early returns — guaranteed to have a matching Send()
	c.startTyping(m.ChannelID)

	logger.DebugCF("discord", "Received message", map[string]any{
		"sender_name": senderName,
		"sender_id":   senderID,
		"preview":     utils.Truncate(content, 50),
	})

	peerKind := "channel"
	peerID := m.ChannelID
	if m.GuildID == "" {
		peerKind = "direct"
		peerID = senderID
	}

	metadata := map[string]string{
		"message_id":   m.ID,
		"user_id":      senderID,
		"username":     m.Author.Username,
		"display_name": senderName,
		"guild_id":     m.GuildID,
		"channel_id":   m.ChannelID,
		"is_dm":        fmt.Sprintf("%t", m.GuildID == ""),
		"peer_kind":    peerKind,
		"peer_id":      peerID,
	}

	c.HandleMessage(senderID, m.ChannelID, content, mediaPaths, metadata)
}

// startTyping starts a continuous typing indicator loop for the given chatID.
// It stops any existing typing loop for that chatID before starting a new one.
func (c *DiscordChannel) startTyping(chatID string) {
	c.typingMu.Lock()
	// Stop existing loop for this chatID if any
	if stop, ok := c.typingStop[chatID]; ok {
		close(stop)
	}
	stop := make(chan struct{})
	c.typingStop[chatID] = stop
	c.typingMu.Unlock()

	go func() {
		if err := c.session.ChannelTyping(chatID); err != nil {
			logger.DebugCF("discord", "ChannelTyping error", map[string]any{"chatID": chatID, "err": err})
		}
		ticker := time.NewTicker(8 * time.Second)
		defer ticker.Stop()
		timeout := time.After(5 * time.Minute)
		for {
			select {
			case <-stop:
				return
			case <-timeout:
				return
			case <-c.ctx.Done():
				return
			case <-ticker.C:
				if err := c.session.ChannelTyping(chatID); err != nil {
					logger.DebugCF("discord", "ChannelTyping error", map[string]any{"chatID": chatID, "err": err})
				}
			}
		}
	}()
}

// stopTyping stops the typing indicator loop for the given chatID.
func (c *DiscordChannel) stopTyping(chatID string) {
	c.typingMu.Lock()
	defer c.typingMu.Unlock()
	if stop, ok := c.typingStop[chatID]; ok {
		close(stop)
		delete(c.typingStop, chatID)
	}
}

func (c *DiscordChannel) downloadAttachment(url, filename string) string {
	return utils.DownloadFile(url, filename, utils.DownloadOptions{
		LoggerPrefix: "discord",
	})
}

// stripBotMention removes the bot mention from the message content.
// Discord mentions have the format <@USER_ID> or <@!USER_ID> (with nickname).
func (c *DiscordChannel) stripBotMention(text string) string {
	if c.botUserID == "" {
		return text
	}
	// Remove both regular mention <@USER_ID> and nickname mention <@!USER_ID>
	text = strings.ReplaceAll(text, fmt.Sprintf("<@%s>", c.botUserID), "")
	text = strings.ReplaceAll(text, fmt.Sprintf("<@!%s>", c.botUserID), "")
	return strings.TrimSpace(text)
}

func (c *DiscordChannel) handleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	data := i.ApplicationCommandData()
	cmdName := data.Name
	options := make(map[string]any)
	for _, opt := range data.Options {
		options[opt.Name] = opt.Value
	}

	// Immediate Acknowledgment
	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("⚙️ Ejecutando `/%s`...", cmdName),
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})

	// Convert to standard text command
	var fullCmd strings.Builder
	fullCmd.WriteString("/")
	fullCmd.WriteString(cmdName)

	if id, ok := options["id"].(string); ok {
		fullCmd.WriteString(" id=")
		fullCmd.WriteString(id)
	}
	if target, ok := options["target"].(string); ok {
		fullCmd.WriteString(" ")
		fullCmd.WriteString(target)
	}

	content := fullCmd.String()
	senderID := ""
	if i.Member != nil && i.Member.User != nil {
		senderID = i.Member.User.ID
	} else if i.User != nil {
		senderID = i.User.ID
	}

	peerKind := "channel"
	peerID := i.ChannelID
	if i.GuildID == "" {
		peerKind = "direct"
		peerID = senderID
	}

	metadata := map[string]string{
		"interaction_id": i.ID,
		"user_id":        senderID,
		"guild_id":       i.GuildID,
		"channel_id":     i.ChannelID,
		"is_interaction": "true",
		"peer_kind":      peerKind,
		"peer_id":        peerID,
	}

	logger.InfoCF("discord", "Slash command interaction", map[string]any{
		"command":   cmdName,
		"sender_id": senderID,
	})

	// Handle Sentinel Commands Directly (Fast-Path)
	switch cmdName {
	case "disable_sentinel":
		duration := ""
		for _, opt := range data.Options {
			if opt.Name == "duration" {
				duration = opt.StringValue()
				break
			}
		}
		if duration == "" {
			c.session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "⚠️ Usage: /disable_sentinel [5m|15m|1h]",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			return
		}

		// Forward to AgentLoop for processing
		content := "/disable_sentinel " + duration
		c.HandleMessage(senderID, i.ChannelID, content, nil, metadata)
		return

	case "activate_sentinel":
		content := "/activate_sentinel"
		c.HandleMessage(senderID, i.ChannelID, content, nil, metadata)
		return

	case "sentinel_status":
		content := "/sentinel_status"
		c.HandleMessage(senderID, i.ChannelID, content, nil, metadata)
		return

	case "restrict_to_workspace":
		action := ""
		for _, opt := range data.Options {
			if opt.Name == "action" {
				action = opt.StringValue()
				break
			}
		}
		if action == "" {
			c.session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "⚠️ Usage: /restrict_to_workspace [activate|deactivate|status]",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			return
		}
		content := "/restrict_to_workspace " + action
		c.HandleMessage(senderID, i.ChannelID, content, nil, metadata)
		return

	case "model":
		// Handle /model command as fast-path (don't pass to LLM)
		var modelName string
		for _, opt := range data.Options {
			if opt.Name == "model_name" {
				modelName = opt.StringValue()
				break
			}
		}

		// Build command string
		commandText := "/model"
		if modelName != "" {
			commandText += " " + modelName
		}

		// Process with ModelCommandHandler (fast-path)
		if c.modelHandler != nil {
			response, _ := c.modelHandler.Handle(commandText, senderID)
			// Send response as ephemeral message
			_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: response,
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
		}
		return
	}

	c.HandleMessage(senderID, i.ChannelID, content, nil, metadata)
}
