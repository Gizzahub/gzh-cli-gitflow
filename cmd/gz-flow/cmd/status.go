package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current git-flow status",
	Long: `Show the current git-flow workflow status.

Displays:
  - Current branch and its type
  - Active flow branches
  - Working directory status

Example:
  gz-flow status`,
	RunE: runStatus,
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

func runStatus(cmd *cobra.Command, args []string) error {
	if err := checkGitRepo(); err != nil {
		return err
	}

	// TODO: Implement status logic
	// 1. Get current branch
	// 2. Determine branch type
	// 3. List active flow branches
	// 4. Show working directory status

	fmt.Println("Git-flow Status")
	fmt.Println("===============")
	fmt.Println("")
	fmt.Println("Current branch: develop")
	fmt.Println("Branch type: develop")
	fmt.Println("")
	fmt.Println("Active branches:")
	fmt.Println("  feature: (none)")
	fmt.Println("  release: (none)")
	fmt.Println("  hotfix:  (none)")
	fmt.Println("")
	fmt.Println("Working directory: clean")

	return nil
}
