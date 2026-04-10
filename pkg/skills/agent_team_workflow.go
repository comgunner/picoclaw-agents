// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)

package skills

import (
	"strings"
)

// AgentTeamWorkflowSkill implements native skill for multi-agent team orchestration.
// Based on PicoClaw's native multi-agent architecture and Claude Code Agent Team patterns.
type AgentTeamWorkflowSkill struct {
	workspace string
}

// NewAgentTeamWorkflowSkill creates a new AgentTeamWorkflowSkill instance.
func NewAgentTeamWorkflowSkill(workspace string) *AgentTeamWorkflowSkill {
	return &AgentTeamWorkflowSkill{
		workspace: workspace,
	}
}

// Name returns the skill identifier name.
func (a *AgentTeamWorkflowSkill) Name() string {
	return "agent_team_workflow"
}

// Description returns a brief description of the skill.
func (a *AgentTeamWorkflowSkill) Description() string {
	return "Multi-Agent Team Orchestrator - Organize optimal agent teams for any task by reading config.json and applying resource rules."
}

// GetInstructions returns the complete team orchestration guidelines.
func (a *AgentTeamWorkflowSkill) GetInstructions() string {
	return agentTeamWorkflowInstructions
}

// GetAntiPatterns returns common multi-agent anti-patterns.
func (a *AgentTeamWorkflowSkill) GetAntiPatterns() string {
	return agentTeamWorkflowAntiPatterns
}

// GetExamples returns concrete agent team examples.
func (a *AgentTeamWorkflowSkill) GetExamples() string {
	return agentTeamWorkflowExamples
}

// BuildSkillContext returns the complete skill context for prompt injection.
func (a *AgentTeamWorkflowSkill) BuildSkillContext() string {
	var parts []string

	parts = append(parts, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	parts = append(parts, "🚀 NATIVE SKILL: Agent Team Workflow Orchestrator")
	parts = append(parts, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	parts = append(parts, "")
	parts = append(parts, "**ROLE:** Multi-Agent Team Orchestrator")
	parts = append(parts, "")
	parts = append(
		parts,
		"**OBJECTIVE:** Organize optimal agent teams for any task type by reading config.json and applying resource rules (CPU, RAM, concurrency).",
	)
	parts = append(parts, "")
	parts = append(parts, a.GetInstructions())
	parts = append(parts, "")
	parts = append(parts, a.GetAntiPatterns())
	parts = append(parts, "")
	parts = append(parts, a.GetExamples())

	return strings.Join(parts, "\n")
}

// BuildSummary returns an XML summary for compact context injection.
func (a *AgentTeamWorkflowSkill) BuildSummary() string {
	return `<skill name="agent_team_workflow" type="native">
  <purpose>Multi-Agent Team Orchestrator</purpose>
  <pattern>Use for organizing agent teams, spawn optimization, resource management</pattern>
  <config>Read ~/.picoclaw/config.json for agents.list[]</config>
  <resources>CPU, RAM, max_concurrent, max_spawn_depth, max_children</resources>
  <patterns>Dev, Content, Image, Social, Research, General teams</patterns>
</skill>`
}

// ============================================================================
// DOCUMENTATION CONSTANTS
// ============================================================================

const agentTeamWorkflowInstructions = `## ROLE & OBJECTIVE

**Role:** Multi-Agent Team Orchestrator

**Objective:** Organize optimal agent teams for any task type by reading config.json and applying resource rules (CPU, RAM, concurrency).

## READING CONFIG.JSON

### Location
` + bt + `~/.picoclaw/config.json` + bt + `

### Structure
` + bt + bt + bt + `json
{
  "agents": {
    "defaults": {
      "workspace": "~/.picoclaw/workspace",
      "restrict_to_workspace": true,
      "model": "deepseek-chat",
      "max_tokens": 8192,
      "max_tool_iterations": 20
    },
    "list": [
      {
        "id": "project_manager",
        "default": true,
        "name": "Project Manager",
        "model": "deepseek-chat",
        "subagents": {
          "allow_agents": ["*", "general_worker"],
          "max_spawn_depth": 3,
          "max_children_per_agent": 5
        }
      },
      {
        "id": "senior_dev",
        "name": "Senior Developer",
        "model": "deepseek-chat",
        "subagents": {
          "allow_agents": ["qa_specialist", "junior_fixer"],
          "max_spawn_depth": 2
        }
      }
    ]
  }
}
` + bt + bt + bt + `

### Internal Methods (Go Code)

Instead of reading JSON manually, use internal registry methods:

` + bt + bt + bt + `go
// Get specific agent
agent, ok := registry.GetAgent("project_manager")

// Check if parent can spawn child
canSpawn := registry.CanSpawnSubagent("project_manager", "senior_dev")

// List all available agents
agentIDs := registry.ListAgentIDs()

// Get default agent (usually project_manager)
defaultAgent := registry.GetDefaultAgent()
` + bt + bt + bt + `

## RESOURCE RULES

### CPU Limits
- **Per-Agent CPU Budget:** Each agent should use ≤25% CPU during active reasoning
- **Total CPU Budget:** Sum of all active agents ≤80% CPU
- **Throttling:** If CPU >80%, pause lower-priority agents

### RAM Limits
- **Per-Agent RAM:** Each agent context ≤500MB
- **Total RAM:** Sum of all agent contexts ≤2GB
- **GC Triggers:** Force garbage collection when total >1.5GB

### Concurrency Limits
- **max_concurrent:** Maximum agents running simultaneously (default: 5)
- **Semaphore Pattern:** Use buffered channel to limit concurrency
- **Queue Excess:** Agents beyond limit wait in queue

` + bt + bt + bt + `go
// Example semaphore pattern
maxConcurrent := 5
semaphore := make(chan struct{}, maxConcurrent)

// Before spawning agent
semaphore <- struct{}{}  // Acquire slot
defer func() { <-semaphore }()  // Release slot

// Spawn agent
go spawnSubagent(...)
` + bt + bt + bt + `

## SPAWN OPTIMIZATION

### allow_agents[] Whitelist

Each agent can define which subagents it's allowed to spawn:

` + bt + bt + bt + `json
{
  "id": "project_manager",
  "subagents": {
    "allow_agents": ["*", "general_worker"]
  }
}
` + bt + bt + bt + `

**Patterns:**
- ` + bt + `"*"` + bt + ` - Wildcard, allows spawning any agent
- ` + bt + `"agent_id"` + bt + ` - Specific agent ID
- Empty/missing - No subagents allowed

### max_spawn_depth (Hierarchy Depth)

Controls how deep the agent hierarchy can go:

` + bt + bt + bt + `
Level 0: project_manager (main agent)
  ├─ Level 1: senior_dev (spawned by PM)
  │    ├─ Level 2: qa_specialist (spawned by senior_dev)
  │    │    └─ Level 3: junior_fixer (spawned by QA)
  │    └─ Level 2: junior_fixer
  └─ Level 1: general_worker
` + bt + bt + bt + `

**Typical Values:**
- project_manager: 3 (can spawn grandchildren)
- senior_dev: 2 (can spawn children)
- workers: 0 or 1 (leaf nodes, no spawning)

### max_children_per_agent (Parallelism)

Controls how many subagents an agent can spawn in parallel:

` + bt + bt + bt + `json
{
  "max_children_per_agent": 5
}
` + bt + bt + bt + `

**Benefits:**
- Prevents resource exhaustion
- Ensures fair scheduling
- Easier debugging (limited parallelism)

### max_concurrent (Global Limit)

Total number of agents running simultaneously across entire system:

` + bt + bt + bt + `json
{
  "max_concurrent": 5
}
` + bt + bt + bt + `

**Implementation:**
- Global semaphore across all agents
- Shared state in AgentRegistry
- First-come, first-served queuing

## TEAM PATTERNS (Multi-Purpose)

### 1. Development Team

**Use Case:** Software development, coding, debugging

**Team Structure:**
` + bt + bt + bt + `
project_manager (default)
├─ senior_dev
│   ├─ qa_specialist
│   └─ junior_fixer
└─ general_worker
` + bt + bt + bt + `

**Task Assignment:**
- **PM:** Requirements, architecture, coordination
- **Senior Dev:** Core logic, complex algorithms, code review
- **QA:** Testing, bug detection, GitHub ops
- **Junior:** Boilerplate, documentation, simple fixes
- **General:** Research, file operations, deployment

### 2. Content Team

**Use Case:** Writing, editing, content creation

**Team Structure:**
` + bt + bt + bt + `
project_manager (default)
├─ script_writer
├─ general_worker (research)
└─ image_creator (for visuals)
` + bt + bt + bt + `

**Task Assignment:**
- **PM:** Outline, coordination, quality control
- **Writer:** Draft creation, storytelling
- **General:** Research, fact-checking
- **Image:** Diagrams, illustrations

### 3. Image Team

**Use Case:** Image generation, editing, publishing

**Team Structure:**
` + bt + bt + bt + `
project_manager (default)
├─ image_creator
├─ image_creator (parallel for batch)
└─ social_manager (for publishing)
` + bt + bt + bt + `

**Task Assignment:**
- **PM:** Prompt engineering, quality review
- **Image Creator:** Generation, variations
- **Social Manager:** Publishing, scheduling

### 4. Social Media Team

**Use Case:** Social media management, posting, engagement

**Team Structure:**
` + bt + bt + bt + `
project_manager (default)
├─ social_manager
├─ content_writer
└─ image_creator
` + bt + bt + bt + `

**Task Assignment:**
- **PM:** Strategy, calendar, analytics
- **Social Manager:** Posting, scheduling, engagement
- **Writer:** Caption creation, hashtag research
- **Image:** Post visuals, stories

### 5. Research Team

**Use Case:** Information gathering, analysis, summarization

**Team Structure:**
` + bt + bt + bt + `
project_manager (default)
├─ general_worker (parallel researchers)
├─ general_worker
└─ general_worker
` + bt + bt + bt + `

**Task Assignment:**
- **PM:** Research questions, synthesis, reporting
- **General Workers:** Parallel research threads, fact-checking

### 6. General Team (Default)

**Use Case:** Any task without specific requirements

**Team Structure:**
` + bt + bt + bt + `
project_manager (default)
└─ general_worker
` + bt + bt + bt + `

**Task Assignment:**
- **PM:** Complex reasoning, coordination
- **General:** Execution, file operations, research

## LOAD BALANCING

### Round-Robin Across Models

When multiple agents use different models, distribute load:

` + bt + bt + bt + `go
// Global counter for round-robin
var rrCounter atomic.Uint64

func selectModel(agentID string) string {
    models := []string{"deepseek-chat", "gpt-4", "claude-sonnet"}
    index := rrCounter.Add(1) % uint64(len(models))
    return models[index]
}
` + bt + bt + bt + `

### Fallback Chains

Each agent can have primary + fallback models:

` + bt + bt + bt + `json
{
  "model": {
    "primary": "deepseek-chat",
    "fallbacks": ["gpt-4", "claude-sonnet"]
  }
}
` + bt + bt + bt + `

**Fallback Logic:**
1. Try primary model
2. If rate limited → fallback 1
3. If error → fallback 2
4. Track failures per model

### Model Selection by Task Type

| Task Type | Recommended Model | Why |
|-----------|------------------|-----|
| Coding | deepseek-chat | Excellent reasoning, cost-effective |
| Writing | claude-sonnet | Natural language, creativity |
| Research | gpt-4 | Broad knowledge, accuracy |
| Image Prompts | gpt-4-vision | Visual understanding |
| Quick Tasks | deepseek-chat | Fast, cheap |

## INTERNAL METHODS (Go Code)

### AgentRegistry Methods

These methods are available in the Go codebase:

**GetAgent:**
` + bt + bt + bt + `go
agent, ok := registry.GetAgent("project_manager")
if ok {
    // Use agent.Workspace, agent.Model, agent.Subagents
}
` + bt + bt + bt + `

**CanSpawnSubagent:**
` + bt + bt + bt + `go
canSpawn := registry.CanSpawnSubagent("project_manager", "senior_dev")
if canSpawn {
    // Proceed with spawn
} else {
    // Reject spawn request
}
` + bt + bt + bt + `

**ListAgentIDs:**
` + bt + bt + bt + `go
agentIDs := registry.ListAgentIDs()
// Returns: ["project_manager", "senior_dev", "qa_specialist", ...]
` + bt + bt + bt + `

**GetDefaultAgent:**
` + bt + bt + bt + `go
defaultAgent := registry.GetDefaultAgent()
// Usually returns project_manager
` + bt + bt + bt + `

### Subagent Spawn Flow

` + bt + bt + bt + `go
func spawnSubagent(parentID, targetID, task string) error {
    // 1. Check if parent can spawn target
    if !registry.CanSpawnSubagent(parentID, targetID) {
        return fmt.Errorf("agent %s cannot spawn %s", parentID, targetID)
    }

    // 2. Check depth limit
    parentDepth := getAgentDepth(parentID)
    if parentDepth >= parent.Subagents.MaxSpawnDepth {
        return fmt.Errorf("max spawn depth reached")
    }

    // 3. Acquire semaphore slot
    select {
    case semaphore <- struct{}{}:
        defer func() { <-semaphore }()
    case <-time.After(30 * time.Second):
        return fmt.Errorf("timeout waiting for agent slot")
    }

    // 4. Spawn agent
    go runAgent(targetID, task)

    return nil
}
` + bt + bt + bt + `

## SECURITY & CONSTRAINTS

### restrict_to_workspace

All file operations limited to agent's workspace:

` + bt + bt + bt + `json
{
  "agents": {
    "defaults": {
      "restrict_to_workspace": true
    }
  }
}
` + bt + bt + bt + `

**Enforcement:**
- File reads/writes checked against workspace path
- Symlink attacks prevented
- Path traversal blocked (../)

### max_tool_iterations

Limit tool execution loops to prevent infinite loops:

` + bt + bt + bt + `json
{
  "agents": {
    "defaults": {
      "max_tool_iterations": 20
    }
  }
}
` + bt + bt + bt + `

**Behavior:**
- Each tool call counts as 1 iteration
- After 20 iterations → force response
- Prevents runaway tool loops

### Token Budget Management

Each agent has token budget:

` + bt + bt + bt + `json
{
  "agents": {
    "defaults": {
      "max_tokens": 8192
    }
  }
}
` + bt + bt + bt + `

**Budget Allocation:**
- Context: ~4000 tokens
- Response: ~2000 tokens
- Tool outputs: ~2000 tokens
- Buffer: ~192 tokens

### Session Isolation (dm_scope)

Each channel-peer pair has isolated session:

` + bt + bt + bt + `json
{
  "session": {
    "dm_scope": "per-channel-peer"
  }
}
` + bt + bt + bt + `

**Benefits:**
- No cross-chat contamination
- Independent conversation history
- Privacy between users
`

const agentTeamWorkflowAntiPatterns = `## TEAM ORGANIZATION ANTI-PATTERNS

### ❌ Spawning Without Config Check
` + bt + bt + bt + `go
// BAD: Spawn without checking config
spawnSubagent("unknown_agent", "task")

// GOOD: Check registry first
agent, ok := registry.GetAgent("unknown_agent")
if !ok {
    return fmt.Errorf("agent not configured")
}
spawnSubagent("unknown_agent", "task")
` + bt + bt + bt + `

### ❌ Ignoring allow_agents Whitelist
` + bt + bt + bt + `go
// BAD: Spawn any agent
spawnSubagent("any_agent", "task")

// GOOD: Check whitelist
if !registry.CanSpawnSubagent(parentID, targetID) {
    return fmt.Errorf("spawn not allowed")
}
` + bt + bt + bt + `

### ❌ Exceeding max_spawn_depth
` + bt + bt + bt + `
Level 0: PM
  ├─ Level 1: Senior Dev
  │    ├─ Level 2: QA
  │    │    ├─ Level 3: Junior  # ❌ Exceeds depth=3
  │    │    └─ Level 3: Intern  # ❌ Exceeds depth=3
` + bt + bt + bt + `

### ❌ No Concurrency Control
` + bt + bt + bt + `go
// BAD: Spawn unlimited agents
for i := 0; i < 100; i++ {
    go spawnSubagent("worker", "task")
}

// GOOD: Use semaphore
semaphore := make(chan struct{}, 5)
for i := 0; i < 100; i++ {
    semaphore <- struct{}{}
    go func() {
        defer func() { <-semaphore }()
        spawnSubagent("worker", "task")
    }()
}
` + bt + bt + bt + `

## RESOURCE ANTI-PATTERNS

### ❌ No CPU Monitoring
` + bt + bt + bt + `go
// BAD: Spawn without CPU check
spawnSubagent("agent", "task")
// System CPU at 100%, agents slow

// GOOD: Check CPU before spawn
if getCPUUsage() > 80 {
    queueTask("agent", "task")  // Queue instead
} else {
    spawnSubagent("agent", "task")
}
` + bt + bt + bt + `

### ❌ Memory Leaks
` + bt + bt + bt + `go
// BAD: Never clean up contexts
for {
    ctx := context.Background()
    runAgent(ctx)
    // Context never released
}

// GOOD: Proper cleanup
for {
    ctx, cancel := context.WithTimeout(parent, 5*time.Minute)
    runAgent(ctx)
    cancel()  // Release resources
}
` + bt + bt + bt + `

### ❌ Ignoring max_concurrent
` + bt + bt + bt + `json
// Config says max_concurrent: 5
{
  "subagents": {"max_concurrent": 5}
}

// BAD: Spawn 20 agents
for i := 0; i < 20; i++ {
    go spawnSubagent("worker", "task")
}

// GOOD: Respect limit
semaphore := make(chan struct{}, 5)
for i := 0; i < 20; i++ {
    semaphore <- struct{}{}
    go func() {
        defer func() { <-semaphore }()
        spawnSubagent("worker", "task")
    }()
}
` + bt + bt + bt + `

## COMMUNICATION ANTI-PATTERNS

### ❌ Cross-Session Contamination
` + bt + bt + bt + `go
// BAD: Share history across sessions
globalHistory = append(globalHistory, message)

// GOOD: Isolate by dm_scope
sessionKey := fmt.Sprintf("%s-%s", channel, chatID)
sessionHistory[sessionKey] = append(sessionHistory[sessionKey], message)
` + bt + bt + bt + `

### ❌ No Task Locks
` + bt + bt + bt + `go
// BAD: Multiple agents edit same file
go agent1.EditFile("config.json")
go agent2.EditFile("config.json")
// File corruption!

// GOOD: Use task locks
lock := tasklock.New("config.json")
if lock.TryAcquire() {
    defer lock.Release()
    EditFile("config.json")
} else {
    // Wait or skip
}
` + bt + bt + bt + `

## MODEL SELECTION ANTI-PATTERNS

### ❌ Using Wrong Model for Task
` + bt + bt + bt + `
// BAD: Use expensive model for simple task
Task: "What's 2+2?"
Model: gpt-4  # Overkill, expensive

// GOOD: Match model to task
Task: "What's 2+2?"
Model: deepseek-chat  # Fast, cheap
` + bt + bt + bt + `

### ❌ No Fallback Handling
` + bt + bt + bt + `go
// BAD: Single model, no fallback
response := callModel(primary)
// If primary fails → error

// GOOD: Fallback chain
response, err := callWithFallback(primary, fallbacks)
if err != nil {
    // All models failed, handle gracefully
}
` + bt + bt + bt + `

## SECURITY ANTI-PATTERNS

### ❌ Bypassing restrict_to_workspace
` + bt + bt + bt + `go
// BAD: Allow arbitrary file access
filePath := userInput  // "/etc/passwd"
readFile(filePath)

// GOOD: Enforce workspace
filePath := filepath.Join(workspace, userInput)
if !strings.HasPrefix(filePath, workspace) {
    return fmt.Errorf("path outside workspace")
}
readFile(filePath)
` + bt + bt + bt + `

### ❌ Infinite Tool Loops
` + bt + bt + bt + `go
// BAD: No iteration limit
for {
    toolResult := callTool()
    if toolResult.NeedMoreInfo {
        continue  // Infinite loop!
    }
}

// GOOD: Limit iterations
for i := 0; i < maxToolIterations; i++ {
    toolResult := callTool()
    if !toolResult.NeedMoreInfo {
        break
    }
}
if i == maxToolIterations {
    return fmt.Errorf("max iterations reached")
}
` + bt + bt + bt + `
`

const agentTeamWorkflowExamples = `## EXAMPLE 1: DEVELOPMENT TEAM WORKFLOW

**Request:** "Fix bug in authentication module"

**Expert Response:**

### Team Organization

Based on config.json, use Development Team pattern:

` + bt + bt + bt + `
project_manager (default, depth=0)
└─ senior_dev (depth=1, can spawn QA + Junior)
    ├─ qa_specialist (depth=2)
    └─ junior_fixer (depth=2)
` + bt + bt + bt + `

### Spawn Flow

` + bt + bt + bt + `go
// 1. PM receives task
task := "Fix bug in authentication module"

// 2. PM spawns senior_dev
if registry.CanSpawnSubagent("project_manager", "senior_dev") {
    spawnSubagent("senior_dev", task)
}

// 3. Senior dev spawns QA and Junior
if registry.CanSpawnSubagent("senior_dev", "qa_specialist") {
    spawnSubagent("qa_specialist", "Test authentication fix")
}
if registry.CanSpawnSubagent("senior_dev", "junior_fixer") {
    spawnSubagent("junior_fixer", "Apply authentication patch")
}
` + bt + bt + bt + `

### Resource Allocation

- **CPU:** 4 agents × 25% = 100% (within 80% limit, throttle if needed)
- **RAM:** 4 agents × 500MB = 2GB (at limit, monitor GC)
- **Concurrency:** 4 agents < max_concurrent(5) ✅

### Task Assignment

1. **PM:** Analyze bug report, assign to senior_dev
2. **Senior Dev:** Debug, implement fix, review code
3. **QA:** Write tests, verify fix, run regression
4. **Junior:** Update documentation, commit changes

---

## EXAMPLE 2: CONTENT CREATION TEAM

**Request:** "Write blog post about AI agents"

**Expert Response:**

### Team Organization

Use Content Team pattern:

` + bt + bt + bt + `
project_manager (default, depth=0)
├─ script_writer (depth=1)
├─ general_worker (depth=1, research)
└─ image_creator (depth=1, visuals)
` + bt + bt + bt + `

### Spawn Flow

` + bt + bt + bt + `go
// PM spawns content team in parallel
agents := []string{"script_writer", "general_worker", "image_creator"}
semaphore := make(chan struct{}, 3)

for _, agentID := range agents {
    semaphore <- struct{}{}
    go func(id string) {
        defer func() { <-semaphore }()
        spawnSubagent(id, "Create content for AI agents blog post")
    }(agentID)
}
` + bt + bt + bt + `

### Resource Allocation

- **CPU:** 3 agents × 25% = 75% (within 80% limit) ✅
- **RAM:** 3 agents × 500MB = 1.5GB (within 2GB limit) ✅
- **Concurrency:** 3 agents < max_concurrent(5) ✅

### Task Assignment

1. **PM:** Outline, coordinate, quality review
2. **Writer:** Draft blog post, storytelling
3. **General:** Research AI agent trends, fact-check
4. **Image:** Create diagrams, illustrations

---

## EXAMPLE 3: IMAGE GENERATION BATCH

**Request:** "Generate 10 images of landscapes"

**Expert Response:**

### Team Organization

Use Image Team pattern with batch queue:

` + bt + bt + bt + `
project_manager (default, depth=0)
├─ image_creator #1 (depth=1, images 1-3)
├─ image_creator #2 (depth=1, images 4-6)
├─ image_creator #3 (depth=1, images 7-9)
└─ image_creator #4 (depth=1, image 10)
` + bt + bt + bt + `

### Spawn Flow with Queue

` + bt + bt + bt + `go
// PM spawns image creators with batching
batchSize := 3
images := 10
semaphore := make(chan struct{}, 4)  // max_concurrent

for i := 0; i < images; i += batchSize {
    // Wait for semaphore slot
    semaphore <- struct{}{}

    go func(start int) {
        defer func() { <-semaphore }()

        end := min(start+batchSize, images)
        task := fmt.Sprintf("Generate images %d to %d of landscapes", start, end)
        spawnSubagent("image_creator", task)
    }(i)
}
` + bt + bt + bt + `

### Resource Allocation

- **CPU:** 4 agents × 25% = 100% (throttle to 80%)
- **RAM:** 4 agents × 500MB = 2GB (at limit, force GC)
- **Concurrency:** 4 agents < max_concurrent(5) ✅

### Queue Integration

Use queue_batch skill for fire-and-forget:

` + bt + bt + bt + `go
// Generate batch ID
batchID := batch_id(prefix="IMG_GEN")

// Enqueue task
queue(action="enqueue", task_type="IMAGE_GEN", payload={
    "script": "scripts/batch_gen.py",
    "args": ["--count", "10", "--prompt", "landscapes"],
    "batch_id": batchID
})

// Notify user
"🔥 Started generating 10 images. ID: " + batchID
"Te aviso cuando termine."
` + bt + bt + bt + `

---

## EXAMPLE 4: RESEARCH TEAM (PARALLEL)

**Request:** "Research top 10 AI companies"

**Expert Response:**

### Team Organization

Use Research Team pattern with parallel workers:

` + bt + bt + bt + `
project_manager (default, depth=0)
├─ general_worker #1 (depth=1, companies 1-3)
├─ general_worker #2 (depth=1, companies 4-6)
├─ general_worker #3 (depth=1, companies 7-9)
└─ general_worker #4 (depth=1, company 10 + synthesis)
` + bt + bt + bt + `

### Spawn Flow

` + bt + bt + bt + `go
// PM spawns parallel researchers
companies := []string{"OpenAI", "Anthropic", "Google DeepMind", ...}
batchSize := 3
semaphore := make(chan struct{}, 4)

for i := 0; i < len(companies); i += batchSize {
    semaphore <- struct{}{}

    go func(start int) {
        defer func() { <-semaphore }()

        end := min(start+batchSize, len(companies))
        batch := companies[start:end]
        task := fmt.Sprintf("Research companies: %v", batch)
        spawnSubagent("general_worker", task)
    }(i)
}
` + bt + bt + bt + `

### Resource Allocation

- **CPU:** 4 agents × 25% = 100% (throttle to 80%)
- **RAM:** 4 agents × 500MB = 2GB (at limit)
- **Concurrency:** 4 agents < max_concurrent(5) ✅

### Synthesis

PM collects results from all workers:

` + bt + bt + bt + `go
// Collect results
var results []string
for _, worker := range workers {
    results = append(results, worker.Result)
}

// Synthesize final report
report := synthesize(results)
return report
` + bt + bt + bt + `

---

## EXAMPLE 5: CONFIG-BASED AGENT SELECTION

**Request:** "How do I know which agents are available?"

**Expert Response:**

### Reading Available Agents

Use internal registry methods (no JSON parsing needed):

` + bt + bt + bt + `go
// List all configured agents
agentIDs := registry.ListAgentIDs()
// Returns: ["project_manager", "senior_dev", "qa_specialist", ...]

// Get specific agent
agent, ok := registry.GetAgent("senior_dev")
if ok {
    fmt.Printf("Agent: %s\n", agent.Name)
    fmt.Printf("Model: %s\n", agent.Model)
    fmt.Printf("Workspace: %s\n", agent.Workspace)
    fmt.Printf("Can spawn: %v\n", agent.Subagents.AllowAgents)
}

// Check spawn permission
canSpawn := registry.CanSpawnSubagent("project_manager", "senior_dev")
if canSpawn {
    fmt.Println("PM can spawn Senior Dev")
}

// Get default agent (usually project_manager)
defaultAgent := registry.GetDefaultAgent()
fmt.Printf("Default agent: %s\n", defaultAgent.ID)
` + bt + bt + bt + `

### Example Output

` + bt + bt + bt + `
Available Agents:
  ✓ project_manager (default) - deepseek-chat
  ✓ senior_dev - deepseek-chat
  ✓ qa_specialist - deepseek-chat
  ✓ junior_fixer - deepseek-chat
  ✓ general_worker - deepseek-chat
  ✓ script_writer - deepseek-chat
  ✓ image_creator - deepseek-chat
  ✓ social_manager - deepseek-chat

Spawn Permissions:
  project_manager → [* (wildcard), general_worker]
  senior_dev → [qa_specialist, junior_fixer]
  qa_specialist → []
  junior_fixer → []
` + bt + bt + bt + `

---

## EXAMPLE 6: RESOURCE-AWARE SPAWNING

**Request:** "Spawn agents but respect system resources"

**Expert Response:**

### Resource Monitoring Before Spawn

` + bt + bt + bt + `go
func spawnWithResourceCheck(agentID, task string) error {
    // 1. Check CPU
    cpuUsage := getCPUUsage()
    if cpuUsage > 80 {
        logger.Warn("CPU too high, queuing task")
        queueTask(agentID, task)
        return nil  // Don't spawn now
    }

    // 2. Check RAM
    ramUsage := getRAMUsage()
    if ramUsage > 1500 {  // MB
        logger.Warn("RAM high, forcing GC")
        runtime.GC()
        // Recheck after GC
        time.Sleep(1 * time.Second)
        ramUsage = getRAMUsage()
        if ramUsage > 1800 {
            queueTask(agentID, task)
            return nil
        }
    }

    // 3. Check concurrency
    select {
    case semaphore <- struct{}{}:
        // Slot acquired
    case <-time.After(30 * time.Second):
        return fmt.Errorf("timeout waiting for agent slot")
    }

    // 4. Spawn agent
    go func() {
        defer func() { <-semaphore }()
        runAgent(agentID, task)
    }()

    return nil
}
` + bt + bt + bt + `

### Usage

` + bt + bt + bt + `go
// Instead of direct spawn
spawnSubagent("senior_dev", "Fix bug")

// Use resource-aware spawn
spawnWithResourceCheck("senior_dev", "Fix bug")
` + bt + bt + bt + `

### Benefits

- **Prevents OOM:** Monitors RAM, forces GC
- **Prevents CPU Saturation:** Queues tasks when CPU >80%
- **Fair Scheduling:** Semaphore ensures max_concurrent respected
- **Graceful Degradation:** Queues instead of failing
`
