# Docs Fix Loop

The release-coordinator orchestrates an automated **review/create/re-review loop** to ensure release documentation is complete before tagging.

## Overview

When documentation validation fails, the release-coordinator automatically invokes the docs-writer agent to create or update documentation, then re-runs validation. This loop continues until all documentation checks pass or maximum attempts are reached.

## Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           DOCS FIX LOOP                                     │
│                    (Orchestrated by release-coordinator)                    │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│                         ┌───────────────────┐                               │
│                         │ release-coordinator│                               │
│                         │   (orchestrator)  │                               │
│                         └─────────┬─────────┘                               │
│                                   │                                         │
│                    ┌──────────────┼──────────────┐                         │
│                    │              │              │                         │
│                    ▼              ▼              ▼                         │
│              ┌──────────┐  ┌──────────┐  ┌──────────┐                      │
│              │  Docs    │  │  Docs    │  │  Docs    │                      │
│              │ Reviewer │  │  Writer  │  │ Reviewer │                      │
│              │(validate)│  │ (create) │  │(re-test) │                      │
│              └────┬─────┘  └────┬─────┘  └────┬─────┘                      │
│                   │             │             │                            │
│                   ▼             ▼             ▼                            │
│              ┌─────────┐  ┌─────────────┐  ┌─────────┐                     │
│              │• README │  │• parse      │  │• README │                     │
│              │• CHANGE-│  │  commits    │  │• CHANGE-│                     │
│              │  LOG    │  │• update     │  │  LOG    │                     │
│              │• release│  │  CHANGELOG  │  │• release│                     │
│              │  notes  │  │• create     │  │  notes  │                     │
│              │• mkdocs │  │  release    │  │• mkdocs │                     │
│              │  nav    │  │  notes      │  │  nav    │                     │
│              └─────────┘  │• update nav │  └─────────┘                     │
│                           └─────────────┘                                  │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

## Flow Diagram

```
                              START
                                │
                                ▼
                    ┌───────────────────────┐
                    │   Determine Version   │
                    │  schangelog parse-    │
                    │  commits --since=tag  │
                    └───────────┬───────────┘
                                │
                                ▼
              ┌─────────────────────────────────────┐
              │       INVOKE DOCS REVIEWER          │
              │  • Check README.md exists           │
              │  • Validate CHANGELOG.json          │
              │  • Check release notes exist        │
              │  • Check mkdocs.yml nav             │
              └─────────────────┬───────────────────┘
                                │
                                ▼
                         ┌──────────────┐
                         │ Docs Status? │
                         └──────┬───────┘
                                │
                 ┌──────────────┼──────────────┐
                 │              │              │
                 ▼              │              ▼
            ┌────────┐         │         ┌────────┐
            │   GO   │         │         │ NO-GO  │
            └────┬───┘         │         └────┬───┘
                 │             │              │
                 │             │              ▼
                 │             │    ┌─────────────────┐
                 │             │    │ attempts < max? │
                 │             │    └────────┬────────┘
                 │             │             │
                 │             │      ┌──────┴──────┐
                 │             │      │             │
                 │             │      ▼             ▼
                 │             │    ┌────┐      ┌──────┐
                 │             │    │YES │      │  NO  │
                 │             │    └──┬─┘      └──┬───┘
                 │             │       │           │
                 │             │       ▼           │
                 │             │  ┌─────────────────────────┐
                 │             │  │  INVOKE DOCS WRITER     │
                 │             │  │  • Parse commits        │
                 │             │  │  • Update CHANGELOG.json│
                 │             │  │  • Generate CHANGELOG.md│
                 │             │  │  • Create release notes │
                 │             │  │  • Update mkdocs.yml    │
                 │             │  └───────────┬─────────────┘
                 │             │              │
                 │             │              │ loop back
                 │             │              │
                 │             └──────────────┘
                 │
                 ▼                             ▼
    ┌────────────────────────┐    ┌────────────────────────┐
    │  PROCEED TO RELEASE    │    │   RAISE TO HUMAN       │
    │  • Deploy gh-pages     │    │   • Report missing     │
    │  • Wait for CI         │    │   • Suggest manual fix │
    │  • Create tag          │    │   • List attempts      │
    │  • Push tag            │    └────────────────────────┘
    └────────────────────────┘
```

## Agent Responsibilities

### Docs Reviewer (Validation)

| Check | File/Command | Required |
|-------|--------------|----------|
| readme | `README.md` exists | Yes |
| changelog-json | `schangelog validate CHANGELOG.json` | Yes |
| changelog-md | `CHANGELOG.md` exists | Yes |
| release-notes | `docs/releases/vX.Y.Z.md` or `RELEASE_NOTES_vX.Y.Z.md` | Yes* |
| mkdocs | `mkdocs.yml` nav includes version | No |

*Release notes required for major/minor releases, optional for patch.

The docs-reviewer agent is **read-only** - it validates and reports but does not modify files.

### Docs Writer (Creation)

The docs-writer agent has **write access** and creates/updates documentation:

| Step | Action | Tool |
|------|--------|------|
| 1 | Parse commits | `schangelog parse-commits --since=<tag>` |
| 2 | Update CHANGELOG.json | Edit file with new version entry |
| 3 | Generate CHANGELOG.md | `schangelog generate CHANGELOG.json -o CHANGELOG.md` |
| 4 | Create release notes | Write `docs/releases/vX.Y.Z.md` |
| 5 | Update mkdocs nav | Edit `mkdocs.yml` |

### Release-Coordinator (Orchestration)

The release-coordinator invokes subagents using the Task tool:

```go
// Step 1: Validate docs
Task(subagent_type="docs-reviewer",
     prompt="Review documentation for v0.15.0 in /path/to/repo")

// Step 2: Create/update docs (if NO-GO)
Task(subagent_type="docs-writer",
     prompt="Create documentation for v0.15.0 (previous: v0.14.0). Findings: <findings>")

// Step 3: Re-validate
Task(subagent_type="docs-reviewer",
     prompt="Re-validate documentation for v0.15.0...")
```

## CHANGELOG.json Format

Use the structured-changelog format from `github.com/grokify/structured-changelog`:

```json
{
  "irVersion": "1.0",
  "project": "project-name",
  "repository": "https://github.com/org/repo",
  "versioning": "semver",
  "commitConvention": "conventional",
  "maintainers": ["maintainer@email.com"],
  "releases": [
    {
      "version": "v0.15.0",
      "date": "2026-06-27",
      "highlights": [
        { "description": "Key feature 1" },
        { "description": "Key feature 2" }
      ],
      "added": [
        { "description": "New feature X", "commit": "abc1234" }
      ],
      "fixed": [
        { "description": "Bug fix Y", "commit": "def5678" }
      ],
      "documentation": [
        { "description": "Updated docs for Z", "commit": "ghi9012" }
      ]
    }
  ]
}
```

### Category Mapping

| Commit Type | Changelog Category |
|-------------|-------------------|
| feat | added |
| fix | fixed |
| refactor | refactored |
| docs | documentation |
| chore(deps) | dependencies |
| test | tests |
| perf | changed |

## Configuration

### Max Attempts

Default: **3 create attempts** before raising to human.

Some documentation cannot be auto-generated:

- Complex migration guides
- Detailed architecture documentation
- Screenshots or diagrams
- Content requiring domain expertise

### Version Detection

Version format depends on language:

| Language | Detect By | Format |
|----------|-----------|--------|
| Go | `go.mod` | v1.2.3 |
| Node.js | `package.json` | 1.2.3 |
| Python | `pyproject.toml` | 1.2.3 |

## Usage

### Via Release-Coordinator

```bash
# Full release workflow (includes docs fix loop)
Task(subagent_type="release-coordinator",
     prompt="Run release workflow for v0.15.0 on /path/to/repo")

# Just docs fix loop
Task(subagent_type="release-coordinator",
     prompt="Run docs fix loop for v0.15.0 (previous v0.14.0) on /path/to/repo")
```

### Manual Steps

```bash
# Step 1: Parse commits
schangelog parse-commits --since=v0.14.0

# Step 2: Update CHANGELOG.json (manually or via editor)

# Step 3: Generate CHANGELOG.md
schangelog generate CHANGELOG.json -o CHANGELOG.md

# Step 4: Create release notes
# Write docs/releases/v0.15.0.md

# Step 5: Validate
schangelog validate CHANGELOG.json
```

## Example Output

```
╔════════════════════════════════════════════════════════════════════════════╗
║                           DOCS FIX LOOP                                    ║
╠════════════════════════════════════════════════════════════════════════════╣
║ Project: github.com/plexusone/omnivoice-core                               ║
║ Target:  v0.15.0                                                            ║
║ Previous: v0.14.0                                                           ║
╠════════════════════════════════════════════════════════════════════════════╣
║ ATTEMPT 1                                                                   ║
╠════════════════════════════════════════════════════════════════════════════╣
║ Docs Review:       🔴 NO-GO (3 issues)                                     ║
║   changelog-json   Missing v0.15.0 entry                                    ║
║   release-notes    Missing docs/releases/v0.15.0.md                         ║
║   mkdocs           v0.15.0 not in nav                                       ║
╠════════════════════════════════════════════════════════════════════════════╣
║ Docs Writer:       ✓ 4 updates applied                                      ║
║   CHANGELOG.json   Added v0.15.0 (25 commits, 17 items)                     ║
║   CHANGELOG.md     Regenerated                                              ║
║   release-notes    Created docs/releases/v0.15.0.md                         ║
║   mkdocs.yml       Added v0.15.0 to nav                                     ║
╠════════════════════════════════════════════════════════════════════════════╣
║ Docs Re-review:    🟢 GO                                                   ║
╠════════════════════════════════════════════════════════════════════════════╣
║                         ✓ DOCS FIX LOOP COMPLETE                            ║
║                           Proceeding to release                             ║
╚════════════════════════════════════════════════════════════════════════════╝
```

## See Also

- [QA Fix Loop](qa-fix-loop.md) - Build/test/lint automation
- [Orchestration](orchestration.md) - Multi-agent coordination
- [Configuration](configuration.md) - Tool configuration
