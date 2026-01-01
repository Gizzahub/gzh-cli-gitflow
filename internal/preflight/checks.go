// Package preflight provides pre-flight checks for git-flow operations.
package preflight

import (
	"context"
	"fmt"
	"strings"
)

// GitExecutor defines the interface for git operations needed by preflight checks
type GitExecutor interface {
	IsClean(ctx context.Context) (bool, error)
	BranchExists(ctx context.Context, branch string) (bool, error)
	CurrentBranch(ctx context.Context) (string, error)
}

// Result represents the result of a single pre-flight check
type Result struct {
	Name   string
	Passed bool
	Error  string
	Hint   string
}

// Results is a collection of pre-flight check results
type Results []Result

// HasErrors returns true if any check failed
func (r Results) HasErrors() bool {
	for _, result := range r {
		if !result.Passed {
			return true
		}
	}
	return false
}

// String formats the results for display
func (r Results) String() string {
	var sb strings.Builder
	sb.WriteString("Pre-flight checks:\n")
	for _, result := range r {
		if result.Passed {
			sb.WriteString(fmt.Sprintf("  ✅ %s\n", result.Name))
		} else {
			sb.WriteString(fmt.Sprintf("  ❌ %s\n", result.Name))
			if result.Error != "" {
				sb.WriteString(fmt.Sprintf("     Error: %s\n", result.Error))
			}
			if result.Hint != "" {
				sb.WriteString(fmt.Sprintf("     💡 %s\n", result.Hint))
			}
		}
	}
	return sb.String()
}

// Checker performs pre-flight checks before git-flow operations
type Checker struct {
	git          GitExecutor
	targetBranch string
}

// NewChecker creates a new Checker
func NewChecker(git GitExecutor, targetBranch string) *Checker {
	return &Checker{
		git:          git,
		targetBranch: targetBranch,
	}
}

// RunAll runs all pre-flight checks
func (c *Checker) RunAll(ctx context.Context) Results {
	var results Results

	// Check 1: Clean working tree
	results = append(results, c.checkCleanTree(ctx))

	// Check 2: Target branch exists (if specified)
	if c.targetBranch != "" {
		results = append(results, c.checkBranchUpToDate(ctx))
	}

	return results
}

// checkCleanTree verifies the working directory is clean
func (c *Checker) checkCleanTree(ctx context.Context) Result {
	clean, err := c.git.IsClean(ctx)
	if err != nil {
		return Result{
			Name:   "Clean working tree",
			Passed: false,
			Error:  err.Error(),
			Hint:   "Failed to check git status",
		}
	}

	if !clean {
		return Result{
			Name:   "Clean working tree",
			Passed: false,
			Hint:   "Commit or stash your changes before finishing",
		}
	}

	return Result{
		Name:   "Clean working tree",
		Passed: true,
	}
}

// checkBranchUpToDate verifies target branch exists
func (c *Checker) checkBranchUpToDate(ctx context.Context) Result {
	exists, err := c.git.BranchExists(ctx, c.targetBranch)
	if err != nil {
		return Result{
			Name:   fmt.Sprintf("Target branch '%s' exists", c.targetBranch),
			Passed: false,
			Error:  err.Error(),
			Hint:   "Failed to check if branch exists",
		}
	}

	if !exists {
		return Result{
			Name:   fmt.Sprintf("Target branch '%s' exists", c.targetBranch),
			Passed: false,
			Hint:   fmt.Sprintf("Create branch '%s' first or check your configuration", c.targetBranch),
		}
	}

	return Result{
		Name:   fmt.Sprintf("Target branch '%s' exists", c.targetBranch),
		Passed: true,
	}
}
