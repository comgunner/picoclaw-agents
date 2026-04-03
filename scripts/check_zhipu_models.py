#!/usr/bin/env python3
"""
Check Zhipu AI (z.ai) Models

This script checks which models are available with your Zhipu AI API key.
It queries the Zhipu API to list all models you have access to.

Usage:
    python check_zhipu_models.py

Requirements:
    - Python 3.6+
    - Zhipu AI API key from https://platform.z.ai/api-keys

The script will automatically read your API key from ~/.picoclaw/auth.json
"""

import json
import os
import http.client
import sys

def get_zhipu_token():
    """
    Read Zhipu AI API token from ~/.picoclaw/auth.json

    Returns:
        str: API token or None if not found
    """
    home = os.path.expanduser("~")
    auth_path = os.path.join(home, ".picoclaw", "auth.json")
    if not os.path.exists(auth_path):
        print(f"Error: No auth file found at {auth_path}")
        print("Run './build/picoclaw-agents auth login --provider zhipu' first.")
        return None

    try:
        with open(auth_path, 'r') as f:
            store = json.load(f)
            creds = store.get("credentials", {})
            zhipu_cred = creds.get("zhipu")
            if not zhipu_cred:
                print("Error: Zhipu credentials not found.")
                print("Run './build/picoclaw-agents auth login --provider zhipu' first.")
                return None

            # Zhipu uses access_token field
            token = zhipu_cred.get("access_token")
            if not token:
                print("Error: No access token found in Zhipu credentials.")
                return None

            return token
    except FileNotFoundError:
        print(f"Error: Auth file not found at {auth_path}")
        return None
    except json.JSONDecodeError as e:
        print(f"Error: Invalid JSON in auth file: {e}")
        return None
    except Exception as e:
        print(f"Error reading auth file: {e}")
        return None

def list_models(token):
    """
    Query Zhipu API to list available models

    Args:
        token (str): Zhipu AI API key
    """
    # Zhipu API endpoint
    host = "api.z.ai"
    endpoint = "/api/paas/v4/models"

    print(f"\n--- Checking Zhipu AI API: {host} ---")
    conn = http.client.HTTPSConnection(host)
    headers = {
        'Authorization': f'Bearer {token}',
        'Content-Type': 'application/json'
    }

    try:
        conn.request("GET", endpoint, headers=headers)
        res = conn.getresponse()
        data = res.read()

        if res.status == 200:
            models_data = json.loads(data)
            models = models_data.get("data", [])

            if not models:
                print("No models returned from API.")
            else:
                print(f"\n✅ Available models ({len(models)}):")
                print("=" * 60)

                # Categorize models
                text_models = []
                vision_models = []
                other_models = []

                for m in models:
                    model_id = m.get('id', 'unknown')
                    if 'vl' in model_id.lower() or 'vision' in model_id.lower():
                        vision_models.append(model_id)
                    elif 'long' in model_id.lower():
                        other_models.append(model_id)
                    else:
                        text_models.append(model_id)

                # Print text models
                if text_models:
                    print("\n📝 Text Models:")
                    for model in sorted(text_models):
                        print(f"   • {model}")

                # Print vision models
                if vision_models:
                    print("\n👁️ Vision-Language Models:")
                    for model in sorted(vision_models):
                        print(f"   • {model}")

                # Print other models
                if other_models:
                    print("\n📚 Other Models:")
                    for model in sorted(other_models):
                        print(f"   • {model}")

                print("=" * 60)

                # Show recommended models
                print("\n⭐ Recommended models:")
                print("   • glm-5         - Latest version (default)")
                print("   • glm-5-turbo   - Fastest")
                print("   • glm-5.1       - Enhanced v5")
                print("   • glm-4.5-air   - Lightweight")

        elif res.status == 401:
            print(f"\n❌ Error 401: Unauthorized")
            print("   Invalid API key. Please check your API key.")
            print("   Get a new one at: https://platform.z.ai/api-keys")
            print(f"\n   Response: {data.decode()}")
        elif res.status == 429:
            print(f"\n❌ Error 429: Rate Limit Exceeded")
            print("   Too many requests. Please wait a moment.")
            print(f"\n   Response: {data.decode()}")
        else:
            print(f"\n❌ Error: {res.status} {res.reason}")
            print(f"   Response: {data.decode()}")

    except http.client.HTTPException as e:
        print(f"\n❌ HTTP Error: {e}")
    except Exception as e:
        print(f"\n❌ Connection error: {e}")
        print("   Make sure you have internet connection.")
    finally:
        conn.close()

def check_free_tier():
    """Display information about Zhipu free tier"""
    print("\n" + "=" * 60)
    print("🆓 Zhipu AI Free Tier Information")
    print("=" * 60)
    print("   • 100% FREE - No credit card required")
    print("   • 60 requests/minute")
    print("   • 100,000 tokens/minute")
    print("   • 1,000,000 tokens/day")
    print("   • Context up to 1M tokens (glm-4-long)")
    print("=" * 60)

if __name__ == "__main__":
    print("=" * 60)
    print("🔍 Zhipu AI (z.ai) Model Checker")
    print("=" * 60)

    token = get_zhipu_token()
    if token:
        # Mask token for output
        masked = token[:8] + "..." + token[-4:] if len(token) > 12 else "***"
        print(f"\n✅ Using API key: {masked}")

        # Show free tier info
        check_free_tier()

        # List models
        list_models(token)

        # Show usage examples
        print("\n💡 Usage examples:")
        print("   ./build/picoclaw-agents agent -m \"Hello\" --model glm-5")
        print("   ./build/picoclaw-agents agent -m \"Hi\" --model glm-5-turbo")
        print("   ./build/picoclaw-agents agent -m \"Test\" --model glm-4.5")
        print("\n📚 Documentation: https://docs.z.ai/guides/overview/quick-start")
        print("🌐 Console: https://platform.z.ai/")
    else:
        print("\n❌ No API key found.")
        print("\n📝 Setup instructions:")
        print("   1. Get API key from: https://platform.z.ai/api-keys")
        print("   2. Run: ./build/picoclaw-agents auth login --provider zhipu")
        print("   3. Paste your API key when prompted")
        print("   4. Run this script again")
        sys.exit(1)
