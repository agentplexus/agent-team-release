# Commands

Release Agent provides seven commands for different stages of the release lifecycle.

## Command Overview

| Command | Description |
|---------|-------------|
| [`check`](check.md) | Run validation checks for detected languages |
| [`validate`](validate.md) | Comprehensive Go/No-Go validation across all areas |
| [`release`](release.md) | Execute the full release workflow |
| [`changelog`](changelog.md) | Generate or update changelog |
| [`readme`](readme.md) | Update README badges and versions |
| [`roadmap`](roadmap.md) | Update roadmap using sroadmap |
| [`version`](version.md) | Show version information |

## Global Flags

These flags are available for all commands:

| Flag | Short | Description |
|------|-------|-------------|
| `--verbose` | `-v` | Show detailed output |
| `--interactive` | `-i` | Enable interactive mode |
| `--json` | | Output as structured data |
| `--format` | | Output format: `toon`, `json`, or `team` (validate only) |

## Common Workflows

### Pre-Push Validation

Run before pushing to catch issues early:

```bash
atrelease check
```

### Release Readiness

Check if the project is ready for release:

```bash
atrelease validate --version=v1.0.0
```

### Full Release

Execute the complete release workflow:

```bash
atrelease release v1.0.0
```

### Generate Documentation

Update changelog and documentation:

```bash
atrelease changelog --since=v0.9.0
atrelease readme --version=v1.0.0
atrelease roadmap
```
