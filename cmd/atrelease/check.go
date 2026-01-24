package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/agentplexus/agent-team-release/pkg/checks"
	"github.com/agentplexus/agent-team-release/pkg/config"
	"github.com/agentplexus/agent-team-release/pkg/detect"
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

	// Build checkers based on detections
	var allResults []checks.Result

	// Go checks
	if detect.HasLanguage(detections, detect.Go) && cfg.IsLanguageEnabled("go") {
		langCfg := cfg.GetLanguageConfig("go")
		opts := checks.Options{
			Test:              !noTest && *langCfg.Test,
			Lint:              !noLint && *langCfg.Lint,
			Format:            !noFormat && *langCfg.Format,
			Coverage:          coverage || (langCfg.Coverage != nil && *langCfg.Coverage),
			Verbose:           cfg.Verbose,
			GoExcludeCoverage: langCfg.ExcludeCoverage,
		}
		if opts.GoExcludeCoverage == "" {
			opts.GoExcludeCoverage = "cmd"
		}

		fmt.Println("Running Go checks...")
		checker := &checks.GoChecker{}
		for _, d := range detect.GetByLanguage(detections, detect.Go) {
			results := checker.Check(d.Path, opts)
			allResults = append(allResults, results...)
		}
		fmt.Println()
	}

	// TypeScript checks
	if detect.HasLanguage(detections, detect.TypeScript) && cfg.IsLanguageEnabled("typescript") {
		langCfg := cfg.GetLanguageConfig("typescript")
		opts := checks.Options{
			Test:    !noTest && *langCfg.Test,
			Lint:    !noLint && *langCfg.Lint,
			Format:  !noFormat && *langCfg.Format,
			Verbose: cfg.Verbose,
		}

		fmt.Println("Running TypeScript checks...")
		checker := &checks.TypeScriptChecker{}
		for _, d := range detect.GetByLanguage(detections, detect.TypeScript) {
			results := checker.Check(d.Path, opts)
			allResults = append(allResults, results...)
		}
		fmt.Println()
	}

	// JavaScript checks (use TypeScript checker)
	if detect.HasLanguage(detections, detect.JavaScript) && cfg.IsLanguageEnabled("javascript") {
		langCfg := cfg.GetLanguageConfig("javascript")
		opts := checks.Options{
			Test:    !noTest && *langCfg.Test,
			Lint:    !noLint && *langCfg.Lint,
			Format:  !noFormat && *langCfg.Format,
			Verbose: cfg.Verbose,
		}

		fmt.Println("Running JavaScript checks...")
		checker := &checks.TypeScriptChecker{}
		for _, d := range detect.GetByLanguage(detections, detect.JavaScript) {
			results := checker.Check(d.Path, opts)
			allResults = append(allResults, results...)
		}
		fmt.Println()
	}

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
