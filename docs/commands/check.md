# check

Run validation checks for all detected languages.

## Usage

```bash
atrelease check [directory] [flags]
```

## Description

The `check` command runs pre-push validation checks for all detected languages in your repository. It automatically detects Go, TypeScript, JavaScript, Python, Rust, and Swift projects and runs appropriate checks for each.

## Arguments

| Argument | Description | Default |
|----------|-------------|---------|
| `directory` | Directory to check | Current directory (`.`) |

## Flags

| Flag | Description |
|------|-------------|
| `--verbose`, `-v` | Show detailed output |
| `--no-test` | Skip test execution |
| `--no-lint` | Skip linting |
| `--no-format` | Skip format checking |
| `--coverage` | Show coverage report (Go only) |
| `--go-no-go` | NASA-style Go/No-Go report |

## Go Checks

When Go is detected (`go.mod` present), the following checks run:

| Check | Type | Description |
|-------|------|-------------|
| no local replace | Hard | Fails if go.mod has local replace directives |
| mod tidy | Hard | Fails if go.mod/go.sum need updating |
| build | Hard | Fails if project doesn't compile |
| gofmt | Hard | Fails if code isn't formatted |
| golangci-lint | Hard | Fails if linter reports issues |
| tests | Hard | Fails if tests fail |
| error handling | Hard | Fails if errors are improperly discarded |
| untracked refs | Soft | Warns if tracked files reference untracked files |
| coverage | Soft | Reports coverage (requires `gocoverbadge`) |

## TypeScript/JavaScript Checks

When TypeScript or JavaScript is detected, the following checks run:

| Check | Type | Description |
|-------|------|-------------|
| eslint | Hard | Fails if linter reports issues |
| prettier | Hard | Fails if code isn't formatted |
| tsc --noEmit | Hard | TypeScript type checking |
| npm test | Hard | Fails if tests fail |

## Examples

```bash
# Check current directory
atrelease check

# Check specific directory
atrelease check ./myproject

# Verbose output
atrelease check --verbose

# Skip tests (faster for quick checks)
atrelease check --no-test

# Show coverage report
atrelease check --coverage

# NASA-style Go/No-Go report
atrelease check --go-no-go
```

## Output

### Default Format

```
=== Pre-push Checks ===

Detecting languages...
  Found: go in .

Running Go checks...

=== Summary ===
✓ Go: no local replace directives
✓ Go: mod tidy
✓ Go: build
✓ Go: gofmt
✓ Go: golangci-lint
✓ Go: tests
✓ Go: error handling compliance

Passed: 7, Failed: 0, Skipped: 0

All pre-push checks passed!
```

### With Warnings

```
=== Summary ===
✓ Go: no local replace directives
✓ Go: mod tidy
✓ Go: build
⚠ Go: untracked references (warning)
  main.go may reference untracked utils.go

Passed: 6, Failed: 0, Skipped: 0, Warnings: 1

Pre-push checks passed with warnings.
```

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | All checks passed (warnings don't affect exit code) |
| 1 | One or more checks failed |
