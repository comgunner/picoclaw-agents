import json
import os
import http.client
import sys

def get_qwen_token():
    home = os.path.expanduser("~")
    auth_path = os.path.join(home, ".picoclaw", "auth.json")
    if not os.path.exists(auth_path):
        print(f"Error: No auth file found at {auth_path}")
        return None

    try:
        with open(auth_path, 'r') as f:
            store = json.load(f)
            creds = store.get("credentials", {})
            qwen_cred = creds.get("qwen")
            if not qwen_cred:
                print("Error: Qwen credentials not found. Run 'picoclaw auth login --provider qwen' first.")
                return None
            return qwen_cred.get("access_token")
    except Exception as e:
        print(f"Error reading auth file: {e}")
        return None

def list_models(token):
    # Try multiple regions to see where models are available
    regions = [
        "dashscope-us.aliyuncs.com",
        "dashscope-intl.aliyuncs.com",
        "dashscope.aliyuncs.com"
    ]

    for host in regions:
        print(f"\n--- Checking region: {host} ---")
        conn = http.client.HTTPSConnection(host)
        headers = {
            'Authorization': f'Bearer {token}',
            'Content-Type': 'application/json'
        }

        try:
            conn.request("GET", "/compatible-mode/v1/models", headers=headers)
            res = conn.getresponse()
            data = res.read()

            if res.status == 200:
                models_data = json.loads(data)
                models = models_data.get("data", [])
                if not models:
                    print("No models returned in this region.")
                else:
                    print(f"Available models ({len(models)}):")
                    for m in models:
                        print(f" - {m['id']}")
            else:
                print(f"Error: {res.status} {res.reason}")
                print(data.decode())
        except Exception as e:
            print(f"Connection error: {e}")
        finally:
            conn.close()

if __name__ == "__main__":
    token = get_qwen_token()
    if token:
        # Mask token for output
        masked = token[:6] + "..." + token[-4:] if len(token) > 10 else "***"
        print(f"Using token: {masked}")
        list_models(token)
