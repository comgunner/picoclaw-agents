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
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/comgunner/picoclaw/pkg/auth"
	"github.com/comgunner/picoclaw/pkg/logger"
	"github.com/comgunner/picoclaw/pkg/utils"
)

// ============== Antigravity Image Generation Tool ==============

// Antigravity image endpoint fallback chain (same pattern as Python client).
var antigravityImageBaseURLs = []string{
	"https://daily-cloudcode-pa.sandbox.googleapis.com", // Priority 1: Sandbox (less rate-limited)
	"https://daily-cloudcode-pa.googleapis.com",         // Priority 2: Daily
	"https://cloudcode-pa.googleapis.com",               // Priority 3: Prod
}

const (
	antigravityImageDefaultModel = "gemini-3.1-flash-image"
	antigravityImageUserAgent    = "antigravity"
	antigravityImageXGoogClient  = "google-cloud-sdk vscode_cloudshelleditor/0.1"

	// Cooldown default values.
	defaultCooldownSeconds = 150 // 2.5 minutes (anti-ban protection)
	minCooldownSeconds     = 60  // Minimum allowed cooldown
)

// Retry delays for 429 rate limiting (matches Python implementation exactly).
var antigravityImageRetryDelays = []time.Duration{
	30 * time.Second,  // 0.5 min
	60 * time.Second,  // 1 min
	120 * time.Second, // 2 min
	300 * time.Second, // 5 min
	600 * time.Second, // 10 min
}

// Safety settings — all OFF (from Antigravity-Manager proxy/handlers).
var antigravitySafetySettings = []map[string]any{
	{"category": "HARM_CATEGORY_HARASSMENT", "threshold": "OFF"},
	{"category": "HARM_CATEGORY_HATE_SPEECH", "threshold": "OFF"},
	{"category": "HARM_CATEGORY_SEXUALLY_EXPLICIT", "threshold": "OFF"},
	{"category": "HARM_CATEGORY_DANGEROUS_CONTENT", "threshold": "OFF"},
	{"category": "HARM_CATEGORY_CIVIC_INTEGRITY", "threshold": "OFF"},
}

// ImageGenAntigravityTool generates images using Google Antigravity via OAuth.
type ImageGenAntigravityTool struct {
	model        string
	aspectRatio  string
	outputDir    string
	workspace    string
	tracker      *utils.ImageGenTracker
	cooldown     *utils.ImageCooldown
	cooldownSecs int
}

// NewImageGenAntigravityTool creates a new instance with defaults.
func NewImageGenAntigravityTool() *ImageGenAntigravityTool {
	return NewImageGenAntigravityToolFromConfig("", "", "", "", 0, nil)
}

// NewImageGenAntigravityToolFromConfig creates an instance from configuration.
func NewImageGenAntigravityToolFromConfig(
	configModel, configAspectRatio, configOutputDir, workspace string,
	cooldownSecs int,
	cooldown *utils.ImageCooldown,
) *ImageGenAntigravityTool {
	model := strings.TrimSpace(configModel)
	if model == "" {
		model = antigravityImageDefaultModel
	}
	// Strip provider prefixes for consistency.
	model = strings.TrimPrefix(model, "antigravity/")
	model = strings.TrimPrefix(model, "google-antigravity/")
	model = strings.TrimPrefix(model, "models/")

	aspectRatio := strings.TrimSpace(configAspectRatio)
	if aspectRatio == "" {
		aspectRatio = "1:1"
	}

	outputDir := strings.TrimSpace(configOutputDir)
	if outputDir == "" {
		outputDir = resolveOutputDir("", workspace)
	} else {
		outputDir = resolvePathInWorkspace(outputDir, workspace)
	}

	// Initialize cooldown — FIX #4: pass workspace, not empty string.
	cd := cooldown
	if cd == nil {
		cd, _ = utils.NewImageCooldown(workspace)
	}

	trackerPath := filepath.Join(outputDir, "tracker.json")
	tracker, _ := utils.NewImageGenTracker(trackerPath)

	// FIX #3: use effectiveCooldown instead of always defaultCooldownSeconds.
	effectiveCooldown := cooldownSecs
	if effectiveCooldown <= 0 {
		effectiveCooldown = defaultCooldownSeconds
	}

	return &ImageGenAntigravityTool{
		model:        model,
		aspectRatio:  aspectRatio,
		outputDir:    outputDir,
		workspace:    workspace,
		tracker:      tracker,
		cooldown:     cd,
		cooldownSecs: effectiveCooldown,
	}
}

func (t *ImageGenAntigravityTool) Name() string {
	return "image_gen_antigravity"
}

func (t *ImageGenAntigravityTool) GetTracker() *utils.ImageGenTracker {
	return t.tracker
}

func (t *ImageGenAntigravityTool) Description() string {
	return fmt.Sprintf(
		"Generate image using Google Antigravity (OAuth). Model: %s. No API key needed — uses stored OAuth credentials. Mandatory cooldown of %ds after each generation (anti-ban). Supports aspect-ratio: 1:1, 16:9, 9:16, 4:5, 3:4.",
		t.model,
		t.cooldownSecs,
	)
}

func (t *ImageGenAntigravityTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"prompt": map[string]any{
				"type":        "string",
				"description": "Text prompt for image generation",
			},
			"model": map[string]any{
				"type":        "string",
				"description": fmt.Sprintf("Image model (default: %s)", t.model),
			},
			"aspect_ratio": map[string]any{
				"type":        "string",
				"description": "Aspect ratio: '1:1', '16:9', '9:16', '4:5', '3:4' (default: 1:1)",
			},
			"script_path": map[string]any{
				"type":        "string",
				"description": "Path to previously generated script file. ALWAYS use if coming from 'text_script_create'.",
			},
		},
		"required": []string{"prompt"},
	}
}

func (t *ImageGenAntigravityTool) Execute(ctx context.Context, args map[string]any) *ToolResult {
	prompt, _ := args["prompt"].(string)
	model, _ := args["model"].(string)
	aspectRatio, _ := args["aspect_ratio"].(string)
	scriptPath, _ := args["script_path"].(string)

	prompt = strings.TrimSpace(prompt)
	if prompt == "" {
		return ErrorResult("prompt is required")
	}

	// Check cooldown BEFORE doing anything.
	if t.cooldown != nil && t.cooldown.IsOnCooldown() {
		remaining := t.cooldown.GetRemaining()
		remainingStr := formatDuration(time.Duration(remaining * float64(time.Second)))
		info := t.cooldown.GetInfo()
		provider := "antigravity"
		if p, ok := info["provider"].(string); ok {
			provider = p
		}
		return UserResult(fmt.Sprintf(
			"⏳ **Cooldown active** — Cannot generate image now.\n\n"+
				"⏱ **Time remaining:** %s\n"+
				"🔐 **Provider:** %s\n\n"+
				"💡 Cooldown is mandatory (%ds) to protect against rate limits.\n"+
				"💡 Meanwhile, you can:\n"+
				"  - Edit the prompt and retry after cooldown\n"+
				"  - Use `image_gen_create` with provider 'gemini' or 'ideogram' if you have API keys",
			remainingStr, provider, t.cooldownSecs,
		))
	}

	if model == "" {
		model = t.model
	} else {
		model = strings.TrimPrefix(model, "antigravity/")
		model = strings.TrimPrefix(model, "google-antigravity/")
		model = strings.TrimPrefix(model, "models/")
	}

	if aspectRatio == "" {
		aspectRatio = t.aspectRatio
	}

	// 1. Generate unique ID.
	id := utils.GenerateID()

	// 2. Register in tracker.
	if t.tracker != nil {
		t.tracker.Add(utils.ImageGenRecord{
			ID:          id,
			DateTime:    time.Now().Format("2006-01-02 15:04:05"),
			Prompt:      prompt,
			Provider:    "antigravity",
			AspectRatio: aspectRatio,
			Model:       model,
			Language:    utils.DetectLanguage(prompt),
			Metadata: map[string]string{
				"status": "pending",
			},
		})
	}

	// 3. Verify OAuth credentials.
	cred, err := auth.GetCredential("google-antigravity")
	if err != nil || cred == nil {
		if t.tracker != nil {
			t.tracker.UpdateMetadata(id, "status", "failed")
			t.tracker.UpdateMetadata(id, "error", "No Antigravity OAuth credentials")
		}
		return UserResult(
			"❌ No Antigravity OAuth credentials configured.\n\nRun: `picoclaw auth login --provider google-antigravity`",
		)
	}

	// Auto-refresh if token is expired or about to expire.
	if (cred.NeedsRefresh() || cred.IsExpired()) && cred.RefreshToken != "" {
		oauthCfg := auth.GoogleAntigravityOAuthConfig()
		refreshed, refreshErr := auth.RefreshAccessToken(cred, oauthCfg)
		if refreshErr == nil {
			refreshed.Email = cred.Email
			if refreshed.ProjectID == "" {
				refreshed.ProjectID = cred.ProjectID
			}
			_ = auth.SetCredential("google-antigravity", refreshed)
			cred = refreshed
			logger.InfoCF("tools.antigravity", "Token refreshed for image generation", nil)
		} else {
			logger.WarnCF("tools.antigravity", "Token refresh failed", map[string]any{
				"error": refreshErr.Error(),
			})
		}
	}

	if cred.IsExpired() {
		if t.tracker != nil {
			t.tracker.UpdateMetadata(id, "status", "failed")
			t.tracker.UpdateMetadata(id, "error", "OAuth token expired, refresh failed")
		}
		return UserResult("❌ Antigravity OAuth token expired. Run: `picoclaw auth login --provider google-antigravity`")
	}

	// 4. Get project ID.
	projectID := cred.ProjectID
	if projectID == "" {
		fetchedID, fetchErr := fetchAntigravityProjectIDImage(cred.AccessToken)
		if fetchErr != nil {
			if t.tracker != nil {
				t.tracker.UpdateMetadata(id, "status", "failed")
				t.tracker.UpdateMetadata(id, "error", fmt.Sprintf("Could not fetch project ID: %v", fetchErr))
			}
			return ErrorResult(fmt.Sprintf("Could not get project ID: %v", fetchErr))
		}
		projectID = fetchedID
		cred.ProjectID = projectID
		_ = auth.SetCredential("google-antigravity", cred)
	}

	if t.tracker != nil {
		t.tracker.UpdateMetadata(id, "status", "in_progress")
	}

	// 5. Ensure output directory.
	if err := os.MkdirAll(t.outputDir, 0o755); err != nil {
		if t.tracker != nil {
			t.tracker.UpdateMetadata(id, "status", "failed")
			t.tracker.UpdateMetadata(id, "error", err.Error())
		}
		return ErrorResult(fmt.Sprintf("Error creating directory: %v", err))
	}

	imageDir := filepath.Join(t.outputDir, id)

	// 6. If script_path, read content.
	var scriptContent string
	if scriptPath != "" {
		data, readErr := os.ReadFile(scriptPath)
		if readErr == nil {
			scriptContent = string(data)
			scriptDir := filepath.Dir(scriptPath)
			if strings.TrimSpace(scriptDir) != "" && scriptDir != "." {
				imageDir = scriptDir
				if base := filepath.Base(scriptDir); strings.TrimSpace(base) != "" {
					id = base
				}
			}
		}
	}

	// 7. Generate image.
	callCtx, cancel := context.WithTimeout(ctx, 180*time.Second)
	defer cancel()

	imagePath := filepath.Join(imageDir, fmt.Sprintf("%s.-imagen.jpg", id))
	imageBytes, genErr := generateImageWithAntigravity(callCtx, cred.AccessToken, projectID, model, prompt, aspectRatio)
	if genErr != nil {
		if t.tracker != nil {
			t.tracker.UpdateMetadata(id, "status", "failed")
			t.tracker.UpdateMetadata(id, "error", genErr.Error())
		}

		if isRateLimitError(genErr) {
			return UserResult(fmt.Sprintf(
				"⏳ **Rate Limited (429)** — Antigravity API limit is ~1-2 images every 10 minutes.\n\n"+
					"❌ Error: %v\n\n"+
					"💡 Automatic retries with delays: 30s, 60s, 120s, 300s, 600s\n"+
					"💡 Wait 5-10 minutes and try again\n"+
					"💡 Or use `image_gen_create` with provider 'gemini' or 'ideogram'",
				genErr,
			))
		}

		return ErrorResult(fmt.Sprintf("image_gen_antigravity failed: %v", genErr)).WithError(genErr)
	}

	// 8. Save image.
	if err := os.MkdirAll(imageDir, 0o755); err != nil {
		return ErrorResult(fmt.Sprintf("Error creating image directory: %v", err))
	}
	if err := os.WriteFile(imagePath, imageBytes, 0o644); err != nil {
		if t.tracker != nil {
			t.tracker.UpdateMetadata(id, "status", "failed")
			t.tracker.UpdateMetadata(id, "error", err.Error())
		}
		return ErrorResult(fmt.Sprintf("Error saving image: %v", err))
	}

	// 9. Save visual prompt.
	promptPath := filepath.Join(imageDir, fmt.Sprintf("%s.-prompt_visual.txt", id))
	_ = os.WriteFile(promptPath, []byte(prompt), 0o644)

	// 10. Save script if exists.
	var savedScriptPath string
	if scriptContent != "" {
		if scriptPath != "" {
			savedScriptPath = scriptPath
		} else {
			savedScriptPath = filepath.Join(imageDir, fmt.Sprintf("%s.-script.txt", id))
			_ = os.WriteFile(savedScriptPath, []byte(scriptContent), 0o644)
		}
	}

	// 11. ACTIVATE GLOBAL COOLDOWN (anti-ban protection).
	if t.cooldown != nil {
		_ = t.cooldown.Set(float64(t.cooldownSecs), "antigravity", model)
	}

	// 12. Update tracker with success.
	if t.tracker != nil {
		t.tracker.UpdateMetadata(id, "status", "generated")
		t.tracker.UpdateMetadata(id, "image_path", imagePath)
		t.tracker.UpdateMetadata(id, "prompt_path", promptPath)
		if savedScriptPath != "" {
			t.tracker.UpdateMetadata(id, "script_path", savedScriptPath)
		}
	}

	trackerFile := "tracker.json"
	if t.tracker != nil {
		trackerFile = t.tracker.TrackerPath
	}

	cooldownInfo := ""
	if t.cooldown != nil {
		cooldownInfo = fmt.Sprintf(
			"\n⏳ **Cooldown activated:** %ds before next image (anti-ban protection)",
			t.cooldownSecs,
		)
	}

	return ImageResult(fmt.Sprintf(`✅ **Image generated successfully (Antigravity OAuth)**

📁 **Files created:**
┌─────────────────────────────────────────────────────────────────┐
│ Image:    %s
│ Prompt:   %s
│ Script:   %s
└─────────────────────────────────────────────────────────────────┘

🎨 **Model:** %s
📐 **Aspect Ratio:** %s
🔐 **Auth:** OAuth (google-antigravity)
📊 **Size:** %s

📊 **Status recorded in:**
%s (status: generated)%s

**What do you want to do with this image?**
[1] ✅ **Like** - Approved for publishing
[2] 🔄 **Dislike** - Regenerate after cooldown
[3] 📝 **Generate text** - Create social copy
[4] ⏰ **Schedule** - Publish later
[5] ❌ **Finish** - End without publishing`,
		imagePath, promptPath, savedScriptPath,
		model, aspectRatio, formatFileSize(len(imageBytes)),
		trackerFile, cooldownInfo),
		imagePath) // ← Media path para adjuntar al canal
}

// generateImageWithAntigravity generates an image via Antigravity API.
// Uses endpoint fallback chain (sandbox → daily → prod) with exponential retry.
func generateImageWithAntigravity(
	ctx context.Context,
	accessToken, projectID, model, prompt, aspectRatio string,
) ([]byte, error) {
	// Strip provider prefixes.
	model = strings.TrimPrefix(model, "antigravity/")
	model = strings.TrimPrefix(model, "google-antigravity/")
	model = strings.TrimPrefix(model, "models/")
	model = strings.ReplaceAll(model, "antigravity/", "")

	requestID := fmt.Sprintf("agent-img-%d-%s", time.Now().UnixMilli(), randomString(9))

	envelope := map[string]any{
		"project": projectID,
		"model":   model,
		"request": map[string]any{
			"contents": []map[string]any{
				{
					"role": "user",
					"parts": []map[string]any{
						{"text": prompt},
					},
				},
			},
			"generationConfig": map[string]any{
				"responseModalities": []string{"IMAGE"},
			},
			"safetySettings": antigravitySafetySettings,
		},
		"requestType": "CHAT",
		"userAgent":   antigravityImageUserAgent,
		"requestId":   requestID,
	}

	bodyBytes, _ := json.Marshal(envelope)

	var lastErr error
	for attempt, delay := range antigravityImageRetryDelays {
		var resp *http.Response
		var respBody []byte

		// Try each endpoint fallback.
		for _, baseURL := range antigravityImageBaseURLs {
			req, reqErr := http.NewRequestWithContext(ctx, "POST",
				fmt.Sprintf("%s/v1internal:generateContent", baseURL),
				bytes.NewReader(bodyBytes),
			)
			if reqErr != nil {
				lastErr = fmt.Errorf("creating request: %w", reqErr)
				continue
			}

			req.Header.Set("Authorization", "Bearer "+accessToken)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("User-Agent", antigravityImageUserAgent)
			req.Header.Set("X-Goog-Api-Client", antigravityImageXGoogClient)

			client := &http.Client{Timeout: 180 * time.Second}
			resp, reqErr = client.Do(req)
			if reqErr != nil {
				lastErr = reqErr
				continue
			}

			respBody, reqErr = io.ReadAll(resp.Body)
			if reqErr != nil {
				resp.Body.Close()
				lastErr = reqErr
				continue
			}

			// If success or non-retryable error, stop trying endpoints.
			if resp.StatusCode != 429 && resp.StatusCode != 500 && resp.StatusCode != 503 {
				break
			}
			resp.Body.Close()
			resp = nil
		}

		if resp == nil {
			lastErr = fmt.Errorf(
				"all Antigravity endpoints failed (attempt %d/%d)",
				attempt+1,
				len(antigravityImageRetryDelays),
			)
			if attempt < len(antigravityImageRetryDelays)-1 {
				logger.WarnCF("tools.antigravity", "All endpoints failed, retrying", map[string]any{
					"attempt": attempt + 1,
					"delay":   delay.String(),
				})
				select {
				case <-time.After(delay):
					continue
				case <-ctx.Done():
					return nil, ctx.Err()
				}
			}
			return nil, lastErr
		}

		if resp.StatusCode == 429 {
			resp.Body.Close()
			logger.WarnCF("tools.antigravity", "Rate limited, waiting", map[string]any{
				"delay":   delay.String(),
				"attempt": attempt + 1,
			})
			select {
			case <-time.After(delay):
				continue
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}

		if resp.StatusCode == 500 {
			resp.Body.Close()
			logger.WarnCF("tools.antigravity", "Server error, waiting", map[string]any{
				"delay":   delay.String(),
				"attempt": attempt + 1,
			})
			select {
			case <-time.After(delay):
				continue
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}

		if resp.StatusCode == 503 {
			resp.Body.Close()
			logger.WarnCF("tools.antigravity", "Service unavailable, waiting", map[string]any{
				"delay":   delay.String(),
				"attempt": attempt + 1,
			})
			select {
			case <-time.After(delay):
				continue
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}

		if resp.StatusCode >= 400 {
			resp.Body.Close()
			return nil, fmt.Errorf("Antigravity API error (HTTP %d): %s", resp.StatusCode, string(respBody))
		}

		// Success — parse response.
		return extractImageFromAntigravityResponse(respBody)
	}

	return nil, fmt.Errorf("image generation failed after %d retries: %w", len(antigravityImageRetryDelays), lastErr)
}

// extractImageFromAntigravityResponse extracts image bytes from the JSON response.
func extractImageFromAntigravityResponse(body []byte) ([]byte, error) {
	var data map[string]any
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}

	// Navigate: response.candidates[0].content.parts[].inlineData.data
	responseObj, ok := data["response"].(map[string]any)
	if !ok {
		responseObj = data
	}

	candidates, ok := responseObj["candidates"].([]any)
	if !ok || len(candidates) == 0 {
		return nil, fmt.Errorf("no candidates in response")
	}

	firstCandidate, ok := candidates[0].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid candidate format")
	}

	content, ok := firstCandidate["content"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid content format")
	}

	parts, ok := content["parts"].([]any)
	if !ok {
		return nil, fmt.Errorf("no parts in content")
	}

	for _, partRaw := range parts {
		part, ok := partRaw.(map[string]any)
		if !ok {
			continue
		}

		if inlineData, hasInline := part["inlineData"].(map[string]any); hasInline {
			if b64Data, hasData := inlineData["data"].(string); hasData {
				imageBytes, decErr := base64.StdEncoding.DecodeString(b64Data)
				if decErr != nil {
					return nil, fmt.Errorf("decoding base64 image: %w", decErr)
				}
				return imageBytes, nil
			}
		}
	}

	return nil, fmt.Errorf("no image data found in response parts")
}

// fetchAntigravityProjectIDImage gets the project ID from the API.
func fetchAntigravityProjectIDImage(accessToken string) (string, error) {
	reqBody, _ := json.Marshal(map[string]any{
		"metadata": map[string]any{
			"ideType":    "IDE_UNSPECIFIED",
			"platform":   "PLATFORM_UNSPECIFIED",
			"pluginType": "GEMINI",
		},
	})

	for _, baseURL := range antigravityImageBaseURLs {
		req, err := http.NewRequest("POST",
			fmt.Sprintf("%s/v1internal:loadCodeAssist", baseURL),
			bytes.NewReader(reqBody),
		)
		if err != nil {
			continue
		}
		req.Header.Set("Authorization", "Bearer "+accessToken)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", antigravityImageUserAgent)
		req.Header.Set("X-Goog-Api-Client", antigravityImageXGoogClient)

		client := &http.Client{Timeout: 15 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			continue
		}

		body, readErr := io.ReadAll(resp.Body)
		resp.Body.Close()
		if readErr != nil {
			continue
		}

		if resp.StatusCode != http.StatusOK {
			continue
		}

		var result struct {
			CloudAICompanionProject string `json:"cloudaicompanionProject"`
		}
		if err := json.Unmarshal(body, &result); err != nil {
			continue
		}

		if result.CloudAICompanionProject != "" {
			return result.CloudAICompanionProject, nil
		}
	}

	return "bamboo-precept-lgxtn", nil // Fallback project ID.
}

// isRateLimitError checks if the error is a 429 rate limit.
func isRateLimitError(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "429") ||
		strings.Contains(msg, "rate limit") ||
		strings.Contains(msg, "rate limited") ||
		strings.Contains(msg, "quota")
}

// formatDuration formats a duration in human-readable form.
func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%.0fs", d.Seconds())
	}
	mins := int(d.Minutes())
	secs := int(d.Seconds()) % 60
	if secs == 0 {
		return fmt.Sprintf("%dm", mins)
	}
	return fmt.Sprintf("%dm %ds", mins, secs)
}

// formatFileSize formats a file size in human-readable form.
func formatFileSize(bytes int) string {
	if bytes < 1024 {
		return fmt.Sprintf("%d B", bytes)
	}
	if bytes < 1024*1024 {
		return fmt.Sprintf("%.1f KB", float64(bytes)/1024)
	}
	return fmt.Sprintf("%.1f MB", float64(bytes)/(1024*1024))
}

// randomString generates a random string for request IDs.
func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
	}
	return string(b)
}
