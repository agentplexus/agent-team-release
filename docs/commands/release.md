# release

Execute the full release workflow.

## Usage

```bash
atrelease release <version> [flags]
```

## Description

The `release` command orchestrates the complete release workflow, from validation through tagging. It ensures all checks pass, generates documentation, waits for CI, and only tags after everything succeeds.

## Arguments

| Argument | Description | Required |
|----------|-------------|----------|
| `version` | Release version (e.g., v1.0.0) | Yes |

## Flags

| Flag | Description |
|------|-------------|
| `--dry-run` | Preview what would happen without making changes |
| `--skip-ci` | Don't wait for CI to pass |
| `--skip-changelog` | Don't generate changelog |
| `--skip-roadmap` | Don't update roadmap |
| `--verbose`, `-v` | Show detailed output |
| `--interactive`, `-i` | Enable interactive mode |

## Workflow Steps

The release command executes these 9 steps in order:

| Step | Action | Description |
|------|--------|-------------|
| 1 | Validate Version | Check version format and availability |
| 2 | Check Directory | Ensure working directory is clean |
| 3 | Run Checks | Execute all validation checks |
| 4 | Generate Changelog | Update CHANGELOG via schangelog |
| 5 | Update Roadmap | Update ROADMAP via sroadmap |
| 6 | Create Commit | Create release commit |
| 7 | Push | Push to remote repository |
| 8 | Wait for CI | Poll GitHub Actions until pass/fail |
| 9 | Create Tag | Create and push release tag |

## Examples

```bash
# Execute full release
atrelease release v1.0.0

# Preview without making changes
atrelease release v1.0.0 --dry-run

# Skip CI waiting (dangerous - tag before CI passes)
atrelease release v1.0.0 --skip-ci

# Verbose output
atrelease release v1.0.0 --verbose

# Interactive mode (for Claude Code)
atrelease release v1.0.0 --interactive
```

## Output

### Successful Release

```
[1/9] Validating version...
      ✓ Version v1.0.0 is valid and available

[2/9] Checking working directory...
      ✓ Working directory is clean

[3/9] Running validation checks...
      ✓ All checks passed

[4/9] Generating changelog...
      ✓ CHANGELOG.md updated

[5/9] Updating roadmap...
      ✓ ROADMAP.md updated

[6/9] Creating release commit...
      ✓ Created commit: chore(release): v1.0.0

[7/9] Pushing to remote...
      ✓ Pushed to origin/main

[8/9] Waiting for CI...
      ⏳ Checking CI status...
      ✓ CI passed

[9/9] Creating tag...
      ✓ Created and pushed tag v1.0.0

Release v1.0.0 complete!
```

### Dry Run

```
[DRY RUN] Would execute the following:

[1/9] Validate version v1.0.0
[2/9] Check working directory
[3/9] Run validation checks
[4/9] Generate changelog
[5/9] Update roadmap
[6/9] Create commit: chore(release): v1.0.0
[7/9] Push to origin/main
[8/9] Wait for CI
[9/9] Create tag v1.0.0

No changes made.
```

## CI Waiting

The release command waits for CI to pass before creating the tag. This prevents tagging code that fails CI.

### Timeout

By default, CI waiting times out after 10 minutes. You can skip CI waiting with `--skip-ci`, but this is not recommended.

### Supported CI Systems

- GitHub Actions (via `gh` CLI)

## Interactive Mode

With `--interactive`, Release Agent can:

- Ask questions when issues arise
- Propose fixes for problems
- Get approval before making changes

This is designed for use with Claude Code:

```bash
atrelease release v1.0.0 --interactive --json
```

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Release completed successfully |
| 1 | Release failed at some step |

## Best Practices

1. **Always use dry-run first**: Preview what will happen before executing
2. **Don't skip CI**: Let CI verify the release before tagging
3. **Use semantic versioning**: Follow semver for version numbers
4. **Keep working directory clean**: Commit or stash changes before releasing
