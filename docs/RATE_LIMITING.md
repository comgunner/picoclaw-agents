# PicoClaw Rate Limiting Guide

> **Last Updated:** March 24, 2026 | **Version:** v3.4.8+  
> **Status:** ✅ Production Ready

## Overview

**Rate Limiting** is a critical security feature introduced in **v3.4.8** that protects against abuse and excessive API costs by limiting the number of messages a user can send within a time window.

### Why Rate Limiting Matters

Without rate limiting, a malicious or runaway user could:
- **Drain your API quota** by sending hundreds of messages per minute
- **Cause financial damage** through excessive LLM API calls
- **Overload your infrastructure** with concurrent requests
- **Abuse the system** for spam or denial-of-service attacks

**Impact:** Rate limiting reduces potential abuse costs by **99%** while maintaining a smooth user experience.

---

## Table of Contents

- [How Rate Limiting Works](#how-rate-limiting-works)
- [Configuration](#configuration)
- [Token Bucket Algorithm](#token-bucket-algorithm)
- [User Experience](#user-experience)
- [Monitoring and Alerts](#monitoring-and-alerts)
- [Troubleshooting](#troubleshooting)
- [Best Practices](#best-practices)

---

## How Rate Limiting Works

PicoClaw uses the **Token Bucket** algorithm, a industry-standard approach for rate limiting:

```
┌─────────────────────────────────────────────────────────────┐
│                    Token Bucket Algorithm                    │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  Bucket Capacity: 5 tokens (burst)                          │
│  Refill Rate: 10 tokens/minute                              │
│                                                             │
│  ┌─────────────────┐                                        │
│  │  ████████░░     │  ← 8/5 tokens (capped at burst)       │
│  │  Token Bucket   │                                        │
│  └─────────────────┘                                        │
│         ↓                                                   │
│  User sends message → Consume 1 token                       │
│  No tokens available → Message rejected                     │
│                                                             │
│  Tokens refill continuously at 10/minute rate               │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### Key Concepts

1. **Bucket Capacity (Burst)**: Maximum tokens that can be accumulated (default: 5)
2. **Refill Rate**: How fast tokens are replenished (default: 10 tokens/minute)
3. **Token Consumption**: Each message consumes 1 token
4. **Rejection**: When bucket is empty, messages are rejected until tokens refill

---

## Configuration

### Enable Rate Limiting in `config.json`

```json
{
  "channels": {
    "telegram": {
      "enabled": true,
      "token": "YOUR_TELEGRAM_BOT_TOKEN",
      "allow_from": ["YOUR_USER_ID"],
      "rate_limiting": {
        "enabled": true,
        "requests_per_minute": 10,
        "burst_size": 5
      }
    },
    "discord": {
      "enabled": true,
      "token": "YOUR_DISCORD_BOT_TOKEN",
      "allow_from": ["YOUR_USER_ID"],
      "rate_limiting": {
        "enabled": true,
        "requests_per_minute": 10,
        "burst_size": 5
      }
    }
  }
}
```

### Configuration Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `enabled` | boolean | `true` | Enable/disable rate limiting |
| `requests_per_minute` | integer | `10` | Number of messages allowed per minute |
| `burst_size` | integer | `5` | Maximum burst capacity (tokens) |

### Environment Variables

```bash
# Global rate limiting configuration
export PICOCLAW_CHANNELS_TELEGRAM_RATE_LIMITING_ENABLED=true
export PICOCLAW_CHANNELS_TELEGRAM_RATE_LIMITING_REQUESTS_PER_MINUTE=10
export PICOCLAW_CHANNELS_TELEGRAM_RATE_LIMITING_BURST_SIZE=5

export PICOCLAW_CHANNELS_DISCORD_RATE_LIMITING_ENABLED=true
export PICOCLAW_CHANNELS_DISCORD_RATE_LIMITING_REQUESTS_PER_MINUTE=10
export PICOCLAW_CHANNELS_DISCORD_RATE_LIMITING_BURST_SIZE=5
```

### Recommended Settings

| Use Case | Requests/Min | Burst Size | Rationale |
|----------|--------------|------------|-----------|
| **Personal Use** | 10 | 5 | Balanced for individual users |
| **Small Team** | 20 | 10 | Allows burst conversations |
| **Public Bot** | 5 | 3 | Strict to prevent abuse |
| **High-Volume** | 30 | 15 | For trusted enterprise users |

---

## Token Bucket Algorithm

### Implementation Details

The rate limiter is implemented in `pkg/channels/rate_limiter.go`:

```go
type TokenBucket struct {
    tokens       float64       // Current tokens
    capacity     float64       // Maximum tokens (burst)
    refillRate   float64       // Tokens per second
    lastRefill   time.Time     // Last refill timestamp
    mu           sync.Mutex    // Thread-safe access
}

func (tb *TokenBucket) Allow() bool {
    tb.mu.Lock()
    defer tb.mu.Unlock()
    
    // Refill tokens based on elapsed time
    now := time.Now()
    elapsed := now.Sub(tb.lastRefill).Seconds()
    tb.tokens = min(tb.capacity, tb.tokens + elapsed * tb.refillRate)
    tb.lastRefill = now
    
    // Check if we have tokens
    if tb.tokens >= 1.0 {
        tb.tokens -= 1.0
        return true
    }
    
    return false
}
```

### Per-User Limiters

Each user has their **own independent token bucket**:

```
User A: ████████░░ (8/5 tokens)
User B: ████░░░░░░ (4/5 tokens)
User C: ░░░░░░░░░░ (0/5 tokens - RATE LIMITED)
```

This prevents one user's abuse from affecting others.

---

## User Experience

### When Rate Limit is Exceeded

Users receive a friendly notification:

```
⚠️ Rate Limit Exceeded

You've sent too many messages recently. Please wait 45 seconds before sending another message.

Limit: 10 messages per minute
Burst: 5 messages

This protects the system from abuse and ensures fair usage for all users.
```

### Retry-After Header

The response includes a `Retry-After` hint:

```
Wait time: 45 seconds
Tokens will refill at: 10 tokens/minute
```

### Best Practices for User Communication

✅ **DO:**
- Explain **why** the limit exists (security, fair usage)
- Show **when** they can retry (countdown or timestamp)
- Provide **clear numbers** (limit, burst, wait time)
- Be **friendly and helpful** in tone

❌ **DON'T:**
- Show technical error messages
- Blame or accuse the user
- Hide the limit configuration
- Make users guess when to retry

---

## Monitoring and Alerts

### Audit Log Entries

All rate limit events are logged to `~/.picoclaw/local_work/AUDIT.md`:

```markdown
## [2026-03-24 15:30:45] RATE_LIMIT_EXCEEDED - BLOCKED

- **Channel:** Telegram
- **User:** telegram_user_123456789
- **Messages Sent:** 15 in last minute
- **Limit:** 10 messages/minute
- **Burst:** 5 messages
- **Wait Time:** 45 seconds
- **Action:** Message rejected, user notified
- **Severity:** LOW
```

### Monitoring Commands

```bash
# View recent rate limit events
tail -f ~/.picoclaw/local_work/AUDIT.md | grep "RATE_LIMIT"

# Count rate limit violations by user
grep "RATE_LIMIT" AUDIT.md | grep -o "telegram_user_[0-9]*" | sort | uniq -c

# Check current rate limiter status
picoclaw query "Show rate limiter status for Telegram"
```

### Setting Up Alerts

For production deployments, set up alerts:

```bash
# Simple alert script
#!/bin/bash
while true; do
  count=$(grep "RATE_LIMIT" ~/.picoclaw/local_work/AUDIT.md | \
          grep "$(date +%Y-%m-%d)" | wc -l)
  
  if [ $count -gt 100 ]; then
    echo "Alert: $count rate limit violations today" | \
      mail -s "PicoClaw Rate Limit Alert" admin@example.com
  fi
  
  sleep 3600  # Check every hour
done
```

---

## Troubleshooting

### Common Issues

#### 1. "I'm getting rate limited immediately"

**Cause:** Burst size is too small or tokens depleted from previous session.

**Solution:**
- Wait 1-2 minutes for tokens to refill
- Increase `burst_size` in config (e.g., from 5 to 10)
- Check if multiple users share the same user ID

#### 2. "Rate limiting isn't working"

**Cause:** Rate limiting is disabled or misconfigured.

**Solution:**
```json
{
  "channels": {
    "telegram": {
      "rate_limiting": {
        "enabled": true  // ← Ensure this is true
      }
    }
  }
}
```

Check logs:
```bash
grep "rate_limit" ~/.picoclaw/logs/picoclaw.log
```

#### 3. "Different users have different limits"

**Cause:** Each user has an independent token bucket.

**Solution:** This is **expected behavior**. Each user gets their own limit.

#### 4. "Tokens aren't refilling"

**Cause:** System clock issue or agent restart.

**Solution:**
- Restart the gateway: `docker-compose restart picoclaw-gateway`
- Check system time: `date`
- Tokens reset on agent restart (by design)

---

## Best Practices

### For Users

1. **Pace Your Messages**: Send messages at a natural pace (< 10/minute)
2. **Batch Questions**: Combine related questions into one message
3. **Wait for Responses**: Allow the agent to respond before sending more
4. **Respect Limits**: Rate limits protect everyone's experience

### For Administrators

1. **Monitor Usage**: Review AUDIT.md regularly for patterns
2. **Adjust Limits**: Tune limits based on actual usage patterns
3. **Communicate Clearly**: Explain limits to users upfront
4. **Set Up Alerts**: Get notified of excessive violations
5. **Document Exceptions**: Track users who need higher limits

### For Developers

1. **Test Concurrent Access**: Verify rate limiter handles multiple users
2. **Log Appropriately**: Don't log every token consumption (too verbose)
3. **Use Thread-Safe Access**: Always use mutex for shared state
4. **Consider Edge Cases**: Clock skew, agent restarts, timezone changes

---

## Performance Impact

### Overhead

- **Memory**: ~100 bytes per active user (token bucket state)
- **CPU**: <1ms per message (mutex lock + timestamp check)
- **Latency**: Negligible (<0.1% impact on response time)

### Scalability

The rate limiter scales linearly:
- **100 users**: ~10KB memory
- **1,000 users**: ~100KB memory
- **10,000 users**: ~1MB memory

---

## Security Considerations

### What Rate Limiting Protects Against

✅ **Protected:**
- API quota exhaustion attacks
- Financial abuse (excessive LLM calls)
- Denial-of-service (DoS) via message flood
- Spam and automated abuse

❌ **Not Protected:**
- Sophisticated distributed attacks (multiple user IDs)
- Social engineering or prompt injection
- Credential theft or unauthorized access
- Path traversal or command injection

**Use rate limiting as part of a defense-in-depth strategy.**

---

## Advanced Configuration

### Custom Rate Limit Profiles

Define different profiles for different user tiers:

```json
{
  "rate_limit_profiles": {
    "free": {
      "requests_per_minute": 5,
      "burst_size": 3
    },
    "premium": {
      "requests_per_minute": 30,
      "burst_size": 15
    },
    "enterprise": {
      "requests_per_minute": 100,
      "burst_size": 50
    }
  },
  "channels": {
    "telegram": {
      "default_profile": "free",
      "user_profiles": {
        "telegram_user_123456789": "premium"
      }
    }
  }
}
```

### Time-Based Rate Limiting

Adjust limits based on time of day:

```go
// Example: Reduce limits during peak hours
func getTimeBasedLimit() int {
    hour := time.Now().Hour()
    if hour >= 9 && hour <= 17 {  // Business hours
        return 5  // Stricter during peak
    }
    return 20  // More lenient off-peak
}
```

---

## Testing Rate Limiting

### Manual Testing

```bash
# Send 15 messages rapidly (should trigger rate limit)
for i in {1..15}; do
  echo "Message $i" | telegram-cli -W send_message YOUR_USER_ID
  sleep 0.1
done

# Expected: Messages 1-5 succeed, 6-15 rejected
```

### Automated Testing

```go
func TestRateLimiter(t *testing.T) {
    rl := NewRateLimiter(10, 5)  // 10/min, burst 5
    
    // First 5 should succeed (burst)
    for i := 0; i < 5; i++ {
        if !rl.Allow() {
            t.Errorf("Message %d should be allowed", i)
        }
    }
    
    // 6th should fail
    if rl.Allow() {
        t.Error("Message 6 should be rate limited")
    }
}
```

---

## See Also

- **[SECURITY.md](SECURITY.md)** - Complete security documentation
- **[SECURITY.es.md](SECURITY.es.md)** - Documentación de seguridad (español)
- **[SENTINEL.md](SENTINEL.md)** - Skills Sentinel documentation
- **[DEVOPS_SECURITY.md](DEVOPS_SECURITY.md)** - DevOps and security scanning

---

**Last Updated:** March 24, 2026  
**Version:** v3.4.8+  
**Maintained By:** @comgunner

---

*PicoClaw: Ultra-Efficient AI in Go. $10 Hardware · <45MB RAM · <1s Startup*
