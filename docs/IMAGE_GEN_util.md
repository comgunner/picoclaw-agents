# Image Generation Util

Quick guide to use image generation tools in PicoClaw from terminal and Telegram.

> **PicoClaw v3.4.1**: Features **Fast-path Slash Commands** for instant bundle management and **Global Tracker** for multi-agent consistency.
>
> **⚠️ IMPORTANT: THIS IS FOR GENERATING STATIC IMAGES, NOT VIDEOS**
>
> - `text_script_create`: Generates **POST TEXT** for Facebook/Twitter (like post copy)
> - `image_gen_create`: Generates **STATIC IMAGES** from text
> - **NO video generation** in this tool

---

## Requirements

Configure your credentials in `~/.picoclaw/config.json`:

```json
{
  "tools": {
    "image_gen": {
      "provider": "gemini",
      "gemini_api_key": "YOUR_GEMINI_API_KEY",
      "gemini_image_model": "gemini-2.5-flash-image-preview",
      "ideogram_api_key": "",
      "ideogram_api_url": "https://api.ideogram.ai/v1/ideogram-v3/generate",
      "aspect_ratio": "4:5",
      "output_dir": "./workspace/image_gen"
    }
  }
}
```

Or use environment variables:

```bash
# Provider
export PICOCLAW_TOOLS_IMAGE_GEN_PROVIDER="gemini"

# Gemini (image-capable model)
export PICOCLAW_TOOLS_IMAGE_GEN_GEMINI_API_KEY="your_api_key"
export PICOCLAW_TOOLS_IMAGE_GEN_GEMINI_IMAGE_MODEL="gemini-2.5-flash-image-preview"

# Ideogram
export PICOCLAW_TOOLS_IMAGE_GEN_IDEOGRAM_API_KEY="your_api_key"
export PICOCLAW_TOOLS_IMAGE_GEN_IDEOGRAM_API_URL="https://api.ideogram.ai/v1/ideogram-v3/generate"

# Aspect Ratio
export PICOCLAW_TOOLS_IMAGE_GEN_ASPECT_RATIO="4:5"

# Output directory
export PICOCLAW_TOOLS_IMAGE_GEN_OUTPUT_DIR="./workspace/image_gen"
```

### Config Priority

- Environment variables override `~/.picoclaw/config.json`.
- The valid variable name is `PICOCLAW_TOOLS_IMAGE_GEN_GEMINI_IMAGE_MODEL`.
- `GEMINI_IMAGEN_MODEL` is not used by PicoClaw.
- `tools.image_gen.output_dir` relative paths are resolved from the agent workspace (`agents.defaults.workspace`), not from process `cwd`.
- After changing config/env, restart the PicoClaw process/service.

---

## Available Tools

### Routing Rules (Important)

- If the user asks only for an image (example: `Generate an image of...`), use only `image_gen_create`.
- Use `text_script_create` only when the user explicitly asks for script/post text, or asks for Script -> Image workflow.
- `prompt_base_img.txt` is used to build a visual prompt from an existing script (not for direct image-only requests).

### `text_script_create`

Generates **TEXT FOR SOCIAL MEDIA POSTS** (Facebook, Twitter, Discord).

**NOT for video** - this is the text that accompanies an image post.

**Parameters:**
- `topic` (required): Post topic
- `category` (optional): 'news', 'story', 'tutorial', 'announcement'
- `tone` (optional): 'professional', 'casual', 'engaging'
- `language` (optional): 'en', 'es' (auto-detected)

**Example:**
```bash
./picoclaw agent -m "Use text_script_create with topic='Artificial Intelligence', category='news'"
```

**Output:** Facebook post text (max 1200 chars, viral style)

---

### `image_gen_create`

Generates **STATIC IMAGES** from text prompt.

**Parameters:**
- `prompt` (required): Image description
- `provider` (optional): 'gemini' or 'ideogram'
- `aspect_ratio` (optional): '4:5', '16:9', '1:1'

**Technical Note:**
- **Gemini:** Uses `gemini_image_model` from config/env (default: `gemini-2.5-flash-image-preview`)
- If Gemini returns model `NOT_FOUND` and Ideogram is configured, PicoClaw falls back automatically to Ideogram
- Some Gemini image-capable models may still return square outputs in certain accounts/deployments; for strict `4:5/16:9`, prefer Ideogram
- **Ideogram:** V3 API (1 image by default)

**Example:**
```bash
./picoclaw agent -m "Use image_gen_create with prompt='Cinematic sunset over mountains'"
```

**Output:** Static JPG image

---

### `image_gen_workflow`

Shows options after generating an image (publish to social, etc.).

**Parameters:**
- `image_path` (required): Path of generated image

---

### `script_to_image_workflow`

Creates post text + generates image based on topic.

**Parameters:**
- `topic` (required): Topic for post and image

**Example:**
```bash
./picoclaw agent -m "Use script_to_image_workflow with topic='Dragon story'"
```

**Output:** Post text + Static image

---

### `community_manager_create_draft`

Creates draft post for social media.

**Parameters:**
- `raw_data` (required): Technical content
- `platform` (required): 'discord', 'twitter', 'facebook', 'blog'

---

### `community_from_image`

Generates post text based on an image.

**Parameters:**
- `image_path` (required): Path of generated image
- `platform` (optional): 'discord', 'twitter', 'facebook'

---

## Usage Examples

### Terminal

```bash
# Generate TEXT FOR FACEBOOK POST
./picoclaw agent -m "Use text_script_create with topic='Artificial Intelligence', category='news'"

# Generate STATIC IMAGE
./picoclaw agent -m "Use image_gen_create with prompt='Cinematic sunset over mountains'"

# Generate with aspect ratio
./picoclaw agent -m "Use image_gen_create with prompt='Product photo', aspect_ratio='16:9'"

# Create POST TEXT + IMAGE
./picoclaw agent -m "Use script_to_image_workflow with topic='Dragon story'"

# Post-generation workflow
./picoclaw agent -m "Use image_gen_workflow with image_path='./workspace/image_gen/20260301_abc/20260301_abc.-imagen.jpg'"

# Generate post text from image
./picoclaw agent -m "Use community_from_image with image_path='./workspace/image_gen/test.jpg', platform='facebook'"
```

---

---

### Telegram

In Telegram, you can use natural language directly or explicit commands.

**Natural language examples:**
- "Generate an image of an astronaut cat on Mars"
- "Create a Facebook post about electric cars and its image"
- "Draw a cyber-punk landscape in 16:9"

**Interactive workflow:**
When you generate an image, PicoClaw will respond with the image and **interactive buttons**:
- `📖 View script`: Shows the full generated post text.
- `📱 Publish`: Opens the social media posting menu.
- `🔄 Regenerate`: Attempts to generate a different version.

**Multi-agent delegation (Recommended):**
- "@picoclaw subagent task='Create a horror story script and its image'"
- "@picoclaw spawn task='Generate astronomy content and publish it'"

---

### Discord

Discord supports the same features as Telegram, with the advantage of richer visualization in specific channels.

**Channel interaction:**
- `!agent Generate an image of an enchanted forest`
- `!agent Create a Twitter script about the new iPhone and generate its image`

**Rich Embeds:**
PicoClaw will send the image inside an **Embed**, allowing you to see metadata (topic, folder, ID) organized alongside **action buttons**.

**Quick commands:**
- `subagent task='Generate a robot image'`
- `image_gen_create prompt='Circular coffee logo' aspect_ratio='1:1'`

## File Structure

```
./workspace/image_gen/
├── tracker.json
├── 20260301_143022_abc123/
│   ├── 20260301_143022_abc123.-script.txt        # POST TEXT (Facebook/Twitter)
│   ├── 20260301_143022_abc123.-prompt_visual.txt # Prompt for image generation
│   └── 20260301_143022_abc123.-imagen.jpg        # STATIC IMAGE generated
└── ...
```

---

## Prompt Templates

### Custom Prompts

You can use your own prompts:

```json
{
  "tools": {
    "image_gen": {
      "image_script_path": "./workspace/prompt_base.txt",
      "image_gen_script_path": "./workspace/prompt_base_img.txt"
    }
  }
}
```

### Default Prompts

**DEFAULT_IMAGE_SCRIPT** (for post text):
- Viral Facebook style
- Max 1200 characters
- Structure: Hook → Development → Closing
- Optimized for engagement (comments, shares)
- Typical custom file: `./workspace/prompt_base.txt`

**DEFAULT_IMAGE_GEN_SCRIPT** (for visual prompt generation from a script):
- Hollywood storyboard format
- Output in English (better image quality)
- Includes: characters, atmosphere, lighting
- Negative prompt: no text, no watermarks
- Typical custom file: `./workspace/prompt_base_img.txt`

---

## Aspect Ratios

- `"4:5"` - Portrait (Instagram, Facebook)
- `"16:9"` - Landscape (YouTube, Twitter)
- `"1:1"` - Square (Instagram feed)
- `"9:16"` - Vertical (Stories, Reels, TikTok)

---

## Workflows

### Workflow 1: Facebook Post with Image

```bash
# Step 1: Generate POST TEXT
./picoclaw agent -m "Use text_script_create with topic='AI revolution'"

# Step 2: Generate IMAGE
./picoclaw agent -m "Use image_gen_create with prompt='Futuristic robot'"

# Step 3: Publish to Facebook
./picoclaw agent -m "Publish text and image to Facebook"
```

### Workflow 2: Direct Generation

```bash
# All in one
./picoclaw agent -m "Generate product image and post to Twitter with text"
```

### Workflow 3: Multi-Agent Delegation (v3.4.1+)

```bash
# Delegate full workflow to subagent
./picoclaw agent -m "spawn task='Generate image about AI, create Twitter post, and publish it' label='ai_campaign'"

# The subagent will:
# 1. Generate image with image_gen_create
# 2. Create post with community_manager_create_draft
# 3. Publish to Twitter with x_post_tweet
# 4. Report back when complete
```

The **Global Tracker** ensures the subagent's generated content is immediately available to the main agent.

---

## ⚡ Fast-path Slash Commands (v3.4.1+)

After receiving an image generation notification (e.g., `#IMA_GEN_...`), use quick commands:

```text
/bundle_approve id=20260302_161740_yiia22
/bundle_regen id=20260302_161740_yiia22
/bundle_edit id=20260302_161740_yiia22
/bundle_publish id=20260302_161740_yiia22 platforms=facebook,twitter
```

**Benefits:**
- ✅ **Zero latency**: No LLM reasoning, instant execution
- ✅ **Consistent syntax**: Works identically on Telegram, Discord, CLI
- ✅ **Safe**: ID validation before execution

See [docs/queue_batch.md](docs/queue_batch.md) for complete documentation.

---

## Rate Limits

| Provider | Limit      |
| -------- | ---------- |
| Gemini   | 60 req/min |
| Ideogram | 20 req/min |

---

## Troubleshooting

### "API key not configured"
Set `gemini_api_key` in config.json

### "Image generation failed"
Check prompt for prohibited content

### "Gemini model NOT_FOUND"
- Verify model in `~/.picoclaw/config.json` under `tools.image_gen.gemini_image_model`
- Verify runtime env does not override it:
  - `PICOCLAW_TOOLS_IMAGE_GEN_GEMINI_IMAGE_MODEL`
  - `PICOCLAW_TOOLS_IMAGE_GEN_PROVIDER`
- Restart PicoClaw after changes.

### "Multiple images generated"
System generates **1 single image** by default:
- Gemini: `sampleCount: 1`
- Ideogram V3: 1 image by default

---

## Technical Notes

### Gemini API

**Model:** `gemini-2.5-flash-image-preview` (default, configurable)

- ✅ Uses Gemini image-capable API path according to model type
- ✅ If unavailable in your account, automatic fallback to Ideogram (when configured)

### Ideogram API

**V3 API (Recommended):**
- Endpoint: `https://api.ideogram.ai/v1/ideogram-v3/generate`
- 1 image by default

---

## Summary

| Tool                       | Output       | Is it video? |
| -------------------------- | ------------ | ------------ |
| `text_script_create`       | POST TEXT    | ❌ NO         |
| `image_gen_create`         | STATIC IMAGE | ❌ NO         |
| `script_to_image_workflow` | TEXT + IMAGE | ❌ NO         |
| `community_from_image`     | POST TEXT    | ❌ NO         |

**ALL FOR SOCIAL MEDIA (Facebook, Twitter, Discord) - NO VIDEO**
