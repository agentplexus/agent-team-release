---
name: commit-classification
description: Classifies commits according to conventional commits specification and maps them to changelog categories
triggers: [commit, classify, conventional]
dependencies: [git, schangelog]
---

# Commit Classification

Classifies commits according to conventional commits specification and maps them to changelog categories

## Instructions

Guidelines for classifying commits using Conventional Commits specification.

## Commit Types

| Type | Description | Changelog Category |
|------|-------------|-------------------|
| `feat` | New feature | Added |
| `fix` | Bug fix | Fixed |
| `docs` | Documentation only | Documentation |
| `style` | Formatting, no code change | - |
| `refactor` | Code change, no feature/fix | Changed |
| `perf` | Performance improvement | Performance |
| `test` | Adding/fixing tests | - |
| `build` | Build system changes | - |
| `ci` | CI configuration | - |
| `chore` | Maintenance tasks | - |

## Changelog Categories

Map commit types to Keep a Changelog categories:

| Changelog Category | Commit Types |
|-------------------|--------------|
| Added | `feat` |
| Changed | `refactor`, `perf` |
| Deprecated | Manual marking |
| Removed | Manual marking |
| Fixed | `fix` |
| Security | `fix` with security scope |

## Scope Guidelines

Scopes provide context:

```
feat(auth): add OAuth2 support
fix(api): handle timeout errors
docs(readme): update installation steps
```

Common scopes:

- `api` - API changes
- `cli` - Command-line interface
- `auth` - Authentication
- `db` - Database
- `ui` - User interface
- `core` - Core functionality

## Breaking Changes

Mark breaking changes with:

1. `!` after type: `feat!: new API`
2. Footer: `BREAKING CHANGE: description`

Both trigger MAJOR version bump.

## Commit Message Format

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

Examples:

```
feat(auth): add OAuth2 login support

Implements OAuth2 flow with support for Google and GitHub providers.

Closes #123
```

```
fix!: change error response format

BREAKING CHANGE: Error responses now use RFC 7807 format.
Migration guide: https://example.com/migrate
```

## Classification Tools

Use schangelog for automatic classification:

```bash
schangelog parse-commits --since=v1.0.0
```

Output provides:

- Parsed type and scope
- Suggested changelog category
- Breaking change detection

