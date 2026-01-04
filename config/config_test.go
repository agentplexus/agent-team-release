package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_NoConfig(t *testing.T) {
	dir := t.TempDir()

	cfg, err := Load(dir)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Should return defaults
	if cfg.Verbose {
		t.Error("expected Verbose to be false by default")
	}
}

func TestLoad_WithConfig(t *testing.T) {
	dir := t.TempDir()

	configContent := `
verbose: true
languages:
  go:
    enabled: true
    test: true
    lint: false
    coverage: true
    exclude_coverage: "cmd,internal"
  typescript:
    enabled: false
`
	if err := os.WriteFile(filepath.Join(dir, ".prepush.yaml"), []byte(configContent), 0600); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(dir)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if !cfg.Verbose {
		t.Error("expected Verbose to be true")
	}

	if !cfg.IsLanguageEnabled("go") {
		t.Error("expected Go to be enabled")
	}

	if cfg.IsLanguageEnabled("typescript") {
		t.Error("expected TypeScript to be disabled")
	}
}

func TestLoad_AlternateNames(t *testing.T) {
	names := []string{".prepush.yaml", ".prepush.yml", "prepush.yaml", "prepush.yml"}

	for _, name := range names {
		t.Run(name, func(t *testing.T) {
			dir := t.TempDir()
			configContent := "verbose: true\n"
			if err := os.WriteFile(filepath.Join(dir, name), []byte(configContent), 0600); err != nil {
				t.Fatal(err)
			}

			cfg, err := Load(dir)
			if err != nil {
				t.Fatalf("Load failed: %v", err)
			}

			if !cfg.Verbose {
				t.Errorf("expected Verbose to be true with config file %s", name)
			}
		})
	}
}

func TestIsLanguageEnabled(t *testing.T) {
	cfg := DefaultConfig()

	// Not configured = enabled (auto-detect)
	if !cfg.IsLanguageEnabled("go") {
		t.Error("expected unconfigured language to be enabled")
	}

	// Explicitly enabled
	enabled := true
	cfg.Languages["go"] = LanguageConfig{Enabled: &enabled}
	if !cfg.IsLanguageEnabled("go") {
		t.Error("expected explicitly enabled language to be enabled")
	}

	// Explicitly disabled
	disabled := false
	cfg.Languages["typescript"] = LanguageConfig{Enabled: &disabled}
	if cfg.IsLanguageEnabled("typescript") {
		t.Error("expected explicitly disabled language to be disabled")
	}

	// Nil enabled = auto-detect = enabled
	cfg.Languages["python"] = LanguageConfig{Enabled: nil}
	if !cfg.IsLanguageEnabled("python") {
		t.Error("expected nil-enabled language to be enabled")
	}
}

func TestGetLanguageConfig(t *testing.T) {
	cfg := DefaultConfig()

	// Unconfigured language gets defaults
	goCfg := cfg.GetLanguageConfig("go")
	if goCfg.Enabled == nil || !*goCfg.Enabled {
		t.Error("expected default Enabled to be true")
	}
	if goCfg.Test == nil || !*goCfg.Test {
		t.Error("expected default Test to be true")
	}
	if goCfg.Lint == nil || !*goCfg.Lint {
		t.Error("expected default Lint to be true")
	}
	if goCfg.Coverage == nil || *goCfg.Coverage {
		t.Error("expected default Coverage to be false")
	}
}

func TestGetLanguageConfig_Partial(t *testing.T) {
	cfg := DefaultConfig()

	// Partially configured
	f := false
	cfg.Languages["go"] = LanguageConfig{
		Lint: &f, // only lint is set
	}

	goCfg := cfg.GetLanguageConfig("go")

	// Lint should be false (configured)
	if goCfg.Lint == nil || *goCfg.Lint {
		t.Error("expected Lint to be false")
	}

	// Test should be true (default)
	if goCfg.Test == nil || !*goCfg.Test {
		t.Error("expected Test to default to true")
	}
}

func TestBoolPtr(t *testing.T) {
	truePtr := BoolPtr(true)
	if truePtr == nil || !*truePtr {
		t.Error("expected BoolPtr(true) to return pointer to true")
	}

	falsePtr := BoolPtr(false)
	if falsePtr == nil || *falsePtr {
		t.Error("expected BoolPtr(false) to return pointer to false")
	}
}
