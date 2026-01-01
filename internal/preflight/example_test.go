package preflight_test

import (
	"context"
	"fmt"

	"github.com/gizzahub/gzh-cli-gitflow/internal/preflight"
	"github.com/gizzahub/gzh-cli-gitflow/internal/preflight/testdata"
)

// ExampleChecker_RunAll_success demonstrates successful pre-flight checks
func ExampleChecker_RunAll_success() {
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

	fmt.Print(results.String())
	// Output:
	// Pre-flight checks:
	//   ✅ Clean working tree
	//   ✅ Target branch 'develop' exists
}

// ExampleChecker_RunAll_failure demonstrates failed pre-flight checks
func ExampleChecker_RunAll_failure() {
	mockGit := &testdata.MockGit{
		IsCleanFunc: func(ctx context.Context) (bool, error) {
			return false, nil
		},
		BranchExistsFunc: func(ctx context.Context, branch string) (bool, error) {
			return false, nil
		},
	}

	checker := preflight.NewChecker(mockGit, "develop")
	results := checker.RunAll(context.Background())

	fmt.Print(results.String())
	// Output:
	// Pre-flight checks:
	//   ❌ Clean working tree
	//      💡 Commit or stash your changes before finishing
	//   ❌ Target branch 'develop' exists
	//      💡 Create branch 'develop' first or check your configuration
}
