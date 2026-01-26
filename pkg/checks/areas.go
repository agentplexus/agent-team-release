// Copyright 2025 John Wang. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package checks

import "fmt"

// ValidationArea represents a department/area of responsibility in the release process.
type ValidationArea string

const (
	// AreaPM represents Product Management validation.
	// Ensures the release scope, versioning, and product decisions are appropriate.
	// Checks: version-recommendation, release-scope, changelog-quality,
	// breaking-changes, roadmap-alignment, deprecation-notices.
	AreaPM ValidationArea = "PM"

	// AreaQA represents Quality Assurance validation.
	// Ensures the software works as expected and meets quality standards.
	// Checks: build, tests, lint, format, error handling compliance.
	AreaQA ValidationArea = "QA"

	// AreaDocumentation represents Documentation validation.
	// Confirms that user guides, release notes, and support materials are updated.
	// Checks: README, PRD, TRD, release notes, API docs, MkDocs site.
	AreaDocumentation ValidationArea = "Documentation"

	// AreaRelease represents Release Management validation.
	// Oversees the technical release process, versioning, and deployment.
	// Checks: version validation, changelog, git status, CI verification.
	AreaRelease ValidationArea = "Release"

	// AreaSecurity represents Security/Compliance validation.
	// Ensures the release complies with security policies and regulations.
	// Checks: dependency audit, license compliance, vulnerability scans.
	AreaSecurity ValidationArea = "Security"
)

// AreaResult represents the validation result for an area.
type AreaResult struct {
	Area    ValidationArea
	Status  AreaStatus
	Results []Result
}

// AreaStatus represents the Go/No-Go status for an area.
type AreaStatus string

const (
	StatusGo   AreaStatus = "GO"
	StatusNoGo AreaStatus = "NO-GO"
	StatusWarn AreaStatus = "WARN"
	StatusSkip AreaStatus = "SKIP"
)

// AreaIcon returns the UTF-8 icon for the area status.
func (s AreaStatus) Icon() string {
	switch s {
	case StatusGo:
		return IconGo
	case StatusNoGo:
		return IconNoGo
	case StatusWarn:
		return IconWarning
	case StatusSkip:
		return IconSkipped
	default:
		return "?"
	}
}

// ValidationReport contains all area results for a release validation.
type ValidationReport struct {
	Version string
	Areas   []AreaResult
}

// IsGo returns true if all areas pass validation.
func (r *ValidationReport) IsGo() bool {
	for _, area := range r.Areas {
		if area.Status == StatusNoGo {
			return false
		}
	}
	return true
}

// ComputeAreaStatus computes the status for an area based on its results.
func ComputeAreaStatus(results []Result) AreaStatus {
	hasNoGo := false
	hasWarn := false
	allSkipped := true

	for _, r := range results {
		if !r.Skipped {
			allSkipped = false
		}
		if r.Skipped {
			continue
		}
		if r.Warning && !r.Passed {
			hasWarn = true
			continue
		}
		if !r.Passed {
			hasNoGo = true
		}
	}

	if allSkipped {
		return StatusSkip
	}
	if hasNoGo {
		return StatusNoGo
	}
	if hasWarn {
		return StatusWarn
	}
	return StatusGo
}

// PrintValidationReport prints a comprehensive Go/No-Go report organized by area.
func PrintValidationReport(report *ValidationReport) {
	fmt.Println()
	fmt.Println("╔══════════════════════════════════════════════════════════════════════════════╗")
	if report.Version != "" {
		fmt.Printf("║             RELEASE VALIDATION: %-42s ║\n", report.Version)
	} else {
		fmt.Println("║                        RELEASE VALIDATION                                    ║")
	}
	fmt.Println("╠══════════════════════════════════════════════════════════════════════════════╣")
	fmt.Println("║  Assumes: Engineering ✅ SIGNED OFF | Product ✅ SIGNED OFF                  ║")
	fmt.Println("╠══════════════════════════════════════════════════════════════════════════════╣")

	for _, area := range report.Areas {
		// Area header
		// Box is 78 chars. Emoji displays as 2 but counts as 1, so reduce padding by 1
		icon := area.Status.Icon()
		fmt.Printf("║ %s %-8s %-62s ║\n", icon, area.Status, string(area.Area))
		fmt.Println("║ ────────────────────────────────────────────────────────────────────────── ║")

		// Individual checks
		for _, r := range area.Results {
			checkIcon := IconGo
			checkStatus := "GO"

			if r.Skipped {
				checkIcon = IconSkipped
				checkStatus = "SKIP"
			} else if r.Warning && !r.Passed {
				checkIcon = IconWarning
				checkStatus = "WARN"
			} else if !r.Passed {
				checkIcon = IconNoGo
				checkStatus = "NO-GO"
			}

			// Truncate name if too long
			name := r.Name
			if len(name) > 52 {
				name = name[:49] + "..."
			}
			// Emoji displays as 2 but counts as 1, so reduce padding by 1
			fmt.Printf("║   %s %-6s %-58s ║\n", checkIcon, checkStatus, name)
		}
		fmt.Println("║                                                                              ║")
	}

	fmt.Println("╠══════════════════════════════════════════════════════════════════════════════╣")

	// Final verdict
	if report.IsGo() {
		fmt.Println("║                         🚀 ALL SYSTEMS GO 🚀                               ║")
		fmt.Println("║                      RELEASE VALIDATION: APPROVED                          ║")
	} else {
		fmt.Println("║                        🛑 NO-GO FOR RELEASE 🛑                              ║")
		fmt.Println("║                    RELEASE VALIDATION: NOT APPROVED                        ║")
	}
	fmt.Println("╚══════════════════════════════════════════════════════════════════════════════╝")
	fmt.Println()
}
