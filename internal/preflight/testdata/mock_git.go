package testdata

import "context"

// MockGit is a mock implementation for testing
type MockGit struct {
	IsCleanFunc       func(ctx context.Context) (bool, error)
	BranchExistsFunc  func(ctx context.Context, branch string) (bool, error)
	CurrentBranchFunc func(ctx context.Context) (string, error)
}

func (m *MockGit) IsClean(ctx context.Context) (bool, error) {
	if m.IsCleanFunc != nil {
		return m.IsCleanFunc(ctx)
	}
	return true, nil
}

func (m *MockGit) BranchExists(ctx context.Context, branch string) (bool, error) {
	if m.BranchExistsFunc != nil {
		return m.BranchExistsFunc(ctx, branch)
	}
	return true, nil
}

func (m *MockGit) CurrentBranch(ctx context.Context) (string, error) {
	if m.CurrentBranchFunc != nil {
		return m.CurrentBranchFunc(ctx)
	}
	return "develop", nil
}
