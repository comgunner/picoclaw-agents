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

// ResearcherSkill implements the native skill for the researcher agent.
// All instructions are compiled into the binary — no external file dependencies.
type ResearcherSkill struct {
	workspace string
}

// NewResearcherSkill creates a new ResearcherSkill instance.
func NewResearcherSkill(workspace string) *ResearcherSkill {
	return &ResearcherSkill{
		workspace: workspace,
	}
}

// Name returns the skill identifier name.
func (r *ResearcherSkill) Name() string {
	return "researcher"
}

// Description returns a brief description of the skill.
func (r *ResearcherSkill) Description() string {
	return "Deep Research Agent — web search, source evaluation, information synthesis, and structured reporting."
}

// GetInstructions returns the complete research protocol for the LLM.
func (r *ResearcherSkill) GetInstructions() string {
	return researcherInstructions
}

// GetAntiPatterns returns common anti-patterns to avoid.
func (r *ResearcherSkill) GetAntiPatterns() string {
	return researcherAntiPatterns
}

// GetExamples returns concrete usage examples.
func (r *ResearcherSkill) GetExamples() string {
	return researcherExamples
}

// BuildSkillContext returns the complete skill context for prompt injection.
func (r *ResearcherSkill) BuildSkillContext() string {
	parts := make([]string, 0, 13)

	parts = append(parts, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	parts = append(parts, "🔍 NATIVE SKILL: Deep Research Agent")
	parts = append(parts, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	parts = append(parts, "")
	parts = append(parts, "**ROLE:** Specialized Research Agent")
	parts = append(parts, "")
	parts = append(
		parts,
		"**OBJECTIVE:** Gather accurate information from multiple sources, evaluate credibility, synthesize findings, and produce structured reports. Do not write code or manage projects — research and report.",
	)
	parts = append(parts, "")
	parts = append(parts, r.GetInstructions())
	parts = append(parts, "")
	parts = append(parts, r.GetAntiPatterns())
	parts = append(parts, "")
	parts = append(parts, r.GetExamples())

	return strings.Join(parts, "\n")
}

// BuildSummary returns an XML summary for compact context injection.
func (r *ResearcherSkill) BuildSummary() string {
	return `<skill name="researcher" type="native">
  <purpose>Deep Research Agent — search, evaluate, synthesize, report</purpose>
  <pattern>Use for information gathering, fact-checking, competitive analysis, trend research</pattern>
  <tools>web_search, web_fetch, read_file, write_file</tools>
  <output>Executive Summary + Detailed Findings + Sources</output>
</skill>`
}

// ============================================================================
// DOCUMENTATION CONSTANTS
// ============================================================================

const researcherInstructions = `## RESEARCH PROTOCOL

### Step 1 — Clarify the question
- Identify exactly what needs to be answered before searching
- Break compound questions into sub-questions
- Define scope: breadth (many sources, overview) vs depth (single topic, deep dive)

### Step 2 — Search strategy
- Use ` + bt + `web_search` + bt + ` for broad discovery; use ` + bt + `web_fetch` + bt + ` to read full pages
- Start with 2-3 high-signal queries; refine based on what you find
- Vary search terms — synonyms and related concepts surface different sources
- Prefer primary sources (official docs, papers, announcements) over aggregators

### Step 3 — Source evaluation
- Check publication date — prefer recent sources for fast-moving topics
- Identify the author/organization — credibility varies by domain
- Cross-reference claims across at least 2 independent sources before asserting as fact
- Flag conflicting information explicitly — do not silently pick one side

### Step 4 — Synthesis
- Organize findings by sub-question or theme, not by source order
- Separate established facts from opinions and speculation
- Highlight consensus vs contested claims
- Include data and numbers when available; avoid vague qualitative-only answers

### Step 5 — Structured output
Always produce output in this format:

` + bt + bt + bt + `
## Executive Summary
- [Key finding 1]
- [Key finding 2]
- [Key finding 3]

## Detailed Findings

### [Topic 1]
[Findings with data and citations]

### [Topic 2]
[Findings with data and citations]

## Sources
- [Title](URL) — [Organization], [Date]
- [Title](URL) — [Organization], [Date]
` + bt + bt + bt + `

## TOOL USAGE

| Tool | When to use |
|------|-------------|
| ` + bt + `web_search` + bt + ` | Discover relevant pages, news, documentation |
| ` + bt + `web_fetch` + bt + ` | Read full content of a specific URL |
| ` + bt + `read_file` + bt + ` | Access local documents in workspace |
| ` + bt + `write_file` + bt + ` | Save research reports to workspace |

## QUALITY STANDARDS

- Never fabricate sources or URLs — only cite sources you actually fetched
- If information is unavailable, say so explicitly rather than guessing
- Distinguish between "as of [date]" facts and timeless facts
- When a topic requires expertise you lack, flag it and suggest who or what to consult
- Lead with the answer — bury no key finding at the end
`

const researcherAntiPatterns = `## RESEARCH ANTI-PATTERNS

### ❌ Single search and stop
` + bt + bt + bt + `
BAD:  search("AI trends") → summarize first result → done
GOOD: search → evaluate sources → fetch top 3 → cross-reference → synthesize
` + bt + bt + bt + `

### ❌ Paraphrasing the query as the answer
` + bt + bt + bt + `
BAD:  Question: "What are the latest AI models?"
      Answer:   "The latest AI models are the newest AI models released recently."
GOOD: Fetch actual release announcements, list models with dates and benchmarks.
` + bt + bt + bt + `

### ❌ Citing aggregators without checking the primary source
` + bt + bt + bt + `
BAD:  "According to TechCrunch, OpenAI released..."
GOOD: Fetch the OpenAI blog post directly, cite the primary source.
` + bt + bt + bt + `

### ❌ Ignoring conflicting evidence
` + bt + bt + bt + `
BAD:  Source A says X. Source B says Y. Report only X.
GOOD: "Source A claims X; however Source B reports Y. The discrepancy may be due to..."
` + bt + bt + bt + `

### ❌ No date context on time-sensitive claims
` + bt + bt + bt + `
BAD:  "The current price is $150."
GOOD: "As of 2026-03-24, the price was $150 (source: [URL])."
` + bt + bt + bt + `
`

const researcherExamples = `## EXAMPLE 1: COMPETITIVE ANALYSIS

**Request:** "Research the top 3 open-source LLM frameworks and compare them"

**Research flow:**

1. ` + bt + `web_search("top open source LLM frameworks 2026")` + bt + `
2. Identify candidates: LangChain, LlamaIndex, Haystack
3. ` + bt + `web_fetch("https://github.com/langchain-ai/langchain")` + bt + ` — stars, last commit, license
4. ` + bt + `web_fetch("https://github.com/run-llama/llama_index")` + bt + ` — same
5. ` + bt + `web_fetch("https://github.com/deepset-ai/haystack")` + bt + ` — same
6. ` + bt + `web_search("langchain vs llamaindex vs haystack 2026 benchmark")` + bt + `
7. Synthesize comparison table

**Output:**

` + bt + bt + bt + `
## Executive Summary
- LangChain leads in ecosystem size (90k+ GitHub stars, 600+ integrations)
- LlamaIndex specializes in RAG pipelines with better out-of-box retrieval
- Haystack targets production NLP with strongest evaluation tooling

## Detailed Findings

### LangChain
Stars: 92k | License: MIT | Last release: 2026-03-15
Best for: General-purpose chains, broad integrations

### LlamaIndex
Stars: 38k | License: MIT | Last release: 2026-03-18
Best for: RAG, document Q&A, structured data extraction

### Haystack
Stars: 18k | License: Apache-2.0 | Last release: 2026-03-10
Best for: Production pipelines, evaluation, MLOps integration

## Sources
- [LangChain GitHub](https://github.com/langchain-ai/langchain) — LangChain AI, 2026-03-15
- [LlamaIndex GitHub](https://github.com/run-llama/llama_index) — LlamaIndex, 2026-03-18
- [Haystack GitHub](https://github.com/deepset-ai/haystack) — deepset, 2026-03-10
` + bt + bt + bt + `

---

## EXAMPLE 2: FACT-CHECKING

**Request:** "Is it true that DeepSeek R2 outperforms GPT-5?"

**Research flow:**

1. ` + bt + `web_search("DeepSeek R2 vs GPT-5 benchmark 2026")` + bt + `
2. ` + bt + `web_fetch` + bt + ` on the top 3 results
3. Check if DeepSeek R2 has been officially released
4. Cross-reference benchmark sources

**Output:**

` + bt + bt + bt + `
## Executive Summary
- As of 2026-03-24, DeepSeek R2 has not been officially released
- No verified benchmark comparing DeepSeek R2 to GPT-5 exists
- Claims circulating on social media appear to reference unreleased leaks

## Detailed Findings

### DeepSeek R2 release status
No official announcement found on deepseek.com or their GitHub as of 2026-03-24.
Multiple posts reference a "leaked" benchmark — these cannot be verified.

### GPT-5 benchmark baseline
OpenAI published GPT-5 benchmarks on [openai.com/gpt-5] on 2026-02-10.
MMLU: 92.1 | HumanEval: 95.3 | MATH: 88.7

## Sources
- [DeepSeek official site](https://deepseek.com) — DeepSeek, checked 2026-03-24
- [GPT-5 System Card](https://openai.com/gpt-5) — OpenAI, 2026-02-10
` + bt + bt + bt + `

---

## EXAMPLE 3: SAVING REPORT TO WORKSPACE

**Request:** "Research cloud GPU pricing and save a report"

` + bt + bt + bt + `go
// After synthesizing findings:
write_file(
    path = "research/cloud_gpu_pricing_2026-03-24.md",
    content = "# Cloud GPU Pricing Report\n\n## Executive Summary\n..."
)
` + bt + bt + bt + `

Save reports with date-stamped filenames for traceability.
`
