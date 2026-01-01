package cmd

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

var releaseCmd = &cobra.Command{
	Use:   "release",
	Short: "Manage release branches",
	Long: `Manage release branches in the git-flow workflow.

Release branches support preparation of a new production release.
They allow for last-minute dotting of i's and crossing t's.

Commands:
  start   - Start a new release branch from develop
  finish  - Finish a release branch (merge to master and develop, tag)`,
}

var releaseStartCmd = &cobra.Command{
	Use:   "start <version>",
	Short: "Start a new release branch",
	Long: `Start a new release branch from the develop branch.

Example:
  gz-flow release start 1.0.0
  gz-flow release start v2.1.0`,
	Args: cobra.ExactArgs(1),
	RunE: runReleaseStart,
}

var releaseFinishCmd = &cobra.Command{
	Use:   "finish <version>",
	Short: "Finish a release branch",
	Long: `Finish a release branch by merging it into master and develop.

This will:
  - Merge the release branch into master
  - Tag the release on master
  - Merge the release branch into develop
  - Delete the release branch

Example:
  gz-flow release finish 1.0.0`,
	Args: cobra.ExactArgs(1),
	RunE: runReleaseFinish,
}

var (
	tagMessage string
	noTag      bool
)

func init() {
	rootCmd.AddCommand(releaseCmd)

	releaseCmd.AddCommand(releaseStartCmd)
	releaseCmd.AddCommand(releaseFinishCmd)

	releaseFinishCmd.Flags().StringVarP(&tagMessage, "message", "m", "", "Tag message")
	releaseFinishCmd.Flags().BoolVar(&noTag, "no-tag", false, "Don't create a tag")
	releaseFinishCmd.Flags().BoolVarP(&keepBranch, "keep", "k", false, "Keep the release branch")
}

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
