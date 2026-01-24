---
description: Execute full release workflow for the specified version
---

# Release

Execute full release workflow for the specified version

## Usage

```
/release <version>
```

## Arguments

- **version** (required): Semantic version for the release

## Process

1. Validate version format and check it doesn't exist
2. Check working directory is clean
3. Run validation checks (build, test, lint, format)
4. Generate changelog via schangelog
5. Update roadmap via sroadmap
6. Create release commit
7. Push to remote
8. Wait for CI to pass
9. Create and push release tag

## Dependencies

- `atrelease`
- `schangelog`
- `git`

## Instructions

Execute the complete release workflow for the specified version.

This command runs the full release process including validation, changelog generation, and git tagging.
