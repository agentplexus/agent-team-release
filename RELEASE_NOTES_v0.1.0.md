# Release Notes - v0.1.0

## Overview

Initial release of `prepush`, a multi-language pre-push hook for Git repositories.

## Highlights

### Multi-Language Support

`prepush` automatically detects languages in your repository and runs appropriate checks:

| Language | Detection | Status |
|----------|-----------|--------|
| Go | `go.mod` | ✅ Full support |
| TypeScript | `package.json` + `tsconfig.json` | ✅ Full support |
| JavaScript | `package.json` | ✅ Full support |
| Python | `pyproject.toml`, `setup.py`, `requirements.txt` | 🔲 Detection only |
| Rust | `Cargo.toml` | 🔲 Detection only |
| Swift | `Package.swift` | 🔲 Detection only |

### Go Checks

- **Format**: `gofmt -l` to detect unformatted files
- **Lint**: `golangci-lint run` for comprehensive linting
- **Test**: `go test ./...` to run all tests
- **Replace**: Detect local `replace` directives in `go.mod`

### TypeScript/JavaScript Checks

- **Lint**: `eslint` or `npm run lint`
- **Format**: `prettier --check` or `npm run format:check`
- **Types**: `tsc --noEmit` for type checking
- **Test**: `npm run test`

### Configuration

Optional `.prepush.yaml` for per-language customization:

```yaml
verbose: true
languages:
  go:
    test: true
    lint: true
    coverage: true
    exclude_coverage: "cmd"
  typescript:
    enabled: true
    paths: ["frontend/"]
```

## Installation

```bash
go install github.com/grokify/prepush@latest
```

## Usage

```bash
# Run in current directory
prepush

# As git hook
ln -s $(which prepush) .git/hooks/pre-push
```

## Links

- [README](README.md)
- [CHANGELOG](CHANGELOG.md)
