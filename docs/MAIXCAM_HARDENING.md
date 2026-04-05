# PicoClaw MaixCam Hardening Guide

> **Last Updated:** March 24, 2026 | **Version:** v3.5.1+  
> **Status:** ✅ Production Ready

## Overview

**MaixCam Hardening** is a critical security enhancement introduced in **v3.5.1** that implements **HMAC-SHA256 authentication** for IoT device messages, protecting against man-in-the-middle (MITM) attacks and unauthorized device access.

### Why MaixCam Hardening Matters

**Before v3.5.1:**
- IoT messages sent in plaintext
- No authentication of device identity
- Vulnerable to MITM attacks on local network
- Any device could impersonate a trusted camera

**After v3.5.1:**
- All messages authenticated with HMAC-SHA256
- Device identity verified before processing
- MITM attacks prevented
- Unauthorized devices rejected

**Impact:** Prevents unauthorized IoT devices from injecting malicious commands or falsified sensor data into your AI agent.

---

## Table of Contents

- [How HMAC Authentication Works](#how-hmac-authentication-works)
- [Configuration](#configuration)
- [Generating HMAC Secrets](#generating-hmac-secrets)
- [Device Setup](#device-setup)
- [Message Format](#message-format)
- [Troubleshooting](#troubleshooting)
- [Security Best Practices](#security-best-practices)

---

## How HMAC Authentication Works

### HMAC-SHA256 Overview

**HMAC** (Hash-based Message Authentication Code) is a cryptographic technique that combines:
- A **secret key** (shared between device and server)
- A **hash function** (SHA-256)
- The **message data**

```
┌─────────────────────────────────────────────────────────────┐
│                    HMAC Authentication Flow                  │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  MaixCam Device                          PicoClaw Server    │
│       │                                       │             │
│       │  1. Create message                    │             │
│       │  2. Compute HMAC signature            │             │
│       │     HMAC = HMAC-SHA256(message, key)  │             │
│       │  3. Send message + HMAC ─────────────>│             │
│       │                                       │             │
│       │                              4. Compute HMAC        │
│       │                                 HMAC' = HMAC-SHA256 │
│       │                                       │             │
│       │                              5. Compare HMAC == HMAC'│
│       │                                       │             │
│       │  6. Accept/Reject <───────────────────│             │
│       │                                       │             │
└─────────────────────────────────────────────────────────────┘
```

### Security Properties

✅ **Protected:**
- **Message Integrity**: Any modification to the message invalidates the HMAC
- **Authentication**: Only devices with the secret key can generate valid HMACs
- **Replay Protection**: Timestamps prevent replay attacks (optional)

❌ **Not Protected:**
- **Confidentiality**: Messages are not encrypted (use HTTPS for that)
- **Key Exchange**: Secret must be shared securely out-of-band

---

## Configuration

### Enable HMAC in `config.json`

```json
{
  "channels": {
    "maixcam": {
      "enabled": true,
      "hmac_secret": "your-secret-key-change-in-production",
      "require_hmac": true,
      "port": 8080,
      "allowed_devices": ["maixcam_001", "maixcam_002"]
    }
  }
}
```

### Configuration Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `enabled` | boolean | `false` | Enable MaixCam channel |
| `hmac_secret` | string | (required) | Shared secret for HMAC |
| `require_hmac` | boolean | `true` | Require HMAC on all messages |
| `port` | integer | `8080` | Port to listen on |
| `allowed_devices` | array | `[]` | List of allowed device IDs |

### Environment Variables

```bash
# MaixCam HMAC configuration
export PICOCLAW_CHANNELS_MAIXCAM_ENABLED=true
export PICOCLAW_CHANNELS_MAIXCAM_SECRET="your-secret-key-change-in-production"
export PICOCLAW_CHANNELS_MAIXCAM_REQUIRE_HMAC=true
export PICOCLAW_CHANNELS_MAIXCAM_PORT=8080
```

### Recommended Settings

| Environment | HMAC Secret | Require HMAC | Allowed Devices |
|-------------|-------------|--------------|-----------------|
| **Development** | `dev-secret-change-me` | `false` (for testing) | `[]` (all) |
| **Production** | `[32+ random bytes]` | `true` | Specific device IDs |
| **High Security** | `[64+ random bytes]` | `true` + timestamps | Specific device IDs + IP whitelist |

---

## Generating HMAC Secrets

### Using OpenSSL (Recommended)

```bash
# Generate 32-byte random secret (256 bits)
openssl rand -hex 32

# Output: a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9t0u1v2w3x4y5z6

# Generate 64-byte random secret (512 bits) for high security
openssl rand -hex 64
```

### Using PicoClaw

```bash
# Generate secret with built-in tool
./build/picoclaw util generate-secret

# Output: Your HMAC secret: [32-byte hex string]
```

### Using Python

```python
import secrets

# Generate 32-byte secret
secret = secrets.token_hex(32)
print(secret)

# Output: a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9t0u1v2w3x4y5z6
```

### Using Go

```go
package main

import (
    "crypto/rand"
    "encoding/hex"
    "fmt"
)

func main() {
    bytes := make([]byte, 32)
    rand.Read(bytes)
    fmt.Println(hex.EncodeToString(bytes))
}
```

---

## Device Setup

### MaixCam Device Configuration

Configure your MaixCam device with the same secret:

```python
# MaixCam device code (MicroPython)
import hmac
import hashlib
import json
import socket

# Configuration
HMAC_SECRET = b"your-secret-key-change-in-production"
SERVER_HOST = "192.168.1.100"
SERVER_PORT = 8080

def create_message(device_id, message_type, data):
    """Create a message with HMAC signature"""
    message = {
        "device_id": device_id,
        "type": message_type,
        "data": data,
        "timestamp": time.time()
    }
    
    # Compute HMAC
    message_json = json.dumps(message, sort_keys=True)
    signature = hmac.new(
        HMAC_SECRET,
        message_json.encode(),
        hashlib.sha256
    ).hexdigest()
    
    # Add signature to message
    message["hmac"] = signature
    
    return message

def send_message(message):
    """Send message to PicoClaw server"""
    sock = socket.socket()
    sock.connect((SERVER_HOST, SERVER_PORT))
    
    # Send as JSON
    message_json = json.dumps(message)
    sock.send(message_json.encode())
    
    # Receive response
    response = sock.recv(1024).decode()
    sock.close()
    
    return json.loads(response)

# Example usage
message = create_message(
    device_id="maixcam_001",
    message_type="motion_detected",
    data={"confidence": 0.95, "bbox": [100, 200, 300, 400]}
)

response = send_message(message)
print(response)
```

### Arduino/ESP32 Example

```cpp
#include <WiFi.h>
#include <HTTPClient.h>
#include <mbedtls/md.h>

const char* hmac_secret = "your-secret-key-change-in-production";
const char* server_url = "http://192.168.1.100:8080";

String computeHMAC(const String& message) {
  const mbedtls_md_info_t* md_info = mbedtls_md_info_from_type(MBEDTLS_MD_SHA256);
  unsigned char hash[32];
  
  mbedtls_md_context_t ctx;
  mbedtls_md_init(&ctx);
  mbedtls_md_setup(&ctx, md_info, 1);
  mbedtls_md_hmac_starts(&ctx, (const unsigned char*)hmac_secret, strlen(hmac_secret));
  mbedtls_md_hmac_update(&ctx, (const unsigned char*)message.c_str(), message.length());
  mbedtls_md_hmac_finish(&ctx, hash);
  mbedtls_md_free(&ctx);
  
  // Convert to hex string
  char hex_string[65];
  for (int i = 0; i < 32; i++) {
    sprintf(&hex_string[i * 2], "%02x", hash[i]);
  }
  hex_string[64] = '\0';
  
  return String(hex_string);
}

void sendMessage(const String& deviceId, const String& messageType, const String& data) {
  // Create message JSON
  String message = "{\"device_id\":\"" + deviceId + 
                   "\",\"type\":\"" + messageType + 
                   "\",\"data\":" + data + "}";
  
  // Compute HMAC
  String hmac = computeHMAC(message);
  
  // Add HMAC to message
  String fullMessage = message.substring(0, message.length() - 1) + 
                       ",\"hmac\":\"" + hmac + "\"}";
  
  // Send to server
  HTTPClient http;
  http.begin(server_url);
  http.POST(fullMessage);
  http.end();
}
```

---

## Message Format

### Standard Message Structure

```json
{
  "device_id": "maixcam_001",
  "type": "motion_detected",
  "data": {
    "confidence": 0.95,
    "bbox": [100, 200, 300, 400],
    "timestamp": 1711234567.890
  },
  "timestamp": 1711234567.890,
  "hmac": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9t0u1v2w3x4y5z6a7b8c9d0e1f2"
}
```

### HMAC Computation

The HMAC is computed over the **canonical JSON representation**:

```python
import json
import hmac
import hashlib

def compute_hmac(message: dict, secret: bytes) -> str:
    """
    Compute HMAC-SHA256 for a message.
    
    The message is serialized to canonical JSON (sorted keys, no whitespace)
    before computing the HMAC.
    """
    # Remove existing HMAC if present
    message_copy = message.copy()
    message_copy.pop('hmac', None)
    
    # Canonical JSON (sorted keys, no whitespace)
    message_json = json.dumps(message_copy, sort_keys=True, separators=(',', ':'))
    
    # Compute HMAC
    signature = hmac.new(
        secret,
        message_json.encode('utf-8'),
        hashlib.sha256
    ).hexdigest()
    
    return signature

# Example
message = {
    "device_id": "maixcam_001",
    "type": "motion_detected",
    "data": {"confidence": 0.95}
}

secret = b"your-secret-key-change-in-production"
hmac_signature = compute_hmac(message, secret)
print(f"HMAC: {hmac_signature}")
```

### Verification on Server Side

```go
// pkg/channels/maixcam.go
func (mc *MaixCamChannel) verifyHMAC(message map[string]interface{}, providedHMAC string) error {
    // Get secret from config
    secret := []byte(mc.config.HMACSecret)
    
    // Remove HMAC from message for verification
    messageCopy := make(map[string]interface{})
    for k, v := range message {
        if k != "hmac" {
            messageCopy[k] = v
        }
    }
    
    // Canonical JSON
    messageJSON, err := json.Marshal(messageCopy)
    if err != nil {
        return fmt.Errorf("failed to marshal message: %w", err)
    }
    
    // Compute expected HMAC
    mac := hmac.New(sha256.New, secret)
    mac.Write(messageJSON)
    expectedHMAC := hex.EncodeToString(mac.Sum(nil))
    
    // Constant-time comparison (prevents timing attacks)
    if !hmac.Equal([]byte(providedHMAC), []byte(expectedHMAC)) {
        return fmt.Errorf("HMAC verification failed")
    }
    
    return nil
}
```

---

## Troubleshooting

### Common Issues

#### 1. "HMAC verification failed"

**Cause:** Mismatched secrets or message format.

**Solution:**
- Verify secret is identical on device and server
- Check JSON serialization (must be canonical: sorted keys, no whitespace)
- Ensure UTF-8 encoding on both sides

**Debug:**
```bash
# Check server logs
tail -f ~/.picoclaw/logs/picoclaw.log | grep "HMAC"

# Expected: "HMAC verified successfully"
# Error: "HMAC verification failed: signature mismatch"
```

#### 2. "HMAC secret not configured"

**Cause:** Missing `hmac_secret` in config.

**Solution:**
```json
{
  "channels": {
    "maixcam": {
      "hmac_secret": "your-secret-key"
    }
  }
}
```

#### 3. "Message rejected: HMAC required"

**Cause:** `require_hmac: true` but message has no HMAC.

**Solution:**
- Add HMAC to message on device
- Or temporarily disable: `require_hmac: false` (development only)

#### 4. "Timestamp expired"

**Cause:** Message timestamp is too old (replay protection).

**Solution:**
- Synchronize device clock with server
- Increase timestamp tolerance in config (if needed)

---

## Security Best Practices

### Secret Management

✅ **DO:**
- Use **32+ byte** random secrets (256+ bits)
- Store secrets in environment variables or secure vault
- Rotate secrets periodically (every 90 days)
- Use different secrets per device/environment

❌ **DON'T:**
- Use default secrets in production
- Commit secrets to version control
- Share secrets over insecure channels
- Reuse secrets across different systems

### Network Security

✅ **DO:**
- Use HTTPS/TLS for message transport
- Implement IP whitelisting for known devices
- Use VLANs to isolate IoT devices
- Monitor network traffic for anomalies

❌ **DON'T:**
- Send messages over plaintext HTTP on public networks
- Allow IoT devices on same network as sensitive systems
- Ignore unusual traffic patterns

### Monitoring and Alerts

Set up monitoring for HMAC failures:

```bash
#!/bin/bash
# Monitor HMAC failures
while true; do
  count=$(grep "HMAC verification failed" ~/.picoclaw/logs/picoclaw.log | \
          grep "$(date +%Y-%m-%d)" | wc -l)
  
  if [ $count -gt 10 ]; then
    echo "Alert: $count HMAC failures today" | \
      mail -s "MaixCam Security Alert" admin@example.com
  fi
  
  sleep 3600  # Check every hour
done
```

---

## Performance Impact

### Overhead

- **CPU**: ~1ms per message for HMAC computation (SHA-256)
- **Memory**: ~100 bytes per message for HMAC state
- **Latency**: <2ms total (compute + verify)

### Scalability

HMAC verification scales linearly:
- **100 messages/sec**: ~100ms CPU time
- **1,000 messages/sec**: ~1 second CPU time
- **10,000 messages/sec**: ~10 seconds CPU time (consider load balancing)

---

## Advanced Configuration

### Time-Based Replay Protection

Enable timestamp validation to prevent replay attacks:

```json
{
  "channels": {
    "maixcam": {
      "hmac_secret": "your-secret-key",
      "require_hmac": true,
      "timestamp_validation": {
        "enabled": true,
        "tolerance_seconds": 300
      }
    }
  }
}
```

### Per-Device Secrets

Use different secrets for each device:

```json
{
  "channels": {
    "maixcam": {
      "require_hmac": true,
      "device_secrets": {
        "maixcam_001": "secret-for-device-1",
        "maixcam_002": "secret-for-device-2",
        "maixcam_003": "secret-for-device-3"
      }
    }
  }
}
```

### Key Rotation

Implement automatic key rotation:

```go
type HMACConfig struct {
    CurrentSecret  string    // Active secret
    PreviousSecret string    // Previous secret (for grace period)
    RotatedAt      time.Time // Last rotation time
}

func (c *HMACConfig) VerifyHMAC(message, providedHMAC string) error {
    // Try current secret first
    if err := verifyWithSecret(message, providedHMAC, c.CurrentSecret); err == nil {
        return nil
    }
    
    // Try previous secret (grace period: 24 hours)
    if time.Since(c.RotatedAt) < 24*time.Hour {
        if err := verifyWithSecret(message, providedHMAC, c.PreviousSecret); err == nil {
            return nil
        }
    }
    
    return fmt.Errorf("HMAC verification failed")
}
```

---

## See Also

- **[SECURITY.md](SECURITY.md)** - Complete security documentation
- **[SENTINEL.md](SENTINEL.md)** - Skills Sentinel documentation
- **[RATE_LIMITING.md](RATE_LIMITING.md)** - Rate limiting guide
- **[DEVOPS_SECURITY.md](DEVOPS_SECURITY.md)** - DevOps and security scanning

---

**Last Updated:** March 24, 2026  
**Version:** v3.5.1+  
**Maintained By:** @comgunner

---

*PicoClaw: Ultra-Efficient AI in Go. $10 Hardware · <45MB RAM · <1s Startup*
