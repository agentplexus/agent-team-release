// Package checks provides pre-push checks for various languages.
package checks

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Result represents the result of a check.
type Result struct {
	Name    string
	Passed  bool
	Output  string
	Error   error
	Skipped bool
	Reason  string
}

// Checker is the interface for language-specific checks.
type Checker interface {
	Name() string
	Check(dir string, opts Options) []Result
}

// Options configures which checks to run.
type Options struct {
	Test     bool
	Lint     bool
	Format   bool
	Coverage bool
	Verbose  bool

	// Language-specific options
	GoExcludeCoverage string // directories to exclude from coverage (e.g., "cmd")
}

// DefaultOptions returns the default check options.
func DefaultOptions() Options {
	return Options{
		Test:              true,
		Lint:              true,
		Format:            true,
		Coverage:          false,
		Verbose:           false,
		GoExcludeCoverage: "cmd",
	}
}

// RunCommand executes a command and returns the result.
func RunCommand(name string, dir string, command string, args ...string) Result {
	cmd := exec.Command(command, args...)
	cmd.Dir = dir

	output, err := cmd.CombinedOutput()

	return Result{
		Name:   name,
		Passed: err == nil,
		Output: strings.TrimSpace(string(output)),
		Error:  err,
	}
}

// CommandExists checks if a command is available in PATH.
func CommandExists(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}

// PrintResults prints check results to stdout.
func PrintResults(results []Result, verbose bool) (passed int, failed int, skipped int) {
	for _, r := range results {
		if r.Skipped {
			fmt.Printf("⊘ %s (skipped: %s)\n", r.Name, r.Reason)
			skipped++
			continue
		}

		if r.Passed {
			fmt.Printf("✓ %s\n", r.Name)
			passed++
		} else {
			fmt.Printf("✗ %s\n", r.Name)
			failed++
		}

		if verbose || !r.Passed {
			if r.Output != "" {
				// Indent output
				lines := strings.Split(r.Output, "\n")
				for _, line := range lines {
					fmt.Printf("  %s\n", line)
				}
			}
			if r.Error != nil && r.Output == "" {
				fmt.Printf("  Error: %v\n", r.Error)
			}
		}
	}

	return passed, failed, skipped
}

// FileExists checks if a file exists.
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
