# Release Agent Team

[![Go CI][go-ci-svg]][go-ci-url]
[![Go Lint][go-lint-svg]][go-lint-url]
[![Go SAST][go-sast-svg]][go-sast-url]
[![Docs][docs-godoc-svg]][docs-godoc-url]
[![Docs][docs-mkdoc-svg]][docs-mkdoc-url]
[![Visualization][viz-svg]][viz-url]
[![License][license-svg]][license-url]

 [go-ci-svg]: https://github.com/plexusone/agent-team-release/actions/workflows/go-ci.yaml/badge.svg?branch=main
 [go-ci-url]: https://github.com/plexusone/agent-team-release/actions/workflows/go-ci.yaml
 [go-lint-svg]: https://github.com/plexusone/agent-team-release/actions/workflows/go-lint.yaml/badge.svg?branch=main
 [go-lint-url]: https://github.com/plexusone/agent-team-release/actions/workflows/go-lint.yaml
 [go-sast-svg]: https://github.com/plexusone/agent-team-release/actions/workflows/go-sast-codeql.yaml/badge.svg?branch=main
 [go-sast-url]: https://github.com/plexusone/agent-team-release/actions/workflows/go-sast-codeql.yaml
 [docs-godoc-svg]: https://pkg.go.dev/badge/github.com/plexusone/agent-team-release
 [docs-godoc-url]: https://pkg.go.dev/github.com/plexusone/agent-team-release
 [docs-mkdoc-svg]: https://img.shields.io/badge/Go-dev%20guide-blue.svg
 [docs-mkdoc-url]: https://plexusone.dev/agent-team-release
 [viz-svg]: https://img.shields.io/badge/Go-visualizaton-blue.svg
 [viz-url]: https://mango-dune-07a8b7110.1.azurestaticapps.net/?repo=plexusone%2Fagent-team-release
 [loc-svg]: https://tokei.rs/b1/github/plexusone/agent-team-release
 [repo-url]: https://github.com/plexusone/agent-team-release
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-url]: https://github.com/plexusone/agent-team-release/blob/main/LICENSE

**Multi-agent release preparation team for multi-language repositories.**

Release Agent Team validates code quality, generates changelogs, updates documentation, and manages the complete release lifecycle. It supports monorepos with multiple languages and integrates with Claude Code as an interactive subagent.

## Features

- 🔍 **Auto-detection**: Detects Go, TypeScript, JavaScript, Python, Rust, Swift
- ✅ **Validation checks**: Build, test, lint, format, security, documentation checks
- 📦 **Monorepo support**: Handles repositories with multiple languages
- 📝 **Changelog generation**: Integrates with schangelog for automated changelogs
- 📄 **Documentation updates**: Updates README badges and version references
- 🚀 **Release workflow**: Full release lifecycle with CI verification
- 💬 **Interactive mode**: Ask questions and propose fixes for Claude Code integration
- 📊 **Multiple output formats**: Human-readable, JSON, and TOON (Token-Oriented Object Notation)
- 🔌 **Claude Code plugin**: Available as a plugin with commands, skills, and agents

## Agent Workflow

Release Agent Team uses a DAG (Directed Acyclic Graph) workflow with specialized agents and automated fix loops:

```mermaid
graph TD
    PM[PM Agent<br/>version & scope] --> QA[QA Agent<br/>build, test, lint]
    PM --> DocsReview[Docs Reviewer<br/>validate docs]
    PM --> Security[Security Agent<br/>license, vulns, secrets]
    PM --> Release[Release Agent<br/>git status, CI config]

    QA -->|NO-GO| CodeFixer[Code-Fixer<br/>fix lint errors]
    CodeFixer -->|re-test| QA

    DocsReview -->|NO-GO| DocsWriter[Docs Writer<br/>create/update docs]
    DocsWriter -->|re-review| DocsReview

    QA --> Coordinator[Release Coordinator<br/>execute release]
    DocsReview --> Coordinator
    Security --> Coordinator
    Release --> Coordinator

    style PM fill:#e1f5fe,color:#01579b
    style QA fill:#fff3e0,color:#e65100
    style DocsReview fill:#e8f5e9,color:#1b5e20
    style DocsWriter fill:#c8e6c9,color:#1b5e20
    style CodeFixer fill:#ffe0b2,color:#e65100
    style Security fill:#fce4ec,color:#880e4f
    style Release fill:#f3e5f5,color:#4a148c
    style Coordinator fill:#fff8e1,color:#f57f17
```

### Validation Agents

| Agent | Role | Checks |
|-------|------|--------|
| **PM** | Product Management | Version recommendation, release scope, changelog quality, breaking changes |
| **QA** | Quality Assurance | Build, tests, lint, format, error handling, mod tidy |
| **Docs Reviewer** | Documentation Validation | README, CHANGELOG.json, release notes, mkdocs nav |
| **Security** | Security | LICENSE, vulnerability scan, dependency audit, secret detection |
| **Release** | Release Management | Version availability, git status, CI configuration |

### Fix Loop Agents

| Agent | Role | Triggered By |
|-------|------|--------------|
| **Code-Fixer** | Fix lint/build errors | QA NO-GO |
| **Docs Writer** | Create/update release docs | Docs Reviewer NO-GO |
| **Coordinator** | Orchestrate fix loops & release | All validations |

### Automated Fix Loops

The release-coordinator orchestrates two automated fix loops:

**QA Fix Loop** (qa → code-fixer → qa):

- Triggered when QA validation fails (lint errors, format issues)
- Code-fixer applies automatic fixes (errcheck, gosec, gofmt)
- Re-runs QA validation until GO or max attempts reached

**Docs Fix Loop** (docs-reviewer → docs-writer → docs-reviewer):

- Triggered when documentation is incomplete
- Docs-writer creates CHANGELOG entries, release notes, mkdocs nav
- Re-validates until all documentation checks pass

See [QA Fix Loop](https://plexusone.dev/agent-team-release/qa-fix-loop) and [Docs Fix Loop](https://plexusone.dev/agent-team-release/docs-fix-loop) for details.

### Loop Specifications

Fix loops are formally defined using REAL/VEAL patterns in `specs/loops/`:

| Loop | Type | Validator | Actor | max_attempts |
|------|------|-----------|-------|--------------|
| `qa-fix` | VEAL | qa | code-fixer | 3 |
| `docs-fix` | VEAL | docs-reviewer | docs-writer | 3 |
| `security-fix` | VEAL | security | code-fixer | 3 |
| `docs-next-version` | REAL | - | docs-writer | 5 |

Loops are bounded by `max_attempts` and escalate to human when exceeded.

### Workflow Guarantees

- **PM runs first** - Validates version and scope before other checks
- **QA, Docs, Security run in parallel** - Independent validation after PM approval
- **Fix loops auto-remediate** - Reduces manual intervention
- **Release runs after all validations** - Confirms release readiness
- **Coordinator executes last** - Only proceeds when all teams report GO

## Installation

```bash
go install github.com/plexusone/agent-team-release/cmd/atrelease@latest
```

### Homebrew

```bash
brew tap plexusone/tap
brew install atrelease
```

## Quick Start

```bash
# Run validation checks in current directory
atrelease check

# Run comprehensive validation with Go/No-Go report
atrelease validate --version=v1.0.0

# Execute full release workflow
atrelease release v1.0.0

# Generate changelog
atrelease changelog --since=v0.9.0

# Show version
atrelease version
```

## Commands

### `atrelease check`

Run validation checks for all detected languages.

```bash
atrelease check [directory]

# With options
atrelease check --verbose
atrelease check --no-test --no-lint
atrelease check --coverage
atrelease check --go-no-go  # NASA-style Go/No-Go report
```

### `atrelease validate`

Run comprehensive validation across all areas (QA, Documentation, Release, Security).

```bash
atrelease validate [directory]

# With version-specific checks
atrelease validate --version=v1.0.0

# Skip specific areas
atrelease validate --skip-qa --skip-docs --skip-security

# Team status report format (template-based)
atrelease validate --format team
```

### `atrelease release`

Execute the full release workflow.

```bash
atrelease release <version>

# Examples
atrelease release v1.0.0
atrelease release v1.0.0 --dry-run      # Preview without changes
atrelease release v1.0.0 --skip-ci      # Don't wait for CI
atrelease release v1.0.0 --verbose
```

**Release workflow steps:**

1. Validate version format and availability
2. Check working directory is clean
3. Run validation checks (build, test, lint, format)
4. Generate changelog via schangelog
5. Update roadmap via sroadmap
6. Create release commit
7. Push to remote
8. Wait for CI to pass
9. Create and push release tag

### `atrelease changelog`

Generate or update changelog using schangelog.

```bash
atrelease changelog [directory]
atrelease changelog --since=v0.9.0
atrelease changelog --dry-run
```

### `atrelease readme`

Update README badges and version references.

```bash
atrelease readme [directory]
atrelease readme --version=v1.0.0
atrelease readme --dry-run
```

### `atrelease roadmap`

Update roadmap using sroadmap.

```bash
atrelease roadmap [directory]
atrelease roadmap --dry-run
```

### `atrelease version`

Show version information.

```bash
atrelease version
```

## Global Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--verbose` | `-v` | Show detailed output |
| `--interactive` | `-i` | Enable interactive mode |
| `--json` | | Output as structured data |
| `--format` | | Output format: `toon`, `json`, or `team` (validate only) |

## Supported Languages

| Language | Detection | Checks |
|----------|-----------|--------|
| **Go** | `go.mod` | `go build`, `go mod tidy`, `gofmt`, `golangci-lint`, `go test`, local replace, untracked refs, error handling |
| **TypeScript** | `package.json` + `tsconfig.json` | `eslint`, `prettier`, `tsc --noEmit`, `npm test` |
| **JavaScript** | `package.json` | `eslint`, `prettier`, `npm test` |
| **Python** | `pyproject.toml`, `setup.py`, `requirements.txt` | Coming soon |
| **Rust** | `Cargo.toml` | Coming soon |
| **Swift** | `Package.swift` | Coming soon |

### Go Checks Detail

| Check | Type | Description |
|-------|------|-------------|
| no local replace | Hard | Fails if go.mod has local replace directives |
| mod tidy | Hard | Fails if go.mod/go.sum need updating |
| build | Hard | Fails if project doesn't compile |
| gofmt | Hard | Fails if code isn't formatted |
| golangci-lint | Hard | Fails if linter reports issues |
| tests | Hard | Fails if tests fail |
| error handling | Hard | Fails if errors are improperly discarded |
| untracked refs | Soft | Warns if tracked files reference untracked files |
| coverage | Soft | Reports coverage (requires `gocoverbadge`) |

### Validation Areas

The `validate` command checks four distinct areas:

| Area | Checks |
|------|--------|
| **QA** | Build, tests, lint, format, error handling compliance |
| **Documentation** | README, PRD, TRD, release notes, CHANGELOG |
| **Release** | Version validation, git status, CI configuration |
| **Security** | LICENSE file, vulnerability scan, dependency audit, secret detection |

## Configuration

Create `.releaseagent.yaml` in your repository root:

```yaml
# Global settings
verbose: false

# Language-specific settings
languages:
  go:
    enabled: true
    test: true
    lint: true
    format: true
    coverage: false
    exclude_coverage: "cmd"  # directories to exclude from coverage

  typescript:
    enabled: true
    paths: ["frontend/"]  # specific paths (empty = auto-detect)
    test: true
    lint: true
    format: true

  javascript:
    enabled: false  # disable for this repo
```

### Configuration Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `enabled` | bool | `true` | Enable/disable language checks |
| `paths` | []string | auto | Specific paths to check |
| `test` | bool | `true` | Run tests |
| `lint` | bool | `true` | Run linter |
| `format` | bool | `true` | Check formatting |
| `coverage` | bool | `false` | Show coverage (Go only) |

## Git Hook Integration

To run agent-team-release automatically before every `git push`:

```bash
# Create the hook
cat > .git/hooks/pre-push << 'EOF'
#!/bin/bash
exec atrelease check
EOF

# Make it executable
chmod +x .git/hooks/pre-push
```

**Bypassing the hook:**

```bash
git push --no-verify
```

## Output Formats

### Human-Readable (Default)

```
=== Summary ===
✓ Go: no local replace directives
✓ Go: mod tidy
✓ Go: build
✓ Go: gofmt
✓ Go: golangci-lint
✓ Go: tests

Passed: 6, Failed: 0, Skipped: 0
```

### JSON (`--json --format=json`)

Standard JSON output for programmatic consumption.

### TOON (`--json --format=toon`)

Token-Oriented Object Notation - approximately 8x more token-efficient than JSON, optimized for LLM consumption.

### Team Status Report (`--format team`)

Structured box report with per-team validation results:

```
╔════════════════════════════════════════════════════════════════════════════╗
║                             TEAM STATUS REPORT                             ║
╠════════════════════════════════════════════════════════════════════════════╣
║ Project: github.com/plexusone/agent-team-release                         ║
║ Target:  v0.3.0                                                            ║
╠════════════════════════════════════════════════════════════════════════════╣
║ RELEASE VALIDATION                                                         ║
╠════════════════════════════════════════════════════════════════════════════╣
║ qa-validation (qa)                                                         ║
║   build                    🟢 GO                                           ║
║   tests                    🟢 GO    42 tests passed                        ║
║   lint                     🟢 GO                                           ║
╠════════════════════════════════════════════════════════════════════════════╣
║ security-validation (security)                                             ║
║   license                  🟢 GO    MIT License                            ║
║   vulnerability-scan       🟡 WARN  1 deprecated                           ║
╠════════════════════════════════════════════════════════════════════════════╣
║                         🚀 TEAM: GO for v0.3.0 🚀                          ║
╚════════════════════════════════════════════════════════════════════════════╝
```

## Claude Code Plugin

Install as a Claude Code plugin for interactive release automation:

```bash
claude plugin add github:plexusone/agent-team-release/plugins/claude
```

The plugin includes:

- **Commands**: `/release-agent:release`, `/release-agent:check`, `/release-agent:changelog`, `/release-agent:version-next`
- **Skills**: Version analysis, commit classification
- **Agents**: Release coordinator subagent for orchestrating complete releases

See [plugins/claude/README.md](plugins/claude/README.md) for full plugin documentation.

### Interactive Mode

Use `--interactive` flag to enable Q&A mode where Release Agent can:

- Ask questions when issues arise
- Propose fixes for lint errors
- Get user approval before making changes

```bash
atrelease check --interactive
atrelease release v1.0.0 --interactive
```

## Examples

### Go Project

```
$ atrelease check
=== Pre-push Checks ===

Detecting languages...
  Found: go in .

Running Go checks...

=== Summary ===
✓ Go: no local replace directives
✓ Go: mod tidy
✓ Go: build
✓ Go: gofmt
✓ Go: golangci-lint
✓ Go: tests
✓ Go: untracked references

Passed: 7, Failed: 0, Skipped: 0

All pre-push checks passed!
```

### Comprehensive Validation

```
$ atrelease validate --version=v1.0.0

=== Release Validation: v1.0.0 ===

┌─────────────────────────────────────────┐
│            VALIDATION REPORT            │
├─────────────────────────────────────────┤
│ QA             │ 🟢 GO                  │
│ Documentation  │ 🟢 GO                  │
│ Release        │ 🟢 GO                  │
│ Security       │ 🟢 GO                  │
├─────────────────────────────────────────┤
│ STATUS: GO FOR LAUNCH                   │
└─────────────────────────────────────────┘
```

### Full Release

```
$ atrelease release v1.0.0 --verbose

[1/9] Validating version...
      ✓ Version v1.0.0 is valid and available

[2/9] Checking working directory...
      ✓ Working directory is clean

[3/9] Running validation checks...
      ✓ All checks passed

[4/9] Generating changelog...
      ✓ CHANGELOG.md updated

[5/9] Updating roadmap...
      ✓ ROADMAP.md updated

[6/9] Creating release commit...
      ✓ Created commit: chore(release): v1.0.0

[7/9] Pushing to remote...
      ✓ Pushed to origin/main

[8/9] Waiting for CI...
      ✓ CI passed

[9/9] Creating tag...
      ✓ Created and pushed tag v1.0.0

Release v1.0.0 complete!
```

## Dependencies

### Required

| Tool | Purpose |
|------|---------|
| `git` | Version control operations |
| `gh` | GitHub CLI for CI status |
| `releasekit` | Language validation (build, test, lint) |

### Language-Specific

These tools are invoked by `releasekit` for language validation:

| Tool | Language | Purpose |
|------|----------|---------|
| `go` | Go | Build and test |
| `golangci-lint` | Go | Linting |
| `node`, `npm` | TypeScript/JS | Build and test |
| `eslint` | TypeScript/JS | Linting |
| `prettier` | TypeScript/JS | Formatting |

### Optional

| Tool | Purpose |
|------|---------|
| `schangelog` | Changelog generation |
| `sroadmap` | Roadmap management |
| `gocoverbadge` | Coverage badge generation |
| `govulncheck` | Vulnerability scanning |

## Documentation

- [PRD.md](PRD.md) - Product requirements
- [TRD.md](TRD.md) - Technical requirements
- [ROADMAP.md](ROADMAP.md) - Development roadmap
- [CHANGELOG.md](CHANGELOG.md) - Release history

## Related Tools

- [releasekit](https://github.com/grokify/releasekit) - Release management toolkit for language validation
- [schangelog](https://github.com/grokify/structured-changelog) - Structured changelog management
- [sroadmap](https://github.com/grokify/structured-roadmap) - Structured roadmap management
- [gocoverbadge](https://github.com/grokify/gocoverbadge) - Generate coverage badges

## License

MIT License - see [LICENSE](LICENSE) for details.
