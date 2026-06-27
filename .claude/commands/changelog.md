---
description: Generate changelog entries for the specified version
---

# Changelog

Generate changelog entries for the specified version

## Usage

```
/changelog <version>
```

## Arguments

- **version** (required): 

## Process

1. Get current version from latest git tag
2. Parse commits since that tag using schangelog
3. Classify commits by type (feat, fix, docs, etc.)
4. Generate changelog entries in CHANGELOG.json
5. Regenerate CHANGELOG.md

## Dependencies

- `schangelog`
- `git`

## Instructions

Generate changelog entries for the specified version based on commits since the last tag.

Uses schangelog to parse conventional commits and generate structured changelog entries.

## Usage

```
/agent-team-release:changelog <version>
```

## Arguments

- **version** (required): Version to generate changelog for (e.g., v1.2.3)

## Examples

Generate changelog:

```
/agent-team-release:changelog v0.9.0
```

Updates CHANGELOG.json and CHANGELOG.md with entries for the specified version.
