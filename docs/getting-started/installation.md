# Installation

## Go Install

The easiest way to install Release Agent Team is using `go install`:

```bash
go install github.com/agentplexus/agent-team-release/cmd/atrelease@latest
```

This installs the `atrelease` binary to your `$GOPATH/bin` directory.

## Homebrew

On macOS and Linux, you can install via Homebrew:

```bash
brew tap agentplexus/tap
brew install atrelease
```

## From Source

Clone the repository and build:

```bash
git clone https://github.com/agentplexus/agent-team-release.git
cd agent-team-release
go build -o atrelease ./cmd/atrelease
```

## Verify Installation

Check that Release Agent is installed correctly:

```bash
atrelease version
```

## Dependencies

### Required

| Tool | Purpose |
|------|---------|
| `git` | Version control operations |
| `gh` | GitHub CLI for CI status checking |

### Language-Specific

| Tool | Language | Purpose |
|------|----------|---------|
| `go` | Go | Build and test |
| `golangci-lint` | Go | Linting |
| `node`, `npm` | TypeScript/JS | Build and test |
| `eslint` | TypeScript/JS | Linting |
| `prettier` | TypeScript/JS | Formatting |

### Optional

| Tool | Purpose |
|------|---------|
| `schangelog` | Changelog generation |
| `sroadmap` | Roadmap management |
| `gocoverbadge` | Coverage badge generation |
| `govulncheck` | Vulnerability scanning |

## Installing Optional Tools

### schangelog

For automated changelog generation:

```bash
go install github.com/grokify/structured-changelog/cmd/schangelog@latest
```

### sroadmap

For roadmap updates:

```bash
go install github.com/grokify/structured-roadmap/cmd/sroadmap@latest
```

### golangci-lint

For Go linting:

```bash
# macOS
brew install golangci-lint

# Linux/Windows
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### govulncheck

For vulnerability scanning:

```bash
go install golang.org/x/vuln/cmd/govulncheck@latest
```
