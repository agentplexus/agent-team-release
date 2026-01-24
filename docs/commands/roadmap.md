# roadmap

Update roadmap using sroadmap.

## Usage

```bash
atrelease roadmap [directory] [flags]
```

## Description

The `roadmap` command updates the project roadmap by integrating with [sroadmap](https://github.com/grokify/structured-roadmap). It regenerates `ROADMAP.md` from `ROADMAP.json`.

## Arguments

| Argument | Description | Default |
|----------|-------------|---------|
| `directory` | Directory to process | Current directory (`.`) |

## Flags

| Flag | Description |
|------|-------------|
| `--dry-run` | Preview changes without writing |
| `--verbose`, `-v` | Show detailed output |

## Requirements

This command requires [sroadmap](https://github.com/grokify/structured-roadmap) to be installed:

```bash
go install github.com/grokify/structured-roadmap/cmd/sroadmap@latest
```

## Examples

```bash
# Update roadmap
atrelease roadmap

# Preview without writing
atrelease roadmap --dry-run

# Verbose output
atrelease roadmap --verbose
```

## Input/Output Files

| File | Description |
|------|-------------|
| `ROADMAP.json` | Structured roadmap data (input) |
| `ROADMAP.md` | Human-readable roadmap (output) |

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Roadmap updated successfully |
| 1 | Error updating roadmap |
