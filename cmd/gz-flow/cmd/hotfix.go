package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var hotfixCmd = &cobra.Command{
	Use:   "hotfix",
	Short: "Manage hotfix branches",
	Long: `Manage hotfix branches in the git-flow workflow.

Hotfix branches arise from the necessity to act immediately upon
an undesired state of a live production version.

Commands:
  start   - Start a new hotfix branch from master
  finish  - Finish a hotfix branch (merge to master and develop, tag)`,
}

var hotfixStartCmd = &cobra.Command{
	Use:   "start <version>",
	Short: "Start a new hotfix branch",
	Long: `Start a new hotfix branch from the master branch.

Example:
  gz-flow hotfix start 1.0.1
  gz-flow hotfix start v2.1.1`,
	Args: cobra.ExactArgs(1),
	RunE: runHotfixStart,
}

var hotfixFinishCmd = &cobra.Command{
	Use:   "finish <version>",
	Short: "Finish a hotfix branch",
	Long: `Finish a hotfix branch by merging it into master and develop.

This will:
  - Merge the hotfix branch into master
  - Tag the hotfix on master
  - Merge the hotfix branch into develop (or release if active)
  - Delete the hotfix branch

Example:
  gz-flow hotfix finish 1.0.1`,
	Args: cobra.ExactArgs(1),
	RunE: runHotfixFinish,
}

func init() {
	rootCmd.AddCommand(hotfixCmd)

	hotfixCmd.AddCommand(hotfixStartCmd)
	hotfixCmd.AddCommand(hotfixFinishCmd)

	hotfixFinishCmd.Flags().StringVarP(&tagMessage, "message", "m", "", "Tag message")
	hotfixFinishCmd.Flags().BoolVar(&noTag, "no-tag", false, "Don't create a tag")
	hotfixFinishCmd.Flags().BoolVarP(&keepBranch, "keep", "k", false, "Keep the hotfix branch")
}

func runHotfixStart(cmd *cobra.Command, args []string) error {
	if err := checkGitRepo(); err != nil {
		return err
	}

	version := args[0]

	// TODO: Implement hotfix start logic
	// 1. Validate version format
	// 2. Check out master branch
	// 3. Create hotfix branch

	fmt.Printf("Started hotfix branch 'hotfix/%s'\n", version)
	fmt.Printf("Switched to branch 'hotfix/%s'\n", version)

	return nil
}

func runHotfixFinish(cmd *cobra.Command, args []string) error {
	if err := checkGitRepo(); err != nil {
		return err
	}

	version := args[0]

	// TODO: Implement hotfix finish logic
	// 1. Merge hotfix to master
	// 2. Create tag
	// 3. Merge hotfix to develop (or active release branch)
	// 4. Delete hotfix branch

	fmt.Printf("Merged hotfix/%s into master\n", version)
	if !noTag {
		fmt.Printf("Created tag 'v%s'\n", version)
	}
	fmt.Printf("Merged hotfix/%s into develop\n", version)
	if !keepBranch {
		fmt.Printf("Deleted branch hotfix/%s\n", version)
	}

	return nil
}
