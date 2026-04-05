# Contributing to PicoClaw

> **Last Updated:** March 2026 | **Version:** v3.4.5+

Thank you for your interest in contributing to PicoClaw! This project is a community-driven effort to build a lightweight and versatile personal AI assistant. We welcome contributions of all kinds: bug fixes, features, documentation, translations, and testing.

PicoClaw itself was substantially developed with AI assistance — we embrace this approach and have built our contribution process around it.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Ways to Contribute](#ways-to-contribute)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Making Changes](#making-changes)
- [AI-Assisted Contributions](#ai-assisted-contributions)
- [Pull Request Process](#pull-request-process)
- [Branch Strategy](#branch-strategy)
- [Code Review](#code-review)
- [Communication](#communication)

---

## Code of Conduct

We are committed to maintaining a welcoming and respectful community. Be kind, constructive, and assume good faith. Harassment or discrimination of any kind will not be tolerated.

**Our Standards:**
- Use welcoming and inclusive language
- Be respectful of differing viewpoints and experiences
- Gracefully accept constructive criticism
- Focus on what is best for the community
- Show empathy towards other community members

**Unacceptable Behavior:**
- The use of sexualized language or imagery
- Trolling, insulting/derogatory comments, and personal or political attacks
- Public or private harassment
- Publishing others' private information without explicit permission
- Other conduct which could reasonably be considered inappropriate in a professional setting

**Reporting:**
If you experience or witness unacceptable behavior, please report it by opening a private issue or contacting a maintainer directly.

---

## Ways to Contribute

### 1. Bug Reports

**When to Report:**
- Unexpected crashes or errors
- Features not working as documented
- Performance issues
- Security vulnerabilities (see [SECURITY.md](./SECURITY.md))

**How to Report:**
1. Search existing issues to avoid duplicates
2. Use the bug report template
3. Include:
   - PicoClaw version
   - OS and hardware
   - Steps to reproduce
   - Expected vs actual behavior
   - Relevant logs (redact secrets)

### 2. Feature Requests

**Before Requesting:**
- Search existing feature requests
- Check if the feature aligns with project goals (lightweight, multi-agent, portable)

**How to Request:**
1. Use the feature request template
2. Describe the use case
3. Explain why it's needed
4. Suggest possible implementations (optional)

### 3. Code Contributions

**Types of Code Contributions:**
- Bug fixes
- New features
- Performance improvements
- Security enhancements
- Test additions
- Refactoring

**Before Coding:**
- Open an issue to discuss the change
- Check if similar work is in progress
- Ensure you have time to complete the PR

### 4. Documentation

**Documentation Needs:**
- README improvements
- API documentation
- Tutorial creation
- Translation to new languages
- Code comment enhancements
- Troubleshooting guides

### 5. Testing

**Testing Opportunities:**
- Test on new hardware platforms
- Test with different LLM providers
- Test new chat channels
- Report compatibility issues
- Write test cases

---

## Getting Started

### 1. Fork the Repository

```bash
# On GitHub, click "Fork" button
# Then clone your fork
git clone https://github.com/<your-username>/picoclaw.git
cd picoclaw
```

### 2. Add Upstream Remote

```bash
git remote add upstream https://github.com/comgunner/picoclaw-agents.git
git remote -v
# Should show both origin (your fork) and upstream
```

### 3. Create a Branch

```bash
git checkout main
git pull upstream main
git checkout -b feature/your-feature-name
```

**Branch Naming:**
- `feature/xyz` — New features
- `fix/xyz` — Bug fixes
- `docs/xyz` — Documentation
- `test/xyz` — Test additions
- `refactor/xyz` — Code refactoring

---

## Development Setup

### Prerequisites

| Tool | Version | Installation |
|------|---------|--------------|
| **Go** | 1.25.8+ | [go.dev](https://go.dev/dl/) |
| **make** | 4.0+ | `apt install make` / `brew install make` |
| **git** | 2.30+ | `apt install git` / `brew install git` |
| **docker** | 24.0+ (optional) | [docker.com](https://docs.docker.com/get-docker/) |

### Build

```bash
# Download dependencies
make deps

# Build binary (runs go generate first)
make build

# Full pre-commit check
make check  # deps + fmt + vet + test
```

### Running Tests

```bash
# Run all tests
make test

# Run specific test
go test -run TestName -v ./pkg/session/

# Run benchmarks
go test -bench=. -benchmem -run='^$' ./...

# Run with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Code Style

```bash
# Format code
make fmt

# Static analysis
make vet

# Full linter run
make lint
```

**All CI checks must pass before a PR can be merged.** Run `make check` locally before pushing to catch issues early.

---

## Making Changes

### Branching

Always branch off `main` and target `main` in your PR. Never push directly to `main` or any `release/*` branch:

```bash
git checkout main
git pull upstream main
git checkout -b your-feature-branch
```

Use descriptive branch names:
- ✅ `fix/telegram-timeout`
- ✅ `feat/ollama-provider`
- ✅ `docs/contributing-guide`
- ❌ `patch-1`
- ❌ `new-feature`
- ❌ `fix`

### Commits

**Commit Message Guidelines:**
- Write clear, concise messages in English
- Use imperative mood: "Add retry logic" not "Added retry logic"
- Reference related issues: `Fix session leak (#123)`
- Keep commits focused: one logical change per commit

**Conventional Commits Format:**
```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation
- `style`: Formatting (no code change)
- `refactor`: Code refactoring
- `test`: Adding tests
- `chore`: Maintenance

**Examples:**
```bash
feat(agents): Add autonomous runtime for background processing

fix(telegram): Fix message timeout in long conversations

docs(security): Add security best practices guide

test(providers): Add integration tests for Antigravity provider
```

### Keeping Up to Date

Rebase your branch onto upstream `main` before opening a PR:

```bash
git fetch upstream
git rebase upstream/main

# Resolve conflicts if any
git add <files>
git rebase --continue

# Force push to your fork
git push -f origin your-feature-branch
```

---

## AI-Assisted Contributions

PicoClaw was built with substantial AI assistance, and we fully embrace AI-assisted development. However, contributors must understand their responsibilities when using AI tools.

### Disclosure Is Required

Every PR must disclose AI involvement using the PR template's **🤖 AI Code Generation** section. There are three levels:

| Level | Description |
|-------|-------------|
| 🤖 **Fully AI-generated** | AI wrote the code; contributor reviewed and validated it |
| 🛠️ **Mostly AI-generated** | AI produced the draft; contributor made significant modifications |
| 👨‍💻 **Mostly Human-written** | Contributor led; AI provided suggestions or none at all |

**Honest disclosure is expected.** There is no stigma attached to any level — what matters is the quality of the contribution.

### You Are Responsible for What You Submit

Using AI to generate code does not reduce your responsibility as the contributor. Before opening a PR with AI-generated code, you must:

1. **Read and understand** every line of the generated code
2. **Test it** in a real environment (see the Test Environment section of the PR template)
3. **Check for security issues** — AI models can generate subtly insecure code (e.g., path traversal, injection, credential exposure). Review carefully.
4. **Verify correctness** — AI-generated logic can be plausible-sounding but wrong. Validate the behavior, not just the syntax.

**PRs where it is clear the contributor has not read or tested the AI-generated code will be closed without review.**

### AI-Generated Code Quality Standards

AI-generated contributions are held to the **same quality bar** as human-written code:

- It must pass all CI checks (`make check`)
- It must be idiomatic Go and consistent with the existing codebase style
- It must not introduce unnecessary abstractions, dead code, or over-engineering
- It must include or update tests where appropriate

### Security Review

AI-generated code requires extra security scrutiny. Pay special attention to:

- **File path handling and sandbox escapes** (see commit `244eb0b` for a real example)
- **External input validation** in channel handlers and tool implementations
- **Credential or secret handling**
- **Command execution** (`exec.Command`, shell invocations)

If you are unsure whether a piece of AI-generated code is safe, say so in the PR — reviewers will help.

---

## Pull Request Process

### Before Opening a PR

**Checklist:**
- [ ] Run `make check` and ensure it passes locally
- [ ] Fill in the PR template completely, including the AI disclosure section
- [ ] Link any related issue(s) in the PR description
- [ ] Keep the PR focused. Avoid bundling unrelated changes together
- [ ] Update documentation if needed
- [ ] Add or update tests
- [ ] Update CHANGELOG.md for user-facing changes

### PR Template Sections

The PR template asks for:

1. **Description** — What does this change do and why?
2. **Type of Change** — Bug fix, feature, docs, or refactor
3. **AI Code Generation** — Disclosure of AI involvement (required)
4. **Related Issue** — Link to the issue this addresses
5. **Technical Context** — Reference URLs and reasoning (skip for pure docs PRs)
6. **Test Environment** — Hardware, OS, model/provider, and channels used for testing
7. **Evidence** — Optional logs or screenshots demonstrating the change works
8. **Checklist** — Self-review confirmation

### PR Size

**Prefer small, reviewable PRs:**
- A PR that changes 200 lines across 5 files is much easier to review than one that changes 2000 lines across 30 files
- If your feature is large, consider splitting it into a series of smaller, logically complete PRs

**Example of Good PR Splitting:**
```
PR 1: Add provider interface
PR 2: Implement OpenAI provider
PR 3: Implement Anthropic provider
PR 4: Add provider factory
```

### Example PR Workflow

```bash
# 1. Create branch
git checkout -b feat/new-provider

# 2. Make changes
# Edit files...

# 3. Stage and commit
git add pkg/providers/new_provider.go
git commit -m "feat(providers): Add new LLM provider

- Implement Provider interface
- Add authentication support
- Include comprehensive tests"

# 4. Push to fork
git push origin feat/new-provider

# 5. Open PR on GitHub
# Navigate to https://github.com/comgunner/picoclaw-agents/pulls
# Click "New Pull Request"
# Fill in template
```

---

## Branch Strategy

### Long-Lived Branches

| Branch | Purpose | Protection |
|--------|---------|------------|
| **`main`** | Active development | Requires 1+ maintainer approval |
| **`release/x.y`** | Stable releases | Strictly protected, no direct pushes |

### Requirements to Merge into `main`

A PR can only be merged when all of the following are satisfied:

1. **CI passes** — All GitHub Actions workflows (lint, test, build) must be green
2. **Reviewer approval** — At least one maintainer has approved the PR
3. **No unresolved review comments** — All review threads must be resolved
4. **PR template is complete** — Including AI disclosure and test environment

### Who Can Merge

**Only maintainers can merge PRs.** Contributors cannot merge their own PRs, even if they have write access.

### Merge Strategy

We use **squash merge** for most PRs to keep the `main` history clean and readable. Each merged PR becomes a single commit referencing the PR number:

```
feat: Add Ollama provider support (#491)
```

If a PR consists of multiple independent, well-separated commits that tell a clear story, a regular merge may be used at the maintainer's discretion.

### Release Branches

When a version is ready, maintainers cut a `release/x.y` branch from `main`. After that point:

- **New features are not backported.** The release branch receives no new functionality after it is cut.
- **Security fixes and critical bug fixes are cherry-picked.** If a fix in `main` qualifies (security vulnerability, data loss, crash), maintainers will cherry-pick the relevant commit(s) onto the affected `release/x.y` branch and issue a patch release.

If you believe a fix in `main` should be backported to a release branch, note it in the PR description or open a separate issue. The decision rests with the maintainers.

**Release branches have stricter protections than `main` and are never directly pushed to under any circumstances.**

---

## Code Review

### For Contributors

**Responsibilities:**
- Respond to review comments within a reasonable time (48 hours preferred)
- When you update a PR in response to feedback, briefly note what changed
- If you disagree with feedback, engage respectfully. Explain your reasoning; reviewers can be wrong too
- Do not force-push after a review has started — it makes it harder for reviewers to see what changed. Use additional commits instead; the maintainer will squash on merge

**Example Response:**
```markdown
@reviewer Thanks for the feedback! I've updated the code to:
- Use `sync.RWMutex` instead of `sync.Mutex` for better read performance
- Add error handling for edge case X
- Update tests to cover the new behavior
```

### For Reviewers

**Review For:**

1. **Correctness**
   - Does the code do what it claims?
   - Are there edge cases?
   - Are there race conditions?

2. **Security**
   - Especially for AI-generated code, tool implementations, and channel handlers
   - Check for path traversal, injection, credential exposure

3. **Architecture**
   - Is the approach consistent with the existing design?
   - Does this add unnecessary complexity?

4. **Simplicity**
   - Is there a simpler solution?
   - Does this introduce over-engineering?

5. **Tests**
   - Are the changes covered by tests?
   - Are existing tests still meaningful?

**Be constructive and specific:**
- ✅ "This could have a race condition if two goroutines call this concurrently — consider using a mutex here"
- ❌ "this looks wrong"

### Reviewer List

Once your PR is submitted, you can reach out to the assigned reviewers:

| Function | Reviewer |
|----------|----------|
| Provider | @yinwm |
| Channel | @yinwm |
| Agent | @lxowalle |
| Tools | @lxowalle |
| Skill | — |
| MCP | — |
| Optimization | @lxowalle |
| Security | — |
| AI CI | @imguoguo |
| UX | — |
| Document | — |

---

## Communication

### Where to Communicate

| Platform | Purpose |
|----------|---------|
| **GitHub Issues** | Bug reports, feature requests, design discussions |
| **GitHub Discussions** | General questions, ideas, community conversation |
| **Pull Request comments** | Code-specific feedback |
| **Discord** | [Coming Soon] |

### When in Doubt

**Open an issue before writing code.** It costs little and prevents wasted effort.

**Good Questions to Ask:**
- "Is this feature aligned with the project goals?"
- "Has someone already worked on this?"
- "What's the best approach for X?"
- "Can you review my design before I implement it?"

### Response Expectations

- **Maintainers:** Aim to respond within 48 hours
- **Contributors:** Respond to review comments within 48 hours
- **Community:** Be patient and understanding — everyone is volunteering their time

---

## A Note on the Project's AI-Driven Origin

PicoClaw's architecture was substantially designed and implemented with AI assistance, guided by human oversight. If you find something that looks odd or over-engineered, it may be an artifact of that process — opening an issue to discuss it is always welcome.

**We believe AI-assisted development done responsibly produces great results. We also believe humans must remain accountable for what they ship. These two beliefs are not in conflict.**

---

## Recognition

Contributors are recognized in the following ways:

1. **CHANGELOG.md** — Notable contributions are mentioned in the changelog
2. **GitHub Contributors Graph** — Visible on the repository
3. **Release Notes** — Major contributors may be mentioned in release notes
4. **Documentation** — Significant contributors may be listed in README

---

## Questions?

If you have questions about contributing, please:

1. Check this document first
2. Search existing issues and discussions
3. Open a new discussion on GitHub
4. Contact a maintainer directly

**Thank you for contributing to PicoClaw!** 🎉

---

*PicoClaw: Ultra-Efficient AI in Go. $10 Hardware · 10MB RAM · <1s Startup*
