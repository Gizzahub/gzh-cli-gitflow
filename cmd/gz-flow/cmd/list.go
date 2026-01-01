package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list [type]",
	Short: "List git-flow branches",
	Long: `List active git-flow branches.

If type is specified, only list branches of that type.
Valid types: feature, release, hotfix

Example:
  gz-flow list           # List all flow branches
  gz-flow list feature   # List only feature branches`,
	Args: cobra.MaximumNArgs(1),
	RunE: runList,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func runList(cmd *cobra.Command, args []string) error {
	if err := checkGitRepo(); err != nil {
		return err
	}

	var branchType string
	if len(args) > 0 {
		branchType = args[0]
	}

	// TODO: Implement list logic
	// 1. Get all branches
	// 2. Filter by type if specified
	// 3. Display with metadata

	if branchType != "" {
		fmt.Printf("Active %s branches:\n", branchType)
	} else {
		fmt.Println("Active git-flow branches:")
	}

	fmt.Println("")
	fmt.Println("Feature branches:")
	fmt.Println("  (none)")
	fmt.Println("")
	fmt.Println("Release branches:")
	fmt.Println("  (none)")
	fmt.Println("")
	fmt.Println("Hotfix branches:")
	fmt.Println("  (none)")

	return nil
}
