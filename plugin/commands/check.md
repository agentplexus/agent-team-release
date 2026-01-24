---
description: Run validation checks on the current repository
---

# Check

Run all validation checks for the current repository without making any changes.

## Usage

```
/release-agent:check
```

## What It Checks

### Go Projects
- `go build ./...` - Compilation
- `go test ./...` - Unit tests
- `golangci-lint run` - Linting
- `gofmt` - Formatting
- `go mod tidy` - Dependencies
- Local replace directives in go.mod

### TypeScript/JavaScript Projects
- `npm test` or `yarn test` - Tests
- `eslint` - Linting
- `prettier --check` - Formatting
- `tsc --noEmit` - Type checking

## Output

Shows pass/fail status for each check with detailed output for failures.

## Example

```bash
atrelease check --verbose
```

Output:
```
=== Release Agent Checks ===
Language: Go

[PASS] go build ./...
[PASS] go test ./...
[PASS] golangci-lint run
[PASS] gofmt check
[PASS] go mod tidy check
[PASS] No local replace directives

All checks passed!
```
