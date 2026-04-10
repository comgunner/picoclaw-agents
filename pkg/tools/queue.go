package tools

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/comgunner/picoclaw/pkg/bus"
	"github.com/comgunner/picoclaw/pkg/utils"
)

// QueueStatus represents the current state of a background task.
type QueueStatus string

const (
	StatusPending    QueueStatus = "pending"
	StatusProcessing QueueStatus = "processing"
	StatusCompleted  QueueStatus = "completed"
	StatusFailed     QueueStatus = "failed"
)

type BatchTask struct {
	ID        string
	Type      string // "IMA_GEN", "TEXT_GEN", "POST_BUNDLE"
	Status    QueueStatus
	Payload   map[string]any
	Result    *ToolResult
	CreatedAt time.Time
	UpdatedAt time.Time
}

// QueueManager handles background tasks without LLM intervention for steps.
type QueueManager struct {
	tasks map[string]*BatchTask
	mu    sync.RWMutex
	bus   *bus.MessageBus
	tools *ToolRegistry
}

var (
	globalQueue *QueueManager
	onceQueue   sync.Once
)

func GetQueueManager(b *bus.MessageBus, r *ToolRegistry) *QueueManager {
	onceQueue.Do(func() {
		globalQueue = &QueueManager{
			tasks: make(map[string]*BatchTask),
			bus:   b,
			tools: r,
		}
	})
	return globalQueue
}

// GetGlobalQueueManager devuelve la instancia ya inicializada.
func GetGlobalQueueManager() *QueueManager {
	return globalQueue
}

func (q *QueueManager) SetTools(r *ToolRegistry) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.tools = r
}

func (q *QueueManager) AddTask(taskType string, payload map[string]any) string {
	q.mu.Lock()
	defer q.mu.Unlock()

	id := utils.GenerateBatchID(taskType)
	task := &BatchTask{
		ID:        id,
		Type:      taskType,
		Status:    StatusPending,
		Payload:   payload,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	q.tasks[id] = task
	return id
}

func (q *QueueManager) GetTask(id string) (*BatchTask, bool) {
	q.mu.RLock()
	defer q.mu.RUnlock()
	t, ok := q.tasks[id]
	return t, ok
}

func (q *QueueManager) UpdateStatus(id string, status QueueStatus, res *ToolResult) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if t, ok := q.tasks[id]; ok {
		t.Status = status
		t.Result = res
		t.UpdatedAt = time.Now()
	}
}

// IsCancelled checks if a task has been canceled.
func (q *QueueManager) IsCancelled(id string) bool {
	q.mu.RLock()
	defer q.mu.RUnlock()
	if t, ok := q.tasks[id]; ok {
		return t.Status == StatusFailed || t.Status == "canceled"
	}
	return false
}

// QueueTool is the user-facing tool to check or interact with the queue.
type QueueTool struct {
	manager *QueueManager
}

func NewQueueTool(manager *QueueManager) *QueueTool {
	return &QueueTool{manager: manager}
}

func (t *QueueTool) Name() string {
	return "queue"
}

func (t *QueueTool) Description() string {
	return "Consulta el estado de tareas en segundo plano o interactúa con IDs específicos (#IMA_GEN...)."
}

func (t *QueueTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"action": map[string]any{
				"type":        "string",
				"description": "Acción: 'list', 'status', 'process'",
				"enum":        []string{"list", "status", "process"},
			},
			"task_id": map[string]any{
				"type":        "string",
				"description": "ID de la tarea (ej: #IMA_GEN_02_03_26_1500)",
			},
		},
	}
}

func (t *QueueTool) Execute(ctx context.Context, args map[string]any) *ToolResult {
	action, _ := args["action"].(string)
	taskID, _ := args["task_id"].(string)

	if action == "" && taskID != "" {
		action = "status"
	}

	switch action {
	case "list":
		t.manager.mu.RLock()
		defer t.manager.mu.RUnlock()
		if len(t.manager.tasks) == 0 {
			return UserResult("📭 La cola está vacía.")
		}
		var sb strings.Builder
		sb.WriteString("📋 **Tareas en Cola:**\n")
		for _, task := range t.manager.tasks {
			sb.WriteString(fmt.Sprintf("- %s [%s] (%s)\n", task.ID, task.Status, task.CreatedAt.Format("15:04")))
		}
		return UserResult(sb.String())

	case "status":
		task, ok := t.manager.GetTask(taskID)
		if !ok {
			return ErrorResult(fmt.Sprintf("Tarea %s no encontrada.", taskID))
		}

		statusEmoji := "⏳"
		if task.Status == StatusCompleted {
			statusEmoji = "✅"
		}
		if task.Status == StatusFailed {
			statusEmoji = "❌"
		}
		if task.Status == StatusProcessing {
			statusEmoji = "⚙️"
		}

		msg := fmt.Sprintf("%s **Estado de %s:** %s", statusEmoji, task.ID, task.Status)
		if task.Result != nil && task.Result.ForUser != "" {
			msg += "\n\n" + task.Result.ForUser
		}

		res := UserResult(msg)
		if task.Result != nil {
			res.Buttons = task.Result.Buttons
		}
		return res

	default:
		return UserResult("Uso: queue(action='list') o queue(task_id='#...')")
	}
}
