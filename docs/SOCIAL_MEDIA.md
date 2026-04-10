# Social Media Integration

PicoClaw includes native tools for Facebook, X (Twitter), and Discord posting.

> **PicoClaw v3.5.0**: Now supports **Antigravity OAuth** image generation with `gemini-3.1-flash-image` — no API key needed! `social_post_bundle` now generates images via OAuth by default. See [IMAGE_GEN_util.md](./IMAGE_GEN_util.md) and [ANTIGRAVITY_IMAGE_GEN.md](./ANTIGRAVITY_IMAGE_GEN.md).
>
> **PicoClaw v3.4.1**: Features **Fast-path Slash Commands** for instant bundle management and **Global Tracker** for multi-agent consistency.

## Supported Tools

- `facebook_post`: Publish to Facebook Page (text-only or image + text, optional comment)
- `x_post_tweet`: Publish tweet (text-only or image, optional reply)
- `discord_post`: Publish to Discord webhook (text-only or image, optional username)

## Facebook Token Model (Important)

Facebook no longer supports `publish_actions`.  
Use a **Page Access Token** with modern Page permissions:

- `pages_manage_posts`
- `pages_read_engagement`
- `pages_show_list`
- optional for moderation: `pages_manage_engagement`

PicoClaw supports two Facebook modes:

1. Direct Page token
2. Auto-refresh Page token (if token expires with code `190`) using:
  - `app_id`
  - `app_secret`
  - `user_token`

## Configuration

Update `~/.picoclaw/config.json`:

```json
{
  "tools": {
    "social_media": {
      "facebook": {
        "default_page_id": "YOUR_FB_PAGE_ID",
        "default_page_token": "YOUR_FB_PAGE_TOKEN",
        "app_id": "YOUR_FB_APP_ID",
        "app_secret": "YOUR_FB_APP_SECRET",
        "user_token": "YOUR_FB_USER_TOKEN"
      },
      "x": {
        "api_key": "YOUR_X_API_KEY",
        "api_secret": "YOUR_X_API_SECRET",
        "access_token": "YOUR_X_ACCESS_TOKEN",
        "access_token_secret": "YOUR_X_ACCESS_TOKEN_SECRET"
      },
      "discord": {
        "webhook_url": "https://discord.com/api/webhooks/YOUR_WEBHOOK_ID/YOUR_WEBHOOK_TOKEN"
      }
    }
  }
}
```

## Image Generation Configuration

For generating images with posts (`social_post_bundle`), configure the image provider:

### Option A: Antigravity OAuth (Default — FREE — No API Key)

Images generated via Antigravity OAuth **do not require an API key and do NOT cost you a cent**. They use your Google account's free quota.

```json
{
  "tools": {
    "image_gen": {
      "provider": "antigravity",
      "antigravity_model": "gemini-3.1-flash-image",
      "cooldown_seconds": 150,
      "aspect_ratio": "1:1",
      "output_dir": "./workspace/image_gen"
    }
  }
}
```

**Login required:**
```bash
picoclaw auth login --provider google-antigravity
```

### Option B: Gemini API Key (Fallback — Paid)

```json
{
  "tools": {
    "image_gen": {
      "provider": "gemini",
      "gemini_api_key": "YOUR_GEMINI_API_KEY",
      "gemini_text_model_name": "gemini-3-flash-agent",
      "gemini_image_model_name": "gemini-2.5-flash-image",
      "aspect_ratio": "1:1",
      "output_dir": "./workspace/image_gen"
    }
  }
}
```

### Option C: Ideogram API Key (Fallback — Paid)

```json
{
  "tools": {
    "image_gen": {
      "provider": "ideogram",
      "ideogram_api_key": "YOUR_IDEOGRAM_API_KEY",
      "ideogram_api_url": "https://api.ideogram.ai/v1/ideogram-v3/generate",
      "ideogram_model_name": "V_3_TURBO",
      "aspect_ratio": "1:1",
      "output_dir": "./workspace/image_gen"
    }
  }
}
```

**Priority order:** Antigravity OAuth (default, FREE) → Gemini API key (paid) → Ideogram API key (paid)

## Environment Variables

```bash
# Facebook
export PICOCLAW_TOOLS_SOCIAL_FACEBOOK_DEFAULT_PAGE_ID="your_page_id"
export PICOCLAW_TOOLS_SOCIAL_FACEBOOK_DEFAULT_PAGE_TOKEN="your_page_token"
export PICOCLAW_TOOLS_SOCIAL_FACEBOOK_APP_ID="your_app_id"
export PICOCLAW_TOOLS_SOCIAL_FACEBOOK_APP_SECRET="your_app_secret"
export PICOCLAW_TOOLS_SOCIAL_FACEBOOK_USER_TOKEN="your_user_token"

# X
export PICOCLAW_TOOLS_SOCIAL_X_API_KEY="your_api_key"
export PICOCLAW_TOOLS_SOCIAL_X_API_SECRET="your_api_secret"
export PICOCLAW_TOOLS_SOCIAL_X_ACCESS_TOKEN="your_access_token"
export PICOCLAW_TOOLS_SOCIAL_X_ACCESS_TOKEN_SECRET="your_access_token_secret"

# Discord
export DISCORD_WEBHOOK_URL="your_webhook_url"
```

## Usage Examples

### Generate Facebook Post with Image (via social_post_bundle)

**User (English):** `Generate a Facebook post with image about nuclear danger and doomsday clock, attach the image`

**User (Spanish):** `genera un post para facebook con imagen sobre peligro nuclear y reloj del juicio final adjunta la imagen`

**What happens:**
1. Agent calls `social_post_bundle` → generates text script via Antigravity OAuth
2. Generates visual prompt from script
3. Generates image via `image_gen_antigravity` (OAuth, no API key)
4. Copies image to bundle output directory
5. Sends post with image attachment to Telegram/Discord

### Generate Simple Image

**User (English):** `Generate an image of a bird with sunglasses, Matrix style`

**User (Spanish):** `genera una imagen de un pajaro con lentes de sol estilo matrix`

**What happens:**
1. Agent calls `image_gen_antigravity` → generates image via OAuth
2. Sends image as photo attachment to Telegram/Discord

### Terminal Usage

```bash
# Facebook text-only
./picoclaw-agents agent -m "Use facebook_post with message='Hello from PicoClaw'"

# Facebook image + text
./picoclaw-agents agent -m "Use facebook_post with message='Launch update', image_path='/tmp/post.jpg'"

# Facebook image + text + comment
./picoclaw-agents agent -m "Use facebook_post with message='Main update', image_path='/tmp/post.jpg', comment='Extra details'"

# Facebook multi-page override
./picoclaw-agents agent -m "Use facebook_post with page_id='123456789', page_token='EAAB...', message='Page-specific update'"

# X text-only
./picoclaw-agents agent -m "Use x_post_tweet with message='Hello X'"

# X with image
./picoclaw-agents agent -m "Use x_post_tweet with message='Check this out', image_path='/tmp/photo.jpg'"

# Discord text-only
./picoclaw-agents agent -m "Use discord_post with message='Hello Discord'"

# Discord with image
./picoclaw-agents agent -m "Use discord_post with message='Check this image', image_path='/tmp/photo.jpg'"

# Post to multiple platforms
./picoclaw-agents agent -m "Post 'Big announcement!' to Facebook, Twitter and Discord"
```

### Telegram Usage

Send messages directly to your bot (with `picoclaw-agents gateway` running):

```text
# Simple posts
Publica en Facebook: "¡Hola desde PicoClaw!"
Publica en Twitter: "Nuevo lanzamiento #PicoClaw"
Publica en Discord: "Anuncio importante"

# With images
Publica en Facebook la imagen /tmp/foto.jpg con mensaje "¡Nuevo producto!"
Publica en Twitter la imagen /tmp/photo.jpg con texto "Mirá esto"

# Multi-platform
Publica en Facebook y Twitter: "Gran noticia hoy"
Publica en todas las redes: "¡Anuncio importante!"
```

### Discord Usage

Send messages to your Discord bot or via commands:

```text
# Direct messages to bot
Post to Facebook: "Hello from our community!"
Post to Twitter: "New feature released #update"
Post image /path/to/image.jpg to Facebook with message "Check this out"

# Multi-platform
Post to all social media: "Major announcement!"
```

### Combined with Image Generation

```bash
# Generate image and post
./picoclaw-agents agent -m "Generate image of new product and post to Facebook with attractive text"

# Full workflow
./picoclaw-agents agent -m "Use script_to_image_workflow with topic='Product launch', then post to Twitter"

# Community manager integration
./picoclaw-agents agent -m "Generate image, create community manager draft for Discord, then publish"
```

### Community Manager Examples

```bash
# Create draft from technical content
./picoclaw-agents agent -m "Use community_manager_create_draft with raw_data='New API endpoints released', platform='discord'"

# Generate text from image
./picoclaw-agents agent -m "Use community_from_image with image_path='./workspace/image_gen/abc/abc.-imagen.jpg', platform='twitter'"

# Full workflow
./picoclaw-agents agent -m "Generate image, create engaging post with community_manager, publish to Facebook"
```

## Notes

- If Facebook comment posting fails with code `368`, PicoClaw merges comment content into the post body.
- If Facebook returns token expiration code `190` and refresh fields are configured, PicoClaw retries with a refreshed Page token.
- For full practical examples, see `docs/SOCIAL_MEDIA_util.md`.

---

## ⚡ Fast-path Slash Commands (v3.4.1+)

After receiving a batch completion notification (e.g., `#IMA_GEN_...` or `#SOCIAL_...`), use quick commands for instant management:

### Bundle Management Commands

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

### Utility Commands

```text
/list pending          # Show all pending tasks
/status                # Show token usage and system status
/help                  # Show interactive help
/show model            # Show active model
/show channel          # Show active channel
```

### Global Tracker (v3.4.1+)

The **Global ImageGenTracker** is shared across all agents:
- **Subagent generates content** → **Main Agent can immediately publish**
- **No "ID not found" errors** across agent boundaries
- **Perfect consistency** in multi-agent workflows

See [docs/queue_batch.md](docs/queue_batch.md) for complete documentation.

---

## 🤖 Multi-Agent Workflows

With v3.4.1, you can delegate complete social media workflows to subagents:

```bash
# Delegate full workflow to subagent
./picoclaw-agents agent -m "spawn task='Generate image about AI, create post for Twitter, and publish it' label='social_campaign'"

# The subagent will:
# 1. Generate image with image_gen_create
# 2. Create post with community_manager_create_draft
# 3. Publish to Twitter with x_post_tweet
# 4. Report back when complete
```

The Global Tracker ensures the subagent's generated content is immediately available to the main agent for approval and publishing.
