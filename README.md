# prepush

[![Go Report Card](https://goreportcard.com/badge/github.com/grokify/prepush)](https://goreportcard.com/report/github.com/grokify/prepush)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

**Multi-language pre-push hook for Git repositories.**

`prepush` automatically detects languages in your repository and runs appropriate checks before you push. It supports monorepos with multiple languages.

## Features

- **Auto-detection**: Detects Go, TypeScript, JavaScript, Python, Rust, Swift
- **Monorepo support**: Handles repositories with multiple languages
- **Configurable**: Optional `.prepush.yaml` for customization
- **Fast**: Only runs checks for detected languages

## Installation

```bash
go install github.com/grokify/prepush@latest
```

## Usage

### Basic Usage

```bash
# Run in current directory
prepush

# Run on specific directory
prepush /path/to/repo

# Verbose output
prepush --verbose

# Skip specific checks
prepush --no-test
prepush --no-lint
prepush --no-format

# Show coverage (Go only)
prepush --coverage
```

### As Git Hook

Create `.git/hooks/pre-push`:

```bash
#!/bin/bash
exec prepush
```

Or symlink:

```bash
ln -s $(which prepush) .git/hooks/pre-push
```

## Supported Languages

| Language | Detection | Checks |
|----------|-----------|--------|
| **Go** | `go.mod` | `gofmt`, `golangci-lint`, `go test`, local replace check |
| **TypeScript** | `package.json` + `tsconfig.json` | `eslint`, `prettier`, `tsc --noEmit`, `npm test` |
| **JavaScript** | `package.json` | `eslint`, `prettier`, `npm test` |
| **Python** | `pyproject.toml`, `setup.py`, `requirements.txt` | Coming soon |
| **Rust** | `Cargo.toml` | Coming soon |
| **Swift** | `Package.swift` | Coming soon |

## Configuration

Create `.prepush.yaml` in your repository root:

```yaml
# Global settings
verbose: false

# Language-specific settings
languages:
  go:
    enabled: true
    test: true
    lint: true
    format: true
    coverage: false
    exclude_coverage: "cmd"  # directories to exclude from coverage

  typescript:
    enabled: true
    paths: ["frontend/"]  # specific paths (empty = auto-detect)
    test: true
    lint: true
    format: true

  javascript:
    enabled: false  # disable for this repo
```

### Configuration Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `enabled` | bool | `true` | Enable/disable language checks |
| `paths` | []string | auto | Specific paths to check |
| `test` | bool | `true` | Run tests |
| `lint` | bool | `true` | Run linter |
| `format` | bool | `true` | Check formatting |
| `coverage` | bool | `false` | Show coverage (Go only) |

## Examples

### Go Project

```
$ prepush
=== Pre-push Checks ===

Detecting languages...
  Found: go in .

Running Go checks...

=== Summary ===
✓ Go: no local replace directives
✓ Go: gofmt
✓ Go: golangci-lint
✓ Go: tests

Passed: 4, Failed: 0, Skipped: 0

All pre-push checks passed!
```

### Monorepo (Go + TypeScript)

```
$ prepush
=== Pre-push Checks ===

Detecting languages...
  Found: go in ./backend
  Found: typescript in ./frontend

Running Go checks...
Running TypeScript checks...

=== Summary ===
✓ Go: no local replace directives
✓ Go: gofmt
✓ Go: golangci-lint
✓ Go: tests
✓ TypeScript: eslint
✓ TypeScript: prettier
✓ TypeScript: type check
✓ TypeScript: tests

Passed: 8, Failed: 0, Skipped: 0

All pre-push checks passed!
```

## Related Tools

- [gocoverbadge](https://github.com/grokify/mogo/tree/main/testing/coverbadge) - Generate coverage badges for Go projects

## License

MIT License - see [LICENSE](LICENSE) for details.
