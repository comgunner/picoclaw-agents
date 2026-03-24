// PicoClaw - Ultra-lightweight personal AI agent
// Agent Management - Unit Tests
// License: MIT

package management_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/comgunner/picoclaw/pkg/agents/management"
	"github.com/comgunner/picoclaw/pkg/config"
)

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func makeConfig(agents ...config.AgentConfig) *config.Config {
	return &config.Config{
		Agents: config.AgentsConfig{List: agents},
	}
}

func agentCfg(id, name string, isDefault bool) config.AgentConfig {
	return config.AgentConfig{
		ID:      id,
		Name:    name,
		Default: isDefault,
	}
}

// ---------------------------------------------------------------------------
// AgentRegistry
// ---------------------------------------------------------------------------

func TestAgentRegistry_GetAgent_Found(t *testing.T) {
	cfg := makeConfig(agentCfg("alpha", "Alpha Agent", false))
	reg := management.NewAgentRegistry(cfg)

	agent, found := reg.GetAgent("alpha")
	require.True(t, found)
	assert.Equal(t, "alpha", agent.ID)
	assert.Equal(t, "Alpha Agent", agent.Name)
}

func TestAgentRegistry_GetAgent_NotFound(t *testing.T) {
	cfg := makeConfig(agentCfg("alpha", "Alpha", false))
	reg := management.NewAgentRegistry(cfg)

	_, found := reg.GetAgent("ghost")
	assert.False(t, found)
}

func TestAgentRegistry_GetAgent_NilConfig(t *testing.T) {
	reg := management.NewAgentRegistry(nil)
	_, found := reg.GetAgent("x")
	assert.False(t, found)
}

func TestAgentRegistry_ListAgentIDs(t *testing.T) {
	cfg := makeConfig(
		agentCfg("alpha", "Alpha", false),
		agentCfg("beta", "Beta", false),
	)
	reg := management.NewAgentRegistry(cfg)

	ids := reg.ListAgentIDs()
	assert.ElementsMatch(t, []string{"alpha", "beta"}, ids)
}

func TestAgentRegistry_GetDefaultAgent_ExplicitDefault(t *testing.T) {
	cfg := makeConfig(
		agentCfg("alpha", "Alpha", false),
		agentCfg("beta", "Beta Agent", true),
	)
	reg := management.NewAgentRegistry(cfg)

	def := reg.GetDefaultAgent()
	require.NotNil(t, def)
	assert.Equal(t, "beta", def.ID)
}

func TestAgentRegistry_GetDefaultAgent_FallbackToFirst(t *testing.T) {
	cfg := makeConfig(
		agentCfg("first", "First", false),
		agentCfg("second", "Second", false),
	)
	reg := management.NewAgentRegistry(cfg)

	def := reg.GetDefaultAgent()
	require.NotNil(t, def)
	assert.Equal(t, "first", def.ID)
}

func TestAgentRegistry_GetDefaultAgent_NoAgents(t *testing.T) {
	reg := management.NewAgentRegistry(makeConfig())
	assert.Nil(t, reg.GetDefaultAgent())
}

func TestAgentRegistry_CanSpawnSubagent_OpenPolicy(t *testing.T) {
	cfg := makeConfig(agentCfg("parent", "Parent", false), agentCfg("child", "Child", false))
	reg := management.NewAgentRegistry(cfg)

	assert.True(t, reg.CanSpawnSubagent("parent", "child"))
}

func TestAgentRegistry_CanSpawnSubagent_RestrictedAllowlist(t *testing.T) {
	cfg := makeConfig(
		agentCfg("mgr", "Manager", false),
		agentCfg("worker1", "Worker 1", false),
		agentCfg("worker2", "Worker 2", false),
	)
	reg := management.NewAgentRegistry(cfg)
	reg.SetSpawnPermission("mgr", []string{"worker1"})

	assert.True(t, reg.CanSpawnSubagent("mgr", "worker1"))
	assert.False(t, reg.CanSpawnSubagent("mgr", "worker2"))
}

func TestAgentRegistry_CanSpawnSubagent_EmptyAllowlist(t *testing.T) {
	cfg := makeConfig(agentCfg("root", "Root", false), agentCfg("child", "Child", false))
	reg := management.NewAgentRegistry(cfg)
	reg.SetSpawnPermission("root", []string{}) // explicitly empty → deny all

	assert.False(t, reg.CanSpawnSubagent("root", "child"))
}

// ---------------------------------------------------------------------------
// AgentInstance & InstanceRegistry
// ---------------------------------------------------------------------------

func TestAgentInstance_Lifecycle(t *testing.T) {
	inst := management.NewAgentInstance("agent1")
	assert.Equal(t, management.StatusIdle, inst.GetStatus())

	inst.UpdateStatus(management.StatusActive)
	assert.Equal(t, management.StatusActive, inst.GetStatus())
}

func TestAgentInstance_TaskManagement(t *testing.T) {
	inst := management.NewAgentInstance("agent1")

	task := management.TaskInfo{
		ID:        "task-1",
		Task:      "Do the thing",
		Status:    "processing",
		StartedAt: time.Now(),
	}
	inst.AddActiveTask(task)

	metrics := inst.GetMetrics()
	assert.Equal(t, 1, metrics.ActiveTasks)

	inst.CompleteTask("task-1", true)
	metrics = inst.GetMetrics()
	assert.Equal(t, 0, metrics.ActiveTasks)
	assert.Equal(t, 1, metrics.CompletedToday)
}

func TestAgentInstance_FailedTask(t *testing.T) {
	inst := management.NewAgentInstance("agent1")
	inst.AddActiveTask(management.TaskInfo{ID: "t1", Task: "fail me", StartedAt: time.Now()})
	inst.CompleteTask("t1", false)

	metrics := inst.GetMetrics()
	assert.Equal(t, 1, metrics.FailedToday)
	assert.Equal(t, 0, metrics.CompletedToday)
}

func TestInstanceRegistry_RegisterAndGet(t *testing.T) {
	reg := management.NewInstanceRegistry()

	inst1 := reg.Register("agent1")
	require.NotNil(t, inst1)

	// Idempotent: second call returns same instance
	inst2 := reg.Register("agent1")
	assert.Equal(t, inst1, inst2)

	got, ok := reg.Get("agent1")
	require.True(t, ok)
	assert.Equal(t, inst1, got)
}

func TestInstanceRegistry_Unregister(t *testing.T) {
	reg := management.NewInstanceRegistry()
	reg.Register("agent1")
	reg.Unregister("agent1")

	_, ok := reg.Get("agent1")
	assert.False(t, ok)
}

func TestInstanceRegistry_List(t *testing.T) {
	reg := management.NewInstanceRegistry()
	reg.Register("a")
	reg.Register("b")

	list := reg.List()
	assert.Len(t, list, 2)
}

// ---------------------------------------------------------------------------
// AgentMessageBus
// ---------------------------------------------------------------------------

func newMsg(id, from, to, typ string) management.AgentMessage {
	rawJSON := []byte(`"hello"`)
	return management.AgentMessage{
		ID:          id,
		SenderID:    from,
		RecipientID: to,
		MessageType: typ,
		Payload:     rawJSON,
		SentAt:      time.Now(),
	}
}

func TestAgentMessageBus_SendAndReceive(t *testing.T) {
	bus := management.NewAgentMessageBus()

	err := bus.Send(newMsg("m1", "alpha", "beta", "info"))
	require.NoError(t, err)

	msgs, err := bus.Receive("beta", nil)
	require.NoError(t, err)
	require.Len(t, msgs, 1)
	assert.Equal(t, "m1", msgs[0].ID)
}

func TestAgentMessageBus_ReceiveEmpty(t *testing.T) {
	bus := management.NewAgentMessageBus()
	msgs, err := bus.Receive("nobody", nil)
	require.NoError(t, err)
	assert.Empty(t, msgs)
}

func TestAgentMessageBus_TypeFilter(t *testing.T) {
	bus := management.NewAgentMessageBus()
	_ = bus.Send(newMsg("m1", "a", "b", "info"))
	_ = bus.Send(newMsg("m2", "a", "b", "task"))

	msgs, err := bus.Receive("b", []string{"info"})
	require.NoError(t, err)
	require.Len(t, msgs, 1)
	assert.Equal(t, "info", msgs[0].MessageType)
}

func TestAgentMessageBus_MarkAsRead(t *testing.T) {
	bus := management.NewAgentMessageBus()
	_ = bus.Send(newMsg("m1", "a", "b", "info"))

	err := bus.MarkAsRead("m1")
	require.NoError(t, err)

	err = bus.MarkAsRead("nonexistent")
	require.Error(t, err)
}

func TestAgentMessageBus_GetStats(t *testing.T) {
	bus := management.NewAgentMessageBus()
	_ = bus.Send(newMsg("m1", "a", "b", "info"))
	_ = bus.Send(newMsg("m2", "a", "b", "task"))

	stats := bus.GetStats()
	assert.Equal(t, 2, stats.TotalMessages)
	assert.Equal(t, 1, stats.ByType["info"])
	assert.Equal(t, 1, stats.ByType["task"])
}

// ---------------------------------------------------------------------------
// RateLimiter
// ---------------------------------------------------------------------------

func TestRateLimiter_AllowsUnlimitedTool(t *testing.T) {
	rl := management.NewRateLimiter()
	// agent_get has no configured rate limit → always allowed
	for i := 0; i < 100; i++ {
		assert.True(t, rl.Check("any_agent", "agent_get"))
	}
}

func TestRateLimiter_EnforcesLimit(t *testing.T) {
	rl := management.NewRateLimiter()
	// agent_broadcast has a limit of 5/min
	limit := management.RateLimits["agent_broadcast"]
	for i := 0; i < limit; i++ {
		require.True(t, rl.Check("agent1", "agent_broadcast"), "call %d should be allowed", i+1)
	}
	assert.False(t, rl.Check("agent1", "agent_broadcast"), "call beyond limit should be denied")
}

// ---------------------------------------------------------------------------
// IsToolAllowed
// ---------------------------------------------------------------------------

func TestIsToolAllowed_Wildcard(t *testing.T) {
	assert.True(t, management.IsToolAllowed("agent_get", "any_role"))
	assert.True(t, management.IsToolAllowed("agent_list", ""))
}

func TestIsToolAllowed_Restricted(t *testing.T) {
	assert.True(t, management.IsToolAllowed("agent_broadcast", "admin"))
	assert.False(t, management.IsToolAllowed("agent_broadcast", "guest"))
}

func TestIsToolAllowed_UnknownTool(t *testing.T) {
	assert.False(t, management.IsToolAllowed("not_a_tool", "admin"))
}
