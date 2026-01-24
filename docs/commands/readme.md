# readme

Update README badges and version references.

## Usage

```bash
atrelease readme [directory] [flags]
```

## Description

The `readme` command updates version references and badges in your README file. It can update coverage badges, version shields, and other dynamic content.

## Arguments

| Argument | Description | Default |
|----------|-------------|---------|
| `directory` | Directory to process | Current directory (`.`) |

## Flags

| Flag | Description |
|------|-------------|
| `--version` | Version to update to |
| `--dry-run` | Preview changes without writing |
| `--verbose`, `-v` | Show detailed output |

## Examples

```bash
# Update README with current version
atrelease readme

# Update to specific version
atrelease readme --version=v1.0.0

# Preview changes
atrelease readme --version=v1.0.0 --dry-run
```

## What Gets Updated

### Coverage Badge

If `gocoverbadge` is available, the coverage badge is regenerated:

```markdown
![Coverage](https://img.shields.io/badge/coverage-85%25-green)
```

### Version References

Version strings in installation instructions:

```markdown
go install github.com/example/project@v1.0.0
```

### Badge URLs

Version-specific badge URLs are updated to reflect the new version.

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | README updated successfully |
| 1 | Error updating README |
