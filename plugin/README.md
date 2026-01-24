# Release Agent Plugin for Claude Code

This plugin integrates Release Agent with Claude Code, enabling automated software release workflows directly from your AI coding assistant.

## Features

- **Automated Releases**: Full release workflow with version validation, changelog generation, and Git tagging
- **Validation Checks**: Build, test, lint, and format verification for Go and TypeScript projects
- **Semantic Versioning**: Automatic version suggestion based on conventional commits
- **CI Integration**: Wait for GitHub Actions CI to pass before tagging
- **Interactive Mode**: Step-by-step approval for each release action

## Installation

### Local Testing (Development)

```bash
# From your project directory
claude --plugin-dir /path/to/release-agent/plugin
```

### Install from Marketplace (Coming Soon)

```bash
claude plugin install release-agent --scope user
```

## Commands

| Command | Description |
|---------|-------------|
| `/release-agent:release <version>` | Execute full release workflow |
| `/release-agent:check` | Run validation checks |
| `/release-agent:changelog <version>` | Generate changelog |
| `/release-agent:version-next` | Suggest next version |

## Agents

### release-coordinator

The main orchestration agent for releases. Invoke by asking Claude to "create a release" or "prepare version X.Y.Z".

## Skills

### version-analysis

Analyzes git history to determine semantic version bumps.

### commit-classification

Classifies commits according to conventional commits specification.

## Requirements

### Required Dependencies

- `git` - Version control
- `gh` - GitHub CLI (for CI status)
- `atrelease` - Release automation CLI
- `schangelog` - Changelog generation
- `sroadmap` - Roadmap generation

### Language-Specific (Auto-detected)

**Go Projects:**
- `go` - Go compiler
- `golangci-lint` - Linter

**Node.js Projects:**
- `node` / `npm` - Runtime and package manager
- `eslint` / `prettier` - Linting and formatting

## Configuration

The plugin uses `.releaseagent.yaml` in your project root for customization:

```yaml
# .releaseagent.yaml
checks:
  go:
    enabled: true
    lint: true
    test: true
  typescript:
    enabled: true
    lint: true
    test: true
```

## Usage Examples

### Create a Release

```
User: Create a release for v1.0.0

Claude: I'll help you create a release for v1.0.0. Let me first run the validation checks...
[Runs atrelease check]
All checks passed! Now executing the release workflow...
[Runs atrelease release v1.0.0]
```

### Check Next Version

```
User: What should the next version be?

Claude: Let me analyze the commits since the last release...
[Uses version-analysis skill]
Based on 3 new features and 2 bug fixes, I suggest v0.9.0.
```

### Dry Run

```
User: Preview the release for v2.0.0

Claude: I'll do a dry run to preview the release...
[Runs atrelease release v2.0.0 --dry-run]
```

## License

Apache-2.0
