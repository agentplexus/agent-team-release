// Package workflow provides workflow orchestration for multi-step releases.
package workflow

import (
	"fmt"
	"strings"
	"time"
)

// StepType defines the type of workflow step.
type StepType int

const (
	// StepTypeFunc is a step that runs a function.
	StepTypeFunc StepType = iota
	// StepTypeComposite is a step that contains sub-steps.
	StepTypeComposite
)

// StepFunc is a function that executes a step.
// It receives the context and returns an error if the step fails.
type StepFunc func(ctx *Context) error

// Step represents a single step in a workflow.
type Step struct {
	Name        string   // Step name for display
	Description string   // Human-readable description
	Type        StepType // Step type
	Required    bool     // If true, workflow fails if step fails
	Func        StepFunc // Function to execute (for StepTypeFunc)
	SubSteps    []Step   // Sub-steps (for StepTypeComposite)
}

// Workflow defines a sequence of steps.
type Workflow struct {
	Name        string
	Description string
	Steps       []Step
}

// Context provides context for step execution.
type Context struct {
	Dir         string            // Working directory
	Version     string            // Target version
	DryRun      bool              // If true, don't make changes
	Verbose     bool              // Show detailed output
	Interactive bool              // Enable interactive mode
	JSONOutput  bool              // Output JSON for Claude Code
	SkipChecks  bool              // Skip validation checks
	SkipCI      bool              // Skip CI wait
	Data        map[string]string // Arbitrary data passed between steps
	Output      *strings.Builder  // Captured output
}

// NewContext creates a new workflow context.
func NewContext(dir string, version string) *Context {
	return &Context{
		Dir:     dir,
		Version: version,
		Data:    make(map[string]string),
		Output:  &strings.Builder{},
	}
}

// Log writes a message to the context output.
func (c *Context) Log(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	c.Output.WriteString(msg)
	if !strings.HasSuffix(msg, "\n") {
		c.Output.WriteString("\n")
	}
}

// StepResult represents the result of a step execution.
type StepResult struct {
	Name     string
	Success  bool
	Skipped  bool
	Error    error
	Output   string
	Duration time.Duration
	SubSteps []StepResult // Results of sub-steps (for composite)
}

// WorkflowResult represents the result of a workflow execution.
type WorkflowResult struct {
	Name     string
	Success  bool
	Steps    []StepResult
	Duration time.Duration
	Output   string
}

// Runner executes workflows.
type Runner struct {
	DryRun      bool
	Verbose     bool
	Interactive bool
	JSONOutput  bool
}

// NewRunner creates a new workflow runner.
func NewRunner() *Runner {
	return &Runner{}
}

// Run executes a workflow and returns the results.
func (r *Runner) Run(w *Workflow, ctx *Context) *WorkflowResult {
	start := time.Now()

	// Apply runner settings to context
	ctx.DryRun = r.DryRun
	ctx.Verbose = r.Verbose
	ctx.Interactive = r.Interactive
	ctx.JSONOutput = r.JSONOutput

	result := &WorkflowResult{
		Name:    w.Name,
		Success: true,
	}

	ctx.Log("=== %s ===\n", w.Name)
	if w.Description != "" {
		ctx.Log("%s\n", w.Description)
	}
	ctx.Log("")

	for _, step := range w.Steps {
		stepResult := r.runStep(&step, ctx)
		result.Steps = append(result.Steps, stepResult)

		if !stepResult.Success && !stepResult.Skipped {
			if step.Required {
				result.Success = false
				ctx.Log("\n❌ Workflow failed at step: %s\n", step.Name)
				break
			}
			ctx.Log("⚠ Step %s failed but is not required, continuing...\n", step.Name)
		}
	}

	result.Duration = time.Since(start)
	result.Output = ctx.Output.String()

	if result.Success {
		ctx.Log("\n✅ %s completed successfully\n", w.Name)
	}

	return result
}

// runStep executes a single step.
func (r *Runner) runStep(step *Step, ctx *Context) StepResult {
	start := time.Now()

	result := StepResult{
		Name: step.Name,
	}

	ctx.Log("→ %s", step.Name)
	if step.Description != "" && ctx.Verbose {
		ctx.Log("  %s", step.Description)
	}

	switch step.Type {
	case StepTypeFunc:
		if step.Func == nil {
			result.Skipped = true
			result.Output = "No function defined"
			ctx.Log(" [skipped]\n")
			return result
		}

		err := step.Func(ctx)
		if err != nil {
			result.Success = false
			result.Error = err
			result.Output = err.Error()
			ctx.Log(" [failed: %v]\n", err)
		} else {
			result.Success = true
			ctx.Log(" [done]\n")
		}

	case StepTypeComposite:
		ctx.Log("\n")
		allSuccess := true
		for _, subStep := range step.SubSteps {
			subResult := r.runStep(&subStep, ctx)
			result.SubSteps = append(result.SubSteps, subResult)
			if !subResult.Success && !subResult.Skipped && subStep.Required {
				allSuccess = false
				break
			}
		}
		result.Success = allSuccess
	}

	result.Duration = time.Since(start)
	return result
}

// Summary returns a summary of the workflow result.
func (wr *WorkflowResult) Summary() string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "Workflow: %s\n", wr.Name)
	fmt.Fprintf(&sb, "Status: %s\n", statusEmoji(wr.Success))
	fmt.Fprintf(&sb, "Duration: %s\n", wr.Duration.Round(time.Millisecond))
	sb.WriteString("\nSteps:\n")

	for _, step := range wr.Steps {
		status := "✓"
		if step.Skipped {
			status = "⊘"
		} else if !step.Success {
			status = "✗"
		}
		fmt.Fprintf(&sb, "  %s %s (%s)\n", status, step.Name, step.Duration.Round(time.Millisecond))

		for _, sub := range step.SubSteps {
			subStatus := "✓"
			if sub.Skipped {
				subStatus = "⊘"
			} else if !sub.Success {
				subStatus = "✗"
			}
			fmt.Fprintf(&sb, "    %s %s\n", subStatus, sub.Name)
		}
	}

	return sb.String()
}

func statusEmoji(success bool) string {
	if success {
		return "✅ Success"
	}
	return "❌ Failed"
}

// JSONResult represents a workflow result in structured format.
type JSONResult struct {
	Type         string           `json:"type" toon:"type"`
	WorkflowName string           `json:"workflow_name" toon:"workflow_name"`
	Success      bool             `json:"success" toon:"success"`
	Duration     string           `json:"duration" toon:"duration"`
	Steps        []JSONStepResult `json:"steps" toon:"steps"`
}

// JSONStepResult represents a step result in structured format.
type JSONStepResult struct {
	Name     string           `json:"name" toon:"name"`
	Success  bool             `json:"success" toon:"success"`
	Skipped  bool             `json:"skipped,omitempty" toon:"skipped,omitempty"`
	Error    string           `json:"error,omitempty" toon:"error,omitempty"`
	Duration string           `json:"duration" toon:"duration"`
	SubSteps []JSONStepResult `json:"sub_steps,omitempty" toon:"sub_steps,omitempty"`
}

// ToJSON converts the workflow result to a JSON-serializable structure.
func (wr *WorkflowResult) ToJSON() JSONResult {
	steps := make([]JSONStepResult, len(wr.Steps))
	for i, step := range wr.Steps {
		steps[i] = stepToJSON(step)
	}

	return JSONResult{
		Type:         "workflow_result",
		WorkflowName: wr.Name,
		Success:      wr.Success,
		Duration:     wr.Duration.Round(time.Millisecond).String(),
		Steps:        steps,
	}
}

func stepToJSON(step StepResult) JSONStepResult {
	result := JSONStepResult{
		Name:     step.Name,
		Success:  step.Success,
		Skipped:  step.Skipped,
		Duration: step.Duration.Round(time.Millisecond).String(),
	}
	if step.Error != nil {
		result.Error = step.Error.Error()
	}
	if len(step.SubSteps) > 0 {
		result.SubSteps = make([]JSONStepResult, len(step.SubSteps))
		for i, sub := range step.SubSteps {
			result.SubSteps[i] = stepToJSON(sub)
		}
	}
	return result
}
