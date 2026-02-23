// Copyright 2025 John Wang. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package report

import (
	"bytes"
	"strings"
	"testing"

	multiagentspec "github.com/agentplexus/multi-agent-spec/sdk/go"
)

func TestStatusIcon(t *testing.T) {
	tests := []struct {
		status multiagentspec.Status
		want   string
	}{
		{multiagentspec.StatusGo, "\U0001F7E2"},   // 🟢
		{multiagentspec.StatusWarn, "\U0001F7E1"}, // 🟡
		{multiagentspec.StatusNoGo, "\U0001F534"}, // 🔴
		{multiagentspec.StatusSkip, "\u26AA"},     // ⚪
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			if got := tt.status.Icon(); got != tt.want {
				t.Errorf("Status(%q).Icon() = %q, want %q", tt.status, got, tt.want)
			}
		})
	}
}

func TestTeamOverallStatus(t *testing.T) {
	tests := []struct {
		name   string
		tasks []multiagentspec.TaskResult
		want   multiagentspec.Status
	}{
		{
			name:   "all GO",
			tasks: []multiagentspec.TaskResult{{Status: multiagentspec.StatusGo}, {Status: multiagentspec.StatusGo}},
			want:   multiagentspec.StatusGo,
		},
		{
			name:   "one NO-GO",
			tasks: []multiagentspec.TaskResult{{Status: multiagentspec.StatusGo}, {Status: multiagentspec.StatusNoGo}},
			want:   multiagentspec.StatusNoGo,
		},
		{
			name:   "one WARN",
			tasks: []multiagentspec.TaskResult{{Status: multiagentspec.StatusGo}, {Status: multiagentspec.StatusWarn}},
			want:   multiagentspec.StatusWarn,
		},
		{
			name:   "all SKIP",
			tasks: []multiagentspec.TaskResult{{Status: multiagentspec.StatusSkip}, {Status: multiagentspec.StatusSkip}},
			want:   multiagentspec.StatusSkip,
		},
		{
			name:   "NO-GO takes precedence over WARN",
			tasks: []multiagentspec.TaskResult{{Status: multiagentspec.StatusWarn}, {Status: multiagentspec.StatusNoGo}},
			want:   multiagentspec.StatusNoGo,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			team := multiagentspec.TeamSection{Tasks: tt.tasks}
			if got := team.OverallStatus(); got != tt.want {
				t.Errorf("TeamSection.OverallStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReportIsGo(t *testing.T) {
	tests := []struct {
		name  string
		teams []multiagentspec.TeamSection
		want  bool
	}{
		{
			name: "all teams GO",
			teams: []multiagentspec.TeamSection{
				{Tasks: []multiagentspec.TaskResult{{Status: multiagentspec.StatusGo}}, Status: multiagentspec.StatusGo},
				{Tasks: []multiagentspec.TaskResult{{Status: multiagentspec.StatusGo}}, Status: multiagentspec.StatusGo},
			},
			want: true,
		},
		{
			name: "one team WARN is still GO",
			teams: []multiagentspec.TeamSection{
				{Tasks: []multiagentspec.TaskResult{{Status: multiagentspec.StatusGo}}, Status: multiagentspec.StatusGo},
				{Tasks: []multiagentspec.TaskResult{{Status: multiagentspec.StatusWarn}}, Status: multiagentspec.StatusWarn},
			},
			want: true,
		},
		{
			name: "one team NO-GO",
			teams: []multiagentspec.TeamSection{
				{Tasks: []multiagentspec.TaskResult{{Status: multiagentspec.StatusGo}}, Status: multiagentspec.StatusGo},
				{Tasks: []multiagentspec.TaskResult{{Status: multiagentspec.StatusNoGo}}, Status: multiagentspec.StatusNoGo},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			report := &multiagentspec.TeamReport{Teams: tt.teams}
			if got := report.IsGo(); got != tt.want {
				t.Errorf("TeamReport.IsGo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRenderer(t *testing.T) {
	report := &multiagentspec.TeamReport{
		Project: "github.com/agentplexus/agent-team-release",
		Version: "v0.3.0",
		Target:  "v0.3.0 (release automation)",
		Phase:   "PHASE 1: REVIEW",
		Teams: []multiagentspec.TeamSection{
			{
				ID:   "qa-validation",
				Name: "qa",
				Tasks: []multiagentspec.TaskResult{
					{ID: "build", Status: multiagentspec.StatusGo, Detail: ""},
					{ID: "tests", Status: multiagentspec.StatusGo, Detail: "42 tests passed"},
					{ID: "lint", Status: multiagentspec.StatusGo, Detail: ""},
				},
				Status: multiagentspec.StatusGo,
			},
			{
				ID:   "security-validation",
				Name: "security",
				Tasks: []multiagentspec.TaskResult{
					{ID: "license", Status: multiagentspec.StatusGo, Detail: "MIT License"},
					{ID: "vulnerability-scan", Status: multiagentspec.StatusWarn, Detail: "1 deprecated"},
				},
				Status: multiagentspec.StatusWarn,
			},
		},
		Status: multiagentspec.StatusWarn,
	}

	var buf bytes.Buffer
	renderer := multiagentspec.NewRenderer(&buf)
	err := renderer.Render(report)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	output := buf.String()

	// Check for expected content
	// Note: teamHeader uses team.Name directly (not "ID (Name)" format)
	expectedContent := []string{
		"TEAM STATUS REPORT",
		"github.com/agentplexus/agent-team-release",
		"v0.3.0 (release automation)",
		"PHASE 1: REVIEW",
		"qa",      // team.Name
		"build",
		"42 tests passed",
		"security", // team.Name
		"MIT License",
		"1 deprecated",
		"GO for v0.3.0",
	}

	for _, expected := range expectedContent {
		if !strings.Contains(output, expected) {
			t.Errorf("Output missing expected content: %q", expected)
		}
	}

	// Check box structure
	if !strings.Contains(output, "╔") {
		t.Error("Output missing top border")
	}
	if !strings.Contains(output, "╚") {
		t.Error("Output missing bottom border")
	}
	if !strings.Contains(output, "║") {
		t.Error("Output missing side borders")
	}
}

func TestRendererNoGo(t *testing.T) {
	report := &multiagentspec.TeamReport{
		Version: "v0.3.0",
		Teams: []multiagentspec.TeamSection{
			{
				ID:   "qa-validation",
				Name: "qa",
				Tasks: []multiagentspec.TaskResult{
					{ID: "build", Status: multiagentspec.StatusNoGo, Detail: "compilation failed"},
				},
				Status: multiagentspec.StatusNoGo,
			},
		},
		Status: multiagentspec.StatusNoGo,
	}

	var buf bytes.Buffer
	renderer := multiagentspec.NewRenderer(&buf)
	err := renderer.Render(report)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	output := buf.String()

	if !strings.Contains(output, "NO-GO for v0.3.0") {
		t.Error("Output should contain NO-GO message")
	}
}
