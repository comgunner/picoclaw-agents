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
	"time"
)

const (
	EnvSocialFacebookPageID    = "PICOCLAW_TOOLS_SOCIAL_FACEBOOK_DEFAULT_PAGE_ID"
	EnvSocialFacebookPageToken = "PICOCLAW_TOOLS_SOCIAL_FACEBOOK_DEFAULT_PAGE_TOKEN"
	EnvSocialFacebookAppID     = "PICOCLAW_TOOLS_SOCIAL_FACEBOOK_APP_ID"
	EnvSocialFacebookAppSecret = "PICOCLAW_TOOLS_SOCIAL_FACEBOOK_APP_SECRET"
	EnvSocialFacebookUserToken = "PICOCLAW_TOOLS_SOCIAL_FACEBOOK_USER_TOKEN"
)

// FBPostRequest define los parámetros para publicar en Facebook
type FBPostRequest struct {
	PageID    string
	PageToken string
	AppID     string
	AppSecret string
	UserToken string
	Message   string
	ImagePath string
	Comment   string // Opcional: comentario a agregar después de publicar
}

// FBResponse representa la respuesta de la API de Facebook
type FBResponse struct {
	ID     string   `json:"id"`
	PostID string   `json:"post_id"`
	Error  *FBError `json:"error,omitempty"`
}

// FBError representa un error de la API de Facebook
type FBError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Type    string `json:"type,omitempty"`
}

// FacebookPost publica una imagen con mensaje en Facebook.
// Si se proporciona un comentario, intenta agregarlo como comentario separado.
// Si Facebook bloquea el comentario (Error 368), fusiona el comentario en el post original.
func FacebookPost(ctx context.Context, req FBPostRequest) (string, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	accessToken := req.PageToken
	if accessToken == "" && req.AppID != "" && req.AppSecret != "" && req.UserToken != "" {
		refreshed, err := refreshFacebookPageToken(ctx, client, req.PageID, req.AppID, req.AppSecret, req.UserToken)
		if err != nil {
			return "", fmt.Errorf("error refreshing facebook page token: %v", err)
		}
		accessToken = refreshed
	}

	// 1. Preparar la subida de la imagen (multipart/form-data)
	file, err := os.Open(req.ImagePath)
	if err != nil {
		return "", fmt.Errorf("error abriendo imagen: %v", err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("source", filepath.Base(req.ImagePath))
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(part, file); err != nil {
		return "", fmt.Errorf("error copiando archivo: %v", err)
	}

	// Agregar campos adicionales
	if err := writer.WriteField("access_token", accessToken); err != nil {
		return "", err
	}
	if err := writer.WriteField("message", req.Message); err != nil {
		return "", err
	}
	if err := writer.Close(); err != nil {
		return "", err
	}

	// 2. Publicar imagen en Facebook Graph API v20.0
	url := fmt.Sprintf("https://graph.facebook.com/v20.0/%s/photos", req.PageID)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, body)
	if err != nil {
		return "", fmt.Errorf("error creando request: %v", err)
	}
	httpReq.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("error de red publicando foto: %v", err)
	}
	defer resp.Body.Close()

	var fbResp FBResponse
	if err := json.NewDecoder(resp.Body).Decode(&fbResp); err != nil {
		return "", fmt.Errorf("error decodificando respuesta: %v", err)
	}

	if fbResp.Error != nil {
		if fbResp.Error.Code == 190 && req.AppID != "" && req.AppSecret != "" && req.UserToken != "" {
			refreshed, refreshErr := refreshFacebookPageToken(
				ctx,
				client,
				req.PageID,
				req.AppID,
				req.AppSecret,
				req.UserToken,
			)
			if refreshErr == nil {
				req.PageToken = refreshed
				req.AppID = ""
				req.AppSecret = ""
				req.UserToken = ""
				return FacebookPost(ctx, req)
			}
		}
		return "", fmt.Errorf("error FB (%d): %s", fbResp.Error.Code, fbResp.Error.Message)
	}

	postID := fbResp.PostID
	if postID == "" {
		postID = fbResp.ID // Fallback si el endpoint devuelve 'id' en lugar de 'post_id'
	}

	// 3. Agregar comentario (si existe)
	if req.Comment != "" && postID != "" {
		commentErr := facebookAddComment(ctx, client, postID, req.Comment, accessToken)
		if commentErr != nil {
			// Manejo del Error 368 de Facebook (Control de seguridad)
			// Intentamos fusionar el comentario en el mensaje original
			fmt.Printf("[FB Warning] Error en comentario: %v. Intentando fusionar...\n", commentErr)
			mergedMessage := fmt.Sprintf("%s\n\n— Signal —\n%s", req.Message, req.Comment)
			if updateErr := facebookUpdatePost(ctx, client, postID, mergedMessage, accessToken); updateErr != nil {
				fmt.Printf("[FB Warning] Error fusionando comentario: %v\n", updateErr)
			}
		}
	}

	return postID, nil
}

// facebookAddComment intenta publicar un comentario en un post existente
func facebookAddComment(ctx context.Context, client *http.Client, postID, message, token string) error {
	apiURL := fmt.Sprintf("https://graph.facebook.com/v20.0/%s/comments", postID)
	return facebookPostAction(ctx, client, apiURL, message, token, "agregando comentario")
}

// facebookUpdatePost actualiza el mensaje de un post existente (fallback para comentarios bloqueados)
func facebookUpdatePost(ctx context.Context, client *http.Client, postID, message, token string) error {
	apiURL := fmt.Sprintf("https://graph.facebook.com/v20.0/%s", postID)
	return facebookPostAction(ctx, client, apiURL, message, token, "actualizando post")
}

func facebookPostAction(ctx context.Context, client *http.Client, apiURL, message, token, actionName string) error {
	data := url.Values{}
	data.Set("message", message)
	data.Set("access_token", token)

	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewReader([]byte(data.Encode())))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error de red %s: %v", actionName, err)
	}
	defer resp.Body.Close()

	var result FBResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("error decodificando respuesta de %s: %v", actionName, err)
	}

	if result.Error != nil {
		return fmt.Errorf("error FB %s (%d): %s", actionName, result.Error.Code, result.Error.Message)
	}

	return nil
}

// FacebookPostTextOnly publica solo texto (sin imagen) en Facebook
func FacebookPostTextOnly(
	ctx context.Context,
	pageID, pageToken, appID, appSecret, userToken, message, comment string,
) (string, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	accessToken := pageToken
	if accessToken == "" && appID != "" && appSecret != "" && userToken != "" {
		refreshed, err := refreshFacebookPageToken(ctx, client, pageID, appID, appSecret, userToken)
		if err != nil {
			return "", fmt.Errorf("error refreshing facebook page token: %v", err)
		}
		accessToken = refreshed
	}

	apiURL := fmt.Sprintf("https://graph.facebook.com/v20.0/%s/feed", pageID)

	data := url.Values{}
	data.Set("message", message)
	data.Set("access_token", accessToken)

	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewReader([]byte(data.Encode())))
	if err != nil {
		return "", fmt.Errorf("error creando request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error de red publicando: %v", err)
	}
	defer resp.Body.Close()

	var fbResp FBResponse
	if err := json.NewDecoder(resp.Body).Decode(&fbResp); err != nil {
		return "", fmt.Errorf("error decodificando respuesta: %v", err)
	}

	if fbResp.Error != nil {
		if fbResp.Error.Code == 190 && appID != "" && appSecret != "" && userToken != "" {
			refreshed, refreshErr := refreshFacebookPageToken(ctx, client, pageID, appID, appSecret, userToken)
			if refreshErr == nil {
				return FacebookPostTextOnly(ctx, pageID, refreshed, "", "", "", message, comment)
			}
		}
		return "", fmt.Errorf("error FB (%d): %s", fbResp.Error.Code, fbResp.Error.Message)
	}

	postID := fbResp.PostID
	if postID == "" {
		postID = fbResp.ID
	}

	// Agregar comentario si existe
	if comment != "" {
		commentErr := facebookAddComment(ctx, client, postID, comment, accessToken)
		if commentErr != nil {
			fmt.Printf("[FB Warning] Error en comentario: %v\n", commentErr)
		}
	}

	return postID, nil
}

func refreshFacebookPageToken(
	ctx context.Context,
	client *http.Client,
	pageID, appID, appSecret, userToken string,
) (string, error) {
	longLivedUserToken, err := exchangeForLongLivedUserToken(ctx, client, appID, appSecret, userToken)
	if err != nil {
		return "", err
	}
	return getPageTokenFromUserToken(ctx, client, pageID, longLivedUserToken)
}

func exchangeForLongLivedUserToken(
	ctx context.Context,
	client *http.Client,
	appID, appSecret, userToken string,
) (string, error) {
	apiURL := "https://graph.facebook.com/v20.0/oauth/access_token"
	q := url.Values{}
	q.Set("grant_type", "fb_exchange_token")
	q.Set("client_id", appID)
	q.Set("client_secret", appSecret)
	q.Set("fb_exchange_token", userToken)

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL+"?"+q.Encode(), nil)
	if err != nil {
		return "", fmt.Errorf("error creating exchange request: %v", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error exchanging user token: %v", err)
	}
	defer resp.Body.Close()

	var payload struct {
		AccessToken string   `json:"access_token"`
		Error       *FBError `json:"error,omitempty"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return "", fmt.Errorf("error decoding exchange response: %v", err)
	}
	if payload.Error != nil {
		return "", fmt.Errorf("error FB exchange (%d): %s", payload.Error.Code, payload.Error.Message)
	}
	if payload.AccessToken == "" {
		return "", fmt.Errorf("facebook exchange response did not include access_token")
	}
	return payload.AccessToken, nil
}

func getPageTokenFromUserToken(ctx context.Context, client *http.Client, pageID, userToken string) (string, error) {
	apiURL := fmt.Sprintf("https://graph.facebook.com/v20.0/%s", pageID)
	q := url.Values{}
	q.Set("fields", "access_token")
	q.Set("access_token", userToken)

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL+"?"+q.Encode(), nil)
	if err != nil {
		return "", fmt.Errorf("error creating page token request: %v", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error fetching page token: %v", err)
	}
	defer resp.Body.Close()

	var payload struct {
		AccessToken string   `json:"access_token"`
		Error       *FBError `json:"error,omitempty"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return "", fmt.Errorf("error decoding page token response: %v", err)
	}
	if payload.Error != nil {
		return "", fmt.Errorf("error FB page token (%d): %s", payload.Error.Code, payload.Error.Message)
	}
	if payload.AccessToken == "" {
		return "", fmt.Errorf("facebook page response did not include access_token")
	}
	return payload.AccessToken, nil
}
