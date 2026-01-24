# Release Agent

**Autonomous release preparation agent for multi-language repositories.**

Release Agent validates code quality, generates changelogs, updates documentation, and manages the complete release lifecycle. It supports monorepos with multiple languages and integrates with Claude Code as an interactive subagent.

## Key Features

- **Auto-detection** - Detects Go, TypeScript, JavaScript, Python, Rust, Swift
- **Validation checks** - Build, test, lint, format, security, documentation checks
- **Monorepo support** - Handles repositories with multiple languages
- **Changelog generation** - Integrates with schangelog for automated changelogs
- **Documentation updates** - Updates README badges and version references
- **Release workflow** - Full release lifecycle with CI verification
- **Interactive mode** - Ask questions and propose fixes for Claude Code integration
- **Multiple output formats** - Human-readable, JSON, and TOON (Token-Oriented Object Notation)
- **Claude Code plugin** - Available as a plugin with commands, skills, and agents

## Quick Example

```bash
# Run validation checks
atrelease check

# Comprehensive Go/No-Go validation
atrelease validate --version=v1.0.0

# Execute full release workflow
atrelease release v1.0.0
```

## Validation Areas

Release Agent validates four distinct areas:

| Area | Checks |
|------|--------|
| **QA** | Build, tests, lint, format, error handling compliance |
| **Documentation** | README, PRD, TRD, release notes, CHANGELOG |
| **Release** | Version validation, git status, CI configuration |
| **Security** | LICENSE file, vulnerability scan, dependency audit, secret detection |

## Supported Languages

| Language | Detection | Status |
|----------|-----------|--------|
| **Go** | `go.mod` | Full support |
| **TypeScript** | `package.json` + `tsconfig.json` | Full support |
| **JavaScript** | `package.json` | Full support |
| **Python** | `pyproject.toml`, `setup.py` | Detection only |
| **Rust** | `Cargo.toml` | Detection only |
| **Swift** | `Package.swift` | Detection only |

## Get Started

- [Installation](getting-started/installation.md) - Install Release Agent
- [Quick Start](getting-started/quickstart.md) - Your first validation run
- [Commands](commands/index.md) - All available commands
- [Configuration](configuration.md) - Configure for your project

## Building the Documentation

To build the docs locally:

```bash
pip install mkdocs-material
mkdocs serve
```

To deploy to GitHub Pages:

```bash
mkdocs gh-deploy
```
