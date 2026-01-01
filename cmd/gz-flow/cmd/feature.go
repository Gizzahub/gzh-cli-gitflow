package cmd

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gizzahub/gzh-cli-gitflow/internal/gitcmd"
	"github.com/gizzahub/gzh-cli-gitflow/internal/preflight"
	"github.com/gizzahub/gzh-cli-gitflow/internal/validator"
	"github.com/gizzahub/gzh-cli-gitflow/pkg/config"
	"github.com/spf13/cobra"
)

var featureCmd = &cobra.Command{
	Use:   "feature",
	Short: "Manage feature branches",
	Long: `Manage feature branches in the git-flow workflow.

Feature branches are used to develop new features for the upcoming
or a distant future release.

Commands:
  start   - Start a new feature branch from develop
  finish  - Finish a feature branch (merge to develop)`,
}

var featureStartCmd = &cobra.Command{
	Use:   "start [name]",
	Short: "Start a new feature branch",
	Long: `Start a new feature branch from the develop branch.

Example:
  gz-flow feature start user-authentication
  gz-flow feature start login-page
  gz-flow feature start api-refactor --from main`,
	Args: cobra.MaximumNArgs(1),
	RunE: runFeatureStart,
}

var featureFinishCmd = &cobra.Command{
	Use:   "finish [name]",
	Short: "Finish a feature branch",
	Long: `Finish a feature branch by merging it into develop.

This will:
  - Run pre-flight checks (clean tree, up-to-date branch)
  - Merge the feature branch into develop
  - Delete the feature branch (unless --keep is specified)

Example:
  gz-flow feature finish user-authentication
  gz-flow feature finish --auto  # Auto-detect from current branch`,
	Args: cobra.MaximumNArgs(1),
	RunE: runFeatureFinish,
}

var (
	keepBranch bool
	fromBranch string
)

func init() {
	rootCmd.AddCommand(featureCmd)

	featureCmd.AddCommand(featureStartCmd)
	featureCmd.AddCommand(featureFinishCmd)

	featureStartCmd.Flags().StringVar(&fromBranch, "from", "", "Base branch to start from (default: develop)")

	featureFinishCmd.Flags().BoolVarP(&keepBranch, "keep", "k", false, "Keep the feature branch after finishing")
}

func runFeatureStart(cmd *cobra.Command, args []string) error {
	// 1. Check git repo
	if err := checkGitRepo(); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	git := gitcmd.New()
	cfg, err := config.LoadFromDir(".")
	if err != nil {
		fmt.Printf("⚠️  Failed to load config, using defaults: %v\n", err)
		cfg = config.Default()
	}

	// 2. Get branch name
	if len(args) == 0 {
		return fmt.Errorf("feature name is required\nUsage: gz-flow feature start <name>")
	}
	name := args[0]

	// 3. Validate branch name
	if err := validator.ValidateBranchName(name); err != nil {
		suggested := validator.SuggestBranchName(name)
		return fmt.Errorf("invalid branch name: %v\n💡 Suggested: %s", err, suggested)
	}

	// 4. Check Guardian rules if enabled
	if cfg.Guardian.Enabled {
		if err := cfg.Guardian.Naming.Validate(name); err != nil {
			return fmt.Errorf("guardian: %v", err)
		}
	}

	// 5. Determine base branch
	baseBranch := cfg.Branches.Develop
	if fromBranch != "" {
		baseBranch = fromBranch
	}

	// 6. Context hint: warn if not on expected branch
	currentBranch, _ := git.CurrentBranch(ctx)
	if currentBranch != baseBranch {
		fmt.Printf("⚠️  You're on '%s', not '%s'\n", currentBranch, baseBranch)
		fmt.Printf("💡 Will checkout '%s' first\n\n", baseBranch)
	}

	// 7. Check if branch already exists
	fullBranchName := cfg.Prefixes.Feature + name
	exists, _ := git.BranchExists(ctx, fullBranchName)
	if exists {
		return fmt.Errorf("branch '%s' already exists", fullBranchName)
	}

	// 8. Execute
	if err := git.Checkout(ctx, baseBranch); err != nil {
		return fmt.Errorf("failed to checkout %s: %v", baseBranch, err)
	}

	if err := git.CreateBranch(ctx, fullBranchName); err != nil {
		return fmt.Errorf("failed to create branch: %v", err)
	}

	fmt.Printf("✅ Started feature branch '%s'\n", fullBranchName)
	fmt.Printf("📍 Switched to branch '%s'\n", fullBranchName)

	return nil
}

func runFeatureFinish(cmd *cobra.Command, args []string) error {
	// 1. Check git repo
	if err := checkGitRepo(); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	git := gitcmd.New()
	cfg, err := config.LoadFromDir(".")
	if err != nil {
		fmt.Printf("⚠️  Failed to load config, using defaults: %v\n", err)
		cfg = config.Default()
	}

	// 2. Determine feature name
	var name string
	if len(args) > 0 {
		name = args[0]
	} else {
		// Auto-detect from current branch
		currentBranch, err := git.CurrentBranch(ctx)
		if err != nil {
			return fmt.Errorf("failed to get current branch: %v", err)
		}

		prefix := cfg.Prefixes.Feature
		if !strings.HasPrefix(currentBranch, prefix) {
			return fmt.Errorf("not on a feature branch (current: %s)\n💡 Use 'gz-flow feature finish <name>' or switch to a feature branch", currentBranch)
		}
		name = strings.TrimPrefix(currentBranch, prefix)
		fmt.Printf("📍 Auto-detected feature: %s\n\n", name)
	}

	fullBranchName := cfg.Prefixes.Feature + name
	targetBranch := cfg.Branches.Develop

	// 3. Pre-flight checks
	checker := preflight.NewChecker(git, targetBranch)
	results := checker.RunAll(ctx)

	fmt.Println("🔍 Pre-flight checks:")
	fmt.Print(results.String())

	if results.HasErrors() {
		return fmt.Errorf("pre-flight checks failed")
	}
	fmt.Println()

	// 4. Check source branch exists
	exists, _ := git.BranchExists(ctx, fullBranchName)
	if !exists {
		return fmt.Errorf("feature branch '%s' does not exist", fullBranchName)
	}

	// 5. Execute merge
	if err := git.Checkout(ctx, targetBranch); err != nil {
		return fmt.Errorf("failed to checkout %s: %v", targetBranch, err)
	}

	if err := git.Merge(ctx, fullBranchName, true); err != nil {
		return fmt.Errorf("merge failed: %v\n💡 Resolve conflicts and run 'git merge --continue'", err)
	}

	fmt.Printf("✅ Merged '%s' into '%s'\n", fullBranchName, targetBranch)

	// 6. Delete branch if requested
	deleteBranch := cfg.Options.DeleteBranchAfterFinish && !keepBranch
	if deleteBranch {
		if err := git.DeleteBranch(ctx, fullBranchName); err != nil {
			fmt.Printf("⚠️  Failed to delete branch: %v\n", err)
		} else {
			fmt.Printf("🗑️  Deleted branch '%s'\n", fullBranchName)
		}
	}

	return nil
}
