# validate

Run comprehensive release validation across all areas of responsibility.

## Usage

```bash
atrelease validate [directory] [flags]
```

## Description

The `validate` command performs comprehensive Go/No-Go validation across four distinct areas: QA, Documentation, Release, and Security. It assumes Engineering and Product have already signed off and validates the remaining areas.

## Arguments

| Argument | Description | Default |
|----------|-------------|---------|
| `directory` | Directory to validate | Current directory (`.`) |

## Flags

| Flag | Description |
|------|-------------|
| `--version` | Target release version (e.g., v1.0.0) |
| `--skip-qa` | Skip QA validation |
| `--skip-docs` | Skip documentation validation |
| `--skip-security` | Skip security validation |
| `--format` | Output format: `default` or `team` |
| `--verbose`, `-v` | Show detailed output |

## Validation Areas

### QA Area

Build, tests, lint, format, and error handling compliance.

| Check | Description |
|-------|-------------|
| build | Project compiles successfully |
| tests | All tests pass |
| lint | No linter issues |
| format | Code is properly formatted |
| error handling | No improperly discarded errors |

### Documentation Area

README, PRD, TRD, release notes, and changelog.

| Check | Description |
|-------|-------------|
| README.md | README exists and is not empty |
| PRD.md | Product requirements document exists |
| TRD.md | Technical requirements document exists |
| release notes | Version-specific release notes exist |
| CHANGELOG.md | Changelog exists |
| MkDocs site | MkDocs documentation (optional) |

### Release Area

Version validation, git status, and CI configuration.

| Check | Description |
|-------|-------------|
| version available | Git tag doesn't already exist |
| git clean | Working directory has no uncommitted changes |
| git remote | Remote repository is configured |
| CI configuration | GitHub Actions or similar configured |

### Security Area

LICENSE, vulnerability scan, dependency audit, and secret detection.

| Check | Description |
|-------|-------------|
| LICENSE | License file exists |
| vulnerability scan | No known vulnerabilities (govulncheck) |
| dependency audit | Dependencies are current |
| secret detection | No hardcoded secrets found |

## Examples

```bash
# Basic validation
atrelease validate

# With version-specific checks
atrelease validate --version=v1.0.0

# Skip QA (already verified manually)
atrelease validate --skip-qa

# Team status report format
atrelease validate --format team

# Combine options
atrelease validate --version=v1.0.0 --format team --verbose
```

## Output Formats

### Default Format

```
╔════════════════════════════════════════════════════════════════════════╗
║           RELEASE VALIDATION: v1.0.0                                   ║
╠════════════════════════════════════════════════════════════════════════╣
║ 🟢 GO       QA                                                         ║
║ ──────────────────────────────────────────────────────────────────── ║
║   🟢 GO     Go: build                                                  ║
║   🟢 GO     Go: tests                                                  ║
║   🟢 GO     Go: lint                                                   ║
╠════════════════════════════════════════════════════════════════════════╣
║                       🚀 ALL SYSTEMS GO 🚀                             ║
╚════════════════════════════════════════════════════════════════════════╝
```

### Team Status Report (`--format team`)

```
╔════════════════════════════════════════════════════════════════════════════╗
║                             TEAM STATUS REPORT                             ║
╠════════════════════════════════════════════════════════════════════════════╣
║ Project: github.com/grokify/release-agent                                  ║
║ Target:  v1.0.0                                                            ║
╠════════════════════════════════════════════════════════════════════════════╣
║ RELEASE VALIDATION                                                         ║
╠════════════════════════════════════════════════════════════════════════════╣
║ qa-validation (qa)                                                         ║
║   build                    🟢 GO                                           ║
║   tests                    🟢 GO    42 tests passed                        ║
║   lint                     🟢 GO                                           ║
╠════════════════════════════════════════════════════════════════════════════╣
║ security-validation (security)                                             ║
║   license                  🟢 GO    MIT License                            ║
║   vulnerability-scan       🟡 WARN  1 deprecated                           ║
╠════════════════════════════════════════════════════════════════════════════╣
║                         🚀 TEAM: GO for v1.0.0 🚀                          ║
╚════════════════════════════════════════════════════════════════════════════╝
```

## Status Icons

| Icon | Status | Meaning |
|------|--------|---------|
| 🟢 | GO | Check passed |
| 🟡 | WARN | Warning (doesn't block) |
| 🔴 | NO-GO | Check failed (blocks release) |
| ⚪ | SKIP | Check skipped |

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | All checks passed (GO) |
| 1 | One or more checks failed (NO-GO) |
