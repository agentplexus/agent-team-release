---
name: docs-release-check
description: Review and fix documentation for release readiness, loop until clean
triggers: [release, docs, pre-release, docs-check]
dependencies: [schangelog, git]
---

# Docs Release Check

Review and fix documentation for release readiness, loop until clean

## Instructions

# Documentation Release Check

Review and fix documentation for release readiness using a VEAL loop (docs-reviewer validates, docs-writer fixes, repeat until GO).

## When to Use

- Before tagging a release, to ensure all docs are updated
- After completing features, before cutting a version
- When `/release` or `/docs-release-check` is invoked

## Inputs

- `version` - Target version (e.g., `v1.2.0`). **Required.** See "Determine Versions First" below — do not omit this and rely on auto-detection unless you are driving this skill interactively via `/docs-release-check` (not via the `Workflow` tool).
- `previous` - Previous version tag. Optional; if omitted, resolved from `git describe --tags --abbrev=0`.

## Determine Versions First (do this before invoking anything)

The `version` you pass is what every check in this skill validates against — the changelog entry, the release notes filename, the mkdocs nav entry. **There is no safe "figure it out later" path.**

1. If the caller (user or orchestrating agent) already gave you a version, use it.
2. Otherwise, resolve it yourself before proceeding:
   - `git tag --sort=-v:refname | head -1` (or `git describe --tags --abbrev=0`) for the previous version
   - `schangelog parse-commits --since=<previous-tag>` to inspect commits since then
   - Apply the [version-analysis](../skills/version-analysis.md) rules (feat→minor, fix→patch, breaking→major) to suggest the next version
   - Confirm the suggested version with the user before proceeding — never guess silently

**Never invoke this skill's `Workflow` step without a resolved `version`.** The underlying `docs-release-check` Workflow accepts `args.version`; if it's omitted, the workflow defaults to the literal string `"unreleased"` and reviews against that — which is not a real version, so nothing meaningful is checked, `ready: true` comes back trivially on the first pass, and no CHANGELOG entry, release notes, or nav update is ever created. This failure is silent: the workflow reports success. Always resolve the version first, per steps 1–2 above, then pass it explicitly.

## Workflow

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

## Invoking via the Workflow Tool

In Claude Code, this skill runs as the built-in `docs-release-check` Workflow. Always pass `version` explicitly in `args` — see "Determine Versions First" above for why:

```
Workflow({
  name: "docs-release-check",
  args: {
    repoPath: "/absolute/path/to/repo",   // default "."
    version: "v1.2.0",                    // REQUIRED — see above
    previousVersion: "v1.1.0",            // optional, auto-detected from git tags if omitted
    features: "one feature per line, for context on what to document",
  },
})
```

`repoPath` alone is not a substitute for `version` — the workflow does not derive a version from it. If you're driving the release from a shell already `cd`'d into the target repo, still pass `repoPath` explicitly (do not rely on ambient working directory, which the harness does not guarantee is preserved between tool calls) and always pass `version`.

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
# Slash command: version auto-detected via version-analysis before the check runs
/docs-release-check

# Specify target version
/docs-release-check v1.2.0

# Specify both versions
/docs-release-check v1.2.0 v1.1.0
```

The bare `/docs-release-check` form is safe only because the command instructions resolve a version first. Do not carry that pattern over to a direct `Workflow({ name: "docs-release-check" })` call with no `args.version` — that path has no auto-detection step and will silently no-op (see "Determine Versions First" above).

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

