---
description: Run validation checks on the current repository
---

# Check

Run validation checks on the current repository

## Process

1. Detect project language(s)
2. Run language-specific build checks
3. Run test suite
4. Run linter
5. Check code formatting
6. Report pass/fail status for each check

## Dependencies

- `atrelease`

## Instructions

Run all validation checks for the current repository without making any changes.

## Go Projects
- `go build ./...` - Compilation
- `go test ./...` - Unit tests
- `golangci-lint run` - Linting
- `gofmt` - Formatting
- `go mod tidy` - Dependencies

## TypeScript/JavaScript Projects
- `npm test` or `yarn test` - Tests
- `eslint` - Linting
- `prettier --check` - Formatting
- `tsc --noEmit` - Type checking
