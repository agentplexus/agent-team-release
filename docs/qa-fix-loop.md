# QA Fix Loop

The release-coordinator orchestrates an automated **build/test/lint/fix/retest loop** to ensure code quality before release.

## Overview

When QA validation fails, the release-coordinator automatically invokes the code-fixer agent to resolve issues, then re-runs QA validation. This loop continues until all checks pass or maximum attempts are reached.

## Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           QA FIX LOOP                                       │
│                    (Orchestrated by release-coordinator)                    │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│                         ┌───────────────────┐                               │
│                         │ release-coordinator│                               │
│                         │   (orchestrator)  │                               │
│                         └─────────┬─────────┘                               │
│                                   │                                         │
│                    ┌──────────────┼──────────────┐                         │
│                    │              │              │                         │
│                    ▼              ▼              ▼                         │
│              ┌──────────┐  ┌──────────┐  ┌──────────┐                      │
│              │ QA Agent │  │Code-Fixer│  │ QA Agent │                      │
│              │(validate)│  │  Agent   │  │(re-test) │                      │
│              └────┬─────┘  └────┬─────┘  └────┬─────┘                      │
│                   │             │             │                            │
│                   ▼             ▼             ▼                            │
│              ┌─────────┐  ┌─────────┐  ┌─────────┐                         │
│              │• build  │  │• errcheck│  │• build  │                        │
│              │• test   │  │• gosec   │  │• test   │                        │
│              │• lint   │  │• format  │  │• lint   │                        │
│              │• format │  │• unused  │  │• format │                        │
│              │• mod    │  │• mod tidy│  │• mod    │                        │
│              └─────────┘  └─────────┘  └─────────┘                         │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

## Flow Diagram

```
                              START
                                │
                                ▼
                    ┌───────────────────────┐
                    │   Pre-flight Checks   │
                    │  (clean working dir)  │
                    └───────────┬───────────┘
                                │
                                ▼
              ┌─────────────────────────────────────┐
              │         INVOKE QA AGENT             │
              │  • go build ./...                   │
              │  • go test -v ./...                 │
              │  • golangci-lint run                │
              │  • gofmt -l .                       │
              │  • go mod tidy -diff                │
              └─────────────────┬───────────────────┘
                                │
                                ▼
                         ┌──────────────┐
                         │  QA Status?  │
                         └──────┬───────┘
                                │
                 ┌──────────────┼──────────────┐
                 │              │              │
                 ▼              │              ▼
            ┌────────┐         │         ┌────────┐
            │   GO   │         │         │ NO-GO  │
            └────┬───┘         │         └────┬───┘
                 │             │              │
                 │             │              ▼
                 │             │    ┌─────────────────┐
                 │             │    │ attempts < max? │
                 │             │    └────────┬────────┘
                 │             │             │
                 │             │      ┌──────┴──────┐
                 │             │      │             │
                 │             │      ▼             ▼
                 │             │    ┌────┐      ┌──────┐
                 │             │    │YES │      │  NO  │
                 │             │    └──┬─┘      └──┬───┘
                 │             │       │           │
                 │             │       ▼           ▼
                 │             │  ┌─────────────────────────┐
                 │             │  │  INVOKE CODE-FIXER      │
                 │             │  │  • Fix errcheck errors  │
                 │             │  │  • Fix gosec issues     │
                 │             │  │  • Fix format issues    │
                 │             │  │  • Remove unused code   │
                 │             │  │  • Run go mod tidy      │
                 │             │  └───────────┬─────────────┘
                 │             │              │
                 │             │              │ loop back
                 │             │              │
                 │             └──────────────┘
                 │
                 ▼                             ▼
    ┌────────────────────────┐    ┌────────────────────────┐
    │  PROCEED TO RELEASE    │    │   RAISE TO HUMAN       │
    │  • Version             │    │   • Report unfixed     │
    │  • Changelog           │    │   • Suggest manual fix │
    │  • Release notes       │    │   • List attempts      │
    │  • Documentation       │    └────────────────────────┘
    │  • Git tag             │
    └────────────────────────┘
```

## Agent Responsibilities

### QA Agent (Validation)

| Check | Command | Required |
|-------|---------|----------|
| build | `go build ./...` | Yes |
| tests | `go test -v ./...` | Yes |
| lint | `golangci-lint run` | Yes |
| format | `gofmt -l .` | Yes |
| mod-tidy | `go mod tidy -diff` | Yes |
| error-handling | Pattern: `_ = err` | Yes |
| local-replace | Pattern: `replace .* => ./` | Yes |

The QA agent is **read-only** - it validates and reports but does not modify code.

### Code-Fixer Agent (Remediation)

The code-fixer agent has **write access** and fixes issues following the error handling priority:

| Priority | Approach | When to Use |
|----------|----------|-------------|
| 1 | Panic | Error should never happen (invariant violation) |
| 2 | Return Error | Function signature can return error |
| 3 | Modify Function | Function can be changed to return error |
| 4 | Log Error | Interface compliance requires no error return |
| 5 | Raise to Human | Cannot be auto-fixed |

Common fixes applied:

| Issue | Fix |
|-------|-----|
| errcheck | Add error handling |
| G306 | Change permissions to 0o600 |
| G115 | Add bounds check or nolint |
| unused | Remove unused code |
| format | Run gofmt |

### Release-Coordinator (Orchestration)

The release-coordinator invokes subagents using the Task tool:

```go
// Step 1: Validate
Task(subagent_type="qa", prompt="Run QA validation...")

// Step 2: Fix (if NO-GO)
Task(subagent_type="code-fixer", prompt="Fix issues: <findings>")

// Step 3: Re-validate
Task(subagent_type="qa", prompt="Re-validate after fixes...")
```

## Configuration

### Max Attempts

Default: **3 fix attempts** before raising to human.

Some issues cannot be auto-fixed:

- Architectural problems
- Missing test coverage
- Complex refactoring required
- Security issues requiring human judgment

### Supported Languages

Currently supports Go projects:

- `golangci-lint` for linting
- `gofmt` for formatting
- `go test` for testing
- `go build` for compilation
- `go mod tidy` for dependencies

## Usage

### Via Release-Coordinator

```bash
# Full release workflow (includes QA fix loop)
Task(subagent_type="release-coordinator",
     prompt="Run release workflow on /path/to/project")

# Just QA fix loop
Task(subagent_type="release-coordinator",
     prompt="Run QA fix loop on /path/to/project")
```

### Manual Steps

```bash
# Step 1: Run QA validation
golangci-lint run
go test -v ./...
gofmt -l .

# Step 2: Fix issues (or invoke code-fixer)
# ... apply fixes ...

# Step 3: Re-validate
golangci-lint run
go test -v ./...
```

## Example Output

```
╔════════════════════════════════════════════════════════════════════════════╗
║                           QA FIX LOOP                                      ║
╠════════════════════════════════════════════════════════════════════════════╣
║ Project: github.com/plexusone/omnivoice-core                              ║
╠════════════════════════════════════════════════════════════════════════════╣
║ ATTEMPT 1                                                                  ║
╠════════════════════════════════════════════════════════════════════════════╣
║ QA Validation:     🔴 NO-GO (17 issues)                                   ║
║   errcheck         9 issues                                                ║
║   gosec G306       3 issues                                                ║
║   gosec G115       4 issues                                                ║
║   unused           1 issue                                                 ║
╠════════════════════════════════════════════════════════════════════════════╣
║ Code-Fixer:        ✓ 17 fixes applied                                     ║
╠════════════════════════════════════════════════════════════════════════════╣
║ QA Re-validation:  🟢 GO                                                  ║
╠════════════════════════════════════════════════════════════════════════════╣
║                         ✓ QA FIX LOOP COMPLETE                            ║
║                           Proceeding to release                            ║
╚════════════════════════════════════════════════════════════════════════════╝
```

## See Also

- [Orchestration](orchestration.md) - Multi-agent coordination
- [Configuration](configuration.md) - Tool configuration
- [Getting Started](getting-started/index.md) - Quick start guide
