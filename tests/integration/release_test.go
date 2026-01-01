// tests/integration/release_test.go

// Package integration provides end-to-end tests for gz-flow CLI
// using real git repositories and binary execution.
package integration

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestReleaseStartValidation(t *testing.T) {
	dir := setupTestRepo(t)

	// Build binary
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	// Navigate up from tests/integration to module root
	moduleRoot := filepath.Join(cwd, "..", "..")

	binary := filepath.Join(t.TempDir(), "gz-flow")
	buildCmd := exec.Command("go", "build", "-o", binary, "./cmd/gz-flow")
	buildCmd.Dir = moduleRoot
	if out, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("Build failed: %v\n%s", err, out)
	}

	tests := []struct {
		name        string
		version     string
		wantErr     bool
		errContains string
	}{
		{
			name:    "valid semver",
			version: "1.0.0",
			wantErr: false,
		},
		{
			name:        "invalid with v prefix",
			version:     "v1.0.0",
			wantErr:     true,
			errContains: "invalid version",
		},
		{
			name:        "invalid incomplete",
			version:     "1.0",
			wantErr:     true,
			errContains: "invalid version",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binary, "release", "start", tt.version)
			cmd.Dir = dir
			output, err := cmd.CombinedOutput()

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none. Output: %s", output)
				}
				if !strings.Contains(string(output), tt.errContains) {
					t.Errorf("Expected error containing %q, got: %s", tt.errContains, output)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v\nOutput: %s", err, output)
				}
			}
		})
	}
}

func TestReleaseWorkflow(t *testing.T) {
	t.Skip("Full workflow test requires git repo setup and binary build")

	// TODO: Implement full workflow test:
	// 1. release start 1.0.0
	// 2. Make commits on release/1.0.0
	// 3. release finish 1.0.0
	// 4. Verify tag exists, branches merged, release deleted
}
