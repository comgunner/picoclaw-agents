# Social Media Util

Quick operational guide for PicoClaw social media tools.

> **PicoClaw v3.4.1**: Features **Fast-path Slash Commands** for instant bundle management and **Global Tracker** for multi-agent consistency.

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
