package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/agentplexus/agent-team-release/pkg/actions"
)

// Changelog command flags
var (
	changelogSince  string
	changelogDryRun bool
)

// changelogCmd represents the changelog command
var changelogCmd = &cobra.Command{
	Use:   "changelog [directory]",
	Short: "Generate or update changelog",
	Long: `Generate or update CHANGELOG.md using schangelog.

This command parses git commits since the specified tag (or latest tag)
and can regenerate CHANGELOG.md from CHANGELOG.json.

Requires schangelog to be installed:
  go install github.com/grokify/schangelog/cmd/schangelog@latest

Examples:
  atrelease changelog                    # Parse commits since latest tag
  atrelease changelog --since=v0.2.0     # Parse commits since v0.2.0
  atrelease changelog --dry-run          # Show what would be done`,
	Args: cobra.MaximumNArgs(1),
	Run:  runChangelog,
}

func init() {
	changelogCmd.Flags().StringVar(&changelogSince, "since", "", "Parse commits since this tag (default: latest tag)")
	changelogCmd.Flags().BoolVar(&changelogDryRun, "dry-run", false, "Show what would be done without making changes")

	rootCmd.AddCommand(changelogCmd)
}

func runChangelog(cmd *cobra.Command, args []string) {
	// Get directory
	dir := "."
	if len(args) > 0 {
		dir = args[0]
	}

	// Make sure directory exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: directory %s does not exist\n", dir)
		os.Exit(1)
	}

	fmt.Println("=== Changelog ===")
	fmt.Println()

	action := &actions.ChangelogAction{}
	opts := actions.Options{
		Since:   changelogSince,
		DryRun:  changelogDryRun,
		Verbose: cfgVerbose,
	}

	result := action.Run(dir, opts)

	if result.Output != "" {
		fmt.Println(result.Output)
	}

	if !result.Success {
		if result.Error != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", result.Error)
		}
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println("Changelog action completed successfully.")
}
