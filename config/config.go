// Package config provides configuration file support for prepush.
package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the .prepush.yaml configuration.
type Config struct {
	// Global settings
	Verbose bool `yaml:"verbose"`

	// Language-specific settings
	Languages map[string]LanguageConfig `yaml:"languages"`
}

// LanguageConfig holds settings for a specific language.
type LanguageConfig struct {
	Enabled  *bool    `yaml:"enabled"`  // nil means auto-detect
	Paths    []string `yaml:"paths"`    // specific paths to check (empty = auto-detect)
	Test     *bool    `yaml:"test"`     // run tests
	Lint     *bool    `yaml:"lint"`     // run linter
	Format   *bool    `yaml:"format"`   // check formatting
	Coverage *bool    `yaml:"coverage"` // show coverage

	// Go-specific
	ExcludeCoverage string `yaml:"exclude_coverage"` // directories to exclude from coverage
}

// DefaultConfig returns a configuration with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Verbose:   false,
		Languages: make(map[string]LanguageConfig),
	}
}

// Load reads configuration from .prepush.yaml in the given directory.
// Returns default config if file doesn't exist.
func Load(dir string) (Config, error) {
	cfg := DefaultConfig()

	// Try multiple config file names
	configFiles := []string{
		dir + "/.prepush.yaml",
		dir + "/.prepush.yml",
		dir + "/prepush.yaml",
		dir + "/prepush.yml",
	}

	var data []byte
	var err error
	for _, f := range configFiles {
		data, err = os.ReadFile(f)
		if err == nil {
			break
		}
	}

	if err != nil {
		// No config file, return defaults
		return cfg, nil
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}

// IsLanguageEnabled checks if a language is enabled in config.
// Returns true if enabled is nil (auto-detect) or explicitly true.
func (c *Config) IsLanguageEnabled(lang string) bool {
	lc, ok := c.Languages[lang]
	if !ok {
		return true // not configured = auto-detect
	}
	if lc.Enabled == nil {
		return true // nil = auto-detect
	}
	return *lc.Enabled
}

// GetLanguageConfig returns the config for a language, with defaults applied.
func (c *Config) GetLanguageConfig(lang string) LanguageConfig {
	lc, ok := c.Languages[lang]
	if !ok {
		// Return defaults
		t := true
		f := false
		return LanguageConfig{
			Enabled:  &t,
			Test:     &t,
			Lint:     &t,
			Format:   &t,
			Coverage: &f,
		}
	}

	// Apply defaults for nil values
	t := true
	if lc.Enabled == nil {
		lc.Enabled = &t
	}
	if lc.Test == nil {
		lc.Test = &t
	}
	if lc.Lint == nil {
		lc.Lint = &t
	}
	if lc.Format == nil {
		lc.Format = &t
	}
	if lc.Coverage == nil {
		f := false
		lc.Coverage = &f
	}

	return lc
}

// BoolPtr returns a pointer to a bool value.
func BoolPtr(b bool) *bool {
	return &b
}
