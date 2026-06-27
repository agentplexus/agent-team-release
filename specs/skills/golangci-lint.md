---
name: golangci-lint
description: Run golangci-lint and fix Go lint errors with proper error handling patterns
triggers: [lint, golangci-lint, linter, go lint]
dependencies: [golangci-lint]
---

# Go Lint Skill

Run `golangci-lint run` on Go projects and fix lint errors following established patterns.

## Running the Linter

```bash
# Run in current directory
golangci-lint run

# Run in specific directory
golangci-lint run ./path/to/dir/...

# Run with verbose output
golangci-lint run -v

# Run specific linters
golangci-lint run --enable=errcheck,gosec
```

## Error Handling Priority Order

**CRITICAL: All errors must be handled. Never assign errors to `_`.**

When fixing unhandled error lint violations, follow this priority order:

### Priority 1: Panic (Programming Errors)

Use `panic()` when an error should never happen (invariant violation, programming error):

```go
// json.Marshal of simple struct should never fail
data, err := json.Marshal(simpleStruct)
if err != nil {
    panic(fmt.Sprintf("json.Marshal failed: %v", err))
}

// regexp.MustCompile equivalent
re, err := regexp.Compile(pattern)
if err != nil {
    panic(fmt.Sprintf("invalid regex pattern: %v", err))
}
```

### Priority 2: Return Error

If the function signature can return an error, return it:

```go
func processData(data []byte) error {
    result, err := parseData(data)
    if err != nil {
        return fmt.Errorf("parse failed: %w", err)
    }
    // ...
    return nil
}
```

### Priority 3: Modify Function to Return Error

If the function doesn't return error but can be modified:

```go
// Before
func loadConfig() Config {
    data, _ := os.ReadFile("config.json")  // BAD
    // ...
}

// After
func loadConfig() (Config, error) {
    data, err := os.ReadFile("config.json")
    if err != nil {
        return Config{}, fmt.Errorf("read config: %w", err)
    }
    // ...
    return cfg, nil
}
```

### Priority 4: Log Error (Interface Compliance)

When function must fulfill an interface without error return, use logging:

```go
import "github.com/grokify/mogo/log/slogutil"

// Get logger from context
logger := slogutil.LoggerFromContext(ctx, nil)

// Log the error
if err != nil {
    slogutil.LogOrNot(ctx, logger, slog.LevelError,
        "operation failed",
        slog.String("error", err.Error()),
        slog.String("operation", "write"))
}

// Alternative with variadic args
slogutil.LogOrNotAny(ctx, logger, slog.LevelError,
    "operation failed", "error", err)
```

### Priority 5: Raise to Human

If none of the above approaches work, raise to the human with explanation:

```
I cannot automatically fix this error because:
- The function implements interface X which has no error return
- There is no context available for logging
- Panic is not appropriate for this recoverable error

Suggested approaches:
1. Add context parameter to the function
2. Create a wrapper that handles the error
3. Use a package-level logger
```

## Common Lint Fixes

### File Permissions (G306)

Always use `0o600` for file permissions, not `0o644` or `0644`:

```go
// BAD
os.WriteFile("file.txt", data, 0644)
os.WriteFile("file.txt", data, 0o644)

// GOOD
os.WriteFile("file.txt", data, 0o600)
```

### Unhandled Errors (errcheck)

```go
// BAD
file.Close()
writer.Flush()

// GOOD - defer with named return
defer func() {
    if cerr := file.Close(); cerr != nil && err == nil {
        err = cerr
    }
}()

// GOOD - inline check
if err := writer.Flush(); err != nil {
    return fmt.Errorf("flush failed: %w", err)
}
```

### HTTP Form Parsing (G120)

Use body limit before parsing forms:

```go
import "github.com/grokify/mogo/lintfix/gosec"

// Before parsing form
if err := gosec.LimitAndParseForm(w, r, gosec.G120MaxBytes.Webhook); err != nil {
    http.Error(w, "Bad Request", http.StatusBadRequest)
    return
}
// Use r.Form.Get("key") NOT r.FormValue("key")
value := r.Form.Get("key")
```

Max bytes constants:
- `gosec.G120MaxBytes.Form` - 1MB
- `gosec.G120MaxBytes.Multipart` - 32MB
- `gosec.G120MaxBytes.Webhook` - 64KB
- `gosec.G120MaxBytes.Twilio` - 64KB

### Gosec Nolint Comments

When suppression is appropriate, use helpers:

```go
import "github.com/grokify/mogo/lintfix/gosec"

// Generate nolint comment
comment := gosec.NolintG117(gosec.CommonReasons.OAuthTokenResponse)
// Output: "//nolint:gosec // G117: OAuth token response per RFC 6749"

// Common reasons available:
gosec.CommonReasons.OAuthTokenResponse    // OAuth token in response
gosec.CommonReasons.ShutdownHandler       // Shutdown runs after context cancelled
gosec.CommonReasons.HttptestServer        // Test uses httptest server URL
gosec.CommonReasons.PathFromCLIFlag       // Path from CLI flag
```

### Integer Overflow (G115)

```go
// BAD - potential overflow
count := int(uint64Value)

// GOOD - with bounds check
if uint64Value > math.MaxInt {
    return fmt.Errorf("value %d exceeds max int", uint64Value)
}
count := int(uint64Value)
```

### Context in Goroutines (G118)

```go
// BAD
go func() {
    doWork(context.Background())  // Using Background in goroutine
}()

// GOOD - pass context
go func(ctx context.Context) {
    doWork(ctx)
}(ctx)
```

## Remediation Database

Load remediation info for specific rules:

```go
import "github.com/grokify/mogo/lintfix"

db := lintfix.MustLoadRemediations()
fix := db.GetGosec("G120")

fmt.Println(fix.Remediation.Summary)
fmt.Println(fix.Remediation.CodeBefore)
fmt.Println(fix.Remediation.CodeAfter)

// Check if helper function exists
if fix.HasHelper() {
    fmt.Printf("Use: %s.%s\n", fix.Remediation.Package, fix.Remediation.Function)
}
```

## Workflow

1. **Run linter**: `golangci-lint run ./...`
2. **Review errors**: Group by type (errcheck, gosec, etc.)
3. **Apply fixes**: Follow priority order above
4. **Verify**: Re-run linter to confirm fixes
5. **Commit**: Use conventional commit: `fix: resolve lint errors`

## Resources

- Error logging: `github.com/grokify/mogo/log/slogutil`
- Lint fixes: `github.com/grokify/mogo/lintfix`
- Gosec helpers: `github.com/grokify/mogo/lintfix/gosec`
- Canonical source: `github.com/plexusone/assistantkit/capabilities/go`
