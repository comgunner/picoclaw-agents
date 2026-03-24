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
	"sync"
	"time"

	"github.com/comgunner/picoclaw/pkg/bus"
	"github.com/comgunner/picoclaw/pkg/config"
	"github.com/comgunner/picoclaw/pkg/providers"
)

type SubagentTask struct {
	ID            string
	Task          string
	Label         string
	AgentID       string
	ExplicitModel string
	OriginChannel string
	OriginChatID  string
	Status        string
	Result        string
	Created       int64
}

type SubagentManager struct {
	tasks          map[string]*SubagentTask
	mu             sync.RWMutex
	provider       providers.LLMProvider
	defaultModel   string
	bus            *bus.MessageBus
	workspace      string
	tools          *ToolRegistry
	maxIterations  int
	maxTokens      int
	temperature    float64
	hasMaxTokens   bool
	hasTemperature bool
	nextID         int
	cfg            *config.Config
}

func NewSubagentManager(
	provider providers.LLMProvider,
	defaultModel, workspace string,
	bus *bus.MessageBus,
	cfg *config.Config,
) *SubagentManager {
	return &SubagentManager{
		tasks:         make(map[string]*SubagentTask),
		provider:      provider,
		defaultModel:  defaultModel,
		bus:           bus,
		workspace:     workspace,
		tools:         NewToolRegistry(),
		maxIterations: 10,
		nextID:        1,
		cfg:           cfg,
	}
}

// SetLLMOptions sets max tokens and temperature for subagent LLM calls.
func (sm *SubagentManager) SetLLMOptions(maxTokens int, temperature float64) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.maxTokens = maxTokens
	sm.hasMaxTokens = true
	sm.temperature = temperature
	sm.hasTemperature = true
}

// SetTools sets the tool registry for subagent execution.
// If not set, subagent will have access to the provided tools.
func (sm *SubagentManager) SetTools(tools *ToolRegistry) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.tools = tools
}

// RegisterTool registers a tool for subagent execution.
func (sm *SubagentManager) RegisterTool(tool Tool) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.tools.Register(tool)
}

func (sm *SubagentManager) Spawn(
	ctx context.Context,
	task, label, agentID, explicitModel, originChannel, originChatID string,
	callback AsyncCallback,
) (string, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	taskID := fmt.Sprintf("subagent-%d", sm.nextID)
	sm.nextID++

	subagentTask := &SubagentTask{
		ID:            taskID,
		Task:          task,
		Label:         label,
		AgentID:       agentID,
		ExplicitModel: explicitModel,
		OriginChannel: originChannel,
		OriginChatID:  originChatID,
		Status:        "running",
		Created:       time.Now().UnixMilli(),
	}
	sm.tasks[taskID] = subagentTask

	// Start task in background with context cancellation support
	go sm.runTask(ctx, subagentTask, callback)

	if label != "" {
		return fmt.Sprintf("Spawned subagent '%s' for task: %s", label, task), nil
	}
	return fmt.Sprintf("Spawned subagent for task: %s", task), nil
}

func (sm *SubagentManager) runTask(ctx context.Context, task *SubagentTask, callback AsyncCallback) {
	task.Status = "running"
	task.Created = time.Now().UnixMilli()

	// Build system prompt for subagent
	systemPrompt := `Eres un subagente de PicoClaw 🦞. Tu objetivo es completar la tarea de forma autónoma.

## REGLAS CRÍTICAS DE EJECUCIÓN:
1. **Flujo de Contenido**: Usa siempre 'text_script_create' -> 'image_gen_create(script_path=...)'.
2. **Social Media (Publicación)**: NUNCA publiques directamente en redes sociales sin mostrar el borrador al usuario primero.
   - Antes de usar 'facebook_post' o 'x_post_tweet', DEBES enviar un mensaje al usuario con la herramienta 'message' incluyendo el texto, la imagen (media) y botones.
3. **Resumen Final**: Indica siempre las rutas de los archivos generados.`

	messages := []providers.Message{
		{
			Role:    "system",
			Content: systemPrompt,
		},
		{
			Role:    "user",
			Content: task.Task,
		},
	}

	// Check if context is already canceled before starting
	select {
	case <-ctx.Done():
		sm.mu.Lock()
		task.Status = "canceled"
		task.Result = "Task canceled before execution"
		sm.mu.Unlock()
		return
	default:
	}

	// Run tool loop with access to tools
	sm.mu.RLock()
	tools := sm.tools
	maxIter := sm.maxIterations
	maxTokens := sm.maxTokens
	temperature := sm.temperature
	hasMaxTokens := sm.hasMaxTokens
	hasTemperature := sm.hasTemperature
	cfg := sm.cfg
	defaultModel := sm.defaultModel
	sm.mu.RUnlock()

	// Resolve model: explicit > agent config > defaults
	modelToUse := task.ExplicitModel
	if modelToUse == "" && task.AgentID != "" && cfg != nil {
		for i := range cfg.Agents.List {
			if cfg.Agents.List[i].ID == task.AgentID {
				if cfg.Agents.List[i].Model != nil && cfg.Agents.List[i].Model.Primary != "" {
					modelToUse = cfg.Agents.List[i].Model.Primary
				}
				break
			}
		}
	}
	if modelToUse == "" {
		modelToUse = defaultModel
	}

	var llmOptions map[string]any
	if hasMaxTokens || hasTemperature {
		llmOptions = map[string]any{}
		if hasMaxTokens {
			llmOptions["max_tokens"] = maxTokens
		}
		if hasTemperature {
			llmOptions["temperature"] = temperature
		}
	}

	loopResult, err := RunToolLoop(ctx, ToolLoopConfig{
		Provider:      sm.provider,
		Model:         modelToUse,
		Tools:         tools,
		MaxIterations: maxIter,
		LLMOptions:    llmOptions,
	}, messages, task.OriginChannel, task.OriginChatID)

	sm.mu.Lock()
	var result *ToolResult
	defer func() {
		sm.mu.Unlock()
		// Call callback if provided and result is set
		if callback != nil && result != nil {
			callback(ctx, result)
		}
	}()

	if err != nil {
		task.Status = "failed"
		task.Result = fmt.Sprintf("Error: %v", err)
		// Check if it was canceled
		if ctx.Err() != nil {
			task.Status = "canceled"
			task.Result = "Task canceled during execution"
		}
		result = &ToolResult{
			ForLLM:  task.Result,
			ForUser: "",
			Silent:  false,
			IsError: true,
			Async:   false,
			Err:     err,
		}
	} else {
		task.Status = "completed"
		task.Result = loopResult.Content
		result = &ToolResult{
			ForLLM: fmt.Sprintf(
				"Subagent '%s' completed (iterations: %d): %s",
				task.Label,
				loopResult.Iterations,
				loopResult.Content,
			),
			ForUser: loopResult.Content,
			Silent:  false,
			IsError: false,
			Async:   false,
		}
	}

	// Send announce message back to main agent
	if sm.bus != nil {
		announceContent := fmt.Sprintf("Task '%s' completed.\n\nResult:\n%s", task.Label, task.Result)
		sm.bus.PublishInbound(bus.InboundMessage{
			Channel:  "system",
			SenderID: fmt.Sprintf("subagent:%s", task.ID),
			// Format: "original_channel:original_chat_id" for routing back
			ChatID:  fmt.Sprintf("%s:%s", task.OriginChannel, task.OriginChatID),
			Content: announceContent,
		})
	}
}

func (sm *SubagentManager) GetTask(taskID string) (*SubagentTask, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	task, ok := sm.tasks[taskID]
	return task, ok
}

func (sm *SubagentManager) ListTasks() []*SubagentTask {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	tasks := make([]*SubagentTask, 0, len(sm.tasks))
	for _, task := range sm.tasks {
		tasks = append(tasks, task)
	}
	return tasks
}

// SubagentTool executes a subagent task synchronously and returns the result.
// Unlike SpawnTool which runs tasks asynchronously, SubagentTool waits for completion
// and returns the result directly in the ToolResult.
type SubagentTool struct {
	manager       *SubagentManager
	originChannel string
	originChatID  string
}

func NewSubagentTool(manager *SubagentManager) *SubagentTool {
	return &SubagentTool{
		manager:       manager,
		originChannel: "cli",
		originChatID:  "direct",
	}
}

func (t *SubagentTool) Name() string {
	return "subagent"
}

func (t *SubagentTool) Description() string {
	return "Execute a specialized subagent task synchronously. You MUST provide a clear 'task' description. Use this for complex sequences like 'create a script and then an image'. The result will be returned to you."
}

func (t *SubagentTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"task": map[string]any{
				"type":        "string",
				"description": "The task for subagent to complete. This is MANDATORY.",
			},
			"agent_id": map[string]any{
				"type":        "string",
				"description": "Optional target agent ID to delegate the task to (e.g. 'script_writer', 'image_creator', 'social_manager')",
			},
			"label": map[string]any{
				"type":        "string",
				"description": "Optional short label for the task (for display)",
			},
		},
		"required": []string{"task"},
	}
}

func (t *SubagentTool) SetContext(channel, chatID string) {
	t.originChannel = channel
	t.originChatID = chatID
}

func (t *SubagentTool) Execute(ctx context.Context, args map[string]any) *ToolResult {
	task, ok := args["task"].(string)
	if !ok || task == "" {
		return ErrorResult(
			"🦞 'task' is required and must be a string. Example: subagent(task='escribe un guion')",
		).WithError(fmt.Errorf("task parameter is missing or invalid"))
	}

	agentID, _ := args["agent_id"].(string)
	label, _ := args["label"].(string)

	if t.manager == nil {
		return ErrorResult("Subagent manager not configured").WithError(fmt.Errorf("manager is nil"))
	}

	// Build messages for subagent
	systemPrompt := "Eres un subagente de PicoClaw 🦞."
	if agentID != "" {
		systemPrompt = fmt.Sprintf("Eres el subagente especializado '%s' de PicoClaw.", agentID)
	}
	systemPrompt += ` Tu objetivo es completar tareas de forma autónoma usando herramientas. 

## REGLAS DE ORO:
1. **Coherencia**: Pasa el 'Path' de 'text_script_create' a 'image_gen_create'.
2. **Aprobación Social**: NUNCA publiques en Facebook o X sin mostrar el borrador primero usando la herramienta 'message'. Es OBLIGATORIO que el usuario vea el contenido antes de la acción final de postear.`

	messages := []providers.Message{
		{
			Role:    "system",
			Content: systemPrompt,
		},
		{
			Role:    "user",
			Content: task,
		},
	}

	// Use RunToolLoop to execute with tools (same as async SpawnTool)
	sm := t.manager
	sm.mu.RLock()
	tools := sm.tools
	maxIter := sm.maxIterations
	maxTokens := sm.maxTokens
	temperature := sm.temperature
	hasMaxTokens := sm.hasMaxTokens
	hasTemperature := sm.hasTemperature
	sm.mu.RUnlock()

	var llmOptions map[string]any
	if hasMaxTokens || hasTemperature {
		llmOptions = map[string]any{}
		if hasMaxTokens {
			llmOptions["max_tokens"] = maxTokens
		}
		if hasTemperature {
			llmOptions["temperature"] = temperature
		}
	}

	loopResult, err := RunToolLoop(ctx, ToolLoopConfig{
		Provider:      sm.provider,
		Model:         sm.defaultModel,
		Tools:         tools,
		MaxIterations: maxIter,
		LLMOptions:    llmOptions,
	}, messages, t.originChannel, t.originChatID)
	if err != nil {
		return ErrorResult(fmt.Sprintf("Subagent execution failed: %v", err)).WithError(err)
	}

	// ForUser: Brief summary for user (truncated if too long)
	userContent := loopResult.Content
	maxUserLen := 500
	if len(userContent) > maxUserLen {
		userContent = userContent[:maxUserLen] + "..."
	}

	// ForLLM: Full execution details
	labelStr := label
	if labelStr == "" {
		labelStr = "(unnamed)"
	}
	llmContent := fmt.Sprintf("Subagent task completed:\nLabel: %s\nIterations: %d\nResult: %s",
		labelStr, loopResult.Iterations, loopResult.Content)

	return &ToolResult{
		ForLLM:  llmContent,
		ForUser: userContent,
		Silent:  false,
		IsError: false,
		Async:   false,
	}
}
