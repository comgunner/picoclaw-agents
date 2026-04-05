// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
)

func loadTemplateConfigFromExample() (*Config, bool) {
	candidates := []string{
		filepath.Join("config", "config.example.json"),
		filepath.Join("..", "..", "config", "config.example.json"),
	}

	if _, file, _, ok := runtime.Caller(0); ok {
		candidates = append(candidates, filepath.Join(filepath.Dir(file), "..", "..", "config", "config.example.json"))
	}

	for _, p := range candidates {
		data, err := os.ReadFile(p)
		if err != nil {
			continue
		}
		cfg := &Config{}
		if err := json.Unmarshal(data, cfg); err != nil {
			continue
		}
		return cfg, true
	}

	return nil, false
}

// DefaultConfig returns the minimal default configuration for PicoClaw.
// This is used as a base for internal state and migrations.
//
// ⚠️  CRITICAL: ContextManager is set to "seahorse" by default. DO NOT remove
// this. It prevents OpenRouter Free tier 402 errors by enabling budget-aware
// context assembly. All template functions (GLMDefaultConfig, OpenAIDefaultConfig,
// etc.) inherit from TemplateDefaultConfig() which also sets this default.
func DefaultConfig() *Config {
	return &Config{
		Agents: AgentsConfig{
			Defaults: AgentDefaults{
				Workspace:           "~/.picoclaw/workspace",
				RestrictToWorkspace: true,
				Provider:            "",
				Model:               "",
				MaxTokens:           32768,
				Temperature:         nil, // nil means use provider default
				MaxToolIterations:   20,
				ContextManager:      "seahorse",
				ContextManagerConfig: map[string]any{
					"context_threshold":       0.75,
					"fresh_tail_count":        16,
					"leaf_target_tokens":      1200,
					"condensed_target_tokens": 2000,
					"max_compact_iterations":  20,
				},
			},
			List: []AgentConfig{
				{
					ID:      "tech_lead",
					Name:    "Tech Lead",
					Default: true,
					Model: &AgentModelConfig{
						Primary: "antigravity-gemini-3-flash",
					},
					Subagents: &SubagentsConfig{
						AllowAgents:         []string{"*"},
						MaxSpawnDepth:       3,
						MaxChildrenPerAgent: 5,
					},
				},
				{
					ID:   "backend_coder",
					Name: "Backend Coder",
					Model: &AgentModelConfig{
						Primary: "antigravity-gemini-3-flash",
					},
					Subagents: &SubagentsConfig{
						AllowAgents:   []string{"data_researcher"},
						MaxSpawnDepth: 1,
					},
				},
				{
					ID:   "data_researcher",
					Name: "Data Researcher",
					Model: &AgentModelConfig{
						Primary: "antigravity-gemini-3-flash",
					},
					Subagents: &SubagentsConfig{
						AllowAgents: []string{},
					},
				},
			},
		},
		Bindings: []AgentBinding{},
		Session: SessionConfig{
			DMScope: "per-channel-peer",
		},
		Channels: ChannelsConfig{
			WhatsApp: WhatsAppConfig{
				Enabled:   false,
				BridgeURL: "ws://localhost:3001",
				AllowFrom: FlexibleStringSlice{},
			},
			Telegram: TelegramConfig{
				Enabled:   false,
				Token:     "",
				AllowFrom: FlexibleStringSlice{},
			},
			Feishu: FeishuConfig{
				Enabled:   false,
				AllowFrom: FlexibleStringSlice{},
			},
			Discord: DiscordConfig{
				Enabled:     false,
				Token:       "",
				AllowFrom:   FlexibleStringSlice{},
				MentionOnly: false,
			},
			MaixCam: MaixCamConfig{
				Enabled:   false,
				Host:      "0.0.0.0",
				Port:      18790,
				AllowFrom: FlexibleStringSlice{},
			},
			QQ: QQConfig{
				Enabled:   false,
				AppID:     "",
				AppSecret: "",
				AllowFrom: FlexibleStringSlice{},
			},
			DingTalk: DingTalkConfig{
				Enabled:      false,
				ClientID:     "",
				ClientSecret: "",
				AllowFrom:    FlexibleStringSlice{},
			},
			Slack: SlackConfig{
				Enabled:   false,
				BotToken:  "",
				AppToken:  "",
				AllowFrom: FlexibleStringSlice{},
			},
			LINE: LINEConfig{
				Enabled:            false,
				ChannelSecret:      "",
				ChannelAccessToken: "",
				WebhookHost:        "0.0.0.0",
				WebhookPort:        18791,
				WebhookPath:        "/webhook/line",
				AllowFrom:          FlexibleStringSlice{},
			},
			OneBot: OneBotConfig{
				Enabled:           false,
				WSUrl:             "ws://127.0.0.1:3001",
				ReconnectInterval: 5,
				AllowFrom:         FlexibleStringSlice{},
			},
			WeCom: WeComConfig{
				Enabled:        false,
				Token:          "",
				EncodingAESKey: "",
				WebhookURL:     "",
				WebhookHost:    "0.0.0.0",
				WebhookPort:    18793,
				WebhookPath:    "/webhook/wecom",
				AllowFrom:      FlexibleStringSlice{},
				ReplyTimeout:   5,
			},
			WeComApp: WeComAppConfig{
				Enabled:        false,
				CorpID:         "",
				CorpSecret:     "",
				AgentID:        0,
				Token:          "",
				EncodingAESKey: "",
				WebhookHost:    "0.0.0.0",
				WebhookPort:    18792,
				WebhookPath:    "/webhook/wecom-app",
				AllowFrom:      FlexibleStringSlice{},
				ReplyTimeout:   5,
			},
		},
		Providers: ProvidersConfig{
			OpenAI: OpenAIProviderConfig{WebSearch: true},
		},
		ModelList: []ModelConfig{
			{
				ModelName: "default",
				Model:     "openai/gpt-4o",
			},
		},
		Gateway: GatewayConfig{
			Host: "127.0.0.1",
			Port: 18790,
		},
		Tools: ToolsConfig{
			Web: WebToolsConfig{
				Proxy: "",
				Brave: BraveConfig{
					Enabled:    false,
					APIKey:     "",
					MaxResults: 5,
				},
				DuckDuckGo: DuckDuckGoConfig{
					Enabled:    true,
					MaxResults: 5,
				},
				Perplexity: PerplexityConfig{
					Enabled:    false,
					APIKey:     "",
					MaxResults: 5,
				},
			},
			Cron: CronToolsConfig{
				ExecTimeoutMinutes: 5,
			},
			Exec: ExecConfig{
				Enabled:            true,
				AllowRemote:        true,
				EnableDenyPatterns: true,
			},
			Skills: SkillsToolsConfig{
				Registries: SkillsRegistriesConfig{
					ClawHub: ClawHubRegistryConfig{
						Enabled: true,
						BaseURL: "https://clawhub.ai",
					},
				},
				MaxConcurrentSearches: 2,
				SearchCache: SearchCacheConfig{
					MaxSize:    50,
					TTLSeconds: 300,
				},
			},
		},
		Heartbeat: HeartbeatConfig{
			Enabled:  true,
			Interval: 30,
		},
		Devices: DevicesConfig{
			Enabled:    false,
			MonitorUSB: true,
		},
		ContextManagement: ContextManagementConfig{
			CompactThreshold:    0.75,
			CriticalThreshold:   0.90,
			MinCompletionTokens: 1024, // SPRINT 1: aumentado de 512 a 1024
			PreserveMessages:    30,   // SPRINT 1: aumentado de 20 a 30
			AutoCompactEnabled:  true,

			// SPRINT 1 FEATURE: Pruning avanzado
			Pruning: ContextPruningConfig{
				Enabled:            true,
				MaxToolResultChars: 8000, // truncar tool results > 8K chars
				ExcludeTools:       []string{"memory_store", "memory_read"},
				AggressiveTools:    []string{"shell", "web_fetch"},
			},

			// SPRINT 1 FEATURE: Compaction avanzado
			Compaction: ContextCompactionConfig{
				Model:               "",   // vacío = usar mismo modelo del agente
				MaxSummaryTokens:    2048, // aumentado de 512 a 2048 (4x más contexto)
				RecentTurnsPreserve: 6,    // preservar últimos 6 turnos verbatim
				MinSummaryQuality:   0.0,  // desactivado por default
				MaxRetries:          2,
			},
		},
	}
}

// TemplateDefaultConfig returns a configuration matching config.example.json.
func TemplateDefaultConfig() *Config {
	if cfg, ok := loadTemplateConfigFromExample(); ok {
		return cfg
	}

	return &Config{
		Agents: AgentsConfig{
			Defaults: AgentDefaults{
				Workspace:           "~/.picoclaw/workspace",
				RestrictToWorkspace: true,
				Provider:            "",
				Model:               "antigravity-gemini-3-flash",
				MaxTokens:           8192,
				Temperature:         nil,
				MaxToolIterations:   20,
				ContextManager:      "seahorse",
				ContextManagerConfig: map[string]any{
					"context_threshold":       0.75,
					"fresh_tail_count":        16,
					"leaf_target_tokens":      1200,
					"condensed_target_tokens": 2000,
					"max_compact_iterations":  20,
				},
			},
			List: []AgentConfig{
				{
					ID:      "project_manager",
					Name:    "Project Manager",
					Default: true,
					Model: &AgentModelConfig{
						Primary: "antigravity-gemini-3-flash",
					},
					Subagents: &SubagentsConfig{
						AllowAgents:         []string{"*", "general_worker"},
						MaxSpawnDepth:       3,
						MaxChildrenPerAgent: 5,
					},
				},
				{
					ID:   "senior_dev",
					Name: "Senior Developer",
					Model: &AgentModelConfig{
						Primary: "antigravity-gemini-3-flash",
					},
					Subagents: &SubagentsConfig{
						AllowAgents:   []string{"qa_specialist", "junior_fixer"},
						MaxSpawnDepth: 2,
					},
				},
				{
					ID:   "qa_specialist",
					Name: "QA Specialist (GitHub Ops)",
					Model: &AgentModelConfig{
						Primary: "antigravity-gemini-3-flash",
					},
					Subagents: &SubagentsConfig{
						AllowAgents: []string{},
					},
				},
				{
					ID:   "junior_fixer",
					Name: "Junior Code Fixer",
					Model: &AgentModelConfig{
						Primary: "antigravity-gemini-3-flash",
					},
					Subagents: &SubagentsConfig{
						AllowAgents: []string{},
					},
				},
				{
					ID:   "general_worker",
					Name: "General Worker",
					Model: &AgentModelConfig{
						Primary: "antigravity-gemini-3-flash",
					},
					Subagents: &SubagentsConfig{
						AllowAgents: []string{},
					},
				},
			},
		},
		Session: SessionConfig{
			DMScope: "per-channel-peer",
		},
		Channels: ChannelsConfig{
			Telegram: TelegramConfig{
				Enabled:   false,
				Token:     "YOUR_TELEGRAM_BOT_TOKEN",
				AllowFrom: FlexibleStringSlice{"YOUR_USER_ID"},
			},
			Discord: DiscordConfig{
				Enabled:   false,
				Token:     "YOUR_DISCORD_BOT_TOKEN",
				AllowFrom: FlexibleStringSlice{},
			},
			QQ: QQConfig{
				Enabled:   false,
				AppID:     "YOUR_QQ_APP_ID",
				AppSecret: "YOUR_QQ_APP_SECRET",
				AllowFrom: FlexibleStringSlice{},
			},
			MaixCam: MaixCamConfig{
				Enabled:   false,
				Host:      "0.0.0.0",
				Port:      18790,
				AllowFrom: FlexibleStringSlice{},
			},
			WhatsApp: WhatsAppConfig{
				Enabled:   false,
				BridgeURL: "ws://localhost:3001",
				AllowFrom: FlexibleStringSlice{},
			},
			Feishu: FeishuConfig{
				Enabled:   false,
				AllowFrom: FlexibleStringSlice{},
			},
			DingTalk: DingTalkConfig{
				Enabled:      false,
				ClientID:     "YOUR_CLIENT_ID",
				ClientSecret: "YOUR_CLIENT_SECRET",
				AllowFrom:    FlexibleStringSlice{},
			},
			Slack: SlackConfig{
				Enabled:   false,
				BotToken:  "xoxb-YOUR-BOT-TOKEN",
				AppToken:  "xapp-YOUR-APP-TOKEN",
				AllowFrom: FlexibleStringSlice{},
			},
			LINE: LINEConfig{
				Enabled:            false,
				ChannelSecret:      "YOUR_LINE_CHANNEL_SECRET",
				ChannelAccessToken: "YOUR_LINE_CHANNEL_ACCESS_TOKEN",
				WebhookHost:        "0.0.0.0",
				WebhookPort:        18791,
				WebhookPath:        "/webhook/line",
				AllowFrom:          FlexibleStringSlice{},
			},
			OneBot: OneBotConfig{
				Enabled:            false,
				WSUrl:              "ws://127.0.0.1:3001",
				ReconnectInterval:  5,
				GroupTriggerPrefix: []string{},
				AllowFrom:          FlexibleStringSlice{},
			},
			WeCom: WeComConfig{
				Enabled:        false,
				Token:          "YOUR_TOKEN",
				EncodingAESKey: "YOUR_43_CHAR_ENCODING_AES_KEY",
				WebhookURL:     "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=YOUR_KEY",
				WebhookHost:    "0.0.0.0",
				WebhookPort:    18793,
				WebhookPath:    "/webhook/wecom",
				AllowFrom:      FlexibleStringSlice{},
				ReplyTimeout:   5,
			},
			WeComApp: WeComAppConfig{
				Enabled:        false,
				CorpID:         "YOUR_CORP_ID",
				CorpSecret:     "YOUR_CORP_SECRET",
				AgentID:        1000002,
				Token:          "YOUR_TOKEN",
				EncodingAESKey: "YOUR_43_CHAR_ENCODING_AES_KEY",
				WebhookHost:    "0.0.0.0",
				WebhookPort:    18792,
				WebhookPath:    "/webhook/wecom-app",
				AllowFrom:      FlexibleStringSlice{},
				ReplyTimeout:   5,
			},
		},
		ModelList: []ModelConfig{
			{
				ModelName: "glm-4.5-flash",
				Model:     "zhipu/glm-4.5-flash",
				APIBase:   "https://api.z.ai/api/paas/v4",
				APIKey:    "",
			},
			{
				ModelName: "gpt-5.2",
				Model:     "openai/gpt-5.2",
				APIBase:   "https://api.openai.com/v1",
				APIKey:    "",
			},
			{
				ModelName: "claude-sonnet-4.6",
				Model:     "anthropic/claude-sonnet-4.6",
				APIBase:   "https://api.anthropic.com/v1",
				APIKey:    "",
			},
			{
				ModelName: "deepseek-chat",
				Model:     "deepseek/deepseek-chat",
				APIBase:   "https://api.deepseek.com/v1",
				APIKey:    "",
			},
			{
				ModelName: "gemini-2.0-flash",
				Model:     "gemini/gemini-2.0-flash-exp",
				APIBase:   "https://generativelanguage.googleapis.com/v1beta/openai/",
				APIKey:    "",
			},
			{
				ModelName: "qwen-plus",
				Model:     "qwen/qwen-plus",
				APIBase:   "https://dashscope.aliyuncs.com/compatible-mode/v1",
				APIKey:    "",
			},
			{
				ModelName: "moonshot-v1-8k",
				Model:     "moonshot/moonshot-v1-8k",
				APIBase:   "https://api.moonshot.cn/v1",
				APIKey:    "",
			},
			{
				ModelName: "llama-3.3-70b",
				Model:     "groq/llama-3.3-70b-versatile",
				APIBase:   "https://api.groq.com/openai/v1",
				APIKey:    "",
			},
			{
				ModelName: "openrouter-auto",
				Model:     "openrouter/auto",
				APIBase:   "https://openrouter.ai/api/v1",
				APIKey:    "",
			},
			{
				ModelName: "openrouter-gpt-5.2",
				Model:     "openrouter/openai/gpt-5.2",
				APIBase:   "https://openrouter.ai/api/v1",
				APIKey:    "",
			},
			{
				ModelName: "nemotron-4-340b",
				Model:     "nvidia/nemotron-4-340b-instruct",
				APIBase:   "https://integrate.api.nvidia.com/v1",
				APIKey:    "",
			},
			{
				ModelName: "cerebras-llama-3.3-70b",
				Model:     "cerebras/llama-3.3-70b",
				APIBase:   "https://api.cerebras.ai/v1",
				APIKey:    "",
			},
			{
				ModelName: "doubao-pro",
				Model:     "volcengine/doubao-pro-32k",
				APIBase:   "https://ark.cn-beijing.volces.com/api/v3",
				APIKey:    "",
			},
			{
				ModelName: "deepseek-v3",
				Model:     "shengsuanyun/deepseek-v3",
				APIBase:   "https://api.shengsuanyun.com/v1",
				APIKey:    "",
			},
			{
				ModelName: "gemini-2.5-flash",
				Model:     "gemini/gemini-2.5-flash",
				APIKey:    "",
			},
			{
				ModelName: "gemini-2.5-flash-image",
				Model:     "gemini/gemini-2.5-flash-image",
				APIKey:    "",
			},
			{
				ModelName:  "gemini-3.1-pro-high",
				Model:      "antigravity/gemini-3.1-pro-high",
				APIKey:     "",
				AuthMethod: "oauth",
			},
			{
				ModelName:  "antigravity-gemini-3-flash",
				Model:      "antigravity/gemini-3-flash",
				APIKey:     "",
				AuthMethod: "oauth",
			},
			{
				ModelName:  "gemini-2.5-pro",
				Model:      "antigravity/gemini-2.5-pro",
				APIKey:     "",
				AuthMethod: "oauth",
			},
			{
				ModelName:  "claude-sonnet-4-6",
				Model:      "antigravity/claude-sonnet-4-6",
				APIKey:     "",
				AuthMethod: "oauth",
			},
			{
				ModelName:  "claude-opus-4-6-thinking",
				Model:      "antigravity/claude-opus-4-6-thinking",
				APIKey:     "",
				AuthMethod: "oauth",
			},
			{
				ModelName:  "copilot-gpt-5.2",
				Model:      "github-copilot/gpt-5.2",
				APIBase:    "http://localhost:4321",
				APIKey:     "",
				AuthMethod: "oauth",
			},
			{
				ModelName: "mistral-small",
				Model:     "mistral/mistral-small-latest",
				APIBase:   "https://api.mistral.ai/v1",
				APIKey:    "",
			},
			{
				ModelName: "local-model",
				Model:     "vllm/custom-model",
				APIBase:   "http://localhost:8000/v1",
				APIKey:    "",
			},
			{
				ModelName: "qwen2.5:0.5b",
				Model:     "qwen2.5:0.5b",
				APIBase:   "http://localhost:11434/v1",
				APIKey:    "ollama",
			},
			{
				ModelName: "qwen2.5-coder:0.5b",
				Model:     "qwen2.5-coder:0.5b",
				APIBase:   "http://localhost:11434/v1",
				APIKey:    "ollama",
			},
			{
				ModelName: "llama3.2:1b",
				Model:     "llama3.2:1b",
				APIBase:   "http://localhost:11434/v1",
				APIKey:    "ollama",
			},
		},
		Tools: ToolsConfig{
			Web: WebToolsConfig{
				Proxy: "",
				Brave: BraveConfig{
					Enabled:    false,
					APIKey:     "YOUR_BRAVE_API_KEY",
					MaxResults: 5,
				},
				DuckDuckGo: DuckDuckGoConfig{
					Enabled:    true,
					MaxResults: 5,
				},
				Perplexity: PerplexityConfig{
					Enabled:    false,
					APIKey:     "pplx-xxx",
					MaxResults: 5,
				},
			},
			Cron: CronToolsConfig{
				ExecTimeoutMinutes: 5,
			},
			Exec: ExecConfig{
				EnableDenyPatterns: false,
				CustomDenyPatterns: []string{},
			},
			Skills: SkillsToolsConfig{
				Registries: SkillsRegistriesConfig{
					ClawHub: ClawHubRegistryConfig{
						Enabled:      true,
						BaseURL:      "https://clawhub.ai",
						SearchPath:   "/api/v1/search",
						SkillsPath:   "/api/v1/skills",
						DownloadPath: "/api/v1/download",
					},
				},
				MaxConcurrentSearches: 2,
				SearchCache: SearchCacheConfig{
					MaxSize:    50,
					TTLSeconds: 300,
				},
			},
		},
		Heartbeat: HeartbeatConfig{
			Enabled:  true,
			Interval: 30,
		},
		Devices: DevicesConfig{
			Enabled:    false,
			MonitorUSB: true,
		},
		ContextManagement: ContextManagementConfig{
			CompactThreshold:    0.75,
			CriticalThreshold:   0.90,
			MinCompletionTokens: 512,
			PreserveMessages:    20,
			AutoCompactEnabled:  true,
		},
		Gateway: GatewayConfig{
			Host: "127.0.0.1",
			Port: 18790,
		},
	}
}

// GLMDefaultConfig returns a configuration template optimized for GLM models.
func GLMDefaultConfig() *Config {
	cfg := TemplateDefaultConfig()
	cfg.Agents.Defaults.Model = "glm-4.5-flash"
	cfg.Agents.Defaults.MaxTokens = 32768
	for i := range cfg.Agents.List {
		cfg.Agents.List[i].Model.Primary = "glm-4.5-flash"
	}
	cfg.ModelList = []ModelConfig{
		{
			ModelName: "glm-4.5-flash",
			Model:     "zhipu/glm-4.5-flash",
			APIBase:   "https://api.z.ai/api/paas/v4",
			APIKey:    "",
		},
		{
			ModelName: "gpt-5.2",
			Model:     "openai/gpt-5.2",
			APIBase:   "https://api.openai.com/v1",
			APIKey:    "",
		},
		{
			ModelName: "claude-sonnet-4.6",
			Model:     "anthropic/claude-sonnet-4.6",
			APIBase:   "https://api.anthropic.com/v1",
			APIKey:    "",
		},
		{
			ModelName: "deepseek-chat",
			Model:     "deepseek/deepseek-chat",
			APIBase:   "https://api.deepseek.com/v1",
			APIKey:    "",
		},
		{
			ModelName: "gemini-2.0-flash",
			Model:     "gemini/gemini-2.0-flash-exp",
			APIBase:   "https://generativelanguage.googleapis.com/v1beta/openai/",
			APIKey:    "",
		},
		{
			ModelName: "qwen-plus",
			Model:     "qwen/qwen-plus",
			APIBase:   "https://dashscope.aliyuncs.com/compatible-mode/v1",
			APIKey:    "",
		},
		{
			ModelName: "moonshot-v1-8k",
			Model:     "moonshot/moonshot-v1-8k",
			APIBase:   "https://api.moonshot.cn/v1",
			APIKey:    "",
		},
		{
			ModelName: "llama-3.3-70b",
			Model:     "groq/llama-3.3-70b-versatile",
			APIBase:   "https://api.groq.com/openai/v1",
			APIKey:    "",
		},
		{
			ModelName: "openrouter-auto",
			Model:     "openrouter/auto",
			APIBase:   "https://openrouter.ai/api/v1",
			APIKey:    "",
		},
		{
			ModelName: "openrouter-gpt-5.2",
			Model:     "openrouter/openai/gpt-5.2",
			APIBase:   "https://openrouter.ai/api/v1",
			APIKey:    "",
		},
		{
			ModelName: "nemotron-4-340b",
			Model:     "nvidia/nemotron-4-340b-instruct",
			APIBase:   "https://integrate.api.nvidia.com/v1",
			APIKey:    "",
		},
		{
			ModelName: "cerebras-llama-3.3-70b",
			Model:     "cerebras/llama-3.3-70b",
			APIBase:   "https://api.cerebras.ai/v1",
			APIKey:    "",
		},
		{
			ModelName: "doubao-pro",
			Model:     "volcengine/doubao-pro-32k",
			APIBase:   "https://ark.cn-beijing.volces.com/api/v3",
			APIKey:    "",
		},
		{
			ModelName: "deepseek-v3",
			Model:     "shengsuanyun/deepseek-v3",
			APIBase:   "https://api.shengsuanyun.com/v1",
			APIKey:    "",
		},
		{
			ModelName: "gemini-2.5-flash",
			Model:     "gemini/gemini-2.5-flash",
			APIKey:    "",
		},
		{
			ModelName: "gemini-2.5-flash-image",
			Model:     "gemini/gemini-2.5-flash-image",
			APIKey:    "",
		},
		{
			ModelName:  "gemini-3.1-pro-high",
			Model:      "antigravity/gemini-3.1-pro-high",
			APIKey:     "",
			AuthMethod: "oauth",
		},
		{
			ModelName:  "antigravity-gemini-3-flash",
			Model:      "antigravity/gemini-3-flash",
			APIKey:     "",
			AuthMethod: "oauth",
		},
		{
			ModelName:  "gemini-2.5-pro",
			Model:      "antigravity/gemini-2.5-pro",
			APIKey:     "",
			AuthMethod: "oauth",
		},
		{
			ModelName:  "claude-sonnet-4-6",
			Model:      "antigravity/claude-sonnet-4-6",
			APIKey:     "",
			AuthMethod: "oauth",
		},
		{
			ModelName:  "claude-opus-4-6-thinking",
			Model:      "antigravity/claude-opus-4-6-thinking",
			APIKey:     "",
			AuthMethod: "oauth",
		},
		{
			ModelName:  "copilot-gpt-5.2",
			Model:      "github-copilot/gpt-5.2",
			APIBase:    "http://localhost:4321",
			APIKey:     "",
			AuthMethod: "oauth",
		},
		{
			ModelName: "mistral-small",
			Model:     "mistral/mistral-small-latest",
			APIBase:   "https://api.mistral.ai/v1",
			APIKey:    "",
		},
		{
			ModelName: "local-model",
			Model:     "vllm/custom-model",
			APIBase:   "http://localhost:8000/v1",
			APIKey:    "",
		},
		{
			ModelName: "qwen2.5:0.5b",
			Model:     "qwen2.5:0.5b",
			APIBase:   "http://localhost:11434/v1",
			APIKey:    "ollama",
		},
		{
			ModelName: "qwen2.5-coder:0.5b",
			Model:     "qwen2.5-coder:0.5b",
			APIBase:   "http://localhost:11434/v1",
			APIKey:    "ollama",
		},
		{
			ModelName: "llama3.2:1b",
			Model:     "llama3.2:1b",
			APIBase:   "http://localhost:11434/v1",
			APIKey:    "ollama",
		},
	}
	return cfg
}

// OpenAIDefaultConfig returns a configuration template optimized for OpenAI models.
func OpenAIDefaultConfig() *Config {
	cfg := TemplateDefaultConfig()
	cfg.Agents.Defaults.Model = "o3-mini-2025-01-31"
	cfg.Agents.Defaults.MaxTokens = 8192
	for i := range cfg.Agents.List {
		cfg.Agents.List[i].Model.Primary = "o3-mini-2025-01-31"
	}
	cfg.ModelList = []ModelConfig{
		{
			ModelName: "o3-mini-2025-01-31",
			Model:     "openai/o3-mini-2025-01-31",
			APIBase:   "https://api.openai.com/v1",
			APIKey:    "YOUR_OPENAI_API_KEY",
		},
		{
			ModelName: "glm-4.5-flash",
			Model:     "zhipu/glm-4.5-flash",
			APIBase:   "https://api.z.ai/api/paas/v4",
			APIKey:    "",
		},
		{
			ModelName: "gpt-5.2",
			Model:     "openai/gpt-5.2",
			APIBase:   "https://api.openai.com/v1",
			APIKey:    "",
		},
		{
			ModelName: "claude-sonnet-4.6",
			Model:     "anthropic/claude-sonnet-4.6",
			APIBase:   "https://api.anthropic.com/v1",
			APIKey:    "",
		},
		{
			ModelName: "deepseek-chat",
			Model:     "deepseek/deepseek-chat",
			APIBase:   "https://api.deepseek.com/v1",
			APIKey:    "",
		},
		{
			ModelName: "gemini-2.0-flash",
			Model:     "gemini/gemini-2.0-flash-exp",
			APIBase:   "https://generativelanguage.googleapis.com/v1beta/openai/",
			APIKey:    "",
		},
		{
			ModelName: "qwen-plus",
			Model:     "qwen/qwen-plus",
			APIBase:   "https://dashscope.aliyuncs.com/compatible-mode/v1",
			APIKey:    "",
		},
		{
			ModelName: "moonshot-v1-8k",
			Model:     "moonshot/moonshot-v1-8k",
			APIBase:   "https://api.moonshot.cn/v1",
			APIKey:    "",
		},
		{
			ModelName: "llama-3.3-70b",
			Model:     "groq/llama-3.3-70b-versatile",
			APIBase:   "https://api.groq.com/openai/v1",
			APIKey:    "",
		},
		{
			ModelName: "openrouter-auto",
			Model:     "openrouter/auto",
			APIBase:   "https://openrouter.ai/api/v1",
			APIKey:    "",
		},
		{
			ModelName: "openrouter-gpt-5.2",
			Model:     "openrouter/openai/gpt-5.2",
			APIBase:   "https://openrouter.ai/api/v1",
			APIKey:    "",
		},
		{
			ModelName: "nemotron-4-340b",
			Model:     "nvidia/nemotron-4-340b-instruct",
			APIBase:   "https://integrate.api.nvidia.com/v1",
			APIKey:    "",
		},
		{
			ModelName: "cerebras-llama-3.3-70b",
			Model:     "cerebras/llama-3.3-70b",
			APIBase:   "https://api.cerebras.ai/v1",
			APIKey:    "",
		},
		{
			ModelName: "doubao-pro",
			Model:     "volcengine/doubao-pro-32k",
			APIBase:   "https://ark.cn-beijing.volces.com/api/v3",
			APIKey:    "",
		},
		{
			ModelName: "deepseek-v3",
			Model:     "shengsuanyun/deepseek-v3",
			APIBase:   "https://api.shengsuanyun.com/v1",
			APIKey:    "",
		},
		{
			ModelName: "gemini-2.5-flash",
			Model:     "gemini/gemini-2.5-flash",
			APIKey:    "",
		},
		{
			ModelName: "gemini-2.5-flash-image",
			Model:     "gemini/gemini-2.5-flash-image",
			APIKey:    "",
		},
		{
			ModelName:  "gemini-3.1-pro-high",
			Model:      "antigravity/gemini-3.1-pro-high",
			APIKey:     "",
			AuthMethod: "oauth",
		},
		{
			ModelName:  "antigravity-gemini-3-flash",
			Model:      "antigravity/gemini-3-flash",
			APIKey:     "",
			AuthMethod: "oauth",
		},
		{
			ModelName:  "gemini-2.5-pro",
			Model:      "antigravity/gemini-2.5-pro",
			APIKey:     "",
			AuthMethod: "oauth",
		},
		{
			ModelName:  "claude-sonnet-4-6",
			Model:      "antigravity/claude-sonnet-4-6",
			APIKey:     "",
			AuthMethod: "oauth",
		},
		{
			ModelName:  "claude-opus-4-6-thinking",
			Model:      "antigravity/claude-opus-4-6-thinking",
			APIKey:     "",
			AuthMethod: "oauth",
		},
		{
			ModelName:  "copilot-gpt-5.2",
			Model:      "github-copilot/gpt-5.2",
			APIBase:    "http://localhost:4321",
			APIKey:     "",
			AuthMethod: "oauth",
		},
		{
			ModelName: "mistral-small",
			Model:     "mistral/mistral-small-latest",
			APIBase:   "https://api.mistral.ai/v1",
			APIKey:    "",
		},
		{
			ModelName: "local-model",
			Model:     "vllm/custom-model",
			APIBase:   "http://localhost:8000/v1",
			APIKey:    "",
		},
		{
			ModelName: "qwen2.5:0.5b",
			Model:     "qwen2.5:0.5b",
			APIBase:   "http://localhost:11434/v1",
			APIKey:    "ollama",
		},
		{
			ModelName: "qwen2.5-coder:0.5b",
			Model:     "qwen2.5-coder:0.5b",
			APIBase:   "http://localhost:11434/v1",
			APIKey:    "ollama",
		},
		{
			ModelName: "llama3.2:1b",
			Model:     "llama3.2:1b",
			APIBase:   "http://localhost:11434/v1",
			APIKey:    "ollama",
		},
	}
	return cfg
}

// QwenDefaultConfig returns a configuration template optimized for Qwen models.
func QwenDefaultConfig(_ bool) *Config {
	cfg := TemplateDefaultConfig()
	cfg.Agents.Defaults.Provider = "qwen"
	cfg.Agents.Defaults.Model = "qwen-plus"
	cfg.Agents.Defaults.MaxTokens = 32768
	for i := range cfg.Agents.List {
		cfg.Agents.List[i].Model.Primary = "qwen-plus"
	}

	cfg.ModelList = []ModelConfig{
		{
			ModelName: "qwen-plus",
			Model:     "qwen/qwen-plus",
			APIBase:   "https://dashscope-us.aliyuncs.com/compatible-mode/v1",
			APIKey:    "YOUR_DASHSCOPE_API_KEY",
		},
	}
	return cfg
}

// OpenRouterDefaultConfig returns a configuration template optimized for OpenRouter.
func OpenRouterDefaultConfig() *Config {
	cfg := TemplateDefaultConfig()
	cfg.Agents.Defaults.Model = "openrouter-auto"
	for i := range cfg.Agents.List {
		cfg.Agents.List[i].Model.Primary = "openrouter-auto"
	}
	cfg.ModelList = []ModelConfig{
		{
			ModelName: "glm-4.5-flash",
			Model:     "zhipu/glm-4.5-flash",
			APIBase:   "https://api.z.ai/api/paas/v4",
			APIKey:    "",
		},
		{
			ModelName: "gpt-5.2",
			Model:     "openai/gpt-5.2",
			APIBase:   "https://api.openai.com/v1",
			APIKey:    "",
		},
		{
			ModelName: "claude-sonnet-4.6",
			Model:     "anthropic/claude-sonnet-4.6",
			APIBase:   "https://api.anthropic.com/v1",
			APIKey:    "",
		},
		{
			ModelName: "deepseek-chat",
			Model:     "deepseek/deepseek-chat",
			APIBase:   "https://api.deepseek.com/v1",
			APIKey:    "",
		},
		{
			ModelName: "gemini-2.0-flash",
			Model:     "gemini/gemini-2.0-flash-exp",
			APIBase:   "https://generativelanguage.googleapis.com/v1beta/openai/",
			APIKey:    "",
		},
		{
			ModelName: "qwen-plus",
			Model:     "qwen/qwen-plus",
			APIBase:   "https://dashscope.aliyuncs.com/compatible-mode/v1",
			APIKey:    "",
		},
		{
			ModelName: "moonshot-v1-8k",
			Model:     "moonshot/moonshot-v1-8k",
			APIBase:   "https://api.moonshot.cn/v1",
			APIKey:    "",
		},
		{
			ModelName: "llama-3.3-70b",
			Model:     "groq/llama-3.3-70b-versatile",
			APIBase:   "https://api.groq.com/openai/v1",
			APIKey:    "",
		},
		{
			ModelName: "openrouter-auto",
			Model:     "openrouter/auto",
			APIBase:   "https://openrouter.ai/api/v1",
			APIKey:    "",
		},
		{
			ModelName: "openrouter-gpt-5.2",
			Model:     "openrouter/openai/gpt-5.2",
			APIBase:   "https://openrouter.ai/api/v1",
			APIKey:    "",
		},
		{
			ModelName: "nemotron-4-340b",
			Model:     "nvidia/nemotron-4-340b-instruct",
			APIBase:   "https://integrate.api.nvidia.com/v1",
			APIKey:    "",
		},
		{
			ModelName: "cerebras-llama-3.3-70b",
			Model:     "cerebras/llama-3.3-70b",
			APIBase:   "https://api.cerebras.ai/v1",
			APIKey:    "",
		},
		{
			ModelName: "doubao-pro",
			Model:     "volcengine/doubao-pro-32k",
			APIBase:   "https://ark.cn-beijing.volces.com/api/v3",
			APIKey:    "",
		},
		{
			ModelName: "deepseek-v3",
			Model:     "shengsuanyun/deepseek-v3",
			APIBase:   "https://api.shengsuanyun.com/v1",
			APIKey:    "",
		},
		{
			ModelName: "gemini-2.5-flash",
			Model:     "gemini/gemini-2.5-flash",
			APIKey:    "",
		},
		{
			ModelName: "gemini-2.5-flash-image",
			Model:     "gemini/gemini-2.5-flash-image",
			APIKey:    "",
		},
		{
			ModelName:  "gemini-3.1-pro-high",
			Model:      "antigravity/gemini-3.1-pro-high",
			APIKey:     "",
			AuthMethod: "oauth",
		},
		{
			ModelName:  "antigravity-gemini-3-flash",
			Model:      "antigravity/gemini-3-flash",
			APIKey:     "",
			AuthMethod: "oauth",
		},
		{
			ModelName:  "gemini-2.5-pro",
			Model:      "antigravity/gemini-2.5-pro",
			APIKey:     "",
			AuthMethod: "oauth",
		},
		{
			ModelName:  "claude-sonnet-4-6",
			Model:      "antigravity/claude-sonnet-4-6",
			APIKey:     "",
			AuthMethod: "oauth",
		},
		{
			ModelName:  "claude-opus-4-6-thinking",
			Model:      "antigravity/claude-opus-4-6-thinking",
			APIKey:     "",
			AuthMethod: "oauth",
		},
		{
			ModelName:  "copilot-gpt-5.2",
			Model:      "github-copilot/gpt-5.2",
			APIBase:    "http://localhost:4321",
			APIKey:     "",
			AuthMethod: "oauth",
		},
		{
			ModelName: "mistral-small",
			Model:     "mistral/mistral-small-latest",
			APIBase:   "https://api.mistral.ai/v1",
			APIKey:    "",
		},
		{
			ModelName: "local-model",
			Model:     "vllm/custom-model",
			APIBase:   "http://localhost:8000/v1",
			APIKey:    "",
		},
		{
			ModelName: "qwen2.5:0.5b",
			Model:     "qwen2.5:0.5b",
			APIBase:   "http://localhost:11434/v1",
			APIKey:    "ollama",
		},
		{
			ModelName: "qwen2.5-coder:0.5b",
			Model:     "qwen2.5-coder:0.5b",
			APIBase:   "http://localhost:11434/v1",
			APIKey:    "ollama",
		},
		{
			ModelName: "llama3.2:1b",
			Model:     "llama3.2:1b",
			APIBase:   "http://localhost:11434/v1",
			APIKey:    "ollama",
		},
	}
	return cfg
}

// OpenRouterFreeDefaultConfig returns a configuration template that uses only
// OpenRouter free-tier models — no API balance required.
// Uses "openrouter/auto" with require_parameters:true so OpenRouter routes only
// to free models that support function/tool calling.
func OpenRouterFreeDefaultConfig() *Config {
	cfg := TemplateDefaultConfig()
	cfg.Agents.Defaults.Model = "openrouter-free"
	for i := range cfg.Agents.List {
		cfg.Agents.List[i].Model = &AgentModelConfig{
			Primary: "openrouter-free",
		}
	}
	cfg.ModelList = []ModelConfig{
		{
			ModelName: "openrouter-free",
			Model:     "openrouter/auto",
			APIBase:   "https://openrouter.ai/api/v1",
			APIKey:    "",
			// require_parameters: true tells OpenRouter to route only to
			// free models that support tool/function calling.
			ExtraBody: map[string]any{
				"provider": map[string]any{
					"require_parameters": true,
				},
			},
		},
	}
	return cfg
}

// GeminiDefaultConfig returns a configuration template optimized for Gemini models.
func GeminiDefaultConfig() *Config {
	cfg := TemplateDefaultConfig()
	cfg.Agents.Defaults.Model = "gemini-2.5-flash"
	cfg.Agents.Defaults.MaxTokens = 8192
	for i := range cfg.Agents.List {
		cfg.Agents.List[i].Model.Primary = "gemini-2.5-flash"
	}
	cfg.ModelList = []ModelConfig{
		{
			ModelName: "glm-4.5-flash",
			Model:     "zhipu/glm-4.5-flash",
			APIBase:   "https://api.z.ai/api/paas/v4",
			APIKey:    "",
		},
		{
			ModelName: "gpt-5.2",
			Model:     "openai/gpt-5.2",
			APIBase:   "https://api.openai.com/v1",
			APIKey:    "",
		},
		{
			ModelName: "claude-sonnet-4.6",
			Model:     "anthropic/claude-sonnet-4.6",
			APIBase:   "https://api.anthropic.com/v1",
			APIKey:    "",
		},
		{
			ModelName: "deepseek-chat",
			Model:     "deepseek/deepseek-chat",
			APIBase:   "https://api.deepseek.com/v1",
			APIKey:    "",
		},
		{
			ModelName: "gemini-2.0-flash",
			Model:     "gemini/gemini-2.0-flash-exp",
			APIBase:   "https://generativelanguage.googleapis.com/v1beta/openai/",
			APIKey:    "",
		},
		{
			ModelName: "qwen-plus",
			Model:     "qwen/qwen-plus",
			APIBase:   "https://dashscope.aliyuncs.com/compatible-mode/v1",
			APIKey:    "",
		},
		{
			ModelName: "moonshot-v1-8k",
			Model:     "moonshot/moonshot-v1-8k",
			APIBase:   "https://api.moonshot.cn/v1",
			APIKey:    "",
		},
		{
			ModelName: "llama-3.3-70b",
			Model:     "groq/llama-3.3-70b-versatile",
			APIBase:   "https://api.groq.com/openai/v1",
			APIKey:    "",
		},
		{
			ModelName: "openrouter-auto",
			Model:     "openrouter/auto",
			APIBase:   "https://openrouter.ai/api/v1",
			APIKey:    "",
		},
		{
			ModelName: "openrouter-gpt-5.2",
			Model:     "openrouter/openai/gpt-5.2",
			APIBase:   "https://openrouter.ai/api/v1",
			APIKey:    "",
		},
		{
			ModelName: "nemotron-4-340b",
			Model:     "nvidia/nemotron-4-340b-instruct",
			APIBase:   "https://integrate.api.nvidia.com/v1",
			APIKey:    "",
		},
		{
			ModelName: "cerebras-llama-3.3-70b",
			Model:     "cerebras/llama-3.3-70b",
			APIBase:   "https://api.cerebras.ai/v1",
			APIKey:    "",
		},
		{
			ModelName: "doubao-pro",
			Model:     "volcengine/doubao-pro-32k",
			APIBase:   "https://ark.cn-beijing.volces.com/api/v3",
			APIKey:    "",
		},
		{
			ModelName: "deepseek-v3",
			Model:     "shengsuanyun/deepseek-v3",
			APIBase:   "https://api.shengsuanyun.com/v1",
			APIKey:    "",
		},
		{
			ModelName: "gemini-2.5-flash",
			Model:     "gemini/gemini-2.5-flash",
			APIKey:    "",
		},
		{
			ModelName: "gemini-2.5-flash-image",
			Model:     "gemini-2.5-flash-image",
			APIKey:    "",
		},
		{
			ModelName:  "gemini-3.1-pro-high",
			Model:      "antigravity/gemini-3.1-pro-high",
			APIKey:     "",
			AuthMethod: "oauth",
		},
		{
			ModelName:  "antigravity-gemini-3-flash",
			Model:      "antigravity/gemini-3-flash",
			APIKey:     "",
			AuthMethod: "oauth",
		},
		{
			ModelName:  "gemini-2.5-pro",
			Model:      "antigravity/gemini-2.5-pro",
			APIKey:     "",
			AuthMethod: "oauth",
		},
		{
			ModelName:  "claude-sonnet-4-6",
			Model:      "antigravity/claude-sonnet-4-6",
			APIKey:     "",
			AuthMethod: "oauth",
		},
		{
			ModelName:  "claude-opus-4-6-thinking",
			Model:      "antigravity/claude-opus-4-6-thinking",
			APIKey:     "",
			AuthMethod: "oauth",
		},
		{
			ModelName:  "copilot-gpt-5.2",
			Model:      "github-copilot/gpt-5.2",
			APIBase:    "http://localhost:4321",
			APIKey:     "",
			AuthMethod: "oauth",
		},
		{
			ModelName: "mistral-small",
			Model:     "mistral/mistral-small-latest",
			APIBase:   "https://api.mistral.ai/v1",
			APIKey:    "",
		},
		{
			ModelName: "local-model",
			Model:     "vllm/custom-model",
			APIBase:   "http://localhost:8000/v1",
			APIKey:    "",
		},
		{
			ModelName: "qwen2.5:0.5b",
			Model:     "qwen2.5:0.5b",
			APIBase:   "http://localhost:11434/v1",
			APIKey:    "ollama",
		},
		{
			ModelName: "qwen2.5-coder:0.5b",
			Model:     "qwen2.5-coder:0.5b",
			APIBase:   "http://localhost:11434/v1",
			APIKey:    "ollama",
		},
		{
			ModelName: "llama3.2:1b",
			Model:     "llama3.2:1b",
			APIBase:   "http://localhost:11434/v1",
			APIKey:    "ollama",
		},
	}
	return cfg
}
