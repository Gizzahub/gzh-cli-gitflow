package preflight_test

import (
	"context"
	"strings"
	"testing"

	"github.com/gizzahub/gzh-cli-gitflow/internal/preflight"
	"github.com/gizzahub/gzh-cli-gitflow/internal/preflight/testdata"
)

func TestChecker_RunAll(t *testing.T) {
	t.Run("all checks pass", func(t *testing.T) {
		mockGit := &testdata.MockGit{
			IsCleanFunc: func(ctx context.Context) (bool, error) {
				return true, nil // Clean tree
			},
			BranchExistsFunc: func(ctx context.Context, branch string) (bool, error) {
				return true, nil
			},
		}

		checker := preflight.New(mockGit).WithTargetBranch("develop")
		results := checker.RunAll(context.Background())

		if results.HasErrors() {
			t.Errorf("Expected no errors, got: %v", results)
		}
	})

	t.Run("dirty working tree", func(t *testing.T) {
		mockGit := &testdata.MockGit{
			IsCleanFunc: func(ctx context.Context) (bool, error) {
				return false, nil
			},
		}

		checker := preflight.New(mockGit)
		results := checker.RunAll(context.Background())

		if !results.HasErrors() {
			t.Error("Expected errors for dirty tree")
		}
	})

	t.Run("target branch does not exist", func(t *testing.T) {
		mockGit := &testdata.MockGit{
			IsCleanFunc: func(ctx context.Context) (bool, error) {
				return true, nil
			},
			BranchExistsFunc: func(ctx context.Context, branch string) (bool, error) {
				return false, nil
			},
		}

		checker := preflight.New(mockGit).WithTargetBranch("develop")
		results := checker.RunAll(context.Background())

		if !results.HasErrors() {
			t.Error("Expected errors when target branch doesn't exist")
		}

		// Verify output contains helpful hints
		output := results.String()
		if !strings.Contains(output, "❌") {
			t.Error("Expected error emoji in output")
		}
		if !strings.Contains(output, "Hint:") {
			t.Error("Expected hint in output")
		}
	})

	t.Run("results string format", func(t *testing.T) {
		mockGit := &testdata.MockGit{
			IsCleanFunc: func(ctx context.Context) (bool, error) {
				return true, nil
			},
			BranchExistsFunc: func(ctx context.Context, branch string) (bool, error) {
				return true, nil
			},
		}

		checker := preflight.New(mockGit).WithTargetBranch("develop")
		results := checker.RunAll(context.Background())

		output := results.String()
		if !strings.Contains(output, "✅") {
			t.Error("Expected success emoji in output")
		}
		if !strings.Contains(output, "Pre-flight checks:") {
			t.Error("Expected header in output")
		}
	})
}
