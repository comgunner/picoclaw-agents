// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

// ============================================================================
// ⚠️  CRITICAL: Token Estimation Logic — DO NOT MODIFY WITHOUT REVIEW
// ============================================================================
//
// This file contains the proactive token validation logic that prevents the
// OpenRouter Free tier 402 error ("Prompt tokens limit exceeded").
//
// THE BUG (2026-04-05):
//   estimateTokens() used to count only Content characters and add a fixed
//   "+2500" overhead for tool definitions. In reality, 60+ tool definitions
//   consume ~15,000 tokens. This caused 21,526 tokens to be sent to models
//   with a ~7,869 token limit, resulting in constant 402 errors.
//
// THE FIX:
//   1. estimateTokens() calls tokenizer.EstimateMessageTokens() for each message
//      (counts Content, ReasoningContent, ToolCalls, ToolCallID, SystemParts)
//   2. Proactive check calls tokenizer.EstimateToolDefsTokens(providerToolDefs)
//      to get REAL tool token counts instead of a fixed "+2500"
//   3. Auto-switches to essential tools if tool tokens > 30% of context window
//   4. Progressive truncation in loop: keep 5→3→2→1 messages, re-estimate each step
//   5. Emergency fallback to PromptLevelMinimal if still over budget
//
// IMMUTABLE RULES:
//   - NEVER use a fixed overhead value (like +2500, +3000) for tool definitions
//   - ALWAYS call tokenizer.EstimateToolDefsTokens() for actual tool token counts
//   - ALWAYS call tokenizer.EstimateMessageTokens() for message token counts
//   - See local_work/MEMORY.md and local_work/openrouter_free_token_fix.md
//
// ============================================================================

package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	managementinit "github.com/comgunner/picoclaw/pkg/agents/management/init"
	"github.com/comgunner/picoclaw/pkg/auth"
	"github.com/comgunner/picoclaw/pkg/bus"
	"github.com/comgunner/picoclaw/pkg/channels"
	"github.com/comgunner/picoclaw/pkg/config"
	"github.com/comgunner/picoclaw/pkg/constants"
	pcontext "github.com/comgunner/picoclaw/pkg/context"
	"github.com/comgunner/picoclaw/pkg/logger"
	"github.com/comgunner/picoclaw/pkg/providers"
	"github.com/comgunner/picoclaw/pkg/routing"
	"github.com/comgunner/picoclaw/pkg/security"
	"github.com/comgunner/picoclaw/pkg/skills"
	"github.com/comgunner/picoclaw/pkg/state"
	"github.com/comgunner/picoclaw/pkg/tasklock"
	"github.com/comgunner/picoclaw/pkg/tokenizer"
	"github.com/comgunner/picoclaw/pkg/tools"
	"github.com/comgunner/picoclaw/pkg/utils"
)

type AgentLoop struct {
	bus            *bus.MessageBus
	cfg            *config.Config
	registry       *AgentRegistry
	state          *state.Manager
	running        atomic.Bool
	summarizing    sync.Map
	fallback       *providers.FallbackChain
	channelManager *channels.Manager
	summaryCache   *utils.SummaryCache
	tasklocks      *tasklock.Manager
	// validator is used to pre-check token usage before calling an LLM.
	tokenValidator providers.TokenValidator
	// sentinel is used for local security checks against prompt injection.
	sentinel *tools.SkillsSentinelTool
	// auditor handles security event logging.
	auditor *security.Auditor
	// configMutex protects concurrent access to config.json during runtime updates.
	configMutex sync.Mutex

	// providerCache stores instantiated providers to avoid repeated creation.
	// Keyed by provider hash (api_base + api_key).
	providerCache map[string]providers.LLMProvider

	contextMiddleware *pcontext.ContextMiddleware
	runtimeMgr        *RuntimeManager

	// contextManager provides pluggable, budget-aware context management.
	// One instance per AgentLoop, selected by config (default: "legacy", optional: "seahorse").
	contextManager ContextManager
}

// GetProvider returns a cached provider for the given model name or creates a new one.
func (al *AgentLoop) GetProvider(modelName string) (providers.LLMProvider, string, error) {
	al.configMutex.Lock()
	defer al.configMutex.Unlock()

	// 1. Get model config to check api_base/api_key
	modelCfg, err := al.cfg.GetModelConfig(modelName)
	if err != nil {
		// Try dynamic provider if not in list
		p, mid, errDyn := providers.CreateProviderForModel(al.cfg, modelName)
		if errDyn != nil {
			// Fall back to the default provider passed to NewAgentLoop
			if al.providerCache["__default__"] != nil {
				return al.providerCache["__default__"], modelName, nil
			}
		}
		return p, mid, errDyn
	}

	// 2. Generate cache key
	cacheKey := fmt.Sprintf("%s:%s", modelCfg.APIBase, modelCfg.APIKey)
	if p, ok := al.providerCache[cacheKey]; ok {
		return p, modelCfg.Model, nil
	}

	// 3. Create new provider
	p, mid, err := providers.CreateProviderFromConfig(modelCfg)
	if err != nil {
		// Fall back to the default provider passed to NewAgentLoop
		if al.providerCache["__default__"] != nil {
			return al.providerCache["__default__"], modelName, nil
		}
		return nil, "", err
	}

	// 4. Cache and return
	if al.providerCache == nil {
		al.providerCache = make(map[string]providers.LLMProvider)
	}
	al.providerCache[cacheKey] = p
	return p, mid, nil
}

// processOptions configures how a message is processed
type processOptions struct {
	SessionKey      string            // Session identifier for history/context
	Channel         string            // Target channel for tool execution
	ChatID          string            // Target chat ID for tool execution
	UserMessage     string            // User message content (may include prefix)
	DefaultResponse string            // Response when LLM returns empty
	EnableSummary   bool              // Whether to trigger summarization
	SendResponse    bool              // Whether to send response via bus
	NoHistory       bool              // If true, don't load session history (for heartbeat)
	Metadata        map[string]string // Message metadata (e.g., model_name from client)
}

func NewAgentLoop(cfg *config.Config, msgBus *bus.MessageBus, provider providers.LLMProvider) *AgentLoop {
	al := &AgentLoop{
		bus:           msgBus,
		cfg:           cfg,
		providerCache: make(map[string]providers.LLMProvider),
	}

	// Store the default provider as fallback for tests and misconfigured environments
	if provider != nil {
		al.providerCache["__default__"] = provider
	}

	// Factory for agent registration
	factory := func(model string) (providers.LLMProvider, string, error) {
		return al.GetProvider(model)
	}

	registry := NewAgentRegistry(cfg, factory)
	al.registry = registry

	// Register shared tools to all agents
	suite := registerSharedTools(cfg, msgBus, registry, provider)

	// Set up shared fallback chain
	cooldown := providers.NewCooldownTracker()
	fallbackChain := providers.NewFallbackChain(cooldown)

	// Create state manager using default agent's workspace for channel recording
	defaultAgent := registry.GetDefaultAgent()
	var stateManager *state.Manager
	if defaultAgent != nil {
		stateManager = state.NewManager(defaultAgent.Workspace)
	}

	validator := providers.NewTokenValidator(utils.NewBasicTokenCounter(), constants.DefaultMaxContextTokens)

	workspace := "."
	if defaultAgent != nil {
		workspace = defaultAgent.Workspace
	}
	cache := utils.NewSummaryCache(workspace)
	tlm := tasklock.NewManager(workspace)

	// Start Watchdog to clean up dead locks (5 min timeout, 1 min check interval)
	tlm.StartWatchdog(5*time.Minute, 1*time.Minute)

	// Rehydrate stalled tasks (Crash Recovery)
	go func() {
		// Allow system to boot up fully before sending recovery messages
		time.Sleep(3 * time.Second)
		active := tlm.GetAllActiveLocks()
		for _, tl := range active {
			if tl.GetStatus() == tasklock.StatusInProgress || tl.GetStatus() == tasklock.StatusNetworkRetry {
				sessionKey := tl.GetSessionKey()

				// Heartbeat tasks are internal and fire-and-forget.
				// If a heartbeat lock is stranded after a crash, silently discard it — never
				// wake up the LLM for an internal monitoring task.
				if sessionKey == "heartbeat" || strings.HasPrefix(tl.TaskID, "heartbeat_") {
					logger.InfoCF(
						"rehydration",
						"Discarding stranded heartbeat lock (internal task, no recovery needed)",
						map[string]any{"task_id": tl.TaskID},
					)
					tlm.RemoveLock(tl.TaskID)
					continue
				}

				tl.UpdateState(tasklock.StatusRecovering, "reboot_recovery", nil)
				logger.InfoCF("rehydration", "Recovering stranded task lock", map[string]any{"task_id": tl.TaskID})
				msgBus.PublishInbound(bus.InboundMessage{
					Channel:    "cli", // Fallback to cli since channel is implicit by session
					ChatID:     "console",
					SessionKey: sessionKey,
					Content:    "[System] A critical crash occurred, but your session has been rehydrated from disk. CRITICAL INSTRUCTION: Do NOT use any internal diagnostic tools (like system_diagnostics, read_file, exec, context_status, queue, etc) right now to investigate the system log. Simply send a short, direct message to the user confirming you are back online, and wait for their input before doing any additional background verification.",
				})
			}
		}
	}()

	sentinel := tools.NewSkillsSentinelTool()
	sentinel.SetWorkspace(workspace)

	cm := pcontext.NewContextMiddleware(workspace, 8192)
	cm.Start()

	al.state = stateManager
	al.registry = registry
	al.summarizing = sync.Map{}
	al.fallback = fallbackChain
	al.tokenValidator = validator
	al.summaryCache = cache
	al.tasklocks = tlm
	al.sentinel = sentinel
	al.auditor = security.NewAuditor(workspace)
	al.contextMiddleware = cm

	al.RegisterTool(tools.NewContextStatusTool(cm))

	// Initialize pluggable ContextManager (default: "legacy", optional: "seahorse")
	// Selection happens per-agent via config: agents.defaults.context_manager
	// Legacy is always available. Seahorse is registered by init() in context_seahorse.go.
	cmName := cfg.Agents.Defaults.ContextManager
	if cmName == "" {
		cmName = "legacy"
	}

	if factory, ok := lookupContextManager(cmName); ok {
		var cmCfg json.RawMessage
		if cfg.Agents.Defaults.ContextManagerConfig != nil {
			cmCfg, _ = json.Marshal(cfg.Agents.Defaults.ContextManagerConfig)
		}
		cm, err := factory(cmCfg, al)
		if err != nil {
			logger.WarnCF("agent", "Failed to create configured ContextManager, falling back to legacy",
				map[string]any{"requested": cmName, "error": err.Error()})
			if legacyFactory, ok2 := lookupContextManager("legacy"); ok2 {
				cm, _ = legacyFactory(nil, al)
			}
		}
		if cm != nil {
			al.contextManager = cm
			logger.InfoCF("agent", "ContextManager initialized", map[string]any{
				"name": cmName,
			})
		}
	} else {
		logger.WarnCF("agent", "Requested ContextManager not found, using legacy",
			map[string]any{"requested": cmName})
		if legacyFactory, ok2 := lookupContextManager("legacy"); ok2 {
			al.contextManager, _ = legacyFactory(nil, al)
		}
	}

	// Initialize the autonomous runtime manager
	al.runtimeMgr = NewRuntimeManager(al, suite.MessageBus, suite.Instances)

	return al
}

// InitSentinelCallbacks initializes sentinel callbacks for auto-reactivation notifications
func (al *AgentLoop) InitSentinelCallbacks() {
	al.sentinel.SetAutoReactivateCallback(func() {
		// Send notification via message bus
		al.bus.PublishOutbound(bus.OutboundMessage{
			Channel: "system",
			ChatID:  "broadcast",
			Content: "🛡️ **SENTINEL AUTO-ACTIVATED** | Security checks restored after timeout",
		})

		// Audit log
		al.auditor.LogSecurityEvent("system", "", "sentinel_auto_activated",
			"Sentinel automatically reactivated after timeout",
			"Auto-reactivation notification sent")
	})
}

// registerSharedTools registers tools that are shared across all agents (web, message, spawn).
func registerSharedTools(
	cfg *config.Config,
	msgBus *bus.MessageBus,
	registry *AgentRegistry,
	provider providers.LLMProvider,
) *managementinit.Suite {
	// 1. Initialize Global Shared Tools/Trackers
	// Resolve Gemini configuration for image tools
	geminiTextAPIKey := cfg.Tools.ImageGen.GeminiAPIKey
	geminiTextModel := ""
	if cfg.Tools.ImageGen.GeminiTextModelName != "" {
		if mc, err := cfg.GetModelConfig(cfg.Tools.ImageGen.GeminiTextModelName); err == nil {
			if mc.APIKey != "" {
				geminiTextAPIKey = mc.APIKey
			}
			geminiTextModel = mc.Model
		}
	}

	geminiImageAPIKey := cfg.Tools.ImageGen.GeminiAPIKey
	geminiImageModel := ""
	if cfg.Tools.ImageGen.GeminiImageModelName != "" {
		if mc, err := cfg.GetModelConfig(cfg.Tools.ImageGen.GeminiImageModelName); err == nil {
			if mc.APIKey != "" {
				geminiImageAPIKey = mc.APIKey
			}
			geminiImageModel = mc.Model
		}
	}

	// Resolve Ideogram configuration
	ideogramAPIKey := cfg.Tools.ImageGen.IdeogramAPIKey
	ideogramAPIURL := cfg.Tools.ImageGen.IdeogramAPIURL
	if cfg.Tools.ImageGen.IdeogramModelName != "" {
		if mc, err := cfg.GetModelConfig(cfg.Tools.ImageGen.IdeogramModelName); err == nil {
			if mc.APIKey != "" {
				ideogramAPIKey = mc.APIKey
			}
			if mc.APIBase != "" {
				ideogramAPIURL = mc.APIBase
			}
		}
	}

	// Create shared ImageGen instances
	// Default workspace used for tool initialization
	workspace := "."
	defaultAgent := registry.GetDefaultAgent()
	if defaultAgent != nil {
		workspace = defaultAgent.Workspace
	}

	// FIX #7: Create ONE shared cooldown instance for ALL agents (outside the loop).
	imageCooldown, _ := utils.NewImageCooldown(workspace)
	cooldownSecs := cfg.Tools.ImageGen.CooldownSeconds
	if cooldownSecs <= 0 {
		cooldownSecs = 150 // Default 2.5 minutes
	}

	// Determine image generation provider with fallback chain.
	imageProvider := strings.TrimSpace(cfg.Tools.ImageGen.Provider)
	if imageProvider == "" {
		imageProvider = os.Getenv("PICOCLAW_TOOLS_IMAGE_GEN_PROVIDER")
	}
	if imageProvider == "" {
		imageProvider = "antigravity" // NEW DEFAULT
	}

	// Check if Antigravity OAuth credentials exist.
	hasAntigravityAuth := false
	if imageProvider == "antigravity" {
		cred, err := auth.GetCredential("google-antigravity")
		if err == nil && cred != nil && !cred.IsExpired() {
			hasAntigravityAuth = true
		}
	}

	// FIX #9: If no OAuth, fallback to gemini for image_gen_create.
	if imageProvider == "antigravity" && !hasAntigravityAuth {
		imageProvider = "gemini"
	}

	imageGenTracker := (*utils.ImageGenTracker)(nil)

	if hasAntigravityAuth {
		// PRIMARY: Antigravity OAuth-based image generation with cooldown.
		antigravityModel := strings.TrimSpace(cfg.Tools.ImageGen.AntigravityModel)
		if antigravityModel == "" {
			antigravityModel = "gemini-3.1-flash-image"
		}
		antigravityTool := tools.NewImageGenAntigravityToolFromConfig(
			antigravityModel,
			cfg.Tools.ImageGen.AspectRatio,
			cfg.Tools.ImageGen.OutputDir,
			workspace,
			cooldownSecs,
			imageCooldown, // Same shared instance for all agents.
		)
		imageGenTracker = antigravityTool.GetTracker()
	}

	// FALLBACK: Always create Gemini/Ideogram API key tools as fallback.
	imageGenTool := tools.NewImageGenCreateToolFromConfig(
		imageProvider,
		geminiImageAPIKey,
		geminiImageModel,
		geminiTextModel,
		ideogramAPIKey,
		ideogramAPIURL,
		cfg.Tools.ImageGen.AspectRatio,
		cfg.Tools.ImageGen.OutputDir,
		cfg.Tools.ImageGen.ImageScriptPath,
		cfg.Tools.ImageGen.ImageGenScriptPath,
		workspace,
	)
	// Use Antigravity tracker if available, otherwise Gemini tracker.
	if imageGenTracker == nil {
		imageGenTracker = imageGenTool.GetTracker()
	}

	queueMgr := tools.GetQueueManager(msgBus, nil) // Will be attached per agent tools later

	// Agent Management Skill: create shared suite once for all agents.
	// The suite owns AgentRegistry, InstanceRegistry, MessageBus and RateLimiter.
	managementSuite := managementinit.NewSuite(cfg)

	for _, agentID := range registry.ListAgentIDs() {
		agent, ok := registry.GetAgent(agentID)
		if !ok {
			continue
		}

		// Subagent tools prep
		subagentManager := tools.NewSubagentManager(provider, agent.Model, agent.Workspace, msgBus, cfg)
		subagentManager.SetLLMOptions(agent.MaxTokens, agent.Temperature)
		subagentManager.SetTools(agent.Tools) // SHARING REGISTRY!

		// Web tools
		if searchTool := tools.NewWebSearchTool(tools.WebSearchToolOptions{
			BraveAPIKey:          cfg.Tools.Web.Brave.APIKey,
			BraveMaxResults:      cfg.Tools.Web.Brave.MaxResults,
			BraveEnabled:         cfg.Tools.Web.Brave.Enabled,
			TavilyAPIKey:         cfg.Tools.Web.Tavily.APIKey,
			TavilyBaseURL:        cfg.Tools.Web.Tavily.BaseURL,
			TavilyMaxResults:     cfg.Tools.Web.Tavily.MaxResults,
			TavilyEnabled:        cfg.Tools.Web.Tavily.Enabled,
			DuckDuckGoMaxResults: cfg.Tools.Web.DuckDuckGo.MaxResults,
			DuckDuckGoEnabled:    cfg.Tools.Web.DuckDuckGo.Enabled,
			PerplexityAPIKey:     cfg.Tools.Web.Perplexity.APIKey,
			PerplexityMaxResults: cfg.Tools.Web.Perplexity.MaxResults,
			PerplexityEnabled:    cfg.Tools.Web.Perplexity.Enabled,
			Proxy:                cfg.Tools.Web.Proxy,
		}); searchTool != nil {
			agent.Tools.Register(searchTool)
		}
		agent.Tools.Register(tools.NewWebFetchToolWithProxy(50000, cfg.Tools.Web.Proxy))
		agent.Tools.Register(tools.NewBinanceTickerPriceToolFromConfig(
			cfg.Tools.Binance.APIKey,
			cfg.Tools.Binance.SecretKey,
		))
		agent.Tools.Register(tools.NewBinanceOrderBookTool())
		agent.Tools.Register(tools.NewBinanceFuturesOrderBookTool())
		agent.Tools.Register(tools.NewBinanceListFuturesVolumeTool())
		agent.Tools.Register(tools.NewBinanceSpotBalanceToolFromConfig(
			cfg.Tools.Binance.APIKey,
			cfg.Tools.Binance.SecretKey,
		))
		agent.Tools.Register(tools.NewBinanceFuturesBalanceToolFromConfig(
			cfg.Tools.Binance.APIKey,
			cfg.Tools.Binance.SecretKey,
		))
		agent.Tools.Register(tools.NewBinanceFuturesOpenPositionToolFromConfig(
			cfg.Tools.Binance.APIKey,
			cfg.Tools.Binance.SecretKey,
		))
		agent.Tools.Register(tools.NewBinanceFuturesClosePositionToolFromConfig(
			cfg.Tools.Binance.APIKey,
			cfg.Tools.Binance.SecretKey,
		))

		// Base Trader tool (comprehensive trading analysis with 6-layer validation)
		agent.Tools.Register(tools.NewBaseTraderTool(10000)) // Default balance: $10,000

		// Social Media tools
		agent.Tools.Register(tools.NewFacebookPostToolFromConfig(
			cfg.Tools.SocialMedia.Facebook.DefaultPageID,
			cfg.Tools.SocialMedia.Facebook.DefaultPageToken,
			cfg.Tools.SocialMedia.Facebook.AppID,
			cfg.Tools.SocialMedia.Facebook.AppSecret,
			cfg.Tools.SocialMedia.Facebook.UserToken,
		))
		agent.Tools.Register(tools.NewXPostTweetToolFromConfig(
			cfg.Tools.SocialMedia.X.APIKey,
			cfg.Tools.SocialMedia.X.APISecret,
			cfg.Tools.SocialMedia.X.AccessToken,
			cfg.Tools.SocialMedia.X.AccessTokenSecret,
		))
		agent.Tools.Register(tools.NewDiscordWebhookToolFromConfig(
			cfg.Tools.SocialMedia.Discord.WebhookURL,
		))

		// Notion tools
		agent.Tools.Register(tools.NewNotionCreatePageToolFromConfig(cfg.Tools.Notion.APIKey))
		agent.Tools.Register(tools.NewNotionQueryDatabaseToolFromConfig(cfg.Tools.Notion.APIKey))
		agent.Tools.Register(tools.NewNotionSearchToolFromConfig(cfg.Tools.Notion.APIKey))
		agent.Tools.Register(tools.NewNotionUpdatePageToolFromConfig(cfg.Tools.Notion.APIKey))

		// Image Generation tools
		agent.Tools.Register(tools.NewTextScriptCreateToolFromConfig(
			geminiTextAPIKey,
			cfg.Tools.ImageGen.OutputDir,
			cfg.Tools.ImageGen.ImageScriptPath,
			agent.Workspace,
			geminiTextModel,
			cfg.Tools.ImageGen.AspectRatio,
			cfg.Tools.ImageGen.Provider,
		))

		// Register Antigravity tool if OAuth is available.
		var antigravityTool *tools.ImageGenAntigravityTool
		if hasAntigravityAuth {
			antigravityModel := strings.TrimSpace(cfg.Tools.ImageGen.AntigravityModel)
			if antigravityModel == "" {
				antigravityModel = "gemini-3.1-flash-image"
			}
			antigravityTool = tools.NewImageGenAntigravityToolFromConfig(
				antigravityModel,
				cfg.Tools.ImageGen.AspectRatio,
				cfg.Tools.ImageGen.OutputDir,
				workspace,
				cooldownSecs,
				imageCooldown, // Shared instance for all agents.
			)
			agent.Tools.Register(antigravityTool)
		} else {
			// No OAuth — register the API key based tool as fallback
			agent.Tools.Register(imageGenTool)
		}
		agent.Tools.Register(tools.NewImageGenWorkflowToolWithTracker(imageGenTracker))
		// Queue and Batch Tools (attach agent tools to queue)
		queueMgr.SetTools(agent.Tools)
		agent.Tools.Register(tools.NewQueueTool(queueMgr))
		agent.Tools.Register(tools.NewBatchIDTool())

		agent.Tools.Register(tools.NewSocialPostBundleTool(
			geminiImageAPIKey,
			geminiTextModel,
			geminiImageModel,
			ideogramAPIKey,
			ideogramAPIURL,
			cfg.Tools.ImageGen.AspectRatio,
			cfg.Tools.ImageGen.OutputDir,
			cfg.Tools.ImageGen.ImageScriptPath,
			cfg.Tools.ImageGen.ImageGenScriptPath,
			agent.Workspace,
			queueMgr,
			msgBus,
			imageGenTracker,
			antigravityTool, // ← Antigravity OAuth for images
		))

		// Social Media tools
		socialManager := tools.NewSocialManagerTool(imageGenTracker)
		socialManager.SetRegistry(agent.Tools) // Set registry to access facebook_post tool
		agent.Tools.Register(socialManager)

		// Community Manager tools
		agent.Tools.Register(tools.NewCommunityManagerToolWithWorkspace(agent.Workspace))
		agent.Tools.Register(tools.NewCommunityFromImageTool())

		// Approval tools (human-in-the-loop)
		agent.Tools.Register(tools.NewTextApprovalToolWithTracker(imageGenTracker))
		agent.Tools.Register(tools.NewImageApprovalToolWithTracker(imageGenTracker))

		// Hardware tools (I2C, SPI) - Linux only, returns error on other platforms
		agent.Tools.Register(tools.NewI2CTool())
		agent.Tools.Register(tools.NewSPITool())

		// Message tool
		messageTool := tools.NewMessageTool()
		messageTool.SetSendCallback(func(channel, chatID, content string, media []string, buttons []bus.Button) error {
			msgBus.PublishOutbound(bus.OutboundMessage{
				Channel: channel,
				ChatID:  chatID,
				Content: content,
				Media:   media,
				Buttons: buttons,
			})
			return nil
		})
		agent.Tools.Register(messageTool)

		// Skill discovery and installation tools
		registryMgr := skills.NewRegistryManagerFromConfig(skills.RegistryConfig{
			MaxConcurrentSearches: cfg.Tools.Skills.MaxConcurrentSearches,
			ClawHub:               skills.ClawHubConfig(cfg.Tools.Skills.Registries.ClawHub),
		})
		searchCache := skills.NewSearchCache(
			cfg.Tools.Skills.SearchCache.MaxSize,
			time.Duration(cfg.Tools.Skills.SearchCache.TTLSeconds)*time.Second,
		)
		agent.Tools.Register(tools.NewFindSkillsTool(registryMgr, searchCache))
		agent.Tools.Register(tools.NewInstallSkillTool(registryMgr, agent.Workspace))

		// Spawn tool with allowlist checker
		spawnTool := tools.NewSpawnTool(subagentManager)
		currentAgentID := agentID
		spawnTool.SetAllowlistChecker(func(targetAgentID string) bool {
			return registry.CanSpawnSubagent(currentAgentID, targetAgentID)
		})
		agent.Tools.Register(spawnTool)
		agent.Tools.Register(tools.NewSubagentTool(subagentManager))
		agent.Tools.Register(tools.NewSubagentListTool(subagentManager))

		// Agent Management Skill — registers all 12 management tools so the LLM
		// can inspect, monitor and communicate with other agents at runtime.
		managementSuite.RegisterAllTools(agent.Tools)
	}

	return managementSuite
}

func (al *AgentLoop) Run(ctx context.Context) error {
	al.running.Store(true)

	// Start all autonomous agent runtimes
	if al.runtimeMgr != nil {
		al.runtimeMgr.StartAll(ctx)
	}

	for al.running.Load() {
		select {
		case <-ctx.Done():
			return nil
		default:
			msg, ok := al.bus.ConsumeInbound(ctx)
			if !ok {
				logger.WarnC("agent", "Message bus closed, stopping agent loop")
				return nil
			}

			// Reset tool states for all agents before processing a new message.
			// This ensures tools like the message tool can track "sent in round" correctly
			// for each individual inbound message.
			al.registry.ResetAllToolStates()

			response, err := al.processMessage(ctx, msg)
			if err != nil {
				response = fmt.Sprintf("Error processing message: %v", err)
			}

			if response != "" {
				// Check if the message tool already sent a response during this round FOR THIS AGENT
				// If so, skip publishing to avoid duplicate messages to the user.
				alreadySent := false

				// Route to determine which agent originally handled the message
				route := al.registry.ResolveRoute(routing.RouteInput{
					Channel:    msg.Channel,
					AccountID:  msg.Metadata["account_id"],
					Peer:       extractPeer(msg),
					ParentPeer: extractParentPeer(msg),
					GuildID:    msg.Metadata["guild_id"],
					TeamID:     msg.Metadata["team_id"],
				})

				agent, ok := al.registry.GetAgent(route.AgentID)
				if !ok {
					agent = al.registry.GetDefaultAgent()
				}

				if agent != nil {
					if tool, ok := agent.Tools.Get("message"); ok {
						if mt, ok := tool.(*tools.MessageTool); ok {
							alreadySent = mt.HasSentInRound()
						}
					}
				}

				if !alreadySent {
					al.bus.PublishOutbound(bus.OutboundMessage{
						Channel: msg.Channel,
						ChatID:  msg.ChatID,
						Content: response,
					})
				}
			}
		}
	}

	return nil
}

func (al *AgentLoop) Stop() {
	if al.contextMiddleware != nil {
		al.contextMiddleware.Stop()
	}
	al.running.Store(false)
}

func (al *AgentLoop) RegisterTool(tool tools.Tool) {
	for _, agentID := range al.registry.ListAgentIDs() {
		if agent, ok := al.registry.GetAgent(agentID); ok {
			agent.Tools.Register(tool)
		}
	}
}

func (al *AgentLoop) SetChannelManager(cm *channels.Manager) {
	al.channelManager = cm
}

// RecordLastChannel records the last active channel for this workspace.
// This uses the atomic state save mechanism to prevent data loss on crash.
func (al *AgentLoop) RecordLastChannel(channel string) error {
	if al.state == nil {
		return nil
	}
	return al.state.SetLastChannel(channel)
}

// RecordLastChatID records the last active chat ID for this workspace.
// This uses the atomic state save mechanism to prevent data loss on crash.
func (al *AgentLoop) RecordLastChatID(chatID string) error {
	if al.state == nil {
		return nil
	}
	return al.state.SetLastChatID(chatID)
}

func (al *AgentLoop) ProcessDirect(ctx context.Context, content, sessionKey string) (string, error) {
	return al.ProcessDirectWithChannel(ctx, content, sessionKey, "cli", "direct")
}

func (al *AgentLoop) ProcessDirectWithChannel(
	ctx context.Context,
	content, sessionKey, channel, chatID string,
) (string, error) {
	msg := bus.InboundMessage{
		Channel:    channel,
		SenderID:   "cron",
		ChatID:     chatID,
		Content:    content,
		SessionKey: sessionKey,
	}

	return al.processMessage(ctx, msg)
}

// ProcessHeartbeat processes a heartbeat request without session history.
// Each heartbeat is independent and doesn't accumulate context.
func (al *AgentLoop) ProcessHeartbeat(ctx context.Context, content, channel, chatID string) (string, error) {
	agent := al.registry.GetDefaultAgent()
	return al.runAgentLoop(ctx, agent, processOptions{
		SessionKey:      "heartbeat",
		Channel:         channel,
		ChatID:          chatID,
		UserMessage:     content,
		DefaultResponse: "I've completed processing but have no response to give.",
		EnableSummary:   false,
		SendResponse:    false,
		NoHistory:       true, // Don't load session history for heartbeat
	})
}

func (al *AgentLoop) processMessage(ctx context.Context, msg bus.InboundMessage) (string, error) {
	// Add message preview to log (show full content for error messages)
	var logContent string
	if strings.Contains(msg.Content, "Error:") || strings.Contains(msg.Content, "error") {
		logContent = msg.Content // Full content for errors
	} else {
		logContent = utils.Truncate(msg.Content, 80)
	}
	logger.InfoCF("agent", fmt.Sprintf("Processing message from %s:%s: %s", msg.Channel, msg.SenderID, logContent),
		map[string]any{
			"channel":     msg.Channel,
			"chat_id":     msg.ChatID,
			"sender_id":   msg.SenderID,
			"session_key": msg.SessionKey,
		})

	// Route system messages to processSystemMessage
	if msg.Channel == "system" {
		return al.processSystemMessage(ctx, msg)
	}

	// Check for commands
	if response, handled := al.handleCommand(ctx, msg); handled {
		return response, nil
	}

	// Route to determine agent and session key
	route := al.registry.ResolveRoute(routing.RouteInput{
		Channel:    msg.Channel,
		AccountID:  msg.Metadata["account_id"],
		Peer:       extractPeer(msg),
		ParentPeer: extractParentPeer(msg),
		GuildID:    msg.Metadata["guild_id"],
		TeamID:     msg.Metadata["team_id"],
	})

	agent, ok := al.registry.GetAgent(route.AgentID)
	if !ok {
		agent = al.registry.GetDefaultAgent()
	}

	// Use routed session key, but honor pre-set agent-scoped keys (for ProcessDirect/cron)
	sessionKey := route.SessionKey
	if msg.SessionKey != "" && strings.HasPrefix(msg.SessionKey, "agent:") {
		sessionKey = msg.SessionKey
	}

	logger.InfoCF("agent", "Routed message",
		map[string]any{
			"agent_id":    agent.ID,
			"session_key": sessionKey,
			"matched_by":  route.MatchedBy,
		})

	// Fast path: explicit tool invocation from user.
	// Example: "Usa text_script_create con topic='IA', category='noticia'"
	if response, handled := al.handleExplicitToolInvoke(ctx, agent, msg.Channel, msg.ChatID, msg.Content); handled {
		return response, nil
	}

	return al.runAgentLoop(ctx, agent, processOptions{
		SessionKey:      sessionKey,
		Channel:         msg.Channel,
		ChatID:          msg.ChatID,
		UserMessage:     msg.Content,
		DefaultResponse: "I've completed processing but have no response to give.",
		EnableSummary:   true,
		SendResponse:    false,
		Metadata:        msg.Metadata,
	})
}

func (al *AgentLoop) processSystemMessage(ctx context.Context, msg bus.InboundMessage) (string, error) {
	if msg.Channel != "system" {
		return "", fmt.Errorf("processSystemMessage called with non-system message channel: %s", msg.Channel)
	}

	logger.InfoCF("agent", "Processing system message",
		map[string]any{
			"sender_id": msg.SenderID,
			"chat_id":   msg.ChatID,
		})

	// Parse origin channel from chat_id (format: "channel:chat_id")
	var originChannel, originChatID string
	if idx := strings.Index(msg.ChatID, ":"); idx > 0 {
		originChannel = msg.ChatID[:idx]
		originChatID = msg.ChatID[idx+1:]
	} else {
		originChannel = "cli"
		originChatID = msg.ChatID
	}

	// Extract subagent result from message content
	// Format: "Task 'label' completed.\n\nResult:\n<actual content>"
	content := msg.Content
	if idx := strings.Index(content, "Result:\n"); idx >= 0 {
		content = content[idx+8:] // Extract just the result part
	}

	// Skip internal channels - only log, don't send to user
	if constants.IsInternalChannel(originChannel) {
		logger.InfoCF("agent", "Subagent completed (internal channel)",
			map[string]any{
				"sender_id":   msg.SenderID,
				"content_len": len(content),
				"channel":     originChannel,
			})
		return "", nil
	}

	// Use default agent for system messages
	agent := al.registry.GetDefaultAgent()

	// Use the origin session for context
	sessionKey := routing.BuildAgentMainSessionKey(agent.ID)

	// Optimization: If the system message is just a "Task completed" notification,
	// we avoid starting a full LLM interaction to prevent auto-responder loops.
	// We'll just proxy the result message.
	isSimpleCompletion := strings.Contains(msg.Content, "completed") && len(content) < 1000

	return al.runAgentLoop(ctx, agent, processOptions{
		SessionKey:      sessionKey,
		Channel:         originChannel,
		ChatID:          originChatID,
		UserMessage:     fmt.Sprintf("[System: %s] %s", msg.SenderID, msg.Content),
		DefaultResponse: "Background task completed.",
		EnableSummary:   false,
		SendResponse:    !isSimpleCompletion, // Only send if it's substantial or complex
	})
}

const HardMaxIterations = 18

// runAgentLoop is the core message processing logic.
func (al *AgentLoop) runAgentLoop(ctx context.Context, agent *AgentInstance, opts processOptions) (string, error) {
	originalUserMessage := opts.UserMessage
	userMessageForLLM := opts.UserMessage
	if !strings.HasPrefix(strings.TrimSpace(userMessageForLLM), "[System:") {
		userMessageForLLM = enrichBinanceShortcutPrompt(userMessageForLLM)
	}

	// 0. Record last channel for heartbeat notifications (skip internal channels)
	if opts.Channel != "" && opts.ChatID != "" {
		// Don't record internal channels (cli, system, subagent)
		if !constants.IsInternalChannel(opts.Channel) {
			channelKey := fmt.Sprintf("%s:%s", opts.Channel, opts.ChatID)
			if err := al.RecordLastChannel(channelKey); err != nil {
				logger.WarnCF("agent", "Failed to record last channel", map[string]any{"error": err.Error()})
			}
		}
	}

	// 1. Internal Security Check (Sentinel)
	if al.sentinel != nil && originalUserMessage != "" {
		// Only check if it's a real user message (not system-initiated or empty)
		if !strings.HasPrefix(originalUserMessage, "[System:") {
			result := al.sentinel.Execute(ctx, map[string]any{"input": originalUserMessage})
			if result.IsError {
				logger.WarnCF("security", "Sentinel blocked malicious input", map[string]any{
					"agent_id":    agent.ID,
					"session_key": opts.SessionKey,
					"pattern":     result.ForLLM,
				})

				if al.auditor != nil {
					al.auditor.LogSecurityEvent(
						agent.ID,
						opts.SessionKey,
						"prompt_injection",
						originalUserMessage,
						result.ForLLM,
					)
				}

				return fmt.Sprintf("🛡️ Security block: %s", result.ForLLM), nil
			}
		}
	}

	// 1. Update tool contexts
	al.updateToolContexts(agent, opts.Channel, opts.ChatID)

	// 2. Build messages — use ContextManager for budget-aware assembly if available
	var history []providers.Message
	var summary string
	if opts.NoHistory {
		// Heartbeat or special cases: skip history
	} else if al.contextManager != nil {
		// ContextManager.Assemble() reads from its own storage (SQLite or session JSONL)
		// and returns budget-aware context (may include summaries, compressed history)
		assembleResp, err := al.contextManager.Assemble(ctx, &AssembleRequest{
			SessionKey: opts.SessionKey,
			Budget:     agent.ContextWindow,
			MaxTokens:  agent.MaxTokens,
		})
		if err != nil {
			logger.WarnCF("agent", "ContextManager.Assemble failed, falling back to raw history",
				map[string]any{"error": err.Error()})
			history = agent.Sessions.GetHistory(opts.SessionKey)
			summary = agent.Sessions.GetSummary(opts.SessionKey)
		} else if assembleResp != nil {
			history = assembleResp.History
			summary = assembleResp.Summary
		}
	} else {
		// Fallback: raw history from session JSONL
		history = agent.Sessions.GetHistory(opts.SessionKey)
		summary = agent.Sessions.GetSummary(opts.SessionKey)
	}
	messages := agent.ContextBuilder.BuildMessages(
		history,
		summary,
		userMessageForLLM,
		nil,
		opts.Channel,
		opts.ChatID,
	)

	// 3. Save user message to session
	agent.Sessions.AddMessage(opts.SessionKey, "user", originalUserMessage)

	// Ingest into ContextManager for SQLite-backed storage (seahorse)
	if al.contextManager != nil && !opts.NoHistory {
		_ = al.contextManager.Ingest(ctx, &IngestRequest{
			SessionKey: opts.SessionKey,
			Message:    providers.Message{Role: "user", Content: originalUserMessage},
		})
	}

	// Check token budget before processing
	if !al.contextMiddleware.BeforeRequest("llm_call") {
		return "Token budget exceeded. Please wait for GC or simplify your request.", nil
	}
	defer al.contextMiddleware.AfterRequest()

	// 4. Run LLM iteration loop
	finalContent, iteration, usedBinancePublicEndpoint, err := al.runLLMIteration(ctx, agent, messages, opts)
	if err != nil {
		return "", err
	}

	// 4.1. Output Sanitization (Sentinel)
	if al.sentinel != nil && finalContent != "" && !constants.IsInternalChannel(opts.Channel) {
		result := al.sentinel.Execute(ctx, map[string]any{"input": finalContent})
		if result.IsError {
			logger.WarnCF("security", "Sentinel blocked malicious content in assistant response", map[string]any{
				"agent_id":    agent.ID,
				"session_key": opts.SessionKey,
				"pattern":     result.ForLLM,
			})

			if al.auditor != nil {
				al.auditor.LogSecurityEvent(
					agent.ID,
					opts.SessionKey,
					"output_leak_detected",
					finalContent,
					result.ForLLM,
				)
			}

			// Redact the response to prevent leak
			finalContent = "🛡️ [REDACTED] Security block: The assistant's response was blocked by the internal sentinel to prevent information disclosure."
		}
	}

	// If last tool had ForUser content and we already sent it, we might not need to send final response
	// This is controlled by the tool's Silent flag and ForUser content

	// 5. Handle empty response
	if finalContent == "" {
		finalContent = opts.DefaultResponse
	}
	if usedBinancePublicEndpoint {
		notice := "[Notice] Binance ticker data was fetched from the public endpoint because API credentials are not configured."
		if !strings.Contains(strings.ToLower(finalContent), "public endpoint") {
			finalContent = strings.TrimSpace(finalContent)
			if finalContent != "" {
				finalContent += "\n\n" + notice
			} else {
				finalContent = notice
			}
		}
	}

	// 6. Save final assistant message to session
	agent.Sessions.AddMessage(opts.SessionKey, "assistant", finalContent)
	agent.Sessions.Save(opts.SessionKey)

	// 7. Optional: summarization
	if opts.EnableSummary {
		al.maybeSummarize(agent, opts.SessionKey, opts.Channel, opts.ChatID)
	}

	// 8. Optional: send response via bus
	if opts.SendResponse {
		outboundMsg := bus.OutboundMessage{
			Channel: opts.Channel,
			ChatID:  opts.ChatID,
			Content: finalContent,
		}
		// Attach any media generated by tools during this session
		if mediaPaths := al.extractMediaFromSession(agent, opts.SessionKey); len(mediaPaths) > 0 {
			outboundMsg.Media = mediaPaths
		}
		al.bus.PublishOutbound(outboundMsg)
	}

	// 9. Log response
	responsePreview := utils.Truncate(finalContent, 120)
	logger.InfoCF("agent", fmt.Sprintf("Response: %s", responsePreview),
		map[string]any{
			"agent_id":     agent.ID,
			"session_key":  opts.SessionKey,
			"iterations":   iteration,
			"final_length": len(finalContent),
		})

	return finalContent, nil
}

func enrichBinanceShortcutPrompt(input string) string {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return input
	}

	lowered := strings.ToLower(trimmed)
	switch {
	case strings.Contains(lowered, "binance_open_futures_position"),
		strings.Contains(lowered, "binance_close_futures_position"),
		strings.Contains(lowered, "binance_get_order_book"),
		strings.Contains(lowered, "binance_get_futures_order_book"),
		strings.Contains(lowered, "binance_get_spot_balance"),
		strings.Contains(lowered, "binance_get_futures_balance"),
		strings.Contains(lowered, "binance_list_futures_volume"):
		return input
	case strings.Contains(lowered, "open futures"):
		return trimmed + "\n\n[Tool shortcut] Use tool `binance_open_futures_position`. Map side LONG/SHORT. If symbol/quantity are missing, ask a concise follow-up before execution. Use confirm=true only after explicit user confirmation."
	case strings.Contains(lowered, "close futures partial"):
		return trimmed + "\n\n[Tool shortcut] Use tool `binance_close_futures_position` with `quantity` for partial close. If symbol/quantity are missing, ask a concise follow-up before execution. Use confirm=true only after explicit user confirmation."
	case strings.Contains(lowered, "close futures"):
		return trimmed + "\n\n[Tool shortcut] Use tool `binance_close_futures_position`. If quantity is provided, do partial close; otherwise close full net position. If symbol is missing, ask a concise follow-up before execution. Use confirm=true only after explicit user confirmation."
	case strings.Contains(lowered, "spot balance"),
		strings.Contains(lowered, "balance spot"),
		strings.Contains(lowered, "saldo spot"):
		return trimmed + "\n\n[Tool shortcut] Use tool `binance_get_spot_balance`."
	case strings.Contains(lowered, "futures balance"),
		strings.Contains(lowered, "balance futures"),
		strings.Contains(lowered, "saldo futures"),
		strings.Contains(lowered, "saldo de futuros"):
		return trimmed + "\n\n[Tool shortcut] Use tool `binance_get_futures_balance`."
	case strings.Contains(lowered, "list future volume"),
		strings.Contains(lowered, "list futures volume"),
		strings.Contains(lowered, "futures volume list"),
		strings.Contains(lowered, "volumen de futuros"),
		strings.Contains(lowered, "top futures volume"):
		return trimmed + "\n\n[Tool shortcut] Use tool `binance_list_futures_volume`. Default top=10 unless user specifies another value."
	case strings.Contains(lowered, "futures order book"),
		strings.Contains(lowered, "order book futures"),
		strings.Contains(lowered, "futures depth"),
		strings.Contains(lowered, "depth futures"),
		strings.Contains(lowered, "libro de futuros"),
		strings.Contains(lowered, "profundidad futures"):
		return trimmed + "\n\n[Tool shortcut] Use tool `binance_get_futures_order_book`. If symbol is missing, ask a concise follow-up."
	case strings.Contains(lowered, "order book"),
		strings.Contains(lowered, "market depth"),
		strings.Contains(lowered, "bids and asks"),
		strings.Contains(lowered, "libro de ordenes"),
		strings.Contains(lowered, "profundidad del mercado"):
		return trimmed + "\n\n[Tool shortcut] Use tool `binance_get_order_book` for spot order book. If user explicitly requests futures, use `binance_get_futures_order_book` instead. If symbol is missing, ask a concise follow-up."
	default:
		return input
	}
}

// runLLMIteration executes the LLM call loop with tool handling.
func (al *AgentLoop) runLLMIteration(
	ctx context.Context,
	agent *AgentInstance,
	messages []providers.Message,
	opts processOptions,
) (string, int, bool, error) {
	iteration := 0
	var finalContent string
	usedBinancePublicEndpoint := false

	// Create a Task Lock for this execution (Disaster Recovery Checkpoint)
	// TaskID uses the session key + timestamp to handle concurrent independent loops
	taskID := fmt.Sprintf("%s_%d", opts.SessionKey, time.Now().UnixNano())
	var tl *tasklock.TaskLock
	if al.tasklocks != nil {
		var err error
		tl, err = al.tasklocks.CreateLock(taskID, agent.ID, "gateway", opts.SessionKey)
		if err != nil {
			logger.WarnCF("agent", "Failed to create task lock", map[string]any{"error": err.Error()})
		} else {
			defer al.tasklocks.RemoveLock(taskID)
		}
	}

	for iteration < agent.MaxIterations && iteration < HardMaxIterations {
		// Check process context before starting iteration
		select {
		case <-ctx.Done():
			return "", iteration, usedBinancePublicEndpoint, ctx.Err()
		default:
		}

		iteration++

		logger.DebugCF("agent", "LLM iteration",
			map[string]any{
				"agent_id":  agent.ID,
				"iteration": iteration,
				"max":       agent.MaxIterations,
			})

		// Build tool definitions
		// For low-context models (OpenRouter free tier), use only essential tools
		// (5 vs 60+) to save ~10,000+ tokens in tool definitions.
		// Check ALL possible sources because the model may change at runtime via WebUI:
		//   1. Client-specified model in request metadata (WebUI model selector)
		//   2. The resolved model alias from client input
		//   3. IsLowContextModel flag (set at agent creation)
		//   4. Resolved model candidates
		//   5. The primary model in agent.Model itself
		var providerToolDefs []providers.ToolDefinition
		isLowContext := false

		// Check client-specified model first (WebUI model selector)
		if clientModel, ok := opts.Metadata["model_name"]; ok && clientModel != "" {
			// Also resolve the alias in case it's "openrouter-free" -> "openrouter/free"
			resolved := resolveModelAlias(clientModel, al.cfg.ModelList)
			for _, m := range []string{clientModel, resolved} {
				cm := strings.ToLower(m)
				isORFree := strings.HasPrefix(cm, "openrouter/free") ||
					strings.HasPrefix(cm, "openrouter-free") ||
					strings.HasPrefix(cm, "openrouter/auto") ||
					cm == "openrouter-free" || cm == "openrouter/auto"
				if isORFree {
					isLowContext = true
					break
				}
			}
			// Log for debugging
			logger.DebugCF("agent", "Client model check", map[string]any{
				"client_model": clientModel,
				"resolved":     resolved,
				"is_low_ctx":   isLowContext,
			})
		}

		// Fall back to agent-level indicators if client model didn't match
		if !isLowContext {
			isLowContext = agent.IsLowContextModel
		}
		if !isLowContext {
			for _, c := range agent.Candidates {
				lc := strings.ToLower(c.Model)
				if strings.HasPrefix(lc, "openrouter/free") ||
					strings.HasPrefix(lc, "openrouter-free") ||
					strings.HasPrefix(lc, "openrouter/auto") ||
					lc == "openrouter-free" || lc == "openrouter/auto" {
					isLowContext = true
					break
				}
			}
		}
		if !isLowContext {
			lm := strings.ToLower(agent.Model)
			if strings.HasPrefix(lm, "openrouter/free") ||
				strings.HasPrefix(lm, "openrouter-free") ||
				strings.HasPrefix(lm, "openrouter/auto") ||
				lm == "openrouter-free" || lm == "openrouter/auto" {
				isLowContext = true
			}
		}

		if isLowContext {
			providerToolDefs = agent.Tools.ToProviderDefsEssential()
		} else {
			providerToolDefs = agent.Tools.ToProviderDefs()
		}

		logger.DebugCF("agent", "Tool definitions built", map[string]any{
			"count":          len(providerToolDefs),
			"is_low_context": isLowContext,
		})

		// Log LLM request details
		logger.DebugCF("agent", "LLM request",
			map[string]any{
				"agent_id":          agent.ID,
				"iteration":         iteration,
				"model":             agent.Model,
				"messages_count":    len(messages),
				"tools_count":       len(providerToolDefs),
				"max_tokens":        agent.MaxTokens,
				"temperature":       agent.Temperature,
				"system_prompt_len": len(messages[0].Content),
			})

		// Log full messages (detailed)
		logger.DebugCF("agent", "Full LLM request",
			map[string]any{
				"iteration":     iteration,
				"messages_json": formatMessagesForLog(messages),
				"tools_json":    formatToolsForLog(providerToolDefs),
			})

		// SPRINT 1 FEATURE: Apply pruning before calling LLM
		// This trims large tool results in-memory (doesn't modify persisted history)
		pruningConfig := al.cfg.ContextManagement.Pruning
		prunedMessages := PruneMessages(messages, pruningConfig)
		if len(prunedMessages) != len(messages) {
			logger.DebugCF("agent", "Pruning applied", map[string]any{
				"original_count": len(messages),
				"pruned_count":   len(prunedMessages),
			})
		}

		// Call LLM with fallback chain if candidates are configured.
		var response *providers.LLMResponse
		var err error

		// Determine model to use: client-specified model from metadata takes priority
		modelToUse := agent.Model
		if clientModel, ok := opts.Metadata["model_name"]; ok && clientModel != "" {
			// Resolve model alias from model_list (e.g., "gpt-4" → "openai/gpt-4")
			modelToUse = resolveModelAlias(clientModel, al.cfg.ModelList)
			logger.DebugCF("agent", "Using client-specified model", map[string]any{
				"model":        modelToUse,
				"agent_id":     agent.ID,
				"session_key":  opts.SessionKey,
				"client_model": clientModel,
			})
		}

		callLLM := func() (*providers.LLMResponse, error) {
			// FIX: Use the correct provider for the model being used.
			// If the client specified a model that belongs to a different provider
			// (e.g., switching from Antigravity to Qwen), we must resolve the provider
			// dynamically to avoid "model not found" errors from the wrong backend.
			providerToUse := agent.Provider
			modelID := modelToUse

			// Guard against nil provider (misconfigured test or missing model_list entry)
			if providerToUse == nil {
				return nil, fmt.Errorf("no provider available for agent %q — check model_list config", agent.ID)
			}

			if clientModel, ok := opts.Metadata["model_name"]; ok && clientModel != "" {
				// Resolve the provider for the client-specified model
				if p, mid, errProv := al.GetProvider(modelToUse); errProv == nil {
					providerToUse = p
					modelID = mid
				} else {
					logger.WarnCF(
						"agent",
						"Failed to resolve provider for client model, falling back to agent provider",
						map[string]any{
							"model": modelToUse,
							"error": errProv.Error(),
						},
					)
				}

				// Client specified a model — use it directly without fallbacks.
				return providerToUse.Chat(ctx, prunedMessages, providerToolDefs, modelID, map[string]any{
					"max_tokens":       agent.MaxTokens,
					"temperature":      agent.Temperature,
					"prompt_cache_key": agent.ID,
				})
			}

			// No client model specified — use agent's configured fallbacks.
			candidates := agent.Candidates
			if len(candidates) > 1 && al.fallback != nil {
				fbResult, fbErr := al.fallback.Execute(ctx, candidates,
					func(ctx context.Context, pName, mName string) (*providers.LLMResponse, error) {
						// Resolve provider for the fallback candidate
						p, mid, errCand := al.GetProvider(mName)
						if errCand != nil {
							// Fallback to the agent's primary provider if resolution fails
							p = agent.Provider
							mid = mName
						}
						if p == nil {
							return nil, fmt.Errorf("no provider available for fallback model %q", mName)
						}
						return p.Chat(ctx, prunedMessages, providerToolDefs, mid, map[string]any{
							"max_tokens":       agent.MaxTokens,
							"temperature":      agent.Temperature,
							"prompt_cache_key": agent.ID,
						})
					},
				)
				if fbErr != nil {
					return nil, fbErr
				}
				if fbResult.Provider != "" && len(fbResult.Attempts) > 0 {
					logger.InfoCF("agent", fmt.Sprintf("Fallback: succeeded with %s/%s after %d attempts",
						fbResult.Provider, fbResult.Model, len(fbResult.Attempts)+1),
						map[string]any{"agent_id": agent.ID, "iteration": iteration})
				}
				return fbResult.Response, nil
			}

			// Standard case: use the resolved provider and model
			p, mid, errStd := al.GetProvider(modelToUse)
			if errStd == nil {
				providerToUse = p
				modelID = mid
			}
			// Final nil check
			if providerToUse == nil {
				return nil, fmt.Errorf("no provider available for model %q", modelToUse)
			}

			return providerToUse.Chat(ctx, prunedMessages, providerToolDefs, modelID, map[string]any{
				"max_tokens":       agent.MaxTokens,
				"temperature":      agent.Temperature,
				"prompt_cache_key": agent.ID,
			})
		}

		// Wrap callLLM with panic recovery
		safeCallLLM := func() (res *providers.LLMResponse, err error) {
			defer func() {
				if r := recover(); r != nil {
					err = fmt.Errorf("panic in provider.Chat: %v", r)
					logger.ErrorCF("agent", "Recovered from panic during LLM call", map[string]any{
						"agent_id": agent.ID,
						"panic":    r,
					})
				}
			}()
			return callLLM()
		}

		// Checkpoint before hitting network
		if tl != nil {
			tl.UpdateState(tasklock.StatusInProgress, "waiting_for_llm_response", messages)
		}

		// PROACTIVE token validation: Check estimated tokens against model's context window
		// BEFORE sending to LLM. If exceeded, force compression immediately.
		// This prevents the 402 "Prompt tokens limit exceeded" error from OpenRouter
		// by catching the issue before the network call.
		//
		// ⚠️  CRITICAL: MUST include toolTokens from EstimateToolDefsTokens().
		// The original bug used "msgTokens + 2500" which underestimated by ~12,000 tokens
		// because 60+ tool definitions consume ~15,000 tokens, not 2,500.
		msgTokens := al.estimateTokens(messages)
		toolTokens := tokenizer.EstimateToolDefsTokens(providerToolDefs)
		totalEstimated := msgTokens + toolTokens
		// Use 90% of context window as safety margin
		safetyLimit := int(float64(agent.ContextWindow) * 0.9)

		if totalEstimated > safetyLimit {
			// If tool definitions alone are consuming too much budget, switch to essential tools
			if toolTokens > int(float64(agent.ContextWindow)*0.3) && !isLowContext {
				logger.WarnCF(
					"agent",
					"Tool definitions exceed budget threshold, switching to essential tools",
					map[string]any{
						"agent_id":       agent.ID,
						"tool_tokens":    toolTokens,
						"context_window": agent.ContextWindow,
					},
				)
				isLowContext = true
				providerToolDefs = agent.Tools.ToProviderDefsEssential()
				toolTokens = tokenizer.EstimateToolDefsTokens(providerToolDefs)
				totalEstimated = msgTokens + toolTokens
			}

			logger.WarnCF("agent", "Proactive token compression triggered", map[string]any{
				"agent_id":         agent.ID,
				"estimated_tokens": totalEstimated,
				"context_window":   agent.ContextWindow,
				"safety_limit":     safetyLimit,
				"session_key":      opts.SessionKey,
			})

			// Force compression to reduce history
			al.forceCompression(agent, opts.SessionKey)
			newHistory := agent.Sessions.GetHistory(opts.SessionKey)
			newSummary := agent.Sessions.GetSummary(opts.SessionKey)
			messages = agent.ContextBuilder.BuildMessages(
				newHistory, newSummary, "",
				nil, opts.Channel, opts.ChatID,
			)

			// Re-estimate after compression
			msgTokens = al.estimateTokens(messages)
			totalEstimated = msgTokens + toolTokens

			if totalEstimated > safetyLimit {
				logger.WarnCF("agent", "Tokens still high after compression, aggressive truncation", map[string]any{
					"agent_id":       agent.ID,
					"estimated":      totalEstimated,
					"context_window": agent.ContextWindow,
					"safety_limit":   safetyLimit,
					"session_key":    opts.SessionKey,
				})

				// Aggressive truncation: keep reducing until under budget
				keepCounts := []int{5, 3, 2, 1}
				for _, keep := range keepCounts {
					if totalEstimated <= safetyLimit {
						break
					}
					logger.WarnCF("agent", "Truncating to last N messages", map[string]any{
						"agent_id":      agent.ID,
						"estimated":     totalEstimated,
						"safety_limit":  safetyLimit,
						"keep_messages": keep,
						"session_key":   opts.SessionKey,
					})
					al.truncateToLastMessages(opts.SessionKey, keep)
					newHistory = agent.Sessions.GetHistory(opts.SessionKey)
					newSummary = agent.Sessions.GetSummary(opts.SessionKey)
					messages = agent.ContextBuilder.BuildMessages(
						newHistory, newSummary, "",
						nil, opts.Channel, opts.ChatID,
					)
					msgTokens = al.estimateTokens(messages)
					totalEstimated = msgTokens + toolTokens
				}

				if totalEstimated > safetyLimit {
					logger.WarnCF("agent", "CRITICAL: Still over budget, emergency minimal", map[string]any{
						"agent_id":         agent.ID,
						"estimated_tokens": totalEstimated,
						"context_window":   agent.ContextWindow,
						"safety_limit":     safetyLimit,
						"session_key":      opts.SessionKey,
					})
					agent.ContextBuilder.SetPromptLevel(PromptLevelMinimal)
					al.truncateToLastMessages(opts.SessionKey, 1)
					newHistory = agent.Sessions.GetHistory(opts.SessionKey)
					newSummary = agent.Sessions.GetSummary(opts.SessionKey)
					messages = agent.ContextBuilder.BuildMessages(
						newHistory, newSummary, "",
						nil, opts.Channel, opts.ChatID,
					)
				}
			}
		}

		// Retry loop for context/token errors
		maxRetries := 4
		for retry := 0; retry <= maxRetries; retry++ {
			response, err = safeCallLLM()
			if err == nil {
				break
			}

			if tl != nil {
				tl.UpdateState(tasklock.StatusNetworkRetry, "llm_error_retry", nil)
				tl.IncrementRetry()
			}

			errMsg := strings.ToLower(err.Error())

			// Handle rate limit errors
			isRateLimit := strings.Contains(errMsg, "rate limit") ||
				strings.Contains(errMsg, "429") ||
				strings.Contains(errMsg, "exhausted") ||
				strings.Contains(errMsg, "too many requests")

			if isRateLimit && retry < maxRetries {
				backoff := time.Duration(1<<retry) * 5 * time.Second // 5s, 10s, 20s, 40s
				logger.WarnCF("agent", "Rate limit hit, backing off before retry", map[string]any{
					"attempt": retry + 1,
					"max":     maxRetries,
					"sleep":   backoff.String(),
					"error":   err.Error(),
				})

				select {
				case <-ctx.Done():
					return "", iteration, usedBinancePublicEndpoint, ctx.Err()
				case <-time.After(backoff):
					continue
				}
			}

			isContextError := strings.Contains(errMsg, "token") ||
				strings.Contains(errMsg, "context") ||
				strings.Contains(errMsg, "invalidparameter") ||
				strings.Contains(errMsg, "length")

			if isContextError && retry < maxRetries {
				logger.WarnCF("agent", "Context window error detected, attempting compression", map[string]any{
					"error": err.Error(),
					"retry": retry,
				})

				if retry == 0 && !constants.IsInternalChannel(opts.Channel) {
					al.bus.PublishOutbound(bus.OutboundMessage{
						Channel: opts.Channel,
						ChatID:  opts.ChatID,
						Content: "Context window exceeded. Compressing history and retrying...",
					})
				}

				// Use ContextManager.Compact() if available (seahorse or legacy)
				if al.contextManager != nil {
					compactErr := al.contextManager.Compact(ctx, &CompactRequest{
						SessionKey: opts.SessionKey,
						Reason:     ContextCompressReasonRetry,
						Budget:     agent.ContextWindow,
					})
					if compactErr != nil {
						logger.ErrorCF("agent", "ContextManager.Compact failed, falling back to forceCompression",
							map[string]any{"error": compactErr.Error()})
						al.forceCompression(agent, opts.SessionKey)
					}
				} else {
					al.forceCompression(agent, opts.SessionKey)
				}

				newHistory := agent.Sessions.GetHistory(opts.SessionKey)
				newSummary := agent.Sessions.GetSummary(opts.SessionKey)
				messages = agent.ContextBuilder.BuildMessages(
					newHistory, newSummary, "",
					nil, opts.Channel, opts.ChatID,
				)
				continue
			}
			break
		}

		if err != nil {
			logger.ErrorCF("agent", "LLM call failed",
				map[string]any{
					"agent_id":  agent.ID,
					"iteration": iteration,
					"error":     err.Error(),
				})
			return "", iteration, usedBinancePublicEndpoint, fmt.Errorf("LLM call failed after retries: %w", err)
		}

		// Check if no tool calls - we're done
		if len(response.ToolCalls) == 0 {
			finalContent = response.Content
			logger.InfoCF("agent", "LLM response without tool calls (direct answer)",
				map[string]any{
					"agent_id":      agent.ID,
					"iteration":     iteration,
					"content_chars": len(finalContent),
				})
			break
		}

		normalizedToolCalls := make([]providers.ToolCall, 0, len(response.ToolCalls))
		for _, tc := range response.ToolCalls {
			normalizedToolCalls = append(normalizedToolCalls, providers.NormalizeToolCall(tc))
		}

		// Log tool calls
		toolNames := make([]string, 0, len(normalizedToolCalls))
		for _, tc := range normalizedToolCalls {
			toolNames = append(toolNames, tc.Name)
		}
		logger.InfoCF("agent", "LLM requested tool calls",
			map[string]any{
				"agent_id":  agent.ID,
				"tools":     toolNames,
				"count":     len(normalizedToolCalls),
				"iteration": iteration,
			})

		// Build assistant message with tool calls
		assistantMsg := providers.Message{
			Role:             "assistant",
			Content:          response.Content,
			ReasoningContent: response.ReasoningContent,
			Source:           "assistant",
		}
		for _, tc := range normalizedToolCalls {
			argumentsJSON, _ := json.Marshal(tc.Arguments)
			// Copy ExtraContent to ensure thought_signature is persisted for Gemini 3
			extraContent := tc.ExtraContent
			thoughtSignature := ""
			if tc.Function != nil {
				thoughtSignature = tc.Function.ThoughtSignature
			}

			assistantMsg.ToolCalls = append(assistantMsg.ToolCalls, providers.ToolCall{
				ID:   tc.ID,
				Type: "function",
				Name: tc.Name,
				Function: &providers.FunctionCall{
					Name:             tc.Name,
					Arguments:        string(argumentsJSON),
					ThoughtSignature: thoughtSignature,
				},
				ExtraContent:     extraContent,
				ThoughtSignature: thoughtSignature,
			})
		}
		messages = append(messages, assistantMsg)

		// Save assistant message with tool calls to session
		agent.Sessions.AddFullMessage(opts.SessionKey, assistantMsg)

		// Checkpoint before executing tools
		if tl != nil {
			tl.UpdateState(tasklock.StatusInProgress, "executing_tools", messages)
		}

		// Execute tool calls with validation to prevent ghost conversations
		for _, tc := range normalizedToolCalls {
			argsJSON, _ := json.Marshal(tc.Arguments)
			argsPreview := utils.Truncate(string(argsJSON), 200)

			// Validate if this tool call seems to originate from user intent vs context confusion
			shouldExecute := true
			if tc.Name == "message" {
				// For message tool, check if the content looks like it came from user intent
				if content, ok := tc.Arguments["content"].(string); ok {
					// Check if the content contains typical reminder/request phrases that might come from context
					lowerContent := strings.ToLower(content)

					// Common phrases that might indicate context confusion vs user request
					ghostIndicators := []string{
						"remind me", "remember me", "i have to", "i need to", "need to",
					}

					for _, indicator := range ghostIndicators {
						if strings.Contains(lowerContent, strings.ToLower(indicator)) {
							logger.InfoCF(
								"agent",
								"Potential ghost conversation detected, validating message tool call",
								map[string]any{
									"agent_id":  agent.ID,
									"tool":      tc.Name,
									"content":   content,
									"iteration": iteration,
								},
							)

							// Additional validation: check if this looks like a user-initiated request
							// For now, we'll log it but still execute - in future iterations we could add more sophisticated checks
							break
						}
					}
				}
			}

			logger.InfoCF("agent", fmt.Sprintf("Tool call: %s(%s)", tc.Name, argsPreview),
				map[string]any{
					"agent_id":  agent.ID,
					"tool":      tc.Name,
					"iteration": iteration,
				})

			// Create async callback for tools that implement AsyncTool
			// NOTE: Following openclaw's design, async tools do NOT send results directly to users.
			// Instead, they notify the agent via PublishInbound, and the agent decides
			// whether to forward the result to the user (in processSystemMessage).
			asyncCallback := func(callbackCtx context.Context, result *tools.ToolResult) {
				// Log the async completion but don't send directly to user
				// The agent will handle user notification via processSystemMessage
				if !result.Silent && result.ForUser != "" {
					logger.InfoCF("agent", "Async tool completed, agent will handle notification",
						map[string]any{
							"tool":        tc.Name,
							"content_len": len(result.ForUser),
						})
				}
			}

			var toolResult *tools.ToolResult
			if shouldExecute {
				toolResult = agent.Tools.ExecuteWithContext(
					ctx,
					tc.Name,
					tc.Arguments,
					opts.Channel,
					opts.ChatID,
					asyncCallback,
				)
			} else {
				// Skip execution and return an appropriate result
				toolResult = tools.NewToolResult("Skipped execution due to potential context confusion validation")
				logger.InfoCF("agent", "Tool execution skipped due to validation",
					map[string]any{
						"agent_id":  agent.ID,
						"tool":      tc.Name,
						"iteration": iteration,
					})
			}

			// Send ForUser content to user immediately if not Silent
			if !toolResult.Silent && toolResult.ForUser != "" && opts.SendResponse {
				outboundMsg := bus.OutboundMessage{
					Channel: opts.Channel,
					ChatID:  opts.ChatID,
					Content: toolResult.ForUser,
					Buttons: toolResult.Buttons,
				}
				if len(toolResult.MediaPaths) > 0 {
					outboundMsg.Media = toolResult.MediaPaths
				}
				al.bus.PublishOutbound(outboundMsg)
				logger.DebugCF("agent", "Sent tool result to user",
					map[string]any{
						"tool":        tc.Name,
						"content_len": len(toolResult.ForUser),
					})
			}

			// Determine content for LLM based on tool result
			contentForLLM := toolResult.ForLLM
			if contentForLLM == "" && toolResult.Err != nil {
				contentForLLM = toolResult.Err.Error()
			}

			// Issue 1 Patch: Filter sensitive data if enabled
			contentForLLM = al.filterSensitiveData(contentForLLM)

			if tc.Name == "binance_get_ticker_price" &&
				strings.Contains(strings.ToLower(contentForLLM), "source: binance public endpoint") {
				usedBinancePublicEndpoint = true
			}

			toolResultMsg := providers.Message{
				Role:       "tool",
				Content:    contentForLLM,
				ToolCallID: tc.ID,
				Source:     "tool_result",
			}
			if len(toolResult.MediaPaths) > 0 {
				toolResultMsg.MediaPaths = toolResult.MediaPaths
			}
			messages = append(messages, toolResultMsg)

			// Save tool result message to session
			agent.Sessions.AddFullMessage(opts.SessionKey, toolResultMsg)
		}
	}

	return finalContent, iteration, usedBinancePublicEndpoint, nil
}

// updateToolContexts updates the context for tools that need channel/chatID info.
func (al *AgentLoop) updateToolContexts(agent *AgentInstance, channel, chatID string) {
	// Use ContextualTool interface instead of type assertions
	if tool, ok := agent.Tools.Get("message"); ok {
		if mt, ok := tool.(tools.ContextualTool); ok {
			mt.SetContext(channel, chatID)
		}
	}
	if tool, ok := agent.Tools.Get("spawn"); ok {
		if st, ok := tool.(tools.ContextualTool); ok {
			st.SetContext(channel, chatID)
		}
	}
	if tool, ok := agent.Tools.Get("subagent"); ok {
		if st, ok := tool.(tools.ContextualTool); ok {
			st.SetContext(channel, chatID)
		}
	}

	// Inject file lock checking tools
	if al.tasklocks != nil {
		for _, toolName := range []string{"write_file", "edit_file", "append_file"} {
			if tool, ok := agent.Tools.Get(toolName); ok {
				if lt, ok := tool.(tools.LockAwareTool); ok {
					lt.SetLockChecker(al.tasklocks)
				}
			}
		}
	}
}

// maybeSummarize triggers summarization if the session history exceeds thresholds.
func (al *AgentLoop) maybeSummarize(agent *AgentInstance, sessionKey, channel, chatID string) {
	newHistory := agent.Sessions.GetHistory(sessionKey)
	tokenEstimate := al.estimateTokens(newHistory)

	compactThreshold := 0.60
	if al.cfg != nil && al.cfg.ContextManagement.CompactThreshold > 0 {
		compactThreshold = al.cfg.ContextManagement.CompactThreshold
	}

	// OpenRouter free models have ~4096 context window.
	// System prompt (~2000) + tool defs (~800) = ~2800 already used.
	// Only ~1300 tokens remain for history. Use 25% threshold to compact early.
	// Use the agent's IsLowContextModel flag for reliable detection.
	if agent.IsLowContextModel {
		compactThreshold = 0.25 // Much earlier compaction for low-context models
	}

	thresholdTokens := int(float64(agent.ContextWindow) * compactThreshold)

	if len(newHistory) > 30 || tokenEstimate > thresholdTokens {
		summarizeKey := agent.ID + ":" + sessionKey
		if _, loading := al.summarizing.LoadOrStore(summarizeKey, true); !loading {
			go func() {
				defer al.summarizing.Delete(summarizeKey)

				if al.cfg != nil && al.cfg.ContextManagement.AutoCompactEnabled &&
					!constants.IsInternalChannel(channel) {
					criticalThreshold := 0.80
					if al.cfg.ContextManagement.CriticalThreshold > 0 {
						criticalThreshold = al.cfg.ContextManagement.CriticalThreshold
					}
					criticalTokens := int(float64(agent.ContextWindow) * criticalThreshold)

					if tokenEstimate > criticalTokens {
						al.bus.PublishOutbound(bus.OutboundMessage{
							Channel: channel,
							ChatID:  chatID,
							Content: "⚠️ Conversation too heavy. Executing automatic compaction...",
						})
					}
					// Only log internally for the 60% threshold to reduce user noise
					logger.DebugCF("agent", "Heavy conversation (60% threshold reached)", map[string]any{
						"token_estimate": tokenEstimate,
						"threshold":      thresholdTokens,
					})
				}

				logger.Debug("Memory threshold reached. Optimizing conversation history...")
				al.summarizeSession(agent, sessionKey)
			}()
		}
	}
}

// forceCompression aggressively reduces context when the limit is hit.
// It drops the oldest 50% of messages (keeping system prompt and last user message).
func (al *AgentLoop) forceCompression(agent *AgentInstance, sessionKey string) {
	history := agent.Sessions.GetHistory(sessionKey)
	if len(history) <= 4 {
		return
	}

	// Keep system prompt (usually [0]) and the very last message (user's trigger)
	// We want to drop the oldest half of the *conversation*
	// Assuming [0] is system, [1:] is conversation
	conversation := history[1 : len(history)-1]
	if len(conversation) == 0 {
		return
	}

	// Helper to find the mid-point of the conversation
	mid := len(conversation) / 2

	// New history structure:
	// 1. System Prompt (with compression note appended)
	// 2. Second half of conversation
	// 3. Last message

	droppedCount := mid
	keptConversation := conversation[mid:]

	newHistory := make([]providers.Message, 0, 1+len(keptConversation)+1)

	// Append compression note to the original system prompt instead of adding a new system message
	// This avoids having two consecutive system messages which some APIs (like Zhipu) reject
	compressionNote := fmt.Sprintf(
		"\n\n[System Note: Emergency compression dropped %d oldest messages due to context limit]",
		droppedCount,
	)
	enhancedSystemPrompt := history[0]
	enhancedSystemPrompt.Content = enhancedSystemPrompt.Content + compressionNote
	newHistory = append(newHistory, enhancedSystemPrompt)

	newHistory = append(newHistory, keptConversation...)
	newHistory = append(newHistory, history[len(history)-1]) // Last message

	// --- START PATCH: Orphan tool cleaning ---
	// Ensure we don't leave responses or tool calls at the beginning
	for len(newHistory) > 0 {
		// 1. Protect index 0 if it's the System Prompt
		checkIdx := 0
		if newHistory[0].Role == "system" {
			if len(newHistory) == 1 {
				break // Solo queda el prompt del sistema, estamos a salvo
			}
			checkIdx = 1
		}

		msg := newHistory[checkIdx]

		// 2. If the message is of type "tool" (response to a tool) without an "assistant" before, it's invalid
		if msg.Role == "tool" {
			newHistory = append(newHistory[:checkIdx], newHistory[checkIdx+1:]...)
			continue
		}

		// 3. Opcional pero recomendado: Si es "assistant" con tool_calls, pero le falta la respuesta ("tool")
		if msg.Role == "assistant" && len(msg.ToolCalls) > 0 {
			if len(newHistory) <= checkIdx+1 || newHistory[checkIdx+1].Role != "tool" {
				newHistory = append(newHistory[:checkIdx], newHistory[checkIdx+1:]...)
				continue
			}
		}

		// If it's not an orphan block, exit the loop
		break
	}
	// --- FIN PARCHE ---

	// Update session
	agent.Sessions.SetHistory(sessionKey, newHistory)
	agent.Sessions.Save(sessionKey)

	logger.WarnCF("agent", "Forced compression executed", map[string]any{
		"session_key":  sessionKey,
		"dropped_msgs": droppedCount,
		"new_count":    len(newHistory),
	})
}

// truncateToLastMessages aggressively truncates session history to keep only the
// last N messages plus the system prompt. Used as a last resort when standard
// compression is insufficient to fit within the model's context window.
func (al *AgentLoop) truncateToLastMessages(sessionKey string, keep int) {
	agent := al.registry.GetDefaultAgent()
	if agent == nil {
		return
	}

	history := agent.Sessions.GetHistory(sessionKey)
	if len(history) <= keep+1 { // +1 for system prompt
		return
	}

	// Keep system prompt and last N messages
	systemMsg := history[0]
	truncationNote := fmt.Sprintf(
		"\n\n[System Note: Aggressive truncation — kept only last %d messages due to context limit]",
		keep,
	)
	systemMsg.Content = systemMsg.Content + truncationNote

	startIdx := len(history) - keep
	if startIdx < 1 {
		startIdx = 1
	}

	newHistory := make([]providers.Message, 0, 1+keep)
	newHistory = append(newHistory, systemMsg)
	newHistory = append(newHistory, history[startIdx:]...)

	// Clean orphan tool messages (same logic as forceCompression)
	for len(newHistory) > 0 {
		checkIdx := 0
		if newHistory[0].Role == "system" {
			if len(newHistory) == 1 {
				break
			}
			checkIdx = 1
		}

		msg := newHistory[checkIdx]
		if msg.Role == "tool" {
			newHistory = append(newHistory[:checkIdx], newHistory[checkIdx+1:]...)
			continue
		}

		if msg.Role == "assistant" && len(msg.ToolCalls) > 0 {
			if len(newHistory) <= checkIdx+1 || newHistory[checkIdx+1].Role != "tool" {
				newHistory = append(newHistory[:checkIdx], newHistory[checkIdx+1:]...)
				continue
			}
		}

		break
	}

	agent.Sessions.SetHistory(sessionKey, newHistory)
	agent.Sessions.Save(sessionKey)

	logger.WarnCF("agent", "Aggressive truncation executed", map[string]any{
		"session_key": sessionKey,
		"kept_msgs":   len(newHistory),
		"dropped":     len(history) - len(newHistory),
	})
}

// GetStartupInfo returns information about loaded tools and skills for logging.
func (al *AgentLoop) GetStartupInfo() map[string]any {
	info := make(map[string]any)

	agent := al.registry.GetDefaultAgent()
	if agent == nil {
		return info
	}

	// Tools info
	toolsList := agent.Tools.List()
	info["tools"] = map[string]any{
		"count": len(toolsList),
		"names": toolsList,
	}

	// Skills info
	info["skills"] = agent.ContextBuilder.GetSkillsInfo()

	// Agents info
	info["agents"] = map[string]any{
		"count": len(al.registry.ListAgentIDs()),
		"ids":   al.registry.ListAgentIDs(),
	}

	return info
}

// formatMessagesForLog formats messages for logging
func formatMessagesForLog(messages []providers.Message) string {
	if len(messages) == 0 {
		return "[]"
	}

	var sb strings.Builder
	sb.WriteString("[\n")
	for i, msg := range messages {
		fmt.Fprintf(&sb, "  [%d] Role: %s\n", i, msg.Role)
		if len(msg.ToolCalls) > 0 {
			sb.WriteString("  ToolCalls:\n")
			for _, tc := range msg.ToolCalls {
				fmt.Fprintf(&sb, "    - ID: %s, Type: %s, Name: %s\n", tc.ID, tc.Type, tc.Name)
				if tc.Function != nil {
					fmt.Fprintf(&sb, "      Arguments: %s\n", utils.Truncate(tc.Function.Arguments, 200))
				}
			}
		}
		if msg.Content != "" {
			content := utils.Truncate(msg.Content, 200)
			fmt.Fprintf(&sb, "  Content: %s\n", content)
		}
		if msg.ToolCallID != "" {
			fmt.Fprintf(&sb, "  ToolCallID: %s\n", msg.ToolCallID)
		}
		sb.WriteString("\n")
	}
	sb.WriteString("]")
	return sb.String()
}

// formatToolsForLog formats tool definitions for logging
func formatToolsForLog(toolDefs []providers.ToolDefinition) string {
	if len(toolDefs) == 0 {
		return "[]"
	}

	var sb strings.Builder
	sb.WriteString("[\n")
	for i, tool := range toolDefs {
		fmt.Fprintf(&sb, "  [%d] Type: %s, Name: %s\n", i, tool.Type, tool.Function.Name)
		fmt.Fprintf(&sb, "      Description: %s\n", tool.Function.Description)
		if len(tool.Function.Parameters) > 0 {
			fmt.Fprintf(&sb, "      Parameters: %s\n", utils.Truncate(fmt.Sprintf("%v", tool.Function.Parameters), 200))
		}
	}
	sb.WriteString("]")
	return sb.String()
}

// summarizeSession summarizes the conversation history for a session.
func (al *AgentLoop) summarizeSession(agent *AgentInstance, sessionKey string) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	history := agent.Sessions.GetHistory(sessionKey)

	compactModel := agent.Model
	if al.cfg != nil && al.cfg.ContextManagement.Compaction.Model != "" {
		compactModel = al.cfg.ContextManagement.Compaction.Model
	}
	var compactionCfg config.ContextCompactionConfig
	if al.cfg != nil {
		compactionCfg = al.cfg.ContextManagement.Compaction
	}
	compactor := NewDefaultContextCompactorWithCfg(al.summaryCache, compactionCfg)
	compacted, err := compactor.CompactMessages(ctx, agent.Provider, compactModel, history, sessionKey)
	if err == nil && len(compacted) < len(history) {
		agent.Sessions.SetHistory(sessionKey, compacted)
		agent.Sessions.Save(sessionKey)
		logger.InfoCF(
			"agent",
			"Context compaction complete",
			map[string]any{"session_key": sessionKey, "original_len": len(history), "compacted_len": len(compacted)},
		)
	} else if err != nil {
		logger.ErrorCF("agent", "Failed to compact messages", map[string]any{"error": err.Error()})
	}
}

// estimateTokens estimates the number of tokens in a message list.
// Uses tokenizer.EstimateMessageTokens() to count Content, ReasoningContent,
// ToolCalls, ToolCallID, and SystemParts — NOT just Content characters.
//
// ⚠️  CRITICAL: DO NOT change this to a simple character count (e.g., utf8.RuneCountInString).
// The original bug used char counting + fixed "+2500" overhead, which underestimated
// by ~18,000 tokens and caused constant OpenRouter 402 errors.
// See file header comment and local_work/MEMORY.md for details.
func (al *AgentLoop) estimateTokens(messages []providers.Message) int {
	total := 0
	for _, m := range messages {
		total += tokenizer.EstimateMessageTokens(m)
	}
	return total
}

func (al *AgentLoop) handleCommand(ctx context.Context, msg bus.InboundMessage) (string, bool) {
	content := strings.TrimSpace(msg.Content)
	if content == "" {
		return "", false
	}

	// Fast Path for Batch IDs: #PREFIX_DD_MM_YY_HH_MM
	if strings.HasPrefix(content, "#") {
		// Route directly to queue tool
		agent := al.registry.GetDefaultAgent()
		if agent != nil {
			logger.InfoCF("agent", "Fast-path triggered for batch ID", map[string]any{"id": content})
			res := agent.Tools.ExecuteWithContext(
				ctx,
				"queue",
				map[string]any{"action": "status", "task_id": content},
				msg.Channel,
				msg.ChatID,
				nil,
			)
			if res != nil {
				return res.ForUser, true
			}
		}
	}

	if !strings.HasPrefix(content, "/") {
		return "", false
	}

	parts := strings.Fields(content)
	if len(parts) == 0 {
		return "", false
	}

	cmd := parts[0]
	args := parts[1:]

	switch cmd {
	case "/bundle_approve", "/bundle_regen", "/bundle_edit", "/bundle_publish", "/bundle_cancel":
		// Direct action on social bundle without LLM
		agent := al.registry.GetDefaultAgent()
		if agent != nil {
			id := ""
			platforms := ""
			for _, arg := range args {
				if strings.HasPrefix(arg, "id=") {
					id = strings.TrimPrefix(arg, "id=")
				} else if strings.HasPrefix(arg, "platforms=") {
					platforms = strings.TrimPrefix(arg, "platforms=")
				}
			}
			if id == "" {
				return "Error: Faltó ID del bundle.", true
			}

			logger.InfoCF("agent", "Fast-path triggered for bundle action", map[string]any{
				"cmd":       cmd,
				"id":        id,
				"platforms": platforms,
			})

			action := "approve"
			if cmd == "/bundle_regen" {
				action = "regenerate"
			} else if cmd == "/bundle_edit" {
				action = "edit"
			} else if cmd == "/bundle_publish" {
				action = "publish"
			} else if cmd == "/bundle_cancel" {
				action = "cancel"
			}

			// Execute tool directly
			toolArgs := map[string]any{
				"action": action,
				"id":     id,
			}
			if platforms != "" {
				toolArgs["platforms"] = platforms
			}

			res := agent.Tools.ExecuteWithContext(ctx, "social_manager", toolArgs, msg.Channel, msg.ChatID, nil)

			if res != nil {
				if res.ForUser != "" {
					return res.ForUser, true
				}
				return res.ForLLM, true
			}
		}
		return "Acción enviada.", true

	case "/show":
		if len(args) < 1 {
			return "Usage: /show [model|channel|agents]", true
		}
		switch args[0] {
		case "model":
			defaultAgent := al.registry.GetDefaultAgent()
			if defaultAgent == nil {
				return "No default agent configured", true
			}
			return fmt.Sprintf("Current model: %s", defaultAgent.Model), true
		case "channel":
			return fmt.Sprintf("Current channel: %s", msg.Channel), true
		case "agents":
			agentIDs := al.registry.ListAgentIDs()
			return fmt.Sprintf("Registered agents: %s", strings.Join(agentIDs, ", ")), true
		default:
			return fmt.Sprintf("Unknown show target: %s", args[0]), true
		}

	case "/status":
		agent := al.registry.GetDefaultAgent()
		if agent == nil {
			return "No default agent configured", true
		}
		history := agent.Sessions.GetHistory(msg.SessionKey)
		tokens := al.estimateTokens(history)
		pct := 0.0
		if agent.ContextWindow > 0 {
			pct = float64(tokens) / float64(agent.ContextWindow) * 100
		}
		return fmt.Sprintf("Tokens used: %d / %d (%.1f%%)", tokens, agent.ContextWindow, pct), true

	case "/list":
		if len(args) < 1 {
			return "Usage: /list [models|channels|agents]", true
		}
		switch args[0] {
		case "models":
			return "Available models: configured in config.json per agent", true
		case "channels":
			if al.channelManager == nil {
				return "Channel manager not initialized", true
			}
			channels := al.channelManager.GetEnabledChannels()
			if len(channels) == 0 {
				return "No channels enabled", true
			}
			return fmt.Sprintf("Enabled channels: %s", strings.Join(channels, ", ")), true
		case "agents":
			agentIDs := al.registry.ListAgentIDs()
			return fmt.Sprintf("Registered agents: %s", strings.Join(agentIDs, ", ")), true
		default:
			return fmt.Sprintf("Unknown list target: %s", args[0]), true
		}

	case "/switch":
		if len(args) < 3 || args[1] != "to" {
			return "Usage: /switch [model|channel] to <name>", true
		}
		target := args[0]
		value := args[2]

		switch target {
		case "model":
			defaultAgent := al.registry.GetDefaultAgent()
			if defaultAgent == nil {
				return "No default agent configured", true
			}
			oldModel := defaultAgent.Model
			defaultAgent.Model = value
			return fmt.Sprintf("Switched model from %s to %s", oldModel, value), true
		case "channel":
			if al.channelManager == nil {
				return "Channel manager not initialized", true
			}
			if _, exists := al.channelManager.GetChannel(value); !exists && value != "cli" {
				return fmt.Sprintf("Channel '%s' not found or not enabled", value), true
			}
			return fmt.Sprintf("Switched target channel to %s", value), true
		default:
			return fmt.Sprintf("Unknown switch target: %s", target), true
		}

	case "/disable_sentinel":
		if len(args) < 1 {
			return "⚠️ Usage: /disable_sentinel [5m|15m|1h]\n\n⚠️ **WARNING**: Sentinel disabled. All files visible.", true
		}

		duration := args[0]
		var d time.Duration

		switch duration {
		case "5m":
			d = 5 * time.Minute
		case "15m":
			d = 15 * time.Minute
		case "1h":
			d = 1 * time.Hour
		default:
			return "⚠️ Invalid duration. Use: 5m, 15m, or 1h", true
		}

		al.sentinel.Disable(d)

		// Audit log
		al.auditor.LogSecurityEvent("system", msg.SessionKey, "sentinel_disabled",
			fmt.Sprintf("User disabled sentinel for %s", duration),
			fmt.Sprintf("Channel: %s, ChatID: %s", msg.Channel, msg.ChatID))

		return "⚠️ **SENTINEL DISABLED** | Caution: All files visible for " + duration, true

	case "/activate_sentinel":
		al.sentinel.Enable()

		// Audit log
		al.auditor.LogSecurityEvent("system", msg.SessionKey, "sentinel_activated",
			"User manually activated sentinel",
			fmt.Sprintf("Channel: %s, ChatID: %s", msg.Channel, msg.ChatID))

		return "✅ **SENTINEL ACTIVATED** | Security checks enabled", true

	case "/sentinel_status":
		status := al.sentinel.GetStatus()
		if status == "active" {
			return "🛡️ Sentinel Status: **ACTIVE**", true
		}
		if strings.HasPrefix(status, "disabled_") {
			remaining := strings.TrimPrefix(status, "disabled_")
			return "⚠️ Sentinel Status: **DISABLED** | Remaining: " + remaining, true
		}
		return "🛡️ Sentinel Status: " + status, true

	// SPRINT 1 FEATURE: Manual context compaction command
	case "/compact":
		agent := al.registry.GetDefaultAgent()
		if agent == nil {
			return "No default agent configured", true
		}

		// Get optional focus instructions
		instructions := strings.TrimSpace(strings.TrimPrefix(content, "/compact"))

		// Load session history
		history := agent.Sessions.GetHistory(msg.SessionKey)
		if len(history) <= 6 {
			return "Context is too short for compaction (need more than 6 messages)", true
		}

		// Trigger compaction
		provider := agent.Provider
		model := agent.Model

		// Use compaction model if configured
		if al.cfg != nil && al.cfg.ContextManagement.Compaction.Model != "" {
			model = al.cfg.ContextManagement.Compaction.Model
		}

		// Build compactor if not exists
		if al.summaryCache == nil {
			al.summaryCache = utils.NewSummaryCache(agent.Workspace)
		}
		compactor := NewDefaultContextCompactor(al.summaryCache)

		// Execute compaction
		compacted, err := compactor.CompactMessages(ctx, provider, model, history, msg.SessionKey)
		if err != nil {
			return fmt.Sprintf("❌ Compaction failed: %v", err), true
		}

		// Replace session history with compacted version
		agent.Sessions.SetHistory(msg.SessionKey, compacted)
		agent.Sessions.Save(msg.SessionKey)

		response := "✅ Context compacted successfully."
		if instructions != "" {
			response += fmt.Sprintf(" Focus: %s", instructions)
		}
		return response, true

	case "/restrict_to_workspace":
		if len(args) < 1 {
			return "⚠️ Usage: /restrict_to_workspace [activate|deactivate|status]\n\n" +
				"🔒 **Security Command**: Controls workspace file access restrictions\n\n" +
				"• `activate`: Agent can ONLY access files within workspace (SAFE)\n" +
				"• `deactivate`: Agent can access ANY system file (DANGEROUS)\n" +
				"• `status`: Show current restriction status", true
		}

		action := strings.ToLower(args[0])

		if !al.isUserAuthorized(msg.Channel, msg.ChatID) {
			al.auditor.LogSecurityEvent("system", msg.SessionKey, "unauthorized_config_change",
				fmt.Sprintf("User attempted to change restrict_to_workspace: %s", action),
				fmt.Sprintf("Channel: %s, ChatID: %s", msg.Channel, msg.ChatID))
			return "❌ **Unauthorized**: Only authorized users can change security settings", true
		}

		return al.handleRestrictToWorkspace(action, msg.Channel, msg.ChatID, msg.SessionKey), true
	}

	if strings.HasPrefix(content, "/") {
		return "⚠️ Commando desconocido. Escribe /help para ver la lista de commandos disponibles.", true
	}

	return "", false
}

// handleRestrictToWorkspace handles the /restrict_to_workspace fast-path command.
// It reads the config file, modifies restrict_to_workspace, and writes atomically.
func (al *AgentLoop) handleRestrictToWorkspace(action, channel, chatID, sessionKey string) string {
	configPath := al.cfg.ConfigPath()
	if configPath == "" {
		return "❌ Config path not available (run via picoclaw binary, not test mode)"
	}

	al.configMutex.Lock()
	defer al.configMutex.Unlock()

	// Read current config
	data, err := os.ReadFile(configPath)
	if err != nil {
		logger.ErrorCF("agent", "Failed to read config for restrict_to_workspace",
			map[string]any{"error": err.Error(), "path": configPath})
		return fmt.Sprintf("❌ Error reading config: %v", err)
	}

	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		logger.ErrorCF("agent", "Failed to parse config JSON", map[string]any{"error": err.Error()})
		return fmt.Sprintf("❌ Error parsing config: %v", err)
	}

	// Navigate to agents.defaults
	agents, ok := raw["agents"].(map[string]any)
	if !ok {
		return "❌ Invalid config structure: missing 'agents' section"
	}
	defaults, ok := agents["defaults"].(map[string]any)
	if !ok {
		return "❌ Invalid config structure: missing 'agents.defaults' section"
	}

	currentValue, _ := defaults["restrict_to_workspace"].(bool)

	if action == "status" {
		statusLabel := "⚠️ DEACTIVATED"
		if currentValue {
			statusLabel = "✅ ACTIVATED"
		}
		return fmt.Sprintf("🔒 **restrict_to_workspace Status**: %s\n\nCurrent value: %v\n\n"+
			"• `activate`: Restrict to workspace (SAFE)\n"+
			"• `deactivate`: Allow full system access (DANGEROUS)", statusLabel, currentValue)
	}

	var newValue bool
	var actionDescription string
	switch action {
	case "activate":
		newValue = true
		actionDescription = "ACTIVATED"
	case "deactivate":
		newValue = false
		actionDescription = "DEACTIVATED"
	default:
		return "❌ Invalid action. Use: activate, deactivate, or status"
	}

	if currentValue == newValue {
		return fmt.Sprintf("ℹ️ **restrict_to_workspace** is already %s", strings.ToLower(actionDescription))
	}

	// Update the JSON map and write atomically
	defaults["restrict_to_workspace"] = newValue

	updatedData, err := json.MarshalIndent(raw, "", "  ")
	if err != nil {
		logger.ErrorCF("agent", "Failed to marshal updated config", map[string]any{"error": err.Error()})
		return fmt.Sprintf("❌ Error serializing config: %v", err)
	}

	tmpPath := configPath + ".tmp"
	if err := os.WriteFile(tmpPath, updatedData, 0o600); err != nil {
		logger.ErrorCF("agent", "Failed to write temp config", map[string]any{"error": err.Error()})
		return fmt.Sprintf("❌ Error writing config: %v", err)
	}
	if err := os.Rename(tmpPath, configPath); err != nil {
		logger.ErrorCF("agent", "Failed to rename temp config", map[string]any{"error": err.Error()})
		os.Remove(tmpPath)
		return fmt.Sprintf("❌ Error saving config: %v", err)
	}

	// Update in-memory config
	al.cfg.Agents.Defaults.RestrictToWorkspace = newValue

	// Propagate to ExecTool of all registered agents
	if al.registry != nil {
		for _, agentID := range al.registry.ListAgentIDs() {
			agent, ok := al.registry.GetAgent(agentID)
			if !ok || agent == nil || agent.Tools == nil {
				continue
			}
			if execTool, exists := agent.Tools.Get("exec"); exists {
				if et, ok := execTool.(*tools.ExecTool); ok {
					et.SetRestrictToWorkspace(newValue)
				}
			}
		}
	}

	// Audit log
	al.auditor.LogSecurityEvent("system", sessionKey, "restrict_to_workspace_"+action,
		fmt.Sprintf("User %s restrict_to_workspace via fast-path", actionDescription),
		fmt.Sprintf("Channel: %s, ChatID: %s, OldValue: %v, NewValue: %v",
			channel, chatID, currentValue, newValue))

	if newValue {
		return "✅ **restrict_to_workspace ACTIVATED**\n\n" +
			"🔒 Agent now operates ONLY within workspace directory\n" +
			"🛡️ Security level: MAXIMUM"
	}
	return "⚠️ **restrict_to_workspace DEACTIVATED**\n\n" +
		"🚨 **WARNING**: Agent can now access ANY system file\n" +
		"⚠️ Use with EXTREME caution\n" +
		"🔒 Run `/restrict_to_workspace activate` to restore security"
}

// isUserAuthorized checks if the chatID is in the allow_from list for the given channel.
// CLI and empty channels are always authorized.
func (al *AgentLoop) isUserAuthorized(channel, chatID string) bool {
	if channel == "cli" || channel == "" {
		return true
	}
	var allowFrom []string
	switch channel {
	case "telegram":
		allowFrom = []string(al.cfg.Channels.Telegram.AllowFrom)
	case "discord":
		allowFrom = []string(al.cfg.Channels.Discord.AllowFrom)
	case "slack":
		allowFrom = []string(al.cfg.Channels.Slack.AllowFrom)
	case "whatsapp":
		allowFrom = []string(al.cfg.Channels.WhatsApp.AllowFrom)
	case "wecom":
		allowFrom = []string(al.cfg.Channels.WeCom.AllowFrom)
	case "wecom_app":
		allowFrom = []string(al.cfg.Channels.WeComApp.AllowFrom)
	case "dingtalk":
		allowFrom = []string(al.cfg.Channels.DingTalk.AllowFrom)
	case "line":
		allowFrom = []string(al.cfg.Channels.LINE.AllowFrom)
	case "onebot":
		allowFrom = []string(al.cfg.Channels.OneBot.AllowFrom)
	default:
		return false
	}
	return containsStr(allowFrom, chatID)
}

// containsStr reports whether value is in slice.
func containsStr(slice []string, value string) bool {
	for _, s := range slice {
		if s == value {
			return true
		}
	}
	return false
}

func (al *AgentLoop) handleExplicitToolInvoke(
	ctx context.Context,
	agent *AgentInstance,
	channel, chatID, content string,
) (string, bool) {
	toolName, args, ok := parseExplicitToolInvocation(content)
	if !ok {
		return "", false
	}
	if _, exists := agent.Tools.Get(toolName); !exists {
		return fmt.Sprintf("Tool no encontrada: %s", toolName), true
	}
	result := agent.Tools.ExecuteWithContext(ctx, toolName, args, channel, chatID, nil)
	if result == nil {
		return "La tool no devolvió resultado.", true
	}
	if result.ForUser != "" && !result.Silent {
		return result.ForUser, true
	}
	return result.ForLLM, true
}

func parseExplicitToolInvocation(content string) (string, map[string]any, bool) {
	trimmed := strings.TrimSpace(content)
	if trimmed == "" {
		return "", nil, false
	}

	lowered := strings.ToLower(trimmed)
	var start int
	switch {
	case strings.HasPrefix(lowered, "usa "):
		start = len("usa ")
	case strings.HasPrefix(lowered, "use "):
		start = len("use ")
	default:
		return "", nil, false
	}

	rest := strings.TrimSpace(trimmed[start:])
	if rest == "" {
		return "", nil, false
	}

	toolName, tail, found := strings.Cut(rest, " ")
	if !found {
		toolName = strings.TrimSpace(rest)
		return toolName, map[string]any{}, toolName != ""
	}
	toolName = strings.TrimSpace(toolName)
	if toolName == "" {
		return "", nil, false
	}

	tail = strings.TrimSpace(tail)
	tailLower := strings.ToLower(tail)
	for _, marker := range []string{" con ", " with "} {
		if idx := strings.Index(tailLower, marker); idx >= 0 {
			tail = strings.TrimSpace(tail[idx+len(marker):])
			break
		}
	}
	if tail == "" {
		return toolName, map[string]any{}, true
	}

	args := make(map[string]any)
	for _, part := range splitArgsCSV(tail) {
		if part == "" {
			continue
		}
		k, v, ok := strings.Cut(part, "=")
		if !ok {
			continue
		}
		key := strings.TrimSpace(k)
		val := parseToolArgValue(strings.TrimSpace(v))
		if key != "" {
			args[key] = val
		}
	}

	return toolName, args, true
}

func splitArgsCSV(s string) []string {
	var parts []string
	var current strings.Builder
	var quote rune

	for _, r := range s {
		switch {
		case (r == '\'' || r == '"') && quote == 0:
			quote = r
			current.WriteRune(r)
		case quote != 0 && r == quote:
			quote = 0
			current.WriteRune(r)
		case r == ',' && quote == 0:
			parts = append(parts, strings.TrimSpace(current.String()))
			current.Reset()
		case r == ' ' && quote == 0:
			// Espacio fuera de comillas: también separa arguments
			if current.Len() > 0 {
				parts = append(parts, strings.TrimSpace(current.String()))
				current.Reset()
			}
		default:
			current.WriteRune(r)
		}
	}
	if current.Len() > 0 {
		parts = append(parts, strings.TrimSpace(current.String()))
	}
	return parts
}

func parseToolArgValue(v string) any {
	v = strings.TrimSpace(v)
	if len(v) >= 2 {
		if (v[0] == '\'' && v[len(v)-1] == '\'') || (v[0] == '"' && v[len(v)-1] == '"') {
			return v[1 : len(v)-1]
		}
	}
	lv := strings.ToLower(v)
	if lv == "true" {
		return true
	}
	if lv == "false" {
		return false
	}
	if n, err := strconv.Atoi(v); err == nil {
		return n
	}
	return v
}

// extractPeer extracts the routing peer from inbound message metadata.
func extractPeer(msg bus.InboundMessage) *routing.RoutePeer {
	peerKind := msg.Metadata["peer_kind"]
	if peerKind == "" {
		return nil
	}
	peerID := msg.Metadata["peer_id"]
	if peerID == "" {
		if peerKind == "direct" {
			peerID = msg.SenderID
		} else {
			peerID = msg.ChatID
		}
	}
	return &routing.RoutePeer{Kind: peerKind, ID: peerID}
}

// extractParentPeer extracts the parent peer (reply-to) from inbound message metadata.
func extractParentPeer(msg bus.InboundMessage) *routing.RoutePeer {
	parentKind := msg.Metadata["parent_peer_kind"]
	parentID := msg.Metadata["parent_peer_id"]
	if parentKind == "" || parentID == "" {
		return nil
	}
	return &routing.RoutePeer{Kind: parentKind, ID: parentID}
}

// filterSensitiveData redacts API keys and secrets from the content if the security filter is enabled.
func (al *AgentLoop) filterSensitiveData(content string) string {
	if al.cfg == nil || !al.cfg.Security.FilterSensitiveData {
		return content
	}

	secrets := al.cfg.GetSensitiveValues()
	for _, secret := range secrets {
		if secret != "" && len(secret) > 3 { // Prevent replacing tiny strings just in case
			content = strings.ReplaceAll(content, secret, "[REDACTED]")
		}
	}
	return content
}

// extractMediaFromSession scans the session messages for tool results that contain media paths.
func (al *AgentLoop) extractMediaFromSession(agent *AgentInstance, sessionKey string) []string {
	var paths []string
	messages := agent.Sessions.GetHistory(sessionKey)
	for _, msg := range messages {
		if msg.Role == "tool" && msg.ToolCallID != "" && len(msg.MediaPaths) > 0 {
			for _, p := range msg.MediaPaths {
				if _, err := os.Stat(p); err == nil {
					paths = append(paths, p)
				}
			}
		}
	}
	return paths
}
