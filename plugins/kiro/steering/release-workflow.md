# Release Workflow

This document describes the complete release workflow for software projects using release-agent.

## Overview

The release process follows these phases:

1. Pre-flight checks
2. Version determination
3. Validation
4. Changelog generation
5. Commit and push
6. CI verification
7. Tag creation

## Pre-flight Checks

Before starting a release:

```bash
# Verify required tools
which agent-team-release schangelog golangci-lint

# Ensure clean working directory
git status --porcelain

# Verify on correct branch
git branch --show-current
```

## Version Determination

Semantic versioning based on conventional commits:

| Commit Type | Version Bump | Example |
|-------------|--------------|---------|
| `feat:` | MINOR | 1.0.0 -> 1.1.0 |
| `fix:` | PATCH | 1.0.0 -> 1.0.1 |
| `feat!:` or `BREAKING CHANGE:` | MAJOR | 1.0.0 -> 2.0.0 |

Use schangelog to analyze commits:

```bash
schangelog parse-commits --since=$(git describe --tags --abbrev=0)
```

## Validation Checks

Run all validation before proceeding:

```bash
atrelease check --verbose
```

This runs:

- Build verification
- Test suite
- Linting (golangci-lint for Go)
- Format checking

## Changelog Generation

Generate changelog entries:

```bash
schangelog generate CHANGELOG.json -o CHANGELOG.md
```

## Release Execution

Full release with all steps:

```bash
# Dry run first
atrelease release v1.2.3 --dry-run --verbose

# Execute release
atrelease release v1.2.3 --verbose
```

## CI Verification

The release process waits for CI to pass:

1. Push commits to remote
2. Wait for GitHub Actions
3. Only tag after CI passes

## Tag Creation

Tags are created after CI verification:

```bash
git tag -a v1.2.3 -m "Release v1.2.3"
git push origin v1.2.3
```

## Error Recovery

If a step fails:

1. Review the error output
2. Fix the issue
3. Re-run from the failed step
4. Never force-tag if validation failed
