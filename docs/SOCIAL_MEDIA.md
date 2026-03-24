# Social Media Integration

PicoClaw includes native tools for Facebook, X (Twitter), and Discord posting.

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

### Terminal Usage

```bash
# Facebook text-only
./picoclaw agent -m "Use facebook_post with message='Hello from PicoClaw'"

# Facebook image + text
./picoclaw agent -m "Use facebook_post with message='Launch update', image_path='/tmp/post.jpg'"

# Facebook image + text + comment
./picoclaw agent -m "Use facebook_post with message='Main update', image_path='/tmp/post.jpg', comment='Extra details'"

# Facebook multi-page override
./picoclaw agent -m "Use facebook_post with page_id='123456789', page_token='EAAB...', message='Page-specific update'"

# X text-only
./picoclaw agent -m "Use x_post_tweet with message='Hello X'"

# X with image
./picoclaw agent -m "Use x_post_tweet with message='Check this out', image_path='/tmp/photo.jpg'"

# Discord text-only
./picoclaw agent -m "Use discord_post with message='Hello Discord'"

# Discord with image
./picoclaw agent -m "Use discord_post with message='Check this image', image_path='/tmp/photo.jpg'"

# Post to multiple platforms
./picoclaw agent -m "Post 'Big announcement!' to Facebook, Twitter and Discord"
```

### Telegram Usage

Send messages directly to your bot (with `picoclaw gateway` running):

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
./picoclaw agent -m "Generate image of new product and post to Facebook with attractive text"

# Full workflow
./picoclaw agent -m "Use script_to_image_workflow with topic='Product launch', then post to Twitter"

# Community manager integration
./picoclaw agent -m "Generate image, create community manager draft for Discord, then publish"
```

### Community Manager Examples

```bash
# Create draft from technical content
./picoclaw agent -m "Use community_manager_create_draft with raw_data='New API endpoints released', platform='discord'"

# Generate text from image
./picoclaw agent -m "Use community_from_image with image_path='./workspace/image_gen/abc/abc.-imagen.jpg', platform='twitter'"

# Full workflow
./picoclaw agent -m "Generate image, create engaging post with community_manager, publish to Facebook"
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
./picoclaw agent -m "spawn task='Generate image about AI, create post for Twitter, and publish it' label='social_campaign'"

# The subagent will:
# 1. Generate image with image_gen_create
# 2. Create post with community_manager_create_draft
# 3. Publish to Twitter with x_post_tweet
# 4. Report back when complete
```

The Global Tracker ensures the subagent's generated content is immediately available to the main agent for approval and publishing.
