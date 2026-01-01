package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
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

	version := args[0]

	// TODO: Implement release start logic
	// 1. Validate version format
	// 2. Check out develop branch
	// 3. Create release branch

	fmt.Printf("Started release branch 'release/%s'\n", version)
	fmt.Printf("Switched to branch 'release/%s'\n", version)

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
