// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package skills

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/comgunner/picoclaw/pkg/logger"
)

var (
	namePattern        = regexp.MustCompile(`^[a-zA-Z0-9]+(-[a-zA-Z0-9]+)*$`)
	reFrontmatter      = regexp.MustCompile(`(?s)^---(?:\r\n|\n|\r)(.*?)(?:\r\n|\n|\r)---`)
	reStripFrontmatter = regexp.MustCompile(`(?s)^---(?:\r\n|\n|\r)(.*?)(?:\r\n|\n|\r)---(?:\r\n|\n|\r)*`)
)

const (
	MaxNameLength        = 64
	MaxDescriptionLength = 1024
)

type SkillMetadata struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type SkillInfo struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	Source      string `json:"source"`
	Description string `json:"description"`
}

func (info SkillInfo) validate() error {
	var errs error
	if info.Name == "" {
		errs = errors.Join(errs, errors.New("name is required"))
	} else {
		if len(info.Name) > MaxNameLength {
			errs = errors.Join(errs, fmt.Errorf("name exceeds %d characters", MaxNameLength))
		}
		if !namePattern.MatchString(info.Name) {
			errs = errors.Join(errs, errors.New("name must be alphanumeric with hyphens"))
		}
	}

	if info.Description == "" {
		errs = errors.Join(errs, errors.New("description is required"))
	} else if len(info.Description) > MaxDescriptionLength {
		errs = errors.Join(errs, fmt.Errorf("description exceeds %d character", MaxDescriptionLength))
	}
	return errs
}

type SkillsLoader struct {
	workspace       string
	workspaceSkills string // workspace skills (project-level)
	globalSkills    string // global skills (~/.picoclaw/skills)
	builtinSkills   string // builtin skills
	embeddedFS      fs.FS  // embedded skills (//go:embed)
	includeNative   bool   // whether to include compiled-in native skills
}

func NewSkillsLoader(workspace string, globalSkills string, builtinSkills string) *SkillsLoader {
	return &SkillsLoader{
		workspace:       workspace,
		workspaceSkills: filepath.Join(workspace, "skills"),
		globalSkills:    globalSkills, // ~/.picoclaw/skills
		builtinSkills:   builtinSkills,
		embeddedFS:      GetEmbeddedSkillsFS(),
		includeNative:   true,
	}
}

// NewSkillsLoaderFiles creates a loader that only scans the file system
// (workspace, global, builtin) without embedded or native compiled-in skills.
func NewSkillsLoaderFiles(workspace string, globalSkills string, builtinSkills string) *SkillsLoader {
	return &SkillsLoader{
		workspace:       workspace,
		workspaceSkills: filepath.Join(workspace, "skills"),
		globalSkills:    globalSkills,
		builtinSkills:   builtinSkills,
		embeddedFS:      nil,
		includeNative:   false,
	}
}

func (sl *SkillsLoader) ListSkills() []SkillInfo {
	skills := make([]SkillInfo, 0)
	seen := make(map[string]bool)

	// Add native compiled-in skills first
	if sl.includeNative {
		nativeSkills := sl.listNativeSkills()
		for _, skill := range nativeSkills {
			seen[skill.Name] = true
			skills = append(skills, skill)
		}
	}

	addSkills := func(dir, source string) {
		if dir == "" {
			return
		}
		dirs, err := os.ReadDir(dir)
		if err != nil {
			return
		}
		for _, d := range dirs {
			if !d.IsDir() {
				continue
			}
			skillFile := filepath.Join(dir, d.Name(), "SKILL.md")
			if _, err := os.Stat(skillFile); err != nil {
				continue
			}
			info := SkillInfo{
				Name:   d.Name(),
				Path:   skillFile,
				Source: source,
			}
			metadata := sl.getSkillMetadata(skillFile)
			if metadata != nil {
				info.Description = metadata.Description
				info.Name = metadata.Name
			}
			if err := info.validate(); err != nil {
				slog.Warn("invalid skill from "+source, "name", info.Name, "error", err)
				continue
			}
			if seen[info.Name] {
				continue
			}
			seen[info.Name] = true
			skills = append(skills, info)
		}
	}

	// Priority: workspace > global > builtin > embedded
	addSkills(sl.workspaceSkills, "workspace")
	addSkills(sl.globalSkills, "global")
	addSkills(sl.builtinSkills, "builtin")

	// Add embedded skills (lowest priority)
	if sl.embeddedFS != nil {
		addEmbeddedSkills(sl.embeddedFS, "embedded", seen, &skills)
	}

	return skills
}

// addEmbeddedSkills adds skills from the embedded filesystem
func addEmbeddedSkills(fsys fs.FS, source string, seen map[string]bool, skills *[]SkillInfo) {
	// Read categories from embedded FS
	categories, err := fs.ReadDir(fsys, ".")
	if err != nil {
		return
	}

	for _, cat := range categories {
		if !cat.IsDir() {
			continue
		}

		skillDirs, err := fs.ReadDir(fsys, cat.Name())
		if err != nil {
			continue
		}

		for _, d := range skillDirs {
			if !d.IsDir() {
				continue
			}

			skillFile := cat.Name() + "/" + d.Name() + "/SKILL.md"
			content, err := fs.ReadFile(fsys, skillFile)
			if err != nil {
				continue
			}

			info := SkillInfo{
				Name:   d.Name(),
				Path:   "embedded://" + skillFile,
				Source: source,
			}

			metadata := parseSkillFrontmatter(string(content))
			if metadata != nil {
				if metadata["name"] != "" {
					info.Name = metadata["name"]
				}
				info.Description = metadata["description"]
			}

			if err := info.validate(); err != nil {
				slog.Warn("invalid skill from "+source, "name", info.Name, "error", err)
				continue
			}

			if seen[info.Name] {
				continue
			}

			seen[info.Name] = true
			*skills = append(*skills, info)
		}
	}
}

func (sl *SkillsLoader) LoadSkill(name string) (string, bool) {
	// 1. load from workspace skills first (project-level)
	if sl.workspaceSkills != "" {
		skillFile := filepath.Join(sl.workspaceSkills, name, "SKILL.md")
		if content, err := os.ReadFile(skillFile); err == nil {
			return sl.stripFrontmatter(string(content)), true
		}
	}

	// 2. then load from global skills (~/.picoclaw/skills)
	if sl.globalSkills != "" {
		skillFile := filepath.Join(sl.globalSkills, name, "SKILL.md")
		if content, err := os.ReadFile(skillFile); err == nil {
			return sl.stripFrontmatter(string(content)), true
		}
	}

	// 3. finally load from builtin skills
	if sl.builtinSkills != "" {
		skillFile := filepath.Join(sl.builtinSkills, name, "SKILL.md")
		if content, err := os.ReadFile(skillFile); err == nil {
			return sl.stripFrontmatter(string(content)), true
		}
	}

	// 4. fallback to embedded FS
	if sl.embeddedFS != nil {
		// Try direct path: {name}/SKILL.md
		skillFile := filepath.Join(name, "SKILL.md")
		if content, err := fs.ReadFile(sl.embeddedFS, skillFile); err == nil {
			return sl.stripFrontmatter(string(content)), true
		}

		// Try category paths: {category}/{name}/SKILL.md
		categories, _ := fs.ReadDir(sl.embeddedFS, ".")
		for _, cat := range categories {
			if !cat.IsDir() {
				continue
			}
			skillFile := filepath.Join(cat.Name(), name, "SKILL.md")
			if content, err := fs.ReadFile(sl.embeddedFS, skillFile); err == nil {
				return sl.stripFrontmatter(string(content)), true
			}
		}
	}

	return "", false
}

func (sl *SkillsLoader) LoadSkillsForContext(skillNames []string) string {
	if len(skillNames) == 0 {
		return ""
	}

	var parts []string
	for _, name := range skillNames {
		content, ok := sl.LoadSkill(name)
		if ok {
			parts = append(parts, fmt.Sprintf("### Skill: %s\n\n%s", name, content))
		}
	}

	return strings.Join(parts, "\n\n---\n\n")
}

func (sl *SkillsLoader) BuildSkillsSummary() string {
	allSkills := sl.ListSkills()
	if len(allSkills) == 0 {
		return ""
	}

	var lines []string
	lines = append(lines, "<skills>")
	for _, s := range allSkills {
		escapedName := escapeXML(s.Name)
		escapedDesc := escapeXML(s.Description)
		escapedPath := escapeXML(s.Path)

		lines = append(lines, fmt.Sprintf("  <skill>"))
		lines = append(lines, fmt.Sprintf("    <name>%s</name>", escapedName))
		lines = append(lines, fmt.Sprintf("    <description>%s</description>", escapedDesc))
		lines = append(lines, fmt.Sprintf("    <location>%s</location>", escapedPath))
		lines = append(lines, fmt.Sprintf("    <source>%s</source>", s.Source))
		lines = append(lines, "  </skill>")
	}
	lines = append(lines, "</skills>")

	return strings.Join(lines, "\n")
}

// ============================================================================
// NATIVE SKILLS REGISTRY
// ============================================================================

// nativeSkillsRegistry holds instances of all native skills compiled into the binary.
var nativeSkillsRegistry = struct {
	queueBatch        *QueueBatchSkill
	binanceMCP        *BinanceMCPSkill
	fullstackDev      *FullStackDeveloperSkill
	n8nWorkflow       *N8NWorkflowSkill
	agentTeamWorkflow *AgentTeamWorkflowSkill
	researcher        *ResearcherSkill       // Fase 0: registrar skill existente
	backendDev        *BackendDeveloperSkill // Fase 5: skills de rol engineering
	frontendDev       *FrontendDeveloperSkill
	devopsEng         *DevOpsEngineerSkill
	securityEng       *SecurityEngineerSkill
	qaEng             *QAEngineerSkill
	dataEng           *DataEngineerSkill
	mlEng             *MLEngineerSkill
	odooDev           *OdooDeveloperSkill // Odoo Developer skill
}{
	queueBatch:        nil, // Initialized on first use
	binanceMCP:        nil, // Initialized on first use
	fullstackDev:      nil, // Initialized on first use
	n8nWorkflow:       nil, // Initialized on first use
	agentTeamWorkflow: nil, // Initialized on first use
	researcher:        nil, // Initialized on first use
	backendDev:        nil, // Initialized on first use
	frontendDev:       nil, // Initialized on first use
	devopsEng:         nil, // Initialized on first use
	securityEng:       nil, // Initialized on first use
	qaEng:             nil, // Initialized on first use
	dataEng:           nil, // Initialized on first use
	mlEng:             nil, // Initialized on first use
	odooDev:           nil, // Initialized on first use
}

// GetQueueBatchSkill returns the singleton instance of QueueBatchSkill.
// Thread-safe lazy initialization.
func GetQueueBatchSkill(workspace string) *QueueBatchSkill {
	if nativeSkillsRegistry.queueBatch == nil {
		nativeSkillsRegistry.queueBatch = NewQueueBatchSkill(workspace)
	}
	return nativeSkillsRegistry.queueBatch
}

// GetBinanceMCPSkill returns the singleton instance of BinanceMCPSkill.
// Thread-safe lazy initialization.
func GetBinanceMCPSkill(workspace string) *BinanceMCPSkill {
	if nativeSkillsRegistry.binanceMCP == nil {
		nativeSkillsRegistry.binanceMCP = NewBinanceMCPSkill(workspace)
	}
	return nativeSkillsRegistry.binanceMCP
}

// GetFullStackDeveloperSkill returns the singleton instance of FullStackDeveloperSkill.
// Thread-safe lazy initialization.
func GetFullStackDeveloperSkill(workspace string) *FullStackDeveloperSkill {
	if nativeSkillsRegistry.fullstackDev == nil {
		nativeSkillsRegistry.fullstackDev = NewFullStackDeveloperSkill(workspace)
	}
	return nativeSkillsRegistry.fullstackDev
}

// GetN8NWorkflowSkill returns the singleton instance of N8NWorkflowSkill.
// Thread-safe lazy initialization.
func GetN8NWorkflowSkill(workspace string) *N8NWorkflowSkill {
	if nativeSkillsRegistry.n8nWorkflow == nil {
		nativeSkillsRegistry.n8nWorkflow = NewN8NWorkflowSkill(workspace)
	}
	return nativeSkillsRegistry.n8nWorkflow
}

// GetAgentTeamWorkflowSkill returns the singleton instance of AgentTeamWorkflowSkill.
// Thread-safe lazy initialization.
func GetAgentTeamWorkflowSkill(workspace string) *AgentTeamWorkflowSkill {
	if nativeSkillsRegistry.agentTeamWorkflow == nil {
		nativeSkillsRegistry.agentTeamWorkflow = NewAgentTeamWorkflowSkill(workspace)
	}
	return nativeSkillsRegistry.agentTeamWorkflow
}

// GetResearcherSkill returns the singleton instance of ResearcherSkill.
// Thread-safe lazy initialization.
func GetResearcherSkill(workspace string) *ResearcherSkill {
	if nativeSkillsRegistry.researcher == nil {
		nativeSkillsRegistry.researcher = NewResearcherSkill(workspace)
	}
	return nativeSkillsRegistry.researcher
}

// GetBackendDeveloperSkill returns the singleton instance of BackendDeveloperSkill.
// Thread-safe lazy initialization.
func GetBackendDeveloperSkill(workspace string) *BackendDeveloperSkill {
	if nativeSkillsRegistry.backendDev == nil {
		nativeSkillsRegistry.backendDev = NewBackendDeveloperSkill(workspace)
	}
	return nativeSkillsRegistry.backendDev
}

// GetFrontendDeveloperSkill returns the singleton instance of FrontendDeveloperSkill.
// Thread-safe lazy initialization.
func GetFrontendDeveloperSkill(workspace string) *FrontendDeveloperSkill {
	if nativeSkillsRegistry.frontendDev == nil {
		nativeSkillsRegistry.frontendDev = NewFrontendDeveloperSkill(workspace)
	}
	return nativeSkillsRegistry.frontendDev
}

// GetDevOpsEngineerSkill returns the singleton instance of DevOpsEngineerSkill.
// Thread-safe lazy initialization.
func GetDevOpsEngineerSkill(workspace string) *DevOpsEngineerSkill {
	if nativeSkillsRegistry.devopsEng == nil {
		nativeSkillsRegistry.devopsEng = NewDevOpsEngineerSkill(workspace)
	}
	return nativeSkillsRegistry.devopsEng
}

// GetSecurityEngineerSkill returns the singleton instance of SecurityEngineerSkill.
// Thread-safe lazy initialization.
func GetSecurityEngineerSkill(workspace string) *SecurityEngineerSkill {
	if nativeSkillsRegistry.securityEng == nil {
		nativeSkillsRegistry.securityEng = NewSecurityEngineerSkill(workspace)
	}
	return nativeSkillsRegistry.securityEng
}

// GetQAEngineerSkill returns the singleton instance of QAEngineerSkill.
// Thread-safe lazy initialization.
func GetQAEngineerSkill(workspace string) *QAEngineerSkill {
	if nativeSkillsRegistry.qaEng == nil {
		nativeSkillsRegistry.qaEng = NewQAEngineerSkill(workspace)
	}
	return nativeSkillsRegistry.qaEng
}

// GetDataEngineerSkill returns the singleton instance of DataEngineerSkill.
// Thread-safe lazy initialization.
func GetDataEngineerSkill(workspace string) *DataEngineerSkill {
	if nativeSkillsRegistry.dataEng == nil {
		nativeSkillsRegistry.dataEng = NewDataEngineerSkill(workspace)
	}
	return nativeSkillsRegistry.dataEng
}

// GetMLEngineerSkill returns the singleton instance of MLEngineerSkill.
// Thread-safe lazy initialization.
func GetMLEngineerSkill(workspace string) *MLEngineerSkill {
	if nativeSkillsRegistry.mlEng == nil {
		nativeSkillsRegistry.mlEng = NewMLEngineerSkill(workspace)
	}
	return nativeSkillsRegistry.mlEng
}

// GetOdooDeveloperSkill returns the singleton instance of OdooDeveloperSkill.
// Thread-safe lazy initialization.
func GetOdooDeveloperSkill(workspace string) *OdooDeveloperSkill {
	if nativeSkillsRegistry.odooDev == nil {
		nativeSkillsRegistry.odooDev = NewOdooDeveloperSkill(workspace)
	}
	return nativeSkillsRegistry.odooDev
}

// LoadNativeQueueBatchSkill returns the complete skill context from the native Go implementation.
// This replaces the file-based loading with compiled-in documentation.
func (sl *SkillsLoader) LoadNativeQueueBatchSkill() string {
	skill := GetQueueBatchSkill(sl.workspace)
	return skill.BuildSkillContext()
}

// BuildNativeQueueBatchSummary returns an XML summary from the native implementation.
func (sl *SkillsLoader) BuildNativeQueueBatchSummary() string {
	skill := GetQueueBatchSkill(sl.workspace)
	return skill.BuildSummary()
}

// LoadNativeBinanceMCPSkill returns the complete skill context from the native Go implementation.
func (sl *SkillsLoader) LoadNativeBinanceMCPSkill() string {
	skill := GetBinanceMCPSkill(sl.workspace)
	return skill.BuildSkillContext()
}

// BuildNativeBinanceMCPSummary returns an XML summary from the native implementation.
func (sl *SkillsLoader) BuildNativeBinanceMCPSummary() string {
	skill := GetBinanceMCPSkill(sl.workspace)
	return skill.BuildSummary()
}

// LoadNativeFullStackDeveloperSkill returns the complete skill context from the native Go implementation.
func (sl *SkillsLoader) LoadNativeFullStackDeveloperSkill() string {
	skill := GetFullStackDeveloperSkill(sl.workspace)
	return skill.BuildSkillContext()
}

// BuildNativeFullStackDeveloperSummary returns an XML summary from the native implementation.
func (sl *SkillsLoader) BuildNativeFullStackDeveloperSummary() string {
	skill := GetFullStackDeveloperSkill(sl.workspace)
	return skill.BuildSummary()
}

// LoadNativeN8NWorkflowSkill returns the complete skill context from the native Go implementation.
func (sl *SkillsLoader) LoadNativeN8NWorkflowSkill() string {
	skill := GetN8NWorkflowSkill(sl.workspace)
	return skill.BuildSkillContext()
}

// BuildNativeN8NWorkflowSummary returns an XML summary from the native implementation.
func (sl *SkillsLoader) BuildNativeN8NWorkflowSummary() string {
	skill := GetN8NWorkflowSkill(sl.workspace)
	return skill.BuildSummary()
}

// LoadNativeAgentTeamWorkflowSkill returns the complete skill context from the native Go implementation.
func (sl *SkillsLoader) LoadNativeAgentTeamWorkflowSkill() string {
	skill := GetAgentTeamWorkflowSkill(sl.workspace)
	return skill.BuildSkillContext()
}

// BuildNativeAgentTeamWorkflowSummary returns an XML summary from the native implementation.
func (sl *SkillsLoader) BuildNativeAgentTeamWorkflowSummary() string {
	skill := GetAgentTeamWorkflowSkill(sl.workspace)
	return skill.BuildSummary()
}

// LoadNativeResearcherSkill returns the complete skill context from the native Go implementation.
func (sl *SkillsLoader) LoadNativeResearcherSkill() string {
	skill := GetResearcherSkill(sl.workspace)
	return skill.BuildSkillContext()
}

// BuildNativeResearcherSummary returns an XML summary from the native implementation.
func (sl *SkillsLoader) BuildNativeResearcherSummary() string {
	skill := GetResearcherSkill(sl.workspace)
	return skill.BuildSummary()
}

// LoadNativeBackendDeveloperSkill returns the complete skill context from the native Go implementation.
func (sl *SkillsLoader) LoadNativeBackendDeveloperSkill() string {
	skill := GetBackendDeveloperSkill(sl.workspace)
	return skill.BuildSkillContext()
}

// BuildNativeBackendDeveloperSummary returns an XML summary from the native implementation.
func (sl *SkillsLoader) BuildNativeBackendDeveloperSummary() string {
	skill := GetBackendDeveloperSkill(sl.workspace)
	return skill.BuildSummary()
}

// LoadNativeFrontendDeveloperSkill returns the complete skill context from the native Go implementation.
func (sl *SkillsLoader) LoadNativeFrontendDeveloperSkill() string {
	skill := GetFrontendDeveloperSkill(sl.workspace)
	return skill.BuildSkillContext()
}

// BuildNativeFrontendDeveloperSummary returns an XML summary from the native implementation.
func (sl *SkillsLoader) BuildNativeFrontendDeveloperSummary() string {
	skill := GetFrontendDeveloperSkill(sl.workspace)
	return skill.BuildSummary()
}

// LoadNativeDevOpsEngineerSkill returns the complete skill context from the native Go implementation.
func (sl *SkillsLoader) LoadNativeDevOpsEngineerSkill() string {
	skill := GetDevOpsEngineerSkill(sl.workspace)
	return skill.BuildSkillContext()
}

// BuildNativeDevOpsEngineerSummary returns an XML summary from the native implementation.
func (sl *SkillsLoader) BuildNativeDevOpsEngineerSummary() string {
	skill := GetDevOpsEngineerSkill(sl.workspace)
	return skill.BuildSummary()
}

// LoadNativeSecurityEngineerSkill returns the complete skill context from the native Go implementation.
func (sl *SkillsLoader) LoadNativeSecurityEngineerSkill() string {
	skill := GetSecurityEngineerSkill(sl.workspace)
	return skill.BuildSkillContext()
}

// BuildNativeSecurityEngineerSummary returns an XML summary from the native implementation.
func (sl *SkillsLoader) BuildNativeSecurityEngineerSummary() string {
	skill := GetSecurityEngineerSkill(sl.workspace)
	return skill.BuildSummary()
}

// LoadNativeQAEngineerSkill returns the complete skill context from the native Go implementation.
func (sl *SkillsLoader) LoadNativeQAEngineerSkill() string {
	skill := GetQAEngineerSkill(sl.workspace)
	return skill.BuildSkillContext()
}

// BuildNativeQAEngineerSummary returns an XML summary from the native implementation.
func (sl *SkillsLoader) BuildNativeQAEngineerSummary() string {
	skill := GetQAEngineerSkill(sl.workspace)
	return skill.BuildSummary()
}

// LoadNativeDataEngineerSkill returns the complete skill context from the native Go implementation.
func (sl *SkillsLoader) LoadNativeDataEngineerSkill() string {
	skill := GetDataEngineerSkill(sl.workspace)
	return skill.BuildSkillContext()
}

// BuildNativeDataEngineerSummary returns an XML summary from the native implementation.
func (sl *SkillsLoader) BuildNativeDataEngineerSummary() string {
	skill := GetDataEngineerSkill(sl.workspace)
	return skill.BuildSummary()
}

// LoadNativeMLEngineerSkill returns the complete skill context from the native Go implementation.
func (sl *SkillsLoader) LoadNativeMLEngineerSkill() string {
	skill := GetMLEngineerSkill(sl.workspace)
	return skill.BuildSkillContext()
}

// BuildNativeMLEngineerSummary returns an XML summary from the native implementation.
func (sl *SkillsLoader) BuildNativeMLEngineerSummary() string {
	skill := GetMLEngineerSkill(sl.workspace)
	return skill.BuildSummary()
}

// LoadNativeOdooDeveloperSkill returns the complete skill context from the native Go implementation.
func (sl *SkillsLoader) LoadNativeOdooDeveloperSkill() string {
	skill := GetOdooDeveloperSkill(sl.workspace)
	return skill.BuildSkillContext()
}

// BuildNativeOdooDeveloperSummary returns an XML summary from the native implementation.
func (sl *SkillsLoader) BuildNativeOdooDeveloperSummary() string {
	skill := GetOdooDeveloperSkill(sl.workspace)
	return skill.BuildSummary()
}

// NativeSkillContent returns the full skill context for a compiled-in native skill by name.
// Returns the content string and true if found, or ("", false) if the name is not a registered native skill.
func (sl *SkillsLoader) NativeSkillContent(name string) (string, bool) {
	switch name {
	case "queue_batch":
		return sl.LoadNativeQueueBatchSkill(), true
	case "binance_mcp":
		return sl.LoadNativeBinanceMCPSkill(), true
	case "fullstack_developer":
		return sl.LoadNativeFullStackDeveloperSkill(), true
	case "n8n_workflow":
		return sl.LoadNativeN8NWorkflowSkill(), true
	case "agent_team_workflow":
		return sl.LoadNativeAgentTeamWorkflowSkill(), true
	case "researcher":
		return sl.LoadNativeResearcherSkill(), true
	case "backend_developer":
		return sl.LoadNativeBackendDeveloperSkill(), true
	case "frontend_developer":
		return sl.LoadNativeFrontendDeveloperSkill(), true
	case "devops_engineer":
		return sl.LoadNativeDevOpsEngineerSkill(), true
	case "security_engineer":
		return sl.LoadNativeSecurityEngineerSkill(), true
	case "qa_engineer":
		return sl.LoadNativeQAEngineerSkill(), true
	case "data_engineer":
		return sl.LoadNativeDataEngineerSkill(), true
	case "ml_engineer":
		return sl.LoadNativeMLEngineerSkill(), true
	case "odoo_developer":
		return sl.LoadNativeOdooDeveloperSkill(), true
	default:
		return "", false
	}
}

// listNativeSkills returns all native compiled-in skills.
func (sl *SkillsLoader) listNativeSkills() []SkillInfo {
	return []SkillInfo{
		{
			Name:        "queue_batch",
			Description: "Delegate heavy tasks to background queue using fire-and-forget pattern",
			Source:      "native",
			Path:        "builtin://queue_batch",
		},
		{
			Name:        "binance_mcp",
			Description: "Connect to Binance MCP server for trading and market data",
			Source:      "native",
			Path:        "builtin://binance_mcp",
		},
		{
			Name:        "fullstack_developer",
			Description: "Expert full-stack development assistant with patterns and best practices",
			Source:      "native",
			Path:        "builtin://fullstack_developer",
		},
		{
			Name:        "n8n_workflow",
			Description: "n8n Workflow Automation Expert - Create production-ready workflows with valid JSON",
			Source:      "native",
			Path:        "builtin://n8n_workflow",
		},
		{
			Name:        "agent_team_workflow",
			Description: "Multi-Agent Team Orchestrator - Organize optimal agent teams for any task",
			Source:      "native",
			Path:        "builtin://agent_team_workflow",
		},
		// Fase 0: researcher (existía sin registrar)
		{
			Name:        "researcher",
			Description: "Deep Research Agent — web search, source evaluation, information synthesis, and structured reporting",
			Source:      "native",
			Path:        "builtin://researcher",
		},
		// Fase 5: Engineering role skills
		{
			Name:        "backend_developer",
			Description: "Backend development expert: REST APIs, databases, microservices, performance, security",
			Source:      "native",
			Path:        "builtin://backend_developer",
		},
		{
			Name:        "frontend_developer",
			Description: "Frontend development expert: React, Vue, performance, accessibility, design systems",
			Source:      "native",
			Path:        "builtin://frontend_developer",
		},
		{
			Name:        "devops_engineer",
			Description: "DevOps expert: CI/CD pipelines, containers, infrastructure as code, monitoring, SRE",
			Source:      "native",
			Path:        "builtin://devops_engineer",
		},
		{
			Name:        "security_engineer",
			Description: "Security expert: OWASP, penetration testing, hardening, threat modeling, compliance",
			Source:      "native",
			Path:        "builtin://security_engineer",
		},
		{
			Name:        "qa_engineer",
			Description: "QA expert: testing strategies, test automation, coverage analysis, quality gates",
			Source:      "native",
			Path:        "builtin://qa_engineer",
		},
		{
			Name:        "data_engineer",
			Description: "Data engineering expert: ETL pipelines, data warehouses, streaming, data quality",
			Source:      "native",
			Path:        "builtin://data_engineer",
		},
		{
			Name:        "ml_engineer",
			Description: "ML/AI expert: model training, deployment, evaluation pipelines, MLOps, feature engineering",
			Source:      "native",
			Path:        "builtin://ml_engineer",
		},
		// Odoo Developer skill
		{
			Name:        "odoo_developer",
			Description: "Principal Odoo Architect & QA Engineer - Odoo ecosystems, Pine Script migration, L10n-Mexico, CFDI 4.0",
			Source:      "native",
			Path:        "builtin://odoo_developer",
		},
	}
}

func (sl *SkillsLoader) getSkillMetadata(skillPath string) *SkillMetadata {
	content, err := os.ReadFile(skillPath)
	if err != nil {
		logger.WarnCF("skills", "Failed to read skill metadata",
			map[string]any{
				"skill_path": skillPath,
				"error":      err.Error(),
			})
		return nil
	}

	frontmatter := sl.extractFrontmatter(string(content))
	if frontmatter == "" {
		return &SkillMetadata{
			Name: filepath.Base(filepath.Dir(skillPath)),
		}
	}

	// Try JSON first (for backward compatibility)
	var jsonMeta struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := json.Unmarshal([]byte(frontmatter), &jsonMeta); err == nil {
		return &SkillMetadata{
			Name:        jsonMeta.Name,
			Description: jsonMeta.Description,
		}
	}

	// Fall back to simple YAML parsing
	yamlMeta := sl.parseSimpleYAML(frontmatter)
	return &SkillMetadata{
		Name:        yamlMeta["name"],
		Description: yamlMeta["description"],
	}
}

// parseSimpleYAML parses simple key: value YAML format
// Example: name: github\n description: "..."
// Normalizes line endings to handle \n (Unix), \r\n (Windows), and \r (classic Mac)
func (sl *SkillsLoader) parseSimpleYAML(content string) map[string]string {
	result := make(map[string]string)

	// Normalize line endings: convert \r\n and \r to \n
	normalized := strings.ReplaceAll(content, "\r\n", "\n")
	normalized = strings.ReplaceAll(normalized, "\r", "\n")

	for _, line := range strings.Split(normalized, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			// Remove quotes if present
			value = strings.Trim(value, "\"'")
			result[key] = value
		}
	}

	return result
}

func (sl *SkillsLoader) extractFrontmatter(content string) string {
	// Support \n (Unix), \r\n (Windows), and \r (classic Mac) line endings for frontmatter blocks
	match := reFrontmatter.FindStringSubmatch(content)
	if len(match) > 1 {
		return match[1]
	}
	return ""
}

// parseSkillFrontmatter extracts and parses YAML frontmatter from skill content
func parseSkillFrontmatter(content string) map[string]string {
	// Find frontmatter between --- delimiters
	match := reFrontmatter.FindStringSubmatch(content)
	if len(match) < 2 {
		return nil
	}

	frontmatter := match[1]
	result := make(map[string]string)

	// Normalize line endings
	normalized := strings.ReplaceAll(frontmatter, "\r\n", "\n")
	normalized = strings.ReplaceAll(normalized, "\r", "\n")

	for _, line := range strings.Split(normalized, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			// Remove quotes if present
			value = strings.Trim(value, "\"'")
			result[key] = value
		}
	}

	return result
}

func (sl *SkillsLoader) stripFrontmatter(content string) string {
	return reStripFrontmatter.ReplaceAllString(content, "")
}

func escapeXML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	return s
}
