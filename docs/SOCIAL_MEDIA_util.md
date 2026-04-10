# Social Media Util

Quick operational guide for PicoClaw social media tools.

> **PicoClaw v3.5.0**: Now supports **Antigravity OAuth** image generation with `gemini-3.1-flash-image` — no API key needed! `social_post_bundle` now generates images via OAuth by default. See [IMAGE_GEN_util.md](./IMAGE_GEN_util.md) and [ANTIGRAVITY_IMAGE_GEN.md](./ANTIGRAVITY_IMAGE_GEN.md).
>
> **PicoClaw v3.4.1**: Features **Fast-path Slash Commands** for instant bundle management and **Global Tracker** for multi-agent consistency.

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

## Minimal Config

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

## Image Generation Config

For generating images with posts (`social_post_bundle`):

### Antigravity OAuth (Default — FREE — No API Key)

Images via Antigravity OAuth **do not require an API key and do NOT cost a cent**.

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

Login: `picoclaw auth login --provider google-antigravity`

### Gemini API Key (Fallback — Paid)

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

### Ideogram API Key (Fallback — Paid)

```json
{
  "tools": {
    "image_gen": {
      "provider": "ideogram",
      "ideogram_api_key": "YOUR_IDEOGRAM_API_KEY",
      "ideogram_api_url": "https://api.ideogram.ai/v1/ideogram-v3/generate",
      "aspect_ratio": "1:1",
      "output_dir": "./workspace/image_gen"
    }
  }
}
```

**Priority:** Antigravity OAuth (FREE) → Gemini API → Ideogram API

## Facebook Behavior

- `facebook_post` supports:
  - text-only post
  - image + text post
  - optional first comment
- If comment fails with code `368`, content is merged into the post body.
- If token expires with code `190` and `app_id/app_secret/user_token` are configured, PicoClaw refreshes and retries.

## CLI Examples

```bash
# Facebook text-only
./picoclaw-agents agent -m "Use facebook_post with message='Hello from PicoClaw'"

# Facebook image post
./picoclaw-agents agent -m "Use facebook_post with message='Launch update', image_path='/tmp/post.jpg'"

# Facebook image post + comment
./picoclaw-agents agent -m "Use facebook_post with message='Main update', image_path='/tmp/post.jpg', comment='Extra context'"

# X text-only
./picoclaw-agents agent -m "Use x_post_tweet with message='Hello X'"

# Discord text-only
./picoclaw-agents agent -m "Use discord_post with message='Hello Discord'"
```

## Telegram / Discord Channel Prompts

```text
Post to Facebook: "Hello from the bot"
Post to Facebook with image /tmp/post.jpg and message "New update"
Post to X: "Release live now"
Post to Discord: "Important announcement"
```

## Facebook Permissions

Use modern Page permissions, not `publish_actions`:

- `pages_manage_posts`
- `pages_read_engagement`
- `pages_show_list`
- optional: `pages_manage_engagement`

---

## ⚡ Fast-path Slash Commands (v3.4.1+)

After receiving a batch completion notification (e.g., `#IMA_GEN_...` or `#SOCIAL_...`), use quick commands for instant management:

### Bundle Commands

```text
/bundle_approve id=ID        # Approve batch and proceed to publication
/bundle_regen id=ID          # Request full batch regeneration (image + text)
/bundle_edit id=ID           # Edit batch text before approving
/bundle_publish id=ID platforms=facebook,twitter  # Publish to platforms
```

### Utility Commands

```text
/list pending          # Show all pending tasks
/status                # Show token usage and system status
/help                  # Show interactive help
/show model            # Show active model
/show channel          # Show active channel
```

**Benefits:**
- ✅ **Zero latency**: No LLM reasoning, instant execution
- ✅ **Consistent syntax**: Works identically on Telegram, Discord, CLI
- ✅ **Safe**: ID validation before execution

### Global Tracker (v3.4.1+)

The **Global ImageGenTracker** is shared across all agents:
- **Subagent generates content** → **Main Agent can immediately publish**
- **No "ID not found" errors** across agent boundaries

See [docs/queue_batch.md](docs/queue_batch.md) for complete documentation.
