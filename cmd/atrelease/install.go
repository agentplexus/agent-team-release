package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/agentplexus/agent-team-release/plugins/kiro"
	"github.com/spf13/cobra"
)

const (
	// DefaultInstallPrefix is the default prefix for installed files.
	// This is compiled into the binary so users don't need to specify it.
	DefaultInstallPrefix = "agent-team-release"
)

var (
	installApply  bool
	installPrefix string
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install agent configurations to local directories",
	Long: `Install pre-built agent configurations for various AI assistant platforms.

Supported platforms:
  kiro    AWS Kiro CLI agents and steering files

By default, shows a plan of what would be installed. Use --apply to install.`,
}

var installKiroCmd = &cobra.Command{
	Use:   "kiro",
	Short: "Install Kiro CLI agents and steering files",
	Long: `Install pre-built Kiro CLI agent configurations to ~/.kiro/

This installs:
  - Agent definitions to ~/.kiro/agents/
  - Steering files to ~/.kiro/steering/

Files are prefixed with the team name to avoid collisions when installing
agents from multiple projects. Default prefix: agent-team-release

Example installed files:
  ~/.kiro/agents/agent-team-release_pm.json
  ~/.kiro/steering/agent-team-release_version-analysis.md

By default, shows a plan of what would be installed. Use --apply to install.`,
	RunE: runInstallKiro,
}

func init() {
	installCmd.AddCommand(installKiroCmd)
	installKiroCmd.Flags().BoolVar(&installApply, "apply", false, "Apply the installation (default: plan only)")
	installKiroCmd.Flags().StringVar(&installPrefix, "prefix", DefaultInstallPrefix, "Prefix for installed files")
	rootCmd.AddCommand(installCmd)
}

// FileAction represents an install action
type FileAction struct {
	Action string // "create", "update", "unchanged"
	Source string
	Dest   string
	Size   int64
}

func runInstallKiro(cmd *cobra.Command, args []string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	kiroDir := filepath.Join(homeDir, ".kiro")
	agentsDir := filepath.Join(kiroDir, "agents")
	steeringDir := filepath.Join(kiroDir, "steering")

	var actions []FileAction

	// Collect agent files
	agentActions, err := planEmbeddedFiles(kiro.AgentFiles, "agents", agentsDir, ".json", installPrefix)
	if err != nil {
		return err
	}
	actions = append(actions, agentActions...)

	// Collect steering files
	steeringActions, err := planEmbeddedFiles(kiro.SteeringFiles, "steering", steeringDir, ".md", installPrefix)
	if err != nil {
		return err
	}
	actions = append(actions, steeringActions...)

	if len(actions) == 0 {
		fmt.Println("No files to install.")
		return nil
	}

	// Display plan
	fmt.Println()
	fmt.Println("Kiro Agent Installation Plan:")
	fmt.Println()

	var toCreate, toUpdate, unchanged int
	for _, action := range actions {
		symbol := "+"
		color := "\033[32m" // green
		switch action.Action {
		case "update":
			symbol = "~"
			color = "\033[33m" // yellow
			toUpdate++
		case "unchanged":
			symbol = " "
			color = "\033[90m" // gray
			unchanged++
		default:
			toCreate++
		}
		reset := "\033[0m"
		fmt.Printf("  %s%s %s%s\n", color, symbol, action.Dest, reset)
	}

	fmt.Println()
	fmt.Printf("%d to create, %d to update, %d unchanged\n", toCreate, toUpdate, unchanged)

	if !installApply {
		fmt.Println()
		fmt.Println("Run with --apply to install.")
		return nil
	}

	// Apply installation
	fmt.Println()
	fmt.Println("Installing...")

	// Create directories
	if err := os.MkdirAll(agentsDir, 0755); err != nil {
		return fmt.Errorf("failed to create agents directory: %w", err)
	}
	if err := os.MkdirAll(steeringDir, 0755); err != nil {
		return fmt.Errorf("failed to create steering directory: %w", err)
	}

	// Copy files
	installed := 0
	for _, action := range actions {
		if action.Action == "unchanged" {
			continue
		}

		var data []byte
		var readErr error
		if strings.HasSuffix(action.Source, ".json") {
			data, readErr = kiro.AgentFiles.ReadFile(action.Source)
			if readErr == nil && installPrefix != "" {
				// Prefix the "name" field inside the agent JSON
				data, readErr = prefixAgentName(data, installPrefix)
			}
		} else {
			data, readErr = kiro.SteeringFiles.ReadFile(action.Source)
		}
		if readErr != nil {
			return fmt.Errorf("failed to read %s: %w", action.Source, readErr)
		}

		if err := os.WriteFile(action.Dest, data, 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", action.Dest, err)
		}
		installed++
	}

	fmt.Printf("\nInstalled %d files to %s\n", installed, kiroDir)
	return nil
}

func planEmbeddedFiles(fsys fs.FS, srcDir, destDir, ext, prefix string) ([]FileAction, error) {
	var actions []FileAction

	err := fs.WalkDir(fsys, srcDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ext) {
			return nil
		}

		filename := filepath.Base(path)
		// Add prefix to filename if provided
		if prefix != "" {
			filename = prefix + "_" + filename
		}
		destPath := filepath.Join(destDir, filename)

		// Read source file
		srcData, err := fs.ReadFile(fsys, path)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", path, err)
		}

		// For JSON files, apply prefix to name field for accurate comparison
		if ext == ".json" && prefix != "" {
			srcData, err = prefixAgentName(srcData, prefix)
			if err != nil {
				return fmt.Errorf("failed to prefix agent name in %s: %w", path, err)
			}
		}

		action := FileAction{
			Source: path,
			Dest:   destPath,
			Size:   int64(len(srcData)),
		}

		// Check if destination exists
		destData, err := os.ReadFile(destPath)
		if os.IsNotExist(err) {
			action.Action = "create"
		} else if err != nil {
			return fmt.Errorf("failed to read %s: %w", destPath, err)
		} else if string(srcData) == string(destData) {
			action.Action = "unchanged"
		} else {
			action.Action = "update"
		}

		actions = append(actions, action)
		return nil
	})

	return actions, err
}

// prefixAgentName modifies the "name" field in a Kiro agent JSON to include the prefix.
func prefixAgentName(data []byte, prefix string) ([]byte, error) {
	var agent map[string]interface{}
	if err := json.Unmarshal(data, &agent); err != nil {
		return nil, err
	}

	if name, ok := agent["name"].(string); ok {
		agent["name"] = prefix + "_" + name
	}

	return json.MarshalIndent(agent, "", "  ")
}
