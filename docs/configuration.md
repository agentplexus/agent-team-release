# Configuration

Release Agent can be configured via a `.releaseagent.yaml` file in your repository root.

## Configuration File

Create `.releaseagent.yaml` in your project root:

```yaml
# Global settings
verbose: false

# Language-specific settings
languages:
  go:
    enabled: true
    test: true
    lint: true
    format: true
    coverage: false
    exclude_coverage: "cmd"

  typescript:
    enabled: true
    paths: ["frontend/"]
    test: true
    lint: true
    format: true

  javascript:
    enabled: false  # disable for this repo
```

## Global Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `verbose` | bool | `false` | Enable verbose output |

## Language Options

Each language can be configured with these options:

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `enabled` | bool | `true` | Enable/disable language checks |
| `paths` | []string | auto | Specific paths to check |
| `test` | bool | `true` | Run tests |
| `lint` | bool | `true` | Run linter |
| `format` | bool | `true` | Check formatting |

### Go-Specific Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `coverage` | bool | `false` | Show coverage report |
| `exclude_coverage` | string | `"cmd"` | Directories to exclude from coverage |

## Example Configurations

### Go Project

```yaml
verbose: false

languages:
  go:
    enabled: true
    test: true
    lint: true
    format: true
    coverage: true
    exclude_coverage: "cmd,internal/generated"
```

### TypeScript Project

```yaml
languages:
  typescript:
    enabled: true
    test: true
    lint: true
    format: true
```

### Monorepo

```yaml
languages:
  go:
    enabled: true
    paths: ["backend/"]

  typescript:
    enabled: true
    paths: ["frontend/", "shared/"]

  javascript:
    enabled: false
```

### CI-Only Testing

For faster local checks, disable tests:

```yaml
# .releaseagent.yaml - local development
languages:
  go:
    test: false  # Run tests in CI only
    lint: true
    format: true
```

## Environment Variables

Some settings can be overridden via environment variables:

| Variable | Description |
|----------|-------------|
| `RELEASEAGENT_VERBOSE` | Enable verbose output |
| `RELEASEAGENT_CONFIG` | Path to config file |

## Command-Line Override

Command-line flags override configuration file settings:

```bash
# Config says test: true, but skip tests for this run
atrelease check --no-test
```

## Legacy Configuration

For backwards compatibility, `.prepush.yaml` is also supported but deprecated. Rename to `.releaseagent.yaml`:

```bash
mv .prepush.yaml .releaseagent.yaml
```

## Validation

To validate your configuration:

```bash
atrelease check --verbose
```

The verbose output shows which configuration options are being used.
