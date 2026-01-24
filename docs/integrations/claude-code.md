# Claude Code Integration

Release Agent integrates with Claude Code as a plugin, providing commands, skills, and an autonomous release agent.

## Installation

The Release Agent plugin is available in the Claude Code marketplace:

```
/install release-agent
```

Or manually add to your Claude Code configuration.

## Available Commands

### `/release-agent:release`

Execute a full release workflow:

```
/release-agent:release v1.0.0
```

### `/release-agent:check`

Run validation checks:

```
/release-agent:check
```

### `/release-agent:changelog`

Generate changelog:

```
/release-agent:changelog --since=v0.9.0
```

### `/release-agent:version-next`

Suggest next version based on commits:

```
/release-agent:version-next
```

## Interactive Mode

When running with `--interactive`, Release Agent can:

- Ask questions when issues arise
- Propose fixes for problems
- Get approval before making changes

### Example Interaction

```
Claude: I'll run the release checks for v1.0.0.

[Running atrelease check --interactive --json]

Release Agent: Found 2 lint issues. Would you like me to:
1. Show the issues and let you fix them
2. Propose automatic fixes
3. Skip lint checks for this release

User: 2

Release Agent: Proposed fixes:
- main.go:42: Add error handling for Close()
- utils.go:18: Remove unused variable

Apply these fixes? [y/n]
```

## Skills

### Version Analysis

Analyzes commits to suggest appropriate version bumps:

- `feat:` commits → minor version
- `fix:` commits → patch version
- `BREAKING CHANGE:` → major version

### Commit Classification

Classifies commits by conventional commit type for changelog generation.

## Release Coordinator Agent

The release coordinator agent can orchestrate complete releases autonomously:

1. Analyzes current state
2. Determines appropriate version
3. Runs all checks
4. Generates documentation
5. Creates release commit
6. Manages CI verification
7. Creates release tag

### Invoking the Agent

```
@release-agent Please prepare the v1.0.0 release
```

## Output Formats for Claude

### TOON Format

For token-efficient communication, use TOON format:

```bash
atrelease check --json --format=toon
```

TOON is approximately 8x more token-efficient than JSON.

### JSON Protocol

For structured communication:

```bash
atrelease check --interactive --json
```

Questions and proposals are returned as structured JSON:

```json
{
  "type": "question",
  "id": "lint-fix-proposal",
  "message": "Found 2 lint issues. How should I proceed?",
  "options": [
    {"id": "show", "label": "Show issues"},
    {"id": "fix", "label": "Auto-fix"},
    {"id": "skip", "label": "Skip"}
  ]
}
```

## Hooks

### SessionStart Hook

The plugin includes a SessionStart hook that checks for required dependencies:

- `git` - Version control
- `gh` - GitHub CLI
- `schangelog` - Changelog generation (optional)
- `golangci-lint` - Go linting (optional)

Missing optional dependencies trigger warnings, not errors.

## Configuration

### Plugin Configuration

```yaml
# .claude/plugins/release-agent.yaml
enabled: true
auto_fix: false  # Require approval for fixes
verbose: false
```

### Environment Variables

| Variable | Description |
|----------|-------------|
| `RELEASEAGENT_INTERACTIVE` | Force interactive mode |
| `RELEASEAGENT_FORMAT` | Default output format |

## Gemini CLI Support

Release Agent also supports Gemini CLI through platform adapters. The same plugin structure works for both platforms.

## Troubleshooting

### Plugin Not Found

Ensure Release Agent is installed:

```bash
which agent-team-release
```

### Permission Errors

The plugin needs permission to:

- Read files
- Execute shell commands
- Write to files (for fixes)

Grant these permissions when prompted.

### CI Waiting Timeout

If CI takes longer than expected:

```bash
atrelease release v1.0.0 --skip-ci
```

Note: This tags before CI passes, which is not recommended.
