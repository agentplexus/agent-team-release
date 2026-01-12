# Release Agent - Product Requirements Document

## Overview

Release Agent is an autonomous release preparation tool that validates code quality, generates changelogs, updates documentation, and manages the complete release lifecycle for multi-language repositories.

## Problem Statement

Releasing software involves multiple repetitive steps that are easy to forget or execute incorrectly:

- Running tests and linters
- Checking for common issues (local replace directives, untracked files)
- Generating changelogs from commit history
- Updating documentation and badges
- Creating properly formatted commits and tags
- Ensuring CI passes before tagging

Currently, developers must remember and manually execute these steps, leading to inconsistent releases, forgotten changelog updates, and releases that fail CI after tagging.

## Solution

Release Agent automates the entire release workflow:

1. **Validation Phase** - Run all pre-push checks (tests, linting, formatting)
2. **Generation Phase** - Generate changelogs, update roadmaps, refresh badges
3. **Documentation Phase** - Update README, ensure docs reflect changes
4. **Release Phase** - Create commits, push, wait for CI, then tag

## Target Users

- Go developers maintaining open source libraries
- Teams with multi-language monorepos (Go, TypeScript, JavaScript, Python)
- Developers using Claude Code for AI-assisted development

## User Stories

### US-1: Pre-Push Validation

As a developer, I want to run all quality checks before pushing so that I catch issues before they reach CI.

**Status:** Complete

**Acceptance Criteria:**

- Detect languages in repository automatically
- Run language-appropriate checks (build, test, lint, format)
- Report results with clear pass/fail status
- Support configuration via `.releaseagent.yaml`

### US-2: Changelog Generation

As a developer, I want to automatically generate changelogs from my commit history so that I don't have to manually write release notes.

**Status:** Complete

**Acceptance Criteria:**

- Integrate with schangelog for commit parsing
- Generate CHANGELOG.md in standard format
- Support conventional commits categorization
- Allow manual override/editing before finalizing

### US-3: README Badge Updates

As a developer, I want my README badges (coverage, version) updated automatically so that documentation stays current.

**Status:** Complete

**Acceptance Criteria:**

- Update coverage badge via gocoverbadge
- Update version references in README
- Detect and update other common badges

### US-4: Interactive Mode for Claude Code

As a Claude Code user, I want release-agent to interactively ask me questions when issues arise so that I can make decisions about fixes.

**Status:** Complete

**Acceptance Criteria:**

- `--interactive` flag enables Q&A mode
- Present lint issues with proposed fixes
- Allow user to approve/modify/skip each fix
- Output structured JSON for Claude Code parsing

### US-5: Full Release Workflow

As a developer, I want a single command to handle the entire release so that I don't forget steps.

**Status:** Complete

**Acceptance Criteria:**

- `releaseagent release v1.2.0` runs complete workflow
- Validates all checks pass before proceeding
- Generates changelog for the release
- Creates release commit with proper message
- Pushes and waits for CI (optional)
- Creates and pushes tag only after CI passes

### US-6: Claude Code Plugin

As a Claude Code user, I want to invoke release-agent via `/release-agent` so that I can use it within my AI workflow.

**Status:** Complete

**Acceptance Criteria:**

- Plugin structure with commands, skills, and agents
- Subagent can ask user questions during execution
- Proper permission handling for file edits
- Clear summary of actions taken

## Features

### Phase 1: Foundation

| Feature | Priority | Status |
|---------|----------|--------|
| Rename to release-agent | P0 | Complete |
| Multi-language detection | P0 | Complete |
| Go checks (build, test, lint, format) | P0 | Complete |
| TypeScript/JavaScript checks | P0 | Complete |
| Configuration file support | P0 | Complete |
| Adopt cobra for CLI | P0 | Complete |
| `releaseagent check` subcommand | P0 | Complete |
| `releaseagent version` subcommand | P0 | Complete |

### Phase 2: Actions

| Feature | Priority | Status |
|---------|----------|--------|
| Action interface | P0 | Complete |
| schangelog integration | P0 | Complete |
| sroadmap integration | P1 | Complete |
| Coverage badge updates | P1 | Complete |
| README version updates | P1 | Complete |
| `releaseagent changelog` command | P0 | Complete |
| `releaseagent readme` command | P1 | Complete |
| `releaseagent roadmap` command | P1 | Complete |

### Phase 3: Release Workflow

| Feature | Priority | Status |
|---------|----------|--------|
| Workflow engine | P0 | Complete |
| `releaseagent release` command | P0 | Complete |
| CI status checking | P1 | Complete |
| Tag creation | P0 | Complete |
| Git wrapper package | P0 | Complete |
| GitHub release creation | P2 | Planned |

### Phase 4: Claude Code Integration

| Feature | Priority | Status |
|---------|----------|--------|
| Interactive prompter interface | P0 | Complete |
| Interactive mode (`--interactive`) | P0 | Complete |
| Lint fix proposals | P0 | Complete |
| JSON output mode | P1 | Complete |
| JSON protocol for questions/proposals | P1 | Complete |
| TOON output format | P1 | Complete |
| Team status report format (`--format team`) | P1 | Complete |
| Report package with text/template renderer | P1 | Complete |
| Claude Code plugin structure | P0 | Complete |
| AGENT.md for subagent | P0 | Complete |
| SessionStart dependency check hook | P1 | Complete |
| Publish to Claude Marketplace | P0 | Complete |

### Phase 5: Distribution

| Feature | Priority | Status |
|---------|----------|--------|
| GoReleaser configuration | P0 | Complete |
| Homebrew formula | P1 | Complete |

### Planned Features

| Feature | Priority | Status |
|---------|----------|--------|
| GitHub release creation | P2 | Planned |
| GitLab support | P2 | Planned |
| Python checks | P2 | Planned |
| Rust checks | P2 | Planned |
| Swift checks | P2 | Planned |

## Non-Goals

- IDE integration (VS Code extension, etc.) - rely on Claude Code
- Support for package managers beyond Go modules
- Deployment automation (use GoReleaser, etc.)
- CI/CD pipeline generation

## Success Metrics

| Metric | Target | Status |
|--------|--------|--------|
| Time to release | 50% reduction vs manual process | Achieved |
| Forgotten changelog updates | 0 (enforced by workflow) | Achieved |
| Failed releases due to CI | 0 (wait for CI before tagging) | Achieved |
| Claude Code plugin adoption | 100+ installs in first quarter | Tracking |

## Dependencies

| Dependency | Purpose | Required |
|------------|---------|----------|
| schangelog | Changelog generation | Yes (for changelog features) |
| sroadmap | Roadmap updates | No (optional) |
| golangci-lint | Go linting | Yes (for Go projects) |
| gocoverbadge | Coverage badges | No (optional) |
| gh CLI | GitHub CI status | Yes (for CI waiting) |
| govulncheck | Vulnerability scanning | No (optional) |

## Appendix

### Configuration Example

```yaml
# .releaseagent.yaml
verbose: false

# Validation settings
languages:
  go:
    enabled: true
    test: true
    lint: true
    format: true
    coverage: true
    exclude_coverage: "cmd"
  typescript:
    enabled: true
    paths: ["frontend/"]

# Release settings
release:
  changelog:
    enabled: true
    tool: schangelog
    output: CHANGELOG.md
  roadmap:
    enabled: false
    tool: sroadmap
  badges:
    coverage: true
    version: true
  ci:
    wait: true
    timeout: 600  # seconds
```

### CLI Reference

```bash
# Validation
releaseagent check [directory]
releaseagent check --no-test --no-lint
releaseagent check --go-no-go
releaseagent validate --version=v1.0.0
releaseagent validate --skip-qa --skip-docs --skip-security
releaseagent validate --format team  # Team status report

# Actions
releaseagent changelog [--since=v1.0.0]
releaseagent readme [--version=v1.1.0]
releaseagent roadmap

# Release workflow
releaseagent release v1.1.0
releaseagent release v1.1.0 --skip-ci
releaseagent release v1.1.0 --dry-run

# Interactive mode (for Claude Code)
releaseagent check --interactive
releaseagent release v1.1.0 --interactive

# Output formats
releaseagent check --json --format=toon
releaseagent release v1.1.0 --json --format=json

# Version
releaseagent version
```

### Validation Areas

The `validate` command runs comprehensive checks across four areas:

| Area | Checks |
|------|--------|
| QA | Build, tests, lint, format, error handling |
| Documentation | README, PRD, TRD, release notes, CHANGELOG |
| Release | Version validation, git status, CI configuration |
| Security | LICENSE, vulnerability scan, dependency audit, secrets |

### Output Formats

| Format | Flag | Description |
|--------|------|-------------|
| Human | (default) | Colored terminal output with symbols |
| JSON | `--json --format=json` | Standard JSON for programmatic use |
| TOON | `--json --format=toon` | Token-optimized format for LLMs (~8x more efficient) |
| Team | `--format team` | Template-based box report with per-team validation results |
