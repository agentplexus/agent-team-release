package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/agentplexus/agent-team-release/pkg/checks"
	"github.com/agentplexus/agent-team-release/pkg/config"
	"github.com/agentplexus/agent-team-release/pkg/detect"
	"github.com/agentplexus/assistantkit/requirements"
)

// Check command flags
var (
	noTest     bool
	noLint     bool
	noFormat   bool
	coverage   bool
	goNoGoMode bool
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check [directory]",
	Short: "Run validation checks",
	Long: `Run validation checks for all detected languages in the repository.

Checks include build, test, lint, and format verification for each
detected language. Results are summarized with pass/fail status.

Examples:
  atrelease check              # Check current directory
  atrelease check /path/to/repo
  atrelease check --verbose    # Show detailed output
  atrelease check --no-test    # Skip tests`,
	Args: cobra.MaximumNArgs(1),
	Run:  runCheck,
}

func init() {
	checkCmd.Flags().BoolVar(&noTest, "no-test", false, "Skip tests")
	checkCmd.Flags().BoolVar(&noLint, "no-lint", false, "Skip linting")
	checkCmd.Flags().BoolVar(&noFormat, "no-format", false, "Skip format checks")
	checkCmd.Flags().BoolVar(&coverage, "coverage", false, "Show coverage (Go only)")
	checkCmd.Flags().BoolVar(&goNoGoMode, "go-no-go", false, "Display NASA-style Go/No-Go validation report")

	rootCmd.AddCommand(checkCmd)
}

func runCheck(cmd *cobra.Command, args []string) {
	// Get directory to check
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

	// Check if releasekit is available, prompt for installation if not
	prompter := requirements.NewCLIPrompter()
	result := requirements.EnsureRequirements([]string{"releasekit"}, prompter)
	if !result.AllSatisfied() {
		fmt.Fprintf(os.Stderr, "Cannot proceed without required tools\n")
		fmt.Fprint(os.Stderr, requirements.FormatMissingError(result))
		os.Exit(1)
	}

	// Detect languages
	fmt.Println("=== Pre-push Checks ===")
	fmt.Println()
	fmt.Println("Detecting languages...")

	detections, err := detect.Detect(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error detecting languages: %v\n", err)
		os.Exit(1)
	}

	if len(detections) == 0 {
		fmt.Println("No supported languages detected.")
		os.Exit(0)
	}

	// Print detected languages
	for _, d := range detections {
		fmt.Printf("  Found: %s in %s\n", d.Language, d.Path)
	}
	fmt.Println()

	// Build options from flags and config
	opts := checks.Options{
		Test:    !noTest,
		Lint:    !noLint,
		Format:  !noFormat,
		Coverage: coverage,
		Verbose: cfg.Verbose,
	}

	// Run releasekit validate (auto-detects languages)
	fmt.Println("Running checks via releasekit...")
	allResults, err := checks.RunReleasekit(dir, opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running releasekit: %v\n", err)
		os.Exit(1)
	}
	fmt.Println()

	// Print summary
	if goNoGoMode {
		// NASA-style Go/No-Go report
		allGo := checks.PrintGoNoGoReport(allResults, cfg.Verbose)
		if !allGo {
			os.Exit(1)
		}
	} else {
		// Standard report
		fmt.Println("=== Summary ===")
		passed, failed, skipped, warnings := checks.PrintResults(allResults, cfg.Verbose)
		fmt.Println()
		if warnings > 0 {
			fmt.Printf("Passed: %d, Failed: %d, Skipped: %d, Warnings: %d\n", passed, failed, skipped, warnings)
		} else {
			fmt.Printf("Passed: %d, Failed: %d, Skipped: %d\n", passed, failed, skipped)
		}

		if failed > 0 {
			fmt.Println()
			fmt.Println("Pre-push checks failed!")
			os.Exit(1)
		}

		fmt.Println()
		if warnings > 0 {
			fmt.Println("Pre-push checks passed with warnings.")
		} else {
			fmt.Println("All pre-push checks passed!")
		}
	}
}
