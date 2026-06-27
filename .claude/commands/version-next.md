---
description: Analyze commits and suggest next semantic version
---

# Version Next

Analyze commits and suggest next semantic version

## Process

1. Get current version from latest git tag
2. Parse commits since that tag
3. Classify commits by type
4. Determine version bump (major/minor/patch)
5. Suggest next version with reasoning

## Dependencies

- `git`
- `schangelog`

## Instructions

Analyze git history since the last tag and suggest the next semantic version based on conventional commits.

Given `vMAJOR.MINOR.PATCH`:

- Breaking changes bump MAJOR
- New features (feat:) bump MINOR
- Bug fixes (fix:) bump PATCH

## Usage

```
/agent-team-release:version-next
```

## Examples

Analyze and suggest version:

```
/agent-team-release:version-next
```

Output: `Current: v0.8.0, Suggested: v0.9.0 (3 new features)`
