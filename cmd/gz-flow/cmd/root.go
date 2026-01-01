package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	version string
	cfgFile string
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "gz-flow",
	Short: "Git-flow workflow automation CLI",
	Long: `gz-flow is a Git-flow workflow automation CLI tool.

It provides commands for managing git-flow branches:
  - feature: Feature branch management
  - release: Release branch management
  - hotfix:  Hotfix branch management

Example:
  gz-flow init                    # Initialize git-flow
  gz-flow feature start my-feat   # Start a feature branch
  gz-flow feature finish my-feat  # Finish the feature`,
}

// Execute adds all child commands to the root command
func Execute() error {
	return rootCmd.Execute()
}

// SetVersion sets the version for the version command
func SetVersion(v string) {
	version = v
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default: ~/.gz/gitflow)")

	// Add version command
	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print the version number",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("gz-flow version %s\n", version)
		},
	})
}

// checkGitRepo checks if current directory is a git repository
func checkGitRepo() error {
	if _, err := os.Stat(".git"); os.IsNotExist(err) {
		return fmt.Errorf("not a git repository (or any of the parent directories)")
	}
	return nil
}
