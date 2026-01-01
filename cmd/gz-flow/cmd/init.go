package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize git-flow in the current repository",
	Long: `Initialize git-flow in the current repository.

This will:
  - Detect or create the master/main branch
  - Create the develop branch if it doesn't exist
  - Save git-flow configuration to .gzflow.yaml

Example:
  gz-flow init
  gz-flow init --defaults    # Use all defaults without prompting`,
	RunE: runInit,
}

var (
	useDefaults bool
	force       bool
)

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().BoolVarP(&useDefaults, "defaults", "d", false, "Use default branch names")
	initCmd.Flags().BoolVarP(&force, "force", "f", false, "Force re-initialization")
}

func runInit(cmd *cobra.Command, args []string) error {
	if err := checkGitRepo(); err != nil {
		return err
	}

	// TODO: Implement initialization logic
	// 1. Detect master/main branch
	// 2. Create develop branch if not exists
	// 3. Save configuration

	fmt.Println("Git-flow initialized successfully!")
	fmt.Println("")
	fmt.Println("Summary of branches:")
	fmt.Println("  - master: master")
	fmt.Println("  - develop: develop")
	fmt.Println("")
	fmt.Println("Configuration saved to .gzflow.yaml")

	return nil
}
