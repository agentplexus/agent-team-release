# changelog

Generate or update changelog using schangelog.

## Usage

```bash
atrelease changelog [directory] [flags]
```

## Description

The `changelog` command generates or updates the changelog by integrating with [schangelog](https://github.com/grokify/structured-changelog). It parses conventional commits and organizes them into a structured changelog format.

## Arguments

| Argument | Description | Default |
|----------|-------------|---------|
| `directory` | Directory to process | Current directory (`.`) |

## Flags

| Flag | Description |
|------|-------------|
| `--since` | Generate from this version/tag |
| `--dry-run` | Preview changes without writing |
| `--verbose`, `-v` | Show detailed output |

## Requirements

This command requires [schangelog](https://github.com/grokify/structured-changelog) to be installed:

```bash
go install github.com/grokify/structured-changelog/cmd/schangelog@latest
```

## Examples

```bash
# Generate changelog for all commits
atrelease changelog

# Generate since a specific version
atrelease changelog --since=v0.9.0

# Preview without writing
atrelease changelog --dry-run

# Verbose output
atrelease changelog --since=v0.9.0 --verbose
```

## Output Files

The command updates:

- `CHANGELOG.json` - Structured changelog data
- `CHANGELOG.md` - Human-readable changelog

## Conventional Commits

The changelog generator categorizes commits based on conventional commit prefixes:

| Prefix | Category |
|--------|----------|
| `feat:` | Added |
| `fix:` | Fixed |
| `docs:` | Documentation |
| `style:` | Changed |
| `refactor:` | Changed |
| `perf:` | Changed |
| `test:` | Tests |
| `build:` | Build |
| `ci:` | CI |
| `chore:` | Chore |

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Changelog generated successfully |
| 1 | Error generating changelog |
