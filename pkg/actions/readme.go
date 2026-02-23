package actions

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// ReadmeAction updates README badges and version references.
type ReadmeAction struct{}

// Name returns the action name.
func (a *ReadmeAction) Name() string {
	return "readme"
}

// Run executes the readme action directly.
func (a *ReadmeAction) Run(dir string, opts Options) Result {
	var output strings.Builder
	var changes []string

	readmePath := filepath.Join(dir, "README.md")
	if !fileExists(readmePath) {
		return Result{
			Name:    "readme",
			Success: false,
			Error:   fmt.Errorf("README.md not found"),
			Output:  "README.md not found in " + dir,
		}
	}

	content, err := os.ReadFile(readmePath)
	if err != nil {
		return Result{
			Name:    "readme",
			Success: false,
			Error:   err,
			Output:  "Failed to read README.md",
		}
	}

	newContent := string(content)

	// Update version references if version is specified
	if opts.Version != "" {
		fmt.Fprintf(&output, "Updating version references to %s...\n", opts.Version)

		// Update @latest or @vX.Y.Z in go install commands
		goInstallRegex := regexp.MustCompile(`go install ([^@]+)@v?[\d.]+`)
		if goInstallRegex.MatchString(newContent) {
			newContent = goInstallRegex.ReplaceAllString(newContent, "go install $1@"+opts.Version)
			changes = append(changes, "Updated go install version")
		}

		// Update version in badges (e.g., version-vX.Y.Z-blue)
		versionBadgeRegex := regexp.MustCompile(`version-v[\d.]+-`)
		if versionBadgeRegex.MatchString(newContent) {
			newContent = versionBadgeRegex.ReplaceAllString(newContent, "version-"+opts.Version+"-")
			changes = append(changes, "Updated version badge")
		}
	}

	// Update coverage badge if gocoverbadge is available
	if commandExists("gocoverbadge") {
		output.WriteString("Updating coverage badge...\n")

		// Run gocoverbadge to generate badge
		excludeArg := ""
		if cfg := opts.Config; cfg != nil {
			langCfg := cfg.GetLanguageConfig("go")
			if langCfg.ExcludeCoverage != "" {
				excludeArg = langCfg.ExcludeCoverage
			}
		}

		args := []string{"-dir", dir, "-badge-only"}
		if excludeArg != "" {
			args = append(args, "-exclude", excludeArg)
		}

		result := runCommand("gocoverbadge", dir, "gocoverbadge", args...)
		if result.Success {
			changes = append(changes, "Updated coverage badge")
			fmt.Fprintf(&output, "%s\n", result.Output)
		} else {
			fmt.Fprintf(&output, "Warning: gocoverbadge failed: %s\n", result.Output)
		}
	}

	// If no changes, report that
	if len(changes) == 0 {
		output.WriteString("No changes needed in README.md\n")
		return Result{
			Name:    "readme",
			Success: true,
			Output:  output.String(),
		}
	}

	// If dry run, don't write
	if opts.DryRun {
		output.WriteString("\n[Dry run] Would make these changes:\n")
		for _, change := range changes {
			fmt.Fprintf(&output, "  - %s\n", change)
		}
		return Result{
			Name:    "readme",
			Success: true,
			Output:  output.String(),
		}
	}

	// Write updated content
	if newContent != string(content) {
		err = os.WriteFile(readmePath, []byte(newContent), 0644)
		if err != nil {
			return Result{
				Name:    "readme",
				Success: false,
				Error:   err,
				Output:  "Failed to write README.md",
			}
		}
		output.WriteString("\nUpdated README.md:\n")
		for _, change := range changes {
			fmt.Fprintf(&output, "  - %s\n", change)
		}
	}

	return Result{
		Name:    "readme",
		Success: true,
		Output:  output.String(),
	}
}

// Propose generates proposals for interactive mode.
func (a *ReadmeAction) Propose(dir string, opts Options) ([]Proposal, error) {
	readmePath := filepath.Join(dir, "README.md")
	if !fileExists(readmePath) {
		return nil, fmt.Errorf("README.md not found")
	}

	content, err := os.ReadFile(readmePath)
	if err != nil {
		return nil, err
	}

	oldContent := string(content)
	newContent := oldContent

	var description strings.Builder
	description.WriteString("Update README.md")

	// Preview version changes
	if opts.Version != "" {
		goInstallRegex := regexp.MustCompile(`go install ([^@]+)@v?[\d.]+`)
		if goInstallRegex.MatchString(newContent) {
			newContent = goInstallRegex.ReplaceAllString(newContent, "go install $1@"+opts.Version)
			fmt.Fprintf(&description, "\n  - Update go install version to %s", opts.Version)
		}

		versionBadgeRegex := regexp.MustCompile(`version-v[\d.]+-`)
		if versionBadgeRegex.MatchString(newContent) {
			newContent = versionBadgeRegex.ReplaceAllString(newContent, "version-"+opts.Version+"-")
			fmt.Fprintf(&description, "\n  - Update version badge to %s", opts.Version)
		}
	}

	if commandExists("gocoverbadge") {
		description.WriteString("\n  - Update coverage badge")
	}

	if newContent == oldContent && !commandExists("gocoverbadge") {
		return nil, fmt.Errorf("no changes to propose")
	}

	return []Proposal{
		{
			Description: description.String(),
			FilePath:    "README.md",
			OldContent:  oldContent,
			NewContent:  newContent,
			Metadata: map[string]string{
				"version": opts.Version,
			},
		},
	}, nil
}

// Apply applies approved proposals.
func (a *ReadmeAction) Apply(dir string, proposals []Proposal) Result {
	if len(proposals) == 0 {
		return Result{
			Name:    "readme",
			Success: true,
			Output:  "No proposals to apply",
		}
	}

	proposal := proposals[0]
	readmePath := filepath.Join(dir, proposal.FilePath)

	err := os.WriteFile(readmePath, []byte(proposal.NewContent), 0644)
	if err != nil {
		return Result{
			Name:    "readme",
			Success: false,
			Error:   err,
			Output:  "Failed to write README.md",
		}
	}

	return Result{
		Name:    "readme",
		Success: true,
		Output:  "Updated README.md",
	}
}
