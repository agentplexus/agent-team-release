package main

import (
	"os"

	"github.com/spf13/cobra"
)

// Version information (set via ldflags)
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// OutputFormat defines the structured output format.
type OutputFormat string

const (
	// OutputFormatTOON is the default Token-Oriented Object Notation format.
	OutputFormatTOON OutputFormat = "toon"
	// OutputFormatJSON is standard JSON format.
	OutputFormatJSON OutputFormat = "json"
)

// Global flags
var (
	cfgVerbose     bool
	cfgInteractive bool
	cfgJSON        bool   // Enable structured output (TOON by default)
	cfgFormat      string // Output format: "toon" or "json"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "atrelease",
	Short: "Autonomous release preparation agent",
	Long: `Release Agent validates code quality, generates changelogs, updates documentation,
and manages the complete release lifecycle for multi-language repositories.

It supports Go, TypeScript, JavaScript, Python, Rust, and Swift with automatic
language detection and monorepo support.`,
	// Run check by default when no subcommand is provided
	Run: func(cmd *cobra.Command, args []string) {
		// Default to running check command
		checkCmd.Run(cmd, args)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Global flags available to all subcommands
	rootCmd.PersistentFlags().BoolVarP(&cfgVerbose, "verbose", "v", false, "Show detailed output")
	rootCmd.PersistentFlags().BoolVarP(&cfgInteractive, "interactive", "i", false, "Enable interactive mode")
	rootCmd.PersistentFlags().BoolVar(&cfgJSON, "json", false, "Enable structured output for LLM integration (TOON format by default)")
	rootCmd.PersistentFlags().StringVar(&cfgFormat, "format", "toon", "Output format when --json is enabled: toon (default) or json")

	// Add subcommands
	rootCmd.AddCommand(checkCmd)
	rootCmd.AddCommand(versionCmd)
}

// GetOutputFormat returns the configured output format.
func GetOutputFormat() OutputFormat {
	if cfgFormat == "json" {
		return OutputFormatJSON
	}
	return OutputFormatTOON
}
