package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/agentplexus/agent-team-release/pkg/actions"
)

// Roadmap command flags
var (
	roadmapDryRun bool
)

// roadmapCmd represents the roadmap command
var roadmapCmd = &cobra.Command{
	Use:   "roadmap [directory]",
	Short: "Generate or update roadmap",
	Long: `Generate or update ROADMAP.md using sroadmap.

This command validates ROADMAP.json and regenerates ROADMAP.md
with deterministic formatting.

Requires sroadmap to be installed:
  go install github.com/grokify/sroadmap/cmd/sroadmap@latest

Examples:
  atrelease roadmap              # Regenerate ROADMAP.md
  atrelease roadmap --dry-run    # Show stats without generating`,
	Args: cobra.MaximumNArgs(1),
	Run:  runRoadmap,
}

func init() {
	roadmapCmd.Flags().BoolVar(&roadmapDryRun, "dry-run", false, "Show what would be done without making changes")

	rootCmd.AddCommand(roadmapCmd)
}

func runRoadmap(cmd *cobra.Command, args []string) {
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

	fmt.Println("=== Roadmap ===")
	fmt.Println()

	action := &actions.RoadmapAction{}
	opts := actions.Options{
		DryRun:  roadmapDryRun,
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
	fmt.Println("Roadmap action completed successfully.")
}
