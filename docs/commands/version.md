# version

Show version information.

## Usage

```bash
atrelease version
```

## Description

The `version` command displays the current version of Release Agent along with build information.

## Examples

```bash
atrelease version
```

## Output

```
atrelease version 0.3.0
```

With verbose flag:

```bash
atrelease version --verbose
```

```
atrelease version 0.3.0
  commit: abc123def
  built:  2026-01-12T10:00:00Z
  go:     go1.21.0
```

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Always succeeds |
