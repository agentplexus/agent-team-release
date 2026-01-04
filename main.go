// Command prepush runs pre-push checks for multi-language repositories.
//
// It auto-detects languages based on project files (go.mod, package.json, etc.)
// and runs appropriate checks for each language found.
//
// Usage:
//
//	prepush              # Run in current directory
//	prepush /path/to/repo
//	prepush --verbose    # Show detailed output
//	prepush --no-test    # Skip tests
//
// Configuration:
//
// Create a .prepush.yaml file to customize behavior:
//
//	verbose: true
//	languages:
//	  go:
//	    lint: true
//	    test: true
//	    coverage: true
//	    exclude_coverage: "cmd"
//	  typescript:
//	    enabled: true
//	    paths: ["frontend/"]
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/grokify/prepush/checks"
	"github.com/grokify/prepush/config"
	"github.com/grokify/prepush/detect"
)

func main() {
	// Parse flags
	verbose := flag.Bool("verbose", false, "Show detailed output")
	noTest := flag.Bool("no-test", false, "Skip tests")
	noLint := flag.Bool("no-lint", false, "Skip linting")
	noFormat := flag.Bool("no-format", false, "Skip format checks")
	coverage := flag.Bool("coverage", false, "Show coverage (Go only)")
	flag.Parse()

	// Get directory to check
	dir := "."
	if flag.NArg() > 0 {
		dir = flag.Arg(0)
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
	if *verbose {
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
			Test:              !*noTest && *langCfg.Test,
			Lint:              !*noLint && *langCfg.Lint,
			Format:            !*noFormat && *langCfg.Format,
			Coverage:          *coverage || (langCfg.Coverage != nil && *langCfg.Coverage),
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
			Test:    !*noTest && *langCfg.Test,
			Lint:    !*noLint && *langCfg.Lint,
			Format:  !*noFormat && *langCfg.Format,
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
			Test:    !*noTest && *langCfg.Test,
			Lint:    !*noLint && *langCfg.Lint,
			Format:  !*noFormat && *langCfg.Format,
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
	fmt.Println("=== Summary ===")
	passed, failed, skipped := checks.PrintResults(allResults, cfg.Verbose)
	fmt.Println()
	fmt.Printf("Passed: %d, Failed: %d, Skipped: %d\n", passed, failed, skipped)

	if failed > 0 {
		fmt.Println()
		fmt.Println("Pre-push checks failed!")
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println("All pre-push checks passed!")
}
