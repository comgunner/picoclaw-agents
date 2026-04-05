// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/dghubble/oauth1"
)

const (
	EnvSocialXAPIKey            = "PICOCLAW_TOOLS_SOCIAL_X_API_KEY"
	EnvSocialXAPISecret         = "PICOCLAW_TOOLS_SOCIAL_X_API_SECRET"
	EnvSocialXAccessToken       = "PICOCLAW_TOOLS_SOCIAL_X_ACCESS_TOKEN"
	EnvSocialXAccessTokenSecret = "PICOCLAW_TOOLS_SOCIAL_X_ACCESS_TOKEN_SECRET"
)

// XPostRequest define los parámetros para publicar en X/Twitter
type XPostRequest struct {
	APIKey            string
	APISecret         string
	AccessToken       string
	AccessTokenSecret string
	Message           string
	ImagePath         string
	ReplyToTweetID    string // Opcional: ID del tweet al que responder
}

// XResponse representa la respuesta de la API de X
type XResponse struct {
	Data *XTweetData `json:"data,omitempty"`
}

// XTweetData contiene los datos del tweet creado
type XTweetData struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

// XErrorResponse representa un error de la API de X
type XErrorResponse struct {
	Errors []XErrorDetail `json:"errors,omitempty"`
	Title  string         `json:"title,omitempty"`
	Detail string         `json:"detail,omitempty"`
}

type XErrorDetail struct {
	Message string `json:"message"`
}

// XRateLimitInfo contiene información sobre rate limits
type XRateLimitInfo struct {
	Limit     int
	Remaining int
	Reset     time.Time
}

// XPostTweet publica un tweet en X/Twitter con imagen opcional
func XPostTweet(ctx context.Context, req XPostRequest) (string, *XRateLimitInfo, error) {
	// Configurar OAuth1
	config := oauth1.NewConfig(req.APIKey, req.APISecret)
	token := oauth1.NewToken(req.AccessToken, req.AccessTokenSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	var mediaIDs []string

	// 1. Subir imagen si se proporciona
	if req.ImagePath != "" {
		mediaID, rateLimit, err := xUploadMedia(ctx, httpClient, req.ImagePath)
		if err != nil {
			return "", rateLimit, fmt.Errorf("error subiendo imagen: %v", err)
		}
		mediaIDs = append(mediaIDs, mediaID)
	}

	// 2. Preparar el payload del tweet
	tweetPayload := map[string]any{
		"text": req.Message,
	}

	if len(mediaIDs) > 0 {
		tweetPayload["media"] = map[string]any{
			"media_ids": mediaIDs,
		}
	}

	if req.ReplyToTweetID != "" {
		tweetPayload["reply"] = map[string]any{
			"in_reply_to_tweet_id": req.ReplyToTweetID,
		}
	}

	// 3. Publicar tweet (API v2)
	tweetID, rateLimit, err := xCreateTweet(ctx, httpClient, tweetPayload)
	if err != nil {
		return "", rateLimit, err
	}

	return tweetID, rateLimit, nil
}

// xUploadMedia sube una imagen a la API de media upload de X (v1.1)
func xUploadMedia(ctx context.Context, client *http.Client, imagePath string) (string, *XRateLimitInfo, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return "", nil, fmt.Errorf("error abriendo imagen: %v", err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("media", filepath.Base(imagePath))
	if err != nil {
		return "", nil, err
	}
	if _, err := io.Copy(part, file); err != nil {
		return "", nil, fmt.Errorf("error copiando archivo: %v", err)
	}
	if err := writer.Close(); err != nil {
		return "", nil, err
	}

	// API v1.1 media upload endpoint
	uploadURL := "https://upload.twitter.com/1.1/media/upload.json"
	httpReq, err := http.NewRequestWithContext(ctx, "POST", uploadURL, body)
	if err != nil {
		return "", nil, fmt.Errorf("error creando request: %v", err)
	}
	httpReq.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := client.Do(httpReq)
	if err != nil {
		return "", nil, fmt.Errorf("error de red subiendo imagen: %v", err)
	}
	defer resp.Body.Close()

	rateLimit := parseXRateLimit(resp.Header)

	// Verificar rate limit
	if resp.StatusCode == http.StatusTooManyRequests {
		return "", rateLimit, fmt.Errorf("rate limit excedido. Reset: %v", rateLimit.Reset)
	}

	var result map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", rateLimit, fmt.Errorf("error decodificando respuesta: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		errorBytes, _ := json.Marshal(result)
		return "", rateLimit, fmt.Errorf("error subiendo imagen (%d): %s", resp.StatusCode, string(errorBytes))
	}

	mediaID, ok := result["media_id_string"].(string)
	if !ok {
		return "", rateLimit, fmt.Errorf("no se pudo obtener media_id de la respuesta")
	}

	return mediaID, rateLimit, nil
}

// xCreateTweet crea un tweet usando la API v2
func xCreateTweet(ctx context.Context, client *http.Client, payload map[string]any) (string, *XRateLimitInfo, error) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", nil, fmt.Errorf("error serializando payload: %v", err)
	}

	// API v2 tweets endpoint
	tweetURL := "https://api.twitter.com/2/tweets"
	httpReq, err := http.NewRequestWithContext(ctx, "POST", tweetURL, bytes.NewReader(jsonData))
	if err != nil {
		return "", nil, fmt.Errorf("error creando request: %v", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(httpReq)
	if err != nil {
		return "", nil, fmt.Errorf("error de red creando tweet: %v", err)
	}
	defer resp.Body.Close()

	rateLimit := parseXRateLimit(resp.Header)

	// Verificar rate limit
	if resp.StatusCode == http.StatusTooManyRequests {
		return "", rateLimit, fmt.Errorf("rate limit excedido. Reset: %v", rateLimit.Reset)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", rateLimit, fmt.Errorf("error leyendo respuesta: %v", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		var errResp XErrorResponse
		if err := json.Unmarshal(body, &errResp); err != nil {
			return "", rateLimit, fmt.Errorf("error creando tweet (%d): %s", resp.StatusCode, string(body))
		}
		if len(errResp.Errors) > 0 {
			return "", rateLimit, fmt.Errorf("error X: %s", errResp.Errors[0].Message)
		}
		return "", rateLimit, fmt.Errorf("error creando tweet (%d): %s", resp.StatusCode, errResp.Detail)
	}

	var result XResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", rateLimit, fmt.Errorf("error decodificando respuesta: %v", err)
	}

	if result.Data == nil {
		return "", rateLimit, fmt.Errorf("respuesta vacía de la API de X")
	}

	return result.Data.ID, rateLimit, nil
}

// parseXRateLimit extrae información de rate limit de los headers de X
func parseXRateLimit(header http.Header) *XRateLimitInfo {
	info := &XRateLimitInfo{
		Limit:     0,
		Remaining: 0,
		Reset:     time.Time{},
	}

	if limit := header.Get("X-Rate-Limit-Limit"); limit != "" {
		info.Limit, _ = strconv.Atoi(limit)
	}

	if remaining := header.Get("X-Rate-Limit-Remaining"); remaining != "" {
		info.Remaining, _ = strconv.Atoi(remaining)
	}

	if resetRaw := header.Get("X-Rate-Limit-Reset"); resetRaw != "" {
		if resetUnix, err := strconv.ParseInt(resetRaw, 10, 64); err == nil {
			info.Reset = time.Unix(resetUnix, 0)
		}
	}

	return info
}

// GetXRateLimitString devuelve una representación legible del rate limit
func GetXRateLimitString(rateLimit *XRateLimitInfo) string {
	if rateLimit == nil {
		return ""
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Rate Limit: %d/%d", rateLimit.Remaining, rateLimit.Limit))
	if !rateLimit.Reset.IsZero() {
		sb.WriteString(fmt.Sprintf(", Reset: %s", rateLimit.Reset.Format(time.RFC3339)))
	}
	return sb.String()
}

// BuildOAuth1Config crea una configuración OAuth1 para X
func BuildOAuth1Config(apiKey, apiSecret, accessToken, accessTokenSecret string) *oauth1.Config {
	return oauth1.NewConfig(apiKey, apiSecret)
}

// BuildOAuth1Token crea un token OAuth1 para X
func BuildOAuth1Token(accessToken, accessTokenSecret string) *oauth1.Token {
	return oauth1.NewToken(accessToken, accessTokenSecret)
}

// CreateOAuth1Client crea un cliente HTTP con autenticación OAuth1
func CreateOAuth1Client(config *oauth1.Config, token *oauth1.Token) *http.Client {
	return config.Client(oauth1.NoContext, token)
}

// URLEncodeParams codifica parámetros como application/x-www-form-urlencoded
func URLEncodeParams(params map[string]string) string {
	values := url.Values{}
	for k, v := range params {
		values.Set(k, v)
	}
	return values.Encode()
}
