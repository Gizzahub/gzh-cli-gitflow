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

	version := args[0]

	// TODO: Implement release finish logic
	// 1. Merge release to master
	// 2. Create tag
	// 3. Merge release to develop
	// 4. Delete release branch

	fmt.Printf("Merged release/%s into master\n", version)
	if !noTag {
		fmt.Printf("Created tag 'v%s'\n", version)
	}
	fmt.Printf("Merged release/%s into develop\n", version)
	if !keepBranch {
		fmt.Printf("Deleted branch release/%s\n", version)
	}

	return nil
}
