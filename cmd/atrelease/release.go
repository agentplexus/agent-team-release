package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/toon-format/toon-go"

	"github.com/agentplexus/agent-team-release/pkg/workflow"
)

// Release command flags
var (
	releaseDryRun     bool
	releaseSkipChecks bool
	releaseSkipCI     bool
)

// releaseCmd represents the release command
var releaseCmd = &cobra.Command{
	Use:   "release <version>",
	Short: "Create a release",
	Long: `Execute the full release workflow for the specified version.

The release workflow includes:
  1. Validate version format and check it doesn't exist
  2. Ensure working directory is clean
  3. Run validation checks (build, test, lint, format)
  4. Generate/update changelog
  5. Update roadmap
  6. Create release commit
  7. Push to remote
  8. Wait for CI to pass
  9. Create and push release tag

Examples:
  atrelease release v0.3.0
  atrelease release v0.3.0 --dry-run     # Preview without changes
  atrelease release v0.3.0 --skip-ci     # Don't wait for CI
  atrelease release v0.3.0 --skip-checks # Skip validation`,
	Args: cobra.ExactArgs(1),
	Run:  runRelease,
}

func init() {
	releaseCmd.Flags().BoolVar(&releaseDryRun, "dry-run", false, "Preview what would be done without making changes")
	releaseCmd.Flags().BoolVar(&releaseSkipChecks, "skip-checks", false, "Skip validation checks (dangerous)")
	releaseCmd.Flags().BoolVar(&releaseSkipCI, "skip-ci", false, "Don't wait for CI to pass before tagging")

	rootCmd.AddCommand(releaseCmd)
}

func runRelease(cmd *cobra.Command, args []string) {
	version := args[0]

	// Get directory
	dir := "."

	// Make sure directory exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: directory %s does not exist\n", dir)
		os.Exit(1)
	}

	// Create workflow context
	ctx := workflow.NewContext(dir, version)
	ctx.SkipChecks = releaseSkipChecks
	ctx.SkipCI = releaseSkipCI

	// Create runner
	runner := workflow.NewRunner()
	runner.DryRun = releaseDryRun
	runner.Verbose = cfgVerbose
	runner.Interactive = cfgInteractive
	runner.JSONOutput = cfgJSON

	// Create and run the release workflow
	wf := workflow.ReleaseWorkflow(version)
	result := runner.Run(wf, ctx)

	// Print output
	if cfgJSON {
		// Output structured result (TOON or JSON based on format flag)
		jsonResult := result.ToJSON()
		if GetOutputFormat() == OutputFormatJSON {
			// JSON format
			encoder := json.NewEncoder(os.Stdout)
			encoder.SetIndent("", "  ")
			if err := encoder.Encode(jsonResult); err != nil {
				fmt.Fprintf(os.Stderr, "Error encoding JSON: %v\n", err)
				os.Exit(1)
			}
		} else {
			// TOON format (default)
			data, err := toon.Marshal(jsonResult, toon.WithIndent(2))
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error encoding TOON: %v\n", err)
				os.Exit(1)
			}
			fmt.Print(string(data))
		}
	} else {
		fmt.Print(result.Output)

		// Print summary
		if cfgVerbose {
			fmt.Println()
			fmt.Print(result.Summary())
		}
	}

	if !result.Success {
		os.Exit(1)
	}
}
