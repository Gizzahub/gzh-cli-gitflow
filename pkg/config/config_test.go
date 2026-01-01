package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefault(t *testing.T) {
	cfg := Default()

	// Test branch config
	if cfg.Branches.Master != "master" {
		t.Errorf("Expected master branch 'master', got '%s'", cfg.Branches.Master)
	}
	if cfg.Branches.Develop != "develop" {
		t.Errorf("Expected develop branch 'develop', got '%s'", cfg.Branches.Develop)
	}

	// Test prefix config
	if cfg.Prefixes.Feature != "feature/" {
		t.Errorf("Expected feature prefix 'feature/', got '%s'", cfg.Prefixes.Feature)
	}
	if cfg.Prefixes.Release != "release/" {
		t.Errorf("Expected release prefix 'release/', got '%s'", cfg.Prefixes.Release)
	}
	if cfg.Prefixes.Hotfix != "hotfix/" {
		t.Errorf("Expected hotfix prefix 'hotfix/', got '%s'", cfg.Prefixes.Hotfix)
	}

	// Test options
	if !cfg.Options.DeleteBranchAfterFinish {
		t.Error("Expected DeleteBranchAfterFinish to be true")
	}
	if cfg.Options.PushAfterFinish {
		t.Error("Expected PushAfterFinish to be false")
	}
	if cfg.Options.TagFormat != "v%s" {
		t.Errorf("Expected TagFormat 'v%%s', got '%s'", cfg.Options.TagFormat)
	}

	// Test Guardian (disabled by default)
	if cfg.Guardian.Enabled {
		t.Error("Expected Guardian to be disabled by default")
	}
}

func TestLoadSave(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".gzflow.yaml")

	// Create and save config
	cfg := Default()
	cfg.Guardian.Enabled = true
	if err := cfg.Save(configPath); err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Load config
	loaded, err := Load(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if !loaded.Guardian.Enabled {
		t.Error("Expected Guardian to be enabled after load")
	}
}

func TestLoad_NonExistentFile(t *testing.T) {
	_, err := Load("/non/existent/path/.gzflow.yaml")
	if err == nil {
		t.Error("Expected error when loading non-existent file")
	}
}

func TestLoad_InvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".gzflow.yaml")

	// Write invalid YAML
	if err := os.WriteFile(configPath, []byte("invalid: yaml: content: ["), 0o644); err != nil {
		t.Fatalf("Failed to write invalid YAML: %v", err)
	}

	_, err := Load(configPath)
	if err == nil {
		t.Error("Expected error when loading invalid YAML")
	}
}

func TestGuardianValidation(t *testing.T) {
	tests := []struct {
		name       string
		rule       NamingRule
		branchName string
		wantErr    bool
	}{
		{
			name: "valid name",
			rule: NamingRule{
				Pattern:   "^[a-z0-9-]+$",
				MaxLength: 50,
			},
			branchName: "add-new-feature",
			wantErr:    false,
		},
		{
			name: "invalid pattern",
			rule: NamingRule{
				Pattern:   "^[a-z0-9-]+$",
				MaxLength: 50,
			},
			branchName: "Add_New_Feature",
			wantErr:    true,
		},
		{
			name: "exceeds max length",
			rule: NamingRule{
				Pattern:   "^[a-z0-9-]+$",
				MaxLength: 10,
			},
			branchName: "very-long-branch-name",
			wantErr:    true,
		},
		{
			name: "forbidden word",
			rule: NamingRule{
				Pattern:   "^[a-z0-9-]+$",
				MaxLength: 50,
				Forbidden: []string{"test", "tmp"},
			},
			branchName: "test-feature",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.rule.Validate(tt.branchName)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLoadFromDir(t *testing.T) {
	t.Run("local config exists", func(t *testing.T) {
		dir := t.TempDir()

		// Create local config
		localCfg := &Config{
			Branches: BranchConfig{
				Master:  "main",
				Develop: "dev",
			},
		}
		localPath := filepath.Join(dir, ".gzflow.yaml")
		if err := localCfg.Save(localPath); err != nil {
			t.Fatalf("Failed to save local config: %v", err)
		}

		// Load should use local
		cfg, err := LoadFromDir(dir)
		if err != nil {
			t.Fatalf("LoadFromDir failed: %v", err)
		}
		if cfg.Branches.Master != "main" {
			t.Errorf("Expected main, got %s", cfg.Branches.Master)
		}
		if cfg.Branches.Develop != "dev" {
			t.Errorf("Expected dev, got %s", cfg.Branches.Develop)
		}
	})

	t.Run("no local, returns defaults", func(t *testing.T) {
		dir := t.TempDir()

		cfg, err := LoadFromDir(dir)
		if err != nil {
			t.Fatalf("LoadFromDir failed: %v", err)
		}
		// Should return defaults
		if cfg.Branches.Master != "master" {
			t.Errorf("Expected default 'master', got %s", cfg.Branches.Master)
		}
		if cfg.Branches.Develop != "develop" {
			t.Errorf("Expected default 'develop', got %s", cfg.Branches.Develop)
		}
	})

	t.Run("invalid home dir handling", func(t *testing.T) {
		// Test with temp dir (no .gz/gitflow)
		dir := t.TempDir()

		cfg, err := LoadFromDir(dir)
		if err != nil {
			t.Fatalf("Should not fail: %v", err)
		}
		// Should still return defaults
		if cfg == nil {
			t.Error("Config should not be nil")
		}
		if cfg.Branches.Master != "master" {
			t.Errorf("Expected default 'master', got %s", cfg.Branches.Master)
		}
	})
}

func TestGuardianConfig_ValidateBranchName(t *testing.T) {
	t.Run("guardian disabled", func(t *testing.T) {
		cfg := &GuardianConfig{Enabled: false}

		err := cfg.ValidateBranchName("ANY-INVALID-NAME-123", "feature/")
		if err != nil {
			t.Errorf("Disabled guardian should not validate: %v", err)
		}
	})

	t.Run("guardian enabled - valid name", func(t *testing.T) {
		cfg := &GuardianConfig{
			Enabled: true,
			Naming: NamingRule{
				Pattern: "^[a-z0-9-]+$",
			},
		}

		err := cfg.ValidateBranchName("my-feature", "")
		if err != nil {
			t.Errorf("Valid name should pass: %v", err)
		}
	})

	t.Run("guardian enabled - invalid name", func(t *testing.T) {
		cfg := &GuardianConfig{
			Enabled: true,
			Naming: NamingRule{
				Pattern: "^[a-z0-9-]+$",
			},
		}

		err := cfg.ValidateBranchName("MyFeature", "")
		if err == nil {
			t.Error("Invalid name should fail validation")
		}
	})

	t.Run("prefix stripping", func(t *testing.T) {
		cfg := &GuardianConfig{
			Enabled: true,
			Naming: NamingRule{
				Pattern: "^[a-z0-9-]+$",
			},
		}

		// Should strip "feature/" prefix before validating
		err := cfg.ValidateBranchName("feature/my-feature", "feature/")
		if err != nil {
			t.Errorf("Should strip prefix: %v", err)
		}
	})
}
