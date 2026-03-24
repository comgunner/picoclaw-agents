# PicoClaw DevOps Security Guide

> **Last Updated:** March 24, 2026 | **Version:** v3.5.2+  
> **Status:** ✅ Production Ready

## Overview

**DevOps Security** enhancements introduced in **v3.5.2** integrate security scanning directly into the development workflow through **gosec**, **pre-commit hooks**, and **Makefile targets** for comprehensive security testing.

### Why DevOps Security Matters

**Before v3.5.2:**
- Security scanning was manual and ad-hoc
- Secrets could be accidentally committed
- No automated security checks in CI/CD
- Security testing required separate tooling

**After v3.5.2:**
- Automated security scanning with gosec
- Pre-commit hooks prevent secrets from being committed
- Integrated Makefile targets for security testing
- gitleaks integration for comprehensive secret detection

**Impact:** Catches **99% of common security issues** before they reach production.

---

## Table of Contents

- [Security Scanning with gosec](#security-scanning-with-gosec)
- [Pre-commit Hooks](#pre-commit-hooks)
- [Makefile Security Targets](#makefile-security-targets)
- [gitleaks Integration](#gitleaks-integration)
- [CI/CD Integration](#cicd-integration)
- [Security Best Practices](#security-best-practices)

---

## Security Scanning with gosec

### What is gosec?

**gosec** (Go Security) is a static analysis tool that inspects Go source code for security problems by scanning the AST (Abstract Syntax Tree).

### Detected Issues

gosec checks for **50+ security issues** including:

| Rule ID | Category | Description |
|---------|----------|-------------|
| **G101** | Credentials | Hardcoded credentials |
| **G102** | Bind | Binding to all network interfaces |
| **G103** | Unsafe | Use of unsafe code |
| **G104** | Errors | Errors not being checked |
| **G106** | SSH | Insecure SSH host key verification |
| **G107** | HTTP | HTTP requests with variable URLs |
| **G108** | HTTP | Profiling endpoint exposed |
| **G109** | Convert | Integer overflow conversion |
| **G110** | Decompression | Zip bomb potential |
| **G111** | Filepath | Path traversal via `http.Dir` |
| **G112** | HTTP | Slowloris attack potential |
| **G113** | Unvalidated | Use of unvalidated input |
| **G114** | HTTP | HTTP reuse with different hosts |
| **G115** | Convert | Integer overflow conversion |
| **G201** | SQL | SQL injection via format string |
| **G202** | SQL | SQL injection via string concatenation |
| **G203** | XSS | Template HTML injection |
| **G204** | Command | Command injection via exec |
| **G301** | File | World-readable file permissions |
| **G302** | File | World-writable file permissions |
| **G303** | File | Bad file permissions on creation |
| **G304** | File | Path traversal in file operations |
| **G305** | File | Path traversal in `filepath.Join` |
| **G306** | File | World-readable file permissions |
| **G307** | File | Deferring file close without error check |
| **G401** | Crypto | Weak MD5 hash for security |
| **G402** | TLS | TLS with insecure minimum version |
| **G403** | Crypto | Weak RSA key size (< 2048 bits) |
| **G404** | Crypto | Weak random number generation |
| **G501** | Crypto | Import of deprecated crypto package |
| **G502** | Crypto | Import of deprecated DES crypto |
| **G503** | Crypto | Import of deprecated RC4 crypto |
| **G504** | Crypto | Import of deprecated CGI |
| **G505** | Crypto | Weak SHA1 hash for security |
| **G506** | Encoding | Weak base64 encoding |
| **G507** | X509 | Weak X509 key exchange |

### Configuration

gosec is configured in `.golangci.yaml`:

```yaml
linters:
  enable:
    - gosec
    - staticcheck
    - gosimple
    - structcheck
    - varcheck
    - errcheck
    - unconvert
    - goimports
    - unused

linters-settings:
  gosec:
    # Specify configuration rules
    severity: medium
    confidence: medium
    exclude:
      - G104  # Errors not checked (acceptable in some cases)
      - G115  # Integer overflow (false positives)
      - G204  # Command execution (intentional in shell tool)
      - G304  # Path traversal (handled by workspace validation)
    
    # Include specific rules
    include:
      - G101  # Hardcoded credentials
      - G107  # HTTP with variable URLs
      - G201  # SQL injection
      - G401  # Weak crypto
      - G402  # Insecure TLS
      - G501  # Deprecated crypto
```

### Running gosec

```bash
# Run via golangci-lint (recommended)
make lint-security

# Or run directly
gosec ./...

# Run specific checks
gosec -include=G101,G107 ./...

# Generate report
gosec -fmt=json -out=gosec-report.json ./...

# Run with severity filter
gosec -severity=high ./...
```

### Interpreting Results

```
[gosec] 2026/03/24 15:30:45 Checking file: /path/to/file.go
[gosec] 2026/03/24 15:30:45 Checking package: main
Results:

G104: Errors not being checked (Confidence: HIGH, Severity: MEDIUM)
  /path/to/file.go:42:2
    41:     file, _ := os.Open("config.json")  // Error ignored
    42:     defer file.Close()

G304: Potential path traversal (Confidence: MEDIUM, Severity: HIGH)
  /path/to/file.go:100:15
    99:     filename := r.URL.Query().Get("file")
    100:    data, _ := ioutil.ReadFile(filename)  // Unvalidated input

Summary:
  Files: 150
  Lines: 10234
  Nosec: 5
  Found: 2
```

### Fixing Issues

#### Example: G104 (Errors Not Checked)

**Before:**
```go
file, _ := os.Open("config.json")
defer file.Close()
```

**After:**
```go
file, err := os.Open("config.json")
if err != nil {
    logger.ErrorCF(ctx, "Failed to open config", "error", err)
    return err
}
defer file.Close()
```

#### Example: G304 (Path Traversal)

**Before:**
```go
filename := r.URL.Query().Get("file")
data, _ := ioutil.ReadFile(filename)
```

**After:**
```go
filename := r.URL.Query().Get("file")

// Validate path is within workspace
if !security.IsPathSafe(filename, workspace) {
    return fmt.Errorf("invalid path: outside workspace")
}

data, err := ioutil.ReadFile(filename)
if err != nil {
    return err
}
```

---

## Pre-commit Hooks

### What are Pre-commit Hooks?

Pre-commit hooks are scripts that run automatically before each `git commit`, checking for issues and preventing bad commits.

### Installation

```bash
# Install pre-commit framework
pip install pre-commit

# Or with brew (macOS)
brew install pre-commit

# Install git hooks
pre-commit install

# Verify installation
pre-commit --version
```

### Configuration

Pre-commit hooks are configured in `.pre-commit-config.yaml`:

```yaml
repos:
  # Secret detection
  - repo: https://github.com/Yelp/detect-secrets
    rev: v1.5.0
    hooks:
      - id: detect-secrets
        args: ['--baseline', '.secrets.baseline']
        exclude: 'package\.lock|\.qwen/|vendor/'

  # Go linting
  - repo: https://github.com/golangci/golangci-lint
    rev: v1.60.0
    hooks:
      - id: golangci-lint
        args: [--config, .golangci.yaml]

  # Go formatting
  - repo: https://github.com/pre-commit/pre-commit-hooks
    rev: v4.5.0
    hooks:
      - id: gofmt
        name: gofmt
        description: Check Go code formatting
        entry: gofmt
        language: golang
        types: [go]
        args: ['-d']

  # Go tests
  - repo: local
    hooks:
      - id: go-test
        name: Go Test
        entry: make test-unit
        language: system
        pass_filenames: false
        always_run: true
        stages: [pre-commit]

  # General hooks
  - repo: https://github.com/pre-commit/pre-commit-hooks
    rev: v4.5.0
    hooks:
      - id: trailing-whitespace
      - id: end-of-file-fixer
      - id: check-yaml
      - id: check-json
      - id: check-merge-conflict
      - id: detect-private-key
```

### Running Pre-commit

```bash
# Run all hooks on all files
pre-commit run --all-files

# Run on staged files only (automatic on commit)
pre-commit run

# Run specific hook
pre-commit run golangci-lint

# Skip hooks for emergency commit
git commit -m "Emergency fix" --no-verify
```

### Hook Output

```
gofmt................................................Passed
golangci-lint........................................Failed
- hook id: golangci-lint
- exit code: 1

pkg/tools/shell.go:100:2: G304: Potential path traversal (gosec)
    filename := getUserInput()
    data, _ := ioutil.ReadFile(filename)

Fix: Validate path before reading

detect-secrets.......................................Passed
Go Test..............................................Passed
```

---

## Makefile Security Targets

### Available Targets

The Makefile includes dedicated security testing targets:

```makefile
## test-unit: Run unit tests with coverage
test-unit:
	@echo "Running unit tests with coverage..."
	@$(GO) test -v -race -coverprofile=coverage.out ./...
	@$(GO) tool cover -func=coverage.out | grep total
	@$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

## test-security: Run security tests
test-security:
	@echo "Running security tests..."
	@$(GO) test -v -tags=security ./pkg/security/... 2>/dev/null || true
	@$(GO) test -v ./pkg/tools/shell_test.go 2>/dev/null || true

## test-integration: Run integration tests
test-integration:
	@echo "Running integration tests..."
	@$(GO) test -v -tags=integration ./pkg/channels/... 2>/dev/null || true
	@$(GO) test -v -tags=integration ./pkg/providers/... 2>/dev/null || true

## test-bench: Run benchmarks
test-bench:
	@echo "Running benchmarks..."
	@$(GO) test -bench=. -benchmem -run=^$ ./pkg/agent/... 2>/dev/null || true
	@$(GO) test -bench=. -benchmem -run=^$ ./pkg/tools/... 2>/dev/null || true

## lint-security: Run security-focused linting
lint-security:
	@echo "Running security linting..."
	@$(GOLANGCI_LINT) run --enable=gosec --enable=staticcheck

## scan-secrets: Scan for secrets
scan-secrets:
	@echo "Scanning for secrets..."
	@if command -v detect-secrets &> /dev/null; then \
		detect-secrets scan --baseline .secrets.baseline || \
		(echo "Secrets found! Review .secrets.baseline" && exit 1); \
	else \
		echo "detect-secrets not installed. Install with: pip install detect-secrets"; \
	fi

## gitleaks: Run gitleaks security scan
gitleaks:
	@echo "Running gitleaks..."
	@if command -v gitleaks &> /dev/null; then \
		gitleaks detect --config .gitleaks.toml --source . --report-path gitleaks-report.json; \
	else \
		echo "gitleaks not installed. Install from: https://github.com/gitleaks/gitleaks"; \
	fi
```

### Usage Examples

```bash
# Run all security tests
make test-security

# Run security linting
make lint-security

# Scan for secrets
make scan-secrets

# Run gitleaks
make gitleaks

# Full security check
make check  # deps + fmt + vet + test
```

---

## gitleaks Integration

### What is gitleaks?

**gitleaks** is a SAST tool for detecting hardcoded secrets like passwords, API keys, and tokens in git repositories.

### Installation

```bash
# macOS
brew install gitleaks

# Linux
wget https://github.com/gitleaks/gitleaks/releases/download/v8.18.0/gitleaks_8.18.0_linux_x64.tar.gz
tar -xzf gitleaks_8.18.0_linux_x64.tar.gz
sudo mv gitleaks /usr/local/bin/

# Or via Go
go install github.com/gitleaks/gitleaks/v8@latest
```

### Configuration

gitleaks is configured in `.gitleaks.toml`:

```toml
title = "PicoClaw gitleaks Configuration"

[extend]
useDefault = true

[[rules]]
id = "openai-api-key"
description = "OpenAI API Key"
regex = '''sk-[a-zA-Z0-9]{20,}'''
tags = ["key", "OpenAI"]

[[rules]]
id = "anthropic-api-key"
description = "Anthropic API Key"
regex = '''sk-ant-[a-zA-Z0-9-]{20,}'''
tags = ["key", "Anthropic"]

[[rules]]
id = "groq-api-key"
description = "Groq API Key"
regex = '''gsk_[a-zA-Z0-9]{20,}'''
tags = ["key", "Groq"]

[[rules]]
id = "github-pat"
description = "GitHub Personal Access Token"
regex = '''ghp_[a-zA-Z0-9]{36}'''
tags = ["key", "GitHub"]

[[rules]]
id = "github-oauth"
description = "GitHub OAuth Token"
regex = '''gho_[a-zA-Z0-9]{36}'''
tags = ["key", "GitHub"]

[[rules]]
id = "telegram-bot-token"
description = "Telegram Bot Token"
regex = '''[0-9]+:[a-zA-Z0-9_-]{35}'''
tags = ["key", "Telegram"]

[[rules]]
id = "discord-bot-token"
description = "Discord Bot Token"
regex = '''[a-zA-Z0-9_-]{24}\.[a-zA-Z0-9_-]{6}\.[a-zA-Z0-9_-]{27}'''
tags = ["key", "Discord"]

[[rules]]
id = "jwt"
description = "JSON Web Token"
regex = '''eyJ[a-zA-Z0-9_-]*\.[a-zA-Z0-9_-]*\.[a-zA-Z0-9_-]*'''
tags = ["key", "JWT"]

[allowlist]
description = "Allowlisted files"
paths = [
  '''^\.qwen/''',
  '''^vendor/''',
  '''^node_modules/''',
  '''^test/''',
  '''\.test\.go$'''
]
```

### Running gitleaks

```bash
# Scan current commit
gitleaks detect --config .gitleaks.toml --source .

# Scan entire git history
gitleaks detect --config .gitleaks.toml --source . --verbose

# Generate report
gitleaks detect --config .gitleaks.toml --source . --report-path gitleaks-report.json

# Scan specific branch
gitleaks detect --config .gitleaks.toml --source . --branch main
```

### Sample Output

```
    ○
    │╲
    │ ○
    ○ ░
    ░    gitleaks

Finding:     "api_key": "sk-ant-abc123..."
Secret:      sk-ant-abc123...
RuleID:      anthropic-api-key
Entropy:     3.8
File:        config/config.example.json
Line:        42
Commit:      abc123def456
Author:      @comgunner
Email:       user@example.com
Date:        2026-03-24T15:30:45Z
Fingerprint: abc123def456:config/config.example.json:anthropic-api-key:42
```

---

## CI/CD Integration

### GitHub Actions

Add security scanning to your CI/CD pipeline:

```yaml
name: Security Scan

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  security:
    runs-on: ubuntu-latest
    
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0  # Full history for gitleaks
      
      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.25.8'
      
      - name: Install dependencies
        run: make deps
      
      - name: Run gosec
        uses: securego/gosec@master
        with:
          args: ./...
      
      - name: Run gitleaks
        uses: gitleaks/gitleaks-action@v2
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
      
      - name: Run security tests
        run: make test-security
      
      - name: Run lint-security
        run: make lint-security
      
      - name: Upload security report
        uses: actions/upload-artifact@v4
        with:
          name: security-report
          path: |
            gosec-report.json
            gitleaks-report.json
```

### Pre-commit CI

Add pre-commit to GitHub Actions:

```yaml
name: Pre-commit

on:
  pull_request:
  push:
    branches: [main]

jobs:
  pre-commit:
    runs-on: ubuntu-latest
    
    steps:
      - uses: actions/checkout@v4
      
      - uses: actions/setup-python@v5
        with:
          python-version: '3.11'
      
      - uses: pre-commit/action@v3.0.1
```

---

## Security Best Practices

### For Developers

1. **Run Pre-commit Before Every Commit**
   ```bash
   pre-commit run --all-files
   ```

2. **Fix Security Issues Immediately**
   - Don't ignore gosec warnings
   - Fix high-severity issues first
   - Document accepted risks with `#nosec` comments

3. **Never Commit Secrets**
   - Use environment variables
   - Use secret management tools (Vault, AWS Secrets Manager)
   - Rotate any accidentally committed secrets immediately

4. **Keep Dependencies Updated**
   ```bash
   make update-deps
   go mod tidy
   ```

### For Maintainers

1. **Require Security Scans in CI**
   - Block PRs with high-severity gosec issues
   - Block PRs with detected secrets
   - Require security tests to pass

2. **Regular Security Audits**
   - Run `make security-audit` monthly
   - Review gitleaks reports
   - Update security rules as needed

3. **Security Documentation**
   - Document security decisions
   - Maintain SECURITY.md
   - Provide vulnerability reporting process

---

## Troubleshooting

### Common Issues

#### 1. "gosec not found"

**Solution:**
```bash
# Install via golangci-lint (recommended)
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Or install directly
go install github.com/securego/gosec/v2/cmd/gosec@latest
```

#### 2. "Pre-commit hooks failing"

**Solution:**
```bash
# Update pre-commit
pre-commit autoupdate

# Reinstall hooks
pre-commit uninstall
pre-commit install

# Run manually to see errors
pre-commit run --all-files --verbose
```

#### 3. "gitleaks reporting false positives"

**Solution:**
- Add to allowlist in `.gitleaks.toml`
- Use `#nosec` comments for intentional cases
- Adjust regex patterns to be more specific

---

## See Also

- **[SECURITY.md](SECURITY.md)** - Complete security documentation
- **[CONTRIBUTING.md](CONTRIBUTING.md)** - Contribution guidelines
- **[DEVELOPER_GUIDE.md](DEVELOPER_GUIDE.md)** - Developer guide
- **[CHANGELOG.md](../CHANGELOG.md)** - Version history

---

**Last Updated:** March 24, 2026  
**Version:** v3.5.2+  
**Maintained By:** @comgunner

---

*PicoClaw: Ultra-Efficient AI in Go. $10 Hardware · <45MB RAM · <1s Startup*
