# Release Agent - Technical Requirements Document

## Architecture Overview

Release Agent provides a modular architecture for release automation with validation checks, mutating actions, workflow orchestration, and multiple output formats for both CLI and AI assistant integration.

```
┌─────────────────────────────────────────────────────────────┐
│                     CLI Layer (Cobra)                       │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────────────┐│
│  │  check   │ │ validate │ │ release  │ │ changelog/readme ││
│  └────┬─────┘ └────┬─────┘ └────┬─────┘ └────────┬─────────┘│
│       └────────────┴────────────┴────────────────┘          │
│                          │                                   │
├──────────────────────────┼───────────────────────────────────┤
│                    Core Layer (pkg/)                         │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐       │
│  │    checks    │  │   actions    │  │   workflow   │       │
│  │  - golang    │  │  - changelog │  │  - release   │       │
│  │  - typescript│  │  - roadmap   │  │  - runner    │       │
│  │  - security  │  │  - readme    │  │  - steps     │       │
│  │  - docs      │  │              │  │              │       │
│  │  - release   │  │              │  │              │       │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘       │
│         │                 │                 │                │
│  ┌──────┴─────────────────┴─────────────────┴──────┐        │
│  │              Supporting Packages                 │        │
│  │  ┌────────┐ ┌────────────┐ ┌────────┐ ┌───────┐ │        │
│  │  │ config │ │interactive │ │ output │ │  git  │ │        │
│  │  └────────┘ └────────────┘ └────────┘ └───────┘ │        │
│  └──────────────────────────────────────────────────┘        │
│                          │                                   │
├──────────────────────────┼───────────────────────────────────┤
│                  Integration Layer                           │
│  ┌────────┐  ┌───────────┐  ┌─────────┐  ┌────────┐         │
│  │  Git   │  │schangelog │  │sroadmap │  │ gh CLI │         │
│  └────────┘  └───────────┘  └─────────┘  └────────┘         │
└─────────────────────────────────────────────────────────────┘
```

## Module Structure

```
github.com/grokify/release-agent/
├── cmd/
│   └── releaseagent/           # CLI entry point
│       ├── main.go             # Main entry
│       ├── root.go             # Root command with global flags
│       ├── check.go            # check subcommand
│       ├── validate.go         # validate subcommand
│       ├── release.go          # release subcommand
│       ├── changelog.go        # changelog subcommand
│       ├── readme.go           # readme subcommand
│       ├── roadmap.go          # roadmap subcommand
│       └── version.go          # version subcommand
├── pkg/
│   ├── checks/                 # Validation checks
│   │   ├── checks.go           # Result, Checker interface, Options
│   │   ├── golang.go           # GoChecker implementation
│   │   ├── typescript.go       # TypeScriptChecker implementation
│   │   ├── security.go         # SecurityChecker implementation
│   │   ├── documentation.go    # DocChecker implementation
│   │   ├── release.go          # ReleaseChecker implementation
│   │   └── areas.go            # Validation areas (QA, Docs, Release, Security)
│   ├── actions/                # Mutating actions
│   │   ├── actions.go          # Action interface, Result, Proposal
│   │   ├── changelog.go        # ChangelogAction (schangelog integration)
│   │   ├── roadmap.go          # RoadmapAction (sroadmap integration)
│   │   └── readme.go           # ReadmeAction (badge/version updates)
│   ├── workflow/               # Release workflow orchestration
│   │   ├── workflow.go         # Workflow, Step, Runner, Context
│   │   └── release.go          # ReleaseWorkflow (9-step process)
│   ├── detect/                 # Language detection
│   │   └── detect.go           # Language enum, Detect function
│   ├── config/                 # Configuration
│   │   └── config.go           # Config, LanguageConfig, Load
│   ├── interactive/            # User interaction
│   │   ├── interactive.go      # Prompter interface, Question, Answer
│   │   ├── cli.go              # CLIPrompter (terminal)
│   │   └── json.go             # JSONPrompter (Claude Code)
│   ├── git/                    # Git operations
│   │   ├── git.go              # Git wrapper, Status, tag/commit ops
│   │   └── ci.go               # CI status checking via gh CLI
│   └── output/                 # Output formatting
│       ├── json.go             # JSON output writer
│       └── toon.go             # TOON output writer
├── plugin/                     # Claude Code plugin (development)
│   ├── commands/               # Command definitions
│   ├── skills/                 # Skill definitions
│   ├── agents/                 # Agent definitions
│   └── hooks/                  # Lifecycle hooks
│       ├── hooks.json          # Hook configuration
│       └── scripts/            # Hook scripts
│           ├── check-dependencies.sh
│           └── validate-release-command.sh
├── plugins/                    # Multi-platform plugins
│   ├── spec/                   # Canonical JSON specifications
│   │   ├── plugin.json         # Plugin metadata
│   │   ├── commands/           # Command specs (JSON)
│   │   ├── skills/             # Skill specs (JSON)
│   │   └── agents/             # Agent specs (JSON)
│   ├── generate/               # Plugin generator
│   │   └── main.go             # Generates Claude/Gemini plugins
│   ├── claude/                 # Generated Claude plugin
│   └── gemini/                 # Generated Gemini plugin
├── .goreleaser.yaml            # Multi-platform builds
├── Formula/                    # Homebrew formula
├── go.mod
├── go.sum
├── README.md
├── CHANGELOG.md
├── PRD.md
├── TRD.md
└── ROADMAP.md
```

## Core Interfaces

### Checker Interface

```go
// pkg/checks/checks.go

// Result represents the outcome of a check.
type Result struct {
    Name    string  // Check name
    Passed  bool    // Whether check passed
    Output  string  // Check output/details
    Error   error   // Error if check failed
    Skipped bool    // Whether check was skipped
    Reason  string  // Reason for skip
    Warning bool    // Soft check (reported but doesn't fail)
}

// Checker is the interface for validation checks.
type Checker interface {
    Name() string
    Check(dir string, opts Options) []Result
}

// Options configures which checks to run.
type Options struct {
    Test              bool
    Lint              bool
    Format            bool
    Coverage          bool
    Verbose           bool
    Interactive       bool
    GoExcludeCoverage string
}

// ValidationArea represents a domain of validation.
type ValidationArea int

const (
    AreaQA ValidationArea = iota
    AreaDocumentation
    AreaRelease
    AreaSecurity
)

// AreaStatus represents Go/No-Go status for an area.
type AreaStatus int

const (
    StatusGo AreaStatus = iota
    StatusNoGo
    StatusWarn
    StatusSkip
)
```

### Action Interface

```go
// pkg/actions/actions.go

// Result represents the outcome of an action.
type Result struct {
    Name    string
    Success bool
    Output  string
    Error   error
    Skipped bool
    Reason  string
}

// Proposal represents a proposed change for user approval.
type Proposal struct {
    Description string            // Human-readable description
    FilePath    string            // File to modify
    OldContent  string            // Current content (for diff)
    NewContent  string            // Proposed new content
    Metadata    map[string]string // Additional context
}

// Action is the interface for mutating operations.
type Action interface {
    Name() string

    // Run executes the action directly (non-interactive)
    Run(dir string, opts Options) Result

    // Propose generates proposals for interactive mode
    Propose(dir string, opts Options) ([]Proposal, error)

    // Apply applies approved proposals
    Apply(dir string, proposals []Proposal) Result
}

// Options configures action behavior.
type Options struct {
    DryRun      bool           // Don't make changes
    Interactive bool           // Enable interactive mode
    Version     string         // Target version
    Since       string         // Since tag (for changelog)
    Verbose     bool           // Verbose output
    Config      *config.Config // Configuration
}
```

### Workflow Interface

```go
// pkg/workflow/workflow.go

// StepType defines the type of workflow step.
type StepType int

const (
    StepTypeFunc StepType = iota      // Executes a function
    StepTypeComposite                  // Contains sub-steps
)

// Step represents a single step in a workflow.
type Step struct {
    Name        string
    Description string
    Type        StepType
    Required    bool              // Workflow fails if step fails
    Func        StepFunc          // Function to execute
    SubSteps    []Step            // For composite steps
}

// StepFunc is the function signature for step execution.
type StepFunc func(ctx *Context) *StepResult

// StepResult represents the result of a step.
type StepResult struct {
    Name    string
    Success bool
    Output  string
    Error   error
    Skipped bool
    Reason  string
}

// Workflow defines a sequence of steps.
type Workflow struct {
    Name        string
    Description string
    Steps       []Step
}

// Context provides execution context for workflow steps.
type Context struct {
    Dir         string
    Version     string
    DryRun      bool
    Verbose     bool
    Interactive bool
    JSONOutput  bool
    SkipChecks  bool
    SkipCI      bool
    Data        map[string]interface{}  // Inter-step data
    Output      *strings.Builder        // Accumulated output
}

// Runner executes workflows.
type Runner struct {
    DryRun      bool
    Verbose     bool
    Interactive bool
    JSONOutput  bool
}

// WorkflowResult represents the complete workflow result.
type WorkflowResult struct {
    WorkflowName string
    Success      bool
    Steps        []StepResult
    Duration     time.Duration
    Output       string
}

// Run executes a workflow.
func (r *Runner) Run(w *Workflow, ctx *Context) (*WorkflowResult, error)
```

### Interactive Interface

```go
// pkg/interactive/interactive.go

// QuestionType defines the type of question.
type QuestionType int

const (
    QuestionTypeSingleChoice QuestionType = iota
    QuestionTypeMultiChoice
    QuestionTypeConfirm
    QuestionTypeText
)

// Question represents a question for the user.
type Question struct {
    ID      string
    Text    string
    Type    QuestionType
    Options []Option
    Default string
    Context string  // Additional context (e.g., code snippet)
}

// Option represents a choice option.
type Option struct {
    ID          string
    Label       string
    Description string
}

// Answer represents a user's response.
type Answer struct {
    QuestionID string
    Selected   []string  // Selected option IDs
    Text       string    // For text questions
    Confirmed  bool      // For confirm questions
}

// ProposalAction represents user's decision on a proposal.
type ProposalAction int

const (
    ProposalActionApply ProposalAction = iota
    ProposalActionSkip
    ProposalActionEdit
    ProposalActionAbort
)

// Prompter handles user interaction.
type Prompter interface {
    Ask(q Question) (Answer, error)
    ShowProposal(p actions.Proposal) error
    Confirm(message string) (bool, error)
    Info(message string)
    Warn(message string)
    Error(message string)
}

// CLIPrompter implements Prompter for terminal.
type CLIPrompter struct {
    reader *bufio.Reader
}

// JSONPrompter implements Prompter for Claude Code (outputs JSON).
type JSONPrompter struct {
    Output io.Writer
    Input  io.Reader
}
```

### Git Interface

```go
// pkg/git/git.go

// Git provides git operations.
type Git struct {
    Dir    string
    Remote string
}

// Status represents repository status.
type Status struct {
    Branch      string
    Ahead       int
    Behind      int
    Staged      []string
    Modified    []string
    Untracked   []string
    HasRemote   bool
    RemoteBranch string
}

// Tag operations
func (g *Git) LatestTag() (string, error)
func (g *Git) AllTags() ([]string, error)
func (g *Git) CreateTag(tag, message string) error
func (g *Git) DeleteTag(tag string) error
func (g *Git) PushTag(tag string) error

// Repository operations
func (g *Git) Status() (*Status, error)
func (g *Git) Push() error
func (g *Git) PushWithUpstream(branch string) error
func (g *Git) CommitAll(message string) error
func (g *Git) Commit(message string, files ...string) error

// Information
func (g *Git) CurrentBranch() (string, error)
func (g *Git) CurrentCommit() (string, error)
func (g *Git) ShortCommit() (string, error)
func (g *Git) RemoteURL() (string, error)
func (g *Git) IsAncestor(ancestor, descendant string) (bool, error)

// pkg/git/ci.go

// CIStatus represents CI check status.
type CIStatus struct {
    State       string       // success, failure, pending
    TotalCount  int
    Statuses    []StatusCheck
    CheckSuites []CheckSuite
}

// CI operations
func (g *Git) GetCIStatus(ref string) (*CIStatus, error)
func (g *Git) WaitForCI(timeout time.Duration) error
func (g *Git) IsCIPassing(ref string) (bool, error)
func (g *Git) GetPRForBranch() (int, error)
func (g *Git) GetPRStatus(prNumber int) (*CIStatus, error)
```

### Output Interface

```go
// pkg/output/json.go and pkg/output/toon.go

// MessageType identifies the type of output message.
type MessageType string

const (
    MessageTypeQuestion MessageType = "question"
    MessageTypeProposal MessageType = "proposal"
    MessageTypeInfo     MessageType = "info"
    MessageTypeWarning  MessageType = "warning"
    MessageTypeError    MessageType = "error"
    MessageTypeResult   MessageType = "result"
    MessageTypeProgress MessageType = "progress"
)

// Writer outputs structured messages.
type Writer interface {
    WriteQuestion(q interactive.Question) error
    WriteProposal(p actions.Proposal) error
    WriteInfo(message string) error
    WriteWarning(message string) error
    WriteError(message string) error
    WriteResult(r actions.Result) error
    WriteProgress(step, total int, name, status string) error
}

// JSONWriter outputs JSON format.
type JSONWriter struct {
    w io.Writer
}

// TOONWriter outputs TOON format (token-optimized).
type TOONWriter struct {
    w io.Writer
}
```

## Configuration Schema

```go
// pkg/config/config.go

type Config struct {
    Verbose   bool                       `yaml:"verbose"`
    Languages map[string]*LanguageConfig `yaml:"languages"`
}

type LanguageConfig struct {
    Enabled         *bool    `yaml:"enabled"`         // nil = auto-detect
    Paths           []string `yaml:"paths"`           // specific paths
    Test            *bool    `yaml:"test"`            // run tests
    Lint            *bool    `yaml:"lint"`            // run linter
    Format          *bool    `yaml:"format"`          // check formatting
    Coverage        *bool    `yaml:"coverage"`        // show coverage
    ExcludeCoverage string   `yaml:"exclude_coverage"` // exclude dirs
}

// Load loads configuration from .releaseagent.yaml
func Load(dir string) (*Config, error)

// IsLanguageEnabled checks if a language is enabled.
func (c *Config) IsLanguageEnabled(lang string) bool

// GetLanguageConfig gets language-specific configuration.
func (c *Config) GetLanguageConfig(lang string) *LanguageConfig
```

## CLI Design

### Command Structure

```
releaseagent
├── check          # Run validation checks
├── validate       # Comprehensive validation (QA, Docs, Release, Security)
├── changelog      # Generate changelog
├── readme         # Update README
├── roadmap        # Update roadmap
├── release        # Full release workflow
└── version        # Show version
```

### Global Flags

```
--verbose, -v       Show detailed output
--interactive, -i   Enable interactive mode
--json              Output as structured data
--format            Output format: "toon" (default) or "json"
```

### Command-Specific Flags

```
check:
  --no-test           Skip tests
  --no-lint           Skip linting
  --no-format         Skip format checks
  --coverage          Show coverage
  --go-no-go          NASA-style Go/No-Go report

validate:
  --version           Target version for validation
  --skip-qa           Skip QA checks
  --skip-docs         Skip documentation checks
  --skip-security     Skip security checks

changelog:
  --since             Start from this tag
  --dry-run           Preview changes

release:
  --dry-run           Preview without changes
  --skip-checks       Skip validation (dangerous)
  --skip-ci           Don't wait for CI

readme:
  --version           Version to update to
  --dry-run           Preview changes

roadmap:
  --dry-run           Preview changes
```

### Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | Check/action failed |
| 2 | Configuration error |
| 3 | Dependency missing |
| 4 | User cancelled (interactive mode) |
| 5 | CI timeout/failure |

## Release Workflow

The `releaseagent release` command executes a 9-step workflow:

```
1. Validate version      [REQUIRED]  Check format, ensure tag doesn't exist
2. Check working dir     [REQUIRED]  Ensure no uncommitted changes
3. Run validation        [REQUIRED]  Build, test, lint, format checks
4. Generate changelog    [OPTIONAL]  Update CHANGELOG.md via schangelog
5. Update roadmap        [OPTIONAL]  Update ROADMAP.md via sroadmap
6. Create release commit [REQUIRED]  Commit with "chore(release): vX.Y.Z"
7. Push to remote        [REQUIRED]  Push commits to origin
8. Wait for CI           [OPTIONAL]  Poll GitHub Actions until pass/fail
9. Create tag            [REQUIRED]  Create and push version tag
```

## Checker Implementations

### GoChecker

Validates Go projects with these checks:

| Check | Type | Command/Logic |
|-------|------|---------------|
| no local replace | Hard | Parse go.mod for local replace directives |
| mod tidy | Hard | `go mod tidy -diff` |
| build | Hard | `go build ./...` |
| gofmt | Hard | `gofmt -l .` |
| golangci-lint | Hard | `golangci-lint run` |
| tests | Hard | `go test -v ./...` |
| error handling | Hard | Detect `_ = err` patterns |
| untracked refs | Soft | Check tracked files for untracked references |
| coverage | Soft | `gocoverbadge` for coverage report |

### TypeScriptChecker

Validates TypeScript/JavaScript projects:

| Check | Type | Command |
|-------|------|---------|
| eslint | Hard | `npm run lint` or `eslint .` |
| prettier | Hard | `prettier --check .` |
| type check | Hard | `tsc --noEmit` |
| tests | Hard | `npm test` |

### SecurityChecker

Security and compliance checks:

| Check | Type | Command/Logic |
|-------|------|---------------|
| LICENSE | Hard | Check LICENSE file exists |
| govulncheck | Hard | `govulncheck ./...` |
| dependency audit | Hard | Check for retracted versions |
| secret detection | Hard | Pattern matching for hardcoded secrets |

### DocChecker

Documentation validation:

| Check | Type | Logic |
|-------|------|-------|
| README.md | Hard | Exists and >100 bytes |
| PRD.md | Soft | Exists if expected |
| TRD.md | Soft | Exists if expected |
| CHANGELOG.md | Soft | Exists and up to date |
| Release notes | Soft | Exists for version |

### ReleaseChecker

Release-specific validation:

| Check | Type | Logic |
|-------|------|-------|
| Version format | Hard | Valid semver format |
| Tag available | Hard | Tag doesn't already exist |
| Git status | Hard | Working directory clean |
| CI config | Soft | CI configuration exists |

## Output Formats

### Human-Readable (Default)

Terminal output with colored symbols:

- `✓` = passed (green)
- `✗` = failed (red)
- `⚠` = warning (yellow)
- `⊘` = skipped (gray)

### JSON Format

Standard JSON for programmatic consumption:

```json
{
  "type": "result",
  "result": {
    "name": "Go: build",
    "success": true,
    "output": "Build completed successfully"
  }
}
```

### TOON Format

Token-Oriented Object Notation for LLM efficiency (~8x more token-efficient):

```
Result:
  Name: "Go: build"
  Success: true
  Output: "Build completed successfully"
```

## Plugin Architecture

### Canonical Specification

Plugins are defined in JSON (`plugins/spec/`) and generated for each platform:

```
plugins/spec/
├── plugin.json           # Plugin metadata
├── commands/             # Command definitions
│   ├── release.json
│   ├── check.json
│   ├── changelog.json
│   └── version-next.json
├── skills/               # Skill definitions
│   ├── version-analysis.json
│   └── commit-classification.json
└── agents/               # Agent definitions
    └── release-coordinator.json
```

### Platform Adapters

The generator (`plugins/generate/main.go`) converts specs to platform formats:

- **Claude**: Markdown with YAML frontmatter
- **Gemini**: TOML configuration files

### Hooks

Lifecycle hooks for Claude Code integration:

- **SessionStart**: Check dependencies are installed
- **PreToolUse**: Validate release commands before execution

## Dependencies

### Go Packages

| Package | Purpose |
|---------|---------|
| `github.com/spf13/cobra` | CLI framework |
| `gopkg.in/yaml.v3` | YAML configuration |
| `github.com/toon-format/toon-go` | TOON output format |

### External Tools

| Tool | Required For |
|------|--------------|
| `git` | All operations |
| `gh` | CI status checking |
| `go` | Go checks |
| `golangci-lint` | Go linting |
| `schangelog` | Changelog generation |
| `sroadmap` | Roadmap updates |
| `gocoverbadge` | Coverage badges |
| `govulncheck` | Vulnerability scanning |

## Security Considerations

1. **No credential storage** - Uses system git/gh credentials
2. **No remote code execution** - Only runs local tools
3. **Validate tag names** - Prevent injection via tag names
4. **Dry-run support** - Preview before destructive operations
5. **CI verification** - Wait for CI before tagging releases

## Performance Considerations

1. **Parallel checks** - Independent checks run concurrently
2. **Incremental changelog** - Only process new commits
3. **Cached detection** - Language detection cached per session
4. **Lazy tool detection** - Only check for tools when needed
5. **TOON output** - ~8x more token-efficient for LLM integration
