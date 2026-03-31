// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package auth

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/comgunner/picoclaw/cmd/picoclaw/internal"
	"github.com/comgunner/picoclaw/pkg/auth"
	"github.com/comgunner/picoclaw/pkg/config"
	"github.com/comgunner/picoclaw/pkg/providers"
)

const supportedProvidersMsg = "supported providers: openai, anthropic, google-antigravity"

func authLoginCmd(provider string, useDeviceCode bool) error {
	switch provider {
	case "openai":
		return authLoginOpenAI(useDeviceCode)
	case "anthropic":
		// Check if user wants browser or token auth
		fmt.Println("Select authentication method:")
		fmt.Println("1) Browser (OAuth)")
		fmt.Println("2) API Key (Token)")
		fmt.Print("Enter choice (1 or 2, default 1): ")
		reader := bufio.NewReader(os.Stdin)
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		if choice == "2" {
			return authLoginPasteToken(provider)
		}
		return authLoginAnthropicBrowser()
	case "google-antigravity", "antigravity":
		return authLoginGoogleAntigravity()
	default:
		return fmt.Errorf("unsupported provider: %s (%s)", provider, supportedProvidersMsg)
	}
}

func authLoginOpenAI(useDeviceCode bool) error {
	cfg := auth.OpenAIOAuthConfig()

	var cred *auth.AuthCredential
	var err error

	if useDeviceCode {
		cred, err = auth.LoginDeviceCode(cfg)
	} else {
		cred, err = auth.LoginBrowser(cfg)
	}

	if err != nil {
		return fmt.Errorf("login failed: %w", err)
	}

	if err = auth.SetCredential("openai", cred); err != nil {
		return fmt.Errorf("failed to save credentials: %w", err)
	}

	appCfg, err := internal.LoadConfig()
	if err == nil {
		// Update Providers (legacy format)
		appCfg.Providers.OpenAI.AuthMethod = "oauth"

		// Use shared function to add models
		addedCount := AddOpenAIModels(appCfg)

		if err = config.SaveConfig(internal.GetConfigPath(), appCfg); err != nil {
			return fmt.Errorf("could not update config: %w", err)
		}

		if addedCount > 0 {
			fmt.Printf("\n✓ Added %d OpenAI models to config\n", addedCount)
		}
	}

	fmt.Println("Login successful!")
	if cred.AccountID != "" {
		fmt.Printf("Account: %s\n", cred.AccountID)
	}
	fmt.Println("Default model set to: gpt-5.2")

	return nil
}

func authLoginAnthropicBrowser() error {
	cfg := auth.AnthropicOAuthConfig()

	cred, err := auth.LoginBrowser(cfg)
	if err != nil {
		return fmt.Errorf("login failed: %w", err)
	}

	cred.Provider = "anthropic"

	if err = auth.SetCredential("anthropic", cred); err != nil {
		return fmt.Errorf("failed to save credentials: %w", err)
	}

	appCfg, err := internal.LoadConfig()
	if err == nil {
		// Update Providers (legacy format)
		appCfg.Providers.Anthropic.AuthMethod = "oauth"

		// Use shared function to add models
		addedCount := AddAnthropicModels(appCfg)

		if err = config.SaveConfig(internal.GetConfigPath(), appCfg); err != nil {
			fmt.Printf("Warning: could not update config: %v\n", err)
		} else {
			fmt.Printf("\n✓ Added %d Anthropic models to config\n", addedCount)
		}
	}

	fmt.Println("\n✓ Anthropic OAuth login successful!")
	fmt.Println("Default model set to: claude-sonnet-4-6")
	if cred.AccountID != "" {
		fmt.Printf("Account: %s\n", cred.AccountID)
	}

	return nil
}

func authLoginGoogleAntigravity() error {
	cfg := auth.GoogleAntigravityOAuthConfig()

	cred, err := auth.LoginBrowser(cfg)
	if err != nil {
		return fmt.Errorf("login failed: %w", err)
	}

	cred.Provider = "google-antigravity"

	// Fetch user email from Google userinfo
	email, err := fetchGoogleUserEmail(cred.AccessToken)
	if err != nil {
		fmt.Printf("Warning: could not fetch email: %v\n", err)
	} else {
		cred.Email = email
		fmt.Printf("Email: %s\n", email)
	}

	// Fetch Cloud Code Assist project ID
	projectID, err := providers.FetchAntigravityProjectID(cred.AccessToken)
	if err != nil {
		fmt.Printf("Warning: could not fetch project ID: %v\n", err)
		fmt.Println("You may need Google Cloud Code Assist enabled on your account.")
	} else {
		cred.ProjectID = projectID
		fmt.Printf("Project: %s\n", projectID)
	}

	if err = auth.SetCredential("google-antigravity", cred); err != nil {
		return fmt.Errorf("failed to save credentials: %w", err)
	}

	appCfg, err := internal.LoadConfig()
	if err == nil {
		// Update Providers (legacy format, for backward compatibility)
		appCfg.Providers.Antigravity.AuthMethod = "oauth"

		// Use shared function to add models
		addedCount := AddAntigravityModels(appCfg)

		if err := config.SaveConfig(internal.GetConfigPath(), appCfg); err != nil {
			fmt.Printf("Warning: could not update config: %v\n", err)
		} else {
			fmt.Printf("\n✓ Added %d Antigravity models to config\n", addedCount)
		}
	}

	fmt.Println("\n✓ Google Antigravity login successful!")
	fmt.Println("Default model set to: gemini-3-flash (fallback: gemini-2.5-flash)")
	fmt.Println("Available models:")
	fmt.Println("  - gemini-3-flash (default)")
	fmt.Println("  - gemini-3-pro-high, gemini-3-pro-low")
	fmt.Println("  - gemini-3.1-pro-high, gemini-3.1-pro-low, gemini-3.1-flash-lite")
	fmt.Println("  - gemini-3-flash-agent, gemini-3-flash-preview")
	fmt.Println("  - gemini-2.5-flash, gemini-2.5-flash-lite, gemini-2.5-flash-thinking, gemini-2.5-pro")
	fmt.Println("  - claude-sonnet-4-6, claude-opus-4-6-thinking")
	fmt.Println("  - gpt-oss-120b-medium")
	fmt.Println("\nTry it: picoclaw-agents agent -m \"Hello world\" --model gemini-3-flash")

	return nil
}

func fetchGoogleUserEmail(accessToken string) (string, error) {
	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("userinfo request failed: %s", string(body))
	}

	var userInfo struct {
		Email string `json:"email"`
	}
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return "", err
	}
	return userInfo.Email, nil
}

func authLoginPasteToken(provider string) error {
	cred, err := auth.LoginPasteToken(provider, os.Stdin)
	if err != nil {
		return fmt.Errorf("login failed: %w", err)
	}

	if err = auth.SetCredential(provider, cred); err != nil {
		return fmt.Errorf("failed to save credentials: %w", err)
	}

	appCfg, err := internal.LoadConfig()
	if err == nil {
		switch provider {
		case "anthropic":
			appCfg.Providers.Anthropic.AuthMethod = "token"
			// Update ModelList
			found := false
			for i := range appCfg.ModelList {
				if isAnthropicModel(appCfg.ModelList[i].Model) {
					appCfg.ModelList[i].AuthMethod = "token"
					found = true
					break
				}
			}
			if !found {
				appCfg.ModelList = append(appCfg.ModelList, config.ModelConfig{
					ModelName:  "claude-sonnet-4.6",
					Model:      "anthropic/claude-sonnet-4.6",
					AuthMethod: "token",
				})
			}
			// Update default model
			appCfg.Agents.Defaults.ModelName = "claude-sonnet-4.6"
		case "openai":
			appCfg.Providers.OpenAI.AuthMethod = "token"
			// Update ModelList
			found := false
			for i := range appCfg.ModelList {
				if isOpenAIModel(appCfg.ModelList[i].Model) {
					appCfg.ModelList[i].AuthMethod = "token"
					found = true
					break
				}
			}
			if !found {
				appCfg.ModelList = append(appCfg.ModelList, config.ModelConfig{
					ModelName:  "gpt-5.2",
					Model:      "openai/gpt-5.2",
					AuthMethod: "token",
				})
			}
			// Update default model
			appCfg.Agents.Defaults.ModelName = "gpt-5.2"
		}
		if err := config.SaveConfig(internal.GetConfigPath(), appCfg); err != nil {
			return fmt.Errorf("could not update config: %w", err)
		}
	}

	fmt.Printf("Token saved for %s!\n", provider)

	if appCfg != nil {
		fmt.Printf("Default model set to: %s\n", appCfg.Agents.Defaults.GetModelName())
	}

	return nil
}

func authLogoutCmd(provider string) error {
	if provider != "" {
		if err := auth.DeleteCredential(provider); err != nil {
			return fmt.Errorf("failed to remove credentials: %w", err)
		}

		appCfg, err := internal.LoadConfig()
		if err == nil {
			// Clear AuthMethod in ModelList
			for i := range appCfg.ModelList {
				switch provider {
				case "openai":
					if isOpenAIModel(appCfg.ModelList[i].Model) {
						appCfg.ModelList[i].AuthMethod = ""
					}
				case "anthropic":
					if isAnthropicModel(appCfg.ModelList[i].Model) {
						appCfg.ModelList[i].AuthMethod = ""
					}
				case "google-antigravity", "antigravity":
					if isAntigravityModel(appCfg.ModelList[i].Model) {
						appCfg.ModelList[i].AuthMethod = ""
					}
				}
			}
			// Clear AuthMethod in Providers (legacy)
			switch provider {
			case "openai":
				appCfg.Providers.OpenAI.AuthMethod = ""
			case "anthropic":
				appCfg.Providers.Anthropic.AuthMethod = ""
			case "google-antigravity", "antigravity":
				appCfg.Providers.Antigravity.AuthMethod = ""
			}
			config.SaveConfig(internal.GetConfigPath(), appCfg)
		}

		fmt.Printf("Logged out from %s\n", provider)

		return nil
	}

	if err := auth.DeleteAllCredentials(); err != nil {
		return fmt.Errorf("failed to remove credentials: %w", err)
	}

	appCfg, err := internal.LoadConfig()
	if err == nil {
		// Clear all AuthMethods in ModelList
		for i := range appCfg.ModelList {
			appCfg.ModelList[i].AuthMethod = ""
		}
		// Clear all AuthMethods in Providers (legacy)
		appCfg.Providers.OpenAI.AuthMethod = ""
		appCfg.Providers.Anthropic.AuthMethod = ""
		appCfg.Providers.Antigravity.AuthMethod = ""
		config.SaveConfig(internal.GetConfigPath(), appCfg)
	}

	fmt.Println("Logged out from all providers")

	return nil
}

func authStatusCmd() error {
	store, err := auth.LoadStore()
	if err != nil {
		return fmt.Errorf("failed to load auth store: %w", err)
	}

	if len(store.Credentials) == 0 {
		fmt.Println("No authenticated providers.")
		fmt.Println("Run: picoclaw auth login --provider <name>")
		return nil
	}

	fmt.Println("\nAuthenticated Providers:")
	fmt.Println("------------------------")
	for provider, cred := range store.Credentials {
		// For OAuth providers with an expired or soon-to-expire token, attempt a
		// silent refresh so the displayed status reflects the real usable state.
		if cred.AuthMethod == "oauth" && (cred.IsExpired() || cred.NeedsRefresh()) && cred.RefreshToken != "" {
			cfg := oauthConfigForProvider(provider)
			if refreshed, refreshErr := auth.RefreshAccessToken(cred, cfg); refreshErr == nil {
				refreshed.Email = cred.Email
				if refreshed.ProjectID == "" {
					refreshed.ProjectID = cred.ProjectID
				}
				if saveErr := auth.SetCredential(provider, refreshed); saveErr == nil {
					cred = refreshed
				}
			}
		}

		status := "active"
		if cred.IsExpired() {
			status = "expired"
		} else if cred.NeedsRefresh() {
			status = "needs refresh"
		}

		fmt.Printf("  %s:\n", provider)
		fmt.Printf("    Method: %s\n", cred.AuthMethod)
		fmt.Printf("    Status: %s\n", status)
		if cred.AccountID != "" {
			fmt.Printf("    Account: %s\n", cred.AccountID)
		}
		if cred.Email != "" {
			fmt.Printf("    Email: %s\n", cred.Email)
		}
		if cred.ProjectID != "" {
			fmt.Printf("    Project: %s\n", cred.ProjectID)
		}
		if !cred.ExpiresAt.IsZero() {
			fmt.Printf("    Expires: %s\n", cred.ExpiresAt.Format("2006-01-02 15:04"))
		}
	}

	return nil
}

// oauthConfigForProvider returns the OAuth config for known providers.
// Returns an empty config for unknown providers (refresh will be skipped by
// RefreshAccessToken when client_id is missing).
func oauthConfigForProvider(provider string) auth.OAuthProviderConfig {
	switch provider {
	case "google-antigravity", "antigravity":
		return auth.GoogleAntigravityOAuthConfig()
	case "openai":
		return auth.OpenAIOAuthConfig()
	default:
		return auth.OAuthProviderConfig{}
	}
}

func authModelsCmd() error {
	cred, err := auth.GetCredential("google-antigravity")
	if err != nil || cred == nil {
		return fmt.Errorf(
			"not logged in to Google Antigravity.\nrun: picoclaw auth login --provider google-antigravity",
		)
	}

	// Refresh token if it's about to expire OR already expired (consistent with provider behavior).
	// Previously only NeedsRefresh() was checked — if token expired during inactivity
	// (>1h idle), `auth models` would fail with a 401 instead of auto-refreshing.
	if (cred.NeedsRefresh() || cred.IsExpired()) && cred.RefreshToken != "" {
		oauthCfg := auth.GoogleAntigravityOAuthConfig()
		refreshed, refreshErr := auth.RefreshAccessToken(cred, oauthCfg)
		if refreshErr == nil {
			cred = refreshed
			_ = auth.SetCredential("google-antigravity", cred)
		}
	}

	projectID := cred.ProjectID
	if projectID == "" {
		return fmt.Errorf("no project id stored. Try logging in again")
	}

	fmt.Printf("Fetching models for project: %s\n\n", projectID)

	models, err := providers.FetchAntigravityModels(cred.AccessToken, projectID)
	if err != nil {
		return fmt.Errorf("error fetching models: %w", err)
	}

	if len(models) == 0 {
		return fmt.Errorf("no models available")
	}

	fmt.Println("Available Antigravity Models:")
	fmt.Println("-----------------------------")
	for _, m := range models {
		status := "✓"
		if m.IsExhausted {
			status = "✗ (quota exhausted)"
		}
		name := m.ID
		if m.DisplayName != "" {
			name = fmt.Sprintf("%s (%s)", m.ID, m.DisplayName)
		}
		fmt.Printf("  %s %s\n", status, name)
	}

	return nil
}

// isAntigravityModel checks if a model string belongs to antigravity provider
func isAntigravityModel(model string) bool {
	return model == "antigravity" ||
		model == "google-antigravity" ||
		strings.HasPrefix(model, "antigravity/") ||
		strings.HasPrefix(model, "google-antigravity/")
}

// isOpenAIModel checks if a model string belongs to openai provider
func isOpenAIModel(model string) bool {
	return model == "openai" ||
		strings.HasPrefix(model, "openai/")
}

// isAnthropicModel checks if a model string belongs to anthropic provider
func isAnthropicModel(model string) bool {
	return model == "anthropic" ||
		strings.HasPrefix(model, "anthropic/")
}

// AddAnthropicModels adds Anthropic models to config with deduplication
func AddAnthropicModels(appCfg *config.Config) int {
	anthropicModels := []config.ModelConfig{
		{ModelName: "claude-sonnet-4-6", Model: "anthropic/claude-sonnet-4-6", AuthMethod: "oauth"},
		{ModelName: "claude-opus-4-6", Model: "anthropic/claude-opus-4-6", AuthMethod: "oauth"},
		{ModelName: "claude-opus-4-6-thinking", Model: "anthropic/claude-opus-4-6-thinking", AuthMethod: "oauth"},
		{ModelName: "claude-3-5-sonnet", Model: "anthropic/claude-3-5-sonnet", AuthMethod: "oauth"},
		{ModelName: "claude-3-5-haiku", Model: "anthropic/claude-3-5-haiku", AuthMethod: "oauth"},
	}

	existingModels := make(map[string]bool)
	for _, m := range appCfg.ModelList {
		existingModels[m.ModelName] = true
	}

	addedCount := 0
	for _, modelCfg := range anthropicModels {
		if !existingModels[modelCfg.ModelName] {
			appCfg.ModelList = append(appCfg.ModelList, modelCfg)
			existingModels[modelCfg.ModelName] = true
			addedCount++
		}
	}

	if addedCount > 0 {
		appCfg.Agents.Defaults.ModelName = "claude-sonnet-4-6"
	}

	return addedCount
}

// AddAntigravityModels adds Antigravity models to config with deduplication
func AddAntigravityModels(appCfg *config.Config) int {
	antigravityModels := []config.ModelConfig{
		{ModelName: "gemini-3-flash", Model: "antigravity/gemini-3-flash", AuthMethod: "oauth"},
		{ModelName: "gemini-3-pro-high", Model: "antigravity/gemini-3-pro-high", AuthMethod: "oauth"},
		{ModelName: "gemini-3-pro-low", Model: "antigravity/gemini-3-pro-low", AuthMethod: "oauth"},
		{ModelName: "gemini-3.1-pro-high", Model: "antigravity/gemini-3.1-pro-high", AuthMethod: "oauth"},
		{ModelName: "gemini-3.1-pro-low", Model: "antigravity/gemini-3.1-pro-low", AuthMethod: "oauth"},
		{ModelName: "gemini-3.1-flash-lite", Model: "antigravity/gemini-3.1-flash-lite", AuthMethod: "oauth"},
		{ModelName: "gemini-3-flash-agent", Model: "antigravity/gemini-3-flash-agent", AuthMethod: "oauth"},
		{ModelName: "gemini-3-flash-preview", Model: "antigravity/gemini-3-flash-preview", AuthMethod: "oauth"},
		{ModelName: "gemini-2.5-flash", Model: "antigravity/gemini-2.5-flash", AuthMethod: "oauth"},
		{ModelName: "gemini-2.5-flash-lite", Model: "antigravity/gemini-2.5-flash-lite", AuthMethod: "oauth"},
		{ModelName: "gemini-2.5-flash-thinking", Model: "antigravity/gemini-2.5-flash-thinking", AuthMethod: "oauth"},
		{ModelName: "gemini-2.5-pro", Model: "antigravity/gemini-2.5-pro", AuthMethod: "oauth"},
		{ModelName: "claude-sonnet-4-6", Model: "antigravity/claude-sonnet-4-6", AuthMethod: "oauth"},
		{ModelName: "claude-opus-4-6-thinking", Model: "antigravity/claude-opus-4-6-thinking", AuthMethod: "oauth"},
		{ModelName: "gpt-oss-120b-medium", Model: "antigravity/gpt-oss-120b-medium", AuthMethod: "oauth"},
	}

	existingModels := make(map[string]bool)
	for _, m := range appCfg.ModelList {
		existingModels[m.ModelName] = true
	}

	addedCount := 0
	for _, modelCfg := range antigravityModels {
		if !existingModels[modelCfg.ModelName] {
			appCfg.ModelList = append(appCfg.ModelList, modelCfg)
			existingModels[modelCfg.ModelName] = true
			addedCount++
		}
	}

	if addedCount > 0 {
		appCfg.Agents.Defaults.ModelName = "gemini-3-flash"
		for i := range appCfg.Agents.List {
			if appCfg.Agents.List[i].Model == nil {
				appCfg.Agents.List[i].Model = &config.AgentModelConfig{}
			}
			appCfg.Agents.List[i].Model.Primary = "gemini-3-flash"
			appCfg.Agents.List[i].Model.Fallbacks = []string{"gemini-2.5-flash"}
		}
	}

	return addedCount
}

// AddOpenAIModels adds OpenAI models to config with deduplication
func AddOpenAIModels(appCfg *config.Config) int {
	openAIModels := []config.ModelConfig{
		{ModelName: "gpt-5.4", Model: "openai/gpt-5.4", AuthMethod: "oauth"},
		{ModelName: "gpt-5", Model: "openai/gpt-5", AuthMethod: "oauth"},
		{ModelName: "o3-mini", Model: "openai/o3-mini", AuthMethod: "oauth"},
		{ModelName: "o3", Model: "openai/o3", AuthMethod: "oauth"},
		{ModelName: "o1", Model: "openai/o1", AuthMethod: "oauth"},
		{ModelName: "o1-mini", Model: "openai/o1-mini", AuthMethod: "oauth"},
		{ModelName: "gpt-4.1", Model: "openai/gpt-4.1", AuthMethod: "oauth"},
		{ModelName: "gpt-4-turbo", Model: "openai/gpt-4-turbo", AuthMethod: "oauth"},
	}

	existingModels := make(map[string]bool)
	for _, m := range appCfg.ModelList {
		existingModels[m.ModelName] = true
	}

	addedCount := 0
	for _, modelCfg := range openAIModels {
		if !existingModels[modelCfg.ModelName] {
			appCfg.ModelList = append(appCfg.ModelList, modelCfg)
			existingModels[modelCfg.ModelName] = true
			addedCount++
		}
	}

	if addedCount > 0 {
		appCfg.Agents.Defaults.ModelName = "gpt-5.4"
	}

	return addedCount
}
