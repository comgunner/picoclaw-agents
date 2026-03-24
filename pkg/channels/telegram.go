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
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegohandler"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"

	"github.com/comgunner/picoclaw/pkg/bus"
	"github.com/comgunner/picoclaw/pkg/config"
	"github.com/comgunner/picoclaw/pkg/logger"
	"github.com/comgunner/picoclaw/pkg/utils"
	"github.com/comgunner/picoclaw/pkg/voice"
)

var (
	reHeading    = regexp.MustCompile(`^#{1,6}\s+(.+)$`)
	reBlockquote = regexp.MustCompile(`^>\s*(.*)$`)
	reLink       = regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)
	reBoldStar   = regexp.MustCompile(`\*\*(.+?)\*\*`)
	reBoldUnder  = regexp.MustCompile(`__(.+?)__`)
	reItalic     = regexp.MustCompile(`_([^_]+)_`)
	reStrike     = regexp.MustCompile(`~~(.+?)~~`)
	reListItem   = regexp.MustCompile(`^[-*]\s+`)
	reCodeBlock  = regexp.MustCompile("```[\\w]*\\n?([\\s\\S]*?)```")
	reInlineCode = regexp.MustCompile("`([^`]+)`")
)

type TelegramChannel struct {
	*BaseChannel
	bot          *telego.Bot
	commands     TelegramCommander
	config       *config.Config
	chatIDs      map[string]int64
	transcriber  *voice.GroqTranscriber
	placeholders sync.Map // chatID -> messageID
	stopThinking sync.Map // chatID -> thinkingCancel
}

type thinkingCancel struct {
	fn context.CancelFunc
}

func (c *thinkingCancel) Cancel() {
	if c != nil && c.fn != nil {
		c.fn()
	}
}

func NewTelegramChannel(cfg *config.Config, bus *bus.MessageBus) (*TelegramChannel, error) {
	var opts []telego.BotOption
	telegramCfg := cfg.Channels.Telegram

	if telegramCfg.Proxy != "" {
		proxyURL, parseErr := url.Parse(telegramCfg.Proxy)
		if parseErr != nil {
			return nil, fmt.Errorf("invalid proxy URL %q: %w", telegramCfg.Proxy, parseErr)
		}
		opts = append(opts, telego.WithHTTPClient(&http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyURL),
			},
		}))
	} else if os.Getenv("HTTP_PROXY") != "" || os.Getenv("HTTPS_PROXY") != "" {
		// Use environment proxy if configured
		opts = append(opts, telego.WithHTTPClient(&http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
			},
		}))
	}

	bot, err := telego.NewBot(telegramCfg.Token, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create telegram bot: %w", err)
	}

	base := NewBaseChannel("telegram", telegramCfg, bus, telegramCfg.AllowFrom)

	return &TelegramChannel{
		BaseChannel:  base,
		commands:     NewTelegramCommands(bot, cfg),
		bot:          bot,
		config:       cfg,
		chatIDs:      make(map[string]int64),
		transcriber:  nil,
		placeholders: sync.Map{},
		stopThinking: sync.Map{},
	}, nil
}

func (c *TelegramChannel) SetTranscriber(transcriber *voice.GroqTranscriber) {
	c.transcriber = transcriber
}

func (c *TelegramChannel) Start(ctx context.Context) error {
	logger.InfoC("telegram", "Starting Telegram bot (polling mode)...")

	updates, err := c.bot.UpdatesViaLongPolling(ctx, &telego.GetUpdatesParams{
		Timeout: 30,
	})
	if err != nil {
		return fmt.Errorf("failed to start long polling: %w", err)
	}

	bh, err := telegohandler.NewBotHandler(c.bot, updates)
	if err != nil {
		return fmt.Errorf("failed to create bot handler: %w", err)
	}

	// Register Bot Commands in Telegram UI
	_ = c.bot.SetMyCommands(ctx, &telego.SetMyCommandsParams{
		Commands: []telego.BotCommand{
			{Command: "start", Description: "Iniciar el bot"},
			{Command: "help", Description: "Ver lista de commandos"},
			{Command: "bundle_approve", Description: "Aprobar un lote de post+imagen (ej: /bundle_approve id=...)"},
			{Command: "bundle_regen", Description: "Regenerar un lote"},
			{Command: "bundle_edit", Description: "Editar el texto de un lote"},
			{Command: "bundle_cancel", Description: "Cancelar un lote y descartar"},
			{Command: "disable_sentinel", Description: "Desactivar sentinel (5m|15m|1h)"},
			{Command: "activate_sentinel", Description: "Activar sentinel"},
			{Command: "sentinel_status", Description: "Ver estado del sentinel"},
			{Command: "restrict_to_workspace", Description: "Controlar restricción de archivos (activate|deactivate|status)"},
			{Command: "show", Description: "Ver configuración actual"},
			{Command: "list", Description: "Listar opciones disponibles"},
			{Command: "status", Description: "Ver estado del contexto"},
		},
	})

	bh.HandleMessage(func(ctx *th.Context, message telego.Message) error {
		if !c.isMessageAllowed(&message) {
			return nil
		}
		c.commands.Help(ctx, message)
		return nil
	}, th.CommandEqual("help"))
	bh.HandleMessage(func(ctx *th.Context, message telego.Message) error {
		if !c.isMessageAllowed(&message) {
			return nil
		}
		return c.commands.Start(ctx, message)
	}, th.CommandEqual("start"))

	bh.HandleMessage(func(ctx *th.Context, message telego.Message) error {
		if !c.isMessageAllowed(&message) {
			return nil
		}
		return c.commands.Show(ctx, message)
	}, th.CommandEqual("show"))

	bh.HandleMessage(func(ctx *th.Context, message telego.Message) error {
		if !c.isMessageAllowed(&message) {
			return nil
		}
		return c.commands.List(ctx, message)
	}, th.CommandEqual("list"))

	bh.HandleMessage(func(ctx *th.Context, message telego.Message) error {
		return c.handleMessage(ctx, &message)
	}, th.AnyMessage())

	bh.HandleCallbackQuery(func(ctx *th.Context, query telego.CallbackQuery) error {
		// 1. Acknowledge the callback immediately to stop the loading spinner
		_ = c.bot.AnswerCallbackQuery(ctx.Context(), &telego.AnswerCallbackQueryParams{
			CallbackQueryID: query.ID,
		})

		// 2. Identify selection for visual feedback
		selection := query.Data
		emoji := "🔘"
		if strings.Contains(selection, "approve") || strings.Contains(selection, "yes") {
			emoji = "✅"
		} else if strings.Contains(selection, "reject") || strings.Contains(selection, "no") {
			emoji = "❌"
		}

		// 3. Edit the original message to remove buttons and show selection
		if query.Message != nil {
			msg, ok := query.Message.(*telego.Message)
			if ok {
				newText := msg.Text
				if newText == "" {
					newText = msg.Caption
				}
				// Append selection line
				newText += fmt.Sprintf("\n\n%s <b>Selección: %s</b>", emoji, selection)

				editMsg := tu.EditMessageText(tu.ID(msg.Chat.ID), msg.MessageID, newText)
				editMsg.ParseMode = telego.ModeHTML
				// No ReplyMarkup = buttons removed
				_, _ = c.bot.EditMessageText(ctx.Context(), editMsg)
			}
		}

		// 4. Publish to bus (as if it were a text message from user)
		senderID := fmt.Sprintf("%d", query.From.ID)
		if query.From.Username != "" {
			senderID = fmt.Sprintf("%d|%s", query.From.ID, query.From.Username)
		}

		c.HandleMessage(senderID, fmt.Sprintf("%d", query.From.ID), query.Data, nil, map[string]string{
			"is_callback": "true",
			"callback_id": query.ID,
		})

		return nil
	})

	c.setRunning(true)
	logger.InfoCF("telegram", "Telegram bot connected", map[string]any{
		"username": c.bot.Username(),
	})

	go bh.Start()

	go func() {
		<-ctx.Done()
		bh.Stop()
	}()

	return nil
}

func (c *TelegramChannel) Stop(ctx context.Context) error {
	logger.InfoC("telegram", "Stopping Telegram bot...")
	c.setRunning(false)
	return nil
}

func (c *TelegramChannel) Send(ctx context.Context, msg bus.OutboundMessage) error {
	if !c.IsRunning() {
		return fmt.Errorf("telegram bot not running")
	}

	chatID, err := parseChatID(msg.ChatID)
	if err != nil {
		return fmt.Errorf("invalid chat ID: %w", err)
	}

	// Stop thinking animation
	if stop, ok := c.stopThinking.Load(msg.ChatID); ok {
		if cf, ok := stop.(*thinkingCancel); ok && cf != nil {
			cf.Cancel()
		}
		c.stopThinking.Delete(msg.ChatID)
	}

	htmlContent := markdownToTelegramHTML(msg.Content)
	if len(htmlContent) > 4000 {
		logger.WarnCF("telegram", "Message too long, truncating", map[string]any{"len": len(htmlContent)})
		htmlContent = htmlContent[:4000] + "\n\n...(truncated)"
	}

	// Prepare buttons if any
	var keyboard *telego.InlineKeyboardMarkup
	if len(msg.Buttons) > 0 {
		var row []telego.InlineKeyboardButton
		for _, btn := range msg.Buttons {
			row = append(row, tu.InlineKeyboardButton(btn.Text).WithCallbackData(btn.Data))
		}
		keyboard = tu.InlineKeyboard(row)
	}

	// If there are media files, handle them
	if len(msg.Media) > 0 {
		caption := htmlContent
		sendTextSeparately := false

		// Check if content has approval buttons
		hasApprovalButtons := false
		for _, btn := range msg.Buttons {
			if strings.Contains(btn.Data, "approve_") || strings.Contains(btn.Data, "reject_") ||
				strings.Contains(btn.Data, "edit_") || strings.Contains(btn.Text, "✅") ||
				strings.Contains(btn.Text, "❌") || strings.Contains(btn.Text, "✏️") {
				hasApprovalButtons = true
				break
			}
		}

		// Always send text separately if it has approval buttons OR if caption is too long
		if hasApprovalButtons || len(caption) > 1024 {
			sendTextSeparately = true
		}

		if sendTextSeparately {
			tgMsg := tu.Message(tu.ID(chatID), htmlContent)
			tgMsg.ParseMode = telego.ModeHTML
			if keyboard != nil {
				tgMsg.ReplyMarkup = keyboard
			}
			_, err = c.bot.SendMessage(ctx, tgMsg)
			if err != nil {
				logger.ErrorCF("telegram", "SendMessage HTML failed, retrying plain text (markdown)", map[string]any{
					"error": err.Error(),
					"chat":  chatID,
				})
				plainMsg := tu.Message(tu.ID(chatID), msg.Content)
				if keyboard != nil {
					plainMsg.ReplyMarkup = keyboard
				}
				_, err = c.bot.SendMessage(ctx, plainMsg)
				if err != nil {
					logger.ErrorCF(
						"telegram",
						"SendMessage plain text fallback failed too",
						map[string]any{"error": err.Error()},
					)
				}
			}
			keyboard = nil
			caption = ""
		}

		// Optimización: Agrupar fotos en un MediaGroup si hay más de una
		var mediaGroup []telego.InputMedia
		isAllPhotos := true

		for _, mediaPath := range msg.Media {
			ext := strings.ToLower(filepath.Ext(mediaPath))
			isPhoto := ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".webp"
			if !isPhoto {
				isAllPhotos = false
				break
			}
		}

		if isAllPhotos && len(msg.Media) > 1 {
			// Enviar como grupo
			for i, path := range msg.Media {
				f, err := os.Open(path)
				if err != nil {
					continue
				}
				defer f.Close()

				photo := &telego.InputMediaPhoto{
					Type:      telego.MediaTypePhoto,
					Media:     tu.File(f),
					ParseMode: telego.ModeHTML,
				}
				if i == 0 {
					photo.Caption = caption
				}
				mediaGroup = append(mediaGroup, photo)
			}
			_, err = c.bot.SendMediaGroup(ctx, tu.MediaGroup(tu.ID(chatID), mediaGroup...))
			return err
		}

		// Fallback o envío individual (si no son fotos o es solo una)
		for i, mediaPath := range msg.Media {
			currentCaption := ""
			if i == 0 {
				currentCaption = caption
			}

			ext := strings.ToLower(filepath.Ext(mediaPath))
			isPhoto := ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".webp"

			f, err := os.Open(mediaPath)
			if err != nil {
				return fmt.Errorf("opening file %s: %w", mediaPath, err)
			}
			defer f.Close()

			if isPhoto {
				tgPhoto := tu.Photo(tu.ID(chatID), tu.File(f))
				tgPhoto.Caption = currentCaption
				tgPhoto.ParseMode = telego.ModeHTML
				if !sendTextSeparately && keyboard != nil && i == len(msg.Media)-1 {
					tgPhoto.ReplyMarkup = keyboard
				}
				_, err = c.bot.SendPhoto(ctx, tgPhoto)
			} else {
				tgDoc := tu.Document(tu.ID(chatID), tu.File(f))
				tgDoc.Caption = currentCaption
				tgDoc.ParseMode = telego.ModeHTML
				if !sendTextSeparately && keyboard != nil && i == len(msg.Media)-1 {
					tgDoc.ReplyMarkup = keyboard
				}
				_, err = c.bot.SendDocument(ctx, tgDoc)
			}
			if err != nil {
				return err
			}
		}
		return nil
	}

	// Try to edit placeholder
	if pID, ok := c.placeholders.Load(msg.ChatID); ok {
		c.placeholders.Delete(msg.ChatID)
		editMsg := tu.EditMessageText(tu.ID(chatID), pID.(int), htmlContent)
		editMsg.ParseMode = telego.ModeHTML
		if keyboard != nil {
			editMsg.ReplyMarkup = keyboard
		}

		if _, err = c.bot.EditMessageText(ctx, editMsg); err != nil {
			logger.ErrorCF("telegram", "EditMessageText HTML failed, retrying plain text (markdown)", map[string]any{
				"error": err.Error(),
			})
			editMsg.Text = msg.Content // Use original markdown
			editMsg.ParseMode = ""
			if _, err = c.bot.EditMessageText(ctx, editMsg); err == nil {
				return nil
			}
		} else {
			return nil
		}
		// Fallback to new message if edit fails
	}

	tgMsg := tu.Message(tu.ID(chatID), htmlContent)
	tgMsg.ParseMode = telego.ModeHTML
	if keyboard != nil {
		tgMsg.ReplyMarkup = keyboard
	}

	if _, err = c.bot.SendMessage(ctx, tgMsg); err != nil {
		logger.ErrorCF("telegram", "HTML parse failed, falling back to plain text", map[string]any{
			"error": err.Error(),
		})
		tgMsg.ParseMode = ""
		_, err = c.bot.SendMessage(ctx, tgMsg)
		return err
	}

	return nil
}

func (c *TelegramChannel) isMessageAllowed(message *telego.Message) bool {
	if message == nil {
		return false
	}
	user := message.From
	if user == nil {
		return false
	}

	senderID := fmt.Sprintf("%d", user.ID)
	if user.Username != "" {
		senderID = fmt.Sprintf("%d|%s", user.ID, user.Username)
	}

	if !c.IsAllowed(senderID) {
		logger.DebugCF("telegram", "Command rejected by allowlist", map[string]any{
			"user_id": senderID,
		})
		return false
	}

	return true
}

func (c *TelegramChannel) handleMessage(ctx context.Context, message *telego.Message) error {
	if !c.isMessageAllowed(message) {
		return nil
	}

	user := message.From
	senderID := fmt.Sprintf("%d", user.ID)
	if user.Username != "" {
		senderID = fmt.Sprintf("%d|%s", user.ID, user.Username)
	}

	chatID := message.Chat.ID
	c.chatIDs[senderID] = chatID

	content := ""
	mediaPaths := []string{}
	localFiles := []string{} // track local files that need cleanup

	// ensure temp files are cleaned up when function returns
	defer func() {
		for _, file := range localFiles {
			if err := os.Remove(file); err != nil {
				logger.DebugCF("telegram", "Failed to cleanup temp file", map[string]any{
					"file":  file,
					"error": err.Error(),
				})
			}
		}
	}()

	if message.Text != "" {
		content += message.Text
	}

	if message.Caption != "" {
		if content != "" {
			content += "\n"
		}
		content += message.Caption
	}

	if len(message.Photo) > 0 {
		photo := message.Photo[len(message.Photo)-1]
		photoPath := c.downloadPhoto(ctx, photo.FileID)
		if photoPath != "" {
			localFiles = append(localFiles, photoPath)
			mediaPaths = append(mediaPaths, photoPath)
			if content != "" {
				content += "\n"
			}
			content += "[image: photo]"
		}
	}

	if message.Voice != nil {
		voicePath := c.downloadFile(ctx, message.Voice.FileID, ".ogg")
		if voicePath != "" {
			localFiles = append(localFiles, voicePath)
			mediaPaths = append(mediaPaths, voicePath)

			var transcribedText string
			if c.transcriber != nil && c.transcriber.IsAvailable() {
				transcriberCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
				defer cancel()

				result, err := c.transcriber.Transcribe(transcriberCtx, voicePath)
				if err != nil {
					logger.ErrorCF("telegram", "Voice transcription failed", map[string]any{
						"error": err.Error(),
						"path":  voicePath,
					})
					transcribedText = "[voice (transcription failed)]"
				} else {
					transcribedText = fmt.Sprintf("[voice transcription: %s]", result.Text)
					logger.InfoCF("telegram", "Voice transcribed successfully", map[string]any{
						"text": result.Text,
					})
				}
			} else {
				transcribedText = "[voice]"
			}

			if content != "" {
				content += "\n"
			}
			content += transcribedText
		}
	}

	if message.Audio != nil {
		audioPath := c.downloadFile(ctx, message.Audio.FileID, ".mp3")
		if audioPath != "" {
			localFiles = append(localFiles, audioPath)
			mediaPaths = append(mediaPaths, audioPath)
			if content != "" {
				content += "\n"
			}
			content += "[audio]"
		}
	}

	if message.Document != nil {
		docPath := c.downloadFile(ctx, message.Document.FileID, "")
		if docPath != "" {
			localFiles = append(localFiles, docPath)
			mediaPaths = append(mediaPaths, docPath)
			if content != "" {
				content += "\n"
			}
			content += "[file]"
		}
	}

	if content == "" {
		content = "[empty message]"
	}

	logger.DebugCF("telegram", "Received message", map[string]any{
		"sender_id": senderID,
		"chat_id":   fmt.Sprintf("%d", chatID),
		"preview":   utils.Truncate(content, 50),
	})

	// Thinking indicator
	err := c.bot.SendChatAction(ctx, tu.ChatAction(tu.ID(chatID), telego.ChatActionTyping))
	if err != nil {
		logger.ErrorCF("telegram", "Failed to send chat action", map[string]any{
			"error": err.Error(),
		})
	}

	// Stop any previous thinking animation
	chatIDStr := fmt.Sprintf("%d", chatID)
	if prevStop, ok := c.stopThinking.Load(chatIDStr); ok {
		if cf, ok := prevStop.(*thinkingCancel); ok && cf != nil {
			cf.Cancel()
		}
	}

	// Create cancel function for thinking state
	_, thinkCancel := context.WithTimeout(ctx, 5*time.Minute)
	c.stopThinking.Store(chatIDStr, &thinkingCancel{fn: thinkCancel})

	// Skip placeholder for fast-path commands
	isCommand := strings.HasPrefix(content, "/") || strings.HasPrefix(content, "#")

	if !isCommand {
		pMsg, err := c.bot.SendMessage(ctx, tu.Message(tu.ID(chatID), "Thinking... 💭"))
		if err == nil {
			pID := pMsg.MessageID
			c.placeholders.Store(chatIDStr, pID)
		}
	}

	peerKind := "direct"
	peerID := fmt.Sprintf("%d", user.ID)
	if message.Chat.Type != "private" {
		peerKind = "group"
		peerID = fmt.Sprintf("%d", chatID)
	}

	metadata := map[string]string{
		"message_id": fmt.Sprintf("%d", message.MessageID),
		"user_id":    fmt.Sprintf("%d", user.ID),
		"username":   user.Username,
		"first_name": user.FirstName,
		"is_group":   fmt.Sprintf("%t", message.Chat.Type != "private"),
		"peer_kind":  peerKind,
		"peer_id":    peerID,
	}

	c.HandleMessage(fmt.Sprintf("%d", user.ID), fmt.Sprintf("%d", chatID), content, mediaPaths, metadata)
	return nil
}

func (c *TelegramChannel) downloadPhoto(ctx context.Context, fileID string) string {
	file, err := c.bot.GetFile(ctx, &telego.GetFileParams{FileID: fileID})
	if err != nil {
		logger.ErrorCF("telegram", "Failed to get photo file", map[string]any{
			"error": err.Error(),
		})
		return ""
	}

	return c.downloadFileWithInfo(file, ".jpg")
}

func (c *TelegramChannel) downloadFileWithInfo(file *telego.File, ext string) string {
	if file.FilePath == "" {
		return ""
	}

	url := c.bot.FileDownloadURL(file.FilePath)
	logger.DebugCF("telegram", "File URL", map[string]any{"url": url})

	// Use FilePath as filename for better identification
	filename := file.FilePath + ext
	return utils.DownloadFile(url, filename, utils.DownloadOptions{
		LoggerPrefix: "telegram",
	})
}

func (c *TelegramChannel) downloadFile(ctx context.Context, fileID, ext string) string {
	file, err := c.bot.GetFile(ctx, &telego.GetFileParams{FileID: fileID})
	if err != nil {
		logger.ErrorCF("telegram", "Failed to get file", map[string]any{
			"error": err.Error(),
		})
		return ""
	}

	return c.downloadFileWithInfo(file, ext)
}

func parseChatID(chatIDStr string) (int64, error) {
	var id int64
	_, err := fmt.Sscanf(chatIDStr, "%d", &id)
	return id, err
}

func markdownToTelegramHTML(text string) string {
	if text == "" {
		return ""
	}

	// 1. Preserve code blocks
	codeBlocks := extractCodeBlocks(text)
	text = codeBlocks.text

	// 2. Preserve inline codes
	inlineCodes := extractInlineCodes(text)
	text = inlineCodes.text

	// 3. Simple block replacements
	text = reHeading.ReplaceAllString(text, "$1")
	text = reBlockquote.ReplaceAllString(text, "$1")
	text = reListItem.ReplaceAllString(text, "• ")

	// 4. Escape general content BEFORE adding HTML tags
	text = escapeHTML(text)

	// 5. Handle Links
	text = reLink.ReplaceAllString(text, `<a href="$2">$1</a>`)

	// 6. Handle Bold (Order matters: Bold before Italic to avoid ambiguity)
	text = reBoldStar.ReplaceAllString(text, "<b>$1</b>")
	text = reBoldUnder.ReplaceAllString(text, "<b>$1</b>")

	// 7. Handle Italic (Strict check to avoid matching inside existing tags)
	// We use a simpler strategy: only replace if not already part of an HTML tag
	text = reItalic.ReplaceAllStringFunc(text, func(s string) string {
		match := reItalic.FindStringSubmatch(s)
		if len(match) < 2 {
			return s
		}
		// If it looks like it's inside a tag (very basic check), skip
		if strings.Contains(match[1], "<") || strings.Contains(match[1], ">") {
			return s
		}
		return "<i>" + match[1] + "</i>"
	})

	// 8. Handle Strikethrough
	text = reStrike.ReplaceAllString(text, "<s>$1</s>")

	// 9. Restore Inline Codes (already escaped)
	for i, code := range inlineCodes.codes {
		escaped := escapeHTML(code)
		text = strings.ReplaceAll(text, fmt.Sprintf("\x00IC%d\x00", i), fmt.Sprintf("<code>%s</code>", escaped))
	}

	// 10. Restore Code Blocks (already escaped)
	for i, code := range codeBlocks.codes {
		escaped := escapeHTML(code)
		text = strings.ReplaceAll(
			text,
			fmt.Sprintf("\x00CB%d\x00", i),
			fmt.Sprintf("<pre><code>%s</code></pre>", escaped),
		)
	}

	return text
}

type codeBlockMatch struct {
	text  string
	codes []string
}

func extractCodeBlocks(text string) codeBlockMatch {
	matches := reCodeBlock.FindAllStringSubmatch(text, -1)

	codes := make([]string, 0, len(matches))
	for _, match := range matches {
		codes = append(codes, match[1])
	}

	i := 0
	text = reCodeBlock.ReplaceAllStringFunc(text, func(m string) string {
		placeholder := fmt.Sprintf("\x00CB%d\x00", i)
		i++
		return placeholder
	})

	return codeBlockMatch{text: text, codes: codes}
}

type inlineCodeMatch struct {
	text  string
	codes []string
}

func extractInlineCodes(text string) inlineCodeMatch {
	matches := reInlineCode.FindAllStringSubmatch(text, -1)

	codes := make([]string, 0, len(matches))
	for _, match := range matches {
		codes = append(codes, match[1])
	}

	i := 0
	text = reInlineCode.ReplaceAllStringFunc(text, func(m string) string {
		placeholder := fmt.Sprintf("\x00IC%d\x00", i)
		i++
		return placeholder
	})

	return inlineCodeMatch{text: text, codes: codes}
}

func escapeHTML(text string) string {
	text = strings.ReplaceAll(text, "&", "&amp;")
	text = strings.ReplaceAll(text, "<", "&lt;")
	text = strings.ReplaceAll(text, ">", "&gt;")
	return text
}
