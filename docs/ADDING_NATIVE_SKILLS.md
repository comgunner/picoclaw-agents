# Adding Native Skills - Developer Guide

**Version:** 3.6.0  
**Last Updated:** March 26, 2026

## Overview

This guide walks you through creating a new **Native Skill** for PicoClaw. Native Skills are role definitions compiled directly into the binary, providing specialized expertise without external file dependencies.

**Time to complete:** 30-60 minutes per skill  
**Difficulty:** Intermediate (requires Go knowledge)

---

## When to Create a Native Skill

Create a Native Skill when:
- ✅ Defining a **specialized role** (e.g., "backend developer", "security engineer")
- ✅ Providing **domain expertise** (e.g., "medical advisor", "legal consultant")
- ✅ Establishing **behavioral guidelines** (e.g., "code reviewer", "mentor")
- ✅ Creating **workflow patterns** (e.g., "agile coach", "product manager")

**Do NOT create a Native Skill for:**
- ❌ Actions that modify state (use **Tools** instead)
- ❌ External API integrations (use **Tools** or **Providers**)
- ❌ Simple prompt templates (use **workspace files** instead)

---

## Architecture

### Skill vs Tool

```
Native Skill                          Tool
┌─────────────────────┐              ┌─────────────────────┐
│ Role Definition     │              │ Action Definition   │
│ "You ARE a backend  │              │ "You CAN read files"│
│  developer"         │              │                     │
│                     │              │                     │
│ Injected into       │              │ Called via          │
│ system prompt       │              │ tool_use            │
└─────────────────────┘              └─────────────────────┘
        │                                      │
        ▼                                      ▼
  LLM adopts personality              LLM executes action
```

### Skill Lifecycle

```
1. Developer creates skill file (backend_developer.go)
         ↓
2. Skill implements Skill interface
         ↓
3. Register in loader.go nativeSkillsRegistry
         ↓
4. Add to listNativeSkills()
         ↓
5. User configures in config.json
         ↓
6. SkillsLoader.LoadNativeXxxSkill() injects into prompt
```

---

## Step-by-Step Guide

### Step 1: Create Skill File

Create a new file in `pkg/skills/`:

```bash
touch pkg/skills/your_skill.go
```

**Naming convention:**
- File: `your_skill.go` (snake_case)
- Type: `YourSkill` (PascalCase)
- ID: `your_skill` (snake_case, matches config.json)

### Step 2: Implement Skill Structure

Use this template:

```go
// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
// Copyright (c) 2026 PicoClaw contributors

package skills

import "strings"

// YourSkill implements the native skill for [role description].
type YourSkill struct {
	workspace string
}

// NewYourSkill creates a new YourSkill instance.
func NewYourSkill(workspace string) *YourSkill {
	return &YourSkill{
		workspace: workspace,
	}
}

// Name returns the skill identifier name.
func (y *YourSkill) Name() string {
	return "your_skill"
}

// Description returns a brief description of the skill.
func (y *YourSkill) Description() string {
	return "Your skill description — key capabilities in 10-15 words"
}

// GetInstructions returns the complete role protocol for the LLM.
func (y *YourSkill) GetInstructions() string {
	return yourSkillInstructions
}

// GetAntiPatterns returns common anti-patterns to avoid.
func (y *YourSkill) GetAntiPatterns() string {
	return yourSkillAntiPatterns
}

// GetExamples returns concrete usage examples.
func (y *YourSkill) GetExamples() string {
	return yourSkillExamples
}

// BuildSkillContext returns the complete skill context for prompt injection.
func (y *YourSkill) BuildSkillContext() string {
	var parts []string

	parts = append(parts, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	parts = append(parts, "🎯 NATIVE SKILL: Your Skill Name")
	parts = append(parts, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	parts = append(parts, "")
	parts = append(parts, "**ROLE:** Expert [role] specializing in [specialization].")
	parts = append(parts, "")
	parts = append(parts, y.GetInstructions())
	parts = append(parts, "")
	parts = append(parts, y.GetAntiPatterns())
	parts = append(parts, "")
	parts = append(parts, y.GetExamples())

	return strings.Join(parts, "\n")
}

// BuildSummary returns an XML summary for compact context injection.
func (y *YourSkill) BuildSummary() string {
	return `<skill name="your_skill" type="native">
  <purpose>Your skill purpose — 10 word summary</purpose>
  <pattern>Use for [use case 1], [use case 2], [use case 3]</pattern>
  <stacks>Technology1, Technology2, Technology3</stacks>
  <practices>Practice1, Practice2, Practice3</practices>
</skill>`
}

// ============================================================================
// DOCUMENTATION CONSTANTS
// ============================================================================

const yourSkillInstructions = `## CORE RESPONSIBILITIES

### 1. Responsibility Area 1
- Detail what the skill does
- Best practices to follow
- Quality standards

### 2. Responsibility Area 2
- More details
- Examples of tasks
- Tools used

## TECHNOLOGY STACK

### Category
- **Primary**: Tool1, Tool2, Tool3
- **Secondary**: Tool4, Tool5

## BEST PRACTICES

### Practice Name
` + bt + bt + bt + `
code example here
` + bt + bt + bt + `

## QUALITY CHECKLIST

Before considering a task complete:

- [ ] Checklist item 1
- [ ] Checklist item 2
- [ ] Checklist item 3
`

const yourSkillAntiPatterns = `## ANTI-PATTERNS

### ❌ Bad Pattern Name
` + bt + bt + bt + `
BAD:  code showing what NOT to do

GOOD: code showing correct approach
` + bt + bt + bt + `

### ❌ Another Anti-Pattern
More explanation of what to avoid
`

const yourSkillExamples = `## EXAMPLE 1: Common Task

**Request:** "Example request description"

**Expert Response:**

` + bt + bt + bt + `
code example showing expert solution
with multiple lines
` + bt + bt + bt + `

## EXAMPLE 2: Advanced Task

**Request:** "More complex request"

**Expert Response:**

` + bt + bt + bt + `
another code example
showing advanced usage
` + bt + bt + bt + `
`
```

### Step 3: Replace Template Placeholders

Replace all occurrences:
- `YourSkill` → Your actual skill name (PascalCase)
- `your_skill` → Your skill ID (snake_case)
- `yourSkillInstructions` → yourSkillNameInstructions
- Role descriptions → Your actual role content

### Step 4: Write Skill Content

#### Instructions Section (2,000-4,000 tokens)

Include:
1. **Role Definition** (1-2 sentences)
2. **Core Responsibilities** (5-7 areas)
3. **Technology Stack** (relevant tools/frameworks)
4. **Best Practices** (with code examples)
5. **Quality Checklist** (actionable items)

**Example:**
```markdown
## CORE RESPONSIBILITIES

### 1. API Design & Development
- Design RESTful APIs with clear contracts
- Implement GraphQL schemas when appropriate
- Ensure proper HTTP status codes and error handling
```

#### Anti-Patterns Section (500-1,000 tokens)

Include:
- 5-10 common mistakes
- BAD/GOOD code comparisons
- Clear explanations of why it's wrong

**Example:**
```markdown
### ❌ Hardcoded Credentials
` + bt + bt + bt + `javascript
// BAD
const API_KEY = "sk-1234567890"

// GOOD
const API_KEY = process.env.API_KEY
` + bt + bt + bt + `
```

#### Examples Section (1,000-2,000 tokens)

Include:
- 3-5 realistic scenarios
- Request/Response format
- Complete, working code

**Example:**
```markdown
## EXAMPLE 1: CREATE REST API ENDPOINT

**Request:** "Create Express endpoint to create users"

**Expert Response:**

` + bt + bt + bt + `javascript
// Complete working code
router.post('/', async (req, res) => {
  // Implementation
})
` + bt + bt + bt + `
```

### Step 5: Register in loader.go

Open `pkg/skills/loader.go` and make these changes:

#### 5.1: Add to nativeSkillsRegistry struct

```go
var nativeSkillsRegistry = struct {
	// ... existing fields ...
	yourSkill *YourSkill  // Add this line
}{
	// ... existing initializations ...
	yourSkill: nil,  // Add this line
}
```

#### 5.2: Add Getter Method

After existing getters, add:

```go
// GetYourSkill returns the singleton instance of YourSkill.
// Thread-safe lazy initialization.
func GetYourSkill(workspace string) *YourSkill {
	if nativeSkillsRegistry.yourSkill == nil {
		nativeSkillsRegistry.yourSkill = NewYourSkill(workspace)
	}
	return nativeSkillsRegistry.yourSkill
}
```

#### 5.3: Add Load and Build Methods

After existing load methods, add:

```go
// LoadNativeYourSkill returns the complete skill context from the native Go implementation.
func (sl *SkillsLoader) LoadNativeYourSkill() string {
	skill := GetYourSkill(sl.workspace)
	return skill.BuildSkillContext()
}

// BuildNativeYourSkillSummary returns an XML summary from the native implementation.
func (sl *SkillsLoader) BuildNativeYourSkillSummary() string {
	skill := GetYourSkill(sl.workspace)
	return skill.BuildSummary()
}
```

#### 5.4: Add to listNativeSkills()

In the `listNativeSkills()` function, add:

```go
{
	Name:        "your_skill",
	Description: "Your skill description — 10-15 words",
	Source:      "native",
	Path:        "builtin://your_skill",
},
```

### Step 6: Write Tests

Create `pkg/skills/your_skill_test.go`:

```go
package skills_test

import (
	"testing"

	"github.com/comgunner/picoclaw/pkg/skills"
	"github.com/stretchr/testify/assert"
)

func TestYourSkillName(t *testing.T) {
	skill := skills.NewYourSkill("/tmp/workspace")
	assert.Equal(t, "your_skill", skill.Name())
}

func TestYourSkillDescription(t *testing.T) {
	skill := skills.NewYourSkill("/tmp/workspace")
	desc := skill.Description()
	assert.NotEmpty(t, desc)
	assert.Contains(t, desc, "YourKeyword")
}

func TestYourSkillGetInstructions(t *testing.T) {
	skill := skills.NewYourSkill("/tmp/workspace")
	instructions := skill.GetInstructions()
	assert.NotEmpty(t, instructions)
	assert.Contains(t, instructions, "Responsibility")
}

func TestYourSkillBuildSkillContext(t *testing.T) {
	skill := skills.NewYourSkill("/tmp/workspace")
	context := skill.BuildSkillContext()
	assert.NotEmpty(t, context)
	assert.Contains(t, context, "ROLE:")
}

func TestYourSkillBuildSummary(t *testing.T) {
	skill := skills.NewYourSkill("/tmp/workspace")
	summary := skill.BuildSummary()
	assert.NotEmpty(t, summary)
	assert.Contains(t, summary, `<skill`)
	assert.Contains(t, summary, "your_skill")
}
```

### Step 7: Run Tests

```bash
# Test your skill
go test ./pkg/skills/your_skill_test.go ./pkg/skills/*.go -v

# Test all skills
go test ./pkg/skills/... -v

# Check for regressions
make test
```

### Step 8: Build and Verify

```bash
# Build binary
make build

# Verify skill is registered
./build/picoclaw-agents skills list | grep your_skill

# Test with query
./build/picoclaw-agents agent -m "As a [your_role], help me with [task]"
```

---

## Content Guidelines

### Writing Effective Instructions

✅ **DO:**
- Use action verbs ("Design", "Implement", "Ensure")
- Provide specific, actionable guidance
- Include technology recommendations
- Add quality checklists

❌ **DON'T:**
- Be vague ("Do good work")
- List every possible technology
- Write walls of text without structure
- Include outdated practices

### Writing Effective Anti-Patterns

✅ **DO:**
- Show concrete BAD/GOOD code comparisons
- Explain WHY it's an anti-pattern
- Cover common mistakes in the domain
- Use recognizable scenarios

❌ **DON'T:**
- List trivial issues
- Show unsafe code without warnings
- Assume advanced knowledge
- Use hypothetical examples

### Writing Effective Examples

✅ **DO:**
- Use realistic scenarios
- Provide complete, working code
- Show best practices in action
- Include comments explaining key decisions

❌ **DON'T:**
- Show incomplete snippets
- Use fake data that breaks logic
- Skip error handling
- Omit imports/dependencies

---

## Testing Checklist

Before submitting your skill:

- [ ] Skill compiles without errors
- [ ] All tests pass (`go test ./pkg/skills/...`)
- [ ] Skill appears in `skills list` output
- [ ] `BuildSkillContext()` returns substantial content (>1000 chars)
- [ ] `BuildSummary()` returns valid XML
- [ ] Anti-patterns contain ❌ markers
- [ ] Examples contain **Request:** format
- [ ] No hardcoded paths or credentials
- [ ] Follows Go formatting (`gofmt`)
- [ ] Passes linter (`golangci-lint run`)

---

## Example: Creating a "Data Analyst" Skill

### 1. Create File

```bash
touch pkg/skills/data_analyst.go
```

### 2. Implement

```go
package skills

import "strings"

type DataAnalystSkill struct {
	workspace string
}

func NewDataAnalystSkill(workspace string) *DataAnalystSkill {
	return &DataAnalystSkill{workspace: workspace}
}

func (d *DataAnalystSkill) Name() string { return "data_analyst" }

func (d *DataAnalystSkill) Description() string {
	return "Data analysis expert: SQL, Python, statistics, visualization, insights"
}

// ... implement other methods following template ...
```

### 3. Register in loader.go

```go
// In nativeSkillsRegistry struct
dataAnalyst *DataAnalystSkill

// In initialization
dataAnalyst: nil,

// Add getter
func GetDataAnalystSkill(workspace string) *DataAnalystSkill { ... }

// Add loader methods
func (sl *SkillsLoader) LoadNativeDataAnalystSkill() string { ... }
func (sl *SkillsLoader) BuildNativeDataAnalystSummary() string { ... }

// Add to listNativeSkills()
{
	Name:        "data_analyst",
	Description: "Data analysis expert: SQL, Python, statistics, visualization",
	Source:      "native",
	Path:        "builtin://data_analyst",
},
```

### 4. Test

```bash
go test ./pkg/skills/data_analyst_test.go ./pkg/skills/*.go -v
```

---

## Troubleshooting

### Skill Not Appearing in List

**Problem:** `skills list` doesn't show your skill

**Solution:**
1. Verify skill is in `listNativeSkills()` return slice
2. Rebuild binary: `make build`
3. Check for compilation errors

### Tests Failing

**Problem:** Tests fail with "nil pointer" or "method not found"

**Solution:**
1. Ensure all interface methods are implemented
2. Check method signatures match interface
3. Verify import paths are correct

### Content Not Showing

**Problem:** Skill context is empty or truncated

**Solution:**
1. Check string constants are properly concatenated
2. Verify backtick strings are closed
3. Look for unclosed quotes in content

---

## Best Practices

### Code Organization

1. **One skill per file** — keeps codebase maintainable
2. **Consistent naming** — follow existing patterns
3. **Group constants** — all content constants at bottom
4. **Use bt shorthand** — `bt = "` + "`" + `"` for triple backticks

### Content Quality

1. **Review for accuracy** — domain experts should validate
2. **Update regularly** — keep practices current
3. **Include diverse examples** — cover common scenarios
4. **Avoid bias** — multiple valid approaches exist

### Performance

1. **Keep it focused** — 2,000-4,000 tokens ideal
2. **Use markdown efficiently** — headers, lists, code blocks
3. **Avoid repetition** — each section should add value
4. **Consider token budget** — remember context window limits

---

## See Also

- **[SKILLS.md](SKILLS.md)**: User guide for using native skills
- **[fullstack_developer.go](pkg/skills/fullstack_developer.go)**: Reference implementation
- **[loader.go](pkg/skills/loader.go)**: Registration patterns
- **[CHANGELOG.md](CHANGELOG.md)**: v3.6.0 release notes

---

## Contributing

To contribute your skill to the main repository:

1. **Fork** the repository
2. **Create branch**: `feat/skill-your_skill_name`
3. **Implement skill** following this guide
4. **Write tests** with >90% coverage
5. **Run linter**: `golangci-lint run`
6. **Update docs**: Add to SKILLS.md table
7. **Submit PR** with description of skill capabilities

**Review criteria:**
- ✅ Follows template structure
- ✅ All tests pass
- ✅ Content is accurate and current
- ✅ Examples are complete and working
- ✅ No security issues in examples
- ✅ Documentation is comprehensive

---

*Happy skill building! 🚀*
