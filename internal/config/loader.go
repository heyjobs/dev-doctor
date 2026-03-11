package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/yourusername/dev-doctor/internal/types"
	"gopkg.in/yaml.v3"
)

// Loader handles loading and parsing diagnostic configuration
type Loader struct {
	configPath string
}

// NewLoader creates a new configuration loader
func NewLoader(configPath string) *Loader {
	return &Loader{
		configPath: configPath,
	}
}

// Load reads and parses the diagnostic configuration file
func (l *Loader) Load() (*types.DiagnosticConfig, error) {
	// If no config path provided, use default
	if l.configPath == "" {
		l.configPath = l.getDefaultConfigPath()
	}

	// Read the file
	data, err := os.ReadFile(l.configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse YAML
	var config types.DiagnosticConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Validate
	if err := l.validate(&config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &config, nil
}

// getDefaultConfigPath returns the default configuration file path
func (l *Loader) getDefaultConfigPath() string {
	// Try to find config relative to executable or in working directory
	candidates := []string{
		"configs/diagnostics.yaml",
		"./diagnostics.yaml",
		"/etc/dev-doctor/diagnostics.yaml",
	}

	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}

	// If none found, return the default
	return "configs/diagnostics.yaml"
}

// validate ensures the configuration is well-formed
func (l *Loader) validate(config *types.DiagnosticConfig) error {
	if len(config.Tests) == 0 {
		return fmt.Errorf("no tests defined in configuration")
	}

	for i, test := range config.Tests {
		if test.Test == "" {
			return fmt.Errorf("test %d: missing test identifier", i)
		}
		if test.Description == "" {
			return fmt.Errorf("test %s: missing description", test.Test)
		}
		if test.Diagnostic == "" {
			return fmt.Errorf("test %s: missing diagnostic identifier", test.Test)
		}
		if test.Severity == "" {
			return fmt.Errorf("test %s: missing severity", test.Test)
		}
	}

	return nil
}

// GetConfigPath returns the configuration file path being used
func (l *Loader) GetConfigPath() string {
	return l.configPath
}

// FilterByProfile returns a new config containing tests that match the given profile.
// All profiles include "basic" tests. Other profiles add their specific tests to basic.
func FilterByProfile(config *types.DiagnosticConfig, profile string) *types.DiagnosticConfig {
	if profile == "" {
		return config
	}

	filtered := &types.DiagnosticConfig{
		Tests: []types.DiagnosticTest{},
	}

	for _, test := range config.Tests {
		// Always include basic tests
		hasBasic := false
		hasProfile := false

		for _, p := range test.Profiles {
			if p == "basic" {
				hasBasic = true
			}
			if p == profile {
				hasProfile = true
			}
		}

		// Include if it's a basic test OR matches the selected profile
		if hasBasic || hasProfile {
			filtered.Tests = append(filtered.Tests, test)
		}
	}

	return filtered
}

// FindConfig attempts to locate the configuration file
func FindConfig() (string, error) {
	// Get executable directory
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	exeDir := filepath.Dir(exe)

	// Get working directory
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Try multiple locations
	candidates := []string{
		filepath.Join(wd, "configs", "diagnostics.yaml"),
		filepath.Join(exeDir, "configs", "diagnostics.yaml"),
		filepath.Join(wd, "diagnostics.yaml"),
		"/etc/dev-doctor/diagnostics.yaml",
	}

	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
	}

	return "", fmt.Errorf("configuration file not found in standard locations")
}
