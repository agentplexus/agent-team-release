package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/agentplexus/agent-team-release/pkg/actions"
	"github.com/agentplexus/agent-team-release/pkg/config"
)

// Readme command flags
var (
	readmeVersion string
	readmeDryRun  bool
)

// readmeCmd represents the readme command
var readmeCmd = &cobra.Command{
	Use:   "readme [directory]",
	Short: "Update README badges and version references",
	Long: `Update README.md with new version references and badges.

This command can update:
  - go install version references
  - Version badges
  - Coverage badges (if gocoverbadge is installed)

Examples:
  atrelease readme --version=v0.3.0    # Update version references
  atrelease readme --dry-run           # Show what would change`,
	Args: cobra.MaximumNArgs(1),
	Run:  runReadme,
}

func init() {
	readmeCmd.Flags().StringVar(&readmeVersion, "version", "", "Version to update references to")
	readmeCmd.Flags().BoolVar(&readmeDryRun, "dry-run", false, "Show what would be done without making changes")

	rootCmd.AddCommand(readmeCmd)
}

func runReadme(cmd *cobra.Command, args []string) {
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

	// Load configuration
	cfg, _ := config.Load(dir)

	fmt.Println("=== README ===")
	fmt.Println()

	action := &actions.ReadmeAction{}
	opts := actions.Options{
		Version: readmeVersion,
		DryRun:  readmeDryRun,
		Verbose: cfgVerbose,
		Config:  &cfg,
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
	fmt.Println("README action completed successfully.")
}
