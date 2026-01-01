package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config represents the complete gitflow configuration
type Config struct {
	Branches BranchConfig   `yaml:"branches"`
	Prefixes PrefixConfig   `yaml:"prefixes"`
	Options  OptionsConfig  `yaml:"options"`
	Guardian GuardianConfig `yaml:"guardian"`
}

// BranchConfig defines the main branch names
type BranchConfig struct {
	Master  string `yaml:"master"`
	Develop string `yaml:"develop"`
}

// PrefixConfig defines the prefixes for each flow type
type PrefixConfig struct {
	Feature string `yaml:"feature"`
	Release string `yaml:"release"`
	Hotfix  string `yaml:"hotfix"`
}

// OptionsConfig defines workflow options
type OptionsConfig struct {
	DeleteBranchAfterFinish bool   `yaml:"delete_branch_after_finish"`
	PushAfterFinish         bool   `yaml:"push_after_finish"`
	TagFormat               string `yaml:"tag_format"`
	RequireCleanTree        bool   `yaml:"require_clean_tree"`
}

// Default returns a Config with default gitflow settings
func Default() *Config {
	return &Config{
		Branches: BranchConfig{
			Master:  "master",
			Develop: "develop",
		},
		Prefixes: PrefixConfig{
			Feature: "feature/",
			Release: "release/",
			Hotfix:  "hotfix/",
		},
		Options: OptionsConfig{
			DeleteBranchAfterFinish: true,
			PushAfterFinish:         false,
			TagFormat:               "v%s",
			RequireCleanTree:        true,
		},
		Guardian: GuardianConfig{
			Enabled: false,
			Mode:    "strict",
			Naming: NamingRule{
				Pattern:   "^[a-z0-9-]+$",
				MaxLength: 50,
				Forbidden: []string{},
			},
			Workflow: WorkflowRules{
				RequireCleanTree:     true,
				RequireUpToDate:      true,
				PreventDirectPush:    false,
				RequireLinearHistory: false,
			},
		},
	}
}

// Load loads configuration from a YAML file
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	cfg := Default()
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return cfg, nil
}

// LoadFromDir loads configuration from a directory
// Checks for local .gzflow.yaml first, then falls back to global ~/.gz/gitflow
func LoadFromDir(dir string) (*Config, error) {
	// Try local config first
	localPath := filepath.Join(dir, ".gzflow.yaml")
	if _, err := os.Stat(localPath); err == nil {
		return Load(localPath)
	}

	// Try global config
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return Default(), nil
	}

	globalPath := filepath.Join(homeDir, ".gz", "gitflow")
	if _, err := os.Stat(globalPath); err == nil {
		return Load(globalPath)
	}

	// No config found, return default
	return Default(), nil
}

// Save saves configuration to a YAML file
func Save(cfg *Config, path string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// Save is a convenience method on Config
func (c *Config) Save(path string) error {
	return Save(c, path)
}
