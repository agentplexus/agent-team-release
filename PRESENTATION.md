---
marp: true
theme: default
paginate: true
header: 'Release Agent'
footer: 'Multi-Agent Release Orchestration'
---

# Release Agent

## Multi-Agent Release Validation & Orchestration

Building reliable software releases with AI-powered specialist agents

---

# The Problem

## Release Day Chaos

- Manual checklists get forgotten or skipped
- Version numbers chosen inconsistently
- Changelogs incomplete or missing commit links
- Tests pass locally but fail in CI
- Security vulnerabilities slip through
- Documentation lags behind code

**Result:** Broken releases, hotfixes, and unhappy users

---

# The Vision

## What if releases could validate themselves?

Specialist AI agents that each own a piece of the release process:

| Agent | Responsibility |
|-------|---------------|
| **PM** | Version recommendation, scope validation |
| **QA** | Build, tests, lint, formatting |
| **Docs** | README, changelog, release notes |
| **Security** | Vulnerabilities, secrets, licenses |
| **Release** | Git state, CI, tagging |

---

# Architecture

## Hierarchical Multi-Agent Orchestration

```
                    ┌─────────────────┐
                    │  pm-validation  │
                    │   (version)     │
                    └────────┬────────┘
                             │
              ┌──────────────┼──────────────┐
              ▼              ▼              ▼
       ┌──────────┐   ┌──────────┐   ┌──────────┐
       │    qa    │   │   docs   │   │ security │
       └────┬─────┘   └────┬─────┘   └────┬─────┘
            │              │              │
            └──────────────┼──────────────┘
                           ▼
                    ┌─────────────────┐
                    │    release      │
                    └─────────────────┘
```

---

# Two-Phase Workflow

## Phase 1: Pre-Commit Review

- PM recommends version based on changes
- Agents review content quality
- CHANGELOG.json entries without commit hashes
- Human-in-the-loop for missing items

## Phase 2: Post-Commit Finalization

- Link commits to changelog entries
- Generate CHANGELOG.md
- Create and push tag

---

# The PM Agent

## Intelligent Version Recommendation

```
┌─────────────────────────────────────────────────────┐
│  1. Check commits since last tag                    │
│     └─► If commits exist → Analyze commits          │
│                                                     │
│  2. If no commits, check CHANGELOG.json             │
│     └─► If target version → Analyze entries         │
│                                                     │
│  3. If no commits AND no CHANGELOG entry            │
│     └─► Analyze uncommitted changes → Recommend     │
└─────────────────────────────────────────────────────┘
```

Adapts to your workflow: commit-first or intent-first

---

# Challenge #1

## The Commit Hash Chicken-and-Egg Problem

**Problem:** CHANGELOG.json needs commit hashes, but commits haven't happened yet during review phase

**Solution:** Two-phase workflow

```
Phase 1: Review           Phase 2: Finalize
─────────────────         ─────────────────
CHANGELOG.json            CHANGELOG.json
├─ version: "0.3.0"       ├─ version: "0.3.0"
├─ entries: [...]         ├─ entries: [...]
└─ commit: ""        →    └─ commit: "abc123"
```

---

# Challenge #2

## Version Format Inconsistency

**Problem:** Different ecosystems use different version formats

- Go: `v1.2.3` (v prefix required for modules)
- Node.js: `1.2.3` (no prefix)
- Python: `1.2.3` (PEP 440)

**Solution:** Language-aware detection

```
if exists("go.mod"):
    version_format = "v{major}.{minor}.{patch}"
else:
    version_format = "{major}.{minor}.{patch}"
```

---

# Challenge #3

## When to Require Release Notes?

**Problem:** Not every release needs detailed release notes

**Solution:** Requirement based on release type

| Release Type | Release Notes |
|--------------|---------------|
| Major (1.0.0 → 2.0.0) | **Required** |
| Minor (1.0.0 → 1.1.0) | **Required** |
| Patch (1.0.0 → 1.0.1) | Optional |

Human-in-the-loop prompts when missing for major/minor

---

# Challenge #4

## Report Formatting for Terminals

**Problem:** Status reports overflow terminal width

**Solution:** 78-character width constraint with guidelines

```
╔══════════════════════════════════════════════════════════════════════════╗
║                              QA VALIDATION                               ║
╠══════════════════════════════════════════════════════════════════════════╣
║ build              🟢 GO                                                 ║
║ tests              🟢 GO    35 tests passed                              ║
║ lint               🟢 GO                                                 ║
╚══════════════════════════════════════════════════════════════════════════╝
```

---

# Challenge #5

## Dependency Issues During QA

**Problem:** QA agent found build failure due to ngrok dependency

```
undefined: logext.RandId
```

**Solution:** Interactive fix mode

- QA agent diagnosed incompatible log15 version
- Auto-fixed: `go get github.com/inconshreveable/log15/v3@v3.0.0-testing.5`
- Ran `go mod tidy`
- Re-validated: all checks passed

---

# The HITL Pattern

## Human-in-the-Loop for Critical Decisions

```
═══════════════════════════════════════════════════════════════════════════
 CHANGELOG ENTRY REQUIRED
═══════════════════════════════════════════════════════════════════════════
 Target version v0.2.0 requires a CHANGELOG.json entry.

 Options:
   [1] Generate draft from uncommitted changes analysis
   [2] I'll provide changelog entries manually
   [3] Abort release

 Select option (1/2/3): _
═══════════════════════════════════════════════════════════════════════════
```

---

# Validation Report Example

```
╔══════════════════════════════════════════════════════════════════════════╗
║                            TEAM STATUS REPORT                            ║
╠══════════════════════════════════════════════════════════════════════════╣
║ Project: github.com/grokify/mcpruntime                                   ║
║ Target:  v0.2.0                                                          ║
╠══════════════════════════════════════════════════════════════════════════╣
║ pm-validation                                                            ║
║   version-recommendation   🟢 GO    v0.2.0 (minor bump)                  ║
║   changelog-quality        🟢 GO    Highlights present                   ║
║   breaking-changes         🟢 GO    None detected                        ║
╠══════════════════════════════════════════════════════════════════════════╣
║ qa-validation                                                            ║
║   build                    🟢 GO                                         ║
║   tests                    🟢 GO    35 tests passed                      ║
║   lint                     🟢 GO                                         ║
╚══════════════════════════════════════════════════════════════════════════╝
```

---

# Key Design Decisions

## What We Learned

1. **Spec-driven agents** - Each agent reads from `validation/specs/*.md`
2. **Structured data first** - CHANGELOG.json over CHANGELOG.md
3. **Parallel where possible** - QA, Docs, Security run concurrently
4. **Fix, don't just report** - Agents can auto-fix issues
5. **GO/NO-GO/WARN** - Clear status semantics

---

# Commit Conventions

## Breaking Up Large Changes

Instead of one big commit:
```
feat: add streaming, auth, tests, and docs  # Too broad
```

Split into topical commits:
```
feat: add OAuth 2.1 PKCE authentication
test: add OAuth and HTTP server tests
build: add ngrok and OAuth dependencies
docs: add v0.2.0 changelog entry
```

Order: implementation → tests → docs → housekeeping

---

# Tools & Dependencies

## The Release Agent Stack

| Tool | Purpose |
|------|---------|
| `schangelog` | Structured changelog management |
| `golangci-lint` | Code quality validation |
| `govulncheck` | Security vulnerability scanning |
| `git` | Version control operations |
| `gh` | GitHub CLI for PRs and releases |

---

# Results

## What We Achieved

- **Automated validation** of 25+ release criteria
- **Consistent versioning** based on semver rules
- **Complete changelogs** with commit linkage
- **Security scanning** before every release
- **Human oversight** at critical decision points

---

# Future Directions

## What's Next

- [ ] Claude Code Marketplace plugin distribution
- [ ] GitHub Actions integration
- [ ] Support for more languages (Node.js, Python, Rust)
- [ ] Custom agent definitions via YAML
- [ ] Release rollback automation

---

# Try It Yourself

## Getting Started

```bash
# Clone the release-agent
git clone https://github.com/grokify/release-agent

# Read the orchestration guide
cat teams/ORCHESTRATION.md

# Run validation on your project
# (via Claude Code with release-agent specs)
```

---

# Thank You

## Links

- **Release Agent:** github.com/grokify/release-agent
- **Structured Changelog:** github.com/grokify/structured-changelog
- **MCP Runtime:** github.com/grokify/mcpruntime

## Questions?
