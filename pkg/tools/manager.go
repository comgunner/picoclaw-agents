// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package tools

import (
	"sort"
	"strings"
	"sync"

	"github.com/comgunner/picoclaw/pkg/providers"
)

// ToolTier represents the priority tier of a tool.
type ToolTier int

const (
	// TierEssential: always included (file ops, exec, messaging)
	TierEssential ToolTier = iota
	// TierCommon: included for most conversations (spawn, web, queue)
	TierCommon
	// TierSpecialized: included only when context allows (binance, image_gen, etc.)
	TierSpecialized
)

// ToolCategory groups tools by domain.
type ToolCategory string

const (
	CategoryFileOps   ToolCategory = "file_ops"
	CategoryExec      ToolCategory = "exec"
	CategoryMessaging ToolCategory = "messaging"
	CategorySubagent  ToolCategory = "subagent"
	CategoryWeb       ToolCategory = "web"
	CategoryQueue     ToolCategory = "queue"
	CategoryBinance   ToolCategory = "binance"
	CategoryImageGen  ToolCategory = "image_gen"
	CategorySocial    ToolCategory = "social"
	CategorySystem    ToolCategory = "system"
	CategoryCron      ToolCategory = "cron"
	CategoryMemory    ToolCategory = "memory"
	CategoryVersion   ToolCategory = "version"
	CategoryOther     ToolCategory = "other"
)

// toolMetadata holds tier and category for a registered tool.
type toolMetadata struct {
	tool     Tool
	tier     ToolTier
	category ToolCategory
}

// ToolsManager manages tools with tiered selection for token optimization.
// Instead of sending all 67 tool definitions to the LLM (consuming ~10K tokens),
// it selects only the tools needed for the current context.
type ToolsManager struct {
	registry *ToolRegistry
	tools    map[string]toolMetadata
	mu       sync.RWMutex
}

// NewToolsManager creates a new ToolsManager wrapping an existing ToolRegistry.
func NewToolsManager(registry *ToolRegistry) *ToolsManager {
	return &ToolsManager{
		registry: registry,
		tools:    make(map[string]toolMetadata),
	}
}

// Register registers a tool with tier and category metadata.
func (tm *ToolsManager) Register(tool Tool, tier ToolTier, category ToolCategory) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	tm.tools[tool.Name()] = toolMetadata{
		tool:     tool,
		tier:     tier,
		category: category,
	}
}

// RegisterAll categorizes and registers all tools from the underlying registry.
// Tools not explicitly categorized default to TierSpecialized/CategoryOther.
func (tm *ToolsManager) RegisterAll() {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	// Define tier/category assignments by tool name
	assignments := map[string]struct {
		tier     ToolTier
		category ToolCategory
	}{
		// TierEssential - always included
		"read_file":  {TierEssential, CategoryFileOps},
		"write_file": {TierEssential, CategoryFileOps},
		"edit_file":  {TierEssential, CategoryFileOps},
		"list_dir":   {TierEssential, CategoryFileOps},
		"exec":       {TierEssential, CategoryExec},
		"message":    {TierEssential, CategoryMessaging},
		// TierCommon - included for most conversations
		"spawn":                 {TierCommon, CategorySubagent},
		"subagent":              {TierCommon, CategorySubagent},
		"subagent_list":         {TierCommon, CategorySubagent},
		"web":                   {TierCommon, CategoryWeb},
		"batch_id":              {TierCommon, CategoryQueue},
		"queue":                 {TierCommon, CategoryQueue},
		"memory_store":          {TierCommon, CategoryMemory},
		"workspace_maintenance": {TierCommon, CategorySystem},
		"context_status":        {TierCommon, CategorySystem},
		"self_diagnostics":      {TierCommon, CategorySystem},
		"version_control":       {TierCommon, CategoryVersion},
		"resource_monitor":      {TierCommon, CategorySystem},
		"system_diagnostics":    {TierCommon, CategorySystem},
		"config_manager":        {TierCommon, CategorySystem},
		"cron":                  {TierCommon, CategoryCron},
		// TierSpecialized - included only when context allows
		"binance":                 {TierSpecialized, CategoryBinance},
		"base_trader":             {TierSpecialized, CategoryBinance},
		"image_gen":               {TierSpecialized, CategoryImageGen},
		"image_gen_antigravity":   {TierSpecialized, CategoryImageGen},
		"image_gen_workflow":      {TierSpecialized, CategoryImageGen},
		"image_approval":          {TierSpecialized, CategoryImageGen},
		"social_media":            {TierSpecialized, CategorySocial},
		"social_manager":          {TierSpecialized, CategorySocial},
		"social_post_bundle":      {TierSpecialized, CategorySocial},
		"community_manager":       {TierSpecialized, CategorySocial},
		"discord_webhook":         {TierSpecialized, CategorySocial},
		"notion":                  {TierSpecialized, CategoryOther},
		"skills_install":          {TierSpecialized, CategoryOther},
		"skills_search":           {TierSpecialized, CategoryOther},
		"skills_sentinel":         {TierSpecialized, CategoryOther},
		"md_audit_tool":           {TierSpecialized, CategoryOther},
		"arch_lint_tool":          {TierSpecialized, CategoryOther},
		"bench_tool":              {TierSpecialized, CategoryOther},
		"reaper_tool":             {TierSpecialized, CategoryOther},
		"codegen":                 {TierSpecialized, CategoryOther},
		"toolloop":                {TierSpecialized, CategoryOther},
		"validate":                {TierSpecialized, CategoryOther},
		"text_script":             {TierSpecialized, CategoryOther},
		"text_script_antigravity": {TierSpecialized, CategoryOther},
		"text_approval":           {TierSpecialized, CategoryOther},
		"mcp_tool":                {TierSpecialized, CategoryOther},
		"spi":                     {TierSpecialized, CategoryOther},
		"i2c":                     {TierSpecialized, CategoryOther},
	}

	for _, name := range tm.registry.sortedToolNames() {
		tool, ok := tm.registry.Get(name)
		if !ok {
			continue
		}
		if a, ok := assignments[name]; ok {
			tm.tools[name] = toolMetadata{tool: tool, tier: a.tier, category: a.category}
		} else {
			tm.tools[name] = toolMetadata{tool: tool, tier: TierSpecialized, category: CategoryOther}
		}
	}
}

// SelectTools returns tool definitions filtered by maximum tier.
// tierLimit: TierEssential=only essential, TierCommon=essential+common, TierSpecialized=all.
func (tm *ToolsManager) SelectTools(tierLimit ToolTier) []providers.ToolDefinition {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	var names []string
	for name, meta := range tm.tools {
		if meta.tier <= tierLimit {
			names = append(names, name)
		}
	}
	sort.Strings(names)

	definitions := make([]providers.ToolDefinition, 0, len(names))
	for _, name := range names {
		meta := tm.tools[name]
		schema := ToolToSchema(meta.tool)
		fn, ok := schema["function"].(map[string]any)
		if !ok {
			continue
		}
		n, _ := fn["name"].(string)
		desc, _ := fn["description"].(string)
		params, _ := fn["parameters"].(map[string]any)
		definitions = append(definitions, providers.ToolDefinition{
			Type: "function",
			Function: providers.ToolFunctionDefinition{
				Name:        n,
				Description: desc,
				Parameters:  params,
			},
		})
	}
	return definitions
}

// SelectToolsForContext selects tools based on context window size.
// Large context (>32K): all tools
// Medium context (8K-32K): essential + common
// Small context (<8K): essential only
func (tm *ToolsManager) SelectToolsForContext(contextWindow int) []providers.ToolDefinition {
	switch {
	case contextWindow >= 32768:
		return tm.SelectTools(TierSpecialized)
	case contextWindow >= 8192:
		return tm.SelectTools(TierCommon)
	default:
		return tm.SelectTools(TierEssential)
	}
}

// SelectToolsForMessage analyzes the user message and includes specialized
// tools only if relevant keywords are detected.
func (tm *ToolsManager) SelectToolsForMessage(contextWindow int, userMessage string) []providers.ToolDefinition {
	// Start with base tier based on context window
	baseTier := TierCommon
	if contextWindow < 8192 {
		baseTier = TierEssential
	} else if contextWindow >= 32768 {
		baseTier = TierSpecialized
	}

	// If already at max, return all
	if baseTier == TierSpecialized {
		return tm.SelectTools(TierSpecialized)
	}

	// Check if message needs specialized tools
	lowerMsg := strings.ToLower(userMessage)
	needsBinance := strings.Contains(lowerMsg, "binance") || strings.Contains(lowerMsg, "trading") ||
		strings.Contains(lowerMsg, "crypto")
	needsImage := strings.Contains(lowerMsg, "imagen") || strings.Contains(lowerMsg, "image") ||
		strings.Contains(lowerMsg, "picture") ||
		strings.Contains(lowerMsg, "dibuja")
	needsSocial := strings.Contains(lowerMsg, "facebook") || strings.Contains(lowerMsg, "twitter") ||
		strings.Contains(lowerMsg, "redes") ||
		strings.Contains(lowerMsg, "social") ||
		strings.Contains(lowerMsg, "post")

	if needsBinance || needsImage || needsSocial {
		return tm.SelectTools(TierSpecialized)
	}

	return tm.SelectTools(baseTier)
}

// GetToolCount returns the number of tools that would be selected for a given tier.
func (tm *ToolsManager) GetToolCount(tierLimit ToolTier) int {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	count := 0
	for _, meta := range tm.tools {
		if meta.tier <= tierLimit {
			count++
		}
	}
	return count
}

// GetToolsByCategory returns tools grouped by category for a given tier limit.
func (tm *ToolsManager) GetToolsByCategory(tierLimit ToolTier) map[ToolCategory][]string {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	result := make(map[ToolCategory][]string)
	for name, meta := range tm.tools {
		if meta.tier <= tierLimit {
			result[meta.category] = append(result[meta.category], name)
		}
	}
	for cat := range result {
		sort.Strings(result[cat])
	}
	return result
}

// GetRegistry returns the underlying ToolRegistry for direct access.
func (tm *ToolsManager) GetRegistry() *ToolRegistry {
	return tm.registry
}
