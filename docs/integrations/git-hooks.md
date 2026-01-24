# Git Hooks Integration

Release Agent can be used as a git pre-push hook to automatically validate code before pushing.

## Setup Options

### Option 1: Simple Script

Create `.git/hooks/pre-push` manually:

```bash
cat > .git/hooks/pre-push << 'EOF'
#!/bin/bash
exec atrelease check
EOF

chmod +x .git/hooks/pre-push
```

### Option 2: Symlink

Create a symlink to the Release Agent binary:

```bash
ln -sf $(which agent-team-release) .git/hooks/pre-push
```

Note: This runs `atrelease` with no arguments, which defaults to `check`.

### Option 3: Shared Hooks Directory

For team-wide hooks, use a shared hooks directory:

```bash
# Create shared hooks directory
mkdir -p .githooks

# Create the hook
cat > .githooks/pre-push << 'EOF'
#!/bin/bash
exec atrelease check
EOF

chmod +x .githooks/pre-push

# Configure git to use shared hooks
git config core.hooksPath .githooks
```

Commit `.githooks/` to your repository so all team members use the same hooks.

## Hook Behavior

When the pre-push hook runs:

1. Release Agent detects languages in your repository
2. Runs all enabled checks for each language
3. Exits with code 0 (success) if all hard checks pass
4. Exits with code 1 (failure) if any hard check fails

### Warnings Don't Block

Soft warnings (like untracked references) don't block the push:

```
=== Summary ===
✓ Go: build
✓ Go: tests
⚠ Go: untracked references (warning)

Passed: 6, Failed: 0, Skipped: 0, Warnings: 1

Pre-push checks passed with warnings.
```

## Bypassing the Hook

To push without running checks (use sparingly):

```bash
git push --no-verify
```

Or:

```bash
git push -n
```

## Customizing the Hook

### Skip Certain Checks

```bash
#!/bin/bash
# Skip tests for faster push (run in CI)
exec atrelease check --no-test
```

### Verbose Output

```bash
#!/bin/bash
exec atrelease check --verbose
```

### Different Config for Hooks

```bash
#!/bin/bash
# Use hook-specific config
RELEASEAGENT_CONFIG=.releaseagent.hooks.yaml exec atrelease check
```

## CI Integration

For complete coverage, run Release Agent in both:

1. **Pre-push hook** - Fast feedback during development
2. **CI pipeline** - Comprehensive checks with full test suite

### GitHub Actions Example

```yaml
name: Release Agent Checks
on: [push, pull_request]

jobs:
  check:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.21'
      - name: Install Release Agent
        run: go install github.com/agentplexus/agent-team-release/cmd/atrelease@latest
      - name: Run checks
        run: atrelease check
```

## Troubleshooting

### Hook Not Running

Check that the hook is executable:

```bash
ls -la .git/hooks/pre-push
# Should show -rwxr-xr-x
```

### Wrong Directory

The hook runs in the repository root. Use absolute paths if needed:

```bash
#!/bin/bash
cd "$(git rev-parse --show-toplevel)"
exec atrelease check
```

### Release Agent Not Found

Ensure Release Agent is in your PATH:

```bash
#!/bin/bash
export PATH="$HOME/go/bin:$PATH"
exec atrelease check
```

### Slow Hooks

For faster feedback, skip time-consuming checks:

```bash
#!/bin/bash
# Fast checks only - full checks run in CI
exec atrelease check --no-test --no-coverage
```

## Team Guidelines

### Recommended Setup

1. Use shared hooks directory (Option 3)
2. Commit `.githooks/` to repository
3. Document in README how to enable:

```markdown
## Development Setup

Enable git hooks:
\`\`\`bash
git config core.hooksPath .githooks
\`\`\`
```

### When to Bypass

Only bypass hooks (`--no-verify`) when:

- Pushing work-in-progress to a feature branch
- CI will catch any issues
- You understand the risks

Never bypass when pushing to `main` or release branches.
