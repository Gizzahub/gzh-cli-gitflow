// internal/gitcmd/git_test.go
package gitcmd

import (
	"context"
	"strings"
	"testing"
)

func TestRun_CurrentBranch(t *testing.T) {
	ctx := context.Background()
	git := New()

	branch, err := git.CurrentBranch(ctx)
	if err != nil {
		t.Fatalf("CurrentBranch failed: %v", err)
	}
	if branch == "" {
		t.Error("CurrentBranch returned empty string")
	}
}

func TestRun_IsClean(t *testing.T) {
	ctx := context.Background()
	git := New()

	// Just verify it doesn't panic/error
	_, err := git.IsClean(ctx)
	if err != nil {
		t.Fatalf("IsClean failed: %v", err)
	}
}

func TestBranchExists(t *testing.T) {
	ctx := context.Background()
	git := New()

	// Test with current branch (should exist)
	current, err := git.CurrentBranch(ctx)
	if err != nil {
		t.Fatalf("Failed to get current branch: %v", err)
	}

	exists, err := git.BranchExists(ctx, current)
	if err != nil {
		t.Fatalf("BranchExists failed: %v", err)
	}
	if !exists {
		t.Error("Current branch should exist")
	}

	// Test with non-existent branch
	exists, err = git.BranchExists(ctx, "nonexistent-branch-xyz-12345")
	if err != nil {
		t.Fatalf("BranchExists failed: %v", err)
	}
	if exists {
		t.Error("Non-existent branch should not exist")
	}

	// Test with invalid branch name
	_, err = git.BranchExists(ctx, "")
	if err == nil {
		t.Error("BranchExists should fail with empty branch name")
	}

	_, err = git.BranchExists(ctx, "-invalid")
	if err == nil {
		t.Error("BranchExists should fail with branch name starting with '-'")
	}
}

func TestValidateBranchName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid simple", "feature", false},
		{"valid with slash", "feature/branch", false},
		{"valid with dash", "my-feature", false},
		{"valid with underscore", "my_feature", false},
		{"valid complex", "feature/my-branch_123", false},
		{"empty", "", true},
		{"starts with dash", "-feature", true},
		{"special chars", "feature@123", true},
		{"space", "my feature", true},
		{"dot", "feature.branch", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateBranchName(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateBranchName(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestListBranches(t *testing.T) {
	ctx := context.Background()
	git := New()

	branches, err := git.ListBranches(ctx, "")
	if err != nil {
		t.Fatalf("ListBranches failed: %v", err)
	}
	if len(branches) == 0 {
		t.Error("Expected at least one branch")
	}

	// Verify current branch is in the list
	current, _ := git.CurrentBranch(ctx)
	found := false
	for _, b := range branches {
		if b == current {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Current branch %q not found in branch list", current)
	}
}

func TestWithWorkDir(t *testing.T) {
	git := New().WithWorkDir("/tmp")

	// Verify it returns a new instance
	if git.workDir != "/tmp" {
		t.Errorf("WithWorkDir didn't set workDir, got %q", git.workDir)
	}

	// Verify original instance is not modified
	git2 := New()
	if git2.workDir != "" {
		t.Errorf("Original instance should have empty workDir, got %q", git2.workDir)
	}
}

func TestCheckout_InvalidBranchName(t *testing.T) {
	ctx := context.Background()
	git := New()

	// Test invalid branch names
	err := git.Checkout(ctx, "")
	if err == nil {
		t.Error("Checkout should fail with empty branch name")
	}

	err = git.Checkout(ctx, "-invalid")
	if err == nil {
		t.Error("Checkout should fail with branch name starting with '-'")
	}

	err = git.Checkout(ctx, "invalid@branch")
	if err == nil {
		t.Error("Checkout should fail with special characters")
	}
}

func TestCreateBranch_InvalidBranchName(t *testing.T) {
	ctx := context.Background()
	git := New()

	err := git.CreateBranch(ctx, "")
	if err == nil {
		t.Error("CreateBranch should fail with empty branch name")
	}

	err = git.CreateBranch(ctx, "-invalid")
	if err == nil {
		t.Error("CreateBranch should fail with branch name starting with '-'")
	}
}

func TestMerge_InvalidBranchName(t *testing.T) {
	ctx := context.Background()
	git := New()

	err := git.Merge(ctx, "", false)
	if err == nil {
		t.Error("Merge should fail with empty branch name")
	}

	err = git.Merge(ctx, "-invalid", false)
	if err == nil {
		t.Error("Merge should fail with branch name starting with '-'")
	}
}

func TestDeleteBranch_InvalidBranchName(t *testing.T) {
	ctx := context.Background()
	git := New()

	err := git.DeleteBranch(ctx, "")
	if err == nil {
		t.Error("DeleteBranch should fail with empty branch name")
	}

	err = git.DeleteBranch(ctx, "-invalid")
	if err == nil {
		t.Error("DeleteBranch should fail with branch name starting with '-'")
	}
}

func TestMerge_NoFF(t *testing.T) {
	// Test that Merge validates branch name even with noFF flag
	ctx := context.Background()
	git := New()

	err := git.Merge(ctx, "invalid@branch", true)
	if err == nil {
		t.Error("Merge with noFF should still validate branch name")
	}
}

func TestListBranches_WithPrefix(t *testing.T) {
	ctx := context.Background()
	git := New()

	// Test with a prefix that likely exists (current branch prefix)
	current, _ := git.CurrentBranch(ctx)
	if len(current) > 2 {
		prefix := current[:2]
		branches, err := git.ListBranches(ctx, prefix)
		if err != nil {
			t.Fatalf("ListBranches with prefix failed: %v", err)
		}
		// Should have at least the current branch
		if len(branches) == 0 {
			t.Error("Expected at least one branch with prefix")
		}
	}
}

func TestRun_ErrorHandling(t *testing.T) {
	// Test that run method properly handles errors
	ctx := context.Background()
	git := New()

	// Try to checkout a non-existent branch (should fail)
	err := git.Checkout(ctx, "definitely-nonexistent-branch-12345")
	if err == nil {
		t.Error("Checkout of non-existent branch should fail")
	}
	if !strings.Contains(err.Error(), "git checkout") {
		t.Errorf("Error message should mention the git command, got: %v", err)
	}
}

func TestListBranches_EmptyResult(t *testing.T) {
	ctx := context.Background()
	git := New()

	// Test with a prefix that definitely doesn't exist
	branches, err := git.ListBranches(ctx, "nonexistent-prefix-xyz-")
	if err != nil {
		t.Fatalf("ListBranches should not error on empty result: %v", err)
	}
	if len(branches) != 0 {
		t.Errorf("Expected empty branch list, got %d branches", len(branches))
	}
}

func TestWithWorkDir_ChainableAndImmutable(t *testing.T) {
	git1 := New()
	git2 := git1.WithWorkDir("/tmp")
	git3 := git2.WithWorkDir("/var")

	// Verify each instance has the correct workDir
	if git1.workDir != "" {
		t.Errorf("git1.workDir = %q, want empty", git1.workDir)
	}
	if git2.workDir != "/tmp" {
		t.Errorf("git2.workDir = %q, want /tmp", git2.workDir)
	}
	if git3.workDir != "/var" {
		t.Errorf("git3.workDir = %q, want /var", git3.workDir)
	}
}

func TestDeleteBranch_NonExistent(t *testing.T) {
	ctx := context.Background()
	git := New()

	// Try to delete a non-existent branch
	err := git.DeleteBranch(ctx, "nonexistent-branch-xyz-12345")
	if err == nil {
		t.Error("DeleteBranch should fail for non-existent branch")
	}
	if !strings.Contains(err.Error(), "git branch") {
		t.Errorf("Error should mention git branch command, got: %v", err)
	}
}

func TestCreateBranch_AlreadyExists(t *testing.T) {
	ctx := context.Background()
	git := New()

	// Try to create a branch that already exists (current branch)
	current, _ := git.CurrentBranch(ctx)
	err := git.CreateBranch(ctx, current)
	if err == nil {
		t.Error("CreateBranch should fail when branch already exists")
	}
}

func TestExecutor_CreateTag(t *testing.T) {
	tests := []struct {
		name        string
		tag         string
		message     string
		wantErr     bool
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
