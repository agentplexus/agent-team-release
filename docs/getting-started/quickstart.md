# Quick Start

This guide walks you through your first use of Release Agent.

## Your First Check

Navigate to a Go, TypeScript, or JavaScript project and run:

```bash
atrelease check
```

Release Agent automatically detects the language(s) in your repository and runs appropriate checks.

### Example Output

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
⚠ Go: untracked references (warning)
  main.go may reference untracked utils.go

Passed: 7, Failed: 0, Skipped: 0, Warnings: 1

Pre-push checks passed with warnings.
```

## Comprehensive Validation

For a full release readiness assessment:

```bash
atrelease validate --version=v1.0.0
```

This runs checks across all four validation areas:

- **QA** - Build, tests, lint, format
- **Documentation** - README, CHANGELOG, release notes
- **Release** - Version availability, git status
- **Security** - LICENSE, vulnerabilities, secrets

### Team Status Report

For a structured team report format:

```bash
atrelease validate --format team --version=v1.0.0
```

Output:

```
╔════════════════════════════════════════════════════════════════════════════╗
║                             TEAM STATUS REPORT                             ║
╠════════════════════════════════════════════════════════════════════════════╣
║ Project: github.com/yourorg/yourproject                                    ║
║ Target:  v1.0.0                                                            ║
╠════════════════════════════════════════════════════════════════════════════╣
║ qa-validation (qa)                                                         ║
║   build                    🟢 GO                                           ║
║   tests                    🟢 GO    42 tests passed                        ║
║   lint                     🟢 GO                                           ║
╠════════════════════════════════════════════════════════════════════════════╣
║                         🚀 TEAM: GO for v1.0.0 🚀                          ║
╚════════════════════════════════════════════════════════════════════════════╝
```

## Execute a Release

When you're ready to release:

```bash
# Preview what will happen
atrelease release v1.0.0 --dry-run

# Execute the release
atrelease release v1.0.0
```

The release workflow:

1. Validates version format and availability
2. Checks working directory is clean
3. Runs all validation checks
4. Generates changelog via schangelog
5. Updates roadmap via sroadmap
6. Creates release commit
7. Pushes to remote
8. Waits for CI to pass
9. Creates and pushes release tag

## Next Steps

- [Commands Reference](../commands/index.md) - Learn all available commands
- [Configuration](../configuration.md) - Customize Release Agent for your project
- [Git Hooks](../integrations/git-hooks.md) - Set up automatic pre-push validation
- [Claude Code Integration](../integrations/claude-code.md) - Use with AI assistants
