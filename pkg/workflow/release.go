package workflow

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/agentplexus/agent-team-release/pkg/actions"
	"github.com/agentplexus/agent-team-release/pkg/checks"
	"github.com/agentplexus/agent-team-release/pkg/detect"
	"github.com/agentplexus/agent-team-release/pkg/git"
	"github.com/agentplexus/assistantkit/requirements"
)

// ReleaseWorkflow creates a workflow for releasing a new version.
func ReleaseWorkflow(version string) *Workflow {
	return &Workflow{
		Name:        "Release " + version,
		Description: "Prepare and create release " + version,
		Steps: []Step{
			{
				Name:        "Validate version",
				Description: "Check version format and ensure it doesn't exist",
				Type:        StepTypeFunc,
				Required:    true,
				Func:        validateVersion,
			},
			{
				Name:        "Check working directory",
				Description: "Ensure no uncommitted changes",
				Type:        StepTypeFunc,
				Required:    true,
				Func:        checkWorkingDirectory,
			},
			{
				Name:        "Run validation checks",
				Description: "Run build, test, lint, format checks",
				Type:        StepTypeFunc,
				Required:    true,
				Func:        runValidationChecks,
			},
			{
				Name:        "Generate changelog",
				Description: "Update CHANGELOG.md with new entries",
				Type:        StepTypeFunc,
				Required:    false,
				Func:        generateChangelog,
			},
			{
				Name:        "Update roadmap",
				Description: "Regenerate ROADMAP.md",
				Type:        StepTypeFunc,
				Required:    false,
				Func:        updateRoadmap,
			},
			{
				Name:        "Create release commit",
				Description: "Commit all changes with release message",
				Type:        StepTypeFunc,
				Required:    true,
				Func:        createReleaseCommit,
			},
			{
				Name:        "Push to remote",
				Description: "Push commits to origin",
				Type:        StepTypeFunc,
				Required:    true,
				Func:        pushToRemote,
			},
			{
				Name:        "Wait for CI",
				Description: "Wait for CI checks to pass",
				Type:        StepTypeFunc,
				Required:    false,
				Func:        waitForCI,
			},
			{
				Name:        "Create tag",
				Description: "Create and push release tag",
				Type:        StepTypeFunc,
				Required:    true,
				Func:        createTag,
			},
		},
	}
}

// validateVersion checks that the version is valid and doesn't already exist.
func validateVersion(ctx *Context) error {
	if ctx.Version == "" {
		return fmt.Errorf("version is required")
	}

	// Check version format (should start with v)
	if ctx.Version[0] != 'v' {
		ctx.Version = "v" + ctx.Version
	}

	// Check if tag already exists
	g := git.New(ctx.Dir)
	tags, err := g.AllTags()
	if err == nil {
		for _, tag := range tags {
			if tag == ctx.Version {
				return fmt.Errorf("tag %s already exists", ctx.Version)
			}
		}
	}

	ctx.Log("  Version: %s", ctx.Version)
	return nil
}

// checkWorkingDirectory ensures there are no uncommitted changes.
func checkWorkingDirectory(ctx *Context) error {
	g := git.New(ctx.Dir)

	dirty, err := g.IsDirty()
	if err != nil {
		return fmt.Errorf("failed to check git status: %w", err)
	}

	if dirty {
		// In dry-run mode, just warn
		if ctx.DryRun {
			ctx.Log("  Warning: working directory has uncommitted changes")
			return nil
		}
		return fmt.Errorf("working directory has uncommitted changes; commit or stash them first")
	}

	ctx.Log("  Working directory is clean")
	return nil
}

// runValidationChecks runs all validation checks using releasekit CLI.
func runValidationChecks(ctx *Context) error {
	if ctx.SkipChecks {
		ctx.Log("  Skipping validation checks (--skip-checks)")
		return nil
	}

	// Check if releasekit is available, prompt for installation if not
	if !checks.ReleasekitAvailable() {
		prompter := requirements.NewCLIPrompter()
		reqResult := requirements.EnsureRequirements([]string{"releasekit"}, prompter)
		if !reqResult.AllSatisfied() {
			ctx.Log("  Warning: releasekit CLI not installed, skipping validation")
			return nil
		}
	}

	// Detect languages to see if there's anything to check
	detections, err := detect.Detect(ctx.Dir)
	if err != nil {
		return fmt.Errorf("failed to detect languages: %w", err)
	}

	if len(detections) == 0 {
		ctx.Log("  No supported languages detected, skipping checks")
		return nil
	}

	ctx.Log("  Running releasekit validate...")

	// Build options
	opts := checks.Options{
		Test:    true,
		Lint:    true,
		Format:  true,
		Verbose: ctx.Verbose,
	}

	// Run releasekit validate (it auto-detects languages)
	results, err := checks.RunReleasekit(ctx.Dir, opts)
	if err != nil {
		return fmt.Errorf("releasekit failed: %w", err)
	}

	// Count results
	failed := 0
	for _, r := range results {
		if !r.Passed && !r.Skipped && !r.Warning {
			failed++
			ctx.Log("    ✗ %s: %s", r.Name, r.Output)
		}
	}

	if failed > 0 {
		return fmt.Errorf("%d checks failed", failed)
	}

	ctx.Log("  All checks passed")
	return nil
}

// generateChangelog updates the changelog.
func generateChangelog(ctx *Context) error {
	action := &actions.ChangelogAction{}

	// Get latest tag for since
	g := git.New(ctx.Dir)
	since, _ := g.LatestTag()

	opts := actions.Options{
		Since:   since,
		Version: ctx.Version,
		DryRun:  ctx.DryRun,
		Verbose: ctx.Verbose,
	}

	result := action.Run(ctx.Dir, opts)
	if !result.Success {
		if result.Error != nil {
			ctx.Log("  Warning: %v", result.Error)
		}
		// Don't fail the workflow for changelog issues
		return nil
	}

	ctx.Log("  Changelog updated")
	return nil
}

// updateRoadmap regenerates the roadmap.
func updateRoadmap(ctx *Context) error {
	action := &actions.RoadmapAction{}

	opts := actions.Options{
		DryRun:  ctx.DryRun,
		Verbose: ctx.Verbose,
	}

	result := action.Run(ctx.Dir, opts)
	if !result.Success {
		if result.Error != nil {
			ctx.Log("  Warning: %v", result.Error)
		}
		// Don't fail the workflow for roadmap issues
		return nil
	}

	ctx.Log("  Roadmap updated")
	return nil
}

// createReleaseCommit commits all changes with a release message.
func createReleaseCommit(ctx *Context) error {
	g := git.New(ctx.Dir)

	// Check if there are changes to commit
	dirty, err := g.IsDirty()
	if err != nil {
		return err
	}

	if !dirty {
		ctx.Log("  No changes to commit")
		return nil
	}

	if ctx.DryRun {
		ctx.Log("  [Dry run] Would create commit: chore(release): %s", ctx.Version)
		return nil
	}

	message := fmt.Sprintf("chore(release): %s", ctx.Version)
	if err := g.CommitAll(message, false); err != nil {
		return fmt.Errorf("failed to create commit: %w", err)
	}

	ctx.Log("  Created commit: %s", message)
	return nil
}

// pushToRemote pushes commits to the remote.
func pushToRemote(ctx *Context) error {
	g := git.New(ctx.Dir)

	if ctx.DryRun {
		ctx.Log("  [Dry run] Would push to origin")
		return nil
	}

	// Check if we need to push
	status, err := g.Status()
	if err != nil {
		return err
	}

	if status.Ahead == 0 {
		ctx.Log("  Already up to date with remote")
		return nil
	}

	if err := g.Push(); err != nil {
		return fmt.Errorf("failed to push: %w", err)
	}

	ctx.Log("  Pushed to origin")
	return nil
}

// waitForCI waits for CI checks to pass.
func waitForCI(ctx *Context) error {
	if ctx.SkipCI {
		ctx.Log("  Skipping CI wait (--skip-ci)")
		return nil
	}

	g := git.New(ctx.Dir)

	// Check if gh CLI is available
	if !commandExists("gh") {
		ctx.Log("  gh CLI not found, skipping CI wait")
		return nil
	}

	if ctx.DryRun {
		ctx.Log("  [Dry run] Would wait for CI")
		return nil
	}

	ctx.Log("  Waiting for CI (timeout: 10 minutes)...")

	timeout := 10 * time.Minute
	if err := g.WaitForCI(timeout); err != nil {
		return fmt.Errorf("CI failed: %w", err)
	}

	ctx.Log("  CI passed")
	return nil
}

// createTag creates and pushes the release tag.
func createTag(ctx *Context) error {
	g := git.New(ctx.Dir)

	if ctx.DryRun {
		ctx.Log("  [Dry run] Would create tag: %s", ctx.Version)
		return nil
	}

	// Create the tag
	message := fmt.Sprintf("Release %s", ctx.Version)
	if err := g.CreateTag(ctx.Version, message, false); err != nil {
		return fmt.Errorf("failed to create tag: %w", err)
	}

	ctx.Log("  Created tag: %s", ctx.Version)

	// Push the tag
	if err := g.PushTag(ctx.Version); err != nil {
		// Try to clean up the local tag
		_ = g.DeleteTag(ctx.Version)
		return fmt.Errorf("failed to push tag: %w", err)
	}

	ctx.Log("  Pushed tag: %s", ctx.Version)
	return nil
}

// commandExists checks if a command is available.
func commandExists(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}
