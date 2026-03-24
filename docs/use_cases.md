# Dynamic Use Case: From Telegram Message to Local Fix

This document describes the **dynamic, real-time flow** of how PicoClaw processes an instruction received via Telegram to execute a full repair in the local repository using the `config_dev_multiple_models.example.json` configuration.

> **PicoClaw v3.4.1**: Now includes **Fast-path Slash Commands** for instant operations and **Global Tracker** for perfect multi-agent consistency.

## Dynamic Interaction Flow

### 1. Input: The Trigger (Telegram)
The user, from their mobile device, notices an issue or decides on an improvement and sends a direct message to the bot:

**User (Telegram):**
> 📲 *"Hey PM, I've noticed user sessions are not expiring correctly in the local database. Please check the code in `pkg/session/manager.go` and ensure the TTL is applied. Once fixed, run the tests and let me know when it's on GitHub."*

### 2. Core Processing (Project Manager - GPT-4o)
PicoClaw receives the Telegram webhook. The **Project Manager** (PM) springs into action:
- **Intent Analysis**: The PM identifies three tasks: Analysis/Fix, Testing, and Deployment.
- **Dynamic Delegation**: The PM doesn't do it all alone. It sends technical instructions to the Senior Dev.

**PM ➔ Senior Dev (DeepSeek Reasoner):**
> 🤖 *"Critical Task: User reports session TTL failures. Analyze `pkg/session/manager.go`, implement the fix, and delegate verification of unit tests to QA."*

### 3. Technical Execution (Senior Developer - DeepSeek Reasoner)
The **Senior Dev** receives the technical context. Thanks to its deep reasoning model (`deepseek-reasoner`), it analyzes the time logic in Go:
- Reads the file: `read_file("pkg/session/manager.go")`.
- Finds that `time.After` is not being reset correctly.
- Applies the fix dynamically: `edit_file(...)`.

### 4. Quality Cycle (QA Engineer - Claude 3.5 Sonnet)
The Senior Dev proactively invokes the **QA Engineer**:
- **QA** runs: `exec("go test ./pkg/session/...")`.
- If tests fail, QA reports back to the Senior Dev for corrections. If they pass, QA proceeds.
- **QA ➔ GitHub**: Executes git commands to push the solution.

### 5. Closing the Loop: Feedback to User (Telegram)
Once the hierarchy finishes, the PM consolidates the final report and sends it back to the user's Telegram chat.

**PM (Telegram):**
> ✅ *"Done! The Senior Dev has fixed the TTL handling bug in `manager.go`. QA has verified the tests (100% pass), and the changes are now in the main branch. Need anything else?"*

## Why it's "Dynamic"

1. **Context Injection**: The system reads the *current* real-time state of your local repository when you send the message.
2. **Autonomous Hierarchy**: The user only talks to the PM; the internal "conversation" between Senior Dev, QA, and Junior Dev happens without human intervention.
3. **Intelligent Multi-Modeling**: Different LLMs are triggered based on the specific complexity of the subtasks arising from your original message.

## Benefits of this Configuration

1. **Cost Efficiency**: Use cheaper models (`DeepSeek Chat`) for simple tasks and reserve powerful models (`DeepSeek Reasoner`, `GPT-4o`) for critical analysis.
2. **Specialization**: Each agent has a clear role, preventing a single model from getting overwhelmed with too much irrelevant information.
3. **Parallelism**: The QA can be verifying one part of the code while the Senior Dev analyzes the next step, accelerating the fix lifecycle.

---

## Use Case 2: Automated Social Media Post Generation and Approval (Social Bundle)

This flow demonstrates the use of the **Queue & Batch System** and **Fast-path Slash Commands** for efficient content management.

### 1. Initial Instruction (Telegram/Discord/CLI)
The user asks the agent to generate social media content:

**User:**
> 📲 *"Generate a post for Facebook and Instagram about our new v2.5 release. I want it to include a professional image and persuasive text."*

### 2. Background Processing (Social Manager)
The **Project Manager** delegates the task to the **Social Manager**:
- **Social Manager** initiates the `social_post_bundle` tool.
- The system immediately returns a tracking ID: `⏳ Task started: #IMA_GEN_02_03_26_1600. I'll notify you when ready.`
- The LLM is freed, saving tokens while the image (via DALL-E/Gemini) and text are generated.

### 3. Delivery Notification
Once the batch is complete, PicoClaw sends the result to the user with the final text and attached image.

**PicoClaw (Notification):**
> 🎨 **Batch Generated: #IMA_GEN_02_03_26_1600**
> [Attached Image]
> *Suggested text: "Great news! PicoClaw v2.5 is here..."*
> 💡 **Options (Copy and paste):**
> 1) `/bundle_approve id=20260302_1600`
> 2) `/bundle_regen id=20260302_1600`
> 3) `/bundle_edit id=20260302_1600`

### 4. Instant Approval (Fast-path)
The user reviews and decides to approve from their mobile or terminal:

**User (Telegram/Discord):**
> `/bundle_approve id=20260302_1600`

**Result:**
- The system intercepts the `/` command (Fast-path).
- **Instant processing**: Without consulting the AI, the system marks the batch as approved and publishes it directly to Facebook/Instagram.
- The user receives immediate confirmation: `✅ Bundle approved and published successfully.`

## Advantages of the Command System
- **Zero Latency**: Response to `/` commands is immediate.
- **Efficiency**: No reasoning tokens consumed for a simple approval.
- **Omnichannel**: Works identically on Discord (Slash Commands), Telegram, and Terminal.

---

## Use Case 3: Multi-Agent Image Generation with Global Tracker (v3.4.1+)

With v3.4.1, the **Global Tracker** ensures perfect consistency across multi-agent workflows:

### Scenario: Subagent generates, Main Agent publishes

1. **User requests content**:
   ```
   @picoclaw spawn task='Generate image about AI and create Twitter post'
   ```

2. **Subagent works**:
   - Generates image with `image_gen_create`
   - Creates post with `community_manager_create_draft`
   - Saves to **Global Shared Workspace**

3. **Main Agent can immediately access**:
   - No "ID not found" errors
   - Instant access to subagent's work
   - Can approve and publish without delays

### Benefits

- ✅ **Shared State**: All agents access the same workspace
- ✅ **No Sync Issues**: Changes are immediately visible
- ✅ **Scalable**: Works with any number of subagents
