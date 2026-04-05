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

// QueueBatchSkill implements native skill for queue/batch delegation.
// It encapsulates all documentation and instructions as Go strings,
// compiled directly into the binary (no external file dependencies).
type QueueBatchSkill struct {
	workspace string
}

// NewQueueBatchSkill creates a new QueueBatchSkill instance.
func NewQueueBatchSkill(workspace string) *QueueBatchSkill {
	return &QueueBatchSkill{
		workspace: workspace,
	}
}

// Name returns the skill identifier name.
func (q *QueueBatchSkill) Name() string {
	return "queue_batch"
}

// Description returns a brief description of the skill.
func (q *QueueBatchSkill) Description() string {
	return "Delegate heavy tasks (Python, Batch, FFmpeg, CUDA) to background queue. Use 'fire and forget' pattern to save 95% tokens on long operations."
}

// GetInstructions returns the complete usage instructions for the LLM.
// This replaces the external SKILL.md file with native Go strings.
func (q *QueueBatchSkill) GetInstructions() string {
	return queueBatchInstructions
}

// GetAntiPatterns returns common anti-patterns to avoid.
func (q *QueueBatchSkill) GetAntiPatterns() string {
	return queueBatchAntiPatterns
}

// GetExamples returns concrete usage examples.
func (q *QueueBatchSkill) GetExamples() string {
	return queueBatchExamples
}

// BuildSkillContext returns the complete skill context for prompt injection.
// This is the method used by loader.go to inject the skill into the system prompt.
func (q *QueueBatchSkill) BuildSkillContext() string {
	var parts []string

	parts = append(parts, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	parts = append(parts, "🚀 NATIVE SKILL: Queue/Batch Delegation")
	parts = append(parts, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	parts = append(parts, "")
	parts = append(parts, "**PURPOSE:** Delegate heavy tasks to background using \"fire and forget\" pattern.")
	parts = append(parts, "")
	parts = append(parts, q.GetInstructions())
	parts = append(parts, "")
	parts = append(parts, q.GetAntiPatterns())
	parts = append(parts, "")
	parts = append(parts, q.GetExamples())

	return strings.Join(parts, "\n")
}

// BuildSummary returns an XML summary for compact context injection.
func (q *QueueBatchSkill) BuildSummary() string {
	return `<skill name="queue_batch" type="native">
  <purpose>Delegate heavy tasks to background queue</purpose>
  <pattern>Fire and forget for operations &gt; 30s</pattern>
  <tools>batch_id(), queue()</tools>
  <savings>~95% tokens on batch operations</savings>
</skill>`
}

// ============================================================================
// DOCUMENTATION CONSTANTS - Encapsulated documentation strings
// ============================================================================

const queueBatchInstructions = `## WHEN TO USE (CRITICAL)

Use this skill **AUTOMATICALLY** when you detect:

### ✅ Signals to Use Queue

1. **Duration > 30 seconds**
   - "generate 10 images" → Each image takes 2 min → Total 20 min
   - "train a model" → Can take hours
   - "render a video" → FFmpeg is heavy

2. **External Scripts**
   - User mentions Python, Bash, Batch, FFmpeg, CUDA
   - You need to execute ` + bt + `.py` + bt + `, ` + bt + `.sh` + bt + `, ` + bt + `.bat` + bt + `
   - Files exist in ` + bt + `workspace/scripts/` + bt + `

3. **Batch Operations**
   - "upload these 20 photos to Facebook"
   - "process all PDFs in this folder"
   - "generate variants of this text"

4. **Long Waits**
   - APIs with rate limiting
   - Large downloads
   - Compilations/builds

### ❌ WHEN NOT TO USE

- Simple queries (< 10s)
- Operations requiring intermediate approval
- Tasks needing LLM context at each step

## USAGE PATTERN (Step by Step)

### Step 1: Identify the Task

` + bt + bt + bt + `
User: "Generate 5 images with CUDA"
You: [Detect it's heavy → Use queue]
` + bt + bt + bt + `

### Step 2: Generate Unique ID

` + bt + bt + bt + `
batch_id(prefix="IMG_GEN")
// Result: #IMG_GEN_02_03_26_1530
` + bt + bt + bt + `

### Step 3: Prepare Payload

` + bt + bt + bt + `python
# LLM must generate this as part of the plan
payload = {
    "script": "scripts/batch_cuda_gen.py",
    "args": ["--count", "5", "--model", "sdxl"],
    "batch_id": "#IMG_GEN_02_03_26_1530"
}
` + bt + bt + bt + `

### Step 4: Queue and Release

` + bt + bt + bt + `
1. queue(action="enqueue", task_type="CUDA_GEN", payload=payload)
2. Message to user: "🔥 Started #IMG_GEN_02_03_26_1530. Estimated duration: 10 min."
3. RELEASE IMMEDIATELY - Don't wait
4. Script will report to QueueManager via /tmp/picoclaw_queue_{BATCH_ID}.json
` + bt + bt + bt + `

### Step 5: Follow-up (On-Demand)

` + bt + bt + bt + `
# Only if user asks
queue(task_id="#IMG_GEN_02_03_26_1530")
// Result: "⚙️ Processing (45% complete)"
` + bt + bt + bt + `

## INTEGRATION WITH PYTHON/BATCH SCRIPTS

### Communication Contract

External scripts MUST follow this contract:

` + bt + bt + bt + `python
#!/usr/bin/env python3
# scripts/batch_cuda_gen.py

import sys
import json
from pathlib import Path

def main():
    # 1. Receive BATCH_ID as first argument
    batch_id = sys.argv[1]  # e.g., "IMG_GEN_02_03_26_1530"

    # 2. State file (contract with QueueManager)
    state_file = Path(f"/tmp/picoclaw_queue_{batch_id}.json")

    # 3. Report start
    state_file.write_text(json.dumps({
        "status": "processing",
        "progress": 0,
        "message": "Starting generation..."
    }))

    # 4. Execute heavy work
    for i in range(5):
        generate_image(i)

        # 5. Update progress
        state_file.write_text(json.dumps({
            "status": "processing",
            "progress": (i + 1) * 20,
            "message": f"Image {i+1}/5 completed"
        }))

    # 6. Report completion
    state_file.write_text(json.dumps({
        "status": "completed",
        "progress": 100,
        "result": {
            "images": ["img1.png", "img2.png", ...]
        }
    }))

if __name__ == "__main__":
    main()
` + bt + bt + bt + `

### Monitoring from Go

The Go QueueManager monitors ` + bt + `/tmp/picoclaw_queue_*.json` + bt + `:
- Reads state without LLM intervention
- Notifies via MessageBus on completion
- LLM only gets involved if user asks
`

const queueBatchAntiPatterns = `## ANTI-PATTERNS TO AVOID

### ❌ Anti-Pattern 1: Unnecessary Polling

` + bt + bt + bt + `
# BAD - LLM asks every 30s
queue(task_id="#IMG_...")  # Done?
queue(task_id="#IMG_...")  # Now?
queue(task_id="#IMG_...")  # Now?

# GOOD - QueueManager notifies automatically
# LLM only checks if user asks
` + bt + bt + bt + `

### ❌ Anti-Pattern 2: Not Releasing

` + bt + bt + bt + `
# BAD - LLM waits blocked
queue(action="enqueue", ...)
[Waiting 10 min doing nothing]

# GOOD - LLM releases
queue(action="enqueue", ...)
"Attend to other users in the meantime"
` + bt + bt + bt + `

### ❌ Anti-Pattern 3: Manual IDs

` + bt + bt + bt + `
# BAD - LLM invents IDs
"Your task is IMG_123"  # Where was it registered?

# GOOD - Use batch_id()
batch_id(prefix="IMG") → #IMG_02_03_26_1630
` + bt + bt + bt + `
`

const queueBatchExamples = `## CONCRETE EXAMPLES

### Example 1: Batch Image Generation

**User:** "Generate 10 landscape images with SDXL"

**LLM (WITH this skill):**
` + bt + bt + bt + `
1. batch_id(prefix="IMG_GEN")
   → #IMG_GEN_02_03_26_1545

2. queue(action="enqueue",
         task_type="IMAGE_GEN",
         payload={
           "script": "scripts/batch_sdxl.py",
           "args": ["--prompt", "landscapes", "--count", "10"],
           "batch_id": "#IMG_GEN_02_03_26_1545"
         })

3. Message: "🔥 Started generating 10 images. ID: #IMG_GEN_02_03_26_1545"
   "Estimated duration: 20 minutes. I'll notify you when done."

4. [LLM releases immediately - ZERO waiting]

5. [20 min later] MessageBus notifies completion
6. LLM: "✅ Complete. All 10 images are in workspace/images/"
` + bt + bt + bt + `

**Tokens saved:** ~1500 (vs. checking every 2 min)

### Example 2: Bulk Social Media Upload

**User:** "Upload these 20 photos to Facebook and Instagram"

**LLM (WITH this skill):**
` + bt + bt + bt + `
1. batch_id(prefix="SOCIAL")
   → #SOCIAL_02_03_26_1600

2. queue(action="enqueue",
         task_type="POST_BUNDLE",
         payload={
           "script": "scripts/batch_post.py",
           "args": ["--platforms", "facebook,instagram", "--count", "20"],
           "batch_id": "#SOCIAL_02_03_26_1600"
         })

3. Message: "📱 Started uploading 20 photos. ID: #SOCIAL_02_03_26_1600"
   "I'll let you know when it's done."

4. [LLM releases]
` + bt + bt + bt + `

### Example 3: Model Training

**User:** "Train a model with this data"

**LLM (WITH this skill):**
` + bt + bt + bt + `
1. batch_id(prefix="TRAIN")
   → #TRAIN_02_03_26_1615

2. queue(action="enqueue",
         task_type="MODEL_TRAIN",
         payload={
           "script": "scripts/train_model.py",
           "args": ["--data", "workspace/data/", "--epochs", "100"],
           "batch_id": "#TRAIN_02_03_26_1615"
         })

3. Message: "🧠 Started training. ID: #TRAIN_02_03_26_1615"
   "Estimated duration: 2 hours. I'll notify you when complete."

4. [LLM releases for 2 hours - ZERO intervention]
` + bt + bt + bt + `

## QUICK SUMMARY

| Scenario | Use Queue | Token Savings |
|----------|-----------|---------------|
| Generate 1 image | ❌ No | - |
| Generate 5+ images | ✅ Yes | ~90% |
| Upload 1 photo | ❌ No | - |
| Upload 20 photos | ✅ Yes | ~95% |
| Train model | ✅ Yes | ~98% |
| Render video | ✅ Yes | ~95% |

**Mnemonic rule:** "If it takes more than 30 seconds or it's a batch, use the queue"
`
