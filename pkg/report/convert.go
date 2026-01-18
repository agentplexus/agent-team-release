// Copyright 2025 John Wang. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package report

import (
	"fmt"
	"strings"
	"time"

	multiagentspec "github.com/agentplexus/multi-agent-spec/sdk/go"
	"github.com/agentplexus/release-agent-team/pkg/checks"
)

// TeamConfig maps validation areas to team IDs, names, and DAG dependencies.
type TeamConfig struct {
	Area      checks.ValidationArea
	ID        string
	Name      string
	DependsOn []string
}

// DefaultTeamConfigs returns the default team configurations with DAG dependencies.
// The DAG structure follows team.json:
// - pm-validation: no dependencies (runs first)
// - qa-validation, docs-validation, security-validation: depend on pm-validation (run in parallel)
// - release-validation: depends on all others (runs last)
func DefaultTeamConfigs() []TeamConfig {
	return []TeamConfig{
		{Area: checks.AreaPM, ID: "pm-validation", Name: "pm", DependsOn: nil},
		{Area: checks.AreaQA, ID: "qa-validation", Name: "qa", DependsOn: []string{"pm-validation"}},
		{Area: checks.AreaDocumentation, ID: "docs-validation", Name: "documentation", DependsOn: []string{"pm-validation"}},
		{Area: checks.AreaSecurity, ID: "security-validation", Name: "security", DependsOn: []string{"pm-validation"}},
		{Area: checks.AreaRelease, ID: "release-validation", Name: "release", DependsOn: []string{"pm-validation", "qa-validation", "docs-validation", "security-validation"}},
	}
}

// FromValidationReport converts a checks.ValidationReport to a multiagentspec.TeamReport.
func FromValidationReport(vr *checks.ValidationReport, project, target, phase string) *multiagentspec.TeamReport {
	configs := DefaultTeamConfigs()
	configMap := make(map[checks.ValidationArea]TeamConfig)
	for _, c := range configs {
		configMap[c.Area] = c
	}

	var teams []multiagentspec.TeamSection
	for _, ar := range vr.Areas {
		config, ok := configMap[ar.Area]
		if !ok {
			config = TeamConfig{
				ID:   strings.ToLower(string(ar.Area)) + "-validation",
				Name: strings.ToLower(string(ar.Area)),
			}
		}

		var teamChecks []multiagentspec.Check
		for _, r := range ar.Results {
			status := multiagentspec.StatusGo
			if r.Skipped {
				status = multiagentspec.StatusSkip
			} else if r.Warning && !r.Passed {
				status = multiagentspec.StatusWarn
			} else if !r.Passed {
				status = multiagentspec.StatusNoGo
			}

			// Extract check ID from name (e.g., "Go: build" -> "build")
			id := r.Name
			if idx := strings.Index(id, ": "); idx >= 0 {
				id = id[idx+2:]
			}
			// Convert to kebab-case
			id = strings.ToLower(strings.ReplaceAll(id, " ", "-"))

			// Use output as detail, truncate if needed
			detail := ""
			if r.Output != "" {
				detail = r.Output
				// Take first line only
				if idx := strings.Index(detail, "\n"); idx >= 0 {
					detail = detail[:idx]
				}
				// Truncate
				if len(detail) > 40 {
					detail = detail[:37] + "..."
				}
			}
			if r.Reason != "" && detail == "" {
				detail = r.Reason
			}

			teamChecks = append(teamChecks, multiagentspec.Check{
				ID:     id,
				Status: status,
				Detail: detail,
			})
		}

		team := multiagentspec.TeamSection{
			ID:        config.ID,
			Name:      config.Name,
			AgentID:   config.Name,
			DependsOn: config.DependsOn,
			Checks:    teamChecks,
		}
		team.Status = team.OverallStatus()
		teams = append(teams, team)
	}

	report := &multiagentspec.TeamReport{
		Schema:      "https://raw.githubusercontent.com/agentplexus/multi-agent-spec/main/schema/report/team-report.schema.json",
		Project:     project,
		Version:     vr.Version,
		Target:      target,
		Phase:       phase,
		Teams:       teams,
		GeneratedAt: time.Now().UTC(),
		GeneratedBy: "release-agent-team",
	}
	report.Status = report.ComputeOverallStatus()

	return report
}

// PMTeam creates a Product Management validation team section.
func PMTeam(version string, roadmapTotal, roadmapCompleted int, hasHighlights, hasBreaking, hasDeprecations bool) multiagentspec.TeamSection {
	teamChecks := []multiagentspec.Check{
		{
			ID:     "version-recommendation",
			Status: multiagentspec.StatusGo,
			Detail: version + " appropriate",
		},
		{
			ID:     "release-scope",
			Status: multiagentspec.StatusGo,
			Detail: "Phase complete",
		},
		{
			ID:     "changelog-quality",
			Status: boolToStatus(hasHighlights),
			Detail: boolToDetail(hasHighlights, "Highlights present", "Missing highlights"),
		},
		{
			ID:     "breaking-changes",
			Status: boolToStatus(!hasBreaking),
			Detail: boolToDetail(!hasBreaking, "None", "Has breaking changes"),
		},
		{
			ID:     "roadmap-alignment",
			Status: multiagentspec.StatusGo,
			Detail: formatFraction(roadmapCompleted, roadmapTotal) + " items completed",
		},
		{
			ID:     "deprecation-notices",
			Status: boolToStatus(!hasDeprecations),
			Detail: boolToDetail(!hasDeprecations, "No deprecations", "Has deprecations"),
		},
	}

	team := multiagentspec.TeamSection{
		ID:      "pm-validation",
		Name:    "pm",
		AgentID: "pm",
		Checks:  teamChecks,
	}
	team.Status = team.OverallStatus()

	return team
}

func boolToStatus(ok bool) multiagentspec.Status {
	if ok {
		return multiagentspec.StatusGo
	}
	return multiagentspec.StatusWarn
}

func boolToDetail(ok bool, okDetail, notOkDetail string) string {
	if ok {
		return okDetail
	}
	return notOkDetail
}

func formatFraction(num, total int) string {
	return fmt.Sprintf("%d/%d", num, total)
}
