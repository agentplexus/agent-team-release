# Release Agent - Gemini CLI Extension

Gemini CLI extension for automated software releases.

## Installation

### Method 1: Install from GitHub

```bash
gemini extensions install https://github.com/grokify/release-agent --ref=main
```

The extension files are in the `plugins/gemini/` directory.

### Method 2: Link for Development

```bash
cd release-agent/plugins/gemini
gemini extensions link .
```

## Usage

After installation, use the release commands:

```bash
# Start Gemini CLI
gemini

# Use release commands
> /release-agent:release v1.2.3
> /release-agent:check
> /release-agent:changelog v1.2.3
> /release-agent:version-next
```

## Available Commands

| Command | Description |
|---------|-------------|
| `/release-agent:release <version>` | Execute full release workflow |
| `/release-agent:check` | Run validation checks (build, test, lint) |
| `/release-agent:changelog <version>` | Generate changelog for version |
| `/release-agent:version-next` | Suggest next semantic version |

## Extension Structure

```
plugins/gemini/
├── gemini-extension.json    # Extension manifest
├── GEMINI.md                # Context for the model
├── README.md                # This file
└── commands/                # TOML command definitions
    ├── release.toml
    ├── check.toml
    ├── changelog.toml
    └── version-next.toml
```

## Command Details

### /release-agent:release

Executes the complete release workflow:

1. Validate version format
2. Check working directory is clean
3. Run validation checks (build, test, lint, format)
4. Generate changelog via schangelog
5. Update roadmap via sroadmap
6. Create release commit
7. Push to remote
8. Wait for CI to pass
9. Create and push release tag

**Usage:**

```
/release-agent:release v1.0.0
```

**Options:**

```bash
# Dry run (preview without changes)
atrelease release v1.0.0 --dry-run

# Interactive mode (approve each step)
atrelease release v1.0.0 --interactive
```

### /release-agent:check

Runs all validation checks without creating a release:

- Build verification
- Test suite execution
- Linting (golangci-lint for Go, ESLint for TypeScript)
- Format checking (gofmt, Prettier)

### /release-agent:changelog

Generates changelog entries for a specific version using schangelog.

### /release-agent:version-next

Analyzes commits since the last tag and suggests the next semantic version based on conventional commits.

## Dependencies

The extension expects these CLI tools to be installed:

| Tool | Purpose |
|------|---------|
| `atrelease` | Release automation CLI |
| `schangelog` | Structured changelog generator |
| `sroadmap` | Roadmap management |
| `golangci-lint` | Go linter (for Go projects) |
| `git` | Version control |
| `gh` | GitHub CLI |

## Supported Languages

| Language | Build | Test | Lint | Format |
|----------|-------|------|------|--------|
| Go | `go build` | `go test` | `golangci-lint` | `gofmt` |
| TypeScript | `tsc` | `npm test` | `ESLint` | `Prettier` |

## Extension Manifest

The `gemini-extension.json` defines the extension:

```json
{
  "name": "release-agent",
  "version": "1.0.0",
  "description": "Automates software release workflows",
  "contextFileName": "GEMINI.md"
}
```

## Updating

```bash
gemini extensions update release-agent
```

## Uninstalling

```bash
gemini extensions uninstall release-agent
```

## Comparison with Other Platforms

| Feature | Gemini CLI | Claude Code | Kiro CLI |
|---------|------------|-------------|----------|
| Commands | TOML files | Markdown files | N/A |
| Context | GEMINI.md | CLAUDE.md | Steering files |
| Sub-agents | No | Yes | Yes |
| Format | JSON manifest | JSON manifest | JSON agents |

## Sources

- [Gemini CLI Extensions](https://github.com/google-gemini/gemini-cli/blob/main/docs/extensions/index.md)
- [Extensions Gallery](https://geminicli.com/extensions/browse/)
