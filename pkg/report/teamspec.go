// Copyright 2025 John Wang. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package report

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// TeamSpec represents the team.json specification.
type TeamSpec struct {
	Schema       string       `json:"$schema,omitempty"`
	Name         string       `json:"name"`
	Version      string       `json:"version"`
	Description  string       `json:"description"`
	Agents       []string     `json:"agents"`
	Orchestrator string       `json:"orchestrator"`
	Workflow     WorkflowSpec `json:"workflow"`
	Context      string       `json:"context"`
}

// WorkflowSpec represents the workflow definition.
type WorkflowSpec struct {
	Type  string     `json:"type"` // "dag" or "sequential"
	Steps []StepSpec `json:"steps"`
}

// StepSpec represents a single workflow step.
type StepSpec struct {
	Name      string   `json:"name"`
	Agent     string   `json:"agent"`
	DependsOn []string `json:"depends_on,omitempty"`
	Outputs   []string `json:"outputs,omitempty"`
}

// LoadTeamSpec loads a team.json file from the given directory.
func LoadTeamSpec(dir string) (*TeamSpec, error) {
	path := filepath.Join(dir, "team.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading team.json: %w", err)
	}

	var spec TeamSpec
	if err := json.Unmarshal(data, &spec); err != nil {
		return nil, fmt.Errorf("parsing team.json: %w", err)
	}

	return &spec, nil
}

// GetValidationSteps returns only the validation steps (excludes execute-release).
func (s *TeamSpec) GetValidationSteps() []StepSpec {
	var steps []StepSpec
	for _, step := range s.Workflow.Steps {
		// Exclude execution steps, only include validation steps
		if step.Name != "execute-release" {
			steps = append(steps, step)
		}
	}
	return steps
}

// GetPhases groups steps into phases based on dependencies.
// Phase 1 (REVIEW): Steps with no dependencies or only pm-validation dependency
// Phase 2 (EXECUTE): Steps that depend on validation completion
func (s *TeamSpec) GetPhases() []Phase {
	phases := []Phase{
		{Name: "PHASE 1: REVIEW", Steps: []StepSpec{}},
		{Name: "PHASE 2: EXECUTE", Steps: []StepSpec{}},
	}

	for _, step := range s.Workflow.Steps {
		if step.Name == "execute-release" {
			phases[1].Steps = append(phases[1].Steps, step)
		} else {
			phases[0].Steps = append(phases[0].Steps, step)
		}
	}

	return phases
}

// Phase represents a workflow phase containing multiple steps.
type Phase struct {
	Name  string
	Steps []StepSpec
}

// BuildReportFromSpec creates a TeamStatusReport ensuring all steps from the spec are present.
// Results that don't have corresponding validation data are marked as SKIP.
func BuildReportFromSpec(spec *TeamSpec, results map[string][]Check, project, version string) *TeamStatusReport {
	var teams []Team

	// Get phases for display
	phases := spec.GetPhases()
	phaseName := "PHASE 1: REVIEW"
	if len(phases) > 0 {
		phaseName = phases[0].Name
	}

	// Build teams from spec, ensuring all steps are represented
	for _, step := range spec.GetValidationSteps() {
		checks, hasResults := results[step.Name]

		team := Team{
			ID:   step.Name,
			Name: step.Agent,
		}

		if hasResults {
			team.Checks = checks
		} else {
			// No results - mark as pending/skipped
			team.Checks = []Check{
				{
					ID:     "validation",
					Status: StatusSkip,
					Detail: "Not executed",
				},
			}
		}

		teams = append(teams, team)
	}

	return &TeamStatusReport{
		Project: project,
		Version: version,
		Target:  version,
		Phase:   phaseName,
		Teams:   teams,
	}
}

// StepResultMap is a helper to collect results by step name.
type StepResultMap map[string][]Check

// NewStepResultMap creates a new empty step result map.
func NewStepResultMap() StepResultMap {
	return make(StepResultMap)
}

// Add adds checks to a step's results.
func (m StepResultMap) Add(stepName string, checks []Check) {
	m[stepName] = checks
}

// AddCheck adds a single check to a step's results.
func (m StepResultMap) AddCheck(stepName string, check Check) {
	m[stepName] = append(m[stepName], check)
}
