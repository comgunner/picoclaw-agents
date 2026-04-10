# Image Generation Util

Quick guide to use image generation tools in PicoClaw from terminal and Telegram.

> **PicoClaw v3.5.0**: Now supports **Antigravity OAuth** image generation with `gemini-3.1-flash-image` — **no API key needed, completely FREE!** Just login with your Google account. Includes mandatory **150s cooldown** for anti-ban protection. See [ANTIGRAVITY_IMAGE_GEN.md](./ANTIGRAVITY_IMAGE_GEN.md).
>
> **PicoClaw v3.4.1**: Features **Fast-path Slash Commands** for instant bundle management and **Global Tracker** for multi-agent consistency.
>
> **⚠️ IMPORTANT: THIS IS FOR GENERATING STATIC IMAGES, NOT VIDEOS**
>
> - `text_script_create`: Generates **POST TEXT** for Facebook/Twitter (like post copy)
> - `image_gen_create`: Generates **STATIC IMAGES** from text
> - `image_gen_antigravity`: Generates images via **Antigravity OAuth** (default, no API key, FREE)
> - **NO video generation** in this tool

---

## 💸 Free Image Generation with Antigravity OAuth

**Images generated via Antigravity OAuth do NOT require an API key and do NOT cost you a cent.** They use your Google account's free quota through the Cloud Code Assist service. No billing setup, no credit card, no charges.

| Method | Cost | API Key Required? |
|--------|------|-------------------|
| **Antigravity OAuth** (default) | ✅ **FREE** | ❌ No — just Google login |
| Gemini API | 💰 Per-token billing | ✅ Yes — requires billing in Google Cloud |
| Ideogram API | 💰 Per-image billing | ✅ Yes — requires paid plan |

**Recommendation:** Use Antigravity OAuth as your default. It's free and works out of the box.

---

## 🆕 Antigravity OAuth Image Generation (Recommended — Default — FREE)

Since v3.5.0, the **default** method uses **Google Antigravity** via OAuth:
- **No API key required** — just login with your Google account
- **Model:** `gemini-3.1-flash-image`
- **Cost:** ✅ **FREE** — no API key, no billing, no charges
- **Cooldown:** 150s (2.5 min) mandatory after each generation (anti-ban)

### Setup

```bash
picoclaw auth login --provider google-antigravity
```

This automatically configures:
- `provider: "antigravity"`
- `antigravity_model: "gemini-3.1-flash-image"`
- `cooldown_seconds: 150`

---

## Configuration Options

### Option A: Antigravity OAuth (Default — Recommended — FREE)

**No API key needed. No cost.** Just OAuth login with your Google account.

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

**Login command:**
```bash
picoclaw auth login --provider google-antigravity
```

---

### Option B: Gemini API Key (Fallback — Paid)

Use if you have a Gemini API key and prefer direct API access. **Requires billing in Google Cloud Console.**

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

**Cost:** Per-token billing via Google Cloud Console.

---

### Option C: Ideogram API Key (Fallback — Paid)

Use if you have an Ideogram API key. **Requires paid Ideogram plan.**

```json
{
  "tools": {
    "image_gen": {
      "provider": "ideogram",
      "ideogram_api_key": "YOUR_IDEOGRAM_API_KEY",
      "ideogram_api_url": "https://api.ideogram.ai/v1/ideogram-v3/generate",
      "ideogram_model_name": "V_3_TURBO",
      "ideogram_aspect_ratio": "4x5",
      "ideogram_rendering_speed": "TURBO",
      "ideogram_style_type": "REALISTIC",
      "aspect_ratio": "1:1",
      "output_dir": "./workspace/image_gen"
    }
  }
}
```

**Cost:** Per-image billing via Ideogram subscription.

---

### Complete Configuration (All Providers)

```json
{
  "tools": {
    "image_gen": {
      "provider": "antigravity",
      "antigravity_model": "gemini-3.1-flash-image",
      "cooldown_seconds": 150,
      "gemini_api_key": "",
      "gemini_text_model_name": "gemini-3-flash-agent",
      "gemini_image_model_name": "gemini-2.5-flash-image",
      "ideogram_api_key": "",
      "ideogram_api_url": "https://api.ideogram.ai/v1/ideogram-v3/generate",
      "aspect_ratio": "1:1",
      "output_dir": "./workspace/image_gen"
    }
  }
}
```

**Priority order:** Antigravity OAuth (default, FREE) → Gemini API key (paid) → Ideogram API key (paid)

---

## Usage Examples

### Example 1: Generate an image (FREE via OAuth)

**User:** `Generate an image of a cute cat wearing sunglasses`

**What happens:** Agent calls `image_gen_antigravity` → generates image via OAuth (FREE) → sends as photo attachment to Telegram/Discord.

### Example 2: Facebook post with image (FREE via OAuth)

**User (English):** `Generate a Facebook post with image about nuclear danger and doomsday clock, attach the image`

**User (Spanish):** `genera un post para facebook con imagen sobre peligro nuclear y reloj del juicio final adjunta la imagen`

**What happens:** Agent calls `social_post_bundle` → generates text script → generates image via Antigravity OAuth (FREE) → copies image to bundle directory → sends post with image attachment.

### Example 3: Simple image generation (FREE via OAuth)

**User (English):** `Generate an image of a bird with sunglasses, Matrix style`

**User (Spanish):** `genera una imagen de un pajaro con lentes de sol estilo matrix`

**What happens:** Agent calls `image_gen_antigravity` → generates image via OAuth (FREE) → sends as photo attachment.

---

## Available Tools

### Routing Rules (Important)

- If the user asks only for an image (example: `Generate an image of...`), use only `image_gen_create` or `image_gen_antigravity`.
- Use `text_script_create` only when the user explicitly asks for script/post text, or asks for Script → Image workflow.
- `prompt_base_img.txt` is used to build a visual prompt from an existing script (not for direct image-only requests).

### `text_script_create`

Generates **TEXT FOR SOCIAL MEDIA POSTS** (Facebook, Twitter, Discord).

**NOT for video** - this is the text that accompanies an image post.

### `image_gen_antigravity`

Generates images via **Antigravity OAuth** (default). **FREE — no API key needed.**

### `image_gen_create`

Generates images via **Gemini API** or **Ideogram API** (fallback if OAuth not configured).

### `image_gen_workflow`

Script-to-image workflow. Generates a text script first, then creates a matching image.

---

## Aspect Ratios

| Ratio | Use Case |
|-------|----------|
| `1:1` | Default, square |
| `16:9` | Widescreen, desktop |
| `9:16` | Stories, mobile vertical |
| `4:5` | Instagram portrait |
| `3:4` | Classic portrait |

---

## Troubleshooting

### "No Antigravity OAuth credentials"

Run the login command:
```bash
picoclaw auth login --provider google-antigravity
```

### "Rate Limited (429)"

Antigravity has a shared quota for all API calls (chat + images). Wait ~5 minutes and retry. The system has automatic retry with exponential backoff (5s → 15s → 30s → 60s → 120s).

### "Method doesn't allow unregistered callers" (403)

This means you're trying to use Gemini API without an API key. Switch to Antigravity OAuth (default) which is FREE and requires no API key.

### "Tool definitions exceed budget threshold"

Set `context_window: 128000` in `~/.picoclaw/config.json` and restart the gateway.

### "no such file or directory" for the image

This is a workspace mismatch bug — fixed in v3.5.0 with automatic `copyFile`. Make sure you're using the latest binary.

---

## References

- [ANTIGRAVITY_IMAGE_GEN.md](./ANTIGRAVITY_IMAGE_GEN.md) — Full Antigravity OAuth guide
- [SOCIAL_MEDIA.md](./SOCIAL_MEDIA.md) — Social media publishing
- [SOCIAL_MEDIA_util.md](./SOCIAL_MEDIA_util.md) — Social media utilities
