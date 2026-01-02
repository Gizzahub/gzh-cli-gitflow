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
	// Setup: Create test repo and build binary
	dir := setupTestRepo(t)

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	moduleRoot := filepath.Join(cwd, "..", "..")

	binary := filepath.Join(t.TempDir(), "gz-flow")
	buildCmd := exec.Command("go", "build", "-o", binary, "./cmd/gz-flow")
	buildCmd.Dir = moduleRoot
	if out, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("Build failed: %v\n%s", err, out)
	}

	// STEP 1: Start release branch
	startCmd := exec.Command(binary, "release", "start", "1.0.0")
	startCmd.Dir = dir
	if out, err := startCmd.CombinedOutput(); err != nil {
		t.Fatalf("release start failed: %v\nOutput: %s", err, out)
	}

	// Verify we're on release/1.0.0
	currentBranch := gitCommand(t, dir, "branch", "--show-current")
	if strings.TrimSpace(currentBranch) != "release/1.0.0" {
		t.Fatalf("Expected to be on release/1.0.0, got: %s", currentBranch)
	}

	// STEP 2: Make commits on release branch
	changelog := filepath.Join(dir, "CHANGELOG.md")
	if err := os.WriteFile(changelog, []byte("# Release 1.0.0\n\n- Feature A\n- Feature B\n"), 0644); err != nil {
		t.Fatalf("Failed to create CHANGELOG: %v", err)
	}
	run(t, dir, "git", "add", "CHANGELOG.md")
	run(t, dir, "git", "commit", "-m", "Add release notes")

	// STEP 3: Finish release
	finishCmd := exec.Command(binary, "release", "finish", "1.0.0")
	finishCmd.Dir = dir
	if out, err := finishCmd.CombinedOutput(); err != nil {
		t.Fatalf("release finish failed: %v\nOutput: %s", err, out)
	}

	// STEP 4: Verify git state
	// TODO: Implement verification logic
	// This is where you decide WHAT to verify and HOW thoroughly
	// See the comment block below for guidance
	verifyReleaseWorkflowState(t, dir)
}

// verifyReleaseWorkflowState checks the git state after release finish
// Strategy: Option B - Thorough verification
// Verifies observable behavior without checking implementation details
func verifyReleaseWorkflowState(t *testing.T, dir string) {
	t.Helper()

	// 1. Verify tag v1.0.0 exists
	tags := gitCommand(t, dir, "tag", "-l")
	if !strings.Contains(tags, "v1.0.0") {
		t.Errorf("Tag v1.0.0 not found. Tags:\n%s", tags)
	}

	// 2. Verify release branch deleted
	branches := gitCommand(t, dir, "branch", "-a")
	if strings.Contains(branches, "release/1.0.0") {
		t.Errorf("Release branch should be deleted. Branches:\n%s", branches)
	}

	// 3. Verify CHANGELOG.md exists in master
	run(t, dir, "git", "checkout", "master")
	if _, err := os.Stat(filepath.Join(dir, "CHANGELOG.md")); os.IsNotExist(err) {
		t.Error("CHANGELOG.md should exist in master branch")
	}

	// 4. Verify CHANGELOG.md exists in develop
	run(t, dir, "git", "checkout", "develop")
	if _, err := os.Stat(filepath.Join(dir, "CHANGELOG.md")); os.IsNotExist(err) {
		t.Error("CHANGELOG.md should exist in develop branch")
	}

	// 5. Verify current branch is develop (release finish should leave us on develop)
	currentBranch := gitCommand(t, dir, "branch", "--show-current")
	if strings.TrimSpace(currentBranch) != "develop" {
		t.Errorf("Should be on develop after release finish, got: %s", currentBranch)
	}
}

// gitCommand is a helper to run git commands and return output
func gitCommand(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v failed: %v\nOutput: %s", args, err, out)
	}
	return string(out)
}
