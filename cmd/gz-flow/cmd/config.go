package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config [key] [value]",
	Short: "Manage git-flow configuration",
	Long: `Manage git-flow configuration.

Without arguments, shows current configuration.
With key, shows the value of that key.
With key and value, sets the configuration.

Example:
  gz-flow config                      # Show all config
  gz-flow config branches.master      # Get master branch name
  gz-flow config branches.master main # Set master branch to 'main'
  gz-flow config --global ...         # Modify global config`,
	Args: cobra.MaximumNArgs(2),
	RunE: runConfig,
}

var globalConfig bool

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.Flags().BoolVarP(&globalConfig, "global", "g", false, "Use global configuration")
}

func runConfig(cmd *cobra.Command, args []string) error {
	// TODO: Implement config logic
	// 1. Load configuration
	// 2. Handle get/set operations
	// 3. Save if modified

	if len(args) == 0 {
		// Show all config
		tagFormat := "v%s" //nolint:govet // literal string, not format
		fmt.Println("Git-flow Configuration")
		fmt.Println("======================")
		fmt.Println("")
		fmt.Println("Branches:")
		fmt.Println("  master:  master")
		fmt.Println("  develop: develop")
		fmt.Println("")
		fmt.Println("Prefixes:")
		fmt.Println("  feature: feature/")
		fmt.Println("  release: release/")
		fmt.Println("  hotfix:  hotfix/")
		fmt.Println("")
		fmt.Println("Options:")
		fmt.Println("  delete_branch_after_finish: true")
		fmt.Println("  push_after_finish: false")
		fmt.Println("  tag_format:", tagFormat)
		return nil
	}

	key := args[0]

	if len(args) == 1 {
		// Get value
		fmt.Printf("%s = (value)\n", key)
		return nil
	}

	value := args[1]
	configScope := "local"
	if globalConfig {
		configScope = "global"
	}

	fmt.Printf("Set %s = %s (%s)\n", key, value, configScope)
	return nil
}
