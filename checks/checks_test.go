package checks

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultOptions(t *testing.T) {
	opts := DefaultOptions()

	if !opts.Test {
		t.Error("expected Test to be true by default")
	}
	if !opts.Lint {
		t.Error("expected Lint to be true by default")
	}
	if !opts.Format {
		t.Error("expected Format to be true by default")
	}
	if opts.Coverage {
		t.Error("expected Coverage to be false by default")
	}
	if opts.GoExcludeCoverage != "cmd" {
		t.Errorf("expected GoExcludeCoverage to be 'cmd', got %s", opts.GoExcludeCoverage)
	}
}

func TestRunCommand_Success(t *testing.T) {
	result := RunCommand("test", ".", "echo", "hello")

	if !result.Passed {
		t.Error("expected command to pass")
	}
	if result.Output != "hello" {
		t.Errorf("expected output 'hello', got %q", result.Output)
	}
	if result.Error != nil {
		t.Errorf("expected no error, got %v", result.Error)
	}
}

func TestRunCommand_Failure(t *testing.T) {
	result := RunCommand("test", ".", "false")

	if result.Passed {
		t.Error("expected command to fail")
	}
	if result.Error == nil {
		t.Error("expected error")
	}
}

func TestRunCommand_NotFound(t *testing.T) {
	result := RunCommand("test", ".", "nonexistent-command-12345")

	if result.Passed {
		t.Error("expected command to fail")
	}
	if result.Error == nil {
		t.Error("expected error for non-existent command")
	}
}

func TestCommandExists(t *testing.T) {
	// echo should exist on all systems
	if !CommandExists("echo") {
		t.Error("expected 'echo' command to exist")
	}

	// This should not exist
	if CommandExists("nonexistent-command-12345") {
		t.Error("expected fake command to not exist")
	}
}

func TestFileExists(t *testing.T) {
	dir := t.TempDir()
	testFile := filepath.Join(dir, "test.txt")

	// File doesn't exist yet
	if FileExists(testFile) {
		t.Error("expected file to not exist")
	}

	// Create file
	if err := os.WriteFile(testFile, []byte("test"), 0600); err != nil {
		t.Fatal(err)
	}

	// Now it should exist
	if !FileExists(testFile) {
		t.Error("expected file to exist")
	}
}

func TestPrintResults(t *testing.T) {
	results := []Result{
		{Name: "test1", Passed: true},
		{Name: "test2", Passed: false, Output: "failed"},
		{Name: "test3", Skipped: true, Reason: "not configured"},
	}

	passed, failed, skipped := PrintResults(results, false)

	if passed != 1 {
		t.Errorf("expected 1 passed, got %d", passed)
	}
	if failed != 1 {
		t.Errorf("expected 1 failed, got %d", failed)
	}
	if skipped != 1 {
		t.Errorf("expected 1 skipped, got %d", skipped)
	}
}
