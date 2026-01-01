// Package gitcmd provides safe git command execution.
package gitcmd

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

// Executor executes git commands safely.
type Executor struct {
	workDir string
}

// New creates a new git command executor.
func New() *Executor {
	return &Executor{}
}

// WithWorkDir sets the working directory
func (e *Executor) WithWorkDir(dir string) *Executor {
	return &Executor{workDir: dir}
}

// validateBranchName performs basic validation on branch names
// to prevent malformed git commands
func validateBranchName(name string) error {
	if name == "" {
		return fmt.Errorf("branch name cannot be empty")
	}
	if strings.HasPrefix(name, "-") {
		return fmt.Errorf("branch name cannot start with '-'")
	}
	// Allow alphanumeric, slash, underscore, hyphen
	if !regexp.MustCompile(`^[a-zA-Z0-9/_-]+$`).MatchString(name) {
		return fmt.Errorf("branch name contains invalid characters")
	}
	return nil
}

// run executes a git command with the given arguments.
// This is the ONLY place where exec.Command should be called.
func (e *Executor) run(ctx context.Context, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, "git", args...)
	if e.workDir != "" {
		cmd.Dir = e.workDir
	}
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("git %s: %w: %s", strings.Join(args, " "), err, stderr.String())
	}

	return strings.TrimSpace(stdout.String()), nil
}

// CurrentBranch returns the current branch name.
func (e *Executor) CurrentBranch(ctx context.Context) (string, error) {
	return e.run(ctx, "branch", "--show-current")
}

// IsClean returns true if the working directory is clean.
func (e *Executor) IsClean(ctx context.Context) (bool, error) {
	out, err := e.run(ctx, "status", "--porcelain")
	if err != nil {
		return false, err
	}
	return out == "", nil
}

// BranchExists checks if a branch exists
func (e *Executor) BranchExists(ctx context.Context, name string) (bool, error) {
	if err := validateBranchName(name); err != nil {
		return false, fmt.Errorf("invalid branch name: %w", err)
	}
	_, err := e.run(ctx, "rev-parse", "--verify", name)
	if err != nil {
		// Check if it's "not found" error vs other errors
		errStr := err.Error()
		if strings.Contains(errStr, "unknown revision") ||
			strings.Contains(errStr, "not a valid ref") ||
			strings.Contains(errStr, "Needed a single revision") {
			return false, nil // branch doesn't exist - expected
		}
		return false, err // other error - unexpected
	}
	return true, nil
}

// Checkout switches to the specified branch.
func (e *Executor) Checkout(ctx context.Context, branch string) error {
	if err := validateBranchName(branch); err != nil {
		return fmt.Errorf("invalid branch name: %w", err)
	}
	_, err := e.run(ctx, "checkout", branch)
	return err
}

// CreateBranch creates a new branch from the current HEAD.
func (e *Executor) CreateBranch(ctx context.Context, branch string) error {
	if err := validateBranchName(branch); err != nil {
		return fmt.Errorf("invalid branch name: %w", err)
	}
	_, err := e.run(ctx, "checkout", "-b", branch)
	return err
}

// Merge merges the specified branch into the current branch.
func (e *Executor) Merge(ctx context.Context, branch string, noFF bool) error {
	if err := validateBranchName(branch); err != nil {
		return fmt.Errorf("invalid branch name: %w", err)
	}
	args := []string{"merge"}
	if noFF {
		args = append(args, "--no-ff")
	}
	args = append(args, branch)
	_, err := e.run(ctx, args...)
	return err
}

// DeleteBranch deletes the specified branch.
func (e *Executor) DeleteBranch(ctx context.Context, name string) error {
	if err := validateBranchName(name); err != nil {
		return fmt.Errorf("invalid branch name: %w", err)
	}
	_, err := e.run(ctx, "branch", "-d", name)
	return err
}

// ListBranches returns all branches matching the prefix.
func (e *Executor) ListBranches(ctx context.Context, prefix string) ([]string, error) {
	out, err := e.run(ctx, "branch", "--list", prefix+"*")
	if err != nil {
		return nil, err
	}

	if out == "" {
		return []string{}, nil
	}

	lines := strings.Split(out, "\n")
	branches := make([]string, 0, len(lines))
	for _, line := range lines {
		// Remove leading "* " or "  " from branch names
		branch := strings.TrimSpace(strings.TrimPrefix(line, "*"))
		if branch != "" {
			branches = append(branches, branch)
		}
	}
	return branches, nil
}
