// tests/integration/feature_test.go

// Package integration provides end-to-end tests for gz-flow CLI
// using real git repositories and binary execution.
package integration

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

const testFileMode = 0o644

func setupTestRepo(t *testing.T) string {
	t.Helper()

	dir := t.TempDir()

	// Initialize git repo
	run(t, dir, "git", "init")
	run(t, dir, "git", "config", "user.email", "test@test.com")
	run(t, dir, "git", "config", "user.name", "Test")

	// Create initial commit
	readme := filepath.Join(dir, "README.md")
	if err := os.WriteFile(readme, []byte("# Test"), testFileMode); err != nil {
		t.Fatalf("Failed to create README: %v", err)
	}
	run(t, dir, "git", "add", ".")
	run(t, dir, "git", "commit", "-m", "Initial commit")

	// Create develop branch
	run(t, dir, "git", "checkout", "-b", "develop")

	return dir
}

func run(t *testing.T, dir string, name string, args ...string) {
	t.Helper()
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("%s %v failed: %v\n%s", name, args, err, out)
	}
}

func TestFeatureStart(t *testing.T) {
	t.Skip("Integration test requires built binary")

	dir := setupTestRepo(t)

	// Build gz-flow
	binary := filepath.Join(t.TempDir(), "gz-flow")
	build := exec.Command("go", "build", "-o", binary, "./cmd/gz-flow")
	moduleRoot, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	build.Dir = moduleRoot
	// This won't work in isolation - need to build from source

	_ = dir // Will be used when test is implemented
}
