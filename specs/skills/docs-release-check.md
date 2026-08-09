---
name: docs-release-check
description: Review and fix documentation for release readiness, loop until clean
triggers: [release, docs, pre-release, docs-check]
dependencies: [schangelog, git]
model: sonnet
tools: [Read, Grep, Glob, Bash, Edit, Write, Agent]
---

# Documentation Release Check

Review and fix documentation for release readiness using a VEAL loop (docs-reviewer validates, docs-writer fixes, repeat until GO).

## When to Use

- Before tagging a release, to ensure all docs are updated
- After completing features, before cutting a version
- When `/release` or `/docs-release-check` is invoked

## Inputs

The skill accepts optional arguments:

- `version` - Target version (e.g., `v1.2.0`). If omitted, uses version analysis to suggest next version.
- `previous` - Previous version tag. If omitted, uses latest git tag.

## Workflow

### Phase 1: Determine Versions

1. If `version` not provided:
   - Run `schangelog parse-commits --since=<latest-tag>` to analyze commits
   - Suggest version based on conventional commits (feat→minor, fix→patch, breaking→major)
   - Confirm with user before proceeding

2. If `previous` not provided:
   - Use `git describe --tags --abbrev=0` to get latest tag

### Phase 2: VEAL Loop (docs-fix)

Execute the docs-fix loop with max 3 iterations:

```
┌─────────────────────────────────────────────────────────────┐
│                      VEAL LOOP: docs-fix                    │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│   ┌─────────────┐                                          │
│   │ docs-reviewer│ ◄── Validates documentation state       │
│   │   (haiku)    │     Reports GO/NO-GO with findings      │
│   └──────┬──────┘                                          │
│          │                                                  │
│          ▼                                                  │
│     ┌────────┐                                             │
│     │  GO?   │───YES──► Exit loop, proceed to Phase 3     │
│     └────┬───┘                                             │
│          │ NO                                               │
│          ▼                                                  │
│   ┌─────────────┐                                          │
│   │ docs-writer │ ◄── Fixes issues from findings          │
│   │   (sonnet)  │     Creates/updates documentation       │
│   └──────┬──────┘                                          │
│          │                                                  │
│          ▼                                                  │
│   ┌─────────────┐                                          │
│   │  Attempt    │───< max?──► Escalate to human           │
│   │  Counter    │                                          │
│   └──────┬──────┘                                          │
│          │                                                  │
│          └──────────────────────────────────────┐          │
│                                                  ▼          │
│                                           [Loop back]       │
└─────────────────────────────────────────────────────────────┘
```

### Phase 3: Report and Gate

After loop completes:

1. **If GO**: Report success, documentation is release-ready
2. **If NO-GO after max attempts**: Escalate to human with remaining issues

## Checks Performed

The docs-reviewer validates:

| Check | Required | Description |
|-------|----------|-------------|
| readme | Yes | README.md exists with adequate content |
| readme-content | Yes | README documents new features from `feat` commits |
| changelog-json | Yes | `schangelog validate CHANGELOG.json` passes |
| changelog-md | Yes | CHANGELOG.md exists with target version |
| release-notes | Yes* | `docs/releases/vX.Y.Z.md` exists |
| mkdocs | No | mkdocs.yml nav includes target version |
| mkdocs-nav-complete | No | All `docs/**/*.md` files have nav entries |

*Release notes required for major/minor, optional for patch.

## Agent Invocations

### docs-reviewer (Validator)

```
Invoke docs-reviewer agent with:
- Target Version: ${target_version}
- Previous Version: ${previous_version}
- Repo Directory: ${cwd}

Expect GO/NO-GO report with structured findings.
```

### docs-writer (Actor)

```
Invoke docs-writer agent with:
- Target Version: ${target_version}
- Previous Version: ${previous_version}
- Repo Directory: ${cwd}
- Findings: ${docs_reviewer_findings}

Expect documentation updates and completion report.
```

## Success Criteria

All documentation checks pass with GO status:

- README.md exists and documents new features
- CHANGELOG.json is valid with target version entry
- CHANGELOG.md is generated and current
- Release notes exist for target version (major/minor)
- MkDocs navigation is updated (if applicable)

## Escalation

If the loop exceeds 3 attempts without reaching GO status:

1. Report remaining NO-GO checks
2. List what was attempted
3. Ask human for guidance:
   - Proceed anyway (acknowledge doc gaps)
   - Manual intervention required
   - Abort release

## Example Usage

```bash
# Auto-detect versions
/docs-release-check

# Specify target version
/docs-release-check v1.2.0

# Specify both versions
/docs-release-check v1.2.0 v1.1.0
```

## Integration with Release Skill

This skill is designed to be invoked as part of a larger release workflow:

```
/release v1.2.0
  ├── /docs-release-check v1.2.0  ◄── This skill
  ├── /golangci-lint
  ├── topical conventional commits
  ├── wait for CI
  └── tag and push
```
