package checks

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

// GoChecker implements checks for Go projects.
type GoChecker struct{}

// Name returns the checker name.
func (c *GoChecker) Name() string {
	return "Go"
}

// Check runs Go checks on the specified directory.
func (c *GoChecker) Check(dir string, opts Options) []Result {
	var results []Result

	// Check for local replace directives
	results = append(results, c.checkNoLocalReplace(dir))

	// Format check
	if opts.Format {
		results = append(results, c.checkFormat(dir))
	}

	// Lint check
	if opts.Lint {
		results = append(results, c.checkLint(dir))
	}

	// Test check
	if opts.Test {
		results = append(results, c.checkTest(dir))
	}

	// Coverage check (informational)
	if opts.Coverage {
		results = append(results, c.checkCoverage(dir, opts.GoExcludeCoverage))
	}

	return results
}

func (c *GoChecker) checkNoLocalReplace(dir string) Result {
	name := "Go: no local replace directives"

	cmd := exec.Command("go", "mod", "edit", "-json")
	cmd.Dir = dir
	output, err := cmd.Output()
	if err != nil {
		return Result{
			Name:   name,
			Passed: false,
			Error:  err,
		}
	}

	// Check for local paths in replace directives
	// Local replaces typically have paths starting with . or /
	localReplacePattern := regexp.MustCompile(`"Path":\s*"[./]`)
	if localReplacePattern.Match(output) {
		return Result{
			Name:   name,
			Passed: false,
			Output: "go.mod contains local replace directives",
		}
	}

	return Result{
		Name:   name,
		Passed: true,
	}
}

func (c *GoChecker) checkFormat(dir string) Result {
	name := "Go: gofmt"

	// Check if any files need formatting
	cmd := exec.Command("gofmt", "-l", ".")
	cmd.Dir = dir
	output, err := cmd.Output()
	if err != nil {
		return Result{
			Name:   name,
			Passed: false,
			Error:  err,
		}
	}

	unformatted := strings.TrimSpace(string(output))
	if unformatted != "" {
		return Result{
			Name:   name,
			Passed: false,
			Output: "Files need formatting:\n" + unformatted,
		}
	}

	return Result{
		Name:   name,
		Passed: true,
	}
}

func (c *GoChecker) checkLint(dir string) Result {
	name := "Go: golangci-lint"

	if !CommandExists("golangci-lint") {
		return Result{
			Name:    name,
			Skipped: true,
			Reason:  "golangci-lint not installed",
		}
	}

	return RunCommand(name, dir, "golangci-lint", "run")
}

func (c *GoChecker) checkTest(dir string) Result {
	name := "Go: tests"
	return RunCommand(name, dir, "go", "test", "./...")
}

func (c *GoChecker) checkCoverage(dir string, exclude string) Result {
	name := "Go: coverage"

	if !CommandExists("gocoverbadge") {
		return Result{
			Name:    name,
			Skipped: true,
			Reason:  "gocoverbadge not installed",
		}
	}

	args := []string{"-dir", dir, "-badge-only"}
	if exclude != "" {
		args = append(args, "-exclude", exclude)
	}

	result := RunCommand(name, dir, "gocoverbadge", args...)
	// Coverage is informational, always passes
	result.Passed = true
	if result.Error == nil {
		fmt.Printf("  %s\n", result.Output)
	}
	return result
}
