# agent-team-release - Kiro CLI Plugin

Multi-agent team for automating software release workflows including versioning, changelog generation, CI verification, security scanning, and Git tagging

## Agents

| Agent | Description |
|-------|-------------|
| `code-fixer` | Fix code quality issues reported by QA validation |
| `docs-reviewer` | Documentation validation for release readiness |
| `docs-writer` | Create and update release documentation |
| `pm` | Product Management specialist for release scoping and version decisions |
| `qa` | Quality Assurance validation for release readiness |
| `release-coordinator` | Orchestrates software releases including CI verification and Git tagging |
| `release` | Release Management validation for deployment readiness |
| `security` | Security and Compliance validation for release readiness |

## Usage

Run an agent with the Kiro CLI:

```bash
kiro-cli chat --agent code-fixer "<your prompt>"
```

Run the full team (coordinator-driven):

```bash
kiro-cli chat --agent release-coordinator "<your prompt>"
```

## Steering Files

Copy steering files to `.kiro/steering/` for automatic context loading:

```bash
mkdir -p .kiro/steering
cp steering/*.md .kiro/steering/
```

## Installation

Copy agents to your Kiro agents directory:

```bash
mkdir -p ~/.kiro/agents
cp agents/*.json ~/.kiro/agents/
```
