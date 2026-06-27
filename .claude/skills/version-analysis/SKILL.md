---
name: version-analysis
description: Analyzes git history to determine semantic version bumps based on conventional commits
triggers: [version, semver, bump]
dependencies: [git, schangelog]
---

# Version Analysis

Analyzes git history to determine semantic version bumps based on conventional commits

## Instructions

Guidelines for determining semantic versions based on changes.

## Semantic Versioning

Format: `MAJOR.MINOR.PATCH`

| Component | When to Increment |
|-----------|-------------------|
| MAJOR | Breaking changes, incompatible API changes |
| MINOR | New features, backward-compatible additions |
| PATCH | Bug fixes, backward-compatible fixes |

## Conventional Commits Mapping

### PATCH Version (Bug Fixes)

```
fix: resolve null pointer in parser
fix(auth): handle expired tokens correctly
perf: improve query performance by 20%
```

### MINOR Version (Features)

```
feat: add user export functionality
feat(api): implement batch processing endpoint
feat(cli): add --verbose flag support
```

### MAJOR Version (Breaking Changes)

```
feat!: redesign authentication API
fix!: change error response format

BREAKING CHANGE: The login endpoint now requires email instead of username
```

## Analysis Commands

Get commits since last tag:

```bash
schangelog parse-commits --since=$(git describe --tags --abbrev=0)
```

Output includes:

- Commit type classification
- Suggested version bump
- Changelog category mapping

## Decision Tree

```
Has BREAKING CHANGE or !: in any commit?
  YES -> MAJOR bump
  NO  -> Has feat: commits?
           YES -> MINOR bump
           NO  -> PATCH bump
```

## Pre-release Versions

For pre-releases, append suffix:

- Alpha: `v1.2.3-alpha.1`
- Beta: `v1.2.3-beta.1`
- RC: `v1.2.3-rc.1`

## Version Validation

Ensure new version:

- Is greater than current tag
- Follows semver format
- Doesn't already exist as a tag

