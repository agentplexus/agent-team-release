---
description: Execute full release workflow for version $ARGUMENTS
---

# Release

Execute the complete release workflow for the specified version.

## Usage

```
/release-agent:release v1.2.3
```

## Process

1. **Validate** version format and check it doesn't exist
2. **Check** working directory is clean
3. **Run** validation checks (build, test, lint, format)
4. **Generate** changelog via schangelog
5. **Update** roadmap via sroadmap
6. **Create** release commit
7. **Push** to remote
8. **Wait** for CI to pass
9. **Create** and push release tag

## Options

For dry-run (preview without changes):
```bash
atrelease release $ARGUMENTS --dry-run
```

For interactive mode (approve each step):
```bash
atrelease release $ARGUMENTS --interactive
```

For JSON/TOON output:
```bash
atrelease release $ARGUMENTS --json
```

## Example

User: `/release-agent:release v0.9.0`

This will execute: `atrelease release v0.9.0 --verbose`
