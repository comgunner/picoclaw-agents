// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package tools

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/comgunner/picoclaw/pkg/utils"
)

// ============== Facebook Tools ==============

type FacebookPostTool struct {
	defaultPageID    string
	defaultPageToken string
	appID            string
	appSecret        string
	userToken        string
}

func NewFacebookPostTool() *FacebookPostTool {
	return NewFacebookPostToolFromConfig("", "", "", "", "")
}

func NewFacebookPostToolFromConfig(
	configPageID, configPageToken, configAppID, configAppSecret, configUserToken string,
) *FacebookPostTool {
	pageID := strings.TrimSpace(os.Getenv(utils.EnvSocialFacebookPageID))
	pageToken := strings.TrimSpace(os.Getenv(utils.EnvSocialFacebookPageToken))
	appID := strings.TrimSpace(os.Getenv(utils.EnvSocialFacebookAppID))
	appSecret := strings.TrimSpace(os.Getenv(utils.EnvSocialFacebookAppSecret))
	userToken := strings.TrimSpace(os.Getenv(utils.EnvSocialFacebookUserToken))
	if pageID == "" {
		pageID = strings.TrimSpace(configPageID)
	}
	if pageToken == "" {
		pageToken = strings.TrimSpace(configPageToken)
	}
	if appID == "" {
		appID = strings.TrimSpace(configAppID)
	}
	if appSecret == "" {
		appSecret = strings.TrimSpace(configAppSecret)
	}
	if userToken == "" {
		userToken = strings.TrimSpace(configUserToken)
	}
	return &FacebookPostTool{
		defaultPageID:    pageID,
		defaultPageToken: pageToken,
		appID:            appID,
		appSecret:        appSecret,
		userToken:        userToken,
	}
}

func (t *FacebookPostTool) Name() string {
	return "facebook_post"
}

func (t *FacebookPostTool) Description() string {
	return "Publicar en Facebook (solo texto o imagen + mensaje, comentario opcional). Soporta multi-página y reintento por expiración de token cuando hay app_id/app_secret/user_token."
}

func (t *FacebookPostTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"page_id": map[string]any{
				"type":        "string",
				"description": "Facebook Page ID. Si se omite, usa el default del config.",
			},
			"page_token": map[string]any{
				"type":        "string",
				"description": "Facebook Page Access Token. Si se omite, usa el default del config.",
			},
			"message": map[string]any{
				"type":        "string",
				"description": "Mensaje principal de la publicación",
			},
			"image_path": map[string]any{
				"type":        "string",
				"description": "Ruta absoluta de la imagen a publicar (opcional). Si no se envía, publica solo texto.",
			},
			"comment": map[string]any{
				"type":        "string",
				"description": "Comentario opcional a agregar después de publicar",
			},
			"confirmed_by_user": map[string]any{
				"type":        "boolean",
				"description": "REQUERIDO: Debes setear esto a true SOLO después de haber mostrado el borrador al usuario con la herramienta 'message' y haber recibido su aprobación. NUNCA publiques sin aprobación previa.",
			},
		},
		"required": []string{"message", "confirmed_by_user"},
	}
}

func (t *FacebookPostTool) Execute(ctx context.Context, args map[string]any) *ToolResult {
	pageID, _ := args["page_id"].(string)
	pageToken, _ := args["page_token"].(string)
	message, _ := args["message"].(string)
	imagePath, _ := args["image_path"].(string)
	comment, _ := args["comment"].(string)

	// Usar defaults si no se proporcionan
	if pageID == "" {
		pageID = t.defaultPageID
	}
	if pageToken == "" {
		pageToken = t.defaultPageToken
	}

	// Validar parámetros requeridos
	message = strings.TrimSpace(message)
	if message == "" {
		return ErrorResult("message es requerido")
	}

	imagePath = strings.TrimSpace(imagePath)

	confirmed, _ := args["confirmed_by_user"].(bool)
	if !confirmed {
		return UserResult(
			"❌ **ERROR DE SEGURIDAD**: No puedes publicar directamente sin aprobación. " +
				"PASOS REQUERIDOS:\n" +
				"1. Usa la herramienta 'message' para mostrar el texto y la imagen al usuario.\n" +
				"2. Pide permiso explícito.\n" +
				"3. Solo después de que el usuario diga 'Adelante' o similar, llama a esta herramienta con `confirmed_by_user=true`.",
		)
	}

	// Validar credenciales: siempre requiere page_id; para token se acepta page_token directo
	// o el flujo de renovación con app_id/app_secret/user_token.
	if pageID == "" {
		return UserResult(
			"Facebook Page ID no configurado. " +
				"Configura en config.json (tools.social_media.facebook) o usa variables de entorno:\n" +
				"  PICOCLAW_TOOLS_SOCIAL_FACEBOOK_DEFAULT_PAGE_ID\n" +
				"Para token puedes usar:\n" +
				"  PICOCLAW_TOOLS_SOCIAL_FACEBOOK_DEFAULT_PAGE_TOKEN\n" +
				"o flujo de refresh:\n" +
				"  PICOCLAW_TOOLS_SOCIAL_FACEBOOK_APP_ID\n" +
				"  PICOCLAW_TOOLS_SOCIAL_FACEBOOK_APP_SECRET\n" +
				"  PICOCLAW_TOOLS_SOCIAL_FACEBOOK_USER_TOKEN\n" +
				"También puedes pasar page_id y page_token como parámetros.",
		)
	}
	if pageToken == "" && (t.appID == "" || t.appSecret == "" || t.userToken == "") {
		return UserResult(
			"Facebook token not configured. Provide either `default_page_token` or all refresh fields: `app_id`, `app_secret`, and `user_token` in tools.social_media.facebook.",
		)
	}

	callCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	comment = strings.TrimSpace(comment)

	var (
		postID string
		err    error
	)
	if imagePath == "" {
		postID, err = utils.FacebookPostTextOnly(
			callCtx,
			pageID,
			pageToken,
			t.appID,
			t.appSecret,
			t.userToken,
			message,
			comment,
		)
	} else {
		req := utils.FBPostRequest{
			PageID:    pageID,
			PageToken: pageToken,
			AppID:     t.appID,
			AppSecret: t.appSecret,
			UserToken: t.userToken,
			Message:   message,
			ImagePath: imagePath,
			Comment:   comment,
		}
		postID, err = utils.FacebookPost(callCtx, req)
	}
	if err != nil {
		errMsg := err.Error()
		if strings.Contains(strings.ToLower(errMsg), "publish_actions") {
			return ErrorResult(
				"Facebook rejected the token/permissions: `publish_actions` is deprecated. " +
					"Use a Page Access Token with `pages_manage_posts` (and optionally `pages_read_engagement` / `pages_manage_engagement`), " +
					"then retry with `tools.social_media.facebook.default_page_id` and `default_page_token`.",
			).WithError(err)
		}
		return ErrorResult(fmt.Sprintf("facebook post falló: %v", err)).WithError(err)
	}

	return UserResult(fmt.Sprintf("Publicación en Facebook exitosa. Post ID: %s", postID))
}

// ============== X (Twitter) Tools ==============

type XPostTweetTool struct {
	apiKey            string
	apiSecret         string
	accessToken       string
	accessTokenSecret string
}

func NewXPostTweetTool() *XPostTweetTool {
	return NewXPostTweetToolFromConfig("", "", "", "")
}

func NewXPostTweetToolFromConfig(
	configAPIKey, configAPISecret, configAccessToken, configAccessTokenSecret string,
) *XPostTweetTool {
	apiKey := strings.TrimSpace(os.Getenv(utils.EnvSocialXAPIKey))
	apiSecret := strings.TrimSpace(os.Getenv(utils.EnvSocialXAPISecret))
	accessToken := strings.TrimSpace(os.Getenv(utils.EnvSocialXAccessToken))
	accessTokenSecret := strings.TrimSpace(os.Getenv(utils.EnvSocialXAccessTokenSecret))

	if apiKey == "" {
		apiKey = strings.TrimSpace(configAPIKey)
	}
	if apiSecret == "" {
		apiSecret = strings.TrimSpace(configAPISecret)
	}
	if accessToken == "" {
		accessToken = strings.TrimSpace(configAccessToken)
	}
	if accessTokenSecret == "" {
		accessTokenSecret = strings.TrimSpace(configAccessTokenSecret)
	}

	return &XPostTweetTool{
		apiKey:            apiKey,
		apiSecret:         apiSecret,
		accessToken:       accessToken,
		accessTokenSecret: accessTokenSecret,
	}
}

func (t *XPostTweetTool) Name() string {
	return "x_post_tweet"
}

func (t *XPostTweetTool) Description() string {
	return "Publicar tweet en X/Twitter con imagen opcional. Soporta replies para crear hilos."
}

func (t *XPostTweetTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"message": map[string]any{
				"type":        "string",
				"description": "Texto del tweet",
			},
			"image_path": map[string]any{
				"type":        "string",
				"description": "Ruta absoluta de la imagen (opcional)",
			},
			"reply_to_tweet_id": map[string]any{
				"type":        "string",
				"description": "ID del tweet al que responder (opcional, para hilos)",
			},
			"confirmed_by_user": map[string]any{
				"type":        "boolean",
				"description": "REQUERIDO: Debes setear esto a true SOLO después de haber mostrado el borrador al usuario con la herramienta 'message' y haber recibido su aprobación.",
			},
		},
		"required": []string{"message", "confirmed_by_user"},
	}
}

func (t *XPostTweetTool) Execute(ctx context.Context, args map[string]any) *ToolResult {
	message, _ := args["message"].(string)
	imagePath, _ := args["image_path"].(string)
	replyToTweetID, _ := args["reply_to_tweet_id"].(string)
	confirmed, _ := args["confirmed_by_user"].(bool)

	if !confirmed {
		return UserResult(
			"❌ **ERROR DE SEGURIDAD**: No puedes publicar en X/Twitter sin aprobación previa. " +
				"Muestra primero el borrador usando 'message' y luego llama a esta herramienta con `confirmed_by_user=true`.",
		)
	}

	// Validar parámetros requeridos
	message = strings.TrimSpace(message)
	if message == "" {
		return ErrorResult("message es requerido")
	}

	// Validar credenciales
	if t.apiKey == "" || t.apiSecret == "" || t.accessToken == "" || t.accessTokenSecret == "" {
		return UserResult(
			"Credenciales de X/Twitter no configuradas. " +
				"Configura en config.json (tools.social_media.x) o usa variables de entorno:\n" +
				"  PICOCLAW_TOOLS_SOCIAL_X_API_KEY\n" +
				"  PICOCLAW_TOOLS_SOCIAL_X_API_SECRET\n" +
				"  PICOCLAW_TOOLS_SOCIAL_X_ACCESS_TOKEN\n" +
				"  PICOCLAW_TOOLS_SOCIAL_X_ACCESS_TOKEN_SECRET",
		)
	}

	callCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	req := utils.XPostRequest{
		APIKey:            t.apiKey,
		APISecret:         t.apiSecret,
		AccessToken:       t.accessToken,
		AccessTokenSecret: t.accessTokenSecret,
		Message:           message,
		ImagePath:         strings.TrimSpace(imagePath),
		ReplyToTweetID:    strings.TrimSpace(replyToTweetID),
	}

	tweetID, rateLimit, err := utils.XPostTweet(callCtx, req)
	if err != nil {
		result := fmt.Sprintf("X post falló: %v", err)
		if rateLimit != nil {
			result += fmt.Sprintf(" (%s)", utils.GetXRateLimitString(rateLimit))
		}
		return ErrorResult(result).WithError(err)
	}

	response := fmt.Sprintf("Tweet publicado exitosamente. Tweet ID: %s", tweetID)
	if rateLimit != nil {
		response += fmt.Sprintf(" (%s)", utils.GetXRateLimitString(rateLimit))
	}

	return UserResult(response)
}
