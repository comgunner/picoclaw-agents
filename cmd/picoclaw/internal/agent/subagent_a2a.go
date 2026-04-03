// cmd/picoclaw/internal/agent/subagent_a2a.go
//
// A2A CLI Commands - Agent-to-Agent Communication CLI
// Based on concepts from icueth's picoclaw-agents fork (@icueth)
//
// Credits: @icueth (https://github.com/icueth)
// License: Same as base project (MIT)

package agent

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/comgunner/picoclaw/pkg/agent"
	"github.com/comgunner/picoclaw/pkg/mailbox"
)

// NewA2ACommand creates A2A orchestrator commands
func NewA2ACommand(orch *agent.A2AOrchestrator, deptRouter *agent.DepartmentRouter) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "a2a",
		Short: "A2A Orchestrator commands",
		Long:  "Commands for Agent-to-Agent communication and coordination",
	}

	cmd.AddCommand(
		newA2AStatusCommand(orch),
		newA2AMessageCommand(orch),
		newA2ADiscoveryCommand(orch),
		newA2ADepartmentCommand(deptRouter),
		newA2ATaskCommand(orch),
	)

	return cmd
}

// newA2AStatusCommand shows A2A orchestration status
func newA2AStatusCommand(orch *agent.A2AOrchestrator) *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show A2A orchestration status",
		RunE: func(cmd *cobra.Command, args []string) error {
			stats := orch.GetOrchestrationStats()

			fmt.Printf("=== A2A Orchestration Status ===\n\n")
			fmt.Printf("Agents: %d\n", stats["agents"])
			fmt.Printf("Message Log: %d\n", stats["message_log"])
			fmt.Printf("Context Size: %d\n\n", stats["context_size"])

			mbStats := stats["mailbox_stats"].(map[string]any)
			fmt.Printf("Mailbox Stats:\n")
			fmt.Printf("  Total Agents: %d\n", mbStats["total_agents"].(int))

			agents := mbStats["agents"].(map[string]any)
			if len(agents) > 0 {
				fmt.Printf("\n  Per-Agent Mailboxes:\n")
				for agentID, agentStats := range agents {
					s := agentStats.(map[string]any)
					fmt.Printf("    %s: size=%d, unread=%d\n",
						agentID,
						s["size"].(int),
						s["unread_count"].(int),
					)
				}
			}

			fmt.Printf("\nTimestamp: %s\n", time.Unix(stats["timestamp"].(int64), 0).Format(time.RFC3339))
			return nil
		},
	}
}

// newA2AMessageCommand sends A2A messages
func newA2AMessageCommand(orch *agent.A2AOrchestrator) *cobra.Command {
	var priority string
	var msgType string

	cmd := &cobra.Command{
		Use:   "message <from> <to> <content>",
		Short: "Send A2A message",
		Args:  cobra.MinimumNArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			from := args[0]
			to := args[1]
			content := args[2]

			// Parse priority
			p := mailbox.PriorityNormal
			switch priority {
			case "critical":
				p = mailbox.PriorityCritical
			case "high":
				p = mailbox.PriorityHigh
			case "low":
				p = mailbox.PriorityLow
			}

			// Parse message type
			t := mailbox.MessageTypeTask
			switch msgType {
			case "question":
				t = mailbox.MessageTypeQuestion
			case "status":
				t = mailbox.MessageTypeStatus
			case "broadcast":
				t = mailbox.MessageTypeBroadcast
			}

			if err := orch.SendMessage(from, to, t, p, content); err != nil {
				return fmt.Errorf("failed to send message: %w", err)
			}

			fmt.Printf("✓ Message sent: %s -> %s [%s, %s]\n", from, to, msgType, priority)
			return nil
		},
	}

	cmd.Flags().StringVarP(&priority, "priority", "p", "normal", "Priority: critical, high, normal, low")
	cmd.Flags().StringVarP(&msgType, "type", "t", "task", "Type: task, question, status, broadcast")

	return cmd
}

// newA2ADiscoveryCommand starts discovery phase
func newA2ADiscoveryCommand(orch *agent.A2AOrchestrator) *cobra.Command {
	return &cobra.Command{
		Use:   "discover",
		Short: "Start discovery phase",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := orch.WorkflowDiscovery(); err != nil {
				return fmt.Errorf("discovery failed: %w", err)
			}

			fmt.Println("✓ Discovery phase completed")
			fmt.Println("  All agents have broadcast their capabilities")
			return nil
		},
	}
}

// newA2ADepartmentCommand manages department model configuration
func newA2ADepartmentCommand(deptRouter *agent.DepartmentRouter) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "department",
		Short: "Department model configuration",
	}

	cmd.AddCommand(
		newDepartmentListCommand(deptRouter),
		newDepartmentAgentCommand(deptRouter),
	)

	return cmd
}

// newDepartmentListCommand lists departments and their models
func newDepartmentListCommand(deptRouter *agent.DepartmentRouter) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List departments and their models",
		RunE: func(cmd *cobra.Command, args []string) error {
			depts := deptRouter.ListDepartments()

			fmt.Printf("=== Department Models ===\n\n")
			fmt.Printf("%-25s %-20s\n", "Department", "Model")
			fmt.Printf("%-25s %-20s\n", "----------", "-----")

			for dept, model := range depts {
				fmt.Printf("%-25s %-20s\n", dept, model)
			}

			defaultModel := deptRouter.GetDefaultModel()
			fmt.Printf("\nDefault (Fallback): %s\n", defaultModel)
			fmt.Printf("Total departments: %d\n", len(depts))

			return nil
		},
	}
}

// newDepartmentAgentCommand shows agents in a department
func newDepartmentAgentCommand(deptRouter *agent.DepartmentRouter) *cobra.Command {
	return &cobra.Command{
		Use:   "agents <department>",
		Short: "List agents in a department",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			department := args[0]

			agents := deptRouter.GetDepartmentAgents(department)

			if len(agents) == 0 {
				fmt.Printf("No agents found in department: %s\n", department)
				return nil
			}

			fmt.Printf("=== Agents in Department: %s ===\n\n", department)
			for _, agentID := range agents {
				fmt.Printf("  - %s\n", agentID)
			}
			fmt.Printf("\nTotal: %d agents\n", len(agents))

			return nil
		},
	}
}

// newA2ATaskCommand manages A2A tasks
func newA2ATaskCommand(orch *agent.A2AOrchestrator) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "task",
		Short: "Task management",
	}

	cmd.AddCommand(
		newTaskAssignCommand(orch),
		newTaskCompleteCommand(orch),
	)

	return cmd
}

// newTaskAssignCommand assigns a task to an agent
func newTaskAssignCommand(orch *agent.A2AOrchestrator) *cobra.Command {
	var priority string

	cmd := &cobra.Command{
		Use:   "assign <from> <to> <task>",
		Short: "Assign task to agent",
		Args:  cobra.MinimumNArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			from := args[0]
			to := args[1]
			task := args[2]

			// Parse priority
			p := mailbox.PriorityNormal
			switch priority {
			case "critical":
				p = mailbox.PriorityCritical
			case "high":
				p = mailbox.PriorityHigh
			case "low":
				p = mailbox.PriorityLow
			}

			if err := orch.AssignTask(from, to, task, p); err != nil {
				return fmt.Errorf("failed to assign task: %w", err)
			}

			fmt.Printf("✓ Task assigned: %s -> %s [%s]\n", from, to, priority)
			fmt.Printf("  Task: %s\n", truncateString(task, 100))
			return nil
		},
	}

	cmd.Flags().StringVarP(&priority, "priority", "p", "high", "Priority: critical, high, normal, low")

	return cmd
}

// newTaskCompleteCommand reports task completion
func newTaskCompleteCommand(orch *agent.A2AOrchestrator) *cobra.Command {
	return &cobra.Command{
		Use:   "complete <agent> <task_key> <result>",
		Short: "Report task completion",
		Args:  cobra.MinimumNArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			agentID := args[0]
			taskKey := args[1]
			result := args[2]

			if err := orch.ReportTaskComplete(agentID, taskKey, result); err != nil {
				return fmt.Errorf("failed to report completion: %w", err)
			}

			fmt.Printf("✓ Task completed: %s\n", taskKey)
			fmt.Printf("  Result: %s\n", truncateString(result, 100))
			return nil
		},
	}
}

// Helper function to truncate strings
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
