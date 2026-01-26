package checks

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	multiagentspec "github.com/agentplexus/multi-agent-spec/sdk/go"
)

// RunReleasekit executes `releasekit validate` and returns the results as checks.Result.
// It shells out to the releasekit CLI and parses the AgentResult JSON output.
func RunReleasekit(dir string, opts Options) ([]Result, error) {
	args := []string{"validate", "--format", "json"}

	if !opts.Lint {
		args = append(args, "--no-lint")
	}
	if !opts.Test {
		args = append(args, "--no-test")
	}
	if opts.Coverage {
		args = append(args, "--coverage")
	}
	if opts.Verbose {
		args = append(args, "--verbose")
	}

	args = append(args, dir)

	cmd := exec.Command("releasekit", args...)
	output, err := cmd.Output()

	// releasekit exits with code 2 for NO-GO, which is not an error for our purposes
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if exitErr.ExitCode() == 2 {
				// NO-GO status, but we have valid output
				output = append(output, exitErr.Stderr...)
			} else {
				return nil, fmt.Errorf("releasekit failed: %w\nstderr: %s", err, string(exitErr.Stderr))
			}
		} else {
			return nil, fmt.Errorf("releasekit not found or failed to execute: %w", err)
		}
	}

	// Parse AgentResult
	agentResult, err := multiagentspec.ParseAgentResult(output)
	if err != nil {
		return nil, fmt.Errorf("failed to parse releasekit output: %w\noutput: %s", err, string(output))
	}

	// Convert TaskResults to checks.Result
	return convertTaskResults(agentResult.Tasks), nil
}

// convertTaskResults converts multiagentspec.TaskResult slice to checks.Result slice.
func convertTaskResults(tasks []multiagentspec.TaskResult) []Result {
	results := make([]Result, 0, len(tasks))

	for _, t := range tasks {
		r := Result{
			Name: t.ID,
		}

		switch t.Status {
		case multiagentspec.StatusGo:
			r.Passed = true
		case multiagentspec.StatusNoGo:
			r.Passed = false
			r.Output = t.Detail
			if t.Metadata != nil {
				if out, ok := t.Metadata["output"].(string); ok {
					r.Output = out
				}
			}
		case multiagentspec.StatusWarn:
			r.Warning = true
			r.Passed = false // Warning with issue
			r.Output = t.Detail
			if t.Metadata != nil {
				if out, ok := t.Metadata["output"].(string); ok {
					r.Output = out
				}
			}
		case multiagentspec.StatusSkip:
			r.Skipped = true
			r.Reason = t.Detail
		}

		results = append(results, r)
	}

	return results
}

// ReleasekitAvailable checks if the releasekit CLI is installed and available.
func ReleasekitAvailable() bool {
	_, err := exec.LookPath("releasekit")
	return err == nil
}

// RunReleasekitRaw executes releasekit and returns the raw AgentResult.
// Use this when you want to work directly with multi-agent-spec types.
func RunReleasekitRaw(dir string, opts Options) (*multiagentspec.AgentResult, error) {
	args := []string{"validate", "--format", "json"}

	if !opts.Lint {
		args = append(args, "--no-lint")
	}
	if !opts.Test {
		args = append(args, "--no-test")
	}
	if opts.Coverage {
		args = append(args, "--coverage")
	}
	if opts.Verbose {
		args = append(args, "--verbose")
	}

	args = append(args, dir)

	cmd := exec.Command("releasekit", args...)
	output, err := cmd.Output()

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if exitErr.ExitCode() == 2 {
				// NO-GO status, combine stdout and stderr
				// Actually releasekit outputs JSON to stdout even on exit 2
				// Just use stdout
			} else {
				return nil, fmt.Errorf("releasekit failed: %w", err)
			}
		} else {
			return nil, fmt.Errorf("releasekit not found: %w", err)
		}
	}

	// Handle case where output might be empty due to error handling
	if len(output) == 0 {
		return nil, fmt.Errorf("releasekit produced no output")
	}

	return multiagentspec.ParseAgentResult(output)
}

// buildReleasekitArgs constructs the command line arguments for releasekit.
func buildReleasekitArgs(dir string, opts Options) []string {
	args := []string{"validate", "--format", "json"}

	if !opts.Lint {
		args = append(args, "--no-lint")
	}
	if !opts.Test {
		args = append(args, "--no-test")
	}
	if opts.Coverage {
		args = append(args, "--coverage")
	}
	if opts.Verbose {
		args = append(args, "--verbose")
	}

	args = append(args, dir)
	return args
}

// parseReleasekitOutput parses JSON output from releasekit into AgentResult.
func parseReleasekitOutput(data []byte) (*multiagentspec.AgentResult, error) {
	var result multiagentspec.AgentResult
	if err := json.Unmarshal(data, &result); err != nil {
		// Try to extract any useful error info
		if strings.Contains(string(data), "error") {
			return nil, fmt.Errorf("releasekit error: %s", string(data))
		}
		return nil, fmt.Errorf("invalid JSON from releasekit: %w", err)
	}
	return &result, nil
}
