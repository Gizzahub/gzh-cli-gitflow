package cmd

import (
	"fmt"

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
	Use:   "start <name>",
	Short: "Start a new feature branch",
	Long: `Start a new feature branch from the develop branch.

Example:
  gz-flow feature start user-authentication
  gz-flow feature start login-page`,
	Args: cobra.ExactArgs(1),
	RunE: runFeatureStart,
}

var featureFinishCmd = &cobra.Command{
	Use:   "finish <name>",
	Short: "Finish a feature branch",
	Long: `Finish a feature branch by merging it into develop.

This will:
  - Merge the feature branch into develop
  - Delete the feature branch (unless --keep is specified)

Example:
  gz-flow feature finish user-authentication`,
	Args: cobra.ExactArgs(1),
	RunE: runFeatureFinish,
}

var (
	keepBranch bool
	squash     bool
)

func init() {
	rootCmd.AddCommand(featureCmd)

	featureCmd.AddCommand(featureStartCmd)
	featureCmd.AddCommand(featureFinishCmd)

	featureFinishCmd.Flags().BoolVarP(&keepBranch, "keep", "k", false, "Keep the feature branch after finishing")
	featureFinishCmd.Flags().BoolVar(&squash, "squash", false, "Squash commits when merging")
}

func runFeatureStart(cmd *cobra.Command, args []string) error {
	if err := checkGitRepo(); err != nil {
		return err
	}

	name := args[0]

	// TODO: Implement feature start logic
	// 1. Validate branch name
	// 2. Check out develop branch
	// 3. Create feature branch

	fmt.Printf("Started feature branch 'feature/%s'\n", name)
	fmt.Printf("Switched to branch 'feature/%s'\n", name)

	return nil
}

func runFeatureFinish(cmd *cobra.Command, args []string) error {
	if err := checkGitRepo(); err != nil {
		return err
	}

	name := args[0]

	// TODO: Implement feature finish logic
	// 1. Check for uncommitted changes
	// 2. Merge feature branch to develop
	// 3. Delete feature branch (unless --keep)

	fmt.Printf("Merged feature/%s into develop\n", name)
	if !keepBranch {
		fmt.Printf("Deleted branch feature/%s\n", name)
	}

	return nil
}
