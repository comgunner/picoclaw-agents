# Using Antigravity Provider in PicoClaw

This guide explains how to set up and use the **Antigravity** (Google Cloud Code Assist) provider in PicoClaw.

## All Supported Providers

| Provider             | Command                                                             | Auth Method                                                | Default Model                                       |
| -------------------- | ------------------------------------------------------------------- | ---------------------------------------------------------- | --------------------------------------------------- |
| `google-antigravity` | `./picoclaw auth antigravity`<br>or `--provider google-antigravity` | OAuth 2.0 + PKCE (browser)                                 | `gemini-flash` (`antigravity/gemini-3-flash`)       |
| `openai`             | `./picoclaw auth login --provider openai`                           | OAuth 2.0 + PKCE (browser)<br>`--device-code` for headless | `gpt-5.2` (`openai/gpt-5.2`)                        |
| `anthropic`          | `./picoclaw auth login --provider anthropic`                        | Paste API key (no OAuth)                                   | `claude-sonnet-4.6` (`anthropic/claude-sonnet-4.6`) |

> [!NOTE]
> `anthropic` uses a **static API key** from [console.anthropic.com](https://console.anthropic.com) — no browser login, no token expiry. `openai` and `google-antigravity` use OAuth with auto-refresh.

## Antigravity Prerequisites

1.  A Google account.
2.  Google Cloud Code Assist enabled. This is usually part of the **Gemini for Google Cloud** onboarding.
3.  **Plan Compatibility**: Antigravity is specifically designed to work with **Google One AI Premium** plans and **Google Workspace** Gemini add-ons.

> [!NOTE]
> This is **different** from the standard Google AI Studio (Gemini API). Instead of an API key, it uses OAuth authentication linked to your Google plan's credits and quotas.

## 1. Authentication

To authenticate with Antigravity, run:

```bash
./picoclaw auth antigravity
```

This opens a browser window for Google OAuth. After login, PicoClaw stores credentials in `~/.picoclaw/auth.json`.

### Manual Authentication (Headless/VPS)
If you are running on a server (Coolify/Docker) and cannot reach `localhost`, follow these steps:
1.  Run the command above.
2.  Copy the URL provided and open it in your local browser.
3.  Complete the login.
4.  Your browser will redirect to a `localhost:51121` URL (which will fail to load).
5.  **Copy that final URL** from your browser's address bar.
6.  **Paste it back into the terminal** where PicoClaw is waiting.

PicoClaw will extract the authorization code and complete the process automatically.

### Token Expiry & Auto-Refresh
The `access_token` expires in **1 hour** (Google OAuth standard). PicoClaw now handles refresh in 3 layers:
1. **Background daemon**: proactively refreshes every 20 min if <30 min remain
2. **On every request**: retries with refresh_token even if token already expired (not just pre-expiry)
3. **`auth models` command**: also recovers from expired tokens automatically

Manual re-auth (`./picoclaw auth antigravity`) is only needed if:
- You revoked access from `myaccount.google.com > Security > Apps with access`
- You changed your Google password
- The refresh_token itself has been inactive for 6+ months

## 2. Managing Models

### List Available Models
```bash
./picoclaw auth models
```

### Validated Working Models (as of 2026-03-04)

Use these exact `model_name` values in your config or with `--model`:

| model_name                   | Description                              |
| ---------------------------- | ---------------------------------------- |
| `antigravity-gemini-3-flash` | Fast, reliable — **recommended default** |
| `gemini-3-flash`             | Same as above (auto-resolves)            |
| `gemini-3-pro-high`          | High reasoning Gemini 3                  |
| `gemini-3.1-pro-high`        | High reasoning Gemini 3.1                |
| `gemini-3.1-flash-image`     | Fast, multimodal image support           |
| `gemini-2.5-pro`             | Gemini 2.5 Pro                           |
| `gemini-2.5-flash`           | Gemini 2.5 Flash                         |
| `gemini-2.5-flash-thinking`  | Flash with reasoning                     |
| `gemini-2.5-flash-lite`      | Lightweight model                        |
| `claude-sonnet-4-6`          | Fast Claude reasoning                    |
| `claude-opus-4-6-thinking`   | Top-tier Claude with thinking blocks     |
| `gpt-oss-120b-medium`        | Open-source GPT alternative              |

### Standard Gemini Models (via API Key)
If you are not using Antigravity Auth (OAuth) and prefer to use a direct API key from Google AI Studio, you can configure the provider by adding the explicit `gemini/` prefix to the `model` directive in your `config.json` to route directly to the Google API, and use any of these public models:

| Model Name                                | Display Name                  | Input Limit | Output Limit |
| :---------------------------------------- | :---------------------------- | :---------- | :----------- |
| `gemini-2.5-flash`                        | Gemini 2.5 Flash              | 1,048,576   | 65,536       |
| `gemini-2.5-pro`                          | Gemini 2.5 Pro                | 1,048,576   | 65,536       |
| `gemini-2.0-flash`                        | Gemini 2.0 Flash              | 1,048,576   | 8,192        |
| `gemini-2.5-flash-lite`                   | Gemini 2.5 Flash-Lite         | 1,048,576   | 65,536       |
| `gemini-2.5-flash-image`                  | Nano Banana (Flash Image)     | 32,768      | 32,768       |
| `gemini-3-pro-preview`                    | Gemini 3 Pro Preview          | 1,048,576   | 65,536       |
| `gemini-3-flash-preview`                  | Gemini 3 Flash Preview        | 1,048,576   | 65,536       |
| `gemini-3.1-pro-preview`                  | Gemini 3.1 Pro Preview        | 1,048,576   | 65,536       |
| `gemini-3.1-flash-lite-preview`           | Gemini 3.1 Flash Lite Preview | 1,048,576   | 65,536       |
| `gemini-3-pro-image-preview`              | Nano Banana Pro               | 131,072     | 32,768       |
| `gemini-3.1-flash-image-preview`          | Nano Banana 2                 | 65,536      | 65,536       |
| `deep-research-pro-preview-12-2025`       | Deep Research Pro Preview     | 131,072     | 65,536       |
| `gemini-2.5-computer-use-preview-10-2025` | Gemini Computer Use Preview   | 131,072     | 65,536       |
| `gemini-robotics-er-1.5-preview`          | Gemini Robotics-ER 1.5        | 1,048,576   | 65,536       |
| `imagen-4.0-generate-001`                 | Imagen 4                      | 480         | 8,192        |
| `imagen-4.0-ultra-generate-001`           | Imagen 4 Ultra                | 480         | 8,192        |
| `veo-3.1-generate-preview`                | Veo 3.1                       | 480         | 8,192        |
| `gemini-flash-latest`                     | Gemini Flash Latest           | 1,048,576   | 65,536       |
| `gemini-pro-latest`                       | Gemini Pro Latest             | 1,048,576   | 65,536       |

### Switch Models
```bash
# Using models verified by Antigravity Auth (OAuth)
./picoclaw agent -m "Hello" --model claude-opus-4-6-thinking
./picoclaw agent -m "Hello" --model antigravity-gemini-3-flash

# Using standard models via public API (Google AI Studio Key)
./picoclaw agent -m "Hello" --model gemini-2.5-flash
./picoclaw agent -m "Hello" --model gemini-3-pro-preview
```

## 3. Configuration

### Default Configuration (deepseek-chat)

The main `config.example.json` file uses **deepseek-chat** as the default model for all agents. This provides a consistent, cost-effective starting point for new users:

```bash
cp config/config.example.json ~/.picoclaw/config.json
# Add your DeepSeek API key in config.json
./picoclaw agent -m "Hello"
```

### Antigravity Configuration

The `config.example_antigravity.json` file is a ready-to-use config where all agents default to `antigravity-gemini-3-flash`. Use this if you prefer Google's Antigravity provider:

```bash
cp config/config.example_antigravity.json ~/.picoclaw/config.json
./picoclaw auth antigravity
./picoclaw agent -m "Hello"
```

### model_list Entries for Antigravity

Each Antigravity model requires `"auth_method": "oauth"` and no `api_base`:

```json
{
  "model_name": "antigravity-gemini-3-flash",
  "model": "antigravity/gemini-3-flash",
  "api_key": "",
  "auth_method": "oauth"
},
{
  "model_name": "claude-sonnet-4-6",
  "model": "antigravity/claude-sonnet-4-6",
  "api_key": "",
  "auth_method": "oauth"
},
{
  "model_name": "claude-opus-4-6-thinking",
  "model": "antigravity/claude-opus-4-6-thinking",
  "api_key": "",
  "auth_method": "oauth"
}
```

### Comparative Examples: Antigravity (OAuth) vs Gemini (API Key)

Depending on whether you use your Google Cloud account quota (Antigravity via OAuth) or your own API Key generated on Google AI Studio, the `model_list` definitions differ specifically in the prefix and the `auth_method` field. Here are 3 clear examples of how to configure the exact same models using the different methods:

#### 1. Gemini 2.5 Flash
```json
// Via Antigravity (OAuth, no API Key needed)
{
  "model_name": "ag-gemini-2.5-flash",
  "model": "antigravity/gemini-2.5-flash",
  "api_key": "",
  "auth_method": "oauth"
}

// Via Public API (Requires API Key)
{
  "model_name": "gemini-2.5-flash",
  "model": "gemini/gemini-2.5-flash",
  "api_key": "YOUR_GEMINI_API_KEY"
}
```

#### 2. Gemini 3 Flash
```json
// Via Antigravity (OAuth, no API Key needed)
{
  "model_name": "ag-gemini-3-flash",
  "model": "antigravity/gemini-3-flash",
  "api_key": "",
  "auth_method": "oauth"
}

// Via Public API (Requires API Key)
{
  "model_name": "gemini-3-flash-preview",
  "model": "gemini/gemini-3-flash-preview",
  "api_key": "YOUR_GEMINI_API_KEY"
}
```

#### 3. Gemini 2.5 Pro
```json
// Via Antigravity (OAuth, no API Key needed)
{
  "model_name": "ag-gemini-2.5-pro",
  "model": "antigravity/gemini-2.5-pro",
  "api_key": "",
  "auth_method": "oauth"
}

// Via Public API (Requires API Key)
{
  "model_name": "gemini-2.5-pro",
  "model": "gemini/gemini-2.5-pro",
  "api_key": "YOUR_GEMINI_API_KEY"
}
```

## 4. Real-world Usage (Coolify/Docker)

If you are deploying via Coolify or Docker:

1.  Log in locally first, then copy credentials to the server:
    ```bash
    scp ~/.picoclaw/auth.json user@your-server:~/.picoclaw/
    ```
2.  *Alternatively*, run `./picoclaw auth antigravity` directly on the server using the headless flow.

## 5. Troubleshooting

| Error                                                     | Cause                              | Fix                                                                            |
| --------------------------------------------------------- | ---------------------------------- | ------------------------------------------------------------------------------ |
| `403 PERMISSION_DENIED / ACCESS_TOKEN_SCOPE_INSUFFICIENT` | Token expired or revoked           | Run `./picoclaw auth antigravity` again                                        |
| `404 NOT_FOUND`                                           | Model alias not resolved correctly | Verify `model_list` entry has correct `model` field and `auth_method: "oauth"` |
| `401 invalid_api_key`                                     | Wrong provider used for model      | Check `model` field has `antigravity/` prefix, not an OpenAI-style key         |
| `429 Rate Limit`                                          | Quota hit                          | PicoClaw shows reset time; wait or switch to another model                     |
| Empty response                                            | Model restricted for project       | Try `antigravity-gemini-3-flash` or `gemini-2.5-flash`                         |

## 6. Model Routing Architecture

PicoClaw uses a 3-step pipeline for model resolution:

### Configuration Definitions
- **`model_name`**: Internal ALIAS — the friendly name you use (e.g. `antigravity-gemini-3-flash`)
- **`model`**: Routing instruction — must contain `provider/model-id` (e.g. `antigravity/gemini-3-flash`)

### The 3-Step Pipeline
1. **Memory Load**: On startup, PicoClaw reads `model_list` from `~/.picoclaw/config.json` into RAM. Changes require a full restart.
2. **Factory Routing**: The alias is looked up → the `model` field is split by `/` → the `antigravity` prefix selects the Antigravity provider.
3. **Prefix Sanitization**: Before calling the API, the provider strips all prefixes:
   - `antigravity/gemini-3-flash` → `gemini-3-flash` ✅
   - `antigravity-gemini-3-flash` → `gemini-3-flash` ✅ (dash prefix also handled)

> [!TIP]
> Both `antigravity/gemini-3-flash` (slash) and `antigravity-gemini-3-flash` (dash) are valid model values in `model_list`. The provider strips both correctly.
