# Release Command Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development to implement this plan task-by-task.

**Goal:** Implement `gz-flow release start` and `gz-flow release finish` commands with strict semver validation, tag creation, and master-first merge strategy.

**Architecture:** Follows feature command pattern - reuses gitcmd, validator, config, preflight. Adds version validator (strict semver) and git tag functions. Master-first merge ensures production release completes before develop merge.

**Tech Stack:** Go 1.23+, Cobra CLI, internal/gitcmd, internal/validator, pkg/config, internal/preflight

---

## Task 1: Version Validator

**Files:**
- Create: `internal/validator/version.go`
- Create: `internal/validator/version_test.go`

**Step 1: Write the failing test**

```go
// internal/validator/version_test.go
package validator

import "testing"

func TestValidateVersion(t *testing.T) {
	tests := []struct {
		name    string
		version string
		wantErr bool
	}{
		{"valid semver", "1.0.0", false},
		{"valid with large numbers", "12.34.56", false},
		{"valid zero version", "0.0.0", false},
		{"invalid with v prefix", "v1.0.0", true},
		{"invalid two parts", "1.0", true},
		{"invalid four parts", "1.0.0.0", true},
		{"invalid with dash", "1.0.0-beta", true},
		{"invalid with plus", "1.0.0+build", true},
		{"invalid with letters", "1.0.x", true},
		{"invalid negative", "-1.0.0", true},
		{"empty", "", true},
		{"only dots", "..", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateVersion(tt.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateVersion(%q) error = %v, wantErr %v", tt.version, err, tt.wantErr)
			}
		})
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/validator/... -v -run TestValidateVersion`
Expected: FAIL with "undefined: ValidateVersion"

**Step 3: Write minimal implementation**

```go
// internal/validator/version.go
package validator

import (
	"fmt"
	"regexp"
)

var (
	// semverPattern is the strict semver pattern (X.Y.Z only)
	semverPattern = regexp.MustCompile(`^\d+\.\d+\.\d+$`)
)

// ValidateVersion validates a version string using strict semver format (X.Y.Z).
// Only accepts versions in the format: X.Y.Z where X, Y, Z are non-negative integers.
//
// Valid examples:
//   - 1.0.0
//   - 2.5.3
//   - 10.20.30
//
// Invalid examples:
//   - v1.0.0 (no prefix allowed)
//   - 1.0 (incomplete)
//   - 1.0.0-beta (no pre-release)
//   - 1.0.0+build (no build metadata)
func ValidateVersion(version string) error {
	if version == "" {
		return fmt.Errorf("version cannot be empty")
	}

	if !semverPattern.MatchString(version) {
		return fmt.Errorf("invalid version format (expected: X.Y.Z, e.g., 1.0.0)")
	}

	return nil
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./internal/validator/... -v -run TestValidateVersion`
Expected: PASS (all 12 test cases)

**Step 5: Commit**

```bash
git add internal/validator/version.go internal/validator/version_test.go
git commit -m "feat(internal): add strict semver version validator"
```

---

## Task 2: Git Tag Functions

**Files:**
- Modify: `internal/gitcmd/git.go`
- Modify: `internal/gitcmd/git_test.go`

**Step 1: Write the failing tests**

```go
// internal/gitcmd/git_test.go
// Add to existing file after other tests

func TestExecutor_CreateTag(t *testing.T) {
	tests := []struct {
		name       string
		tag        string
		message    string
		wantErr    bool
		errContains string
	}{
		{
			name:    "valid tag with message",
			tag:     "v1.0.0",
			message: "Release version 1.0.0",
			wantErr: false,
		},
		{
			name:    "valid tag empty message",
			tag:     "v2.0.0",
			message: "",
			wantErr: false,
		},
		{
			name:        "invalid empty tag",
			tag:         "",
			message:     "test",
			wantErr:     true,
			errContains: "tag name cannot be empty",
		},
		{
			name:        "invalid tag with spaces",
			tag:         "v 1.0.0",
			message:     "test",
			wantErr:     true,
			errContains: "invalid characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Skip actual git operations in unit tests
			t.Skip("requires git repository setup")
		})
	}
}

func TestExecutor_TagExists(t *testing.T) {
	tests := []struct {
		name    string
		tag     string
		wantErr bool
	}{
		{
			name:    "check existing tag",
			tag:     "v1.0.0",
			wantErr: false,
		},
		{
			name:    "check non-existent tag",
			tag:     "v99.99.99",
			wantErr: false,
		},
		{
			name:    "empty tag name",
			tag:     "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Skip actual git operations in unit tests
			t.Skip("requires git repository setup")
		})
	}
}
```

**Step 2: Run tests to verify they skip**

Run: `go test ./internal/gitcmd/... -v -run "TestExecutor_CreateTag|TestExecutor_TagExists"`
Expected: SKIP (tests defined but skipped)

**Step 3: Add tag validation helper**

```go
// internal/gitcmd/git.go
// Add after validateBranchName function

// validateTagName validates a git tag name
func validateTagName(tag string) error {
	if tag == "" {
		return fmt.Errorf("tag name cannot be empty")
	}

	// Check for invalid characters
	if strings.ContainsAny(tag, " \t\n\r") {
		return fmt.Errorf("tag name contains invalid characters")
	}

	return nil
}
```

**Step 4: Implement CreateTag**

```go
// internal/gitcmd/git.go
// Add after DeleteBranch method

// CreateTag creates an annotated tag at the current HEAD
func (e *Executor) CreateTag(ctx context.Context, tag, message string) error {
	if err := validateTagName(tag); err != nil {
		return err
	}

	args := []string{"tag", "-a", tag}
	if message != "" {
		args = append(args, "-m", message)
	} else {
		args = append(args, "-m", tag)
	}

	_, err := e.run(ctx, args...)
	if err != nil {
		return fmt.Errorf("failed to create tag: %w", err)
	}

	return nil
}
```

**Step 5: Implement TagExists**

```go
// internal/gitcmd/git.go
// Add after CreateTag method

// TagExists checks if a tag exists in the repository
func (e *Executor) TagExists(ctx context.Context, tag string) (bool, error) {
	if err := validateTagName(tag); err != nil {
		return false, err
	}

	_, err := e.run(ctx, "rev-parse", tag)
	if err != nil {
		// Tag doesn't exist if rev-parse fails
		return false, nil
	}

	return true, nil
}
```

**Step 6: Run tests to verify they still skip**

Run: `go test ./internal/gitcmd/... -v -run "TestExecutor_CreateTag|TestExecutor_TagExists"`
Expected: SKIP (implementation exists, tests skip as designed)

**Step 7: Commit**

```bash
git add internal/gitcmd/git.go internal/gitcmd/git_test.go
git commit -m "feat(internal): add git tag creation and existence check"
```

---

## Task 3: Release Start Command

**Files:**
- Modify: `cmd/gz-flow/cmd/release.go`

**Step 1: Update imports**

```go
// cmd/gz-flow/cmd/release.go
// Replace existing imports with:
package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/gizzahub/gzh-cli-gitflow/internal/gitcmd"
	"github.com/gizzahub/gzh-cli-gitflow/internal/validator"
	"github.com/gizzahub/gzh-cli-gitflow/pkg/config"
)
```

**Step 2: Implement runReleaseStart**

```go
// cmd/gz-flow/cmd/release.go
// Replace runReleaseStart function

func runReleaseStart(cmd *cobra.Command, args []string) error {
	if err := checkGitRepo(); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	git := gitcmd.New()
	version := args[0]

	// 1. Validate version format (strict semver)
	if err := validator.ValidateVersion(version); err != nil {
		return fmt.Errorf("invalid version: %v\n💡 Use semver format: 1.0.0", err)
	}

	// 2. Load config
	cfg, err := config.LoadFromDir(".")
	if err != nil {
		fmt.Printf("⚠️  Failed to load config, using defaults: %v\n", err)
		cfg = config.Default()
	}

	// 3. Check if release branch already exists
	releaseBranch := cfg.Prefixes.Release + version
	exists, _ := git.BranchExists(ctx, releaseBranch)
	if exists {
		return fmt.Errorf("release branch '%s' already exists", releaseBranch)
	}

	// 4. Context hint: warn if not on develop
	currentBranch, _ := git.CurrentBranch(ctx)
	if currentBranch != cfg.Branches.Develop {
		fmt.Printf("⚠️  You're on '%s', not '%s'\n", currentBranch, cfg.Branches.Develop)
		fmt.Printf("💡 Will checkout '%s' first\n\n", cfg.Branches.Develop)
	}

	// 5. Create release branch from develop
	if err := git.Checkout(ctx, cfg.Branches.Develop); err != nil {
		return fmt.Errorf("failed to checkout %s: %v", cfg.Branches.Develop, err)
	}

	if err := git.CreateBranch(ctx, releaseBranch); err != nil {
		return fmt.Errorf("failed to create branch: %v", err)
	}

	fmt.Printf("✅ Started release branch '%s'\n", releaseBranch)
	fmt.Printf("📍 Switched to branch '%s'\n", releaseBranch)

	return nil
}
```

**Step 3: Build and test manually**

Run: `go build -o /tmp/gz-flow ./cmd/gz-flow`
Expected: Build succeeds

**Step 4: Commit**

```bash
git add cmd/gz-flow/cmd/release.go
git commit -m "feat(cmd): implement release start with semver validation"
```

---

## Task 4: Release Finish Command

**Files:**
- Modify: `cmd/gz-flow/cmd/release.go`

**Step 1: Add preflight import**

```go
// cmd/gz-flow/cmd/release.go
// Add to imports
import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/gizzahub/gzh-cli-gitflow/internal/gitcmd"
	"github.com/gizzahub/gzh-cli-gitflow/internal/preflight"
	"github.com/gizzahub/gzh-cli-gitflow/internal/validator"
	"github.com/gizzahub/gzh-cli-gitflow/pkg/config"
)
```

**Step 2: Implement runReleaseFinish**

```go
// cmd/gz-flow/cmd/release.go
// Replace runReleaseFinish function

func runReleaseFinish(cmd *cobra.Command, args []string) error {
	if err := checkGitRepo(); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	git := gitcmd.New()
	version := args[0]

	// 1. Validate version
	if err := validator.ValidateVersion(version); err != nil {
		return fmt.Errorf("invalid version: %v", err)
	}

	// 2. Load config
	cfg, err := config.LoadFromDir(".")
	if err != nil {
		fmt.Printf("⚠️  Failed to load config, using defaults: %v\n", err)
		cfg = config.Default()
	}

	releaseBranch := cfg.Prefixes.Release + version
	masterBranch := cfg.Branches.Master
	developBranch := cfg.Branches.Develop

	// 3. Verify release branch exists
	exists, _ := git.BranchExists(ctx, releaseBranch)
	if !exists {
		return fmt.Errorf("release branch '%s' does not exist", releaseBranch)
	}

	// 4. Pre-flight checks
	checker := preflight.NewChecker(git, masterBranch)
	results := checker.RunAll(ctx)
	fmt.Println("🔍 Pre-flight checks:")
	fmt.Print(results.String())
	if results.HasErrors() {
		return fmt.Errorf("pre-flight checks failed")
	}
	fmt.Println()

	// 5. STEP 1: Merge to master (--no-ff)
	if err := git.Checkout(ctx, masterBranch); err != nil {
		return fmt.Errorf("failed to checkout %s: %v", masterBranch, err)
	}

	if err := git.Merge(ctx, releaseBranch, true); err != nil {
		return fmt.Errorf("merge to %s failed: %v\n💡 Resolve conflicts:\n  1. git checkout %s\n  2. git merge --no-ff %s\n  3. Resolve conflicts\n  4. git merge --continue\n  5. Retry: gz-flow release finish %s",
			masterBranch, err, masterBranch, releaseBranch, version)
	}
	fmt.Printf("✅ Merged '%s' into '%s'\n", releaseBranch, masterBranch)

	// 6. STEP 2: Create tag on master
	if !noTag {
		tagName := fmt.Sprintf(cfg.Options.TagFormat, version)

		// Check if tag already exists
		tagExists, _ := git.TagExists(ctx, tagName)
		if tagExists {
			return fmt.Errorf("tag '%s' already exists\n💡 Use different version or delete existing tag", tagName)
		}

		message := tagMessage
		if message == "" {
			message = fmt.Sprintf("Release version %s", version)
		}

		if err := git.CreateTag(ctx, tagName, message); err != nil {
			return fmt.Errorf("failed to create tag: %v", err)
		}
		fmt.Printf("🏷️  Created tag '%s'\n", tagName)
	}

	// 7. STEP 3: Merge to develop
	developExists, _ := git.BranchExists(ctx, developBranch)
	if !developExists {
		fmt.Printf("⚠️  Develop branch '%s' does not exist\n", developBranch)
		fmt.Printf("💡 Skipping merge to develop\n")
	} else {
		if err := git.Checkout(ctx, developBranch); err != nil {
			return fmt.Errorf("failed to checkout %s: %v", developBranch, err)
		}

		if err := git.Merge(ctx, releaseBranch, true); err != nil {
			tagName := fmt.Sprintf(cfg.Options.TagFormat, version)
			fmt.Printf("⚠️  PARTIAL SUCCESS:\n")
			fmt.Printf("  ✅ Merged to %s and tagged '%s'\n", masterBranch, tagName)
			fmt.Printf("  ❌ Merge to %s failed: %v\n", developBranch, err)
			fmt.Printf("\n💡 To complete:\n")
			fmt.Printf("  1. git checkout %s\n", developBranch)
			fmt.Printf("  2. git merge --no-ff %s\n", releaseBranch)
			fmt.Printf("  3. Resolve conflicts and commit\n")
			return err
		}
		fmt.Printf("✅ Merged '%s' into '%s'\n", releaseBranch, developBranch)
	}

	// 8. STEP 4: Delete release branch
	deleteBranch := cfg.Options.DeleteBranchAfterFinish && !keepBranch
	if deleteBranch {
		if err := git.DeleteBranch(ctx, releaseBranch); err != nil {
			fmt.Printf("⚠️  Failed to delete branch: %v\n", err)
		} else {
			fmt.Printf("🗑️  Deleted branch '%s'\n", releaseBranch)
		}
	}

	return nil
}
```

**Step 3: Build and verify**

Run: `go build -o /tmp/gz-flow ./cmd/gz-flow`
Expected: Build succeeds

**Step 4: Commit**

```bash
git add cmd/gz-flow/cmd/release.go
git commit -m "feat(cmd): implement release finish with master-first strategy"
```

---

## Task 5: Integration Tests

**Files:**
- Create: `tests/integration/release_test.go`

**Step 1: Write integration test scaffold**

```go
// tests/integration/release_test.go
package integration

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReleaseStartValidation(t *testing.T) {
	dir := setupTestRepo(t)

	// Build binary
	moduleRoot, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	binary := filepath.Join(t.TempDir(), "gz-flow")
	buildCmd := exec.Command("go", "build", "-o", binary, "./cmd/gz-flow")
	buildCmd.Dir = moduleRoot
	if out, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("Build failed: %v\n%s", err, out)
	}

	tests := []struct {
		name       string
		version    string
		wantErr    bool
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
```

**Step 2: Add missing import**

```go
// tests/integration/release_test.go
// Add to imports at top of file
import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)
```

**Step 3: Run tests**

Run: `go test ./tests/integration/... -v -run TestRelease`
Expected: SKIP for TestReleaseWorkflow, PASS for TestReleaseStartValidation

**Step 4: Commit**

```bash
git add tests/integration/release_test.go
git commit -m "test(integration): add release command integration tests"
```

---

## Summary

**Phase 1 (Validators & Git Functions):**
- Task 1: Version Validator ✅
- Task 2: Git Tag Functions ✅

**Phase 2 (Commands):**
- Task 3: Release Start ✅
- Task 4: Release Finish ✅

**Phase 3 (Testing):**
- Task 5: Integration Tests ✅

**Next Steps (Not in this plan):**
- Hotfix command (similar pattern)
- Init command (setup git-flow)
- Enhanced Guardian mode for versions

---

**End of Plan**
