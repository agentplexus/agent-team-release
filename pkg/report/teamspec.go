// Copyright 2025 John Wang. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package report

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	multiagentspec "github.com/agentplexus/multi-agent-spec/sdk/go"
)

// LoadTeamSpec loads a team.json file from the given directory.
// Returns a multiagentspec.Team parsed from the JSON file.
func LoadTeamSpec(dir string) (*multiagentspec.Team, error) {
	path := filepath.Join(dir, "team.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading team.json: %w", err)
	}

	var spec multiagentspec.Team
	if err := json.Unmarshal(data, &spec); err != nil {
		return nil, fmt.Errorf("parsing team.json: %w", err)
	}

	return &spec, nil
}

// GetValidationSteps returns only the validation steps (excludes execute-release).
func GetValidationSteps(team *multiagentspec.Team) []multiagentspec.Step {
	if team.Workflow == nil {
		return nil
	}
	var steps []multiagentspec.Step
	for _, step := range team.Workflow.Steps {
		// Exclude execution steps, only include validation steps
		if step.Name != "execute-release" {
			steps = append(steps, step)
		}
	}
	return steps
}

// Phase represents a workflow phase containing multiple steps.
type Phase struct {
	Name  string
	Steps []multiagentspec.Step
}

// GetPhases groups steps into phases based on dependencies.
// Phase 1 (REVIEW): Steps with no dependencies or only pm-validation dependency
// Phase 2 (EXECUTE): Steps that depend on validation completion
func GetPhases(team *multiagentspec.Team) []Phase {
	phases := []Phase{
		{Name: "PHASE 1: REVIEW", Steps: []multiagentspec.Step{}},
		{Name: "PHASE 2: EXECUTE", Steps: []multiagentspec.Step{}},
	}

	if team.Workflow == nil {
		return phases
	}

	for _, step := range team.Workflow.Steps {
		if step.Name == "execute-release" {
			phases[1].Steps = append(phases[1].Steps, step)
		} else {
			phases[0].Steps = append(phases[0].Steps, step)
		}
	}

	return phases
}

// BuildReportFromSpec creates a TeamReport ensuring all steps from the spec are present.
// Results that don't have corresponding validation data are marked as SKIP.
func BuildReportFromSpec(spec *multiagentspec.Team, results map[string][]multiagentspec.TaskResult, project, version string) *multiagentspec.TeamReport {
	var teams []multiagentspec.TeamSection

	// Get phases for display
	phases := GetPhases(spec)
	phaseName := "PHASE 1: REVIEW"
	if len(phases) > 0 {
		phaseName = phases[0].Name
	}

	// Build teams from spec, ensuring all steps are represented
	for _, step := range GetValidationSteps(spec) {
		tasks, hasResults := results[step.Name]

		team := multiagentspec.TeamSection{
			ID:      step.Name,
			Name:    step.Agent,
			AgentID: step.Agent,
		}

		if hasResults {
			team.Tasks = tasks
		} else {
			// No results - mark as pending/skipped
			team.Tasks = []multiagentspec.TaskResult{
				{
					ID:     "validation",
					Status: multiagentspec.StatusSkip,
					Detail: "Not executed",
				},
			}
		}
		team.Status = team.OverallStatus()

		teams = append(teams, team)
	}

	report := &multiagentspec.TeamReport{
		Schema:      "https://raw.githubusercontent.com/agentplexus/multi-agent-spec/main/schema/report/team-report.schema.json",
		Project:     project,
		Version:     version,
		Target:      version,
		Phase:       phaseName,
		Teams:       teams,
		GeneratedAt: time.Now().UTC(),
		GeneratedBy: "release-agent-team",
	}
	report.Status = report.ComputeOverallStatus()

	return report
}

// StepResultMap is a helper to collect task results by step name.
type StepResultMap map[string][]multiagentspec.TaskResult

// NewStepResultMap creates a new empty step result map.
func NewStepResultMap() StepResultMap {
	return make(StepResultMap)
}

// Add adds task results to a step.
func (m StepResultMap) Add(stepName string, tasks []multiagentspec.TaskResult) {
	m[stepName] = tasks
}

// AddTask adds a single task result to a step.
func (m StepResultMap) AddTask(stepName string, task multiagentspec.TaskResult) {
	m[stepName] = append(m[stepName], task)
}
