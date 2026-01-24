# Output Formats

Release Agent supports multiple output formats for different use cases.

## Available Formats

| Format | Flag | Description |
|--------|------|-------------|
| Human | (default) | Colored terminal output with symbols |
| JSON | `--json --format=json` | Standard JSON for programmatic use |
| TOON | `--json --format=toon` | Token-optimized format for LLMs |
| Team | `--format team` | Template-based box report (validate only) |

## Human-Readable (Default)

The default format provides colored terminal output with Unicode symbols:

```
=== Pre-push Checks ===

Detecting languages...
  Found: go in .

Running Go checks...

=== Summary ===
✓ Go: no local replace directives
✓ Go: mod tidy
✓ Go: build
✓ Go: gofmt
✓ Go: golangci-lint
✓ Go: tests
⚠ Go: untracked references (warning)

Passed: 6, Failed: 0, Skipped: 0, Warnings: 1

Pre-push checks passed with warnings.
```

### Symbols

| Symbol | Meaning |
|--------|---------|
| ✓ | Check passed |
| ✗ | Check failed |
| ⚠ | Warning |
| ⊘ | Skipped |

## JSON Format

Standard JSON output for programmatic consumption:

```bash
atrelease check --json --format=json
```

```json
{
  "results": [
    {
      "name": "Go: build",
      "passed": true,
      "skipped": false,
      "warning": false,
      "output": ""
    },
    {
      "name": "Go: tests",
      "passed": true,
      "skipped": false,
      "warning": false,
      "output": "ok  \tgithub.com/example/project\t0.042s"
    }
  ],
  "summary": {
    "passed": 7,
    "failed": 0,
    "skipped": 0,
    "warnings": 1
  }
}
```

## TOON Format

Token-Oriented Object Notation is approximately 8x more token-efficient than JSON, optimized for LLM consumption:

```bash
atrelease check --json --format=toon
```

```
RESULTS
name:Go: build|passed:true|skipped:false|warning:false
name:Go: tests|passed:true|output:ok github.com/example/project 0.042s
SUMMARY
passed:7|failed:0|skipped:0|warnings:1
```

### When to Use TOON

- When integrating with Claude Code or other LLMs
- When token efficiency matters (API costs)
- When output will be parsed by AI assistants

## Team Status Report

The team format provides a structured box report for release validation:

```bash
atrelease validate --format team --version=v1.0.0
```

```
╔════════════════════════════════════════════════════════════════════════════╗
║                             TEAM STATUS REPORT                             ║
╠════════════════════════════════════════════════════════════════════════════╣
║ Project: github.com/grokify/release-agent                                  ║
║ Target:  v1.0.0                                                            ║
╠════════════════════════════════════════════════════════════════════════════╣
║ RELEASE VALIDATION                                                         ║
╠════════════════════════════════════════════════════════════════════════════╣
║ qa-validation (qa)                                                         ║
║   build                    🟢 GO                                           ║
║   tests                    🟢 GO    42 tests passed                        ║
║   lint                     🟢 GO                                           ║
║   format                   🟢 GO                                           ║
╠════════════════════════════════════════════════════════════════════════════╣
║ docs-validation (documentation)                                            ║
║   readme                   🟢 GO                                           ║
║   changelog                🟢 GO                                           ║
║   release-notes            🟢 GO                                           ║
╠════════════════════════════════════════════════════════════════════════════╣
║ security-validation (security)                                             ║
║   license                  🟢 GO    MIT License                            ║
║   vulnerability-scan       🟡 WARN  1 deprecated                           ║
╠════════════════════════════════════════════════════════════════════════════╣
║                         🚀 TEAM: GO for v1.0.0 🚀                          ║
╚════════════════════════════════════════════════════════════════════════════╝
```

### Status Icons

| Icon | Status | Meaning |
|------|--------|---------|
| 🟢 | GO | Check passed |
| 🟡 | WARN | Warning (doesn't block) |
| 🔴 | NO-GO | Check failed |
| ⚪ | SKIP | Check skipped |

### When to Use Team Format

- Release readiness reviews
- Team status meetings
- Documentation of release decisions
- Audit trails

## Combining Formats

Some flags can be combined:

```bash
# JSON output with TOON format
atrelease check --json --format=toon

# Verbose human-readable
atrelease check --verbose

# Team format with verbose details
atrelease validate --format team --verbose
```

## Parsing Output

### JSON in Shell

```bash
# Get pass/fail status
atrelease check --json --format=json | jq '.summary.failed'

# List failed checks
atrelease check --json --format=json | jq '.results[] | select(.passed == false)'
```

### TOON in Python

```python
def parse_toon(output):
    results = []
    for line in output.split('\n'):
        if '|' in line:
            fields = dict(f.split(':') for f in line.split('|'))
            results.append(fields)
    return results
```
