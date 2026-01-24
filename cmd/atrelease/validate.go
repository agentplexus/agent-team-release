// Copyright 2025 John Wang. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	multiagentspec "github.com/agentplexus/multi-agent-spec/sdk/go"
	"github.com/spf13/cobra"

	"github.com/agentplexus/agent-team-release/pkg/checks"
	"github.com/agentplexus/agent-team-release/pkg/config"
	"github.com/agentplexus/agent-team-release/pkg/detect"
	"github.com/agentplexus/agent-team-release/pkg/report"
)

// Validate command flags
var (
	validateVersion  string
	validateSkipPM   bool
	validateSkipQA   bool
	validateSkipDocs bool
	validateSkipSec  bool
	validateFormat   string
)

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:   "validate [directory]",
	Short: "Run comprehensive release validation",
	Long: `Run comprehensive release validation across all areas of responsibility.

Validation Areas:
  PM            Version recommendation, release scope, changelog quality, breaking changes
  QA            Build, tests, lint, format, error handling compliance
  Documentation README, PRD, TRD, release notes, CHANGELOG
  Release       Version availability, git status, CI configuration
  Security      LICENSE, vulnerability scan, dependency audit

The PM agent runs first and produces the version recommendation. Other agents depend on PM.

Examples:
  atrelease validate                    # Validate current directory
  atrelease validate --version v0.2.0   # Include version-specific checks
  atrelease validate --skip-qa          # Skip QA checks
  atrelease validate --format team      # Team status report format
  atrelease validate -v                 # Verbose output`,
	Args: cobra.MaximumNArgs(1),
	Run:  runValidate,
}

func init() {
	validateCmd.Flags().StringVar(&validateVersion, "version", "", "Target release version (e.g., v0.2.0)")
	validateCmd.Flags().BoolVar(&validateSkipPM, "skip-pm", false, "Skip PM validation")
	validateCmd.Flags().BoolVar(&validateSkipQA, "skip-qa", false, "Skip QA checks")
	validateCmd.Flags().BoolVar(&validateSkipDocs, "skip-docs", false, "Skip documentation checks")
	validateCmd.Flags().BoolVar(&validateSkipSec, "skip-security", false, "Skip security checks")
	validateCmd.Flags().StringVar(&validateFormat, "format", "default", "Output format (default, team)")

	rootCmd.AddCommand(validateCmd)
}

func runValidate(cmd *cobra.Command, args []string) {
	// Get directory to validate
	dir := "."
	if len(args) > 0 {
		dir = args[0]
	}

	// Make sure directory exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: directory %s does not exist\n", dir)
		os.Exit(1)
	}

	// Load configuration
	cfg, err := config.Load(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: error loading config: %v\n", err)
	}

	// Override config with flags
	if cfgVerbose {
		cfg.Verbose = true
	}

	// Create validation report
	validationReport := &checks.ValidationReport{
		Version: validateVersion,
	}

	// Detect languages for QA checks
	detections, err := detect.Detect(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: error detecting languages: %v\n", err)
	}

	fmt.Println("╔════════════════════════════════════════════════════════════════════════╗")
	fmt.Println("║                     RELEASE VALIDATION STARTING                        ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════════════╝")
	fmt.Println()

	// PM Area (runs first - other agents depend on PM)
	if !validateSkipPM {
		fmt.Println("▶ Running PM validation...")
		pmChecker := &checks.PMChecker{}
		pmResults := pmChecker.Check(dir, checks.PMOptions{
			Version: validateVersion,
			Verbose: cfg.Verbose,
		})
		pmStatus := checks.ComputeAreaStatus(pmResults)
		validationReport.Areas = append(validationReport.Areas, checks.AreaResult{
			Area:    checks.AreaPM,
			Status:  pmStatus,
			Results: pmResults,
		})

		if pmStatus == checks.StatusNoGo {
			fmt.Println("  ⚠ PM validation failed - other agents will still run but release is blocked")
		}
	}

	// QA Area
	if !validateSkipQA {
		fmt.Println("▶ Running QA validation...")
		qaResults := runQAChecks(dir, detections, &cfg)
		validationReport.Areas = append(validationReport.Areas, checks.AreaResult{
			Area:    checks.AreaQA,
			Status:  checks.ComputeAreaStatus(qaResults),
			Results: qaResults,
		})
	}

	// Documentation Area
	if !validateSkipDocs {
		fmt.Println("▶ Running Documentation validation...")
		docChecker := &checks.DocChecker{}
		docResults := docChecker.Check(dir, checks.DocOptions{
			Version: validateVersion,
			Verbose: cfg.Verbose,
		})
		validationReport.Areas = append(validationReport.Areas, checks.AreaResult{
			Area:    checks.AreaDocumentation,
			Status:  checks.ComputeAreaStatus(docResults),
			Results: docResults,
		})
	}

	// Release Management Area
	fmt.Println("▶ Running Release Management validation...")
	releaseChecker := &checks.ReleaseChecker{}
	releaseResults := releaseChecker.Check(dir, checks.ReleaseOptions{
		Version: validateVersion,
		Verbose: cfg.Verbose,
	})
	validationReport.Areas = append(validationReport.Areas, checks.AreaResult{
		Area:    checks.AreaRelease,
		Status:  checks.ComputeAreaStatus(releaseResults),
		Results: releaseResults,
	})

	// Security Area
	if !validateSkipSec {
		fmt.Println("▶ Running Security validation...")
		secChecker := &checks.SecurityChecker{}
		secResults := secChecker.Check(dir, checks.SecurityOptions{
			Verbose: cfg.Verbose,
		})
		validationReport.Areas = append(validationReport.Areas, checks.AreaResult{
			Area:    checks.AreaSecurity,
			Status:  checks.ComputeAreaStatus(secResults),
			Results: secResults,
		})
	}

	// Print comprehensive report
	if validateFormat == "team" {
		printTeamStatusReport(validationReport, dir)
	} else {
		checks.PrintValidationReport(validationReport)
	}

	// Exit with error if validation failed
	if !validationReport.IsGo() {
		os.Exit(1)
	}
}

// printTeamStatusReport prints the validation report in team status format.
func printTeamStatusReport(vr *checks.ValidationReport, dir string) {
	// Determine project name from git remote
	project := getGitRemoteProject(dir)
	if project == "" {
		// Fall back to directory path
		if dir == "." {
			if cwd, err := os.Getwd(); err == nil {
				project = cwd
			}
		} else {
			project = dir
		}
	}

	// Build target string
	target := vr.Version
	if target == "" {
		target = "release validation"
	}

	// Try to load team spec for phase information
	phase := "PHASE 1: REVIEW"
	if spec, err := report.LoadTeamSpec(dir); err == nil {
		phases := report.GetPhases(spec)
		if len(phases) > 0 {
			phase = phases[0].Name
		}
	}

	// Convert to team status report (using multi-agent-spec types)
	teamReport := report.FromValidationReport(vr, project, target, phase)

	// Render the report using multi-agent-spec renderer
	renderer := multiagentspec.NewRenderer(os.Stdout)
	if err := renderer.Render(teamReport); err != nil {
		fmt.Fprintf(os.Stderr, "Error rendering report: %v\n", err)
	}
}

// getGitRemoteProject extracts the project path from git remote origin.
func getGitRemoteProject(dir string) string {
	// Try to get git remote URL using git command
	cmd := exec.Command("git", "-C", dir, "remote", "get-url", "origin")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	url := strings.TrimSpace(string(output))

	// Convert various URL formats to github.com/org/repo format
	// Handle: https://github.com/org/repo.git
	//         git@github.com:org/repo.git
	//         https://github.com/org/repo

	url = strings.TrimSuffix(url, ".git")

	if strings.HasPrefix(url, "https://") {
		return strings.TrimPrefix(url, "https://")
	}

	if strings.HasPrefix(url, "git@") {
		// git@github.com:org/repo -> github.com/org/repo
		url = strings.TrimPrefix(url, "git@")
		url = strings.Replace(url, ":", "/", 1)
		return url
	}

	return url
}

// runQAChecks runs all QA checks for detected languages.
func runQAChecks(dir string, detections []detect.Detection, cfg *config.Config) []checks.Result {
	var results []checks.Result

	// Go checks
	if detect.HasLanguage(detections, detect.Go) && cfg.IsLanguageEnabled("go") {
		langCfg := cfg.GetLanguageConfig("go")
		opts := checks.Options{
			Test:              *langCfg.Test,
			Lint:              *langCfg.Lint,
			Format:            *langCfg.Format,
			Coverage:          langCfg.Coverage != nil && *langCfg.Coverage,
			Verbose:           cfg.Verbose,
			GoExcludeCoverage: langCfg.ExcludeCoverage,
		}
		if opts.GoExcludeCoverage == "" {
			opts.GoExcludeCoverage = "cmd"
		}

		checker := &checks.GoChecker{}
		for _, d := range detect.GetByLanguage(detections, detect.Go) {
			results = append(results, checker.Check(d.Path, opts)...)
		}
	}

	// TypeScript checks
	if detect.HasLanguage(detections, detect.TypeScript) && cfg.IsLanguageEnabled("typescript") {
		langCfg := cfg.GetLanguageConfig("typescript")
		opts := checks.Options{
			Test:    *langCfg.Test,
			Lint:    *langCfg.Lint,
			Format:  *langCfg.Format,
			Verbose: cfg.Verbose,
		}

		checker := &checks.TypeScriptChecker{}
		for _, d := range detect.GetByLanguage(detections, detect.TypeScript) {
			results = append(results, checker.Check(d.Path, opts)...)
		}
	}

	// JavaScript checks
	if detect.HasLanguage(detections, detect.JavaScript) && cfg.IsLanguageEnabled("javascript") {
		langCfg := cfg.GetLanguageConfig("javascript")
		opts := checks.Options{
			Test:    *langCfg.Test,
			Lint:    *langCfg.Lint,
			Format:  *langCfg.Format,
			Verbose: cfg.Verbose,
		}

		checker := &checks.TypeScriptChecker{}
		for _, d := range detect.GetByLanguage(detections, detect.JavaScript) {
			results = append(results, checker.Check(d.Path, opts)...)
		}
	}

	return results
}
