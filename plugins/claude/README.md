# Release Agent Team - Claude Code Plugin

A multi-agent team plugin for Claude Code that automates software release workflows with specialized agents for PM, QA, Documentation, Security, and Release coordination.

## Installation

### From GitHub (Recommended)

```bash
claude plugin add github:agentplexus/agent-team-release/plugins/claude
```

### From Local Path

```bash
claude plugin add /path/to/agent-team-release/plugins/claude
```

### Manual Installation

Download and copy to your project:

```bash
curl -sL https://github.com/agentplexus/agent-team-release/archive/main.tar.gz | tar -xz
cp -r agent-team-release-main/plugins/claude your-project/plugins/
```

## Team Agents

This plugin includes 6 specialized agents that work together:

| Agent | Role | Model |
|-------|------|-------|
| `release-coordinator` | Orchestrates release workflow | opus |
| `pm` | Product scope and versioning | sonnet |
| `qa` | Build, test, lint validation | sonnet |
| `documentation` | Docs completeness checks | haiku |
| `release` | Git state and CI validation | sonnet |
| `security` | Vulnerability and compliance | sonnet |

### Workflow DAG

```
pm-validation
    ├──> qa-validation ──────────────┐
    ├──> docs-validation ────────────┤
    ├──> security-validation ────────┼──> execute-release
    └──> release-validation ─────────┘
         (depends on pm + qa)
```

## Usage

After installation, use the release commands:

```bash
# Start Claude Code
claude

# Use release commands
> /agent-team-release:release v1.2.3
> /agent-team-release:check
> /agent-team-release:changelog v1.2.3
> /agent-team-release:version-next
```

Or invoke agents directly:

```
Use the release-coordinator agent to prepare a release
Use the qa agent to validate code quality
Use the security agent to check for vulnerabilities
```

## Available Commands

| Command | Description |
|---------|-------------|
| `/agent-team-release:release <version>` | Execute full release workflow |
| `/agent-team-release:check` | Run validation checks (build, test, lint) |
| `/agent-team-release:changelog <version>` | Generate changelog for version |
| `/agent-team-release:version-next` | Suggest next semantic version |

## Plugin Structure

```
plugins/claude/
├── .claude-plugin/
│   └── plugin.json              # Plugin manifest
├── README.md                    # This file
├── agents/                      # Agent definitions (6 agents)
│   ├── release-coordinator.md   # Orchestrator
│   ├── pm.md                    # Product Manager
│   ├── qa.md                    # Quality Assurance
│   ├── documentation.md         # Documentation
│   ├── release.md               # Release mechanics
│   └── security.md              # Security & compliance
├── commands/                    # Markdown command definitions
│   ├── release.md
│   ├── check.md
│   ├── changelog.md
│   └── version-next.md
└── skills/                      # Reusable skill definitions
    ├── version-analysis/
    │   └── SKILL.md
    └── commit-classification/
        └── SKILL.md
```

## Agent Responsibilities

### release-coordinator (Orchestrator)

Coordinates all agents and executes the final release:

- Delegates validation tasks to specialized agents
- Collects and synthesizes validation results
- Makes go/no-go release decisions
- Executes final release steps (tag, push)

### pm (Product Manager)

Validates release from a product perspective:

- Version recommendation based on commits
- Release scope verification
- Changelog quality assessment
- Breaking changes identification
- Roadmap alignment

### qa (Quality Assurance)

Technical quality validation:

- Build verification
- Test execution
- Lint checks
- Format verification
- Error handling audit

### documentation

Documentation completeness:

- README validation
- CHANGELOG.md/json verification
- Release notes check
- PRD/TRD existence

### release

Release mechanics validation:

- Version tag availability
- Git clean state
- Remote configuration
- CI configuration
- Changelog finalization

### security

Security and compliance:

- License compliance
- Vulnerability scanning (govulncheck)
- Dependency audit
- Secret detection
- Environment file checks

## Dependencies

The plugin expects these CLI tools to be installed:

| Tool | Purpose |
|------|---------|
| `atrelease` | Release automation CLI |
| `schangelog` | Structured changelog generator |
| `sroadmap` | Roadmap management |
| `golangci-lint` | Go linter (for Go projects) |
| `govulncheck` | Go vulnerability scanner |
| `git` | Version control |
| `gh` | GitHub CLI |

## Supported Languages

| Language | Build | Test | Lint | Format |
|----------|-------|------|------|--------|
| Go | `go build` | `go test` | `golangci-lint` | `gofmt` |
| TypeScript | `tsc` | `npm test` | `ESLint` | `Prettier` |

## Usage Examples

### Quick Release

```bash
# Check project is ready
/agent-team-release:check

# Get suggested version
/agent-team-release:version-next

# Execute release
/agent-team-release:release v1.2.0
```

### Full Team Review

Ask the release-coordinator to orchestrate a complete team review:

```
Use the release-coordinator agent to prepare v2.0.0 with full validation from all team members
```

### Individual Agent Tasks

```
Use the security agent to scan for vulnerabilities
Use the qa agent to run all quality checks
Use the pm agent to recommend the next version
```

## Multi-Agent Spec

This plugin is generated from [multi-agent-spec](https://github.com/agentplexus/multi-agent-spec) canonical definitions. The same agent team can be deployed to:

- Claude Code (this plugin)
- Kiro CLI
- AWS Bedrock AgentCore
- AgentKit local servers

## License

Apache-2.0
