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

		checker := preflight.NewChecker(mockGit, "develop")
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

		checker := preflight.NewChecker(mockGit, "")
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

		checker := preflight.NewChecker(mockGit, "develop")
		results := checker.RunAll(context.Background())

		if !results.HasErrors() {
			t.Error("Expected errors when target branch doesn't exist")
		}

		// Verify output contains helpful hints
		output := results.String()
		if !strings.Contains(output, "❌") {
			t.Error("Expected error emoji in output")
		}
		if !strings.Contains(output, "💡") {
			t.Error("Expected hint emoji in output")
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

		checker := preflight.NewChecker(mockGit, "develop")
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

func TestResult_HasErrors(t *testing.T) {
	tests := []struct {
		name    string
		results preflight.Results
		want    bool
	}{
		{
			name: "no errors",
			results: preflight.Results{
				{Name: "clean", Passed: true},
				{Name: "uptodate", Passed: true},
			},
			want: false,
		},
		{
			name: "one error",
			results: preflight.Results{
				{Name: "clean", Passed: true},
				{Name: "uptodate", Passed: false, Error: "", Hint: "not up to date"},
			},
			want: true,
		},
		{
			name: "all errors",
			results: preflight.Results{
				{Name: "clean", Passed: false, Error: "", Hint: "dirty"},
				{Name: "uptodate", Passed: false, Error: "", Hint: "not up to date"},
			},
			want: true,
		},
		{
			name:    "empty results",
			results: preflight.Results{},
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.results.HasErrors()
			if got != tt.want {
				t.Errorf("HasErrors() = %v, want %v", got, tt.want)
			}
		})
	}
}
